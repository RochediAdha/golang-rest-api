## Seeder

Seeder tidak lagi jalan saat server start. Admin menjalankannya sendiri, per file.

```
go run ./cmd/seed roles
go run ./cmd/seed books
go run ./cmd/seed menus
```

atau:

```
make seed name=menus
```

Tanpa nama seeder, perintah menampilkan daftar yang tersedia: roles, books, menus.

File seedernya tetap terpisah:

```
internal/infrastructure/seeder/role.go
internal/infrastructure/seeder/book.go
internal/infrastructure/seeder/menu.go
```

Data yang sudah ada tetap dilewati, jadi aman dijalankan ulang.
