package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"golang-rest-api/internal/domain"
)

type UserUseCase struct {
	repo domain.UserRepository
	now  func() time.Time
}

func NewUserUseCase(repo domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		now:  time.Now,
	}
}

func (s *UserUseCase) Create(ctx context.Context, input domain.CreateUserInput) (domain.User, error) {
	username, email, name, createdBy, err := normalizeUser(input.Username, input.Email, input.Name, input.CreatedBy)
	if err != nil {
		return domain.User{}, err
	}

	if err := s.ensureUnique(ctx, "", username, email); err != nil {
		return domain.User{}, err
	}

	now := s.now()
	user := domain.User{
		ID:        newUUID(),
		Username:  username,
		Email:     email,
		Name:      name,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	return s.repo.Create(ctx, user)
}

func (s *UserUseCase) Get(ctx context.Context, id string) (domain.User, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.User{}, domain.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *UserUseCase) List(ctx context.Context, filter domain.ListFilter) ([]domain.User, int, error) {
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

func (s *UserUseCase) Update(ctx context.Context, id string, input domain.UpdateUserInput) (domain.User, error) {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.User{}, domain.ErrInvalidInput
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	username, email, name, updatedBy, err := normalizeUser(input.Username, input.Email, input.Name, input.UpdatedBy)
	if err != nil {
		return domain.User{}, err
	}
	if err := s.ensureUnique(ctx, current.ID, username, email); err != nil {
		return domain.User{}, err
	}

	current.Username = username
	current.Email = email
	current.Name = name
	current.UpdatedAt = s.now()
	current.UpdatedBy = updatedBy
	if input.IsActive != nil {
		current.IsActive = *input.IsActive
	}

	return s.repo.Update(ctx, current)
}

func (s *UserUseCase) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isUUID(id) {
		return domain.ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

func (s *UserUseCase) ensureUnique(ctx context.Context, currentID, username, email string) error {
	if existing, err := s.repo.GetByUsername(ctx, username); err == nil && existing.ID != currentID {
		return domain.ErrDuplicateUsername
	} else if err != nil && err != domain.ErrNotFound {
		return err
	}

	if existing, err := s.repo.GetByEmail(ctx, email); err == nil && existing.ID != currentID {
		return domain.ErrDuplicateEmail
	} else if err != nil && err != domain.ErrNotFound {
		return err
	}

	return nil
}

func normalizeUser(username, email, name string, actor *string) (string, string, string, *string, error) {
	username = strings.TrimSpace(username)
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)

	if username == "" || email == "" || name == "" {
		return "", "", "", nil, domain.ErrInvalidInput
	}
	if strings.ContainsAny(username, " \t") || utf8.RuneCountInString(username) > 80 {
		return "", "", "", nil, domain.ErrInvalidInput
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") || utf8.RuneCountInString(email) > 254 {
		return "", "", "", nil, domain.ErrInvalidInput
	}
	if utf8.RuneCountInString(name) > 120 {
		return "", "", "", nil, domain.ErrInvalidInput
	}

	var actorID *string
	if actor != nil {
		value := strings.TrimSpace(*actor)
		if value != "" {
			if !isUUID(value) {
				return "", "", "", nil, domain.ErrInvalidInput
			}
			actorID = &value
		}
	}

	return username, email, name, actorID, nil
}

func newUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func isUUID(v string) bool {
	if len(v) != 36 {
		return false
	}
	for i, c := range v {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !isHex(c) {
				return false
			}
		}
	}
	return true
}

func isHex(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
