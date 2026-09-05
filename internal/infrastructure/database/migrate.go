package database

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const advisoryLockID int64 = 20260905

type Migration struct {
	Version int64
	Name    string
	SQL     string
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, "SELECT pg_advisory_lock($1)", advisoryLockID); err != nil {
		return fmt.Errorf("lock migrations: %w", err)
	}
	defer func() {
		if _, unlockErr := conn.Exec(ctx, "SELECT pg_advisory_unlock($1)", advisoryLockID); unlockErr != nil {
			slog.Error("unlock migrations", "err", unlockErr)
		}
	}()

	if _, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	for _, m := range migrations {
		var applied bool
		if err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, m.Version).Scan(&applied); err != nil {
			return fmt.Errorf("check migration %d: %w", m.Version, err)
		}
		if applied {
			continue
		}

		if _, err := conn.Exec(ctx, m.SQL); err != nil {
			return fmt.Errorf("apply %d_%s: %w", m.Version, m.Name, err)
		}
		if _, err := conn.Exec(ctx, `INSERT INTO schema_migrations (version, name) VALUES ($1, $2)`, m.Version, m.Name); err != nil {
			return fmt.Errorf("record %d_%s: %w", m.Version, m.Name, err)
		}
		slog.Info("migration applied", "version", m.Version, "name", m.Name)
	}

	return nil
}

func loadMigrations() ([]Migration, error) {
	entries, err := fs.ReadDir(migrationFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("read migrations: %w", err)
	}

	migrations := make([]Migration, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		m, err := parseMigration(entry.Name())
		if err != nil {
			return nil, err
		}

		sql, err := fs.ReadFile(migrationFS, path.Join("migrations", entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		m.SQL = string(sql)
		migrations = append(migrations, m)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	return migrations, nil
}

func parseMigration(filename string) (Migration, error) {
	base := strings.TrimSuffix(filename, ".sql")
	versionPart, name, ok := strings.Cut(base, "_")
	if !ok {
		return Migration{}, fmt.Errorf("invalid migration filename %q, want {version}_{name}.sql", filename)
	}
	version, err := strconv.ParseInt(versionPart, 10, 64)
	if err != nil {
		return Migration{}, fmt.Errorf("invalid migration version in %q: %w", filename, err)
	}
	if name == "" {
		return Migration{}, fmt.Errorf("invalid migration name in %q", filename)
	}
	return Migration{Version: version, Name: name}, nil
}
