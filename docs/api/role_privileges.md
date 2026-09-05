# Role privileges

Penugasan privilege ke role pada suatu menu. Dipakai sebagai dasar otorisasi. ID berupa UUID. Hapus bersifat permanen.

Satu kombinasi `roleId` + `menuId` + `privilegeId` hanya boleh ada sekali. Jika role, menu, atau privilege dihapus fisik, baris terkait ikut terhapus (`ON DELETE CASCADE`). Role atau menu yang di-soft-delete tidak bisa ditugaskan.

## Object

| Field | Tipe | Keterangan |
| --- | --- | --- |
| `id` | uuid | Diisi server |
| `roleId` | uuid | Wajib, harus merujuk role yang belum dihapus |
| `menuId` | uuid | Wajib, harus merujuk menu yang belum dihapus |
| `privilegeId` | uuid | Wajib, harus merujuk privilege yang ada |
| `createdAt` | datetime | Diisi server |
| `createdBy` | uuid | Opsional |

Tidak ada `updatedAt` dan tidak ada endpoint update. Ubah penugasan dengan hapus lalu create ulang.

## `GET /api/v1/role-privileges`

Query: `roleId`, `menuId`, `privilegeId`, `limit`, `offset`.

```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-4789-abcd-1234567890ab",
      "roleId": "11111111-1111-4111-8111-111111111111",
      "menuId": "22222222-2222-4222-8222-222222222222",
      "privilegeId": "33333333-3333-4333-8333-333333333333",
      "createdAt": "2026-09-05T08:00:00Z"
    }
  ],
  "meta": { "total": 1, "limit": 20, "offset": 0 }
}
```

## `GET /api/v1/role-privileges/{id}`

`id` harus UUID. Response `200` object role privilege. `400` / `404` jika tidak valid atau tidak ada.

## `POST /api/v1/role-privileges`

```json
{
  "roleId": "11111111-1111-4111-8111-111111111111",
  "menuId": "22222222-2222-4222-8222-222222222222",
  "privilegeId": "33333333-3333-4333-8333-333333333333",
  "createdBy": "11111111-1111-4111-8111-111111111111"
}
```

`createdBy` opsional.

Response `201` object role privilege.

Error: `400 invalid_input`, `400 invalid_role`, `400 invalid_menu`, `400 invalid_privilege`, `409 duplicate_role_privilege`.

## `DELETE /api/v1/role-privileges/{id}`

Hapus permanen. Response `204` tanpa body.
