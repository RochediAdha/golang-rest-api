package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"golang-rest-api/internal/domain"
)

type BookRepository struct {
	mu     sync.RWMutex
	books  map[string]domain.Book
	byISBN map[string]string
}

func NewBookRepository() *BookRepository {
	return &BookRepository{
		books:  make(map[string]domain.Book),
		byISBN: make(map[string]string),
	}
}

func (s *BookRepository) Create(_ context.Context, book domain.Book) (domain.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if book.ISBN != "" {
		if _, exists := s.byISBN[book.ISBN]; exists {
			return domain.Book{}, domain.ErrDuplicateISBN
		}
		s.byISBN[book.ISBN] = book.ID
	}

	s.books[book.ID] = book
	return book, nil
}

func (s *BookRepository) GetByID(_ context.Context, id string) (domain.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.books[id]
	if !ok {
		return domain.Book{}, domain.ErrNotFound
	}
	return book, nil
}

func (s *BookRepository) GetByISBN(_ context.Context, isbn string) (domain.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.byISBN[isbn]
	if !ok {
		return domain.Book{}, domain.ErrNotFound
	}
	return s.books[id], nil
}

func (s *BookRepository) List(_ context.Context, filter domain.ListFilter) ([]domain.Book, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := strings.ToLower(strings.TrimSpace(filter.Query))
	matched := make([]domain.Book, 0, len(s.books))

	for _, book := range s.books {
		if query == "" || strings.Contains(strings.ToLower(book.Title), query) || strings.Contains(strings.ToLower(book.Author), query) {
			matched = append(matched, book)
		}
	}

	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})

	total := len(matched)
	start := filter.Offset
	if start > total {
		start = total
	}
	end := start + filter.Limit
	if end > total {
		end = total
	}

	return matched[start:end], total, nil
}

func (s *BookRepository) Update(_ context.Context, book domain.Book) (domain.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.books[book.ID]
	if !ok {
		return domain.Book{}, domain.ErrNotFound
	}

	if current.ISBN != "" && current.ISBN != book.ISBN {
		delete(s.byISBN, current.ISBN)
	}
	if book.ISBN != "" {
		if owner, exists := s.byISBN[book.ISBN]; exists && owner != book.ID {
			return domain.Book{}, domain.ErrDuplicateISBN
		}
		s.byISBN[book.ISBN] = book.ID
	}

	s.books[book.ID] = book
	return book, nil
}

func (s *BookRepository) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	book, ok := s.books[id]
	if !ok {
		return domain.ErrNotFound
	}

	delete(s.books, id)
	if book.ISBN != "" {
		delete(s.byISBN, book.ISBN)
	}
	return nil
}

var _ domain.BookRepository = (*BookRepository)(nil)
