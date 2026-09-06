package seeder

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

type menuSeed struct {
	ParentCode  string
	Code        string
	Name        string
	Path        string
	Icon        string
	Description string
	SortOrder   int
	Type        string
}

var defaultMenus = []menuSeed{
	{Code: "dashboard", Name: "Dashboard", Path: "/dashboard", Icon: "home", Description: "Halaman utama", SortOrder: 1, Type: domain.MenuTypeItem},
	{Code: "master", Name: "Master Data", Path: "/master", Icon: "folder", Description: "Data referensi", SortOrder: 2, Type: domain.MenuTypeGroup},
	
	{Code: "user-management", Name: "User Management", Path: "/user-management", Icon: "users", Description: "Manajemen user", SortOrder: 98, Type: domain.MenuTypeGroup},
	{ParentCode: "user-management", Code: "users", Name: "Users", Path: "/master/users", Icon: "users", Description: "Manajemen user", SortOrder: 1, Type: domain.MenuTypeItem},
	{ParentCode: "user-management", Code: "roles", Name: "Roles", Path: "/master/roles", Icon: "shield", Description: "Manajemen role", SortOrder: 2, Type: domain.MenuTypeItem},
	{ParentCode: "user-management", Code: "permission", Name: "Permission", Path: "/master/permission", Icon: "shield", Description: "Manajemen permission", SortOrder: 3, Type: domain.MenuTypeItem},
	
	{Code: "system", Name: "System", Path: "/system", Icon: "cog", Description: "Pengaturan sistem", SortOrder: 99, Type: domain.MenuTypeGroup},
	{ParentCode: "system", Code: "menus", Name: "Menus", Path: "/menus", Icon: "menu", Description: "Manajemen menu", SortOrder: 1, Type: domain.MenuTypeItem},
}

func Menus(ctx context.Context, uc *usecase.MenuUseCase) error {
	ids := make(map[string]string, len(defaultMenus))

	for _, sample := range defaultMenus {
		input := domain.CreateMenuInput{
			Code:        sample.Code,
			Name:        sample.Name,
			Path:        sample.Path,
			Icon:        sample.Icon,
			Description: sample.Description,
			SortOrder:   intPtr(sample.SortOrder),
			Type:        sample.Type,
		}

		if sample.ParentCode != "" {
			parentID, ok := ids[sample.ParentCode]
			if !ok {
				parentID, ok = lookupMenuID(ctx, uc, sample.ParentCode)
			}
			if !ok {
				return fmt.Errorf("seed menu %s: parent %s not found", sample.Code, sample.ParentCode)
			}
			input.ParentID = &parentID
		}

		menu, err := uc.Create(ctx, input)
		if err != nil {
			if errors.Is(err, domain.ErrDuplicateMenuCode) {
				id, ok := lookupMenuID(ctx, uc, sample.Code)
				if !ok {
					return fmt.Errorf("seed menu %s: already exists but could not be loaded", sample.Code)
				}
				ids[sample.Code] = id
				continue
			}
			slog.Error("seed menu failed", "code", sample.Code, "err", err)
			return err
		}

		ids[sample.Code] = menu.ID
		slog.Info("seed menu created", "code", sample.Code)
	}

	return nil
}

func lookupMenuID(ctx context.Context, uc *usecase.MenuUseCase, code string) (string, bool) {
	menus, _, err := uc.List(ctx, domain.MenuListFilter{Limit: 100})
	if err != nil {
		return "", false
	}
	for _, menu := range menus {
		if strings.EqualFold(menu.Code, code) {
			return menu.ID, true
		}
	}
	return "", false
}

func intPtr(v int) *int {
	return &v
}
