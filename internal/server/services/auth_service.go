package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"yamony/internal/database/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

func (s *service) RegisterUser(ctx context.Context, username, email, password string) (*sqlc.CreateUserRow, string, error) {
	_, err := s.db.GetQueries().GetUserByEmail(ctx, email)
	if err == nil {
		return nil, "", ErrEmailAlreadyExists
	}
	if err != pgx.ErrNoRows {
		return nil, "", fmt.Errorf("failed to check existing user: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
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
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate session token: %w", err)
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
		return nil, "", fmt.Errorf("failed to create session: %w", err)
	}

	return &user, sessionToken, nil
}

func (s *service) LoginUser(ctx context.Context, email, password string) (*sqlc.GetUserByEmailRow, string, error) {
	user, err := s.db.GetQueries().GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate session token: %w", err)
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
		return nil, "", fmt.Errorf("failed to create session: %w", err)
	}

	return &user, sessionToken, nil
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

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
