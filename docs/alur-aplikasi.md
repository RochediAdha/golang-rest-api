# Alur aplikasi

Dokumen ini menjelaskan bagaimana aplikasi berjalan: dari server start, request HTTP, migrasi, sampai seeder.

## 1. Server start

Perintah: `go run ./cmd/api`

Urutan di `cmd/api/main.go`:

1. Logger JSON diaktifkan.
2. Config dibaca dari environment dan file `.env` (`DB_*`, atau `DATABASE_URL` jika diisi). Interpolasi `${VAR}` di `.env` tidak di-expand.
3. Connection pool PostgreSQL dibuka (`database.Open`). Kalau gagal, proses berhenti.
4. Migrasi dijalankan (`database.Migrate`). File SQL yang belum ada di `schema_migrations` dieksekusi, lalu dicatat.
5. Repository PostgreSQL dibuat, lalu diikat ke usecase:
   - `UserUseCase`
   - `RoleUseCase`
   - `MenuUseCase`
   - `UserRoleUseCase`
   - `PrivilegeUseCase`
   - `RolePrivilegeUseCase`
6. Router HTTP dipasang (`adapter/http`). Pool dipakai juga untuk `GET /health`.
7. Server listen di `ADDR` (default `:8080`).
8. Proses menunggu `SIGINT` / `SIGTERM`. Saat berhenti, `Shutdown` dipanggil lalu pool ditutup.

Seeder **tidak** ikut di langkah ini.

```
cmd/api
  → config.Load()
  → database.Open()
  → database.Migrate()
  → usecase + postgres repository
  → http.ListenAndServe
  → shutdown
```

## 2. Request HTTP

Contoh: `POST /api/v1/roles`

```
Client
  → middleware (recover → logging → CORS)
  → router
  → handler (adapter/http)
  → usecase
  → repository (adapter/postgres)
  → PostgreSQL
  → JSON response
```

Rincian:

1. Request masuk `ServeMux`.
2. Middleware menangani panic, menulis log, dan CORS. `OPTIONS` dijawab `204`.
3. Handler membaca body JSON dan path `{id}`.
4. Usecase memvalidasi data, membuat UUID/ID, dan menerapkan aturan (duplikat, pagination, soft delete).
5. Repository menjalankan SQL.
6. Handler mengubah hasil atau error domain menjadi JSON + status HTTP:
   - input salah → `400`
   - tidak ketemu → `404`
   - data duplikat → `409`
   - error lain → `500`

`GET /health` tidak lewat usecase. Handler langsung ping pool. Kalau database tidak merespons, status `degraded` / `503`.

## 3. Alur tiap modul

Semua modul CRUD mengikuti pola yang sama. Perbedaan ada di aturan bisnis.

| Modul | Path | Catatan |
| --- | --- | --- |
| Users | `/api/v1/users` | ID UUID; username dan email unik |
| Roles | `/api/v1/roles` | ID UUID; nama unik; delete adalah soft delete |
| Menus | `/api/v1/menus` | ID UUID; `code` unik; `parentId` opsional; delete adalah soft delete |
| User roles | `/api/v1/user-roles` | ID UUID; pasangan `userId`+`roleId` unik; list menampilkan nama user/role, deskripsi, dan `number`; show menampilkan semua role per user; delete permanen |
| Privileges | `/api/v1/privileges` | ID UUID; `code` unik; delete permanen |
| Role privileges | `/api/v1/role-privileges` | ID UUID; kombinasi role+menu+privilege unik; delete permanen |

Soft delete role:

1. Handler `DELETE /api/v1/roles/{id}` (body `deletedBy` opsional).
2. Usecase mengambil role yang belum terhapus, mengisi `deletedAt` / `deletedBy`.
3. Repository `UPDATE` baris, bukan `DELETE` fisik.
4. Response `200` berisi data role yang baru dihapus.
5. `GET` list/detail tidak menampilkan baris yang sudah punya `deletedAt`.

## 4. Migrasi

Dipanggil saat API atau seeder start.

1. Ambil advisory lock supaya tidak bentrok.
2. Buat tabel `schema_migrations` jika belum ada.
3. Baca file `internal/infrastructure/database/migrations/*.sql` (urutan nomor).
4. Lewati versi yang sudah tercatat.
5. Jalankan SQL baru, lalu `INSERT` ke `schema_migrations`.

Contoh: `000003_create_roles.sql` → `version=3`, `name=create_roles`.

## 5. Seeder (manual)

Tidak otomatis. Admin memilih file yang dijalankan:

```bash
go run ./cmd/seed roles
go run ./cmd/seed menus
go run ./cmd/seed privileges
```

Alur `cmd/seed`:

1. Baca nama seeder dari argumen.
2. Load config, buka database, jalankan migrasi.
3. Panggil fungsi di `internal/infrastructure/seeder` sesuai nama.
4. Data yang sudah ada dilewati.

Detail perintah: [seeder.md](seeder.md).

## 6. Tes

```
go test ./...
```

Tes HTTP dan usecase memakai `adapter/memory`, jadi tidak butuh PostgreSQL. Tes migrasi hanya memeriksa parsing file SQL, tidak menulis ke database.

## 7. Deploy

Di server, proses start sama (`cmd/api` → migrasi → listen). Perbedaannya hanya cara menjalankan proses: Docker Compose atau systemd. Lihat [deploy.md](deploy.md).
