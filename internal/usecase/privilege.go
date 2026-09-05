package usecase

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"golang-rest-api/internal/domain"
)

type PrivilegeUseCase struct {
	repo domain.PrivilegeRepository
	now  func() time.Time
}

func NewPrivilegeUseCase(repo domain.PrivilegeRepository) *PrivilegeUseCase {
	return &PrivilegeUseCase{
		repo: repo,
		now:  time.Now,
	}
}

func (s *PrivilegeUseCase) Create(ctx context.Context, input domain.CreatePrivilegeInput) (domain.Privilege, error) {
	code, name, description, createdBy, err := normalizePrivilege(input.Code, input.Name, input.Description, input.CreatedBy)
	if err != nil {
		return domain.Privilege{}, err
	}
	if err := s.ensureUniqueCode(ctx, "", code); err != nil {
		return domain.Privilege{}, err
	}

	now := s.now()
	privilege := domain.Privilege{
		ID:          newUUID(),
		Code:        code,
		Name:        name,
		Description: description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if input.IsActive != nil {
		privilege.IsActive = *input.IsActive
	}

	return s.repo.Create(ctx, privilege)
}

func (s *PrivilegeUseCase) Get(ctx context.Context, id string) (domain.Privilege, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.Privilege{}, domain.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *PrivilegeUseCase) List(ctx context.Context, filter domain.ListFilter) ([]domain.Privilege, int, error) {
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

func (s *PrivilegeUseCase) Update(ctx context.Context, id string, input domain.UpdatePrivilegeInput) (domain.Privilege, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.Privilege{}, domain.ErrInvalidInput
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Privilege{}, err
	}

	code, name, description, updatedBy, err := normalizePrivilege(input.Code, input.Name, input.Description, input.UpdatedBy)
	if err != nil {
		return domain.Privilege{}, err
	}
	if err := s.ensureUniqueCode(ctx, current.ID, code); err != nil {
		return domain.Privilege{}, err
	}

	current.Code = code
	current.Name = name
	current.Description = description
	current.UpdatedAt = s.now()
	current.UpdatedBy = updatedBy
	if input.IsActive != nil {
		current.IsActive = *input.IsActive
	}

	return s.repo.Update(ctx, current)
}

func (s *PrivilegeUseCase) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

func (s *PrivilegeUseCase) ensureUniqueCode(ctx context.Context, currentID, code string) error {
	existing, err := s.repo.GetByCode(ctx, code)
	if err == nil && existing.ID != currentID {
		return domain.ErrDuplicatePrivilegeCode
	}
	if err != nil && err != domain.ErrNotFound {
		return err
	}
	return nil
}

func normalizePrivilege(code, name, description string, actor *string) (string, string, string, *string, error) {
	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	if code == "" || name == "" {
		return "", "", "", nil, domain.ErrInvalidInput
	}
	if strings.ContainsAny(code, " \t") || utf8.RuneCountInString(code) > 80 {
		return "", "", "", nil, domain.ErrInvalidInput
	}
	if utf8.RuneCountInString(name) > 120 || utf8.RuneCountInString(description) > 500 {
		return "", "", "", nil, domain.ErrInvalidInput
	}

	actorID, err := optionalActor(actor)
	if err != nil {
		return "", "", "", nil, err
	}

	return code, name, description, actorID, nil
}
