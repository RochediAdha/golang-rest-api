package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"golang-rest-api/internal/domain"
)

type RolePrivilegeRepository struct {
	pool *pgxpool.Pool
}

func NewRolePrivilegeRepository(pool *pgxpool.Pool) *RolePrivilegeRepository {
	return &RolePrivilegeRepository{pool: pool}
}

func (s *RolePrivilegeRepository) Create(ctx context.Context, rolePrivilege domain.RolePrivilege) (domain.RolePrivilege, error) {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO role_privileges (id, "roleId", "menuId", "privilegeId", "createdAt", "createdBy")
		VALUES ($1, $2, $3, $4, $5, $6)
	`, rolePrivilege.ID, rolePrivilege.RoleID, rolePrivilege.MenuID, rolePrivilege.PrivilegeID, rolePrivilege.CreatedAt, nullUUID(rolePrivilege.CreatedBy))
	if err != nil {
		return domain.RolePrivilege{}, mapRolePrivilegeError(err)
	}
	return rolePrivilege, nil
}

func (s *RolePrivilegeRepository) GetByID(ctx context.Context, id string) (domain.RolePrivilege, error) {
	return scanRolePrivilege(s.pool.QueryRow(ctx, `
		SELECT id, "roleId", "menuId", "privilegeId", "createdAt", "createdBy"
		FROM role_privileges
		WHERE id = $1
	`, id))
}

func (s *RolePrivilegeRepository) GetByRoleMenuPrivilege(ctx context.Context, roleID, menuID, privilegeID string) (domain.RolePrivilege, error) {
	return scanRolePrivilege(s.pool.QueryRow(ctx, `
		SELECT id, "roleId", "menuId", "privilegeId", "createdAt", "createdBy"
		FROM role_privileges
		WHERE "roleId" = $1 AND "menuId" = $2 AND "privilegeId" = $3
	`, roleID, menuID, privilegeID))
}

func (s *RolePrivilegeRepository) List(ctx context.Context, filter domain.RolePrivilegeListFilter) ([]domain.RolePrivilege, int, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, "roleId", "menuId", "privilegeId", "createdAt", "createdBy", COUNT(*) OVER() AS total
		FROM role_privileges
		WHERE ($1 = '' OR "roleId"::text = $1)
		  AND ($2 = '' OR "menuId"::text = $2)
		  AND ($3 = '' OR "privilegeId"::text = $3)
		ORDER BY "createdAt" DESC
		LIMIT $4 OFFSET $5
	`, filter.RoleID, filter.MenuID, filter.PrivilegeID, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list role privileges: %w", err)
	}
	defer rows.Close()

	items := make([]domain.RolePrivilege, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanRolePrivilegeWithTotal(rows)
		if err != nil {
			return nil, 0, err
		}
		total = rowTotal
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list role privileges: %w", err)
	}
	return items, total, nil
}

func (s *RolePrivilegeRepository) Delete(ctx context.Context, id string) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM role_privileges WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete role privilege: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanRolePrivilege(row rowScanner) (domain.RolePrivilege, error) {
	item, _, err := scanRolePrivilegeRow(row, false)
	return item, err
}

func scanRolePrivilegeWithTotal(row rowScanner) (domain.RolePrivilege, int, error) {
	return scanRolePrivilegeRow(row, true)
}

func scanRolePrivilegeRow(row rowScanner, withTotal bool) (domain.RolePrivilege, int, error) {
	var (
		item      domain.RolePrivilege
		createdBy *string
		total     int
	)

	dest := []any{
		&item.ID,
		&item.RoleID,
		&item.MenuID,
		&item.PrivilegeID,
		&item.CreatedAt,
		&createdBy,
	}
	if withTotal {
		dest = append(dest, &total)
	}

	if err := row.Scan(dest...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RolePrivilege{}, 0, domain.ErrNotFound
		}
		return domain.RolePrivilege{}, 0, fmt.Errorf("scan role privilege: %w", err)
	}
	item.CreatedBy = createdBy
	return item, total, nil
}

func mapRolePrivilegeError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	name := strings.ToLower(pgErr.ConstraintName)
	switch pgErr.Code {
	case "23505":
		return domain.ErrDuplicateRolePrivilege
	case "23503":
		if strings.Contains(name, "roleid") {
			return domain.ErrInvalidRole
		}
		if strings.Contains(name, "menuid") {
			return domain.ErrInvalidMenu
		}
		if strings.Contains(name, "privilegeid") {
			return domain.ErrInvalidPrivilege
		}
	}
	return err
}

var _ domain.RolePrivilegeRepository = (*RolePrivilegeRepository)(nil)
