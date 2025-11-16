package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"yamony/internal/database/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

func (s *service) RegisterUser(ctx context.Context, username, email, password string) (*sqlc.CreateUserRow, string, int32, error) {
	_, err := s.db.GetQueries().GetUserByEmail(ctx, email)
	if err == nil {
		return nil, "", 0, ErrEmailAlreadyExists
	}
	if err != pgx.ErrNoRows {
		return nil, "", 0, fmt.Errorf("failed to check existing user: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to hash password: %w", err)
	}

	params := sqlc.CreateUserParams{
		Username:      username,
		Email:         email,
		PasswordHash:  string(hashedPassword),
		EmailVerified: false,
		Image:         "",
	}

	user, err := s.db.GetQueries().CreateUser(ctx, params)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to create user: %w", err)
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to generate session token: %w", err)
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	var expiresAtPg pgtype.Timestamp
	expiresAtPg.Time = expiresAt
	expiresAtPg.Valid = true

	sessionParams := sqlc.CreateSessionParams{
		UserID:       user.ID,
		SessionToken: sessionToken,
		ExpiresAt:    expiresAtPg,
	}

	_, err = s.db.GetQueries().CreateSession(ctx, sessionParams)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to create session: %w", err)
	}

	return &user, sessionToken, 0, nil
}

func (s *service) LoginUser(ctx context.Context, email, password string) (*sqlc.GetUserByEmailRow, string, int32, error) {
	user, err := s.db.GetQueries().GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, "", 0, ErrInvalidCredentials
		}
		return nil, "", 0, fmt.Errorf("failed to get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", 0, ErrInvalidCredentials
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to generate session token: %w", err)
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	var expiresAtPg pgtype.Timestamp
	expiresAtPg.Time = expiresAt
	expiresAtPg.Valid = true

	var activePageID int32
	recentPage, err := s.db.GetQueries().GetUserMostRecentPage(ctx, user.ID)
	if err == nil {
		activePageID = recentPage.ID
	}

	sessionParams := sqlc.CreateSessionParams{
		UserID:       user.ID,
		SessionToken: sessionToken,
		ExpiresAt:    expiresAtPg,
	}

	session, err := s.db.GetQueries().CreateSession(ctx, sessionParams)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to create session: %w", err)
	}

	if activePageID > 0 {
		var activePageIDPg pgtype.Int4
		activePageIDPg.Int32 = activePageID
		activePageIDPg.Valid = true

		_, err = s.db.GetQueries().UpdateSessionWithActivePage(ctx, sqlc.UpdateSessionWithActivePageParams{
			ID:           session.ID,
			ActivePageID: activePageIDPg,
		})
		if err != nil {
			fmt.Printf("Warning: failed to update session with active page: %v\n", err)
		}
	}

	return &user, sessionToken, activePageID, nil
}

func (s *service) ValidateSession(ctx context.Context, sessionToken string) (*sqlc.GetUserByIDRow, error) {
	session, err := s.db.GetQueries().GetSessionByToken(ctx, sessionToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session.ExpiresAt.Valid && time.Now().After(session.ExpiresAt.Time) {
		_ = s.db.GetQueries().DeleteSession(ctx, session.ID)
		return nil, errors.New("session expired")
	}

	user, err := s.db.GetQueries().GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *service) LogoutUser(ctx context.Context, sessionToken string) error {
	session, err := s.db.GetQueries().GetSessionByToken(ctx, sessionToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return fmt.Errorf("failed to get session: %w", err)
	}

	err = s.db.GetQueries().DeleteSession(ctx, session.ID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *service) GetSessionByUserID(ctx context.Context, userID int32) ([]sqlc.Session, error) {
	sessions, err := s.db.GetQueries().GetSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions by user ID: %w", err)
	}

	return sessions, nil
}

func (s *service) SyncActivePageToSession(ctx context.Context, sessionToken string, pageID int32) error {
	session, err := s.db.GetQueries().GetSessionByToken(ctx, sessionToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("session not found")
		}
		return fmt.Errorf("failed to get session: %w", err)
	}

	var activePageIDPg pgtype.Int4
	activePageIDPg.Int32 = pageID
	activePageIDPg.Valid = true

	_, err = s.db.GetQueries().UpdateSessionWithActivePage(ctx, sqlc.UpdateSessionWithActivePageParams{
		ID:           session.ID,
		ActivePageID: activePageIDPg,
	})
	if err != nil {
		return fmt.Errorf("failed to update session with active page: %w", err)
	}

	return nil
}

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *service) GetGoogleOAuthConfig() *oauth2.Config {
	return s.googleOAuthConfig
}

