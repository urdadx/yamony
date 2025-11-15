package services

import (
	"context"
	"yamony/internal/database"
	"yamony/internal/database/sqlc"
)

type Service interface {
	RegisterUser(ctx context.Context, username, email, password string) (*sqlc.CreateUserRow, string, error)
	LoginUser(ctx context.Context, email, password string) (*sqlc.GetUserByEmailRow, string, error)
	ValidateSession(ctx context.Context, sessionToken string) (*sqlc.GetUserByIDRow, error)
	LogoutUser(ctx context.Context, sessionToken string) error
	CreatePage(ctx context.Context, userID int32, handle string) (*sqlc.Page, error)
	UpdatePage(ctx context.Context, activePageID int32, isActive bool, name, handle, image, bannerImage, bio string) (*sqlc.Page, error)
	CheckHandleExists(ctx context.Context, handle string) (bool, error)
	GetPageByID(ctx context.Context, id int32) (*sqlc.Page, error)
	DeletePage(ctx context.Context, id int32, userID int32) error
	SetActivePage(ctx context.Context, pageID int32) error
	GetPageByHandle(ctx context.Context, handle string) (*sqlc.Page, error)
	SetNextPageAsActive(ctx context.Context, user_id int32) (*sqlc.Page, error)
}

type service struct {
	db database.Service
}

func New(db database.Service) Service {
	return &service{
		db: db,
	}
}
