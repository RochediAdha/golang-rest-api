package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"golang-rest-api/internal/domain"
)

type UserRepository struct {
	mu         sync.RWMutex
	users      map[string]domain.User
	byUsername map[string]string
	byEmail    map[string]string
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:      make(map[string]domain.User),
		byUsername: make(map[string]string),
		byEmail:    make(map[string]string),
	}
}

func (s *UserRepository) Create(_ context.Context, user domain.User) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.byUsername[strings.ToLower(user.Username)]; exists {
		return domain.User{}, domain.ErrDuplicateUsername
	}
	if _, exists := s.byEmail[strings.ToLower(user.Email)]; exists {
		return domain.User{}, domain.ErrDuplicateEmail
	}

	s.users[user.ID] = user
	s.byUsername[strings.ToLower(user.Username)] = user.ID
	s.byEmail[strings.ToLower(user.Email)] = user.ID
	return user, nil
}

func (s *UserRepository) GetByID(_ context.Context, id string) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

func (s *UserRepository) GetByUsername(_ context.Context, username string) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.byUsername[strings.ToLower(username)]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return s.users[id], nil
}

func (s *UserRepository) GetByEmail(_ context.Context, email string) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.byEmail[strings.ToLower(email)]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return s.users[id], nil
}

func (s *UserRepository) List(_ context.Context, filter domain.ListFilter) ([]domain.User, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := strings.ToLower(strings.TrimSpace(filter.Query))
	matched := make([]domain.User, 0, len(s.users))

	for _, user := range s.users {
		if query == "" ||
			strings.Contains(strings.ToLower(user.Username), query) ||
			strings.Contains(strings.ToLower(user.Email), query) ||
			strings.Contains(strings.ToLower(user.Name), query) {
			matched = append(matched, user)
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

func (s *UserRepository) Update(_ context.Context, user domain.User) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.users[user.ID]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}

	if owner, exists := s.byUsername[strings.ToLower(user.Username)]; exists && owner != user.ID {
		return domain.User{}, domain.ErrDuplicateUsername
	}
	if owner, exists := s.byEmail[strings.ToLower(user.Email)]; exists && owner != user.ID {
		return domain.User{}, domain.ErrDuplicateEmail
	}

	delete(s.byUsername, strings.ToLower(current.Username))
	delete(s.byEmail, strings.ToLower(current.Email))
	s.byUsername[strings.ToLower(user.Username)] = user.ID
	s.byEmail[strings.ToLower(user.Email)] = user.ID
	s.users[user.ID] = user
	return user, nil
}

func (s *UserRepository) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return domain.ErrNotFound
	}

	delete(s.users, id)
	delete(s.byUsername, strings.ToLower(user.Username))
	delete(s.byEmail, strings.ToLower(user.Email))
	return nil
}

var _ domain.UserRepository = (*UserRepository)(nil)
