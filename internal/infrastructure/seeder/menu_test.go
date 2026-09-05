package seeder

import (
	"context"
	"testing"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

func TestMenusSeederIsIdempotent(t *testing.T) {
	ctx := context.Background()
	uc := usecase.NewMenuUseCase(memory.NewMenuRepository())

	if err := Menus(ctx, uc); err != nil {
		t.Fatalf("first seed: %v", err)
	}
	if err := Menus(ctx, uc); err != nil {
		t.Fatalf("second seed: %v", err)
	}

	menus, total, err := uc.List(ctx, domain.MenuListFilter{Limit: 20})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != len(defaultMenus) || len(menus) != len(defaultMenus) {
		t.Fatalf("want %d menus, got total=%d len=%d", len(defaultMenus), total, len(menus))
	}

	roots, rootTotal, err := uc.List(ctx, domain.MenuListFilter{RootOnly: true, Limit: 20})
	if err != nil {
		t.Fatalf("list roots: %v", err)
	}
	wantRoots := 0
	for _, sample := range defaultMenus {
		if sample.ParentCode == "" {
			wantRoots++
		}
	}
	if rootTotal != wantRoots || len(roots) != wantRoots {
		t.Fatalf("want %d root menus, got total=%d len=%d", wantRoots, rootTotal, len(roots))
	}
}
