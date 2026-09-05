package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"golang-rest-api/internal/domain"
)

type PrivilegeRepository struct {
	mu     sync.RWMutex
	items  map[string]domain.Privilege
	byCode map[string]string
}

func NewPrivilegeRepository() *PrivilegeRepository {
	return &PrivilegeRepository{
		items:  make(map[string]domain.Privilege),
		byCode: make(map[string]string),
	}
}

func (s *PrivilegeRepository) Create(_ context.Context, privilege domain.Privilege) (domain.Privilege, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.byCode[strings.ToLower(privilege.Code)]; exists {
		return domain.Privilege{}, domain.ErrDuplicatePrivilegeCode
	}

	s.items[privilege.ID] = privilege
	s.byCode[strings.ToLower(privilege.Code)] = privilege.ID
	return privilege, nil
}

func (s *PrivilegeRepository) GetByID(_ context.Context, id string) (domain.Privilege, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return domain.Privilege{}, domain.ErrNotFound
	}
	return item, nil
}

func (s *PrivilegeRepository) GetByCode(_ context.Context, code string) (domain.Privilege, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.byCode[strings.ToLower(code)]
	if !ok {
		return domain.Privilege{}, domain.ErrNotFound
	}
	return s.items[id], nil
}

func (s *PrivilegeRepository) List(_ context.Context, filter domain.ListFilter) ([]domain.Privilege, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := strings.ToLower(strings.TrimSpace(filter.Query))
	matched := make([]domain.Privilege, 0, len(s.items))

	for _, item := range s.items {
		if query == "" ||
			strings.Contains(strings.ToLower(item.Code), query) ||
			strings.Contains(strings.ToLower(item.Name), query) ||
			strings.Contains(strings.ToLower(item.Description), query) {
			matched = append(matched, item)
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

func (s *PrivilegeRepository) Update(_ context.Context, privilege domain.Privilege) (domain.Privilege, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.items[privilege.ID]
	if !ok {
		return domain.Privilege{}, domain.ErrNotFound
	}

	if owner, exists := s.byCode[strings.ToLower(privilege.Code)]; exists && owner != privilege.ID {
		return domain.Privilege{}, domain.ErrDuplicatePrivilegeCode
	}

	delete(s.byCode, strings.ToLower(current.Code))
	s.byCode[strings.ToLower(privilege.Code)] = privilege.ID
	s.items[privilege.ID] = privilege
	return privilege, nil
}

func (s *PrivilegeRepository) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return domain.ErrNotFound
	}

	delete(s.items, id)
	delete(s.byCode, strings.ToLower(item.Code))
	return nil
}

var _ domain.PrivilegeRepository = (*PrivilegeRepository)(nil)
