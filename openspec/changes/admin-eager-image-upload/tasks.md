## 1. Size limits

- [x] 1.1 Client: original max 3MiB; thumb max 0.5MiB (aligned with compress)
- [x] 1.2 Server: `MaxOriginalBytes` 3MiB; `MaxThumbBytes` 0.5MiB; update error copy

## 2. Upload queue

- [x] 2.1 Add slotId + status model and concurrency-3 queue (enqueue, abort/cancel by slotId, ignore stale completions)
- [x] 2.2 Wire ImageListEditor: add/replace/batch → enqueue; badges 处理中/排队中/上传中/已上传/失败+重试
- [x] 2.3 Delete single/batch cancels queue or aborts in-flight; safe no-op on late success
- [x] 2.4 DualImageField cover: same eager upload + status

## 3. Save path

- [x] 3.1 Portfolio save: block while uploading/queued/failed; submit URLs without re-upload when done
- [x] 3.2 Rebuild admin, sync `resource/public/admin`; smoke batch queue, delete mid-upload, save gate, 3MiB limit
