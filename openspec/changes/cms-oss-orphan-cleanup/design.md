## Context

每次上传生成唯一 `muhou/{category}/{guid}_o|_t` 对象；CMS 保存只改 MySQL 中的 CDN URL。`service.DeleteObject` 已支持从 CDN URL 反解 key 并限制 `muhou/` 前缀，但公司、作品、资费、删除作品路径均未调用，导致换图/删图产生孤儿对象。

## Goals / Non-Goals

**Goals:**

- 全部 CMS 图片路径在写库成功后，同步删除不再引用的 OSS 对象（thumb + original）
- 删除整条作品时清理其全部关联图片对象
- OSS 删除失败仅日志，不失败、不回滚业务写库
- 接受 CDN 短暂仍可命中已删对象

**Non-Goals:**

- 异步队列、定时扫桶对账、引用计数
- CDN 主动 purge / 刷新缓存
- 前端直连 OSS 删除或新增删除 API
- 清理历史已存在的孤儿（可选后续运维脚本，不在本变更）

## Decisions

### 1. 服务端 diff，先 DB 后 OSS

- **选择：** 读旧 URL → 事务/写库成功 → `toDelete = oldURLs \ newURLs` → best-effort `DeleteObject`
- **备选：** 先删 OSS 再写库 — 否决（写库失败会裂图）
- **备选：** 仅前端在换图时调删除 API — 否决（密钥与一致性应在服务端）

### 2. 全 CMS 挂点

| 路径 | 旧集合 | 新集合 |
|------|--------|--------|
| `UpdateCompany` | 旧 logo + logoHor | 新 logo + logoHor |
| `UpdatePricing` | 旧 dual | 新 dual |
| `UpdatePortfolio` | 旧 cover（+ 若本次带 gallery 则含旧 gallery） | 新 cover（+ 新 gallery） |
| `SavePortfolioGallery` | 旧 renders+reals | 新 renders+reals |
| `DeletePortfolio` | 封面 + 全部图库 | 空 |

说明：当前 admin 编辑常分「更新字段/封面」与「put gallery」；两处各自 diff，避免漏删。

### 3. URL 集合差

- 将每个 `DualURL` 的非空 `thumb`、`original` 加入集合
- 空串与空 dual 跳过（不成组删除目标）
- 未变化的 URL 保留在新集合中 → 不删
- 上传 key 唯一，默认不做跨实体引用检查

### 4. Best-effort 删除辅助

- 封装如 `deleteOSSURLsBestEffort(ctx, urls []string)`：对每个 URL 调 `DeleteObject`，失败 `g.Log().Warning`（或等价），继续下一个
- 不向调用方返回错误（或仅内部统计）

### 5. 删除作品顺序

- 先加载封面与 gallery URL 到内存
- DB delete（级联 `portfolio_image`）
- 再 best-effort 删收集到的 URL  
若 DB 删失败则不删 OSS。

## Risks / Trade-offs

- [OSS 删失败残留孤儿] → 接受；日志可人工补删；日后可加对账
- [CDN 短暂命中旧图] → 接受（产品已确认）
- [误删仍被他处引用的 URL] → 当前上传不复用 key，风险低；若将来复用需引用计数
- [UpdatePortfolio 不带 gallery 时] → 只 diff cover，避免把未提交的 gallery 当成「已清空」

## Migration Plan

1. 部署含清理逻辑的新镜像即可；无 schema 变更
2. 回滚：回退镜像后行为恢复为「只写库不删 OSS」；已删对象无法自动恢复
3. 历史孤儿不在本变更清理

## Open Questions

- 无阻塞项。
