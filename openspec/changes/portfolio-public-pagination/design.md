## Context

公开 `GET /api/v1/portfolios` 与小程序 `pages/home` 当前一次拉取并渲染全部作品。后台 `ListPortfoliosAdmin` 仍全量，且依赖整表 reorder——本次保持不变。

小程序首页主内容在 **`scroll-view`（`scroll-y`）** 内，不是页面级滚动；因此「触底 / 下拉刷新」必须挂在 `scroll-view` 上，而不是仅依赖 `Page.onReachBottom` / `Page.onPullDownRefresh`（页面本身往往不滚动，回调不会触发）。

涉及两个仓库：`go-compane-profile`（API）与 `wx-compane-profile-link`（小程序）。

## Goals / Non-Goals

**Goals:**

- 公开列表分页：默认每页 10 条，返回 `list` + `total` + `page` + `pageSize`
- 小程序首页：首屏一页、触底追加、底部「加载中 / 没有更多」、下拉刷新回到第一页
- 排序规则不变：`sort_order ASC, id ASC`
- 与小程序同步发版，避免旧客户端只看到第一页

**Non-Goals:**

- 后台作品列表分页
- 改 reorder / sort_order 交互
- 公开详情、公司、资费接口变更
- 真瀑布流按高度分列（继续 `index % 2`）

## Decisions

### 1. 仅公开 API 分页；后台保持全量

- **选择**：`ListPortfoliosPublic` 增加 page/pageSize；admin 不动  
- **理由**：运营侧依赖整表排序；作品量对后台可接受  
- **备选**：两边统一分页 → 引入跨页排序复杂度，本期不需要  

### 2. 查询参数与默认值

| 参数 | 默认 | 说明 |
|------|------|------|
| `page` | `1` | `< 1` 时按 `1` |
| `pageSize` | `10` | 上限 `50`（防刷）；非法/超限钳制到合法范围 |

响应：

```json
{
  "list": [ /* PortfolioListItem */ ],
  "total": 86,
  "page": 1,
  "pageSize": 10
}
```

列表项字段不变（`id`=slug、`cover` thumb 等）。空库：`list=[]`，`total=0`。

- **BREAKING**：无参不再返回全量，只返回第 1 页  
- **备选**：无参全量兼容旧客户端 → 拖长双行为；选定同步发版后直接切  

### 3. Service 层实现

- `Count` + `Page`/`Limit`（GoFrame Model）按同一 Order 查询  
- Controller 从 query 解析 page/pageSize，填入 Res  

### 4. 小程序：scroll-view 触底 + refresher（对齐产品意图）

产品要求「触底文案 + onPullDownRefresh」。在现有 `scroll-view` 布局下：

| 产品意图 | 实现 |
|----------|------|
| 触底加载 | `bindscrolltolower`（可设 `lower-threshold`） |
| 下拉刷新 | `refresher-enabled` + `bindrefresherrefresh`（必要时 `refresher-triggered`） |
| 文案 | 列表下方：加载中 → 没有更多了；首屏空列表可沿用现有失败 toast |

状态机：

```
idle ──触底且 hasMore──▶ loadingMore ──成功──▶ idle（追加）
                         └──失败──▶ idle（toast，可重试）
idle ──下拉──▶ refreshing ──成功──▶ idle（重置为第 1 页）
```

- 并发：loading/refreshing 期间忽略重复触底  
- 分列：对**累计** `portfolioList` 再 `_splitColumns`（全局 `index % 2`），勿按页单独分列后拼接  
- `fetchPortfolioList({ page, pageSize })` 返回 `{ list, total, page, pageSize }`  
- `hasMore = accumulated.length < total`  

**备选**：改回页面滚动 + `onPullDownRefresh` / `onReachBottom` → 要动导航/布局；本期不改结构。

### 5. 仓库与发版

1. 先发后端（或同窗口）  
2. 再发小程序；旧小程序在后端已切分页后只会显示 10 条——需控制发布窗口  

## Risks / Trade-offs

- **[Risk] 旧小程序未更新 → 用户只见 10 条** → 同步发版；必要时短暂 feature 开关（本期不做开关，靠发布节奏）  
- **[Risk] scroll-view refresher 机型差异** → 用官方 refresher API；真机抽测 iOS/Android  
- **[Risk] 触底连触发重复请求** → `loading` 锁 + `hasMore` 判断  
- **[Trade-off] 无参打破全量契约** → 换取实现简单；依赖同发版  

## Migration Plan

1. 合并并部署 `go-compane-profile`（公开列表分页）  
2. 发布 `wx-compane-profile-link` 首页分页版  
3. 回滚：后端可临时恢复全量（若需）；小程序回退上一版  

## Open Questions

- 无（pageSize=10、后台不分页、触底文案、下拉刷新已确认；下拉在 scroll-view 上用 refresher 实现）
