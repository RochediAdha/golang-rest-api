## Seeder

Seeder tidak jalan saat server start. Admin menjalankannya sendiri, per file.

```
go run ./cmd/seed roles
go run ./cmd/seed menus
go run ./cmd/seed privileges
```

atau:

```
make seed name=menus
```

Lewat Docker Compose:

```
docker compose run --rm --entrypoint /app/seed api roles
```

atau `make docker-seed name=menus`.

Tanpa nama seeder, perintah menampilkan daftar yang tersedia: roles, menus, privileges.

Data yang sudah ada dilewati, jadi aman dijalankan ulang.

## Isi data awal

### `roles`

- Super Admin
- Admin
- Lecturer
- Student

### `privileges`

CREATE, EXECUTE, APPROVE, VIEW, PRINT, PUBLISH, DELETE, EXPORT, UPDATE, REJECT, IMPORT

### `menus`

| Code | Parent | Tipe |
| --- | --- | --- |
| `dashboard` | — | ITEM |
| `master` | — | GROUP |
| `user-management` | — | GROUP |
| `users` | `user-management` | ITEM |
| `roles` | `user-management` | ITEM |
| `permission` | `user-management` | ITEM |
| `system` | — | GROUP |
| `menus` | `system` | ITEM |

File seedernya:

```
internal/infrastructure/seeder/role.go
internal/infrastructure/seeder/menu.go
internal/infrastructure/seeder/privilege.go
```
