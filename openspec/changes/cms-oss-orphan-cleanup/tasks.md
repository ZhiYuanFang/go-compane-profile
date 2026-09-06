## 1. Shared cleanup helpers

- [x] 1.1 Add helpers to collect non-empty URLs from `DualURL` / slices and compute set difference (old − new)
- [x] 1.2 Add `deleteOSSURLsBestEffort` that calls existing `DeleteObject` per URL, logs failures, never returns a failing error to callers

## 2. Wire CMS write paths

- [x] 2.1 `UpdateCompany`: after successful DB write, diff-delete old logo / logoHor URLs no longer referenced
- [x] 2.2 `UpdatePricing`: after successful DB write, diff-delete previous pricing dual URLs when replaced or cleared
- [x] 2.3 `UpdatePortfolio`: after successful write, diff-delete replaced cover URLs; if gallery lists are included in the update, also diff-delete removed gallery URLs
- [x] 2.4 `SavePortfolioGallery`: after successful replace, diff-delete gallery URLs present before save but absent after
- [x] 2.5 `DeletePortfolio`: load cover + gallery URLs first, delete DB row, then best-effort delete collected OSS URLs

## 3. Verification

- [x] 3.1 Manually or with a focused test: replace/clear one image type and confirm orphan URLs are attempted for delete without failing the API on OSS errors
