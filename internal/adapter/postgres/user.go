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

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (s *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO users (id, username, email, name, "isActive", "createdAt", "updatedAt", "createdBy", "updatedBy")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, user.ID, user.Username, user.Email, user.Name, user.IsActive, user.CreatedAt, user.UpdatedAt, nullUUID(user.CreatedBy), nullUUID(user.UpdatedBy))
	if err != nil {
		return domain.User{}, mapUserError(err)
	}
	return user, nil
}

func (s *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	return scanUser(s.pool.QueryRow(ctx, `
		SELECT id, username, email, name, "isActive", "createdAt", "updatedAt", "createdBy", "updatedBy"
		FROM users
		WHERE id = $1
	`, id))
}

func (s *UserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	return scanUser(s.pool.QueryRow(ctx, `
		SELECT id, username, email, name, "isActive", "createdAt", "updatedAt", "createdBy", "updatedBy"
		FROM users
		WHERE lower(username) = lower($1)
	`, username))
}

func (s *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return scanUser(s.pool.QueryRow(ctx, `
		SELECT id, username, email, name, "isActive", "createdAt", "updatedAt", "createdBy", "updatedBy"
		FROM users
		WHERE lower(email) = lower($1)
	`, email))
}

func (s *UserRepository) List(ctx context.Context, filter domain.ListFilter) ([]domain.User, int, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, username, email, name, "isActive", "createdAt", "updatedAt", "createdBy", "updatedBy", COUNT(*) OVER() AS total
		FROM users
		WHERE ($1 = '' OR username ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
		ORDER BY "createdAt" DESC
		LIMIT $2 OFFSET $3
	`, filter.Query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	total := 0
	for rows.Next() {
		user, rowTotal, err := scanUserWithTotal(rows)
		if err != nil {
			return nil, 0, err
		}
		total = rowTotal
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	return users, total, nil
}

func (s *UserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE users
		SET username = $2, email = $3, name = $4, "isActive" = $5, "updatedAt" = $6, "updatedBy" = $7
		WHERE id = $1
	`, user.ID, user.Username, user.Email, user.Name, user.IsActive, user.UpdatedAt, nullUUID(user.UpdatedBy))
	if err != nil {
		return domain.User{}, mapUserError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

func (s *UserRepository) Delete(ctx context.Context, id string) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanUser(row rowScanner) (domain.User, error) {
	user, _, err := scanUserRow(row, false)
	return user, err
}

func scanUserWithTotal(row rowScanner) (domain.User, int, error) {
	return scanUserRow(row, true)
}

func scanUserRow(row rowScanner, withTotal bool) (domain.User, int, error) {
	var (
		user      domain.User
		createdBy *string
		updatedBy *string
		total     int
	)

	dest := []any{
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&createdBy,
		&updatedBy,
	}
	if withTotal {
		dest = append(dest, &total)
	}

	if err := row.Scan(dest...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, 0, domain.ErrNotFound
		}
		return domain.User{}, 0, fmt.Errorf("scan user: %w", err)
	}
	user.CreatedBy = createdBy
	user.UpdatedBy = updatedBy
	return user, total, nil
}

func mapUserError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		if strings.Contains(strings.ToLower(pgErr.ConstraintName), "username") {
			return domain.ErrDuplicateUsername
		}
		if strings.Contains(strings.ToLower(pgErr.ConstraintName), "email") {
			return domain.ErrDuplicateEmail
		}
	}
	return err
}

func nullUUID(v *string) any {
	if v == nil || *v == "" {
		return nil
	}
	return *v
}

var _ domain.UserRepository = (*UserRepository)(nil)
