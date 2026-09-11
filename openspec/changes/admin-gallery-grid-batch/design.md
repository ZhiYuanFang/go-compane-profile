## Context

Portfolio edit uses `ImageListEditor` + `DualImageField`: single-column full-width rows, batch select runs `processDualImage` sequentially then appends all at once, and `previewUrl` ignores `pending` blobs—so batch rows show blank. Delete is per-row only. Product locked: two-column **thumbnail-first** grid with **click-to-edit**, plus **勾选 / 全选 / 删除所选**.

## Goals / Non-Goals

**Goals:**

- Dense two-column gallery grid with thumbnail as primary surface
- Edit one slot via overlay/dialog containing DualImageField (replace/clear/full preview)
- Batch import shows pending previews and progressive progress
- Multi-select delete with select-all and confirmation for bulk remove
- Preserve existing dual-image model and submit-time upload

**Non-Goals:**

- Changing cover upload UX (except shared DualImageField pending-preview fix if reused)
- Reworking PairImageEditor render/real pairing
- Server-side batch APIs or parallel OSS upload on save (can be a follow-up)
- Drag-and-drop reorder in the grid (keep 上移/下移 in edit dialog or compact ops if needed)
- Mini-program changes

## Decisions

1. **Thumbnail grid in ImageListEditor**  
   - CSS `grid-template-columns: 1fr 1fr` (stack to 1 col on narrow widths if needed).  
   - Cell: checkbox, `#n`, thumbnail (`thumb` URL or object URL from `pending.thumb`), empty placeholder if no image.  
   - Clicking the thumbnail (not the checkbox) opens edit UI.  
   - Alternative: keep full DualImageField inline in two columns — rejected (too tall; product wants thumbnail-first).

2. **Click-to-edit = modal/dialog**  
   - Simple overlay/dialog with DualImageField for the selected index; optional 上移/下移/关闭.  
   - Closing dialog keeps list model in sync via existing v-model updates.  
   - Alternative: side drawer — modal is enough for CMS density.

3. **Pending preview fix in DualImageField**  
   - `previewUrl` falls back to `URL.createObjectURL(pending.thumb)` (revoke on clear/unmount/pending change).  
   - Fixes batch blank previews and benefits cover/single pick consistently.

4. **Progressive batch import**  
   - For each file: process → append one item immediately → update `已处理 i/n`.  
   - Keep sequential process for simplicity and memory; progressive UI fixes “frozen then blank.”  
   - Optional limited concurrency (2) only if sequential still feels too slow after progressive + preview—default sequential in v1.

5. **Selection + batch delete**  
   - Selection state is UI-only (Set of indices), cleared after delete or when list length changes in ways that invalidate indices (rebuild selection carefully after splice).  
   - Toolbar: 全选 / 取消全选 (or toggle), 删除所选 (disabled when none selected).  
   - Confirm when deleting ≥1 selected:「确定删除所选 N 张？」  
   - After delete, selection cleared.

6. **Chinese admin copy**  
   - Keep Simplified Chinese for new buttons/labels (全选、取消全选、删除所选、编辑、处理中 i/n).

## Risks / Trade-offs

- [Object URL leaks] → Centralize revoke on pending change / unmount / clear  
- [Index selection after reorder/delete] → Clear selection on batch delete; after single ops, drop selected indices that no longer match  
- [Modal + file input focus] → Standard dialog pattern; disable batch while processing  
- [Save still slow for many uploads] → Accepted for this change; document as follow-up

## Migration Plan

1. Deploy admin build only (no DB/API migration)  
2. Rollback: previous admin assets

## Open Questions

- None; thumbnail-first and select-all/delete-selected locked by product.
