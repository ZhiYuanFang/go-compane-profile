## Why

运营在后台录入作品时，Slug 多余、效果图与实景图分开添加易错位、选图后预览被裁切且无法核对原图，保存后还停在编辑页打断流程。小程序详情页页内 pinch 缩放因 swiper 与 movable-view 抢手势实际失效，需要改为原生预览。

## What Changes

- 管理后台图片预览改为框内完整展示（`object-fit: contain`），预览框尺寸不加大
- 所有 `DualImageField` 预览支持点击遮罩放大查看（优先原图）；作品列表封面不改
- 新建/编辑作品表单隐藏 Slug 输入；创建时传空，沿用后端自动 `p-{id}`
- 效果图与实景图改为成组编辑：添加一组时可同时选效果图与对应实景图；实景可空；保存时按同下标拆成 `renders` / `reals`（空槽占位对齐）
- 新建与编辑保存成功后均返回作品列表
- 兄弟仓小程序详情：非封面主图点击调用 `wx.previewImage`（原生全屏 + 手势缩放）；封面不进预览；修正空实景按下标配对逻辑；可移除失效的页内 scale

## Capabilities

### New Capabilities

- `admin-image-preview`: 管理后台双图字段的 contain 预览与点击放大
- `admin-portfolio-pair-gallery`: 作品编辑成组图库、隐藏 Slug、保存后回列表
- `miniprogram-detail-native-preview`: 详情非封面原生预览与实景空槽配对

### Modified Capabilities

- （无）`openspec/specs/` 尚未归档主规格；本变更以新 capability 规格为准

## Impact

- 本仓：`admin/src/components/DualImageField.vue`、新建配对编辑组件、`PortfolioEditView.vue`；可能触及 `ImageListEditor.vue` 的替换或弃用
- 本仓后端：预期无 API/schema 变更（空 slug 与双数组图库已支持）；实现时确认空 URL 图库行可按占位保存
- 兄弟仓：`wx-compane-profile-link` 的 `pages/detail/*`（本仓 tasks 标注跨仓，实现在兄弟仓执行）
- 构建产物：`admin` 构建后需同步到 `resource/public/admin`（若项目现有流程要求）
