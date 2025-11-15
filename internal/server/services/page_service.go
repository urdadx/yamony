package services

import (
	"context"
	"errors"
	"fmt"
	"yamony/internal/database/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrPageAlreadyExists   = errors.New("page already exists")
	ErrPageNotFound        = errors.New("page not found")
	ErrHandleAlreadyExists = errors.New("page handle already exists")
)

func (s *service) CreatePage(ctx context.Context, userID int32, handle string, is_active bool) (*sqlc.Page, error) {
	exists, err := s.db.GetQueries().CheckHandleExists(ctx, sqlc.CheckHandleExistsParams{
		Handle: handle,
		ID:     0,
	})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPageAlreadyExists
	}

	params := sqlc.CreatePageParams{
		UserID:   userID,
		Handle:   handle,
		IsActive: is_active,
	}

	page, err := s.db.GetQueries().CreatePage(ctx, params)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

func (s *service) CheckHandleExists(ctx context.Context, handle string) (bool, error) {
	exist, error := s.db.GetQueries().CheckHandleExists(ctx, sqlc.CheckHandleExistsParams{
		Handle: handle,
		ID:     0,
	})
	if error != nil {
		return false, error
	}
	return exist, nil
}

func (s *service) UpdatePage(ctx context.Context, activePageId int32, is_active bool, name, handle, image, banner_image, bio string) (*sqlc.Page, error) {
	_, error := s.db.GetQueries().GetPageByID(ctx, activePageId)
	if error != nil {
		return nil, error
	}

	params := sqlc.UpdatePageParams{
		ID:          activePageId,
		Name:        pgtype.Text{String: name, Valid: true},
		Handle:      pgtype.Text{String: handle, Valid: true},
		Bio:         pgtype.Text{String: bio, Valid: true},
		IsActive:    pgtype.Bool{Bool: is_active, Valid: true},
		Image:       pgtype.Text{String: image, Valid: true},
		BannerImage: pgtype.Text{String: banner_image, Valid: true},
	}

	updated_page, err := s.db.GetQueries().UpdatePage(ctx, params)

	return &updated_page, err

}

func (s *service) GetPageByHandle(ctx context.Context, handle string) (*sqlc.Page, error) {
	page, error := s.db.GetQueries().GetPageByHandle(ctx, handle)
	if error != nil {
		return nil, ErrPageNotFound
	}
	return &page, nil
}

func (s *service) GetPageByID(ctx context.Context, id int32) (*sqlc.Page, error) {
	page, error := s.db.GetQueries().GetPageByID(ctx, id)
	if error != nil {
		return nil, ErrPageNotFound
	}
	return &page, nil
}

func (s *service) DeletePage(ctx context.Context, id int32, userID int32) error {
	_, err := s.db.GetQueries().GetPageByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrPageNotFound
		}
		return fmt.Errorf("failed to get page: %w", err)
	}

	pages, err := s.db.GetQueries().GetPagesByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user pages: %w", err)
	}

	if len(pages) <= 1 {
		return errors.New("cannot delete the only page")
	}

	err = s.db.GetQueries().DeletePage(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete page: %w", err)
	}

	return nil
}

func (s *service) GetAllUserPage(ctx context.Context, user_id int32) ([]sqlc.Page, error) {
	result, err := s.db.GetQueries().GetPagesByUserID(ctx, user_id)

	if err != nil {
		return nil, err
	}

	return result, nil

}

func (s *service) SetActivePage(ctx context.Context, id int32) error {
	_, err := s.db.GetQueries().GetPageByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return err
		}
		return err
	}

	_, err = s.db.GetQueries().UpdatePage(ctx, sqlc.UpdatePageParams{
		ID:       id,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	})

	if err != nil {
		return err
	}

	return nil

}

func (s *service) SetNextPageAsActive(ctx context.Context, user_id int32) (*sqlc.Page, error) {
	pages, err := s.db.GetQueries().GetPagesByUserID(ctx, user_id)
	if err != nil {
		return nil, err
	}

	if len(pages) == 0 {
		return nil, ErrPageNotFound
	}

	// Get the first page
	nextPage := pages[0]

	page, err := s.db.GetQueries().UpdatePage(ctx, sqlc.UpdatePageParams{
		ID:       nextPage.ID,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	return &page, nil
}
