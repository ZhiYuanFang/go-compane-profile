## Context

`processDualImage` compresses thumbs with `compressUnder(file, 0.5, 1280)` (success ≤ `0.5*1024*1024`) but then rejects with `THUMB_MAX = 500*1024`, causing false errors in the ~512–524KB band. Portfolio save resolves pending gallery duals via sequential `await resolveDualImage` in `PortfolioEditView`, so N uploads take ~N× latency. Product locked: fix size align (A) + parallel uploads with concurrency limit 3; not eager upload-on-select.

## Goals / Non-Goals

**Goals:**

- Single consistent thumb/original byte ceiling matching compress targets
- Parallel pending uploads on save with concurrency ≤ 3
- Visible progress while uploading many images
- Preserve upload API and dual-image model

**Non-Goals:**

- Upload-on-select / background pre-upload
- New batch upload API or multipart multi-file endpoint
- Changing OSS server code
- Parallelizing client-side compression in this change (optional later)

## Decisions

1. **Thumb max = `0.5 * 1024 * 1024`**  
   - Set `THUMB_MAX` equal to compress threshold.  
   - Keep post-check or drop it; if kept, must use the same constant. Prefer one source of truth: `const THUMB_MAX_MB = 0.5` and derive bytes.  
   - Original already aligned at `5 * 1024 * 1024`.

2. **Concurrency pool default 3**  
   - Small helper `mapPool(items, concurrency, fn)` (or inline) used by `resolveImages`.  
   - Skip empty duals before pooling; preserve **output order** matching input list order (important for gallery sort).  
   - Alternative: unlimited `Promise.all` — rejected (can overwhelm browser/OSS).

3. **Progress UX on portfolio save**  
   - While resolving: show「上传中 i/n」(or similar) on the save button / message area.  
   - `n` = count of items that need upload (pending); already-uploaded URLs don’t count.

4. **Failure behavior**  
   - If any upload fails, reject the whole save (same as today); no partial gallery write for that submit.  
   - In-flight uploads may complete; orphan OSS cleanup remains existing ops concern.

5. **Reuse**  
   - Helper can live in `utils/asyncPool.js` or next to `client.js`; portfolio is required; activity/company single image unchanged but can call the same resolve helper unchanged.

## Risks / Trade-offs

- [Higher concurrent bandwidth] → Cap at 3  
- [Order bugs if pool maps unordered] → Explicit index → result array  
- [Progress flicker] → Only update on each completion  

## Migration Plan

1. Deploy admin build only  
2. Rollback previous admin assets  

## Open Questions

- None; A + concurrency 3 locked.
