## Why

Admin portfolio gallery editing is painful at scale: one full-width image per row wastes space; batch select compresses everything before showing rows and never previews `pending` blobs (blank “待上传”); deleting many images requires one-by-one 移除. Operators need a denser grid, visible progressive batch import, and multi-select delete.

## What Changes

- Gallery list becomes a **two-column thumbnail grid**; primary surface is a small preview (checkbox + index); **click thumbnail opens an edit dialog/panel** with the existing DualImageField controls (replace/clear/preview)
- Fix batch import: show **local preview from `pending`**, append rows **progressively** with progress feedback (e.g. `3/20`), keep compression bounded so the UI stays responsive
- Add **勾选 + 全选/取消全选 + 删除所选** (confirm when deleting multiple)
- Keep save-time upload model (images still upload on form submit); no API/schema **BREAKING** changes
- Cover field and PairImageEditor (render/real pairs) are out of scope unless they share DualImageField preview fix

## Capabilities

### New Capabilities

- `admin-gallery-grid`: Two-column thumbnail grid with click-to-edit for portfolio gallery slots
- `admin-gallery-batch-import`: Progressive batch select with pending previews and progress
- `admin-gallery-batch-delete`: Checkbox selection, select-all, delete selected with confirmation

### Modified Capabilities

- (none — no main `openspec/specs/` deltas; additive admin UX)

## Impact

- `admin/src/components/ImageListEditor.vue` (primary)
- `admin/src/components/DualImageField.vue` (pending → blob preview)
- Possibly small modal/drawer component or inline overlay for edit
- `PortfolioEditView.vue` only if wiring/labels change
- Rebuild/sync `resource/public/admin`
- No backend/API changes required for this change
