package book

import (
	"context"
	"testing"
)

func TestService_CreateAndSearchBooks(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	newBook := &Book{
		Title:         "Designing Data-Intensive Applications",
		Author:        "Martin Kleppmann",
		Genre:         "Distributed Systems",
		PublishedYear: 2017,
		Publisher:     "O'Reilly",
	}

	if err := svc.CreateBook(ctx, newBook); err != nil {
		t.Fatalf("erro inesperado ao criar livro: %v", err)
	}

	if newBook.ID == "" {
		t.Fatal("esperava um ID gerado automaticamente")
	}

	results, err := svc.SearchBooks(ctx, SearchFilter{Author: "Kleppmann", PublishedYear: 2017})
	if err != nil {
		t.Fatalf("erro na pesquisa: %v", err)
	}

	if len(results) != 1 || results[0].Title != newBook.Title {
		t.Errorf("esperava encontrar 1 livro correspondente, obteve %d", len(results))
	}
}
