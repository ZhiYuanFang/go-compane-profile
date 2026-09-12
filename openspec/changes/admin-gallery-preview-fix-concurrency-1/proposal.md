## Why

Batch gallery select still shows empty cells until earlier compress/upload slots free, because rapid `patchSlot` updates race on stale `props.modelValue` and drop other slots’ `localPreview`. Concurrently, upload concurrency 3 contributes to gateway **504** timeouts on slow API→OSS paths. Fix preview loss and serialize uploads to concurrency **1**.

## What Changes

- Fix batch (and any multi-slot) preview updates so every selected file keeps an immediate `localPreview` until replaced/done/deleted
- Change shared **upload** and **compress** pool concurrency from 3 → **1**
- Rebuild/sync admin static assets
- No API contract changes; no OSS direct upload; nginx timeout tuning left to ops (optional note only)

## Capabilities

### New Capabilities

- `admin-gallery-batch-preview`: Batch/multi-slot patches must not lose per-slot local previews; all selected files show previews promptly
- `admin-upload-concurrency-1`: Compress and upload queues run with maximum concurrency 1

### Modified Capabilities

- (none in main `openspec/specs/`)

## Impact

- `admin/src/components/ImageListEditor.vue` — stable list patching / batch localPreview seed
- `admin/src/utils/imageUploadQueue.js` — `DEFAULT_CONCURRENCY = 1`
- Build + sync `resource/public/admin`
- Ops: may still raise nginx `proxy_read_timeout` separately if single large uploads hit 60s
