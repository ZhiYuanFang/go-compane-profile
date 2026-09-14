## Why

Admin currently uploads dual images to OSS immediately on file select. Operators who remove a slot before saving leave orphan objects in the bucket. Internal-network OSS is now fast enough that eager upload is unnecessary; deferring HTTP upload until save keeps progress UX while eliminating pre-save orphans.

## What Changes

- Stop auto-uploading to OSS on image select across all admin CMS forms (company, portfolio cover/gallery, activity, pricing, and any other dual-image fields).
- On select: keep instant local preview + compress into `pending` blobs only; do not call `/admin/api/upload/stream` until save.
- On save: upload all pending slots (retain existing per-slot OSS progress UI), then write DB payloads with CDN URLs; fail save if any upload fails (no partial write of image URLs for that submit).
- Removing or clearing a not-yet-saved slot discards local blobs only—no delete API / no OSS traffic.
- Compress-on-select remains; only the HTTP/OSS phase moves to save time.
- Out-of-scope: `beforeunload` / leave-page warnings for unsaved pending images.
- Server upload endpoint and OSS orphan cleanup on successful replace/delete remain unchanged.

## Capabilities

### New Capabilities

- `admin-deferred-upload`: Select compresses locally; OSS upload runs at save with progress; unsaved removals create no OSS objects

### Modified Capabilities

- (none in main `openspec/specs/`; this supersedes prior eager-upload admin behavior from archived/experimental changes)

## Impact

- `admin/src/utils/imageUploadQueue.js` — split compress vs upload; select stops after `pending`
- `admin/src/utils/dualSlot.js` — ready/pending status labels; save blocking rules
- `admin/src/components/DualImageField.vue`, `ImageListEditor.vue`, related editors — wire deferred path; honor/remove unused `eagerUpload`
- Admin save handlers: `PortfolioEditView`, `CompanyView`, `ActivityEditView`, pricing (and any other dual-image saves) — flush pending uploads before API write
- `admin/src/api/client.js` — reuse stream upload + optional shared “flush pendings” helper; `resolveDualImage` may be central
- Backend: no required API changes (existing stream upload + CMS OSS cleanup)
