## 1. Size limit align

- [x] 1.1 Align `THUMB_MAX` to `0.5 * 1024 * 1024` (same as compress); keep original/thumb checks consistent with `compressUnder`

## 2. Parallel upload on save

- [x] 2.1 Add ordered concurrency helper (`mapPool`, concurrency 3)
- [x] 2.2 Portfolio `resolveImages`: parallel resolve pending duals with pool; preserve list order
- [x] 2.3 Show save-time upload progress「上传中 i/n」; fail save if any upload errors

## 3. Ship

- [x] 3.1 Rebuild admin and sync `resource/public/admin`
- [x] 3.2 Smoke: thumb in former 512–524KB band works; multi pending save uploads with concurrency ≤3 + progress
