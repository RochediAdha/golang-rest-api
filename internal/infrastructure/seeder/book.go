package seeder

import (
	"context"
	"errors"
	"log/slog"

	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

var defaultBooks = []domain.CreateBookInput{
	{Title: "The Go Programming Language", Author: "Alan A. A. Donovan", ISBN: "9780134190440", Year: 2015},
	{Title: "Concurrency in Go", Author: "Katherine Cox-Buday", ISBN: "9781491941195", Year: 2017},
}

func Books(ctx context.Context, uc *usecase.BookUseCase) error {
	for _, sample := range defaultBooks {
		if _, err := uc.Create(ctx, sample); err != nil {
			if errors.Is(err, domain.ErrDuplicateISBN) {
				continue
			}
			slog.Error("seed book failed", "title", sample.Title, "err", err)
			return err
		}
		slog.Info("seed book created", "title", sample.Title)
	}
	return nil
}