func (s *service) GoogleOAuthLogin(ctx context.Context, code string) (*sqlc.GetUserByEmailRow, string, int32, error) {
	// Exchange code for token
	token, err := s.googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to exchange code: %w", err)
	}

	client := s.googleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, "", 0, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Check if user already exists
	user, err := s.db.GetQueries().GetUserByEmail(ctx, googleUser.Email)
	if err != nil && err != pgx.ErrNoRows {
		return nil, "", 0, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Create user if doesn't exist
	if err == pgx.ErrNoRows {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(googleUser.ID), bcrypt.DefaultCost)
		if err != nil {
			return nil, "", 0, fmt.Errorf("failed to hash password: %w", err)
		}

		params := sqlc.CreateUserParams{
			Username:      googleUser.Name,
			Email:         googleUser.Email,
			PasswordHash:  string(hashedPassword),
			EmailVerified: googleUser.VerifiedEmail,
			Image:         googleUser.Picture,
		}

		newUser, err := s.db.GetQueries().CreateUser(ctx, params)
		if err != nil {
			return nil, "", 0, fmt.Errorf("failed to create user: %w", err)
		}

		sessionToken, err := generateSessionToken()
		if err != nil {
			return nil, "", 0, fmt.Errorf("failed to generate session token: %w", err)
		}

		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		var expiresAtPg pgtype.Timestamp
		expiresAtPg.Time = expiresAt
		expiresAtPg.Valid = true

		sessionParams := sqlc.CreateSessionParams{
			UserID:       newUser.ID,
			SessionToken: sessionToken,
			ExpiresAt:    expiresAtPg,
		}

		_, err = s.db.GetQueries().CreateSession(ctx, sessionParams)
		if err != nil {
			return nil, "", 0, fmt.Errorf("failed to create session: %w", err)
		}

		userRow := &sqlc.GetUserByEmailRow{
			ID:            newUser.ID,
			Username:      newUser.Username,
			Email:         newUser.Email,
			PasswordHash:  string(hashedPassword),
			EmailVerified: newUser.EmailVerified,
			Image:         newUser.Image,
			CreatedAt:     newUser.CreatedAt,
			UpdatedAt:     newUser.UpdatedAt,
		}

		return userRow, sessionToken, 0, nil
	}

	// User exists, create session
	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to generate session token: %w", err)
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	var expiresAtPg pgtype.Timestamp
	expiresAtPg.Time = expiresAt
	expiresAtPg.Valid = true

	var activePageID int32
	recentPage, err := s.db.GetQueries().GetUserMostRecentPage(ctx, user.ID)
	if err == nil {
		activePageID = recentPage.ID
	}

	sessionParams := sqlc.CreateSessionParams{
		UserID:       user.ID,
		SessionToken: sessionToken,
		ExpiresAt:    expiresAtPg,
	}

	session, err := s.db.GetQueries().CreateSession(ctx, sessionParams)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to create session: %w", err)
	}

	if activePageID > 0 {
		var activePageIDPg pgtype.Int4
		activePageIDPg.Int32 = activePageID
		activePageIDPg.Valid = true

		_, err = s.db.GetQueries().UpdateSessionWithActivePage(ctx, sqlc.UpdateSessionWithActivePageParams{
			ID:           session.ID,
			ActivePageID: activePageIDPg,
		})
		if err != nil {
			fmt.Printf("Warning: failed to update session with active page: %v\n", err)
		}
	}

	return &user, sessionToken, activePageID, nil
}
