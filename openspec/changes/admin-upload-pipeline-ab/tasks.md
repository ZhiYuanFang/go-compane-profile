## 1. Pipeline helpers

- [x] 1.1 Add compress pool (concurrency 3) with cancel-by-slotId; wire cancel to abort/ignore late results
- [x] 1.2 Replace/refactor `enqueueFileUpload` into compress-then-`enqueuePendingUpload` pipeline; keep upload queue concurrency 3 for HTTP only
- [x] 1.3 Support instant `localPreview` (File object URL) on slot model; revoke on replace/delete/unmount/done as designed

## 2. Gallery UI

- [x] 2.1 `ImageListEditor`: on add/replace/batch set local preview immediately, then start pipeline
- [x] 2.2 `cellPreview` prefers remote → pending → localPreview; allow lightbox when any preview exists while busy
- [x] 2.3 Ensure delete/batch delete cancels both compress and upload for each slotId

## 3. Cover + ship

- [x] 3.1 Align `DualImageField` with shared pipeline helper where cheap (same cancel + preview semantics)
- [x] 3.2 Rebuild admin and sync `resource/public/admin`; smoke batch preview, compress while uploading, delete mid-compress, lightbox while busy
