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

func (s *UserRoleUseCase) Get(ctx context.Context, id string) (domain.UserRoleView, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.UserRoleView{}, domain.ErrInvalidInput
	}

	user, err := s.users.GetByID(ctx, id)
	if err == nil {
		return s.viewByUser(ctx, user)
	}
	if err != domain.ErrNotFound {
		return domain.UserRoleView{}, err
	}

	assignment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.UserRoleView{}, err
	}
	user, err = s.users.GetByID(ctx, assignment.UserID)
	if err != nil {
		return domain.UserRoleView{}, err
	}
	return s.viewByUser(ctx, user)
}

func (s *UserRoleUseCase) viewByUser(ctx context.Context, user domain.User) (domain.UserRoleView, error) {
	items, err := s.repo.ListByUserID(ctx, user.ID)
	if err != nil {
		return domain.UserRoleView{}, err
	}
	if len(items) == 0 {
		return domain.UserRoleView{}, domain.ErrNotFound
	}

	roles := make([]domain.UserRoleItem, 0, len(items))
	for _, item := range items {
		role, err := s.roles.GetByID(ctx, item.RoleID)
		if err != nil {
			if err == domain.ErrNotFound {
				continue
			}
			return domain.UserRoleView{}, err
		}
		roles = append(roles, domain.UserRoleItem{
			ID:        item.ID,
			RoleID:    item.RoleID,
			Name:      role.Name,
			Number:    item.Number,
			CreatedAt: item.CreatedAt,
			CreatedBy: item.CreatedBy,
		})
	}
	if len(roles) == 0 {
		return domain.UserRoleView{}, domain.ErrNotFound
	}

	return domain.UserRoleView{
		ID:     user.ID,
		UserID: user.ID,
		Name:   user.Name,
		Roles:  roles,
	}, nil
}

func (s *UserRoleUseCase) List(ctx context.Context, filter domain.UserRoleListFilter) ([]domain.UserRoleListItem, int, error) {
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

	items, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	result := make([]domain.UserRoleListItem, 0, len(items))
	for _, item := range items {
		row, err := s.toListItem(ctx, item)
		if err != nil {
			if err == domain.ErrNotFound {
				continue
			}
			return nil, 0, err
		}
		result = append(result, row)
	}
	return result, total, nil
}

func (s *UserRoleUseCase) toListItem(ctx context.Context, item domain.UserRole) (domain.UserRoleListItem, error) {
	user, err := s.users.GetByID(ctx, item.UserID)
	if err != nil {
		return domain.UserRoleListItem{}, err
	}
	role, err := s.roles.GetByID(ctx, item.RoleID)
	if err != nil {
		return domain.UserRoleListItem{}, err
	}

	row := domain.UserRoleListItem{
		ID:              item.ID,
		UserID:          item.UserID,
		UserName:        user.Name,
		RoleID:          item.RoleID,
		RoleName:        role.Name,
		RoleDescription: role.Description,
		Number:          item.Number,
		CreatedAt:       item.CreatedAt,
	}
	if item.CreatedBy != nil {
		actor := domain.ActorRef{ID: *item.CreatedBy}
		if creator, err := s.users.GetByID(ctx, *item.CreatedBy); err == nil {
			actor.Name = creator.Name
		}
		row.CreatedBy = &actor
	}
	return row, nil
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
