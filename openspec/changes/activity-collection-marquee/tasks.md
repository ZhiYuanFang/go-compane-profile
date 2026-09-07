## 1. Schema and API cutover

- [x] 1.1 Replace `activity` schema with collection fields (`title`, `body_html`, image dual, `sort_order`); update seed; document clear/rebuild (no singleton migration)
- [x] 1.2 Implement activity entity/service: list (ordered), get, create (image required), update, delete, reorder; sanitize body HTML subset
- [x] 1.3 Wire public `GET /api/v1/activities` and `GET /api/v1/activities/{id}`; remove singleton public/admin activity endpoints
- [x] 1.4 Wire admin CRUD + reorder API types and controllers

## 2. Admin CMS

- [x] 2.1 Add activity list view with up/down reorder, edit, delete, create entry (Simplified Chinese copy)
- [x] 2.2 Add activity edit/create form: title, required DualImage, rich-text body (font size / bold / color only)
- [x] 2.3 Update admin nav/router/client helpers; remove old singleton ActivityView cover/body-only flow

## 3. Mini-program browse

- [x] 3.1 Data layer: `fetchActivities`, `fetchActivityById`; drop singleton fetch
- [x] 3.2 Home: vertical marquee (1s auto + manual swipe), one-line ellipsis, hide when empty, tap → detail
- [x] 3.3 Activity detail: image full-width height-auto + preview zoom; bold title; rich-text body; badge 「查看全部活动」→ list
- [x] 3.4 Activity list page: two-column waterfall cards (image + title ≤2 lines); register in `app.json`

## 4. Verification

- [x] 4.1 Smoke API/admin: CRUD, reorder, image required, HTML styles round-trip
- [x] 4.2 Smoke mini-program: marquee, detail layout/preview/badge, list waterfall, empty home hides bar
