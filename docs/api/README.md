# Spec API

Base URL default: `http://localhost:8080`

Semua request/response body memakai `Content-Type: application/json`.

`GET /` menampilkan daftar path modul. `GET /health` mem-ping database.

| Modul | File |
| --- | --- |
| Health | [health.md](health.md) |
| Users | [users.md](users.md) |
| Roles | [roles.md](roles.md) |
| Menus | [menus.md](menus.md) |
| User roles | [user_roles.md](user_roles.md) |
| Privileges | [privileges.md](privileges.md) |
| Role privileges | [role_privileges.md](role_privileges.md) |

## Konvensi

- List memakai pagination: `limit` (default `20`, max `100`) dan `offset` (default `0`).
- Pencarian teks lewat query `q`.
- Field yang tidak dikirim atau kosong pada create/update mengikuti aturan tiap modul.
- Belum ada autentikasi.

## Format list

```json
{
  "data": [],
  "meta": {
    "total": 0,
    "limit": 20,
    "offset": 0
  }
}
```

## Format error

```json
{
  "error": {
    "code": "invalid_input",
    "message": "request body or parameters are invalid"
  }
}
```

| HTTP | `error.code` | Arti |
| --- | --- | --- |
| `400` | `invalid_input` | Body atau parameter tidak valid |
| `400` | `invalid_parent` | `parentId` menu tidak valid |
| `400` | `invalid_user` | `userId` tidak merujuk user yang ada |
| `400` | `invalid_role` | `roleId` tidak merujuk role yang ada |
| `400` | `invalid_menu` | `menuId` tidak merujuk menu yang ada |
| `400` | `invalid_privilege` | `privilegeId` tidak merujuk privilege yang ada |
| `404` | `not_found` | Data tidak ditemukan |
| `409` | `duplicate_username` | Username sudah dipakai |
| `409` | `duplicate_email` | Email sudah dipakai |
| `409` | `duplicate_role_name` | Nama role sudah dipakai |
| `409` | `duplicate_menu_code` | Kode menu sudah dipakai |
| `409` | `menu_has_children` | Menu masih punya anak |
| `409` | `duplicate_user_role` | User sudah punya role ini |
| `409` | `duplicate_user_role_number` | `number` sudah dipakai untuk role ini |
| `409` | `duplicate_privilege_code` | Kode privilege sudah dipakai |
| `409` | `duplicate_role_privilege` | Role sudah punya privilege ini pada menu tersebut |
| `500` | `internal_error` | Error server |
| `503` | — | Database tidak tersedia (`GET /health`) |
