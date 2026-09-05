package seeder

import (
	"context"
	"testing"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

func TestPrivilegesSeederIsIdempotent(t *testing.T) {
	ctx := context.Background()
	uc := usecase.NewPrivilegeUseCase(memory.NewPrivilegeRepository())

	if err := Privileges(ctx, uc); err != nil {
		t.Fatalf("first seed: %v", err)
	}
	if err := Privileges(ctx, uc); err != nil {
		t.Fatalf("second seed: %v", err)
	}

	items, total, err := uc.List(ctx, domain.ListFilter{Limit: 20})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != len(defaultPrivileges) || len(items) != len(defaultPrivileges) {
		t.Fatalf("want %d privileges, got total=%d len=%d", len(defaultPrivileges), total, len(items))
	}
}
