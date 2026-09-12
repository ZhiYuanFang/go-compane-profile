## Why

Admins hit a false “缩略图超过 500KB” error after successful compression because thumb size limits disagree (`500×1024` vs `0.5×1024×1024`). Separately, saving a portfolio with many pending images is very slow because each dual upload runs strictly sequentially on submit.

## What Changes

- Align thumbnail max size with compression: use `0.5 * 1024 * 1024` (same as `compressUnder(0.5, …)`); remove or unify redundant post-checks so false failures stop
- On portfolio (and shared) save-time image resolve: upload pending duals **in parallel with limited concurrency** (default 3) instead of one-by-one `await`
- Show simple upload progress during save (e.g. 已上传 i/n) where the sequential resolve loop lives today
- No API/**BREAKING** contract changes; still `POST /admin/api/upload` per dual
- Out of scope: eager upload-on-select, backend batch upload endpoint

## Capabilities

### New Capabilities

- `admin-image-size-align`: Client image processing uses consistent byte limits for original/thumb vs compress targets
- `admin-upload-concurrency`: Save-time pending dual uploads run with a bounded concurrency pool and progress feedback

### Modified Capabilities

- (none)

## Impact

- `admin/src/utils/imageProcess.js`
- `admin/src/views/PortfolioEditView.vue` (`resolveImages` / submit UX)
- Possibly small shared helper in `admin/src/api/client.js` or `admin/src/utils/` for `mapPool`
- Activity/company single-image resolve may reuse helper later; portfolio multi-image is the primary win
- Rebuild/sync `resource/public/admin`
