## Context

作品图库以「效果图 + 实景图」成组编辑，保存为等长 `renders` / `reals`。当前 admin 在 `resolvePairs` 中强制每组有效果图；小程序详情用 `renders.map` 建页，仅实景组会空白或丢页。封面仍由现有封面字段/上传流程约束，本 change 不放宽。

## Goals / Non-Goals

**Goals:**

- 成组两侧均可空；至少一侧有图则保留；两侧都空静默丢弃
- 详情正确展示仅效果 / 仅实景 / 双有三种组
- 保持等长占位与公开 API 形状不变

**Non-Goals:**

- 放宽封面必填
- 改后端表结构或公开字段名
- 拆掉成组编辑模型（仍用 pair）

## Decisions

### 1. Keep equal-length pair arrays

- 缺侧继续写空 `{ thumb: '', original: '' }`，index 对齐不变
- **Alternatives:** 仅存非空、放弃等长 — 否决（破坏现有对齐约定）

### 2. Admin save rule

- `resolvePairs`：两侧都空 → `continue`；否则两侧都 `resolveDualImage`（空侧得到空 URL）
- 去掉「每组须包含效果图」
- `PairImageEditor` 文案改为两侧均可选（如「效果图 #n（可选）」「实景图 #n（可选）」），可加简短 hint「至少一侧」

### 3. Detail page assembly

- 按 `max(renders.length, reals.length)` 或假定等长遍历每一组
- `hasRender` / `hasReal`；默认：有 render 则显示 render，否则显示 real（`showReal = !hasRender && hasReal`）
- Toggle 仅当 `hasRender && hasReal`
- `previewImage` urls 跳过空侧；badge 文案随当前显示侧
- 跳过 `!hasRender && !hasReal` 的空组（防御性）

### 4. Cover unchanged

- 封面必填逻辑保持现状（本 change 不改 DualImageField/封面提交）

## Risks / Trade-offs

- [历史数据仅空 render 占位] → 详情按 hasRender/hasReal 处理即可
- [运营误加空组] → 静默丢弃，无提示（产品已确认）
- [仅实景时 badge 仍写「效果图」] → 随 `showReal` / 默认侧更新文案

## Migration Plan

1. 先发小程序详情（兼容旧「必有效果图」数据）
2. 再发 admin（允许仅实景保存）
3. Rollback：回退 admin 即可停止写入仅实景；旧详情对空 render 页可能仍差，故详情改动应优先或同发

## Open Questions

- （无；封面必填、空组静默丢弃已确认）
