## Seeder

Seeder tidak lagi jalan saat server start. Admin menjalankannya sendiri, per file.

```
go run ./cmd/seed roles
go run ./cmd/seed menus
go run ./cmd/seed privileges
```

atau:

```
make seed name=menus
```

Tanpa nama seeder, perintah menampilkan daftar yang tersedia: roles, menus, privileges.

File seedernya tetap terpisah:

```
internal/infrastructure/seeder/role.go
internal/infrastructure/seeder/menu.go
internal/infrastructure/seeder/privilege.go
```

Data yang sudah ada tetap dilewati, jadi aman dijalankan ulang.
