## 1. Public API (go-compane-profile)

- [x] 1.1 Extend `api/v1` `ListPortfoliosReq` with `page` / `pageSize` query fields and `ListPortfoliosRes` with `total`, `page`, `pageSize` (keep `list`)
- [x] 1.2 Change `ListPortfoliosPublic` to count + paginated select (order `sort_order ASC, id ASC`); clamp `page` ≥ 1 and `pageSize` default 10 max 50
- [x] 1.3 Wire public controller to pass pagination params and return metadata
- [x] 1.4 Smoke-check: default first page 10 items; `page=2`; empty DB; oversized `pageSize` clamped

## 2. Miniprogram data layer (wx-compane-profile-link)

- [x] 2.1 Update `fetchPortfolioList` to accept `{ page, pageSize }` and return `{ list, total, page, pageSize }` from `/api/v1/portfolios`

## 3. Miniprogram home feed

- [x] 3.1 Track accumulated list, `page`, `total`, `loading` / `refreshing` / `hasMore` on `pages/home`
- [x] 3.2 Initial load: page 1 size 10; `_splitColumns` on full accumulated list
- [x] 3.3 `scroll-view` `bindscrolltolower`: load next page when `hasMore` and not in flight; show「加载中」; when exhausted show「没有更多」
- [x] 3.4 Enable `scroll-view` refresher: on refresh reload page 1, replace list, reset pagination, end refresher
- [x] 3.5 Ignore duplicate load-more while loading/refreshing; toast on failure

## 4. Verify

- [x] 4.1 Manual: first paint ≤10 cards; scroll loads more; footer copy; pull refresh resets
- [x] 4.2 Confirm admin portfolio list/reorder still full-list (no regression)
