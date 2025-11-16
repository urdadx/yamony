package services

import (
	"context"
	"os"
	"yamony/internal/database"
	"yamony/internal/database/sqlc"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Service interface {
	RegisterUser(ctx context.Context, username, email, password string) (*sqlc.CreateUserRow, string, int32, error)
	LoginUser(ctx context.Context, email, password string) (*sqlc.GetUserByEmailRow, string, int32, error)
	ValidateSession(ctx context.Context, sessionToken string) (*sqlc.GetUserByIDRow, error)
	LogoutUser(ctx context.Context, sessionToken string) error
	SyncActivePageToSession(ctx context.Context, sessionToken string, pageID int32) error
	CreatePage(ctx context.Context, userID int32, handle string, is_active bool) (*sqlc.Page, error)
	UpdatePage(ctx context.Context, activePageID int32, isActive bool, name, handle, image, bannerImage, bio string) (*sqlc.Page, error)
	CheckHandleExists(ctx context.Context, handle string) (bool, error)
	GetPageByID(ctx context.Context, id int32) (*sqlc.Page, error)
	DeletePage(ctx context.Context, id int32, userID int32) error
	SetActivePage(ctx context.Context, pageID int32) error
	GetPageByHandle(ctx context.Context, handle string) (*sqlc.Page, error)
	SetNextPageAsActive(ctx context.Context, user_id int32) (*sqlc.Page, error)
	GetSessionByUserID(ctx context.Context, userID int32) ([]sqlc.Session, error)
	GetAllUserPage(ctx context.Context, userID int32) ([]sqlc.Page, error)
	GetGoogleOAuthConfig() *oauth2.Config
	GoogleOAuthLogin(ctx context.Context, code string) (*sqlc.GetUserByEmailRow, string, int32, error)
}

type service struct {
	db                database.Service
	googleOAuthConfig *oauth2.Config
}

func New(db database.Service) Service {
	googleOAuthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:3000/api/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &service{
		db:                db,
		googleOAuthConfig: googleOAuthConfig,
	}
}
