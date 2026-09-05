package memory

import (
	"context"
	"sort"
	"sync"

	"golang-rest-api/internal/domain"
)

type RolePrivilegeRepository struct {
	mu    sync.RWMutex
	items map[string]domain.RolePrivilege
}

func NewRolePrivilegeRepository() *RolePrivilegeRepository {
	return &RolePrivilegeRepository{
		items: make(map[string]domain.RolePrivilege),
	}
}

func (s *RolePrivilegeRepository) Create(_ context.Context, rolePrivilege domain.RolePrivilege) (domain.RolePrivilege, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.items {
		if existing.RoleID == rolePrivilege.RoleID && existing.MenuID == rolePrivilege.MenuID && existing.PrivilegeID == rolePrivilege.PrivilegeID {
			return domain.RolePrivilege{}, domain.ErrDuplicateRolePrivilege
		}
	}

	s.items[rolePrivilege.ID] = rolePrivilege
	return rolePrivilege, nil
}

func (s *RolePrivilegeRepository) GetByID(_ context.Context, id string) (domain.RolePrivilege, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return domain.RolePrivilege{}, domain.ErrNotFound
	}
	return item, nil
}

func (s *RolePrivilegeRepository) GetByRoleMenuPrivilege(_ context.Context, roleID, menuID, privilegeID string) (domain.RolePrivilege, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.items {
		if item.RoleID == roleID && item.MenuID == menuID && item.PrivilegeID == privilegeID {
			return item, nil
		}
	}
	return domain.RolePrivilege{}, domain.ErrNotFound
}

func (s *RolePrivilegeRepository) List(_ context.Context, filter domain.RolePrivilegeListFilter) ([]domain.RolePrivilege, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	matched := make([]domain.RolePrivilege, 0, len(s.items))
	for _, item := range s.items {
		if filter.RoleID != "" && item.RoleID != filter.RoleID {
			continue
		}
		if filter.MenuID != "" && item.MenuID != filter.MenuID {
			continue
		}
		if filter.PrivilegeID != "" && item.PrivilegeID != filter.PrivilegeID {
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

func (s *RolePrivilegeRepository) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.items, id)
	return nil
}

var _ domain.RolePrivilegeRepository = (*RolePrivilegeRepository)(nil)
