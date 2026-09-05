package usecase

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"golang-rest-api/internal/domain"
)

type MenuUseCase struct {
	repo domain.MenuRepository
	now  func() time.Time
}

func NewMenuUseCase(repo domain.MenuRepository) *MenuUseCase {
	return &MenuUseCase{
		repo: repo,
		now:  time.Now,
	}
}

func (s *MenuUseCase) Create(ctx context.Context, input domain.CreateMenuInput) (domain.Menu, error) {
	menu, err := normalizeMenu(input.ParentID, input.Code, input.Name, input.Path, input.Icon, input.Description, input.Type, input.SortOrder, input.CreatedBy)
	if err != nil {
		return domain.Menu{}, err
	}
	if err := s.ensureUniqueCode(ctx, "", menu.Code); err != nil {
		return domain.Menu{}, err
	}
	if err := s.ensureParent(ctx, "", menu.ParentID); err != nil {
		return domain.Menu{}, err
	}

	now := s.now()
	menu.ID = newUUID()
	menu.IsActive = true
	menu.CreatedAt = now
	menu.UpdatedAt = now
	menu.CreatedBy = menu.UpdatedBy
	if input.IsActive != nil {
		menu.IsActive = *input.IsActive
	}

	return s.repo.Create(ctx, menu)
}

func (s *MenuUseCase) Get(ctx context.Context, id string) (domain.Menu, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.Menu{}, domain.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *MenuUseCase) List(ctx context.Context, filter domain.MenuListFilter) ([]domain.Menu, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	if filter.Limit > maxLimit {
		filter.Limit = maxLimit
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	if filter.ParentID != "" && !isUUID(filter.ParentID) {
		return nil, 0, domain.ErrInvalidInput
	}
	return s.repo.List(ctx, filter)
}

func (s *MenuUseCase) Update(ctx context.Context, id string, input domain.UpdateMenuInput) (domain.Menu, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.Menu{}, domain.ErrInvalidInput
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Menu{}, err
	}

	next, err := normalizeMenu(input.ParentID, input.Code, input.Name, input.Path, input.Icon, input.Description, input.Type, input.SortOrder, input.UpdatedBy)
	if err != nil {
		return domain.Menu{}, err
	}
	if err := s.ensureUniqueCode(ctx, current.ID, next.Code); err != nil {
		return domain.Menu{}, err
	}
	if err := s.ensureParent(ctx, current.ID, next.ParentID); err != nil {
		return domain.Menu{}, err
	}

	current.ParentID = next.ParentID
	current.Code = next.Code
	current.Name = next.Name
	current.Path = next.Path
	current.Icon = next.Icon
	current.Description = next.Description
	current.SortOrder = next.SortOrder
	current.Type = next.Type
	current.UpdatedAt = s.now()
	current.UpdatedBy = next.UpdatedBy
	if input.IsActive != nil {
		current.IsActive = *input.IsActive
	}

	return s.repo.Update(ctx, current)
}

func (s *MenuUseCase) Delete(ctx context.Context, id string, input domain.DeleteMenuInput) (domain.Menu, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.Menu{}, domain.ErrInvalidInput
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Menu{}, err
	}

	children, err := s.repo.CountChildren(ctx, current.ID)
	if err != nil {
		return domain.Menu{}, err
	}
	if children > 0 {
		return domain.Menu{}, domain.ErrMenuHasChildren
	}

	deletedBy, err := optionalActor(input.DeletedBy)
	if err != nil {
		return domain.Menu{}, err
	}

	now := s.now()
	current.DeletedAt = &now
	current.DeletedBy = deletedBy
	current.UpdatedAt = now
	current.UpdatedBy = deletedBy

	if err := s.repo.Delete(ctx, current); err != nil {
		return domain.Menu{}, err
	}
	return current, nil
}

func (s *MenuUseCase) ensureUniqueCode(ctx context.Context, currentID, code string) error {
	existing, err := s.repo.GetByCode(ctx, code)
	if err == nil && existing.ID != currentID {
		return domain.ErrDuplicateMenuCode
	}
	if err != nil && err != domain.ErrNotFound {
		return err
	}
	return nil
}

func (s *MenuUseCase) ensureParent(ctx context.Context, menuID string, parentID *string) error {
	if parentID == nil {
		return nil
	}
	if menuID != "" && *parentID == menuID {
		return domain.ErrInvalidParent
	}

	seen := map[string]struct{}{}
	current := parentID
	for current != nil {
		if menuID != "" && *current == menuID {
			return domain.ErrInvalidParent
		}
		if _, ok := seen[*current]; ok {
			return domain.ErrInvalidParent
		}
		seen[*current] = struct{}{}

		parent, err := s.repo.GetByID(ctx, *current)
		if err != nil {
			if err == domain.ErrNotFound {
				return domain.ErrInvalidParent
			}
			return err
		}
		current = parent.ParentID
	}
	return nil
}

func normalizeMenu(parentID *string, code, name, path, icon, description, menuType string, sortOrder *int, actor *string) (domain.Menu, error) {
	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)
	path = strings.TrimSpace(path)
	icon = strings.TrimSpace(icon)
	description = strings.TrimSpace(description)
	menuType = strings.ToUpper(strings.TrimSpace(menuType))

	if code == "" || name == "" {
		return domain.Menu{}, domain.ErrInvalidInput
	}
	if strings.ContainsAny(code, " \t") || utf8.RuneCountInString(code) > 80 {
		return domain.Menu{}, domain.ErrInvalidInput
	}
	if utf8.RuneCountInString(name) > 120 || utf8.RuneCountInString(path) > 255 || utf8.RuneCountInString(icon) > 80 || utf8.RuneCountInString(description) > 500 {
		return domain.Menu{}, domain.ErrInvalidInput
	}
	if menuType == "" {
		menuType = domain.MenuTypeItem
	}
	if menuType != domain.MenuTypeItem && menuType != domain.MenuTypeGroup && menuType != domain.MenuTypeCollapse {
		return domain.Menu{}, domain.ErrInvalidInput
	}

	parent, err := optionalActor(parentID)
	if err != nil {
		return domain.Menu{}, domain.ErrInvalidParent
	}
	actorID, err := optionalActor(actor)
	if err != nil {
		return domain.Menu{}, err
	}

	order := 0
	if sortOrder != nil {
		if *sortOrder < 0 {
			return domain.Menu{}, domain.ErrInvalidInput
		}
		order = *sortOrder
	}

	return domain.Menu{
		ParentID:    parent,
		Code:        code,
		Name:        name,
		Path:        path,
		Icon:        icon,
		Description: description,
		SortOrder:   order,
		Type:        menuType,
		UpdatedBy:   actorID,
	}, nil
}
