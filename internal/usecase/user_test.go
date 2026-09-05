package usecase

import (
	"context"
	"testing"
	"time"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/domain"
)

func TestUserUseCaseCRUD(t *testing.T) {
	ctx := context.Background()
	uc := NewUserUseCase(memory.NewUserRepository())
	uc.now = func() time.Time { return time.Date(2026, 9, 5, 4, 0, 0, 0, time.UTC) }

	created, err := uc.Create(ctx, domain.CreateUserInput{
		Username: "  rochedi  ",
		Email:    "Rochedi@Example.com",
		Name:     "Rochedi Adha",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !isUUID(created.ID) || created.Username != "rochedi" || created.Email != "rochedi@example.com" || !created.IsActive {
		t.Fatalf("unexpected created user: %+v", created)
	}

	got, err := uc.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("get mismatch: %+v", got)
	}

	inactive := false
	updated, err := uc.Update(ctx, created.ID, domain.UpdateUserInput{
		Username: "rochedi",
		Email:    "rochedi@example.com",
		Name:     "Rochedi A.",
		IsActive: &inactive,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Rochedi A." || updated.IsActive {
		t.Fatalf("unexpected update: %+v", updated)
	}

	list, total, err := uc.List(ctx, domain.ListFilter{Query: "rochedi", Limit: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("list want 1 got total=%d len=%d", total, len(list))
	}

	if err := uc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := uc.Get(ctx, created.ID); err != domain.ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestUserUseCaseValidation(t *testing.T) {
	ctx := context.Background()
	uc := NewUserUseCase(memory.NewUserRepository())

	_, err := uc.Create(ctx, domain.CreateUserInput{Username: "", Email: "a@b.com", Name: "A"})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}

	first, err := uc.Create(ctx, domain.CreateUserInput{Username: "admin", Email: "admin@example.com", Name: "Admin"})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}

	_, err = uc.Create(ctx, domain.CreateUserInput{Username: "admin", Email: "other@example.com", Name: "Other"})
	if err != domain.ErrDuplicateUsername {
		t.Fatalf("expected duplicate username, got %v", err)
	}

	_, err = uc.Create(ctx, domain.CreateUserInput{Username: "other", Email: "admin@example.com", Name: "Other"})
	if err != domain.ErrDuplicateEmail {
		t.Fatalf("expected duplicate email, got %v", err)
	}

	if _, err := uc.Get(ctx, "not-a-uuid"); err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
	_ = first
}
