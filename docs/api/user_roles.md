# User roles

Penugasan role ke user. ID berupa UUID. Hapus bersifat permanen.

Satu pasangan `userId` + `roleId` hanya boleh ada sekali. Jika user atau role dihapus fisik, baris terkait ikut terhapus (`ON DELETE CASCADE`). Role yang di-soft-delete tidak bisa ditugaskan.

`number` dipakai sebagai identitas join (NIM mahasiswa / NIP dosen). Wajib untuk role `dosen`, `lecturer`, `mahasiswa`, atau `student`. Unik per role: satu NIM tidak boleh dipakai dua mahasiswa.

## Object create / list

| Field | Tipe | Keterangan |
| --- | --- | --- |
| `id` | uuid | ID penugasan, diisi server |
| `userId` | uuid | Wajib, harus merujuk user yang ada |
| `roleId` | uuid | Wajib, harus merujuk role yang belum dihapus |
| `number` | string | Opsional kecuali role dosen/mahasiswa; max 50; unik per role |
| `createdAt` | datetime | Diisi server |
| `createdBy` | uuid | Opsional |

Tidak ada `updatedAt` dan tidak ada endpoint update. Ubah penugasan dengan hapus lalu create ulang.

## Object show

`GET /api/v1/user-roles/{id}` menampilkan satu user beserta semua role-nya. `{id}` boleh `userId` atau ID penugasan.

| Field | Tipe | Keterangan |
| --- | --- | --- |
| `id` | uuid | ID user |
| `userId` | uuid | ID user |
| `name` | string | Nama user |
| `roles` | array | Daftar role user tersebut |

Tiap item `roles`:

| Field | Tipe | Keterangan |
| --- | --- | --- |
| `id` | uuid | ID penugasan (untuk delete) |
| `roleId` | uuid | ID role |
| `name` | string | Nama role |
| `number` | string | NIM/NIP jika ada |
| `createdAt` | datetime | Waktu penugasan |
| `createdBy` | uuid | Opsional |

## `GET /api/v1/user-roles`

Query: `userId`, `roleId`, `number`, `limit`, `offset`.

`number` dipakai untuk mencari penugasan saat join, misalnya `GET /api/v1/user-roles?number=2301001`.

```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-4789-abcd-1234567890ab",
      "userId": "11111111-1111-4111-8111-111111111111",
      "roleId": "22222222-2222-4222-8222-222222222222",
      "number": "2301001",
      "createdAt": "2026-09-05T08:00:00Z"
    }
  ],
  "meta": { "total": 1, "limit": 20, "offset": 0 }
}
```

## `GET /api/v1/user-roles/{id}`

`id` harus UUID user atau ID penugasan. Response `200` object show. `400` / `404` jika tidak valid atau user belum punya role.

```json
{
  "id": "11111111-1111-4111-8111-111111111111",
  "userId": "11111111-1111-4111-8111-111111111111",
  "name": "Rochedi Adha",
  "roles": [
    {
      "id": "a1b2c3d4-e5f6-4789-abcd-1234567890ab",
      "roleId": "22222222-2222-4222-8222-222222222222",
      "name": "admin",
      "createdAt": "2026-09-05T08:00:00Z"
    },
    {
      "id": "b2c3d4e5-f6a7-4890-bcde-234567890abc",
      "roleId": "33333333-3333-4333-8333-333333333333",
      "name": "mahasiswa",
      "number": "2301001",
      "createdAt": "2026-09-05T08:05:00Z"
    }
  ]
}
```

## `POST /api/v1/user-roles`

```json
{
  "userId": "11111111-1111-4111-8111-111111111111",
  "roleId": "22222222-2222-4222-8222-222222222222",
  "number": "2301001",
  "createdBy": "11111111-1111-4111-8111-111111111111"
}
```

`createdBy` opsional. `number` wajib jika role adalah dosen/mahasiswa.

Response `201` object penugasan (satu role).

Error: `400 invalid_input`, `400 invalid_user`, `400 invalid_role`, `409 duplicate_user_role`, `409 duplicate_user_role_number`.

## `DELETE /api/v1/user-roles/{id}`

Hapus satu penugasan. `{id}` adalah ID penugasan (`roles[].id`), bukan `userId`. Response `204` tanpa body.
