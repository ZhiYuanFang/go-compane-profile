## Context

当前活动为单例表行：`cover_*` + `body_*` 双图；公开 `GET /api/v1/activity`；首页 80rpx 封面条；详情仅滚正文图。产品要改为多活动集合、HTML 正文、首页纵向标题轮播与全部活动列表。旧数据允许清空重录。公司简介与心流不做富文本。

## Goals / Non-Goals

**Goals:**

- 活动多行 CRUD + `sort_order`（列表/轮播同一顺序）
- 字段：标题、必填 dual 图、正文 HTML（字号/加粗/颜色）；图文分离
- 公开全量列表 + 详情；不分页
- 小程序：首页纵向跑马灯、详情（图→标题→文）、列表瀑布流；简体文案

**Non-Goals:**

- 公司 `introBody` / 作品 `heartFlow` / 奖项富文本
- 正文内插图
- 活动分页、排期、上下架开关（有标题+图即可出现在列表；删除即下线）
- 繁体文案体系
- 旧单例数据自动迁移

## Decisions

### 1. Replace singleton schema

- 新 `activity` 行模型：`id`, `title`, `body_html`, `image_original_url`, `image_thumb_url`, `sort_order`, timestamps
- 去掉 cover/body 双图列语义；部署时 DROP/重建或清空后 ALTER（运营重录）
- **BREAKING**：删除公开/后台单例 activity 读写，改为 `/activities` 与 `/activities/{id}`

### 2. Public IDs

- 列表/详情公开使用数字 `id`（字符串化亦可）；admin 同作品集用 id
- **Alternatives:** slug — 可选但非必须；默认数字 id 更简单

### 3. Body HTML

- 存 sanitized HTML 子集：`p`, `br`, `span`, `strong`, `b`；`style` 仅 `color` / `font-size` / `font-weight`
- Admin：轻量编辑器（如 Quill/TipTap/wangEditor 之一）工具栏限字号/粗/色
- 小程序：`<rich-text>` 展示
- 图不进 HTML；独立 DualImage，创建/更新校验非空

### 4. Ordering

- Admin 列表上下移 → `reorder` 写回 `0..n-1`
- 公开 list `ORDER BY sort_order ASC, id ASC`
- 首页轮播、列表页同一顺序

### 5. Mini-program UX

- 首页：有列表才渲染条；纵向切换（swiper vertical 或等效）；`interval`/`autoplay` 约 1000ms；支持手势；标题单行 ellipsis；点击进 `/pages/activity/activity?id=`
- 详情：顶图 `widthFix` 或宽 100% 高度自适应；`previewImage` 放大；标题加粗；下接 rich-text；图右下角 badge「查看全部活动」（对齐首页作品封面地址 badge 视觉）
- 列表页：双列瀑布流卡片，上图下标题（最多两行）；点进详情

### 6. Copy language

- UI 硬编码与运营向提示使用简体（如「查看全部活动」）

## Risks / Trade-offs

- [HTML XSS] → 服务端或保存前白名单消毒
- [编辑器体积] → 选轻量包；仅活动编辑页加载
- [rich-text 样式差异] → 限制工具栏能力，避免复杂标签
- [纵向 swiper 与首页手势] → 控制条高度、避免抢作品列表横滑（首页活动为竖向区域）
- [旧客户端仍调 /activity] → 同窗口发版；接口 404 可接受

## Migration Plan

1. 部署新 schema（重建 activity 表）
2. 发 API + admin
3. 发小程序
4. 运营重新录入活动
5. Rollback：回退代码；表结构需保留备份若要恢复单例（默认不保留）

## Open Questions

- （无阻塞项；图必填、排序、跑马灯 1s、不分页、简体、清空重录已确认）
