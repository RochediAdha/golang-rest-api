package seeder

import (
	"context"
	"errors"
	"log/slog"

	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

var defaultRoles = []domain.CreateRoleInput{
	{Name: "Super Admin", Description: "Full access to all features"},
	{Name: "Admin", Description: "Can create and update content"},
	{Name: "Lecturer", Description: "Standard authenticated user"},
	{Name: "Student", Description: "Limited read-only access"},
}

func Roles(ctx context.Context, uc *usecase.RoleUseCase) error {
	for _, sample := range defaultRoles {
		if _, err := uc.Create(ctx, sample); err != nil {
			if errors.Is(err, domain.ErrDuplicateRoleName) {
				continue
			}
			slog.Error("seed role failed", "name", sample.Name, "err", err)
			return err
		}
		slog.Info("seed role created", "name", sample.Name)
	}
	return nil
}
