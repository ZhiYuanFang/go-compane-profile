# Tasks

## 1. 前端直拉 CDN

- [x] 1.1 改写 `admin/src/utils/remoteImage.js`：对 http(s) CDN URL 使用浏览器 `fetch`（无 credentials、不走 `/admin/api/media/proxy`），校验 image blob 后返回 `File`；用代码检索确认无残留 proxy 路径
- [x] 1.2 确认 `DualImageField` / `ImageListEditor` 仍只调用 `fetchRemoteImageFile` 且无需改签名；本地快速阅读调用处无回归注释或死代码

## 2. 移除服务端媒体代理

- [x] 2.1 删除 `internal/controller/admin/media_proxy.go`，并从 `internal/cmd/cmd.go` 去掉 `GET /media/proxy` 注册；`rg media/proxy` 无匹配
- [x] 2.2 删除仅服务于代理的 `ReadObjectBytes` / `MaxProxyBytes`（及无引用辅助），`go build ./...` 通过

## 3. 校验

- [x] 3.1 `openspec validate cdn-direct-fetch-watermark --strict` 通过
- [ ] 3.2 在已放行 CORS 的管理端域名上：已上线图点「加水印」网络面板可见对 CDN 的直连请求、无 `/admin/api/media/proxy`，且进入压缩/水印管线

## Workflow follow-up

- 实现完成后用 `/opsx:archive` 归档并合并 delta 到主规格
