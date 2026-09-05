package memory

import (
	"context"
	"sort"
	"sync"

	"golang-rest-api/internal/domain"
)

type UserRoleRepository struct {
	mu    sync.RWMutex
	items map[string]domain.UserRole
}

func NewUserRoleRepository() *UserRoleRepository {
	return &UserRoleRepository{
		items: make(map[string]domain.UserRole),
	}
}

func (s *UserRoleRepository) Create(_ context.Context, userRole domain.UserRole) (domain.UserRole, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.items {
		if existing.UserID == userRole.UserID && existing.RoleID == userRole.RoleID {
			return domain.UserRole{}, domain.ErrDuplicateUserRole
		}
		if userRole.Number != "" && existing.RoleID == userRole.RoleID && existing.Number == userRole.Number {
			return domain.UserRole{}, domain.ErrDuplicateUserRoleNumber
		}
	}

	s.items[userRole.ID] = userRole
	return userRole, nil
}

func (s *UserRoleRepository) GetByID(_ context.Context, id string) (domain.UserRole, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return domain.UserRole{}, domain.ErrNotFound
	}
	return item, nil
}

func (s *UserRoleRepository) GetByUserAndRole(_ context.Context, userID, roleID string) (domain.UserRole, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.items {
		if item.UserID == userID && item.RoleID == roleID {
			return item, nil
		}
	}
	return domain.UserRole{}, domain.ErrNotFound
}

func (s *UserRoleRepository) GetByRoleAndNumber(_ context.Context, roleID, number string) (domain.UserRole, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.items {
		if item.RoleID == roleID && item.Number == number {
			return item, nil
		}
	}
	return domain.UserRole{}, domain.ErrNotFound
}

func (s *UserRoleRepository) List(_ context.Context, filter domain.UserRoleListFilter) ([]domain.UserRole, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	matched := make([]domain.UserRole, 0, len(s.items))
	for _, item := range s.items {
		if filter.UserID != "" && item.UserID != filter.UserID {
			continue
		}
		if filter.RoleID != "" && item.RoleID != filter.RoleID {
			continue
		}
		if filter.Number != "" && item.Number != filter.Number {
			continue
		}
		matched = append(matched, item)
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

func (s *UserRoleRepository) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.items, id)
	return nil
}

var _ domain.UserRoleRepository = (*UserRoleRepository)(nil)
