package usecase

import (
	"context"
	"testing"
	"time"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/domain"
)

func TestPrivilegeUseCaseCRUD(t *testing.T) {
	ctx := context.Background()
	uc := NewPrivilegeUseCase(memory.NewPrivilegeRepository())
	uc.now = func() time.Time { return time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC) }

	created, err := uc.Create(ctx, domain.CreatePrivilegeInput{
		Code:        "  user.read  ",
		Name:        " Read User ",
		Description: "View user data",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !isUUID(created.ID) || created.Code != "user.read" || created.Name != "Read User" || !created.IsActive {
		t.Fatalf("unexpected created privilege: %+v", created)
	}

	got, err := uc.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("get mismatch: %+v", got)
	}

	inactive := false
	updated, err := uc.Update(ctx, created.ID, domain.UpdatePrivilegeInput{
		Code:        "user.read",
		Name:        "Read Users",
		Description: "View users",
		IsActive:    &inactive,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Read Users" || updated.IsActive {
		t.Fatalf("unexpected update: %+v", updated)
	}

	list, total, err := uc.List(ctx, domain.ListFilter{Query: "user.read", Limit: 10})
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
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestPrivilegeUseCaseValidation(t *testing.T) {
	ctx := context.Background()
	uc := NewPrivilegeUseCase(memory.NewPrivilegeRepository())

	_, err := uc.Create(ctx, domain.CreatePrivilegeInput{Code: "", Name: "Read"})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}

	_, err = uc.Create(ctx, domain.CreatePrivilegeInput{Code: "user read", Name: "Read"})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid code, got %v", err)
	}

	first, err := uc.Create(ctx, domain.CreatePrivilegeInput{Code: "user.read", Name: "Read"})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}

	_, err = uc.Create(ctx, domain.CreatePrivilegeInput{Code: "User.Read", Name: "Other"})
	if err != domain.ErrDuplicatePrivilegeCode {
		t.Fatalf("expected duplicate code, got %v", err)
	}

	if _, err := uc.Get(ctx, "not-a-uuid"); err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
	_ = first
}
