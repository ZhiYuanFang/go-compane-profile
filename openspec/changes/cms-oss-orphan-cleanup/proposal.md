## Why

CMS 换图、清空或删除作品时只更新 MySQL URL，OSS 上旧的 `muhou/` 对象成为孤儿并持续占存储。上传已用唯一 key，且 `DeleteObject` 已实现但从未被业务调用，需要在所有 CMS 图片变更路径接上同步清理。

## What Changes

- 公司 Logo、资费图、作品封面与图库在保存替换/清空后，对「旧有、新无」的 thumb/original CDN URL 做同步 OSS 删除
- 删除整条作品时，删除其封面与全部图库对应的 OSS 对象
- OSS 删除失败仅记录日志，不回滚、不失败业务写库
- 不引入异步队列、定时扫桶或 CDN 主动 purge

## Capabilities

### New Capabilities

- `cms-oss-cleanup`: 全部 CMS 图片变更时的同步 diff 删除与 best-effort 语义

### Modified Capabilities

- （无）`openspec/specs/` 尚未归档主规格；本变更新增 capability

## Impact

- 本仓后端：`internal/service/oss.go`（复用 `DeleteObject`）、`company.go`、`portfolio.go`、pricing 更新路径
- Admin API 契约不变；前端无需改动（服务端在保存/删除时 diff）
- 依赖现有 Aliyun OSS 凭证与 `muhou/` 前缀安全校验
