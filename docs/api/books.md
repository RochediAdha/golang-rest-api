# Books

Manajemen buku. ID berupa hex 32 karakter. Field waktu memakai `snake_case`.

## Object

| Field | Tipe | Keterangan |
| --- | --- | --- |
| `id` | string | Diisi server |
| `title` | string | Wajib, max 200 karakter |
| `author` | string | Wajib, max 120 karakter |
| `isbn` | string | Opsional; 10 atau 13 digit, unik |
| `year` | number | Opsional; 1000 sampai tahun depan |
| `created_at` | datetime | Diisi server |
| `updated_at` | datetime | Diisi server |

## `GET /api/v1/books`

Query: `q`, `limit`, `offset`.

```json
{
  "data": [
    {
      "id": "a1b2c3d4e5f6789012345678abcdef01",
      "title": "Learning Go",
      "author": "Jon Bodner",
      "isbn": "9781492077213",
      "year": 2021,
      "created_at": "2026-09-05T08:00:00Z",
      "updated_at": "2026-09-05T08:00:00Z"
    }
  ],
  "meta": { "total": 1, "limit": 20, "offset": 0 }
}
```

## `GET /api/v1/books/{id}`

Response `200` object buku. `404` jika tidak ada.

## `POST /api/v1/books`

```json
{
  "title": "Learning Go",
  "author": "Jon Bodner",
  "isbn": "978-1492077213",
  "year": 2021
}
```

Response `201` object buku. Tanda hubung ISBN dihapus otomatis.

Error: `400 invalid_input`, `409 duplicate_isbn`.

## `PUT /api/v1/books/{id}`

Body sama seperti create. `title` dan `author` wajib.

Response `200` object buku.

## `DELETE /api/v1/books/{id}`

Hapus permanen. Response `204` tanpa body.
