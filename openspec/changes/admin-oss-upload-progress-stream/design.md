## Context

Today `POST /admin/api/upload` buffers original+thumb, then `UploadDualImageBytes` runs two sequential `PutObject` calls and returns JSON once. Admin XHR `upload.onprogress` only tracks browser→API bytes, so the UI reaches 100% while OSS work still runs (often the slow part, including 504 risk). Product chose option ②: **server pushes progress while writing OSS**, not browser-direct OSS.

## Goals / Non-Goals

**Goals:**

- Stream progress events while original and thumb are uploaded to OSS
- Admin badge shows meaningful overall % through the OSS phase until URLs land
- Abort still cancels in-flight work when possible
- Proxy-friendly headers so chunks flush to the client

**Non-Goals:**

- STS / browser PostObject / multipart direct to OSS
- WebSocket binary upload (prefer HTTP stream on the same POST)
- Perfect byte-for-byte OSS network accounting (SDK may buffer; counting reader on Put body is “good enough”)
- Changing dual size limits (3MiB / 0.5MiB)

## Decisions

1. **Transport: chunked NDJSON on POST (not WebSocket)**  
   - Client `POST multipart` to e.g. `/admin/api/upload/stream` (auth cookie same as today).  
   - After multipart is parsed and validated, response is `Content-Type: application/x-ndjson` (or `text/event-stream` if SSE framing is easier in GoFrame—prefer **NDJSON lines** for simpler fetch parsing).  
   - Events:
     - `{"type":"progress","pct":0-100,"phase":"oss_original"|"oss_thumb"}`
     - `{"type":"done","original":"<cdn>","thumb":"<cdn>"}`
     - `{"type":"error","message":"..."}`  
   - Alternative WebSocket rejected for this change: more infra/auth surface; NDJSON-on-POST matches existing FormData upload.

2. **Progress source**  
   - Do not wait until both Puts finish to speak.  
   - Upload with `io.Reader` wrappers that report `read/total` as OSS SDK reads the body for `PutObject`.  
   - Overall `pct` = weighted by known sizes:  
     `pct = round(100 * (origRead + thumbRead) / (origSize + thumbSize))`  
     with phases for labeling.  
   - Flush after each progress tick (throttle ~100–200ms or every ≥1% change).

3. **Client consumption**  
   - Prefer `fetch` + `ReadableStream` (or XHR with incremental `responseText` parse) for the stream endpoint.  
   - Map `progress.pct` → `uploadProgress`; keep status `uploading` until `done`.  
   - Optional label: still `上传中 N%` (now OSS-weighted). After body received and before first OSS event, may show `写入中…` / `pct:0` briefly.  
   - Browser→API phase: optional secondary signal (keep light); **OSS stream % is the primary badge number**.

4. **Compat**  
   - Keep existing `/admin/api/upload` JSON endpoint for rollback/scripts.  
   - Admin SPA uses `/upload/stream` only.

5. **Proxy**  
   - Set `X-Accel-Buffering: no` on stream responses.  
   - Document nginx: `proxy_buffering off` for `/admin/api/upload/stream` (and raise `proxy_read_timeout` if still needed).

6. **Cancel**  
   - Client abort closes the request; handler should stop further Puts when context cancelled.  
   - Partial OSS objects: best-effort delete on failure (same as thumb-fail cleanup today).

## Risks / Trade-offs

- [Proxy buffers NDJSON → UI stuck at 0%] → Explicit buffering headers + deploy doc  
- [SDK reads ahead → progress jumps] → Accept; throttle updates  
- [GoFrame may wrap writers] → May need raw `http.ResponseWriter` flush (`http.Flusher`)  
- [Long OSS still 504] → Progress does not remove gateway timeouts; still raise timeouts / speed Puts  

## Migration Plan

1. Ship API stream endpoint + flush headers  
2. Ship admin consumer + rebuild `resource/public/admin`  
3. Tune nginx buffering/timeouts  
4. Rollback: point SPA back to JSON `/upload`  

## Open Questions

- None for start—NDJSON-on-POST locked; WebSocket out of scope.
