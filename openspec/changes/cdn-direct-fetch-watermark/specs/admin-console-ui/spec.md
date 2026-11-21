# Spec Delta

## MODIFIED Requirements

### Requirement: 补水印时展示远程下载状态
为已上传图片补加水印时，客户端 SHALL 在浏览器直拉 CDN 原图期间将对应槽位置为下载状态并显示「下载原图…」。下载完成后 SHALL 自动进入后续排队/压缩/添加水印状态；下载失败（含 CORS 拒绝）时槽位 MUST 恢复为「已上传」，并向运营展示可定位到具体序号的错误。

#### Scenario: 已上线图片补水印显示下载状态
- **WHEN** an admin selects already-uploaded images and clicks 加水印
- **THEN** every selected slot immediately shows 下载原图…
- **AND** as each download finishes, its slot moves on to 排队中/处理中…/添加水印中…

#### Scenario: 下载失败恢复已上传
- **WHEN** fetching a selected slot's remote original from CDN fails
- **THEN** that slot returns to 已上传 instead of entering a failed state
- **AND** the error message names the slot by its index
- **AND** the other selected slots continue processing
