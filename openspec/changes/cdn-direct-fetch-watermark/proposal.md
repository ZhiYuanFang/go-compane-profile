# Proposal

## Why

补加水印时拉取已上线原图目前经 `/admin/api/media/proxy` 由服务端代下 CDN，只为绕过跨域。CDN 已对管理端域名配置 CORS 后，代理成为多余跳转，增加服务端带宽与内存占用，应改为浏览器直拉 CDN。

## What Changes

- 管理端 `fetchRemoteImageFile` 改为对已有 CDN URL 做浏览器直连 `fetch`（不经 admin API）
- **BREAKING（对内）**：移除 `GET /admin/api/media/proxy` 及服务端 `ReadObjectBytes` 代理读图能力
- 补水印后的压缩/水印/待保存/上传管线不变；UI「下载原图…」状态保留
- 明确接受：仅生产管理端 Origin（如 `https://muhou.cuplay.top`）可直拉；本机 Vite Origin 被 CORS 拦截时不提供本地代理兜底

## Capabilities

### New Capabilities

（无）

### Modified Capabilities

- `admin-media-upload`: 已上线图补水印的原图字节 MUST 由浏览器对 CDN URL 直拉获得；系统 MUST NOT 再依赖同源媒体代理端点
- `admin-console-ui`: 补水印「下载原图…」状态语义改为浏览器直拉 CDN；失败恢复行为不变

## Impact

- 前端：`admin/src/utils/remoteImage.js`（及仍调用它的 DualImageField / ImageListEditor）
- 后端：删除 `internal/controller/admin/media_proxy.go`、路由注册、`service.ReadObjectBytes` / `MaxProxyBytes`（若无其它引用）
- 依赖：生产 CDN CORS 已配置（`Access-Control-Allow-Origin` 含管理端域名）；本仓库不负责 CDN 控制台配置
- 本地开发：本机 Origin 直拉失败为预期，不做 Vite/API 代理兜底
