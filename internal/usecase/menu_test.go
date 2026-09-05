package usecase

import (
	"context"
	"testing"
	"time"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/domain"
)

func TestMenuUseCaseCRUD(t *testing.T) {
	ctx := context.Background()
	uc := NewMenuUseCase(memory.NewMenuRepository())
	uc.now = func() time.Time { return time.Date(2026, 9, 5, 8, 0, 0, 0, time.UTC) }

	parent, err := uc.Create(ctx, domain.CreateMenuInput{
		Code: "  dashboard  ",
		Name: "Dashboard",
		Path: "/dashboard",
		Icon: "home",
	})
	if err != nil {
		t.Fatalf("create parent: %v", err)
	}
	if parent.Code != "dashboard" || parent.ParentID != nil || !parent.IsActive || parent.Type != domain.MenuTypeItem {
		t.Fatalf("unexpected parent: %+v", parent)
	}

	child, err := uc.Create(ctx, domain.CreateMenuInput{
		ParentID:  &parent.ID,
		Code:      "users",
		Name:      "Users",
		Path:      "/users",
		SortOrder: intPtr(1),
		Type:      domain.MenuTypeItem,
	})
	if err != nil {
		t.Fatalf("create child: %v", err)
	}
	if child.ParentID == nil || *child.ParentID != parent.ID || child.SortOrder != 1 {
		t.Fatalf("unexpected child: %+v", child)
	}

	list, total, err := uc.List(ctx, domain.MenuListFilter{ParentID: parent.ID, Limit: 10})
	if err != nil {
		t.Fatalf("list children: %v", err)
	}
	if total != 1 || list[0].ID != child.ID {
		t.Fatalf("list children = %+v total=%d", list, total)
	}

	updated, err := uc.Update(ctx, child.ID, domain.UpdateMenuInput{
		ParentID: &parent.ID,
		Code:     "users",
		Name:     "User Management",
		Path:     "/users",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "User Management" {
		t.Fatalf("unexpected update: %+v", updated)
	}

	if _, err := uc.Delete(ctx, parent.ID, domain.DeleteMenuInput{}); err != domain.ErrMenuHasChildren {
		t.Fatalf("expected children conflict, got %v", err)
	}

	deleted, err := uc.Delete(ctx, child.ID, domain.DeleteMenuInput{})
	if err != nil {
		t.Fatalf("delete child: %v", err)
	}
	if deleted.DeletedAt == nil {
		t.Fatal("expected deletedAt")
	}

	if _, err := uc.Delete(ctx, parent.ID, domain.DeleteMenuInput{}); err != nil {
		t.Fatalf("delete parent: %v", err)
	}
}

func TestMenuUseCaseValidation(t *testing.T) {
	ctx := context.Background()
	uc := NewMenuUseCase(memory.NewMenuRepository())

	_, err := uc.Create(ctx, domain.CreateMenuInput{Code: "", Name: "A"})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}

	_, err = uc.Create(ctx, domain.CreateMenuInput{Code: "settings", Name: "Settings", Type: "LINK"})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid type, got %v", err)
	}

	first, err := uc.Create(ctx, domain.CreateMenuInput{Code: "settings", Name: "Settings", Type: domain.MenuTypeGroup})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}

	_, err = uc.Create(ctx, domain.CreateMenuInput{Code: "Settings", Name: "Other"})
	if err != domain.ErrDuplicateMenuCode {
		t.Fatalf("expected duplicate code, got %v", err)
	}

	missing := "11111111-1111-4111-8111-111111111111"
	_, err = uc.Create(ctx, domain.CreateMenuInput{ParentID: &missing, Code: "child", Name: "Child"})
	if err != domain.ErrInvalidParent {
		t.Fatalf("expected invalid parent, got %v", err)
	}

	_, err = uc.Update(ctx, first.ID, domain.UpdateMenuInput{ParentID: &first.ID, Code: "settings", Name: "Settings"})
	if err != domain.ErrInvalidParent {
		t.Fatalf("expected self parent invalid, got %v", err)
	}
}

func intPtr(v int) *int { return &v }
