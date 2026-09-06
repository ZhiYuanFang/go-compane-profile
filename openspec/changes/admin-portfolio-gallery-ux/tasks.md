## 1. Admin image preview

- [x] 1.1 Change `DualImageField` preview CSS to `object-fit: contain` without changing preview box size
- [x] 1.2 Add click-to-enlarge lightbox on non-empty `DualImageField` preview (prefer original / pending preview; dismissible)
- [x] 1.3 Verify enlarge works on company logos, portfolio cover, pricing, and gallery pair slots

## 2. Admin portfolio pair gallery

- [x] 2.1 Add pair gallery editor component (`pairs` of render + optional real) with add / reorder / remove whole pairs
- [x] 2.2 Wire `PortfolioEditView` to zip/unzip `renders`/`reals` with equal-length empty dual placeholders for empty reals
- [x] 2.3 Remove slug field from create/edit UI; submit empty slug on create
- [x] 2.4 After successful create or update, navigate to `/portfolios`
- [x] 2.5 Remove or stop using separate dual `ImageListEditor` instances on the portfolio editor

## 3. Admin build

- [x] 3.1 Build admin and sync static assets to `resource/public/admin` if that is the serving path

## 4. Mini-program detail native preview (sibling repo `wx-compane-profile-link`)

- [x] 4.1 On non-cover gallery main image tap, call `wx.previewImage` with `current` and a URL list of non-cover non-empty originals (exclude cover and empty reals)
- [x] 4.2 Keep cover free of this native-preview handler; keep real/render toggle on `catchtap`
- [x] 4.3 Fix pair building: index-align reals; never fall back to `reals[0]`; show toggle only when that page’s real is non-empty
- [x] 4.4 Remove ineffective page-local `movable-view` scale / zoom-lock if no longer needed; confirm vertical swiper paging still works
