package config

import "testing"

func TestDatabaseDSNFromURL(t *testing.T) {
	cfg := Database{URL: "postgres://app:secret@db:5432/appdb?sslmode=require"}
	if got := cfg.DSN(); got != cfg.URL {
		t.Fatalf("DSN() = %q", got)
	}
}

func TestDatabaseDSNFromParts(t *testing.T) {
	cfg := Database{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		Name:     "golang_rest_api",
		SSLMode:  "disable",
	}
	want := "postgres://postgres:postgres@localhost:5432/golang_rest_api?sslmode=disable"
	if got := cfg.DSN(); got != want {
		t.Fatalf("DSN() = %q, want %q", got, want)
	}
}
