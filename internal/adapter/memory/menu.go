package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"golang-rest-api/internal/domain"
)

type MenuRepository struct {
	mu     sync.RWMutex
	menus  map[string]domain.Menu
	byCode map[string]string
}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{
		menus:  make(map[string]domain.Menu),
		byCode: make(map[string]string),
	}
}

func (s *MenuRepository) Create(_ context.Context, menu domain.Menu) (domain.Menu, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.byCode[strings.ToLower(menu.Code)]; exists {
		return domain.Menu{}, domain.ErrDuplicateMenuCode
	}

	s.menus[menu.ID] = menu
	s.byCode[strings.ToLower(menu.Code)] = menu.ID
	return menu, nil
}

func (s *MenuRepository) GetByID(_ context.Context, id string) (domain.Menu, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	menu, ok := s.menus[id]
	if !ok || menu.DeletedAt != nil {
		return domain.Menu{}, domain.ErrNotFound
	}
	return menu, nil
}

func (s *MenuRepository) GetByCode(_ context.Context, code string) (domain.Menu, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.byCode[strings.ToLower(code)]
	if !ok {
		return domain.Menu{}, domain.ErrNotFound
	}
	menu := s.menus[id]
	if menu.DeletedAt != nil {
		return domain.Menu{}, domain.ErrNotFound
	}
	return menu, nil
}

func (s *MenuRepository) List(_ context.Context, filter domain.MenuListFilter) ([]domain.Menu, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := strings.ToLower(strings.TrimSpace(filter.Query))
	matched := make([]domain.Menu, 0, len(s.menus))

	for _, menu := range s.menus {
		if menu.DeletedAt != nil {
			continue
		}
		if filter.RootOnly && menu.ParentID != nil {
			continue
		}
		if filter.ParentID != "" && (menu.ParentID == nil || *menu.ParentID != filter.ParentID) {
			continue
		}
		if query == "" ||
			strings.Contains(strings.ToLower(menu.Code), query) ||
			strings.Contains(strings.ToLower(menu.Name), query) ||
			strings.Contains(strings.ToLower(menu.Path), query) {
			matched = append(matched, menu)
		}
	}

	sort.Slice(matched, func(i, j int) bool {
		if matched[i].SortOrder != matched[j].SortOrder {
			return matched[i].SortOrder < matched[j].SortOrder
		}
		return matched[i].CreatedAt.Before(matched[j].CreatedAt)
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

func (s *MenuRepository) CountChildren(_ context.Context, parentID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, menu := range s.menus {
		if menu.DeletedAt == nil && menu.ParentID != nil && *menu.ParentID == parentID {
			count++
		}
	}
	return count, nil
}

func (s *MenuRepository) Update(_ context.Context, menu domain.Menu) (domain.Menu, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.menus[menu.ID]
	if !ok || current.DeletedAt != nil {
		return domain.Menu{}, domain.ErrNotFound
	}

	if owner, exists := s.byCode[strings.ToLower(menu.Code)]; exists && owner != menu.ID {
		existing := s.menus[owner]
		if existing.DeletedAt == nil {
			return domain.Menu{}, domain.ErrDuplicateMenuCode
		}
	}

	delete(s.byCode, strings.ToLower(current.Code))
	s.byCode[strings.ToLower(menu.Code)] = menu.ID
	s.menus[menu.ID] = menu
	return menu, nil
}

func (s *MenuRepository) Delete(_ context.Context, menu domain.Menu) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.menus[menu.ID]
	if !ok || current.DeletedAt != nil {
		return domain.ErrNotFound
	}

	s.menus[menu.ID] = menu
	delete(s.byCode, strings.ToLower(current.Code))
	return nil
}

var _ domain.MenuRepository = (*MenuRepository)(nil)
