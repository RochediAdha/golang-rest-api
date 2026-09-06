# Golang REST API

REST API untuk manajemen user, role, menu, privilege, dan otorisasi. Dibangun dengan Go standard library (`net/http`) dan PostgreSQL, mengikuti Clean Architecture.

Dokumentasi lengkap ada di folder [`docs/`](docs/README.md):

- [Struktur folder](docs/struktur-folder.md)
- [Alur aplikasi](docs/alur-aplikasi.md)
- [Seeder](docs/seeder.md)
- [Deploy di server](docs/deploy.md)
- [Spec API](docs/api/README.md)

## Struktur

```
cmd/api                         # composition root: wiring config, DB, usecase, HTTP
cmd/seed                        # seeder manual (roles, menus, privileges)
internal/
  domain/                       # entitas, error, dan port repository
  usecase/                      # aturan bisnis
  adapter/
    http/                       # inbound: handler, router, middleware
    postgres/                   # outbound: penyimpanan PostgreSQL
    memory/                     # outbound: penyimpanan in-memory (tes)
  infrastructure/
    config/                     # environment & DSN
    database/                   # connection pool & migrasi
      migrations/               # file SQL berversi
```

Request mengalir: `adapter/http` → `usecase` → `adapter/postgres` → PostgreSQL.

## Database

API memakai PostgreSQL (lokal atau server terpisah). Salin konfigurasi lalu sesuaikan host, user, password, dan nama database:

```bash
cp .env.example .env
```

Variabel yang dibaca aplikasi:

```
ADDR=:8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=golang_rest_api
DB_SSLMODE=disable
DB_MAX_CONNS=10
DB_MIN_CONNS=1
```

File `.env` tidak men-expand `${VAR}`; tulis nilainya langsung. `DATABASE_URL` opsional: jika diisi, mengalahkan `DB_*`.

Saat server start, migrasi dijalankan dan dicatat di tabel `schema_migrations` (`version`, `name`, `applied_at`). File SQL ada di `internal/infrastructure/database/migrations/`. Migrasi yang sudah tercatat tidak dijalankan ulang.

## Menjalankan

```bash
go run ./cmd/api
```

Server default: `http://localhost:8080`

## Docker

Image hanya berisi API. PostgreSQL di server terpisah; isi `DB_HOST` (bukan `localhost` jika DB di mesin lain).

```bash
cp .env.example .env
docker compose up -d --build
```

Health: `http://127.0.0.1:8080/health`. Panduan server: [docs/deploy.md](docs/deploy.md).

## Seeder

Seeder tidak jalan otomatis saat server start. Admin menjalankan per file:

```bash
go run ./cmd/seed roles
go run ./cmd/seed menus
go run ./cmd/seed privileges
```

atau `make seed name=roles`. Data yang sudah ada dilewati (tidak diduplikasi).

## Endpoint

# Health
| Method | Path | Deskripsi |
| --- | --- | --- |
| `GET` | `/` | Informasi singkat API |
| `GET` | `/health` | Health check (termasuk ping database) |

# Users
| Method | Path | Deskripsi |
| --- | --- | --- |
| `GET` | `/api/v1/users` | Daftar user (`q`, `limit`, `offset`) |
| `GET` | `/api/v1/users/{id}` | Detail user |
| `POST` | `/api/v1/users` | Tambah user |
| `PUT` | `/api/v1/users/{id}` | Ubah user |
| `DELETE` | `/api/v1/users/{id}` | Hapus user |

# Roles
| Method | Path | Deskripsi |
| --- | --- | --- |
| `GET` | `/api/v1/roles` | Daftar role (`q`, `limit`, `offset`) |
| `GET` | `/api/v1/roles/{id}` | Detail role |
| `POST` | `/api/v1/roles` | Tambah role |
| `PUT` | `/api/v1/roles/{id}` | Ubah role |
| `DELETE` | `/api/v1/roles/{id}` | Soft delete role |

# Menus
| Method | Path | Deskripsi |
| --- | --- | --- |
| `GET` | `/api/v1/menus` | Daftar menu (`q`, `parentId`, `limit`, `offset`) |
| `GET` | `/api/v1/menus/{id}` | Detail menu |
| `POST` | `/api/v1/menus` | Tambah menu |
| `PUT` | `/api/v1/menus/{id}` | Ubah menu |
| `DELETE` | `/api/v1/menus/{id}` | Soft delete menu |

# User roles
| Method | Path | Deskripsi |
| --- | --- | --- |
| `GET` | `/api/v1/user-roles` | Daftar penugasan (nama user/role, `number`; filter `userId`, `roleId`, `number`) |
| `GET` | `/api/v1/user-roles/{id}` | Detail user beserta semua role-nya (`{id}` = userId atau assignment id) |
| `POST` | `/api/v1/user-roles` | Tugaskan role ke user |
| `DELETE` | `/api/v1/user-roles/{id}` | Hapus penugasan |

# Privileges
| Method | Path | Deskripsi |
| --- | --- | --- |
| `GET` | `/api/v1/privileges` | Daftar privilege (`q`, `limit`, `offset`) |
| `GET` | `/api/v1/privileges/{id}` | Detail privilege |
| `POST` | `/api/v1/privileges` | Tambah privilege |
| `PUT` | `/api/v1/privileges/{id}` | Ubah privilege |
| `DELETE` | `/api/v1/privileges/{id}` | Hapus privilege |

# Role privileges
| Method | Path | Deskripsi |
| --- | --- | --- |
| `GET` | `/api/v1/role-privileges` | Daftar otorisasi (`roleId`, `menuId`, `privilegeId`, `limit`, `offset`) |
| `GET` | `/api/v1/role-privileges/{id}` | Detail otorisasi |
| `POST` | `/api/v1/role-privileges` | Tugaskan privilege ke role pada menu |
| `DELETE` | `/api/v1/role-privileges/{id}` | Hapus otorisasi |

## Tes

```bash
go test ./...
```
