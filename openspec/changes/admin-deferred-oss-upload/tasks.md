## 1. Slot model and queue split

- [x] 1.1 Add ready-with-pending status (e.g. `SLOT_STATUS.READY`) and update `dualStatusLabel` / `dualStatusTone` (e.g. “待保存”, wait tone—not “已上传”)
- [x] 1.2 Update `dualIsBlockingSave` / `dualsBlockingReason` so ready pending does not block save; still block on compressing, uploading, and error
- [x] 1.3 Change `enqueueFileUpload` / `startSlotPipeline` to stop after compress + `pending` (no `enqueuePendingUpload` on select)
- [x] 1.4 Add shared `flushPendingDualUploads` (or equivalent) that uploads all pending duals via existing stream queue, patches progress, resolves when all DONE or any ERROR

## 2. Components

- [x] 2.1 Wire `DualImageField` to deferred pipeline; remove unused `eagerUpload` (always deferred); retry only re-uploads pending or re-picks file
- [x] 2.2 Wire `ImageListEditor` select/replace/batch import to compress-only; remove immediate upload on select
- [x] 2.3 Confirm `PairImageEditor` inherits deferred behavior via `DualImageField` (no extra eager path)

## 3. Form save paths

- [x] 3.1 `PortfolioEditView`: before create/update/gallery write, flush cover + images pendings; only then build URL payloads
- [x] 3.2 `CompanyView`: flush logo / logoHor pendings before update
- [x] 3.3 `ActivityEditView`: flush image pending before create/update
- [x] 3.4 Saving UX: disable double-submit; button/status reflects upload+save; on upload failure do not persist incomplete image payloads for that submit

## 4. Verification

- [ ] 4.1 Manual: select image → Network has no `/upload/stream` until save; remove before save → still no OSS upload/delete
- [ ] 4.2 Manual: save with pending images → per-slot progress then successful persist; replace already-saved image → after save old OSS cleaned via existing server diff
- [ ] 4.3 Smoke company + activity + portfolio cover/gallery (and pricing if dual-image UI exists)
