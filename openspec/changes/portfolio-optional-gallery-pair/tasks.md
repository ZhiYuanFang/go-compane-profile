## 1. Admin pair gallery

- [x] 1.1 Update `PairImageEditor` labels so both render and real are optional
- [x] 1.2 Change `PortfolioEditView.resolvePairs` to accept real-only pairs; keep silent skip when both empty; keep equal-length empty dual placeholders
- [x] 1.3 Confirm cover required behavior is unchanged

## 2. Mini-program detail

- [x] 2.1 Build gallery pages from aligned pairs with `hasRender` / `hasReal`; skip both-empty; default to render when present else real
- [x] 2.2 Show render/real toggle only when both sides exist; update badge/copy for the displayed side
- [x] 2.3 Fix native preview URL list and `current` to skip empty sides (including real-only pages)

## 3. Verification

- [x] 3.1 Smoke admin: save render-only, real-only, both, and mixed empty pairs (empty dropped)
- [x] 3.2 Smoke detail: real-only page displays; toggle only on dual pairs; preview works
