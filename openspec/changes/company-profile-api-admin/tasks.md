## 1. Project foundation & config

- [x] 1.1 Add `.gitignore` entries for `.env.prod`, admin `node_modules`/`dist`, and local secrets
- [x] 1.2 Create `manifest/docker/env/.env.example` with placeholders (MySQL, OSS, CDN, ADMIN_PASSWORD, REGISTRY, IMAGE_TAG)
- [x] 1.3 Wire GoFrame config to read DB/OSS/admin settings from environment (no secrets in tracked yaml)
- [x] 1.4 Create MySQL database `compane_profile` and document connection via `.env.prod`

## 2. Schema, seed & domain models

- [x] 2.1 Add SQL/migration for `company`, `portfolio`, `portfolio_image`, `pricing` tables
- [x] 2.2 Seed company copy text only (from former mini-program mock); leave portfolios/pricing empty
- [x] 2.3 Add Go entity/DAO/logic stubs for company, portfolio, portfolio images, pricing

## 3. Public portfolio API

- [x] 3.1 Implement `GET /api/v1/company` with dual-URL logo fields
- [x] 3.2 Implement `GET /api/v1/portfolios` ordered by `sort_order`, cover = thumb
- [x] 3.3 Implement `GET /api/v1/portfolios/:id` with cover/renders/reals as `{ thumb, original }`
- [x] 3.4 Implement `GET /api/v1/pricing` with dual URLs
- [x] 3.5 Register routes with standard GoFrame response middleware; verify not-found for missing portfolio

## 4. Admin auth

- [x] 4.1 Implement `POST /admin/api/login` validating `ADMIN_PASSWORD` from env
- [x] 4.2 Implement session/JWT middleware protecting `/admin/api/*` except login
- [x] 4.3 Implement `POST /admin/api/logout` invalidating the session

## 5. OSS dual-image upload

- [x] 5.1 Add Aliyun OSS client using shared bucket/endpoint/CDN and prefix `muhou/`
- [x] 5.2 Implement authenticated `POST /admin/api/upload` accepting original+thumb pair, enforcing 5MB/500KB limits, distinct object keys, returning CDN URLs
- [x] 5.3 Reject keys outside `muhou/` for any delete helpers added later

## 6. Admin CMS APIs

- [x] 6.1 Implement company GET/PUT for authenticated admins
- [x] 6.2 Implement portfolio CRUD + reorder API
- [x] 6.3 Implement portfolio image gallery save (renders/reals ordered dual URLs)
- [x] 6.4 Implement pricing GET/PUT for dual URLs

## 7. Vue 3 admin console

- [x] 7.1 Scaffold `admin/` with Vue 3 + Vite; glassmorphism + high-end gray theme tokens
- [x] 7.2 Build login page and auth-aware API client against `/admin/api`
- [x] 7.3 Build company editor (copy + logo dual-image fields)
- [x] 7.4 Build portfolio list/reorder and portfolio editor (cover + renders/reals)
- [x] 7.5 Build pricing editor with dual-image field
- [x] 7.6 Implement client-side image pipeline: process to ≤5MB + ≤500KB, allow re-select, upload pairs only on submit
- [x] 7.7 Configure Go to serve `admin/dist` at `/admin/`

## 8. Deploy & ACR

- [x] 8.1 Multi-stage Dockerfile: build admin + Go binary; runtime serves API and static admin
- [x] 8.2 Add simplified `.github/workflows/docker-acr.yml` (single image, tag/dispatch, VPC-strip REGISTRY)
- [x] 8.3 Add production `.env.prod` template notes / runbook snippet for `muhou.cuplay.top` and ACR namespace `muhou`

## 9. Mini-program API wire (sibling repo)

- [x] 9.1 Add request helper + baseURL `https://muhou.cuplay.top` in `微信小程序/compane-profile`
- [x] 9.2 Replace welcome/home/detail/pricing data loading with public APIs; remove `data/company.js` and `data/portfolio.js` mocks
- [x] 9.3 Home cards use thumb cover; detail/pricing progressive thumb → original swap
- [x] 9.4 Document WeChat合法域名：`muhou.cuplay.top`（request）与 `resorce.cuplay.top`（downloadFile）
