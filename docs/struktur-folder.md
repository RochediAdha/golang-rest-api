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
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
├── .env.example
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
| `main.go` | Menerima nama seeder (`roles`, `menus`, atau `privileges`), konek database, lalu memanggil fungsi di `internal/infrastructure/seeder`. |

## `internal/`

Kode inti aplikasi. Package di luar repo tidak bisa meng-import folder ini.

### `internal/domain/`

Lapisan terdalam: entitas, error, dan kontrak penyimpanan (port). Tidak tahu HTTP atau SQL.

| File | Fungsi |
| --- | --- |
| `user.go` | Entitas `User` dan input create/update |
| `role.go` | Entitas `Role` dan input create/update/delete |
| `menu.go` | Entitas `Menu` dan input create/update/delete |
| `user_role.go` | Entitas `UserRole`, list (`UserRoleListItem`), show (`UserRoleView`) |
| `privilege.go` | Entitas `Privilege` dan input create/update |
| `role_privilege.go` | Entitas `RolePrivilege` dan input create |
| `errors.go` | Error domain (`ErrNotFound`, `ErrInvalidInput`, duplikat, dll.) |
| `repository.go` | Interface repository tiap modul |

### `internal/usecase/`

Aturan bisnis. Hanya bergantung pada `domain`.

| File | Fungsi |
| --- | --- |
| `user.go` | Validasi user, UUID, cek username/email unik |
| `role.go` | Validasi role, soft delete (`deletedAt` / `deletedBy`) |
| `menu.go` | Validasi menu, parent, `sortOrder`, soft delete |
| `user_role.go` | Validasi penugasan, `number` (NIM/NIP), list/show dengan nama user dan role |
| `privilege.go` | Validasi privilege, `code` unik |
| `role_privilege.go` | Validasi otorisasi role-menu-privilege, kombinasi unik |
| `*_test.go` | Tes usecase memakai `adapter/memory` |

### `internal/adapter/http/`

Adapter inbound: terjemahkan HTTP menjadi pemanggilan usecase.

| File | Fungsi |
| --- | --- |
| `router.go` | Daftar route dan middleware |
| `user.go` | Handler CRUD `/api/v1/users` |
| `role.go` | Handler CRUD `/api/v1/roles` |
| `menu.go` | Handler CRUD `/api/v1/menus` |
| `user_role.go` | Handler `/api/v1/user-roles` (create, list, get, delete) |
| `privilege.go` | Handler CRUD `/api/v1/privileges` |
| `role_privilege.go` | Handler `/api/v1/role-privileges` (create, list, get, delete) |
| `health.go` | `GET /health`, ping database |
| `response.go` | Format JSON sukses/error |
| `middleware.go` | Recover panic, logging, CORS |
| `*_test.go` | Tes HTTP dengan repository in-memory |

### `internal/adapter/postgres/`

Adapter outbound: implementasi repository ke PostgreSQL.

| File | Fungsi |
| --- | --- |
| `user.go` | SQL tabel `users` |
| `role.go` | SQL tabel `roles` (termasuk soft delete) |
| `menu.go` | SQL tabel `menus` (parent + soft delete) |
| `user_role.go` | SQL tabel `user_roles` |
| `privilege.go` | SQL tabel `privileges` |
| `role_privilege.go` | SQL tabel `role_privileges` |

### `internal/adapter/memory/`

Implementasi repository di memori. Dipakai tes, bukan production.

| File | Fungsi |
| --- | --- |
| `user.go` | Store user in-memory |
| `role.go` | Store role in-memory |
| `menu.go` | Store menu in-memory |
| `user_role.go` | Store user role in-memory |
| `privilege.go` | Store privilege in-memory |
| `role_privilege.go` | Store role privilege in-memory |

### `internal/infrastructure/config/`

Membaca environment dan `.env`.

| File | Fungsi |
| --- | --- |
| `config.go` | `ADDR`, `DB_*` (atau `DATABASE_URL`), timeout, ukuran pool |
| `config_test.go` | Tes penyusunan DSN |

### `internal/infrastructure/database/`

Koneksi dan migrasi skema.

| File | Fungsi |
| --- | --- |
| `postgres.go` | Buka `pgx` pool dan ping |
| `migrate.go` | Jalankan file SQL yang belum tercatat |
| `embed.go` | Embed folder `migrations/` ke binary |
| `migrate_test.go` | Tes parsing nama file migrasi |
| `migrations/000002_create_users.sql` | Tabel `users` |
| `migrations/000003_create_roles.sql` | Tabel `roles` |
| `migrations/000004_create_menus.sql` | Tabel `menus` + enum `MenuType` |
| `migrations/000005_recreate_menus.sql` | Recreate `menus` jika tabel sudah di-drop |
| `migrations/000006_create_user_roles.sql` | Tabel `user_roles` (termasuk `number`) |
| `migrations/000007_create_privileges.sql` | Tabel `privileges` |
| `migrations/000008_create_role_privileges.sql` | Tabel `role_privileges` |

Progress migrasi disimpan di tabel `schema_migrations` (`version`, `name`, `applied_at`).

### `internal/infrastructure/seeder/`

Data awal. Tidak dipanggil saat server start.

| File | Fungsi |
| --- | --- |
| `role.go` | Seed role default |
| `menu.go` | Seed menu default (termasuk hierarki parent) |
| `privilege.go` | Seed privilege default |
| `*_test.go` | Tes seeder idempotent |

Cara menjalankan: lihat [seeder.md](seeder.md).

## File di root

| File | Fungsi |
| --- | --- |
| `go.mod` / `go.sum` | Modul Go (`golang-rest-api`) dan dependency |
| `Makefile` | `run`, `test`, `tidy`, `seed`, `docker-up`, `docker-down`, `docker-seed` |
| `Dockerfile` | Multi-stage build binary `api` dan `seed` |
| `docker-compose.yml` | Container API saja; Postgres di server terpisah |
| `.dockerignore` | File yang tidak masuk image |
| `.env` / `.env.example` | `DB_*` (`.env` tidak di-commit) |
| `README.md` | Ringkasan cara jalan dan endpoint |

## Pemetaan singkat ke MVC

| MVC | Folder di project ini |
| --- | --- |
| Controller | `internal/adapter/http` |
| Model | `internal/domain` + `internal/usecase` + `internal/adapter/postgres` |
| View | Response JSON (tidak ada folder view) |
