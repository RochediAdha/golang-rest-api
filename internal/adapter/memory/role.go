package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"golang-rest-api/internal/domain"
)

type RoleRepository struct {
	mu     sync.RWMutex
	roles  map[string]domain.Role
	byName map[string]string
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		roles:  make(map[string]domain.Role),
		byName: make(map[string]string),
	}
}

func (s *RoleRepository) Create(_ context.Context, role domain.Role) (domain.Role, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.byName[strings.ToLower(role.Name)]; exists {
		return domain.Role{}, domain.ErrDuplicateRoleName
	}

	s.roles[role.ID] = role
	s.byName[strings.ToLower(role.Name)] = role.ID
	return role, nil
}

func (s *RoleRepository) GetByID(_ context.Context, id string) (domain.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	role, ok := s.roles[id]
	if !ok || role.DeletedAt != nil {
		return domain.Role{}, domain.ErrNotFound
	}
	return role, nil
}

func (s *RoleRepository) GetByName(_ context.Context, name string) (domain.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.byName[strings.ToLower(name)]
	if !ok {
		return domain.Role{}, domain.ErrNotFound
	}
	role := s.roles[id]
	if role.DeletedAt != nil {
		return domain.Role{}, domain.ErrNotFound
	}
	return role, nil
}

func (s *RoleRepository) List(_ context.Context, filter domain.ListFilter) ([]domain.Role, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := strings.ToLower(strings.TrimSpace(filter.Query))
	matched := make([]domain.Role, 0, len(s.roles))

	for _, role := range s.roles {
		if role.DeletedAt != nil {
			continue
		}
		if query == "" ||
			strings.Contains(strings.ToLower(role.Name), query) ||
			strings.Contains(strings.ToLower(role.Description), query) {
			matched = append(matched, role)
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

func (s *RoleRepository) Update(_ context.Context, role domain.Role) (domain.Role, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.roles[role.ID]
	if !ok || current.DeletedAt != nil {
		return domain.Role{}, domain.ErrNotFound
	}

	if owner, exists := s.byName[strings.ToLower(role.Name)]; exists && owner != role.ID {
		existing := s.roles[owner]
		if existing.DeletedAt == nil {
			return domain.Role{}, domain.ErrDuplicateRoleName
		}
	}

	delete(s.byName, strings.ToLower(current.Name))
	s.byName[strings.ToLower(role.Name)] = role.ID
	s.roles[role.ID] = role
	return role, nil
}

func (s *RoleRepository) Delete(_ context.Context, role domain.Role) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.roles[role.ID]
	if !ok || current.DeletedAt != nil {
		return domain.ErrNotFound
	}

	s.roles[role.ID] = role
	delete(s.byName, strings.ToLower(current.Name))
	return nil
}

var _ domain.RoleRepository = (*RoleRepository)(nil)
