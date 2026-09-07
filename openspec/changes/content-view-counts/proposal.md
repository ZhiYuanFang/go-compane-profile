## Why

Operators need a simple sense of which portfolios, activities, and the About Us page are being opened in the mini-program. There is no view tracking today; adding lightweight counters (incremented by dedicated public APIs) gives CMS visibility without coupling counts to content GET responses.

## What Changes

- Add `view_count` on `portfolio` and `activity`, and `about_view_count` on `company` (default 0)
- Add three public increment APIs the mini-program calls once per page `onLoad`
- Surface counts in admin: portfolio/activity lists as「热度 N」; company page module「查阅关于我们次数为 x」
- Mini-program: fire-and-forget increment on portfolio detail, activity detail, and about page `onLoad`
- Counts are admin-only display; mini-program UI does not show them

## Capabilities

### New Capabilities

- `view-count-api`: Public POST endpoints to atomically increment portfolio, activity, and about-us view counts
- `view-count-admin`: Admin list/company UI and list payloads expose view counts with the agreed Chinese copy
- `view-count-miniprogram`: Mini-program calls the three increment APIs on relevant page `onLoad`

### Modified Capabilities

- (none — no main-spec requirement deltas; behavior is additive)

## Impact

- DB migration + entity/service/list DTO fields
- Public API (`api/v1`) + public controller
- Admin portfolio/activity list views and company view; admin list API responses
- Mini-program `pages/detail`, `pages/activity`, `pages/about` (+ thin data helpers)
- No auth on increment endpoints (public, best-effort counters; no rate limit in this change)
