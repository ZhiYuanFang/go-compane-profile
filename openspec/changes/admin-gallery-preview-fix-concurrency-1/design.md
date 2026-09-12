## Context

After `admin-upload-pipeline-ab`, gallery slots set `localPreview` immediately then compress/upload in separate pools (was concurrency 3). Batch select still shows empty cells because `ImageListEditor.patchSlot` reads `props.modelValue` and emits on every patch; synchronous multi-slot patches race and overwrite each other’s `localPreview`. Separately, concurrency 3 on upload through Nginx→API→OSS correlates with **504** gateway timeouts. Product locked: fix batch preview loss + set concurrency to **1**.

## Goals / Non-Goals

**Goals:**

- Every batch-selected (and replace) slot retains its local preview without waiting for compress/upload
- Compress pool and upload pool both use concurrency **1**
- Delete/cancel behavior unchanged

**Non-Goals:**

- Changing nginx timeouts in-repo (document in deploy note only if cheap)
- OSS direct upload
- Raising concurrency again
- Client 504 retry (can be follow-up)

## Decisions

1. **Stable list accumulator for patches**  
   - Keep a local mutable list (or functional update) inside `ImageListEditor`: each `patchSlot` applies against the latest local copy, then emit once per patch from that copy.  
   - Sync local copy from `props.modelValue` when parent replaces the array (load/reorder/external).  
   - Alternative: seed all `localPreview` URLs in one `setList` before enqueue — also valid; combine with accumulator for async compress/upload patches.  
   - Prefer **accumulator + batch seed**: on batch, create slots with `localPreview` already set in one emit, then pipeline patches status/pending/URLs only.

2. **Concurrency 1**  
   - `DEFAULT_CONCURRENCY = 1` for both compress and upload shared queues in `imageUploadQueue.js`.  
   - Serial compress then serial upload reduces gateway/OSS pile-up; total wall time may increase; preview must not depend on queue progress (decision 1).

3. **Revoke safety**  
   - When replacing `localPreview`, revoke only the previous URL for that slot; accumulator must not drop other slots’ URLs.

## Risks / Trade-offs

- [Serial uploads slower for large batches] → Accepted for 504 stability; previews stay instant  
- [Single long upload still 504 if > proxy timeout] → Ops may raise `proxy_read_timeout`; optional follow-up retry  
- [Mirror list drift vs parent] → Re-sync on prop identity/length/slotId changes  

## Migration Plan

1. Ship admin build only  
2. Sync `resource/public/admin`  
3. Hard-refresh admin SPA  

## Open Questions

- None — preview fix + concurrency 1 locked.
