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

type RoleRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{pool: pool}
}

func (s *RoleRepository) Create(ctx context.Context, role domain.Role) (domain.Role, error) {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO roles (id, name, description, "isActive", "createdAt", "updatedAt", "deletedAt", "createdBy", "updatedBy", "deletedBy")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, role.ID, role.Name, nullString(role.Description), role.IsActive, role.CreatedAt, role.UpdatedAt, nullTime(role.DeletedAt), nullUUID(role.CreatedBy), nullUUID(role.UpdatedBy), nullUUID(role.DeletedBy))
	if err != nil {
		return domain.Role{}, mapRoleError(err)
	}
	return role, nil
}

func (s *RoleRepository) GetByID(ctx context.Context, id string) (domain.Role, error) {
	return scanRole(s.pool.QueryRow(ctx, `
		SELECT id, name, description, "isActive", "createdAt", "updatedAt", "deletedAt", "createdBy", "updatedBy", "deletedBy"
		FROM roles
		WHERE id = $1 AND "deletedAt" IS NULL
	`, id))
}

func (s *RoleRepository) GetByName(ctx context.Context, name string) (domain.Role, error) {
	return scanRole(s.pool.QueryRow(ctx, `
		SELECT id, name, description, "isActive", "createdAt", "updatedAt", "deletedAt", "createdBy", "updatedBy", "deletedBy"
		FROM roles
		WHERE lower(name) = lower($1) AND "deletedAt" IS NULL
	`, name))
}

func (s *RoleRepository) List(ctx context.Context, filter domain.ListFilter) ([]domain.Role, int, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, description, "isActive", "createdAt", "updatedAt", "deletedAt", "createdBy", "updatedBy", "deletedBy", COUNT(*) OVER() AS total
		FROM roles
		WHERE "deletedAt" IS NULL
		  AND ($1 = '' OR name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
		ORDER BY "createdAt" DESC
		LIMIT $2 OFFSET $3
	`, filter.Query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list roles: %w", err)
	}
	defer rows.Close()

	roles := make([]domain.Role, 0)
	total := 0
	for rows.Next() {
		role, rowTotal, err := scanRoleWithTotal(rows)
		if err != nil {
			return nil, 0, err
		}
		total = rowTotal
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list roles: %w", err)
	}
	return roles, total, nil
}

func (s *RoleRepository) Update(ctx context.Context, role domain.Role) (domain.Role, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE roles
		SET name = $2, description = $3, "isActive" = $4, "updatedAt" = $5, "updatedBy" = $6
		WHERE id = $1 AND "deletedAt" IS NULL
	`, role.ID, role.Name, nullString(role.Description), role.IsActive, role.UpdatedAt, nullUUID(role.UpdatedBy))
	if err != nil {
		return domain.Role{}, mapRoleError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.Role{}, domain.ErrNotFound
	}
	return role, nil
}

func (s *RoleRepository) Delete(ctx context.Context, role domain.Role) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE roles
		SET "deletedAt" = $2, "deletedBy" = $3, "updatedAt" = $4, "updatedBy" = $5
		WHERE id = $1 AND "deletedAt" IS NULL
	`, role.ID, nullTime(role.DeletedAt), nullUUID(role.DeletedBy), role.UpdatedAt, nullUUID(role.UpdatedBy))
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanRole(row rowScanner) (domain.Role, error) {
	role, _, err := scanRoleRow(row, false)
	return role, err
}

func scanRoleWithTotal(row rowScanner) (domain.Role, int, error) {
	return scanRoleRow(row, true)
}

func scanRoleRow(row rowScanner, withTotal bool) (domain.Role, int, error) {
	var (
		role        domain.Role
		description *string
		deletedAt   *time.Time
		createdBy   *string
		updatedBy   *string
		deletedBy   *string
		total       int
	)

	dest := []any{
		&role.ID,
		&role.Name,
		&description,
		&role.IsActive,
		&role.CreatedAt,
		&role.UpdatedAt,
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
			return domain.Role{}, 0, domain.ErrNotFound
		}
		return domain.Role{}, 0, fmt.Errorf("scan role: %w", err)
	}
	if description != nil {
		role.Description = *description
	}
	role.DeletedAt = deletedAt
	role.CreatedBy = createdBy
	role.UpdatedBy = updatedBy
	role.DeletedBy = deletedBy
	return role, total, nil
}

func mapRoleError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(strings.ToLower(pgErr.ConstraintName), "name") {
		return domain.ErrDuplicateRoleName
	}
	return err
}

func nullTime(v *time.Time) any {
	if v == nil {
		return nil
	}
	return *v
}

var _ domain.RoleRepository = (*RoleRepository)(nil)
