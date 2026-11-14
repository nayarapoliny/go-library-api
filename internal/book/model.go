package book

import (
	"context"
	"errors"
	"time"
)

var (
	ErrBookNotFound = errors.New("livro não encontrado")
	ErrInvalidInput = errors.New("dados do livro inválidos: título, autor e género são obrigatórios")
)

type Book struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Genre         string    `json:"genre"`
	PublishedYear int       `json:"published_year"`
	Publisher     string    `json:"publisher"`
	CreatedAt     time.Time `json:"created_at"`
}

type SearchFilter struct {
	Title         string
	Author        string
	Genre         string
	Publisher     string
	PublishedYear int
}

type Repository interface {
	Create(ctx context.Context, b *Book) error
	GetByID(ctx context.Context, id string) (*Book, error)
	Search(ctx context.Context, filter SearchFilter) ([]Book, error)
	Delete(ctx context.Context, id string) error
}
