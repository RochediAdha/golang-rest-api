package database

import "testing"

func TestParseMigration(t *testing.T) {
	m, err := parseMigration("000001_create_books.sql")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if m.Version != 1 || m.Name != "create_books" {
		t.Fatalf("got version=%d name=%q", m.Version, m.Name)
	}
}

func TestParseMigrationInvalid(t *testing.T) {
	if _, err := parseMigration("create_books.sql"); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadMigrations(t *testing.T) {
	migrations, err := loadMigrations()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(migrations) == 0 {
		t.Fatal("expected at least one migration")
	}
	if migrations[0].Version != 1 || migrations[0].Name != "create_books" {
		t.Fatalf("first migration = %+v", migrations[0])
	}
	if migrations[0].SQL == "" {
		t.Fatal("expected sql body")
	}
}
