package usecase

import (
	"context"
	"strings"
	"time"

	"golang-rest-api/internal/domain"
)

type RolePrivilegeUseCase struct {
	repo       domain.RolePrivilegeRepository
	roles      domain.RoleRepository
	menus      domain.MenuRepository
	privileges domain.PrivilegeRepository
	now        func() time.Time
}

func NewRolePrivilegeUseCase(repo domain.RolePrivilegeRepository, roles domain.RoleRepository, menus domain.MenuRepository, privileges domain.PrivilegeRepository) *RolePrivilegeUseCase {
	return &RolePrivilegeUseCase{
		repo:       repo,
		roles:      roles,
		menus:      menus,
		privileges: privileges,
		now:        time.Now,
	}
}

func (s *RolePrivilegeUseCase) Create(ctx context.Context, input domain.CreateRolePrivilegeInput) (domain.RolePrivilege, error) {
	roleID, menuID, privilegeID, createdBy, err := normalizeRolePrivilege(input.RoleID, input.MenuID, input.PrivilegeID, input.CreatedBy)
	if err != nil {
		return domain.RolePrivilege{}, err
	}
	if err := s.ensureRole(ctx, roleID); err != nil {
		return domain.RolePrivilege{}, err
	}
	if err := s.ensureMenu(ctx, menuID); err != nil {
		return domain.RolePrivilege{}, err
	}
	if err := s.ensurePrivilege(ctx, privilegeID); err != nil {
		return domain.RolePrivilege{}, err
	}
	if err := s.ensureUnique(ctx, roleID, menuID, privilegeID); err != nil {
		return domain.RolePrivilege{}, err
	}

	return s.repo.Create(ctx, domain.RolePrivilege{
		ID:          newUUID(),
		RoleID:      roleID,
		MenuID:      menuID,
		PrivilegeID: privilegeID,
		CreatedAt:   s.now(),
		CreatedBy:   createdBy,
	})
}

func (s *RolePrivilegeUseCase) Get(ctx context.Context, id string) (domain.RolePrivilege, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.RolePrivilege{}, domain.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *RolePrivilegeUseCase) List(ctx context.Context, filter domain.RolePrivilegeListFilter) ([]domain.RolePrivilege, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	if filter.Limit > maxLimit {
		filter.Limit = maxLimit
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	filter.RoleID = strings.TrimSpace(filter.RoleID)
	filter.MenuID = strings.TrimSpace(filter.MenuID)
	filter.PrivilegeID = strings.TrimSpace(filter.PrivilegeID)
	if filter.RoleID != "" && !isUUID(filter.RoleID) {
		return nil, 0, domain.ErrInvalidInput
	}
	if filter.MenuID != "" && !isUUID(filter.MenuID) {
		return nil, 0, domain.ErrInvalidInput
	}
	if filter.PrivilegeID != "" && !isUUID(filter.PrivilegeID) {
		return nil, 0, domain.ErrInvalidInput
	}

	return s.repo.List(ctx, filter)
}

func (s *RolePrivilegeUseCase) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

func (s *RolePrivilegeUseCase) ensureRole(ctx context.Context, roleID string) error {
	if _, err := s.roles.GetByID(ctx, roleID); err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrInvalidRole
		}
		return err
	}
	return nil
}

func (s *RolePrivilegeUseCase) ensureMenu(ctx context.Context, menuID string) error {
	if _, err := s.menus.GetByID(ctx, menuID); err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrInvalidMenu
		}
		return err
	}
	return nil
}

func (s *RolePrivilegeUseCase) ensurePrivilege(ctx context.Context, privilegeID string) error {
	if _, err := s.privileges.GetByID(ctx, privilegeID); err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrInvalidPrivilege
		}
		return err
	}
	return nil
}

func (s *RolePrivilegeUseCase) ensureUnique(ctx context.Context, roleID, menuID, privilegeID string) error {
	_, err := s.repo.GetByRoleMenuPrivilege(ctx, roleID, menuID, privilegeID)
	if err == nil {
		return domain.ErrDuplicateRolePrivilege
	}
	if err != domain.ErrNotFound {
		return err
	}
	return nil
}

func normalizeRolePrivilege(roleID, menuID, privilegeID string, actor *string) (string, string, string, *string, error) {
	roleID = strings.TrimSpace(roleID)
	menuID = strings.TrimSpace(menuID)
	privilegeID = strings.TrimSpace(privilegeID)
	if !isUUID(roleID) || !isUUID(menuID) || !isUUID(privilegeID) {
		return "", "", "", nil, domain.ErrInvalidInput
	}

	createdBy, err := optionalActor(actor)
	if err != nil {
		return "", "", "", nil, err
	}
	return roleID, menuID, privilegeID, createdBy, nil
}
