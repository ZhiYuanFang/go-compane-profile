## Why

Eager upload already queues gallery images with concurrency 3, but compress and upload share one pool slot. Later images stay empty (no preview) until they acquire a slot, so uploading feels like it blocks other images’ processing and display. Cover (`DualImageField`) already compresses then uploads; gallery should match that UX.

## What Changes

- **A — Instant preview**: On add/replace/batch, show a local preview from the selected `File` immediately (before compress finishes).
- **B — Decouple compress from upload queue**: Compress outside the upload pool (with its own concurrency limit); enqueue only network upload via `enqueuePendingUpload`. Upload concurrency remains 3.
- Busy slots with a preview remain viewable (lightbox); delete/cancel still abort queue work.
- Original max 3MiB / thumb 0.5MiB stay as already shipped — no further size-limit change in this change.
- No OSS browser direct upload.

## Capabilities

### New Capabilities

- `admin-gallery-instant-preview`: Gallery cells show local File/blob preview as soon as a file is chosen, independent of queue position.
- `admin-compress-upload-decouple`: Compress runs in a separate limited pool; the shared upload queue only performs HTTP upload of already-compressed pending blobs.

### Modified Capabilities

- (none in main `openspec/specs/`; builds on `admin-eager-image-upload` behavior)

## Impact

- `admin/src/utils/imageUploadQueue.js` — split or add compress pool; `enqueueFileUpload` becomes compress-then-pending-upload or replaced by a two-phase helper
- `admin/src/components/ImageListEditor.vue` — instant preview field/cache; don’t wait for compress for cell display; relax thumb disable when preview exists
- `admin/src/components/DualImageField.vue` — already close to B; may reuse shared compress helper
- `admin/src/utils/dualSlot.js` — optional `localPreview` / status semantics if needed
- Build + sync `resource/public/admin`
- No API / OSS contract changes
