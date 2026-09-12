## Context

`admin-eager-image-upload` introduced a shared upload queue (concurrency 3) keyed by `slotId`. Gallery path uses `enqueueFileUpload`, which runs **compress + upload inside the same pool job**. Until a job starts, cells have no `pending` blobs, so `cellPreview` is empty and later images look stuck. Cover (`DualImageField`) already compresses first for local preview, then `enqueuePendingUpload`. Product chose explore options **A+B**: instant File preview + compress outside the upload pool; keep upload concurrency 3. Size limits 3MiB / 0.5MiB already shipped.

## Goals / Non-Goals

**Goals:**

- Every selected gallery file shows a local preview immediately
- Compress does not occupy an upload-queue slot; upload queue only does HTTP
- Compress itself is parallel with a small concurrency cap (default 3) so CPU does not melt
- Upload concurrency remains 3; front uploading while back still compressing/queued
- Lightbox allowed when any preview exists, even while busy
- Delete/cancel still abort compress and upload for that `slotId`

**Non-Goals:**

- Dual upload pools beyond compress-pool + upload-pool (option C fully separate orchestration UI)
- Changing 3MiB / 0.5MiB limits
- OSS browser direct upload
- Speeding up single-file OSS latency beyond decoupling wait

## Decisions

1. **Instant preview from selected File (A)**  
   - On add/replace/batch, store a revocable blob URL (or keep `File` + `createObjectURL`) on the slot, e.g. `localPreview` / `previewUrl`, before compress.  
   - `cellPreview` priority: remote thumb → pending thumb → local File preview.  
   - Revoke on replace, delete, unmount, and after remote URLs land (optional once CDN thumb exists).  
   - Alternative: wait for compress for preview — rejected; causes empty cells.

2. **Two-phase pipeline (B)**  
   - Phase 1: `enqueueCompress` (or inline helper) → `processDualImage` → set `pending`, clear compressing status toward upload.  
   - Phase 2: `enqueuePendingUpload` into existing upload queue (concurrency 3).  
   - Replace `enqueueFileUpload`’s single-job compress+upload with this pipeline (keep a thin wrapper for callers).  
   - Alternative: keep compress inside upload job — rejected; blocks display and compress of later slots.

3. **Compress concurrency**  
   - Shared compress pool concurrency **3** (same number as upload unless profiling says otherwise).  
   - Status: File chosen → optional brief idle/preview → `compressing` when compress job runs → `queued`/`uploading` for network.  
   - Slots waiting for compress pool show preview + 「排队中」or 「处理中」only when actively compressing.

4. **UI interaction**  
   - Do not disable thumb click when a preview URL exists (allow lightbox while uploading).  
   - Empty-slot pick still allowed when idle; replace/delete remain available during busy.  
   - Global `picking` only for the brief batch-enqueue setup, not for entire upload lifetime.

5. **Cancel**  
   - Cancel must abort compress job (AbortSignal or cancelled-flag check around `processDualImage`) and upload fetch for that `slotId`.  
   - Late compress/upload results ignored if slot gone or cancelled.

## Risks / Trade-offs

- [Many blob URLs for large batches] → Revoke aggressively; cache per blob; limit retention after DONE  
- [CPU spike from parallel compress] → Cap compress concurrency at 3  
- [Status flicker queued→compressing→queued→uploading] → Accept; labels already exist  
- [DualImageField vs list drift] → Prefer shared `startSlotPipeline(file)` helper used by both  

## Migration Plan

1. Ship admin JS only (no API change)  
2. Rebuild/sync `resource/public/admin`  
3. Rollback: previous admin build  

## Open Questions

- None — A+B and concurrency 3 locked.
