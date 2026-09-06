# Deploy di server

PostgreSQL berjalan di server terpisah. Image Docker hanya berisi API (dan binary seeder).

Cara menjalankan API:

1. **Docker Compose** — container API saja, konek ke Postgres lewat `DB_*` di `.env`
2. **systemd** — binary dijalankan langsung, variabel yang sama

Migrasi jalan otomatis saat API start. Seeder tetap manual.

## Persiapan Postgres (server DB)

Buat user dan database di server PostgreSQL:

```sql
CREATE USER api WITH PASSWORD 'password-kuat';
CREATE DATABASE golang_rest_api OWNER api;
```

Izinkan server API konek (`pg_hba.conf` + `listen_addresses`). Port `5432` hanya untuk IP server API, bukan publik.

## Persiapan API

- Server Linux (Ubuntu/Debian)
- Domain opsional (Nginx + HTTPS)
- Salin `.env.example` menjadi `.env`

```bash
cp .env.example .env
```

Isi host server DB di `DB_HOST`, bukan `localhost` (di dalam container, `localhost` adalah container itu sendiri):

```
ADDR=:8080
DB_HOST=db.contoh.com
DB_PORT=5432
DB_USER=api
DB_PASSWORD=password-kuat
DB_NAME=golang_rest_api
DB_SSLMODE=require
DB_MAX_CONNS=10
DB_MIN_CONNS=1
```

Jangan tulis `${VAR}` di `.env`; aplikasi tidak men-expand interpolasi. `DATABASE_URL` opsional: jika diisi, mengalahkan `DB_*`.

`DB_SSLMODE=require` disarankan antar server. Port API di-bind ke `127.0.0.1:8080`; dari luar lewat Nginx.

## 1. Docker Compose

### Instal Docker

```bash
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER
```

Logout/login agar grup `docker` aktif. Cek: `docker compose version`.

### Clone dan jalankan

```bash
sudo mkdir -p /opt/golang-rest-api
sudo chown $USER:$USER /opt/golang-rest-api
cd /opt/golang-rest-api
git clone <url-repo> .
cp .env.example .env
```

Edit `.env` agar `DB_HOST` mengarah ke server Postgres.

```bash
docker compose up -d --build
curl http://127.0.0.1:8080/health
```

Response `200` dengan `"database": "ok"` berarti API sudah konek ke Postgres.

### Seeder

```bash
docker compose run --rm --entrypoint /app/seed api roles
docker compose run --rm --entrypoint /app/seed api menus
docker compose run --rm --entrypoint /app/seed api privileges
```

atau `make docker-seed name=roles`. Data yang sudah ada dilewati.

### Update

```bash
cd /opt/golang-rest-api
git pull
docker compose up -d --build
```

Migrasi baru ikut jalan saat container API start.

### Log dan status

```bash
docker compose ps
docker compose logs -f api
```

Backup dilakukan di server Postgres (`pg_dump`), bukan dari Compose.

## 2. systemd (tanpa Docker)

`.env` di `/opt/golang-rest-api` memakai `DB_*` yang sama ke server DB.

### Build

Install Go 1.25, atau build di laptop lalu copy binary.

Di server:

```bash
cd /opt/golang-rest-api
go test ./...
go build -o /opt/golang-rest-api/bin/api ./cmd/api
go build -o /opt/golang-rest-api/bin/seed ./cmd/seed
```

Dari laptop:

```bash
GOOS=linux GOARCH=amd64 go build -o api ./cmd/api
GOOS=linux GOARCH=amd64 go build -o seed ./cmd/seed
scp api seed user@server:/opt/golang-rest-api/bin/
```

### Unit

`/etc/systemd/system/golang-rest-api.service`:

```ini
[Unit]
Description=Golang REST API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/golang-rest-api
ExecStart=/opt/golang-rest-api/bin/api
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

`WorkingDirectory` wajib agar file `.env` terbaca.

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now golang-rest-api
curl http://127.0.0.1:8080/health
```

Seeder:

```bash
cd /opt/golang-rest-api
./bin/seed roles
./bin/seed menus
./bin/seed privileges
```

Update: `git pull`, build ulang binary, `sudo systemctl restart golang-rest-api`.

## Nginx + HTTPS

Sama untuk Compose maupun systemd. API tetap di `127.0.0.1:8080`.

```nginx
server {
    listen 80;
    server_name api.contoh.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
sudo apt install -y nginx certbot python3-certbot-nginx
sudo ln -s /etc/nginx/sites-available/golang-rest-api /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
sudo certbot --nginx -d api.contoh.com
```

Jangan buka port `8080` ke publik. Port Postgres hanya dari server API.

## Ringkasan

| Hal | Perilaku |
| --- | --- |
| Postgres | Server terpisah, lewat `DB_HOST` / `DB_*` |
| Migrasi | Otomatis saat API start |
| Seeder | Manual, idempotent |
| Health | `GET /health` |
| Compose | `docker compose up -d --build` |
| systemd | `systemctl enable --now golang-rest-api` |
| Secret | `.env` tidak di-commit |
