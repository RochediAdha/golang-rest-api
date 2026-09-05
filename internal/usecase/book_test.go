package usecase

import (
	"context"
	"testing"
	"time"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/domain"
)

func TestBookUseCaseCRUD(t *testing.T) {
	ctx := context.Background()
	svc := NewBookUseCase(memory.NewBookRepository())
	svc.now = func() time.Time { return time.Date(2026, 9, 5, 3, 0, 0, 0, time.UTC) }

	created, err := svc.Create(ctx, domain.CreateBookInput{
		Title:  "  Learning Go  ",
		Author: "Jon Bodner",
		ISBN:   "978-1492077213",
		Year:   2021,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" || created.Title != "Learning Go" || created.ISBN != "9781492077213" {
		t.Fatalf("unexpected created book: %+v", created)
	}

	got, err := svc.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("get mismatch: %+v", got)
	}

	updated, err := svc.Update(ctx, created.ID, domain.UpdateBookInput{
		Title:  "Learning Go, 2nd Edition",
		Author: "Jon Bodner",
		ISBN:   "9781492077213",
		Year:   2024,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Title != "Learning Go, 2nd Edition" || updated.Year != 2024 {
		t.Fatalf("unexpected update: %+v", updated)
	}

	list, total, err := svc.List(ctx, domain.ListFilter{Query: "go", Limit: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("list want 1 got total=%d len=%d", total, len(list))
	}

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.Get(ctx, created.ID); err != domain.ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestBookUseCaseValidation(t *testing.T) {
	ctx := context.Background()
	svc := NewBookUseCase(memory.NewBookRepository())

	_, err := svc.Create(ctx, domain.CreateBookInput{Title: "", Author: "A"})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}

	first, err := svc.Create(ctx, domain.CreateBookInput{Title: "A", Author: "B", ISBN: "1234567890"})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	_, err = svc.Create(ctx, domain.CreateBookInput{Title: "C", Author: "D", ISBN: "1234567890"})
	if err != domain.ErrDuplicateISBN {
		t.Fatalf("expected duplicate isbn, got %v", err)
	}

	if _, err := svc.Get(ctx, "missing"); err != domain.ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	_ = first
}
