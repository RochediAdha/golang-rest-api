package seeder

import (
	"context"
	"testing"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

func TestRolesSeederIsIdempotent(t *testing.T) {
	ctx := context.Background()
	uc := usecase.NewRoleUseCase(memory.NewRoleRepository())

	if err := Roles(ctx, uc); err != nil {
		t.Fatalf("first seed: %v", err)
	}
	if err := Roles(ctx, uc); err != nil {
		t.Fatalf("second seed: %v", err)
	}

	roles, total, err := uc.List(ctx, domain.ListFilter{Limit: 20})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != len(defaultRoles) || len(roles) != len(defaultRoles) {
		t.Fatalf("want %d roles, got total=%d len=%d", len(defaultRoles), total, len(roles))
	}
}
