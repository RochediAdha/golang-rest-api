# Golang REST API

REST API untuk manajemen user, role, menu, privilege, dan otorisasi. Dibangun dengan Go standard library (`net/http`) dan PostgreSQL, mengikuti Clean Architecture.

Dokumentasi lengkap ada di folder [`docs/`](docs/README.md):

- [Struktur folder](docs/struktur-folder.md)
- [Alur aplikasi](docs/alur-aplikasi.md)
- [Seeder](docs/seeder.md)
- [Spec API](docs/api/README.md)

## Struktur

```
cmd/api                         # composition root: wiring config, DB, usecase, HTTP
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

API memakai PostgreSQL lokal. Salin konfigurasi lalu sesuaikan user, password, dan nama database:

```bash
cp .env.example .env
```

Default DSN:

```
postgres://postgres:postgres@localhost:5432/golang_rest_api?sslmode=disable
```

Saat server start, migrasi dijalankan dan dicatat di tabel `schema_migrations` (`version`, `name`, `applied_at`). File SQL ada di `internal/infrastructure/database/migrations/`. Migrasi yang sudah tercatat tidak dijalankan ulang.

## Menjalankan

```bash
go run ./cmd/api
```

Server default: `http://localhost:8080`

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
| `GET` | `/api/v1/user-roles` | Daftar penugasan (`userId`, `roleId`, `number`, `limit`, `offset`) |
| `GET` | `/api/v1/user-roles/{id}` | Detail penugasan |
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
