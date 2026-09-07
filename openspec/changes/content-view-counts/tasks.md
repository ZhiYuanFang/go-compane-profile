## 1. Schema and backend

- [x] 1.1 Add migration + entity fields: `portfolio.view_count`, `activity.view_count`, `company.about_view_count` (DEFAULT 0); update schema if present
- [x] 1.2 Implement atomic increment services for portfolio, activity, and company about
- [x] 1.3 Wire public POST `/api/v1/portfolios/{id}/view`, `/api/v1/activities/{id}/view`, `/api/v1/company/about/view`
- [x] 1.4 Include `viewCount` / `aboutViewCount` in admin list and company payloads (public mini GETs need not expose counts)

## 2. Admin CMS

- [x] 2.1 Portfolio list: show「热度 {n}」per row
- [x] 2.2 Activity list: show「热度 {n}」per row
- [x] 2.3 Company page: read-only module「查阅关于我们次数为 {n}」; rebuild/sync admin assets

## 3. Mini-program

- [x] 3.1 Add thin helpers to POST the three view endpoints
- [x] 3.2 Portfolio detail `onLoad`: fire-and-forget increment (errors ignored)
- [x] 3.3 Activity detail `onLoad`: fire-and-forget increment
- [x] 3.4 About page `onLoad`: fire-and-forget increment

## 4. Verification

- [x] 4.1 Smoke: POST increments update DB; missing id → 404; GET detail does not increment
- [x] 4.2 Smoke: admin lists/company show counts; mini `onLoad` triggers POSTs without blocking UI
