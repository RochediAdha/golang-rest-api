package usecase

import (
	"context"
	"testing"
	"time"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/domain"
)

func TestRolePrivilegeUseCaseCRUD(t *testing.T) {
	ctx := context.Background()
	roleRepo := memory.NewRoleRepository()
	menuRepo := memory.NewMenuRepository()
	privilegeRepo := memory.NewPrivilegeRepository()
	roleUC := NewRoleUseCase(roleRepo)
	menuUC := NewMenuUseCase(menuRepo)
	privilegeUC := NewPrivilegeUseCase(privilegeRepo)
	uc := NewRolePrivilegeUseCase(memory.NewRolePrivilegeRepository(), roleRepo, menuRepo, privilegeRepo)
	uc.now = func() time.Time { return time.Date(2026, 9, 5, 13, 0, 0, 0, time.UTC) }

	role, err := roleUC.Create(ctx, domain.CreateRoleInput{Name: "admin"})
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	menu, err := menuUC.Create(ctx, domain.CreateMenuInput{Code: "users", Name: "Users"})
	if err != nil {
		t.Fatalf("create menu: %v", err)
	}
	privilege, err := privilegeUC.Create(ctx, domain.CreatePrivilegeInput{Code: "VIEW", Name: "View"})
	if err != nil {
		t.Fatalf("create privilege: %v", err)
	}

	created, err := uc.Create(ctx, domain.CreateRolePrivilegeInput{
		RoleID:      role.ID,
		MenuID:      menu.ID,
		PrivilegeID: privilege.ID,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !isUUID(created.ID) || created.RoleID != role.ID || created.MenuID != menu.ID || created.PrivilegeID != privilege.ID {
		t.Fatalf("unexpected created role privilege: %+v", created)
	}

	got, err := uc.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("get mismatch: %+v", got)
	}

	list, total, err := uc.List(ctx, domain.RolePrivilegeListFilter{RoleID: role.ID, Limit: 10})
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

func TestRolePrivilegeUseCaseValidation(t *testing.T) {
	ctx := context.Background()
	roleRepo := memory.NewRoleRepository()
	menuRepo := memory.NewMenuRepository()
	privilegeRepo := memory.NewPrivilegeRepository()
	roleUC := NewRoleUseCase(roleRepo)
	menuUC := NewMenuUseCase(menuRepo)
	privilegeUC := NewPrivilegeUseCase(privilegeRepo)
	uc := NewRolePrivilegeUseCase(memory.NewRolePrivilegeRepository(), roleRepo, menuRepo, privilegeRepo)

	role, err := roleUC.Create(ctx, domain.CreateRoleInput{Name: "admin"})
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	menu, err := menuUC.Create(ctx, domain.CreateMenuInput{Code: "users", Name: "Users"})
	if err != nil {
		t.Fatalf("create menu: %v", err)
	}
	privilege, err := privilegeUC.Create(ctx, domain.CreatePrivilegeInput{Code: "VIEW", Name: "View"})
	if err != nil {
		t.Fatalf("create privilege: %v", err)
	}

	_, err = uc.Create(ctx, domain.CreateRolePrivilegeInput{RoleID: "", MenuID: menu.ID, PrivilegeID: privilege.ID})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}

	missing := "11111111-1111-4111-8111-111111111111"
	_, err = uc.Create(ctx, domain.CreateRolePrivilegeInput{RoleID: missing, MenuID: menu.ID, PrivilegeID: privilege.ID})
	if err != domain.ErrInvalidRole {
		t.Fatalf("expected invalid role, got %v", err)
	}
	_, err = uc.Create(ctx, domain.CreateRolePrivilegeInput{RoleID: role.ID, MenuID: missing, PrivilegeID: privilege.ID})
	if err != domain.ErrInvalidMenu {
		t.Fatalf("expected invalid menu, got %v", err)
	}
	_, err = uc.Create(ctx, domain.CreateRolePrivilegeInput{RoleID: role.ID, MenuID: menu.ID, PrivilegeID: missing})
	if err != domain.ErrInvalidPrivilege {
		t.Fatalf("expected invalid privilege, got %v", err)
	}

	if _, err := uc.Create(ctx, domain.CreateRolePrivilegeInput{RoleID: role.ID, MenuID: menu.ID, PrivilegeID: privilege.ID}); err != nil {
		t.Fatalf("create first: %v", err)
	}
	_, err = uc.Create(ctx, domain.CreateRolePrivilegeInput{RoleID: role.ID, MenuID: menu.ID, PrivilegeID: privilege.ID})
	if err != domain.ErrDuplicateRolePrivilege {
		t.Fatalf("expected duplicate role privilege, got %v", err)
	}

	if _, err := uc.Get(ctx, "not-a-uuid"); err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}
