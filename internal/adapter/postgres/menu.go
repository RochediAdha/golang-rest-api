package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"golang-rest-api/internal/domain"
)

type MenuRepository struct {
	pool *pgxpool.Pool
}

func NewMenuRepository(pool *pgxpool.Pool) *MenuRepository {
	return &MenuRepository{pool: pool}
}

func (s *MenuRepository) Create(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO menus (id, "parentId", code, name, path, icon, description, "sortOrder", type, "isActive", "createdAt", "updatedAt", "deletedAt", "createdBy", "updatedBy", "deletedBy")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`, menu.ID, nullUUID(menu.ParentID), menu.Code, menu.Name, nullString(menu.Path), nullString(menu.Icon), nullString(menu.Description), menu.SortOrder, menu.Type, menu.IsActive, menu.CreatedAt, menu.UpdatedAt, nullTime(menu.DeletedAt), nullUUID(menu.CreatedBy), nullUUID(menu.UpdatedBy), nullUUID(menu.DeletedBy))
	if err != nil {
		return domain.Menu{}, mapMenuError(err)
	}
	return menu, nil
}

func (s *MenuRepository) GetByID(ctx context.Context, id string) (domain.Menu, error) {
	return scanMenu(s.pool.QueryRow(ctx, `
		SELECT id, "parentId", code, name, path, icon, description, "sortOrder", type, "isActive", "createdAt", "updatedAt", "deletedAt", "createdBy", "updatedBy", "deletedBy"
		FROM menus
		WHERE id = $1 AND "deletedAt" IS NULL
	`, id))
}

func (s *MenuRepository) GetByCode(ctx context.Context, code string) (domain.Menu, error) {
	return scanMenu(s.pool.QueryRow(ctx, `
		SELECT id, "parentId", code, name, path, icon, description, "sortOrder", type, "isActive", "createdAt", "updatedAt", "deletedAt", "createdBy", "updatedBy", "deletedBy"
		FROM menus
		WHERE lower(code) = lower($1) AND "deletedAt" IS NULL
	`, code))
}

func (s *MenuRepository) List(ctx context.Context, filter domain.MenuListFilter) ([]domain.Menu, int, error) {
	var parent any
	if filter.ParentID != "" {
		parent = filter.ParentID
	}

	rows, err := s.pool.Query(ctx, `
		SELECT id, "parentId", code, name, path, icon, description, "sortOrder", type, "isActive", "createdAt", "updatedAt", "deletedAt", "createdBy", "updatedBy", "deletedBy", COUNT(*) OVER() AS total
		FROM menus
		WHERE "deletedAt" IS NULL
		  AND ($1 = '' OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%' OR path ILIKE '%' || $1 || '%')
		  AND (NOT $2 OR "parentId" IS NULL)
		  AND ($3::uuid IS NULL OR "parentId" = $3)
		ORDER BY "sortOrder" ASC, "createdAt" ASC
		LIMIT $4 OFFSET $5
	`, filter.Query, filter.RootOnly, parent, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list menus: %w", err)
	}
	defer rows.Close()

	menus := make([]domain.Menu, 0)
	total := 0
	for rows.Next() {
		menu, rowTotal, err := scanMenuWithTotal(rows)
		if err != nil {
			return nil, 0, err
		}
		total = rowTotal
		menus = append(menus, menu)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list menus: %w", err)
	}
	return menus, total, nil
}

func (s *MenuRepository) CountChildren(ctx context.Context, parentID string) (int, error) {
	var count int
	if err := s.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM menus
		WHERE "parentId" = $1 AND "deletedAt" IS NULL
	`, parentID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count menu children: %w", err)
	}
	return count, nil
}

func (s *MenuRepository) Update(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE menus
		SET "parentId" = $2, code = $3, name = $4, path = $5, icon = $6, description = $7, "sortOrder" = $8, type = $9, "isActive" = $10, "updatedAt" = $11, "updatedBy" = $12
		WHERE id = $1 AND "deletedAt" IS NULL
	`, menu.ID, nullUUID(menu.ParentID), menu.Code, menu.Name, nullString(menu.Path), nullString(menu.Icon), nullString(menu.Description), menu.SortOrder, menu.Type, menu.IsActive, menu.UpdatedAt, nullUUID(menu.UpdatedBy))
	if err != nil {
		return domain.Menu{}, mapMenuError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.Menu{}, domain.ErrNotFound
	}
	return menu, nil
}

func (s *MenuRepository) Delete(ctx context.Context, menu domain.Menu) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE menus
		SET "deletedAt" = $2, "deletedBy" = $3, "updatedAt" = $4, "updatedBy" = $5
		WHERE id = $1 AND "deletedAt" IS NULL
	`, menu.ID, nullTime(menu.DeletedAt), nullUUID(menu.DeletedBy), menu.UpdatedAt, nullUUID(menu.UpdatedBy))
	if err != nil {
		return mapMenuError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanMenu(row rowScanner) (domain.Menu, error) {
	menu, _, err := scanMenuRow(row, false)
	return menu, err
}

func scanMenuWithTotal(row rowScanner) (domain.Menu, int, error) {
	return scanMenuRow(row, true)
}

func scanMenuRow(row rowScanner, withTotal bool) (domain.Menu, int, error) {
	var (
		menu        domain.Menu
		parentID    *string
		path        *string
		icon        *string
		description *string
		deletedAt   *time.Time
		createdBy   *string
		updatedBy   *string
		deletedBy   *string
		total       int
	)

	dest := []any{
		&menu.ID,
		&parentID,
		&menu.Code,
		&menu.Name,
		&path,
		&icon,
		&description,
		&menu.SortOrder,
		&menu.Type,
		&menu.IsActive,
		&menu.CreatedAt,
		&menu.UpdatedAt,
		&deletedAt,
		&createdBy,
		&updatedBy,
		&deletedBy,
	}
	if withTotal {
		dest = append(dest, &total)
	}

	if err := row.Scan(dest...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Menu{}, 0, domain.ErrNotFound
		}
		return domain.Menu{}, 0, fmt.Errorf("scan menu: %w", err)
	}
	menu.ParentID = parentID
	if path != nil {
		menu.Path = *path
	}
	if icon != nil {
		menu.Icon = *icon
	}
	if description != nil {
		menu.Description = *description
	}
	menu.DeletedAt = deletedAt
	menu.CreatedBy = createdBy
	menu.UpdatedBy = updatedBy
	menu.DeletedBy = deletedBy
	return menu, total, nil
}

func mapMenuError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case "23505":
		if strings.Contains(strings.ToLower(pgErr.ConstraintName), "code") {
			return domain.ErrDuplicateMenuCode
		}
	case "23503":
		return domain.ErrInvalidParent
	case "22P02":
		return domain.ErrInvalidInput
	}
	return err
}

var _ domain.MenuRepository = (*MenuRepository)(nil)
