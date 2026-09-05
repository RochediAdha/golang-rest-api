# Menus

Manajemen menu hierarkis. ID berupa UUID. Delete bersifat soft delete. `parentId` merujuk ke menu lain.

## Object

| Field | Tipe | Keterangan |
| --- | --- | --- |
| `id` | uuid | Diisi server |
| `parentId` | uuid | Opsional; menu induk yang masih aktif |
| `code` | string | Wajib, tanpa spasi, max 80, unik |
| `name` | string | Wajib, max 120 |
| `path` | string | Opsional, max 255 |
| `icon` | string | Opsional, max 80 |
| `description` | string | Opsional, max 500 |
| `sortOrder` | number | Default `0`, tidak boleh negatif |
| `type` | string | `ITEM`, `GROUP`, atau `COLLAPSE`. Default `ITEM` |
| `isActive` | boolean | Default `true` |
| `createdAt` | datetime | Diisi server |
| `updatedAt` | datetime | Diisi server |
| `deletedAt` | datetime | Terisi setelah soft delete |
| `createdBy` | uuid | Opsional |
| `updatedBy` | uuid | Opsional |
| `deletedBy` | uuid | Opsional |

Urutan list: `sortOrder` naik, lalu `createdAt` naik. Menu yang sudah dihapus tidak muncul di list/detail.

`parentId` tidak boleh:
- menunjuk diri sendiri
- menunjuk menu yang tidak ada / sudah dihapus
- membentuk siklus

## `GET /api/v1/menus`

Query:

| Query | Arti |
| --- | --- |
| `q` | Cari di `code`, `name`, `path` |
| `parentId` | Filter anak. `root`, `null`, atau kosong = menu paling atas |
| `limit` | Default `20` |
| `offset` | Default `0` |

Tanpa `parentId`: semua menu aktif.

```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-4789-abcd-1234567890ab",
      "code": "dashboard",
      "name": "Dashboard",
      "path": "/dashboard",
      "icon": "home",
      "sortOrder": 0,
      "type": "ITEM",
      "isActive": true,
      "createdAt": "2026-09-05T08:00:00Z",
      "updatedAt": "2026-09-05T08:00:00Z"
    }
  ],
  "meta": { "total": 1, "limit": 20, "offset": 0 }
}
```

## `GET /api/v1/menus/{id}`

`id` harus UUID. Response `200` object menu.

## `POST /api/v1/menus`

Menu induk:

```json
{
  "code": "dashboard",
  "name": "Dashboard",
  "path": "/dashboard",
  "icon": "home",
  "description": "Halaman utama",
  "sortOrder": 0,
  "type": "ITEM",
  "isActive": true,
  "createdBy": "11111111-1111-4111-8111-111111111111"
}
```

Submenu:

```json
{
  "parentId": "a1b2c3d4-e5f6-4789-abcd-1234567890ab",
  "code": "users",
  "name": "Users",
  "path": "/users",
  "sortOrder": 1,
  "type": "ITEM"
}
```

Response `201` object menu.

Error: `400 invalid_input`, `400 invalid_parent`, `409 duplicate_menu_code`.

## `PUT /api/v1/menus/{id}`

```json
{
  "parentId": "a1b2c3d4-e5f6-4789-abcd-1234567890ab",
  "code": "users",
  "name": "User Management",
  "path": "/users",
  "icon": "users",
  "sortOrder": 1,
  "type": "ITEM",
  "isActive": true,
  "updatedBy": "11111111-1111-4111-8111-111111111111"
}
```

`code` dan `name` wajib. Jika `isActive` tidak dikirim, nilai lama tetap dipakai. Kirim `parentId` `null` atau kosong untuk jadi menu paling atas.

Response `200` object menu.

## `DELETE /api/v1/menus/{id}`

Soft delete. Gagal jika masih ada anak aktif (`409 menu_has_children`).

Body opsional:

```json
{
  "deletedBy": "11111111-1111-4111-8111-111111111111"
}
```

Response `200` object menu yang baru dihapus, termasuk `deletedAt`.

Data awal menu: jalankan `go run ./cmd/seed menus`. Lihat [seeder](../seeder.md).
