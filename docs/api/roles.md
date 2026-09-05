# Roles

Manajemen role. ID berupa UUID. Delete bersifat soft delete.

## Object

| Field | Tipe | Keterangan |
| --- | --- | --- |
| `id` | uuid | Diisi server |
| `name` | string | Wajib, max 80, unik (yang belum dihapus) |
| `description` | string | Opsional, max 500 |
| `isActive` | boolean | Default `true` |
| `createdAt` | datetime | Diisi server |
| `updatedAt` | datetime | Diisi server |
| `deletedAt` | datetime | Terisi setelah soft delete |
| `createdBy` | uuid | Opsional |
| `updatedBy` | uuid | Opsional |
| `deletedBy` | uuid | Opsional |

Role yang sudah punya `deletedAt` tidak muncul di list/detail.

## `GET /api/v1/roles`

Query: `q` (nama, deskripsi), `limit`, `offset`.

```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-4789-abcd-1234567890ab",
      "name": "Super Admin",
      "description": "Full access to all features",
      "isActive": true,
      "createdAt": "2026-09-05T08:00:00Z",
      "updatedAt": "2026-09-05T08:00:00Z"
    }
  ],
  "meta": { "total": 1, "limit": 20, "offset": 0 }
}
```

## `GET /api/v1/roles/{id}`

`id` harus UUID. Response `200` object role. `400` / `404` jika tidak valid atau tidak ada.

## `POST /api/v1/roles`

```json
{
  "name": "Admin",
  "description": "Can create and update content",
  "isActive": true,
  "createdBy": "11111111-1111-4111-8111-111111111111"
}
```

`description`, `isActive`, dan `createdBy` opsional.

Response `201` object role.

Error: `400 invalid_input`, `409 duplicate_role_name`.

## `PUT /api/v1/roles/{id}`

```json
{
  "name": "Admin",
  "description": "Administrator",
  "isActive": false,
  "updatedBy": "11111111-1111-4111-8111-111111111111"
}
```

`name` wajib. Jika `isActive` tidak dikirim, nilai lama tetap dipakai.

Response `200` object role.

## `DELETE /api/v1/roles/{id}`

Soft delete. Body opsional:

```json
{
  "deletedBy": "11111111-1111-4111-8111-111111111111"
}
```

Response `200` object role yang baru dihapus, termasuk `deletedAt`.

```json
{
  "id": "a1b2c3d4-e5f6-4789-abcd-1234567890ab",
  "name": "Admin",
  "description": "Administrator",
  "isActive": false,
  "createdAt": "2026-09-05T08:00:00Z",
  "updatedAt": "2026-09-05T08:10:00Z",
  "deletedAt": "2026-09-05T08:10:00Z",
  "deletedBy": "11111111-1111-4111-8111-111111111111"
}
```

Data awal role: jalankan `go run ./cmd/seed roles`. Lihat [seeder](../seeder.md).
