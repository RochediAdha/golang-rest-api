package seeder

import (
	"context"
	"errors"
	"log/slog"

	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

var defaultPrivileges = []domain.CreatePrivilegeInput{
	{Code: "CREATE", Name: "Create", Description: "Create new records"},
	{Code: "EXECUTE", Name: "Execute", Description: "Run a process or action"},
	{Code: "APPROVE", Name: "Approve", Description: "Approve a pending request"},
	{Code: "VIEW", Name: "View", Description: "View data and details"},
	{Code: "PRINT", Name: "Print", Description: "Print documents or reports"},
	{Code: "PUBLISH", Name: "Publish", Description: "Publish content to be visible"},
	{Code: "DELETE", Name: "Delete", Description: "Delete existing records"},
	{Code: "EXPORT", Name: "Export", Description: "Export data to a file"},
	{Code: "UPDATE", Name: "Update", Description: "Update existing records"},
	{Code: "REJECT", Name: "Reject", Description: "Reject a pending request"},
	{Code: "IMPORT", Name: "Import", Description: "Import data from a file"},
}

func Privileges(ctx context.Context, uc *usecase.PrivilegeUseCase) error {
	for _, sample := range defaultPrivileges {
		if _, err := uc.Create(ctx, sample); err != nil {
			if errors.Is(err, domain.ErrDuplicatePrivilegeCode) {
				continue
			}
			slog.Error("seed privilege failed", "code", sample.Code, "err", err)
			return err
		}
		slog.Info("seed privilege created", "code", sample.Code)
	}
	return nil
}
