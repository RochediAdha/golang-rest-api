# Privileges

Manajemen privilege. ID berupa UUID. Hapus bersifat permanen.

## Object

| Field | Tipe | Keterangan |
| --- | --- | --- |
| `id` | uuid | Diisi server |
| `code` | string | Wajib, tanpa spasi, max 80, unik (case-insensitive) |
| `name` | string | Wajib, max 120 |
| `description` | string | Opsional, max 500 |
| `isActive` | boolean | Default `true` |
| `createdAt` | datetime | Diisi server |
| `updatedAt` | datetime | Diisi server |
| `createdBy` | uuid | Opsional |
| `updatedBy` | uuid | Opsional |

## `GET /api/v1/privileges`

Query: `q` (code, nama, deskripsi), `limit`, `offset`.

```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-4789-abcd-1234567890ab",
      "code": "user.read",
      "name": "Read User",
      "description": "View user data",
      "isActive": true,
      "createdAt": "2026-09-05T08:00:00Z",
      "updatedAt": "2026-09-05T08:00:00Z"
    }
  ],
  "meta": { "total": 1, "limit": 20, "offset": 0 }
}
```

## `GET /api/v1/privileges/{id}`

`id` harus UUID. Response `200` object privilege. `400` / `404` jika tidak valid atau tidak ada.

## `POST /api/v1/privileges`

```json
{
  "code": "user.read",
  "name": "Read User",
  "description": "View user data",
  "isActive": true,
  "createdBy": "11111111-1111-4111-8111-111111111111"
}
```

`description`, `isActive`, dan `createdBy` opsional. Jika `createdBy` dikirim, `updatedBy` ikut diisi nilai yang sama.

Response `201` object privilege.

Error: `400 invalid_input`, `409 duplicate_privilege_code`.

## `PUT /api/v1/privileges/{id}`

```json
{
  "code": "user.read",
  "name": "Read Users",
  "description": "View users",
  "isActive": false,
  "updatedBy": "11111111-1111-4111-8111-111111111111"
}
```

`code` dan `name` wajib. Jika `isActive` tidak dikirim, nilai lama tetap dipakai.

Response `200` object privilege.

## `DELETE /api/v1/privileges/{id}`

Hapus permanen. Response `204` tanpa body.

Data awal privilege: jalankan `go run ./cmd/seed privileges`. Lihat [seeder](../seeder.md).
