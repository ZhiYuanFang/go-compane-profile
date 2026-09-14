## Context

Admin dual-image pipeline today: select file → shared compress queue → immediately `enqueuePendingUpload` → `/admin/api/upload/stream` → OSS → slot `DONE` with CDN URLs. Save is blocked until all slots finish uploading (`dualIsBlockingSave`). Removing a slot after upload but before save does not delete OSS objects (cleanup only diffs URLs already in DB on successful write).

Internal OSS endpoint makes uploads fast. Product decision: keep compress-on-select and progress UI; move HTTP/OSS to save time, site-wide. No leave-page guard.

## Goals / Non-Goals

**Goals:**

- All admin dual-image UIs defer OSS upload until form save
- Select still produces local preview + compressed `pending` original/thumb
- Save flushes pending uploads with existing per-slot progress, then persists CMS payloads
- Unsaved remove/clear creates no OSS objects
- Preserve stream progress UX and concurrency controls already in the queue

**Non-Goals:**

- `beforeunload` / route-leave warnings for pending blobs
- Backend API changes or new delete-on-slot endpoints
- Changing compress size limits or dual-image processing rules
- Historical orphan cleanup / bucket sweeps
- Browser-direct OSS (STS)

## Decisions

### 1. Split pipeline: compress on select, upload on save

- **Choice:** `enqueueFileUpload` / `startSlotPipeline` stops after setting `pending` + ready status; does not call `enqueuePendingUpload`.
- **Save path:** shared helper (e.g. `flushPendingDualUploads(duals, category, patch)`) enqueues pending slots through the existing upload queue + stream client, waits for all to reach `DONE` or any `ERROR`.
- **Alternatives:** compress+upload only on save — rejected (worse preview, longer save, user asked to keep select-time compress). Eager upload — rejected (orphan problem).

### 2. Site-wide, one behavior

- **Choice:** Default deferred everywhere (`DualImageField`, `ImageListEditor`, `PairImageEditor` if still used, company/activity/pricing/portfolio).
- Remove or ignore unused `eagerUpload` prop (always deferred); avoid dual modes.
- **Alternative:** portfolio-only — rejected (user chose full-site; shared components would fork).

### 3. Slot status for “compressed, not uploaded”

- **Choice:** After successful compress, set status to a distinct ready state (prefer new `SLOT_STATUS.READY` or reuse a clear convention: `status` + non-null `pending` without remote URLs). Labels: e.g. “待保存” / tone wait—not “已上传”.
- **Save blocking:** block while `COMPRESSING` / `QUEUED` (compress) / `UPLOADING` / `ERROR`; **do not** block merely because `pending` exists in ready state.
- During save, slots move `READY` → `QUEUED`/`UPLOADING` → `DONE` (or `ERROR`).

### 4. Save transaction shape (client)

- **Choice:** Upload **all** pending duals for the form first; only if every upload succeeds, call create/update/gallery APIs with final URLs. On upload failure, stop, surface error, leave successful slots as `DONE` (URLs in memory) and failed as `ERROR` so retry can re-upload only failures; **do not** write incomplete gallery/cover payloads for that submit.
- **Alternative:** write DB then upload — rejected (DB would reference missing URLs).
- Server-side OSS cleanup of replaced URLs continues to run after successful DB write (unchanged).

### 5. Progress UX during save

- **Choice:** Keep per-slot badges/`uploadProgress` from NDJSON stream. Form `saving` true; primary button label indicates upload+save (e.g. “上传并保存中…”). Optional aggregate copy is nice-to-have, not required.
- Concurrency: keep current shared upload queue limit (today `1` unless changed elsewhere).

### 6. Remove / clear / replace before save

- Cancel in-flight compress; revoke object URLs; drop `pending`; no OSS delete call.
- Replacing a slot that already had remote URLs: new file → compress to new `pending`; old remote URLs stay in memory until save succeeds with new URLs, then server diff deletes old OSS objects (existing cleanup).

### 7. Out of scope for navigation

- No `beforeunload`. Refresh/close may lose pending work and any in-memory `DONE` URLs not yet saved—accepted.

## Risks / Trade-offs

- [Many large pending blobs until save] → Mitigate: compress-on-select already shrinks; concurrency 1 limits parallel memory spikes during upload phase
- [Long save with many gallery images] → Accept; progress per slot; internal network keeps latency low
- [Partial upload success then save abort] → Retry uploads remaining `ERROR`/`pending`; already-`DONE` slots skip re-upload
- [Operator expects instant OSS after select] → Status label “待保存” clarifies; copy on edit pages if needed
- [Regression if any form still assumes URLs before save] → Audit all admin saves that use dual slots

## Migration Plan

1. Ship admin SPA change only (no DB/API migration)
2. Deploy as usual with frontend assets
3. Rollback: revert admin build; eager behavior returns; no data migration
4. Objects uploaded under deferred path only appear after successful save (desired)

## Open Questions

- None blocking (status enum name `READY` vs convention can be decided during implement).
