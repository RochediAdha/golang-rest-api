package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"golang-rest-api/internal/adapter/postgres"
	"golang-rest-api/internal/infrastructure/config"
	"golang-rest-api/internal/infrastructure/database"
	"golang-rest-api/internal/infrastructure/seeder"
	"golang-rest-api/internal/usecase"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/seed <seeder>")
		fmt.Fprintln(os.Stderr, "available: roles, books, menus")
		os.Exit(1)
	}

	name := os.Args[1]
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := database.Open(ctx, cfg.Database)
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := database.Migrate(ctx, pool); err != nil {
		slog.Error("database migration failed", "err", err)
		os.Exit(1)
	}

	switch name {
	case "roles":
		err = seeder.Roles(ctx, usecase.NewRoleUseCase(postgres.NewRoleRepository(pool)))
	case "books":
		err = seeder.Books(ctx, usecase.NewBookUseCase(postgres.NewBookRepository(pool)))
	case "menus":
		err = seeder.Menus(ctx, usecase.NewMenuUseCase(postgres.NewMenuRepository(pool)))
	default:
		fmt.Fprintf(os.Stderr, "unknown seeder %q\n", name)
		fmt.Fprintln(os.Stderr, "available: roles, books, menus")
		os.Exit(1)
	}

	if err != nil {
		slog.Error("seeder failed", "name", name, "err", err)
		os.Exit(1)
	}
	slog.Info("seeder finished", "name", name)
}
