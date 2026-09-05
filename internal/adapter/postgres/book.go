package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"golang-rest-api/internal/domain"
)

type BookRepository struct {
	pool *pgxpool.Pool
}

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{pool: pool}
}

func (s *BookRepository) Create(ctx context.Context, book domain.Book) (domain.Book, error) {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO books (id, title, author, isbn, year, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, book.ID, book.Title, book.Author, nullString(book.ISBN), nullInt(book.Year), book.CreatedAt, book.UpdatedAt)
	if err != nil {
		return domain.Book{}, mapError(err)
	}
	return book, nil
}

func (s *BookRepository) GetByID(ctx context.Context, id string) (domain.Book, error) {
	return scanBook(s.pool.QueryRow(ctx, `
		SELECT id, title, author, isbn, year, created_at, updated_at
		FROM books
		WHERE id = $1
	`, id))
}

func (s *BookRepository) GetByISBN(ctx context.Context, isbn string) (domain.Book, error) {
	if isbn == "" {
		return domain.Book{}, domain.ErrNotFound
	}
	return scanBook(s.pool.QueryRow(ctx, `
		SELECT id, title, author, isbn, year, created_at, updated_at
		FROM books
		WHERE isbn = $1
	`, isbn))
}

func (s *BookRepository) List(ctx context.Context, filter domain.ListFilter) ([]domain.Book, int, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, title, author, isbn, year, created_at, updated_at, COUNT(*) OVER() AS total
		FROM books
		WHERE ($1 = '' OR title ILIKE '%' || $1 || '%' OR author ILIKE '%' || $1 || '%')
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, filter.Query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list books: %w", err)
	}
	defer rows.Close()

	books := make([]domain.Book, 0)
	total := 0
	for rows.Next() {
		book, rowTotal, err := scanBookWithTotal(rows)
		if err != nil {
			return nil, 0, err
		}
		total = rowTotal
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list books: %w", err)
	}
	return books, total, nil
}

func (s *BookRepository) Update(ctx context.Context, book domain.Book) (domain.Book, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE books
		SET title = $2, author = $3, isbn = $4, year = $5, updated_at = $6
		WHERE id = $1
	`, book.ID, book.Title, book.Author, nullString(book.ISBN), nullInt(book.Year), book.UpdatedAt)
	if err != nil {
		return domain.Book{}, mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.Book{}, domain.ErrNotFound
	}
	return book, nil
}

func (s *BookRepository) Delete(ctx context.Context, id string) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanBook(row rowScanner) (domain.Book, error) {
	book, _, err := scanBookRow(row, false)
	return book, err
}

func scanBookWithTotal(row rowScanner) (domain.Book, int, error) {
	return scanBookRow(row, true)
}

func scanBookRow(row rowScanner, withTotal bool) (domain.Book, int, error) {
	var (
		book  domain.Book
		isbn  *string
		year  *int
		total int
	)

	dest := []any{&book.ID, &book.Title, &book.Author, &isbn, &year, &book.CreatedAt, &book.UpdatedAt}
	if withTotal {
		dest = append(dest, &total)
	}

	if err := row.Scan(dest...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Book{}, 0, domain.ErrNotFound
		}
		return domain.Book{}, 0, fmt.Errorf("scan book: %w", err)
	}
	if isbn != nil {
		book.ISBN = *isbn
	}
	if year != nil {
		book.Year = *year
	}
	return book, total, nil
}

func mapError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrDuplicateISBN
	}
	return err
}

func nullString(v string) any {
	if v == "" {
		return nil
	}
	return v
}

func nullInt(v int) any {
	if v == 0 {
		return nil
	}
	return v
}

var _ domain.BookRepository = (*BookRepository)(nil)
