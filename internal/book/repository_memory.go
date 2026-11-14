package book

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type MemoryRepository struct {
	mu      sync.RWMutex
	books   map[string]Book
	counter uint64
}

func NewMemoryRepository() *MemoryRepository {
	repo := &MemoryRepository{
		books: make(map[string]Book),
	}
	// Dados iniciais para facilitar a validação da API
	seed := []Book{
		{Title: "The Go Programming Language", Author: "Alan Donovan", Genre: "Technology", PublishedYear: 2015, Publisher: "Addison-Wesley"},
		{Title: "Concurrency in Go", Author: "Katherine Cox-Buday", Genre: "Technology", PublishedYear: 2017, Publisher: "O'Reilly"},
		{Title: "Clean Architecture", Author: "Robert C. Martin", Genre: "Software Engineering", PublishedYear: 2017, Publisher: "Prentice Hall"},
	}
	for _, b := range seed {
		copyBook := b
		_ = repo.Create(context.Background(), &copyBook)
	}
	return repo
}

func (r *MemoryRepository) Create(_ context.Context, b *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	idNum := atomic.AddUint64(&r.counter, 1)
	b.ID = fmt.Sprintf("BOOK-%04d", idNum)
	b.CreatedAt = time.Now().UTC()
	r.books[b.ID] = *b
	return nil
}

func (r *MemoryRepository) GetByID(_ context.Context, id string) (*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	b, exists := r.books[id]
	if !exists {
		return nil, ErrBookNotFound
	}
	return &b, nil
}

func (r *MemoryRepository) Search(_ context.Context, f SearchFilter) ([]Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]Book, 0)
	for _, b := range r.books {
		if f.Title != "" && !strings.Contains(strings.ToLower(b.Title), strings.ToLower(f.Title)) {
			continue
		}
		if f.Author != "" && !strings.Contains(strings.ToLower(b.Author), strings.ToLower(f.Author)) {
			continue
		}
		if f.Genre != "" && !strings.EqualFold(b.Genre, f.Genre) {
			continue
		}
		if f.Publisher != "" && !strings.Contains(strings.ToLower(b.Publisher), strings.ToLower(f.Publisher)) {
			continue
		}
		if f.PublishedYear > 0 && b.PublishedYear != f.PublishedYear {
			continue
		}
		results = append(results, b)
	}
	return results, nil
}

func (r *MemoryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return ErrBookNotFound
	}
	delete(r.books, id)
	return nil
}
