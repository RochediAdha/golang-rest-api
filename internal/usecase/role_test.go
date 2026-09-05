package usecase

import (
	"context"
	"testing"
	"time"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/domain"
)

func TestRoleUseCaseCRUD(t *testing.T) {
	ctx := context.Background()
	uc := NewRoleUseCase(memory.NewRoleRepository())
	uc.now = func() time.Time { return time.Date(2026, 9, 5, 6, 0, 0, 0, time.UTC) }

	created, err := uc.Create(ctx, domain.CreateRoleInput{
		Name:        "  admin  ",
		Description: "Full access",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !isUUID(created.ID) || created.Name != "admin" || created.Description != "Full access" || !created.IsActive {
		t.Fatalf("unexpected created role: %+v", created)
	}

	got, err := uc.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("get mismatch: %+v", got)
	}

	inactive := false
	updated, err := uc.Update(ctx, created.ID, domain.UpdateRoleInput{
		Name:        "admin",
		Description: "Administrator",
		IsActive:    &inactive,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Description != "Administrator" || updated.IsActive {
		t.Fatalf("unexpected update: %+v", updated)
	}

	list, total, err := uc.List(ctx, domain.ListFilter{Query: "admin", Limit: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("list want 1 got total=%d len=%d", total, len(list))
	}

	deleted, err := uc.Delete(ctx, created.ID, domain.DeleteRoleInput{})
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if deleted.DeletedAt == nil {
		t.Fatal("expected deletedAt after soft delete")
	}
	if _, err := uc.Get(ctx, created.ID); err != domain.ErrNotFound {
		t.Fatalf("expected not found after soft delete, got %v", err)
	}
}

func TestRoleUseCaseValidation(t *testing.T) {
	ctx := context.Background()
	uc := NewRoleUseCase(memory.NewRoleRepository())

	_, err := uc.Create(ctx, domain.CreateRoleInput{Name: ""})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}

	first, err := uc.Create(ctx, domain.CreateRoleInput{Name: "admin"})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}

	_, err = uc.Create(ctx, domain.CreateRoleInput{Name: "Admin"})
	if err != domain.ErrDuplicateRoleName {
		t.Fatalf("expected duplicate role name, got %v", err)
	}

	if _, err := uc.Get(ctx, "not-a-uuid"); err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
	_ = first
}
