## Context

Gallery/cover images are compressed client-side then uploaded only on portfolio save via `POST /admin/api/upload` (browser→API→OSS), which feels like「上传中 0/1」stalls. Badges only show「待上传」when `pending` is set. Product locked: **eager silent upload** through the existing API (not OSS direct), **queue with concurrency 3**, per-slot statuses, **delete/abort safety**, **original max 3MiB**, thumb **0.5MiB** aligned client+server.

## Goals / Non-Goals

**Goals:**

- Upload ASAP after add/replace/batch with visible queue (front uploading, back waiting)
- Clear per-slot status including 已上传 for URL-backed images
- Safe delete/cancel during queued or in-flight work
- Save primarily persists URLs; gate on incomplete/failed uploads
- Shrink original ceiling to 3MiB end-to-end

**Non-Goals:**

- Browser-direct OSS (STS / PostObject / presigned PUT)
- Changing CDN URL shape or gallery DB schema
- Guaranteed immediate OSS delete on remove (orphan cleanup may remain)

## Decisions

1. **Eager upload via existing API**  
   - Reuse `uploadDualImage` / multipart to `/admin/api/upload`.  
   - Alternative: OSS direct — rejected as too large for this change.

2. **Queue + concurrency 3**  
   - Shared pool: each job = compress (if needed) + upload, or compress first then queue upload—prefer **one job occupying a pool slot for compress+upload** for simplicity.  
   - Batch add enqueues many jobs; at most 3 run; others show 排队中.  
   - Reuse/extend `mapPool` or a small dedicated upload queue module keyed by `slotId`.

3. **Stable `slotId`**  
   - Generate on create (e.g. `crypto.randomUUID()` / guid).  
   - Never use array index as cancel key.  
   - Dual model: `{ slotId, thumb, original, pending, status, error? }` (status can be derived but explicit is clearer for UI).

4. **Delete / abort**  
   - Queued: remove job, never start.  
   - Compressing/uploading: `AbortController.abort()` on fetch if supported; mark slot removed.  
   - Completion callback: if `slotId` not in list, discard (may leave OSS orphan).  
   - Batch delete: same per selected slotId.

5. **Save behavior**  
   - If any slot is compressing/queued/uploading: wait until idle **or** block with「仍有图片上传中」 (prefer **wait with timeout message**, or disable save until idle—**disable/block with clear copy** is simpler and safer).  
   - If any `error`: block save until retry/remove.  
   - `resolveDualImage` on save becomes mostly no-op for already-uploaded URLs; keep as fallback for leftover pending.

6. **Size limits**  
   - Client: `ORIGINAL_MAX_MB = 3`, `THUMB_MAX_MB = 0.5`.  
   - Server: `MaxOriginalBytes = 3<<20`, `MaxThumbBytes = 0.5 * 1024 * 1024` (align with client).  
   - Update error strings accordingly.

7. **Cover**  
   - DualImageField participates in same eager path and shows status text.

## Risks / Trade-offs

- [Abort may not stop server/OSS after request left browser] → Accept orphans; cleanup job  
- [Leave page mid-upload] → Orphans / lost pending; acceptable  
- [CPU+network contention] → Cap concurrency at 3  
- [Save disabled while uploading] → Better than partial URL lists  

## Migration Plan

1. Deploy API with 3MiB / 0.5MiB limits  
2. Deploy admin with eager queue UI  
3. Rollback: previous builds; stricter 3MiB may reject older 5MiB clients briefly if mixed  

## Open Questions

- None for implementation; product choices locked.
