# Users

Manajemen user. ID berupa UUID.

## Object

| Field | Tipe | Keterangan |
| --- | --- | --- |
| `id` | uuid | Diisi server |
| `username` | string | Wajib, tanpa spasi, max 80, unik |
| `email` | string | Wajib, disimpan huruf kecil, unik |
| `name` | string | Wajib, max 120 |
| `isActive` | boolean | Default `true` |
| `createdAt` | datetime | Diisi server |
| `updatedAt` | datetime | Diisi server |
| `createdBy` | uuid | Opsional |
| `updatedBy` | uuid | Opsional |

## `GET /api/v1/users`

Query: `q` (username, email, nama), `limit`, `offset`.

```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-4789-abcd-1234567890ab",
      "username": "rochedi",
      "email": "rochedi@example.com",
      "name": "Rochedi Adha",
      "isActive": true,
      "createdAt": "2026-09-05T08:00:00Z",
      "updatedAt": "2026-09-05T08:00:00Z"
    }
  ],
  "meta": { "total": 1, "limit": 20, "offset": 0 }
}
```

## `GET /api/v1/users/{id}`

`id` harus UUID. Response `200` object user. `400` jika format ID salah. `404` jika tidak ada.

## `POST /api/v1/users`

```json
{
  "username": "rochedi",
  "email": "Rochedi@Example.com",
  "name": "Rochedi Adha",
  "isActive": true,
  "createdBy": "11111111-1111-4111-8111-111111111111"
}
```

`isActive` dan `createdBy` opsional. `createdBy` harus UUID. Jika `createdBy` dikirim, `updatedBy` ikut diisi nilai yang sama.

Response `201` object user.

Error: `400 invalid_input`, `409 duplicate_username`, `409 duplicate_email`.

## `PUT /api/v1/users/{id}`

```json
{
  "username": "rochedi",
  "email": "rochedi@example.com",
  "name": "Rochedi A.",
  "isActive": false,
  "updatedBy": "11111111-1111-4111-8111-111111111111"
}
```

`username`, `email`, dan `name` wajib. Jika `isActive` tidak dikirim, nilai lama tetap dipakai.

Response `200` object user.

## `DELETE /api/v1/users/{id}`

Hapus permanen. Response `204` tanpa body.
