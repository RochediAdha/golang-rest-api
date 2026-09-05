package usecase

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"golang-rest-api/internal/domain"
)

type RoleUseCase struct {
	repo domain.RoleRepository
	now  func() time.Time
}

func NewRoleUseCase(repo domain.RoleRepository) *RoleUseCase {
	return &RoleUseCase{
		repo: repo,
		now:  time.Now,
	}
}

func (s *RoleUseCase) Create(ctx context.Context, input domain.CreateRoleInput) (domain.Role, error) {
	name, description, createdBy, err := normalizeRole(input.Name, input.Description, input.CreatedBy)
	if err != nil {
		return domain.Role{}, err
	}
	if err := s.ensureUniqueName(ctx, "", name); err != nil {
		return domain.Role{}, err
	}

	now := s.now()
	role := domain.Role{
		ID:          newUUID(),
		Name:        name,
		Description: description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if input.IsActive != nil {
		role.IsActive = *input.IsActive
	}

	return s.repo.Create(ctx, role)
}

func (s *RoleUseCase) Get(ctx context.Context, id string) (domain.Role, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.Role{}, domain.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *RoleUseCase) List(ctx context.Context, filter domain.ListFilter) ([]domain.Role, int, error) {
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

func (s *RoleUseCase) Update(ctx context.Context, id string, input domain.UpdateRoleInput) (domain.Role, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.Role{}, domain.ErrInvalidInput
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Role{}, err
	}

	name, description, updatedBy, err := normalizeRole(input.Name, input.Description, input.UpdatedBy)
	if err != nil {
		return domain.Role{}, err
	}
	if err := s.ensureUniqueName(ctx, current.ID, name); err != nil {
		return domain.Role{}, err
	}

	current.Name = name
	current.Description = description
	current.UpdatedAt = s.now()
	current.UpdatedBy = updatedBy
	if input.IsActive != nil {
		current.IsActive = *input.IsActive
	}

	return s.repo.Update(ctx, current)
}

func (s *RoleUseCase) Delete(ctx context.Context, id string, input domain.DeleteRoleInput) (domain.Role, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.Role{}, domain.ErrInvalidInput
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Role{}, err
	}

	deletedBy, err := optionalActor(input.DeletedBy)
	if err != nil {
		return domain.Role{}, err
	}

	now := s.now()
	current.DeletedAt = &now
	current.DeletedBy = deletedBy
	current.UpdatedAt = now
	current.UpdatedBy = deletedBy

	if err := s.repo.Delete(ctx, current); err != nil {
		return domain.Role{}, err
	}
	return current, nil
}

func (s *RoleUseCase) ensureUniqueName(ctx context.Context, currentID, name string) error {
	existing, err := s.repo.GetByName(ctx, name)
	if err == nil && existing.ID != currentID {
		return domain.ErrDuplicateRoleName
	}
	if err != nil && err != domain.ErrNotFound {
		return err
	}
	return nil
}

func normalizeRole(name, description string, actor *string) (string, string, *string, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	if name == "" {
		return "", "", nil, domain.ErrInvalidInput
	}
	if utf8.RuneCountInString(name) > 80 {
		return "", "", nil, domain.ErrInvalidInput
	}
	if utf8.RuneCountInString(description) > 500 {
		return "", "", nil, domain.ErrInvalidInput
	}

	actorID, err := optionalActor(actor)
	if err != nil {
		return "", "", nil, err
	}

	return name, description, actorID, nil
}

func optionalActor(actor *string) (*string, error) {
	if actor == nil {
		return nil, nil
	}
	value := strings.TrimSpace(*actor)
	if value == "" {
		return nil, nil
	}
	if !isUUID(value) {
		return nil, domain.ErrInvalidInput
	}
	return &value, nil
}
