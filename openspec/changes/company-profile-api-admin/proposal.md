## Why

微信小程序「目后空间设计」公司名片目前全靠本地 mock，无法在线上更新公司文案、作品集与资费图。需要本仓提供公开只读 API、可运营的后台 CMS，并把小程序接到真实数据源。

## What Changes

- 新建 MySQL 库 `compane_profile` 与公司 / 作品 / 资费等内容表
- 提供小程序公开只读 API（公司、作品列表/详情、资费），图片一律返回 CDN 上的原图+缩略图双 URL
- 提供固定密码门控的管理 API（登录、CRUD、经服务端代理的双图上传）
- 新建独立 Vue 3 管理端（玻璃拟态 + 高阶灰）：选图后浏览器端预处理双尺寸，提交时再统一上传有效图片
- 图片存阿里云 OSS（与 go_ai_talk ucg 同 bucket/endpoint/CDN，前缀 `muhou/`）
- 敏感配置进 `.env.prod`；单镜像 ACR 构建 workflow（简化自 go_ai_talk）
- 公网域名 `https://muhou.cuplay.top`（API + `/admin` SPA）
- 兄弟仓小程序：去掉本地 mock，接线 API；列表用缩略图，详情先缩略图再静默换原图
- 种子数据仅灌入公司文案；不灌 mock 作品；删除小程序侧现有 mock 数据模块

## Capabilities

### New Capabilities

- `public-portfolio-api`: 小程序公开只读接口与双 URL 响应约定
- `admin-auth`: 后台固定密码登录与会话保护
- `admin-cms`: 公司资料、作品集（含有序图集）、资费的管理读写
- `oss-dual-image`: 浏览器预处理 + Go 代理上传 OSS，原图≤5MB / 缩略图≤500KB，独立 objectKey
- `admin-console-ui`: Vue 3 玻璃拟态管理控制台
- `miniprogram-api-wire`: 兄弟仓小程序 data 层换 API、删 mock、渐进换图
- `deploy-acr`: `.env.prod`、单镜像 Dockerfile、GitHub Actions 推送 ACR

### Modified Capabilities

- （无；本仓尚无既有 specs）

## Impact

- **本仓** `go-compane-profile`：从 GoFrame hello 模板扩展为业务 API + 静态托管管理端；新增 MySQL、OSS SDK、Vue 3 前端工程、Docker/CI
- **兄弟仓** `微信小程序/compane-profile`：`data/*` mock 删除并改为网络请求；合法域名需配置 `muhou.cuplay.top` 与 `resorce.cuplay.top`
- **基础设施**：MySQL 新建库；OSS `pang-bao` 下新增 `muhou/` 前缀；Nginx/域名指向 `muhou.cuplay.top`；ACR 命名空间 `muhou`
- **密钥**：OSS AK、MySQL、`ADMIN_PASSWORD`、ACR 凭证仅存 `.env.prod` / GitHub Secrets，不进仓库明文
