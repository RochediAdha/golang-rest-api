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

	bookUC := usecase.NewBookUseCase(postgres.NewBookRepository(pool))
	userUC := usecase.NewUserUseCase(postgres.NewUserRepository(pool))
	roleUC := usecase.NewRoleUseCase(postgres.NewRoleRepository(pool))
	menuUC := usecase.NewMenuUseCase(postgres.NewMenuRepository(pool))

	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      httpadapter.NewRouter(bookUC, userUC, roleUC, menuUC, pool),
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
