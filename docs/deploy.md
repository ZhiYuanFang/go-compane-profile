# Deploy runbook — go-compane-profile

## Endpoints

| Surface | URL |
|---------|-----|
| Public API | `https://muhou.cuplay.top/api/v1/*` |
| Admin SPA | `https://muhou.cuplay.top/admin/` |
| Admin API | `https://muhou.cuplay.top/admin/api/*` |
| CDN images | `https://resorce.cuplay.top/muhou/...` |
| Container port | **9100** |

## Secrets (`.env.prod`)

Copy from `manifest/docker/env/.env.example` to `manifest/docker/env/.env.prod` on the server (gitignored). Fill:

- `MYSQL_LINK` → `compane_profile`（密码保持原始字符，勿 URL 编码 `+`）
- `OSS_*` / `CDN_BASE` → same as go_ai_talk ucg (`pang-bao`, prefix `muhou/`)
- `ADMIN_PASSWORD` / `ADMIN_SESSION_SECRET`
- `REGISTRY` / `IMAGE_TAG` for compose pull

Never commit real `.env.prod`.

## Docker Compose

Compose project name remains `go-compane-profile`; the application **service key** is `muhou` (`container_name: muhou`).

```bash
cd manifest/docker
docker compose --env-file env/.env.prod -f docker-compose.yml pull
docker compose --env-file env/.env.prod -f docker-compose.yml up -d --no-build
```

Useful: `docker compose --env-file env/.env.prod logs -f muhou`

If renaming from an older service key (`compane-profile`), run `down` once before `up` to avoid container-name conflicts.

Container listens on host **9100**; env vars come from `env/.env.prod`.

## ACR

- Namespace: `compane-profile`
- Image repo: `muhou` → `{REGISTRY}/muhou:{IMAGE_TAG}` (e.g. `compane-profile/muhou:v0.0.1`)
- Example pull (VPC):  
  `crpi-lff3xynwzvqxxxjk-vpc.cn-hangzhou.personal.cr.aliyuncs.com/compane-profile/muhou:v0.0.1`
- CI: push tag `v0.0.1` or run workflow_dispatch; configure GitHub Secrets `ACR_USERNAME` / `ACR_PASSWORD` and Environment `REGISTRY`.

## Nginx sketch

```nginx
server {
  listen 443 ssl;
  server_name muhou.cuplay.top;
  # ssl_certificate ...;

  # Streaming OSS upload progress (NDJSON). Do not buffer the response body.
  location = /admin/api/upload/stream {
    proxy_pass http://127.0.0.1:9100;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    client_max_body_size 12m;
    proxy_buffering off;
    proxy_cache off;
    proxy_http_version 1.1;
    proxy_read_timeout 300s;
    proxy_send_timeout 300s;
  }

  location / {
    proxy_pass http://127.0.0.1:9100;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    client_max_body_size 12m;
    proxy_read_timeout 300s;
  }
}
```

The stream handler also sends `X-Accel-Buffering: no`. After deploy, confirm progress keeps advancing **after** the browser finishes sending the multipart body (OSS put phase).

## WeChat mini-program domains

| Type | Domain |
|------|--------|
| request 合法域名 | `muhou.cuplay.top` |
| downloadFile 合法域名 | `resorce.cuplay.top` |

## First-time DB

```bash
mysql -h ... -u root -p < manifest/sql/schema.sql
mysql -h ... -u root -p compane_profile < manifest/sql/seed_company.sql
```

Seed inserts company copy only; upload logos / portfolios / pricing via `/admin/`.
