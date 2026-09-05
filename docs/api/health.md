# Health

Cek status aplikasi dan koneksi PostgreSQL.

## `GET /health`

Tidak lewat usecase. Handler langsung ping connection pool.

### Response `200`

```json
{
  "status": "ok",
  "database": "ok",
  "time": "2026-09-05T08:00:00Z"
}
```

### Response `503`

Database tidak merespons.

```json
{
  "status": "degraded",
  "database": "unavailable",
  "time": "2026-09-05T08:00:00Z"
}
```

## `GET /`

Informasi singkat API.

```json
{
  "name": "golang-rest-api",
  "version": "v1",
  "docs": {
    "health": "GET /health",
    "users": "GET /api/v1/users",
    "roles": "GET /api/v1/roles",
    "menus": "GET /api/v1/menus"
  }
}
```
