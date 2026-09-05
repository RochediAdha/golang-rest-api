package usecase

import (
	"context"
	"testing"
	"time"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/domain"
)

func TestUserRoleUseCaseCRUD(t *testing.T) {
	ctx := context.Background()
	userRepo := memory.NewUserRepository()
	roleRepo := memory.NewRoleRepository()
	userUC := NewUserUseCase(userRepo)
	roleUC := NewRoleUseCase(roleRepo)
	uc := NewUserRoleUseCase(memory.NewUserRoleRepository(), userRepo, roleRepo)
	uc.now = func() time.Time { return time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC) }

	user, err := userUC.Create(ctx, domain.CreateUserInput{Username: "rochedi", Email: "rochedi@example.com", Name: "Rochedi"})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	role, err := roleUC.Create(ctx, domain.CreateRoleInput{Name: "admin"})
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	created, err := uc.Create(ctx, domain.CreateUserRoleInput{UserID: user.ID, RoleID: role.ID})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !isUUID(created.ID) || created.UserID != user.ID || created.RoleID != role.ID {
		t.Fatalf("unexpected created user role: %+v", created)
	}

	got, err := uc.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("get mismatch: %+v", got)
	}

	student, err := roleUC.Create(ctx, domain.CreateRoleInput{Name: "mahasiswa"})
	if err != nil {
		t.Fatalf("create student role: %v", err)
	}
	assigned, err := uc.Create(ctx, domain.CreateUserRoleInput{UserID: user.ID, RoleID: student.ID, Number: " 2301001 "})
	if err != nil {
		t.Fatalf("create student assignment: %v", err)
	}
	if assigned.Number != "2301001" {
		t.Fatalf("expected trimmed number, got %q", assigned.Number)
	}

	list, total, err := uc.List(ctx, domain.UserRoleListFilter{Number: "2301001", Limit: 10})
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

func TestUserRoleUseCaseValidation(t *testing.T) {
	ctx := context.Background()
	userRepo := memory.NewUserRepository()
	roleRepo := memory.NewRoleRepository()
	userUC := NewUserUseCase(userRepo)
	roleUC := NewRoleUseCase(roleRepo)
	uc := NewUserRoleUseCase(memory.NewUserRoleRepository(), userRepo, roleRepo)

	user, err := userUC.Create(ctx, domain.CreateUserInput{Username: "rochedi", Email: "rochedi@example.com", Name: "Rochedi"})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	role, err := roleUC.Create(ctx, domain.CreateRoleInput{Name: "admin"})
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	_, err = uc.Create(ctx, domain.CreateUserRoleInput{UserID: "", RoleID: role.ID})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}

	missing := "11111111-1111-4111-8111-111111111111"
	_, err = uc.Create(ctx, domain.CreateUserRoleInput{UserID: missing, RoleID: role.ID})
	if err != domain.ErrInvalidUser {
		t.Fatalf("expected invalid user, got %v", err)
	}
	_, err = uc.Create(ctx, domain.CreateUserRoleInput{UserID: user.ID, RoleID: missing})
	if err != domain.ErrInvalidRole {
		t.Fatalf("expected invalid role, got %v", err)
	}

	if _, err := uc.Create(ctx, domain.CreateUserRoleInput{UserID: user.ID, RoleID: role.ID}); err != nil {
		t.Fatalf("create first: %v", err)
	}
	_, err = uc.Create(ctx, domain.CreateUserRoleInput{UserID: user.ID, RoleID: role.ID})
	if err != domain.ErrDuplicateUserRole {
		t.Fatalf("expected duplicate user role, got %v", err)
	}

	student, err := roleUC.Create(ctx, domain.CreateRoleInput{Name: "Student"})
	if err != nil {
		t.Fatalf("create student role: %v", err)
	}
	_, err = uc.Create(ctx, domain.CreateUserRoleInput{UserID: user.ID, RoleID: student.ID})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected number required for student, got %v", err)
	}

	other, err := userUC.Create(ctx, domain.CreateUserInput{Username: "andi", Email: "andi@example.com", Name: "Andi"})
	if err != nil {
		t.Fatalf("create other user: %v", err)
	}
	if _, err := uc.Create(ctx, domain.CreateUserRoleInput{UserID: user.ID, RoleID: student.ID, Number: "2301001"}); err != nil {
		t.Fatalf("create student number: %v", err)
	}
	_, err = uc.Create(ctx, domain.CreateUserRoleInput{UserID: other.ID, RoleID: student.ID, Number: "2301001"})
	if err != domain.ErrDuplicateUserRoleNumber {
		t.Fatalf("expected duplicate number, got %v", err)
	}

	if _, err := uc.Get(ctx, "not-a-uuid"); err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}
