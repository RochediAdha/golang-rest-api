package config

import (
	"bufio"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	Database        Database
}

type Database struct {
	URL             string
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxConnIdleTime time.Duration
	MaxConnLifetime time.Duration
}

func Load() Config {
	loadDotEnv(".env")

	return Config{
		Addr:            env("ADDR", ":8080"),
		ReadTimeout:     durationEnv("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    durationEnv("WRITE_TIMEOUT", 10*time.Second),
		ShutdownTimeout: durationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		Database: Database{
			URL:             env("DATABASE_URL", ""),
			Host:            env("DB_HOST", "localhost"),
			Port:            env("DB_PORT", "5432"),
			User:            env("DB_USER", "postgres"),
			Password:        env("DB_PASSWORD", "postgres"),
			Name:            env("DB_NAME", "golang_rest_api"),
			SSLMode:         env("DB_SSLMODE", "disable"),
			MaxConns:        int32Env("DB_MAX_CONNS", 10),
			MinConns:        int32Env("DB_MIN_CONNS", 1),
			MaxConnIdleTime: durationEnv("DB_MAX_CONN_IDLE_TIME", 5*time.Minute),
			MaxConnLifetime: durationEnv("DB_MAX_CONN_LIFETIME", time.Hour),
		},
	}
}

func (d Database) DSN() string {
	if d.URL != "" {
		return d.URL
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Password),
		Host:   net.JoinHostPort(d.Host, d.Port),
		Path:   "/" + d.Name,
	}
	q := url.Values{}
	q.Set("sslmode", d.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	if seconds, err := strconv.Atoi(raw); err == nil {
		return time.Duration(seconds) * time.Second
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return d
}

func int32Env(key string, fallback int32) int32 {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return int32(n)
}

func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}
