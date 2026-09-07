## Why

作品图库成组编辑当前强制每组必须有效果图，无法只传实景图。部分项目只有实景素材，需要两侧都可选，并保证小程序详情在「仅实景」时仍能正常出页、展示与预览。

## What Changes

- 后台成组图库：效果图与实景图均非必填；一组内至少一侧有图才保留；两侧都空时继续静默丢弃该组
- 保存仍提交等长的 `renders` / `reals`，缺侧用空 dual 占位对齐
- 小程序详情按组组装图页：支持仅效果、仅实景、双侧都有；默认优先效果图，否则实景；切换钮仅双侧都有时出现
- 封面仍必填（本 change 不放宽封面）
- 公开 API 数组契约不变（仍为 index 对齐的 `renders` / `reals`）

## Capabilities

### New Capabilities

- `optional-gallery-pair`: 后台成组双边可选与静默丢空组的保存规则
- `miniprogram-detail-optional-pair`: 详情页对仅实景/仅效果/双有组的展示与预览行为

### Modified Capabilities

- （无现有主规格目录；本 change 以新增能力规格描述行为。历史 `admin-portfolio-pair-gallery` 中「效果图必填、实景可选」由本 change 取代。）

## Impact

- **Admin** (`go-compane-profile/admin`): `PortfolioEditView.resolvePairs`、`PairImageEditor` 文案
- **Miniprogram** (`wx-compane-profile-link`): `pages/detail` 图页组装、默认图、toggle、preview URL
- **Backend**: 预期无契约变更；不改封面必填逻辑
