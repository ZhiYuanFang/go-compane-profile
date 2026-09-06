## Context

管理后台作品编辑目前用两个独立 `ImageListEditor`（`renders` / `reals`），Slug 手填，`DualImageField` 预览为 `object-fit: cover` 且不可放大。创建成功后 `replace` 到编辑页。后端创建时空 slug 已自动生成 `p-{id}`；图库仍为两列有序数组。

兄弟仓小程序详情已用纵向 swiper + `movable-view` scale，但与 swiper 抢手势导致缩放实际不可用。详情组装时 `reals[index] || reals[0]` 在空实景时会错位。

## Goals / Non-Goals

**Goals:**

- 后台图片预览完整可见（contain）并可点击放大核对
- 作品编辑按「效果↔实景」成组录入，实景可空且下标对齐
- 隐藏 Slug；保存成功回列表
- 小程序详情非封面用 `wx.previewImage` 提供可靠全屏缩放

**Non-Goals:**

- 不改公开 API 契约（仍返回 `renders` / `reals` 数组）
- 不改为数字 id 对外暴露；slug 仍自动生成
- 不加大后台预览框尺寸
- 作品列表封面不增加点击放大
- 不重做小程序详情竖滑分页结构
- 不引入新后端依赖或 OSS 流程变更

## Decisions

### 1. 预览 contain + 框尺寸不变

- **选择：** `DualImageField` 的 `.preview img` 改为 `object-fit: contain`；宽高保持现有约 `140×100`
- **备选：** 加大预览框 — 否决（产品明确不要）
- **备选：** `cover` 保留 + 仅放大查看完整图 — 否决（框内也要完整展示）

### 2. 后台放大用组件内 lightbox，非 `window.open`

- **选择：** 点击预览打开全屏/半屏遮罩，优先显示 `original`（已保存）或 pending 本地预览 URL；遮罩/Esc/关闭退出
- **备选：** 新标签打开原图 — UX 较差，否决
- **范围：** 所有使用 `DualImageField` 的字段；列表封面除外

### 3. Slug：仅藏 UI

- **选择：** 新建/编辑表单不展示 Slug；提交 `slug: ''`（新建）或沿用已有（编辑不传改 slug）
- **备选：** 公开 id 改数字 — 否决（影响已有链接与小程序）

### 4. 成组 UI，API 仍拆两数组

- **选择：** 新组件（如 `PairImageEditor`）管理 `pairs: [{ render, real }]`；提交时：
  - `renders = pairs.map(p => resolve(p.render))`
  - `reals = pairs.map(p => resolveOrEmpty(p.real))` 保证 **等长**，空实景为 `{ thumb: '', original: '' }`
- **加载：** `zip(renders, reals)`，缺侧补空 dual
- **备选：** DB 增加 pair_id — 过度，否决
- **效果图：** 提交时该组须有效果图；实景可选

### 5. 空 URL 图库行

- **选择：** 保存时保留空 dual 占位行，保证与效果图下标对齐；若现有 `replaceGallery` 过滤空 URL，则调整为允许空行或写入占位（实现时核对）
- **小程序：** 空 URL 视为无实景，禁止回退到 `reals[0]`

### 6. 保存后导航

- **选择：** 新建与更新成功均 `router.push`/`replace` 到 `/portfolios`
- **备选：** 仅新建回列表 — 否决（产品要求编辑也回）

### 7. 小程序原生预览（方案 B）

- **选择：** 非封面主图 `bindtap` → `wx.previewImage`；`urls` 为本作品所有非封面可展示原图（效果 + 非空实景，按浏览顺序），`current` 为点击项；封面不入 `urls`
- **圆卡切换：** 保持 `catchtap`，不触发预览
- **清理：** 可移除 `movable-view` 的 scale / `zoomLocked` / `_scaleMap`，主图改为静态 `image`，降低与 swiper 冲突
- **备选：** 继续修 movable scale — 已验证不可靠，否决
- **实现仓：** `wx-compane-profile-link`（跨仓任务）

## Risks / Trade-offs

- [空 URL 被后端丢弃] → 实现前读 `replaceGalleryTx`；若过滤则改为保留占位或同步改小程序仅按 renders 分页并独立挂实景
- [小程序 `previewImage` 域名未配置] → 确认 downloadFile/业务域名已含 CDN；与现网图片加载同一域名
- [成组编辑后旧数据不等长] → zip 时按 `max(len(renders), len(reals))` 补空，避免丢图
- [跨仓任务遗漏] → tasks 明确标注兄弟仓路径；本仓 apply 可先做 admin，小程序单独开 PR

## Migration Plan

1. 合并并部署本仓 admin 静态资源（及如有的后端占位修复）
2. 运营无需数据迁移；编辑旧作品时自动 zip 成组
3. 兄弟仓发布小程序详情预览改动
4. 回滚：回退 admin 构建与小程序版本；DB 无破坏性变更

## Open Questions

- 无阻塞项。`previewImage` 的 `urls` 采用「整组非封面可左右滑」（已定为推荐默认）。
