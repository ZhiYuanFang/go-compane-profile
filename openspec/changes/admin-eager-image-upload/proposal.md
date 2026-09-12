## Why

Saving a portfolio blocks on uploading pending images through the API (often ~1 minute per dual), while gallery badges only show「待上传」and never「已上传」for URL-backed slots. Operators need images to upload in the background as soon as they are added (with a visible queue), so save mostly writes metadata—without the larger change of browser-direct OSS.

## What Changes

- After select/replace/batch-add: compress then **silently upload** via existing `POST /admin/api/upload` (no OSS direct upload)
- Upload **queue with concurrency 3**: earlier slots can be uploading while later ones wait (queued / compressing)
- Per-slot status under each thumb: 处理中 / 排队中 / 上传中 / 已上传 / 失败（可重试）; existing CDN URLs show 已上传
- Stable **slotId** per gallery item; delete (single or batch) removes from queue or **aborts** in-flight upload; late success callbacks ignored if slot gone
- Save blocks or waits while uploads still running / failed; prefer submitting URLs only
- **Original max 3MiB** (client + server); **thumb max aligned** to `0.5 * 1024 * 1024` on both ends (fix server `500*1024` mismatch)
- Cover (`DualImageField`) uses the same eager-upload + status behavior
- Not in scope: STS / PostObject / browser→OSS direct

## Capabilities

### New Capabilities

- `admin-eager-upload-queue`: Background upload queue after image add/replace/batch with concurrency, statuses, and delete/abort safety
- `admin-image-slot-status`: Per-image Simplified Chinese status badges including 已上传 for remote URLs
- `admin-image-size-3m`: Original size ceiling 3MiB and thumb ceiling aligned at 0.5MiB client+server

### Modified Capabilities

- (none)

## Impact

- `admin/src/utils/imageProcess.js`, upload helpers, `ImageListEditor.vue`, `DualImageField.vue`, `PortfolioEditView.vue` (and activity/company if they share DualImageField)
- `internal/service/oss.go` (`MaxOriginalBytes` / `MaxThumbBytes`) and any mirrored validation messages
- Rebuild/sync admin assets; no DB schema change
