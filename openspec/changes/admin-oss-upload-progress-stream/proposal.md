## Why

Browser→API upload progress hits 100% in seconds, then the admin waits a long time with no feedback while the API uploads original + thumb to Aliyun OSS. Operators need **real progress for the OSS write phase**, which requires the server to stream progress while putting objects—not browser `xhr.upload` alone.

## What Changes

- Add a **streaming upload** path: after receiving the dual multipart body, the server uploads to OSS while **pushing progress events** to the client (NDJSON/SSE-style chunked response preferred over a separate WebSocket for this app).
- Admin eager-upload queue consumes stream events and shows overall **OSS-weighted** percent (and clear phase labels when useful).
- Wrap OSS `PutObject` readers with byte counters so progress reflects bytes written toward OSS, not just “request finished”.
- Document nginx/`X-Accel-Buffering` so proxies do not buffer the stream.
- Keep or thin-wrap the existing JSON `/admin/api/upload` for compatibility during migration; admin SPA switches to the stream endpoint.
- **Not** browser-direct OSS (STS) in this change.

## Capabilities

### New Capabilities

- `admin-oss-progress-stream`: Server streams OSS upload progress; admin shows truthful progress through the OSS phase until CDN URLs are returned

### Modified Capabilities

- (none in main `openspec/specs/`)

## Impact

- `internal/service/oss.go` — progress-aware upload (counting reader / phased Put)
- `internal/controller/admin` + `api/admin/v1` — new stream upload endpoint (or streaming variant of upload)
- `admin/src/api/client.js` — fetch/stream consumer with progress callbacks
- `admin/src/utils/imageUploadQueue.js` + `dualSlot.js` — map stream progress into slot UI
- Deploy docs: proxy buffering off for the stream path
- Nginx / compose if needed for buffering headers
