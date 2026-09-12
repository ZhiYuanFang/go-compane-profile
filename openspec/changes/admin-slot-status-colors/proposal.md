## Why

Gallery and cover status badges mostly share the same accent color, so queued / compressing / uploading / done are hard to scan at a glance. Operators need clearer visual tiers for status text and cell borders.

## What Changes

- Introduce **four status tone tiers**: waiting (排队中), busy (处理中 / 上传中), success (已上传), failure (失败)
- Apply distinct **text colors** for status copy on gallery badges and cover `DualImageField`
- Apply matching **border/outline** on gallery cells (and cover dual container where practical) by the same tone
- Shared tone helper so both UIs stay consistent
- Rebuild/sync admin assets

## Capabilities

### New Capabilities

- `admin-slot-status-tones`: Four-tier status colors for label text and borders on gallery slots and cover dual fields

### Modified Capabilities

- (none in main `openspec/specs/`)

## Impact

- `admin/src/utils/dualSlot.js` — tone helper (e.g. `dualStatusTone`)
- `admin/src/components/ImageListEditor.vue` — badge + cell border classes
- `admin/src/components/DualImageField.vue` — status + outer border classes
- Build + sync `resource/public/admin`
