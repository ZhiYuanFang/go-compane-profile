## Context

CMS already exposes public GETs for company, portfolio detail, and activity detail; admin lists show cover/title/reorder actions. Operators want rough engagement signals. Explore locked: three dedicated increment APIs (not GET side-effects), mini-program `onLoad` fire-and-forget, admin display「热度 N」on portfolio/activity lists and「查阅关于我们次数为 x」on company page.

## Goals / Non-Goals

**Goals:**

- Persist and atomically increment view counters for portfolio, activity, and about-us
- Public POST APIs for mini-program to call on page `onLoad`
- Admin-only visibility of counts with agreed Simplified Chinese copy
- Migration with `DEFAULT 0` so existing rows start at zero

**Non-Goals:**

- Showing counts in the mini-program UI
- Deduplication / UV / daily unique (every `onLoad` counts)
- Rate limiting or anti-abuse beyond “best-effort public counters”
- Manual reset/edit of counts in admin
- Embedding increments into existing GET detail handlers
- Analytics dashboards or time-series

## Decisions

1. **Storage columns, not a separate stats table**  
   - `portfolio.view_count INT NOT NULL DEFAULT 0`  
   - `activity.view_count INT NOT NULL DEFAULT 0`  
   - `company.about_view_count INT NOT NULL DEFAULT 0`  
   - Rationale: company is already a singleton; list queries stay simple.  
   - Alternative: event log table — rejected as overkill for this scope.

2. **Three public POST endpoints**  
   - `POST /api/v1/portfolios/{id}/view`  
   - `POST /api/v1/activities/{id}/view`  
   - `POST /api/v1/company/about/view`  
   - Atomic `UPDATE ... SET view_count = view_count + 1` (or about column).  
   - Missing portfolio/activity → 404; about always updates the singleton company row.  
   - Response: empty/success body is enough (no need to return new count).  
   - Alternative: query `?track=1` on GET — rejected per product (separate APIs).

3. **Increment on mini-program `onLoad` only**  
   - Detail/about pages call once per cold navigation into the page; returning via `onShow` without reload does not double-count for the same page instance.  
   - Errors swallowed client-side so tracking never blocks content.

4. **Admin display**  
   - Portfolio/activity list rows: `热度 {n}` (e.g. `热度 12`).  
   - Company page: read-only module「查阅关于我们次数为 {n}」using count from existing company GET (include `aboutViewCount` in admin company payload).  
   - Public list/detail GETs need not expose counts to mini-program; admin list/company responses MUST include them.

5. **SQL migration**  
   - Add ALTER migration file under `manifest/sql/`; update `schema.sql` if present/restored in repo.  
   - No backfill beyond DEFAULT 0.

## Risks / Trade-offs

- [Public POST can be scraped/inflated] → Accept for internal heat signal; no rate limit in this change  
- [Concurrent increments] → Use SQL `column = column + 1`, not read-modify-write in app  
- [Admin preview of mini pages] → N/A (admin is web CMS, not mini-program)  
- [schema.sql drift] → Ship explicit migration operators can run

## Migration Plan

1. Run ALTER to add the three columns with DEFAULT 0  
2. Deploy API with increment endpoints + admin field exposure  
3. Deploy admin UI  
4. Deploy mini-program with `onLoad` calls  
5. Rollback: stop calling APIs / redeploy previous builds; columns can remain (harmless)

## Open Questions

- None for implementation; copy and `onLoad` locked by product.
