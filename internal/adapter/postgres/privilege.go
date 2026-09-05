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

type PrivilegeRepository struct {
	pool *pgxpool.Pool
}

func NewPrivilegeRepository(pool *pgxpool.Pool) *PrivilegeRepository {
	return &PrivilegeRepository{pool: pool}
}

func (s *PrivilegeRepository) Create(ctx context.Context, privilege domain.Privilege) (domain.Privilege, error) {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO privileges (id, code, name, description, "isActive", "createdAt", "updatedAt", "createdBy", "updatedBy")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, privilege.ID, privilege.Code, privilege.Name, nullString(privilege.Description), privilege.IsActive, privilege.CreatedAt, privilege.UpdatedAt, nullUUID(privilege.CreatedBy), nullUUID(privilege.UpdatedBy))
	if err != nil {
		return domain.Privilege{}, mapPrivilegeError(err)
	}
	return privilege, nil
}

func (s *PrivilegeRepository) GetByID(ctx context.Context, id string) (domain.Privilege, error) {
	return scanPrivilege(s.pool.QueryRow(ctx, `
		SELECT id, code, name, description, "isActive", "createdAt", "updatedAt", "createdBy", "updatedBy"
		FROM privileges
		WHERE id = $1
	`, id))
}

func (s *PrivilegeRepository) GetByCode(ctx context.Context, code string) (domain.Privilege, error) {
	return scanPrivilege(s.pool.QueryRow(ctx, `
		SELECT id, code, name, description, "isActive", "createdAt", "updatedAt", "createdBy", "updatedBy"
		FROM privileges
		WHERE lower(code) = lower($1)
	`, code))
}

func (s *PrivilegeRepository) List(ctx context.Context, filter domain.ListFilter) ([]domain.Privilege, int, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, code, name, description, "isActive", "createdAt", "updatedAt", "createdBy", "updatedBy", COUNT(*) OVER() AS total
		FROM privileges
		WHERE ($1 = '' OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
		ORDER BY "createdAt" DESC
		LIMIT $2 OFFSET $3
	`, filter.Query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list privileges: %w", err)
	}
	defer rows.Close()

	items := make([]domain.Privilege, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanPrivilegeWithTotal(rows)
		if err != nil {
			return nil, 0, err
		}
		total = rowTotal
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list privileges: %w", err)
	}
	return items, total, nil
}

func (s *PrivilegeRepository) Update(ctx context.Context, privilege domain.Privilege) (domain.Privilege, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE privileges
		SET code = $2, name = $3, description = $4, "isActive" = $5, "updatedAt" = $6, "updatedBy" = $7
		WHERE id = $1
	`, privilege.ID, privilege.Code, privilege.Name, nullString(privilege.Description), privilege.IsActive, privilege.UpdatedAt, nullUUID(privilege.UpdatedBy))
	if err != nil {
		return domain.Privilege{}, mapPrivilegeError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.Privilege{}, domain.ErrNotFound
	}
	return privilege, nil
}

func (s *PrivilegeRepository) Delete(ctx context.Context, id string) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM privileges WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete privilege: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanPrivilege(row rowScanner) (domain.Privilege, error) {
	item, _, err := scanPrivilegeRow(row, false)
	return item, err
}

func scanPrivilegeWithTotal(row rowScanner) (domain.Privilege, int, error) {
	return scanPrivilegeRow(row, true)
}

func scanPrivilegeRow(row rowScanner, withTotal bool) (domain.Privilege, int, error) {
	var (
		item        domain.Privilege
		description *string
		createdBy   *string
		updatedBy   *string
		total       int
	)

	dest := []any{
		&item.ID,
		&item.Code,
		&item.Name,
		&description,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
		&createdBy,
		&updatedBy,
	}
	if withTotal {
		dest = append(dest, &total)
	}

	if err := row.Scan(dest...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Privilege{}, 0, domain.ErrNotFound
		}
		return domain.Privilege{}, 0, fmt.Errorf("scan privilege: %w", err)
	}
	if description != nil {
		item.Description = *description
	}
	item.CreatedBy = createdBy
	item.UpdatedBy = updatedBy
	return item, total, nil
}

func mapPrivilegeError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(strings.ToLower(pgErr.ConstraintName), "code") {
		return domain.ErrDuplicatePrivilegeCode
	}
	return err
}

var _ domain.PrivilegeRepository = (*PrivilegeRepository)(nil)
