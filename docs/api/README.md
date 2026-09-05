# Spec API

Base URL default: `http://localhost:8080`

Semua request/response body memakai `Content-Type: application/json`.

| Modul | File |
| --- | --- |
| Health | [health.md](health.md) |
| Books | [books.md](books.md) |
| Users | [users.md](users.md) |
| Roles | [roles.md](roles.md) |
| Menus | [menus.md](menus.md) |

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
| `404` | `not_found` | Data tidak ditemukan |
| `409` | `duplicate_isbn` | ISBN buku sudah dipakai |
| `409` | `duplicate_username` | Username sudah dipakai |
| `409` | `duplicate_email` | Email sudah dipakai |
| `409` | `duplicate_role_name` | Nama role sudah dipakai |
| `409` | `duplicate_menu_code` | Kode menu sudah dipakai |
| `409` | `menu_has_children` | Menu masih punya anak |
| `500` | `internal_error` | Error server |
| `503` | — | Database tidak tersedia (`GET /health`) |
