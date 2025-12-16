package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadToGCS(c *gin.Context) {
	bucketName := os.Getenv("GCS_BUCKET")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GCS_BUCKET is not configured"})
		return
	}

	// Build credentials JSON if not provided directly
	credsJSON := os.Getenv("GOOGLE_CLOUD_CREDENTIALS_JSON")
	if credsJSON == "" {
		clientEmail := os.Getenv("GOOGLE_CLIENT_EMAIL")
		privateKey := os.Getenv("GOOGLE_PRIVATE_KEY")
		privateKeyID := os.Getenv("GOOGLE_PRIVATE_KEY_ID")
		projectID := os.Getenv("GOOGLE_PROJECT_ID")

		if clientEmail != "" && privateKey != "" {
			privateKey = normalizePrivateKey(privateKey)
			type sa struct {
				Type                    string `json:"type"`
				ProjectID               string `json:"project_id"`
				PrivateKeyID            string `json:"private_key_id,omitempty"`
				PrivateKey              string `json:"private_key"`
				ClientEmail             string `json:"client_email"`
				TokenURI                string `json:"token_uri"`
				AuthURI                 string `json:"auth_uri"`
				AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url,omitempty"`
				ClientX509CertURL       string `json:"client_x509_cert_url,omitempty"`
			}
			composed := sa{
				Type:         "service_account",
				ProjectID:    projectID,
				PrivateKeyID: privateKeyID,
				PrivateKey:   privateKey,
				ClientEmail:  clientEmail,
				// Default URIs expected by Google auth libraries
				TokenURI: "https://oauth2.googleapis.com/token",
				AuthURI:  "https://accounts.google.com/o/oauth2/auth",
			}
			if b, err := json.Marshal(composed); err == nil {
				credsJSON = string(b)
			}
		}
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required", "details": err.Error()})
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file", "details": err.Error()})
		return
	}
	defer src.Close()

	folderPath := c.PostForm("folderPath")
	if folderPath != "" {
		if folderPath[len(folderPath)-1] != '/' {
			folderPath = folderPath + "/"
		}
		folderPath = path.Clean(folderPath)
		if folderPath == "." {
			folderPath = ""
		}
	}

	ctx := context.Background()
	var client *storage.Client
	if credsJSON != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsJSON([]byte(credsJSON)))
	} else {
		client, err = storage.NewClient(ctx)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize storage client", "details": err.Error()})
		return
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	timestamp := time.Now().UnixMilli()
	safeRe := regexp.MustCompile(`[^\w\.-]+`)
	safeName := safeRe.ReplaceAllString(fileHeader.Filename, "_")
	objectName := fmt.Sprintf("%s%d_%s", folderPath, timestamp, safeName)

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	obj := bucket.Object(objectName)
	wc := obj.NewWriter(ctx)
	wc.ContentType = contentType

	if _, err := io.Copy(wc, src); err != nil {
		_ = wc.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload to storage", "details": err.Error()})
		return
	}
	if err := wc.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to finalize upload", "details": err.Error()})
		return
	}

	if credsJSON == "" {
		clientEmail := os.Getenv("GOOGLE_CLIENT_EMAIL")
		privateKey := normalizePrivateKey(os.Getenv("GOOGLE_PRIVATE_KEY"))
		if clientEmail == "" || privateKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "service account credentials are required to generate signed URLs"})
			return
		}
		// Create a minimal JSON for downstream parsing
		minimal := struct {
			ClientEmail string `json:"client_email"`
			PrivateKey  string `json:"private_key"`
		}{ClientEmail: clientEmail, PrivateKey: privateKey}
		if b, err := json.Marshal(minimal); err == nil {
			credsJSON = string(b)
		}
	}

	type saCreds struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
	}
	var cred saCreds
	if err := json.Unmarshal([]byte(credsJSON), &cred); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid service account JSON", "details": err.Error()})
		return
	}
	if cred.ClientEmail == "" || cred.PrivateKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service account JSON missing client_email/private_key"})
		return
	}
	cred.PrivateKey = normalizePrivateKey(cred.PrivateKey)

	expiryMinutes := 15
	if v := c.PostForm("expiryMinutes"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 24*60 {
			expiryMinutes = n
		}
	}
	expires := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)

	signedURL, err := storage.SignedURL(bucketName, objectName, &storage.SignedURLOptions{
		GoogleAccessID: cred.ClientEmail,
		PrivateKey:     []byte(cred.PrivateKey),
		Method:         "GET",
		Expires:        expires,
		ContentType:    contentType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate signed URL", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"signedUrl":   signedURL,
		"bucket":      bucketName,
		"objectName":  objectName,
		"contentType": contentType,
		"expiresAt":   expires.UTC().Format(time.RFC3339),
	})
}

// normalizePrivateKey ensures the private key contains actual newlines.
// If the value contains literal "\n" sequences, they are converted to newline characters.
func normalizePrivateKey(pk string) string {
	if pk == "" {
		return pk
	}

	if len(pk) >= 2 && pk[0] == '"' && pk[len(pk)-1] == '"' {
		pk = pk[1 : len(pk)-1]
	}
	pk = regexp.MustCompile(`\\r\\n|\\n`).ReplaceAllString(pk, "\n")
	return pk
}
