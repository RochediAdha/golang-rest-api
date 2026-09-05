package usecase

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"golang-rest-api/internal/domain"
)

type UserRoleUseCase struct {
	repo  domain.UserRoleRepository
	users domain.UserRepository
	roles domain.RoleRepository
	now   func() time.Time
}

func NewUserRoleUseCase(repo domain.UserRoleRepository, users domain.UserRepository, roles domain.RoleRepository) *UserRoleUseCase {
	return &UserRoleUseCase{
		repo:  repo,
		users: users,
		roles: roles,
		now:   time.Now,
	}
}

func (s *UserRoleUseCase) Create(ctx context.Context, input domain.CreateUserRoleInput) (domain.UserRole, error) {
	userID, roleID, number, createdBy, err := normalizeUserRole(input.UserID, input.RoleID, input.Number, input.CreatedBy)
	if err != nil {
		return domain.UserRole{}, err
	}
	if err := s.ensureUser(ctx, userID); err != nil {
		return domain.UserRole{}, err
	}
	role, err := s.roles.GetByID(ctx, roleID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.UserRole{}, domain.ErrInvalidRole
		}
		return domain.UserRole{}, err
	}
	if requiresRoleNumber(role.Name) && number == "" {
		return domain.UserRole{}, domain.ErrInvalidInput
	}
	if err := s.ensureUnique(ctx, userID, roleID); err != nil {
		return domain.UserRole{}, err
	}
	if err := s.ensureUniqueNumber(ctx, roleID, number); err != nil {
		return domain.UserRole{}, err
	}

	return s.repo.Create(ctx, domain.UserRole{
		ID:        newUUID(),
		UserID:    userID,
		RoleID:    roleID,
		Number:    number,
		CreatedAt: s.now(),
		CreatedBy: createdBy,
	})
}

func (s *UserRoleUseCase) Get(ctx context.Context, id string) (domain.UserRole, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.UserRole{}, domain.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *UserRoleUseCase) List(ctx context.Context, filter domain.UserRoleListFilter) ([]domain.UserRole, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	if filter.Limit > maxLimit {
		filter.Limit = maxLimit
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	filter.UserID = strings.TrimSpace(filter.UserID)
	filter.RoleID = strings.TrimSpace(filter.RoleID)
	filter.Number = strings.TrimSpace(filter.Number)
	if filter.UserID != "" && !isUUID(filter.UserID) {
		return nil, 0, domain.ErrInvalidInput
	}
	if filter.RoleID != "" && !isUUID(filter.RoleID) {
		return nil, 0, domain.ErrInvalidInput
	}

	return s.repo.List(ctx, filter)
}

func (s *UserRoleUseCase) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

func (s *UserRoleUseCase) ensureUser(ctx context.Context, userID string) error {
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrInvalidUser
		}
		return err
	}
	return nil
}

func (s *UserRoleUseCase) ensureUnique(ctx context.Context, userID, roleID string) error {
	_, err := s.repo.GetByUserAndRole(ctx, userID, roleID)
	if err == nil {
		return domain.ErrDuplicateUserRole
	}
	if err != domain.ErrNotFound {
		return err
	}
	return nil
}

func (s *UserRoleUseCase) ensureUniqueNumber(ctx context.Context, roleID, number string) error {
	if number == "" {
		return nil
	}
	_, err := s.repo.GetByRoleAndNumber(ctx, roleID, number)
	if err == nil {
		return domain.ErrDuplicateUserRoleNumber
	}
	if err != domain.ErrNotFound {
		return err
	}
	return nil
}

func normalizeUserRole(userID, roleID, number string, actor *string) (string, string, string, *string, error) {
	userID = strings.TrimSpace(userID)
	roleID = strings.TrimSpace(roleID)
	number = strings.TrimSpace(number)
	if !isUUID(userID) || !isUUID(roleID) {
		return "", "", "", nil, domain.ErrInvalidInput
	}
	if utf8.RuneCountInString(number) > 50 {
		return "", "", "", nil, domain.ErrInvalidInput
	}

	createdBy, err := optionalActor(actor)
	if err != nil {
		return "", "", "", nil, err
	}
	return userID, roleID, number, createdBy, nil
}

func requiresRoleNumber(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "dosen", "lecturer", "mahasiswa", "student":
		return true
	default:
		return false
	}
}
