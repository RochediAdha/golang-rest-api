# Struktur folder

Project ini memakai Clean Architecture. Aturan bisnis tidak bergantung pada HTTP atau PostgreSQL. Yang berubah hanya adapter dan infrastructure.

```
.
├── cmd/
│   ├── api/
│   └── seed/
├── docs/
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── adapter/
│   │   ├── http/
│   │   ├── postgres/
│   │   └── memory/
│   └── infrastructure/
│       ├── config/
│       ├── database/
│       │   └── migrations/
│       └── seeder/
├── Makefile
├── go.mod
└── README.md
```

## `cmd/`

Titik masuk program. Package di sini adalah `main`.

### `cmd/api/`

Menjalankan REST API.

| File | Fungsi |
| --- | --- |
| `main.go` | Baca config, buka koneksi PostgreSQL, jalankan migrasi, rangkai usecase + router, lalu listen HTTP. Saat `SIGINT`/`SIGTERM`, server shutdown lalu pool ditutup. |

Tidak menjalankan seeder. Seeder dipanggil lewat `cmd/seed`.

### `cmd/seed/`

Menjalankan satu file seeder secara manual.

| File | Fungsi |
| --- | --- |
| `main.go` | Menerima nama seeder (`roles`, `books`, atau `menus`), konek database, lalu memanggil fungsi di `internal/infrastructure/seeder`. |

## `internal/`

Kode inti aplikasi. Package di luar repo tidak bisa meng-import folder ini.

### `internal/domain/`

Lapisan terdalam: entitas, error, dan kontrak penyimpanan (port). Tidak tahu HTTP atau SQL.

| File | Fungsi |
| --- | --- |
| `book.go` | Entitas `Book` dan input create/update |
| `user.go` | Entitas `User` dan input create/update |
| `role.go` | Entitas `Role` dan input create/update/delete |
| `menu.go` | Entitas `Menu` dan input create/update/delete |
| `errors.go` | Error domain (`ErrNotFound`, `ErrInvalidInput`, duplikat, dll.) |
| `repository.go` | Interface `BookRepository`, `UserRepository`, `RoleRepository`, `MenuRepository` |

### `internal/usecase/`

Aturan bisnis. Hanya bergantung pada `domain`.

| File | Fungsi |
| --- | --- |
| `book.go` | Validasi buku, generate ID, pagination |
| `user.go` | Validasi user, UUID, cek username/email unik |
| `role.go` | Validasi role, soft delete (`deletedAt` / `deletedBy`) |
| `menu.go` | Validasi menu, parent, `sortOrder`, soft delete |
| `*_test.go` | Tes usecase memakai `adapter/memory` |

### `internal/adapter/http/`

Adapter inbound: terjemahkan HTTP menjadi pemanggilan usecase.

| File | Fungsi |
| --- | --- |
| `router.go` | Daftar route dan middleware |
| `book.go` | Handler CRUD `/api/v1/books` |
| `user.go` | Handler CRUD `/api/v1/users` |
| `role.go` | Handler CRUD `/api/v1/roles` |
| `menu.go` | Handler CRUD `/api/v1/menus` |
| `health.go` | `GET /health`, ping database |
| `response.go` | Format JSON sukses/error |
| `middleware.go` | Recover panic, logging, CORS |
| `*_test.go` | Tes HTTP dengan repository in-memory |

### `internal/adapter/postgres/`

Adapter outbound: implementasi repository ke PostgreSQL.

| File | Fungsi |
| --- | --- |
| `book.go` | SQL tabel `books` |
| `user.go` | SQL tabel `users` |
| `role.go` | SQL tabel `roles` (termasuk soft delete) |
| `menu.go` | SQL tabel `menus` (parent + soft delete) |

### `internal/adapter/memory/`

Implementasi repository di memori. Dipakai tes, bukan production.

| File | Fungsi |
| --- | --- |
| `book.go` | Store buku in-memory |
| `user.go` | Store user in-memory |
| `role.go` | Store role in-memory |
| `menu.go` | Store menu in-memory |

### `internal/infrastructure/config/`

Membaca environment dan `.env`.

| File | Fungsi |
| --- | --- |
| `config.go` | `ADDR`, DSN PostgreSQL, timeout, ukuran pool |
| `config_test.go` | Tes penyusunan DSN |

### `internal/infrastructure/database/`

Koneksi dan migrasi skema.

| File | Fungsi |
| --- | --- |
| `postgres.go` | Buka `pgx` pool dan ping |
| `migrate.go` | Jalankan file SQL yang belum tercatat |
| `embed.go` | Embed folder `migrations/` ke binary |
| `migrate_test.go` | Tes parsing nama file migrasi |
| `migrations/000001_create_books.sql` | Tabel `books` |
| `migrations/000002_create_users.sql` | Tabel `users` |
| `migrations/000003_create_roles.sql` | Tabel `roles` |
| `migrations/000004_create_menus.sql` | Tabel `menus` + enum `MenuType` |
| `migrations/000005_recreate_menus.sql` | Recreate `menus` jika tabel sudah di-drop |

Progress migrasi disimpan di tabel `schema_migrations` (`version`, `name`, `applied_at`).

### `internal/infrastructure/seeder/`

Data awal. Tidak dipanggil saat server start.

| File | Fungsi |
| --- | --- |
| `role.go` | Seed role default |
| `book.go` | Seed buku contoh |
| `menu.go` | Seed menu default (termasuk hierarki parent) |
| `*_test.go` | Tes seeder idempotent |

Cara menjalankan: lihat [seeder.md](seeder.md).

## File di root

| File | Fungsi |
| --- | --- |
| `go.mod` / `go.sum` | Modul Go (`golang-rest-api`) dan dependency |
| `Makefile` | `run`, `test`, `tidy`, `seed` |
| `.env` / `.env.example` | Konfigurasi lokal (`.env` tidak di-commit) |
| `README.md` | Ringkasan cara jalan dan endpoint |

## Pemetaan singkat ke MVC

| MVC | Folder di project ini |
| --- | --- |
| Controller | `internal/adapter/http` |
| Model | `internal/domain` + `internal/usecase` + `internal/adapter/postgres` |
| View | Response JSON (tidak ada folder view) |
