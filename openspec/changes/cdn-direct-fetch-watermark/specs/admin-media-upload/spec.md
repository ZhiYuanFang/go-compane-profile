# Spec Delta

## ADDED Requirements

### Requirement: 补水印原图浏览器直拉 CDN
为已上线内容图补加水印时，客户端 SHALL 使用浏览器对槽位已有 CDN URL（优先 original，其次 thumb）发起跨域 `fetch` 取得图片字节，再进入既有压缩/水印管线。系统 MUST NOT 再通过同源 admin 媒体代理端点代下这些字节。本机或其它未获 CDN CORS 放行的 Origin 直拉失败时，客户端 MUST 按现有下载失败路径提示错误，且 MUST NOT 回退到服务端代理。

#### Scenario: 生产管理端直拉 CDN 成功
- **WHEN** 运营在已配置 CORS 的管理端域名上对已上传内容图点击加水印
- **THEN** 浏览器直接向该图 CDN URL 请求图片字节
- **AND** 取得的字节进入后续压缩与水印，不请求 `/admin/api/media/proxy`

#### Scenario: 未放行 Origin 失败且不回退代理
- **WHEN** 当前页面 Origin 未被 CDN CORS 允许，导致直拉失败
- **THEN** 该槽位按下载失败处理并向运营展示错误
- **AND** 客户端不调用 admin 媒体代理作为兜底

## REMOVED Requirements

（无）
