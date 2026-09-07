## 1. Schema and shared enums

- [x] 1.1 Add `portfolio.category` column and document allowed values (`residential|commercial|office|installation`) in `manifest/sql/schema.sql` (plus any deploy migration note)
- [x] 1.2 Add `company.awards`, `company.phone`, `company.wechat` columns in schema/seed as needed
- [x] 1.3 Add `activity` singleton table (cover + body dual URL fields) in schema; seed empty or omit row
- [x] 1.4 Remove `pricing` table from schema/seed (and note drop for existing DBs)

## 2. Backend portfolio category

- [x] 2.1 Extend portfolio entity/DTOs with `category`; validate on create/update
- [x] 2.2 Public `ListPortfolios`: optional `category` filter; invalid → 4xx; order by `sort_order`, `id`
- [x] 2.3 Admin list filtered by category; create defaults category from request; reorder scoped to category only
- [x] 2.4 Implement category change: append to new category order; resequence old category

## 3. Backend activity

- [x] 3.1 Implement activity entity/service: get, upsert images, clear/delete
- [x] 3.2 Public `GET /api/v1/activity` empty-safe response
- [x] 3.3 Admin `GET/PUT` (and delete/clear) activity endpoints wired in controllers/API types

## 4. Backend company about + remove pricing

- [x] 4.1 Extend company entity/service/API with `awards`, `phone`, `wechat` (preserve newlines in `awards`)
- [x] 4.2 Remove pricing service, controllers, API types, and routes (public + admin)

## 5. Admin console UI

- [x] 5.1 Replace single 作品集 nav with four category entries; routes pass fixed category into list/edit
- [x] 5.2 Update portfolio list/edit/reorder client calls to be category-scoped; support changing category on edit
- [x] 5.3 Add Activity admin page (cover + body upload/replace/delete) and sidebar link
- [x] 5.4 Extend CompanyView with awards textarea, phone, wechat fields
- [x] 5.5 Remove PricingView, route, nav link, and client helpers

## 6. Mini-program data layer

- [x] 6.1 Update `fetchPortfolioList` to accept optional `category`
- [x] 6.2 Add `fetchActivity` data helper for `GET /api/v1/activity`
- [x] 6.3 Extend company fetch typing/usage for awards/phone/wechat; remove pricing data module

## 7. Mini-program home category tabs + activity

- [x] 7.1 Refactor home to Tab row (全部 + four categories) with selected state larger/bolder
- [x] 7.2 Add synced swiper/viewpager; per-page list state, category filter, and pagination
- [x] 7.3 Load activity on home; show 80rpx aspectFill banner above tabs when present; hide when empty
- [x] 7.4 Add activity body page (scrollable single image `aspectFit`); wire banner tap navigation
- [x] 7.5 Replace pricing FAB with About Us entry; register about page; remove pricing page from `app.json` and files

## 8. Mini-program about page + multiline copy

- [x] 8.1 Build About page: scrollable awards with newline rendering; fixed footer address/copy, phone/dial, wechat/copy
- [x] 8.2 Fix detail heartFlow to render newlines via split-lines or `<text>` strategy (not view pre-wrap alone)
- [x] 8.3 Apply the same newline strategy to awards on About page

## 9. Verification

- [x] 9.1 Smoke: public list all vs by category; admin reorder per category; change category end-to-end
- [x] 9.2 Smoke: activity set/clear → home banner + body page; empty hides banner
- [x] 9.3 Smoke: company about fields → About page actions; pricing routes/UI gone; heartFlow multiline on detail
