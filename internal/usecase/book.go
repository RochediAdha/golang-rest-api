package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
	"unicode/utf8"

	"golang-rest-api/internal/domain"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type BookUseCase struct {
	repo domain.BookRepository
	now  func() time.Time
}

func NewBookUseCase(repo domain.BookRepository) *BookUseCase {
	return &BookUseCase{
		repo: repo,
		now:  time.Now,
	}
}

func (s *BookUseCase) Create(ctx context.Context, input domain.CreateBookInput) (domain.Book, error) {
	title, author, isbn, year, err := normalizeBook(input.Title, input.Author, input.ISBN, input.Year)
	if err != nil {
		return domain.Book{}, err
	}

	now := s.now().UTC()
	book := domain.Book{
		ID:        newID(),
		Title:     title,
		Author:    author,
		ISBN:      isbn,
		Year:      year,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.Create(ctx, book)
}

func (s *BookUseCase) Get(ctx context.Context, id string) (domain.Book, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.Book{}, domain.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *BookUseCase) List(ctx context.Context, filter domain.ListFilter) ([]domain.Book, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	if filter.Limit > maxLimit {
		filter.Limit = maxLimit
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return s.repo.List(ctx, filter)
}

func (s *BookUseCase) Update(ctx context.Context, id string, input domain.UpdateBookInput) (domain.Book, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.Book{}, domain.ErrInvalidInput
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Book{}, err
	}

	title, author, isbn, year, err := normalizeBook(input.Title, input.Author, input.ISBN, input.Year)
	if err != nil {
		return domain.Book{}, err
	}

	current.Title = title
	current.Author = author
	current.ISBN = isbn
	current.Year = year
	current.UpdatedAt = s.now().UTC()

	return s.repo.Update(ctx, current)
}

func (s *BookUseCase) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

func normalizeBook(title, author, isbn string, year int) (string, string, string, int, error) {
	title = strings.TrimSpace(title)
	author = strings.TrimSpace(author)
	isbn = strings.ReplaceAll(strings.TrimSpace(isbn), "-", "")

	if title == "" || author == "" {
		return "", "", "", 0, domain.ErrInvalidInput
	}
	if utf8.RuneCountInString(title) > 200 || utf8.RuneCountInString(author) > 120 {
		return "", "", "", 0, domain.ErrInvalidInput
	}
	if isbn != "" && (len(isbn) != 10 && len(isbn) != 13) {
		return "", "", "", 0, domain.ErrInvalidInput
	}
	if year != 0 && (year < 1000 || year > time.Now().Year()+1) {
		return "", "", "", 0, domain.ErrInvalidInput
	}

	return title, author, isbn, year, nil
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
