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

type UserRoleRepository struct {
	pool *pgxpool.Pool
}

func NewUserRoleRepository(pool *pgxpool.Pool) *UserRoleRepository {
	return &UserRoleRepository{pool: pool}
}

func (s *UserRoleRepository) Create(ctx context.Context, userRole domain.UserRole) (domain.UserRole, error) {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO user_roles (id, "userId", "roleId", number, "createdAt", "createdBy")
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userRole.ID, userRole.UserID, userRole.RoleID, nullString(userRole.Number), userRole.CreatedAt, nullUUID(userRole.CreatedBy))
	if err != nil {
		return domain.UserRole{}, mapUserRoleError(err)
	}
	return userRole, nil
}

func (s *UserRoleRepository) GetByID(ctx context.Context, id string) (domain.UserRole, error) {
	return scanUserRole(s.pool.QueryRow(ctx, `
		SELECT id, "userId", "roleId", number, "createdAt", "createdBy"
		FROM user_roles
		WHERE id = $1
	`, id))
}

func (s *UserRoleRepository) GetByUserAndRole(ctx context.Context, userID, roleID string) (domain.UserRole, error) {
	return scanUserRole(s.pool.QueryRow(ctx, `
		SELECT id, "userId", "roleId", number, "createdAt", "createdBy"
		FROM user_roles
		WHERE "userId" = $1 AND "roleId" = $2
	`, userID, roleID))
}

func (s *UserRoleRepository) GetByRoleAndNumber(ctx context.Context, roleID, number string) (domain.UserRole, error) {
	return scanUserRole(s.pool.QueryRow(ctx, `
		SELECT id, "userId", "roleId", number, "createdAt", "createdBy"
		FROM user_roles
		WHERE "roleId" = $1 AND number = $2
	`, roleID, number))
}

func (s *UserRoleRepository) List(ctx context.Context, filter domain.UserRoleListFilter) ([]domain.UserRole, int, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, "userId", "roleId", number, "createdAt", "createdBy", COUNT(*) OVER() AS total
		FROM user_roles
		WHERE ($1 = '' OR "userId"::text = $1)
		  AND ($2 = '' OR "roleId"::text = $2)
		  AND ($3 = '' OR number = $3)
		ORDER BY "createdAt" DESC
		LIMIT $4 OFFSET $5
	`, filter.UserID, filter.RoleID, filter.Number, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list user roles: %w", err)
	}
	defer rows.Close()

	items := make([]domain.UserRole, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanUserRoleWithTotal(rows)
		if err != nil {
			return nil, 0, err
		}
		total = rowTotal
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list user roles: %w", err)
	}
	return items, total, nil
}

func (s *UserRoleRepository) Delete(ctx context.Context, id string) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM user_roles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user role: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanUserRole(row rowScanner) (domain.UserRole, error) {
	item, _, err := scanUserRoleRow(row, false)
	return item, err
}

func scanUserRoleWithTotal(row rowScanner) (domain.UserRole, int, error) {
	return scanUserRoleRow(row, true)
}

func scanUserRoleRow(row rowScanner, withTotal bool) (domain.UserRole, int, error) {
	var (
		item      domain.UserRole
		number    *string
		createdBy *string
		total     int
	)

	dest := []any{
		&item.ID,
		&item.UserID,
		&item.RoleID,
		&number,
		&item.CreatedAt,
		&createdBy,
	}
	if withTotal {
		dest = append(dest, &total)
	}

	if err := row.Scan(dest...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserRole{}, 0, domain.ErrNotFound
		}
		return domain.UserRole{}, 0, fmt.Errorf("scan user role: %w", err)
	}
	if number != nil {
		item.Number = *number
	}
	item.CreatedBy = createdBy
	return item, total, nil
}

func mapUserRoleError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	name := strings.ToLower(pgErr.ConstraintName)
	switch pgErr.Code {
	case "23505":
		if strings.Contains(name, "number") {
			return domain.ErrDuplicateUserRoleNumber
		}
		return domain.ErrDuplicateUserRole
	case "23503":
		if strings.Contains(name, "userid") {
			return domain.ErrInvalidUser
		}
		if strings.Contains(name, "roleid") {
			return domain.ErrInvalidRole
		}
	}
	return err
}

var _ domain.UserRoleRepository = (*UserRoleRepository)(nil)
