## Why

管理后台选图只能通过「选择图片」打开系统文件对话框，运营从桌面或资源管理器拖图进来更自然。需要在统一的双图字段上支持拖拽，并保留点击选文件通道。

## What Changes

- `DualImageField` 支持将图片拖入字段区域完成选择（走现有本地压缩管道，提交时再上传）
- 拖到**已有预览**上替换时，须经操作员确认后再替换
- **空预览**单击即可打开文件选择（与按钮、拖拽并列）
- 已有预览单击仍打开放大查看（不改为选文件）
- 保留现有「选择图片 / 重新选择」按钮

## Capabilities

### New Capabilities

- `admin-image-drag-drop`: 双图字段拖拽选图、覆盖确认、空态点击选文件

### Modified Capabilities

- （无）`openspec/specs/` 尚未归档主规格

## Impact

- 本仓前端：`admin/src/components/DualImageField.vue`（公司 / 资费 / 作品封面 / 成组图库共用）
- 后端与 OSS 上传契约不变
- 构建后需同步 `resource/public/admin`（若按现有部署流程）
