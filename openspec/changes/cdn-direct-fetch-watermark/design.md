# Design

## Context

补水印路径：`DualImageField` / `ImageListEditor` → `fetchRemoteImageFile` → 压缩/水印管线。当前 `fetchRemoteImageFile` 请求同源 `/admin/api/media/proxy`，由 `HandleMediaProxy` + `ReadObjectBytes` 从公共 CDN 拉字节再回传。CDN 已对 `https://muhou.cuplay.top` 返回 `Access-Control-Allow-Origin` 与 `GET,HEAD`。动机见 `proposal.md` Why。

## Goals / Non-Goals

**Goals:**

- 浏览器对已有 CDN URL 直拉图片字节，再复用现有压缩/水印入口
- 删除 admin 媒体代理端点与仅为其服务的服务端读图代码
- 保持槽位「下载原图…」与失败回「已上传」的 UX

**Non-Goals:**

- 不为本机 Vite Origin 配置 CORS 或本地代理兜底
- 不改水印绘制、压缩上限、保存上传、OSS 清理
- 不在本仓配置阿里云 CDN CORS（视为运维已完成）

## Decisions

1. **直拉入口仍集中在 `remoteImage.js`**  
   只改 `fetchRemoteImageFile`：对绝对 CDN URL 使用 `fetch(url)`（不加 `credentials`），校验 `Content-Type` / blob 为 image，再包成 `File`。调用方不变。  
   *备选*：各组件各自 `fetch` — 拒绝，避免重复错误处理。

2. **URL 形态**  
   槽位已存完整 CDN URL（`https://resorce.cuplay.top/muhou/...`）。直拉只处理可 `fetch` 的 http(s) URL；若仅有 object key（无 scheme），报错「没有可处理的线上图片」或等价文案，不拼 CDN_BASE（前端不应硬编码密钥侧配置；生产数据本就是完整 URL）。  
   *备选*：前端读 `CDN_BASE` 拼 URL — 当前无此环境变量暴露，不引入。

3. **删除代理整条链路**  
   移除路由 `GET /media/proxy`、`media_proxy.go`、`ReadObjectBytes`、`MaxProxyBytes`。  
   *备选*：保留代理作 fallback — 明确拒绝（规格禁止回退）。

4. **错误文案**  
   网络/CORS/`!ok`/非图片时抛出可读中文错误；去掉「GoFrame JSON 鉴权 200」分支（直拉 CDN 不再走 admin 会话包装）。

## Risks / Trade-offs

- [CORS 配置回退或漏 Origin] → 生产补水印失败；依赖运维保持 CDN CORS；失败走既有 UI 提示  
- [本机开发无法测补水印直拉] → 接受；在正式管理端验证  
- [超大原图直拉占浏览器内存] → 与原先经代理回传相同量级；仍受后续压缩上限约束  
- [公开 CDN 可被任意放行 Origin 的页面 fetch] → 图片本就公开可读；去掉代理鉴权可接受

## Migration Plan

1. 部署前确认生产 CDN CORS 仍放行管理端 Origin  
2. 发布含直拉前端 + 去掉 proxy 的后端同一版本  
3. 回滚：恢复上一版二进制/静态资源即可恢复代理路径（CDN CORS 可保留无害）

## Open Questions

（无）
