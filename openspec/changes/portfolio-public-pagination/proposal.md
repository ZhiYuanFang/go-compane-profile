## Why

作品集公开列表当前一次返回全部条目。作品变多后，小程序首页首屏请求与渲染成本过高。需要公开侧分页；后台作品量可控且依赖整表排序，本次不分页。

## What Changes

- 公开 `GET /api/v1/portfolios` 支持 `page` / `pageSize` 查询参数，响应增加 `total`、`page`、`pageSize`（列表仍为 `list`）
- 默认 `page=1`、`pageSize=10`；对 `pageSize` 设合理上限
- **BREAKING**（对仍假设一次拉全量的公开客户端）：不传参时改为返回第一页而非全量；与小程序同步发版
- 小程序首页：触底加载更多、底部「加载中 / 没有更多」文案、`onPullDownRefresh` 重置第一页
- 后台 `/admin/api/portfolios` 与排序/reorder **不变**（仍全量）

## Capabilities

### New Capabilities

- `public-portfolio-pagination`: 公开作品列表分页契约（查询参数、响应字段、默认与上限）
- `miniprogram-home-portfolio-feed`: 小程序首页分页加载、触底与下拉刷新的产品行为

### Modified Capabilities

- （无现有主规格；本 change 以新增能力规格描述行为）

## Impact

- **Backend** (`go-compane-profile`): `api/v1/public.go`、`internal/controller/public`、`internal/service.ListPortfoliosPublic`
- **Miniprogram** (`wx-compane-profile-link`): `data/portfolio.js`、`pages/home/*`（含 `enablePullDownRefresh`）
- **Admin**: 无变更
- **Deploy**: 后端与小程序需同窗口发布，避免旧小程序只拿到第一页
