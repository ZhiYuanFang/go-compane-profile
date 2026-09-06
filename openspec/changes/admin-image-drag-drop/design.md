## Context

所有 CMS 图片字段共用 `DualImageField`：仅有隐藏 `input[type=file]` +「选择图片」按钮；预览点击已用于 lightbox 放大。运营希望从桌面拖入图片，并在空态单击预览也能选文件。

## Goals / Non-Goals

**Goals:**

- 在 `DualImageField` 上支持拖拽选图，复用现有 `processDualImage` → pending → 提交上传
- 拖到已有图上替换前必须确认
- 空预览单击打开文件选择；有图单击仍放大
- 保留按钮选文件

**Non-Goals:**

- 一次拖多张进成组编辑器自动拆分效果/实景
- 剪贴板粘贴（可后续）
- 改后端或上传 API
- 拖拽时立即上传 OSS

## Decisions

### 1. Drop zone = 整个 `.dual` 卡片

- **选择：** `dragenter` / `dragover` / `dragleave` / `drop` 挂在根容器；拖入时高亮边框
- **备选：** 仅预览区 — 命中面积偏小，否决

### 2. 单文件、仅 image/*

- **选择：** 从 `dataTransfer.files` 取第一张 `type` 以 `image/` 开头的文件；无则设错误提示
- 与现有 `accept="image/*"` 一致

### 3. 覆盖确认

- **选择：** 若当前已有 `previewUrl`（已保存或 pending），drop 后先 `window.confirm`（或轻量自定义确认）；取消则不处理；确认后走与 `onFile` 相同处理
- **空态 drop：** 无需确认，直接处理
- **备选：** 无确认直接替换 — 否决（产品要求确认）

### 4. 空预览单击 = 选文件；有图单击 = 放大

```
click preview:
  if !previewUrl → trigger hidden file input click
  else → openLightbox
```

按钮「选择图片 / 重新选择」行为不变（重新选择不额外确认，与今日一致；仅**拖拽覆盖**需确认）。

### 5. 实现落点

- 仅改 `DualImageField.vue`（+ 必要时少量样式）
- 抽 `applyFile(File)` 供 input change 与 drop 共用

## Risks / Trade-offs

- [浏览器默认打开拖入的文件] → `preventDefault` on dragover/drop
- [dragleave 闪烁] → 用计数或检测 relatedTarget 是否仍在容器内
- [confirm 简陋] → 可接受；若日后统一 modal 再换
- [有图时误点想换图] → 仍用「重新选择」按钮或拖入（拖入有确认）

## Migration Plan

1. 改前端并 `npm run build`，同步 `resource/public/admin`
2. 无 API/DB 变更；回滚即回退 admin 静态资源

## Open Questions

- 无阻塞项。「重新选择」按钮是否也加确认：默认**不加**（与现网一致，仅拖拽覆盖确认）。
