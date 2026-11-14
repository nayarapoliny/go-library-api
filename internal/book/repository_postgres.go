package book

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) (*PostgresRepository, error) {
	schema := `
	CREATE SEQUENCE IF NOT EXISTS book_seq START 1;
	CREATE TABLE IF NOT EXISTS books (
		id VARCHAR(32) PRIMARY KEY DEFAULT ('BOOK-' || LPAD(nextval('book_seq')::TEXT, 4, '0')),
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		genre VARCHAR(100) NOT NULL,
		published_year INT NOT NULL,
		publisher VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Create(ctx context.Context, b *Book) error {
	b.CreatedAt = time.Now().UTC()
	query := `
		INSERT INTO books (title, author, genre, published_year, publisher, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	return r.db.QueryRowContext(ctx, query, b.Title, b.Author, b.Genre, b.PublishedYear, b.Publisher, b.CreatedAt).Scan(&b.ID)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Book, error) {
	var b Book
	query := `SELECT id, title, author, genre, published_year, publisher, created_at FROM books WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&b.ID, &b.Title, &b.Author, &b.Genre, &b.PublishedYear, &b.Publisher, &b.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBookNotFound
	}
	return &b, err
}

func (r *PostgresRepository) Search(ctx context.Context, f SearchFilter) ([]Book, error) {
	query := `SELECT id, title, author, genre, published_year, publisher, created_at FROM books WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if f.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argID)
		args = append(args, "%"+f.Title+"%")
		argID++
	}
	if f.Author != "" {
		query += fmt.Sprintf(" AND author ILIKE $%d", argID)
		args = append(args, "%"+f.Author+"%")
		argID++
	}
	if f.Genre != "" {
		query += fmt.Sprintf(" AND genre ILIKE $%d", argID)
		args = append(args, f.Genre)
		argID++
	}
	if f.Publisher != "" {
		query += fmt.Sprintf(" AND publisher ILIKE $%d", argID)
		args = append(args, "%"+f.Publisher+"%")
		argID++
	}
	if f.PublishedYear > 0 {
		query += fmt.Sprintf(" AND published_year = $%d", argID)
		args = append(args, f.PublishedYear)
		argID++
	}

	query += " ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Book, 0)
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Genre, &b.PublishedYear, &b.Publisher, &b.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, b)
	}
	return results, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrBookNotFound
	}
	return nil
}
