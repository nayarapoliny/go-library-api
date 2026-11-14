package book

import (
	"context"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateBook(ctx context.Context, b *Book) error {
	b.Title = strings.TrimSpace(b.Title)
	b.Author = strings.TrimSpace(b.Author)
	b.Genre = strings.TrimSpace(b.Genre)
	b.Publisher = strings.TrimSpace(b.Publisher)

	if b.Title == "" || b.Author == "" || b.Genre == "" {
		return ErrInvalidInput
	}
	return s.repo.Create(ctx, b)
}

func (s *Service) GetBook(ctx context.Context, id string) (*Book, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) SearchBooks(ctx context.Context, f SearchFilter) ([]Book, error) {
	return s.repo.Search(ctx, f)
}

func (s *Service) DeleteBook(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
