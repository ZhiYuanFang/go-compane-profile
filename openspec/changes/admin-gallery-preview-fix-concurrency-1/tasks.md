## 1. Preview fix

- [x] 1.1 Add stable list accumulator (or equivalent) in `ImageListEditor` so `patchSlot` never applies against stale props
- [x] 1.2 Batch select: seed all slots with `localPreview` in one list update before starting pipelines; ensure sibling patches cannot drop previews

## 2. Concurrency

- [x] 2.1 Set shared compress and upload `DEFAULT_CONCURRENCY` to 1 in `imageUploadQueue.js`

## 3. Ship

- [x] 3.1 Rebuild admin, sync `resource/public/admin`; smoke batch preview for many files + serial upload (one at a time)
