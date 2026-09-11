## 1. Preview and batch import

- [x] 1.1 DualImageField: preview from `pending` via object URL; revoke on change/unmount/clear
- [x] 1.2 ImageListEditor batch select: progressive append + `i/n` progress; keep partial results on error; lock controls while processing

## 2. Thumbnail grid and click-to-edit

- [x] 2.1 Replace full-width row list with two-column thumbnail grid (checkbox, index, thumb/placeholder)
- [x] 2.2 Click thumbnail opens dialog with DualImageField (+ 上移/下移/移除 as needed); checkbox does not open editor
- [x] 2.3 Keep 添加一张 / empty state; narrow viewport may stack to one column

## 3. Batch delete

- [x] 3.1 Per-cell selection + toolbar 全选 / 取消全选
- [x] 3.2 删除所选 with confirm count; clear selection after delete; disabled when none selected

## 4. Ship

- [x] 4.1 Rebuild admin and sync `resource/public/admin`
- [x] 4.2 Smoke: two-column thumbs, click-to-edit, batch preview/progress, select-all + delete selected
