package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "golang-rest-api/internal/adapter/http"
	"golang-rest-api/internal/adapter/postgres"
	"golang-rest-api/internal/infrastructure/config"
	"golang-rest-api/internal/infrastructure/database"
	"golang-rest-api/internal/usecase"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg := config.Load()

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	pool, err := database.Open(dbCtx, cfg.Database)
	dbCancel()
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := database.Migrate(context.Background(), pool); err != nil {
		slog.Error("database migration failed", "err", err)
		os.Exit(1)
	}

	userRepo := postgres.NewUserRepository(pool)
	roleRepo := postgres.NewRoleRepository(pool)
	menuRepo := postgres.NewMenuRepository(pool)
	privilegeRepo := postgres.NewPrivilegeRepository(pool)
	userUC := usecase.NewUserUseCase(userRepo)
	roleUC := usecase.NewRoleUseCase(roleRepo)
	menuUC := usecase.NewMenuUseCase(menuRepo)
	userRoleUC := usecase.NewUserRoleUseCase(postgres.NewUserRoleRepository(pool), userRepo, roleRepo)
	privilegeUC := usecase.NewPrivilegeUseCase(privilegeRepo)
	rolePrivilegeUC := usecase.NewRolePrivilegeUseCase(postgres.NewRolePrivilegeRepository(pool), roleRepo, menuRepo, privilegeRepo)

	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      httpadapter.NewRouter(userUC, roleUC, menuUC, userRoleUC, privilegeUC, rolePrivilegeUC, pool),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		slog.Info("server started", "addr", cfg.Addr, "database", cfg.Database.Name)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}
