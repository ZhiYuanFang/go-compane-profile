## 1. Server stream upload

- [x] 1.1 Add OSS upload helper with counting readers + overall weighted pct callback (original then thumb)
- [x] 1.2 Add `POST /admin/api/upload/stream` that validates multipart, streams NDJSON progress/done/error with flush (`X-Accel-Buffering: no`)
- [x] 1.3 Honor request context cancel; best-effort cleanup of partial OSS objects on failure

## 2. Admin client

- [x] 2.1 Implement stream upload client (fetch/readable stream or equivalent) mapping events to `onProgress` / resolve URLs
- [x] 2.2 Point eager upload queue at stream endpoint; keep badge「上传中 N%」bound to OSS-weighted pct
- [x] 2.3 Preserve abort/cancel behavior for in-flight stream uploads

## 3. Deploy + ship

- [x] 3.1 Document nginx `proxy_buffering off` (and timeout guidance) for `/admin/api/upload/stream` in deploy docs
- [x] 3.2 Rebuild admin, sync `resource/public/admin`; smoke progress advancing after browser body 100% until 已上传
