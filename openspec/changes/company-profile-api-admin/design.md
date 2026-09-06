## Context

本仓是 GoFrame 单仓模板（仅 `GET /hello`）。兄弟仓微信小程序「目后空间设计」公司名片已完成 UI，数据全在本地 `data/*.js` mock，design 约定「日后接 API：替换 data 模块，页面交互不变」。

同机已有 `go_ai_talk`：OSS bucket `pang-bao`、endpoint `oss-cn-beijing.aliyuncs.com`、CDN `https://resorce.cuplay.top`；ACR 与 `.env.prod` 发布习惯可复用并简化为单镜像。

约束：库名 `compane_profile`；域名 `muhou.cuplay.top`；后台固定密码 `abby`（经 env）；管理端 Vue 3 玻璃拟态+高阶灰；图片双份独立 objectKey 防 CDN 窜缓存。

## Goals / Non-Goals

**Goals:**

- 公开只读 API 对齐小程序字段，列表/详情使用 thumb + original
- 密码门控 CMS：公司、有序作品（含 renders/reals）、资费；双图上传
- Vue 3 管理端：选图本地预处理（原图≤5MB、缩略≤500KB），提交才上传有效图
- OSS 与 ucg 同套，前缀 `muhou/`；密钥仅 `.env.prod`
- 单镜像部署 + 简化 ACR workflow；小程序删 mock 并接线

**Non-Goals:**

- 小程序端微信登录、用户体系、留言/预约/支付
- 多管理员 RBAC、复杂权限
- 独立第二套 OSS / 自建对象存储
- 多语言 CMS、内容审核流水线
- 复刻 go_ai_talk 十微服务 matrix（本仓单服务即可）

## Decisions

### 1. 单进程托管 API + Admin SPA

Go 提供 `/api/v1/*`、`/admin/api/*`，并静态托管 Vue 构建产物于 `/admin/`。一个 Docker 镜像、一个 ACR 仓库名（如 `compane-profile`）。

**Alternatives considered:** 管理端独立容器 → 多一次编排与 CORS；本阶段收益低。

### 2. 管理端 Vue 3 + Vite

独立目录 `admin/`，玻璃拟态 + 高阶灰（深灰底、半透明面板、冷灰银强调；避免紫霓虹默认风）。构建产物拷入镜像由 Go 服务。

### 3. 数据模型

```
company (单例一行)
  intro_title, intro_body, values, years_label, address
  logo_original_url, logo_thumb_url
  logo_hor_original_url, logo_hor_thumb_url

portfolio
  id (bigint PK), slug/code (对外 id 字符串), sort_order
  cover_* , address, area, style, heart_flow
  created_at, updated_at

portfolio_image
  portfolio_id, kind (render|real), sort_order
  original_url, thumb_url

pricing (单例一行)
  original_url, thumb_url
```

作品对外 `id` 可用稳定 slug（后台可编辑）或数字转字符串；列表顺序严格按 `sort_order`，客户端不排序。

### 4. 公开 API 形状（示意）

- `GET /api/v1/company`
- `GET /api/v1/portfolios` → 每项 `cover` 为 thumb；可附带 `coverOriginal` 供详情预载（列表页可不使用）
- `GET /api/v1/portfolios/:id` → cover/renders/reals 均为 `{ thumb, original }`
- `GET /api/v1/pricing` → `{ thumb, original }`

无鉴权。响应字段名尽量贴近现有小程序 mock，降低接线成本。

### 5. 后台鉴权

`POST /admin/api/login` 校验 `ADMIN_PASSWORD`（默认 env 中为 `abby`），签发短期 session cookie 或 JWT；`/admin/api/*`（除 login）中间件校验。无用户表。

### 6. 双图上传流水线

```
管理端选图 → browser-image-compression（或 canvas）
  → 原图 Blob ≤5MB + 缩略 Blob ≤500KB（内存/IndexedDB 暂存）
  → 可反复更换，丢弃旧 Blob
  → 提交：multipart 成对 POST /admin/api/upload
  → Go 用 OSS SDK PutObject 至 muhou/{category}/{uuid}_o.ext 与 _t.ext
  → 返回 CDN URL 对，再写入业务表
```

禁止浏览器持有 AK；禁止仅依赖 OSS 图片处理同一 key + query（避免 CDN 窜缓存）。资费长图同样强制双图。

### 7. OSS / CDN

与 ucg 一致：`bucket=pang-bao`，`region=cn-beijing`，`endpoint=oss-cn-beijing.aliyuncs.com`，`cdnBaseUrl=https://resorce.cuplay.top`，`objectKeyPrefix=muhou/`。配置来自 env。

### 8. 配置与密钥

业务 yaml 不写密钥。服务器 `manifest/docker/env/.env.prod`（gitignore）含：MySQL、OSS AK、CDN、`ADMIN_PASSWORD`、`REGISTRY`、`IMAGE_TAG`。GitHub Actions：`ACR_USERNAME` / `ACR_PASSWORD` / Environment `REGISTRY`；workflow 从 pull 地址去掉 `-vpc` 得 push 公网地址（对齐 go_ai_talk 逻辑，单 matrix 一项）。

### 9. 小程序接线（兄弟仓）

- 新增 `utils/request.js` + `config`（baseURL=`https://muhou.cuplay.top`）
- 替换并删除 `data/company.js`、`data/portfolio.js` mock
- home：cover 用 thumb；detail：先 thumb，`Image`/`getImageInfo` 原图就绪后换 src
- 微信公众平台配置 request：`muhou.cuplay.top`；downloadFile：`resorce.cuplay.top`

### 10. 种子数据

迁移/启动脚本仅插入公司文案（来自现有 mock 文本）。logo/作品/资费为空，由后台上传。不导入三件 mock 作品。

## Risks / Trade-offs

- [长图缩略仍偏大] 资费图缩到 ≤500KB 可能损失可读性 → 管理端提示压缩质量；必要时提高缩略目标但规范仍以 500KB 为硬上限，可迭代算法
- [固定密码安全性低] 单人运营可接受 → HTTPS + 仅管理路径鉴权；密码可经 env 轮换无需改代码
- [空作品列表] 上线初期瀑布流为空 → 运营先在后台录入再发小程序，或接受短暂空态
- [同 bucket 误删] 前缀隔离 `muhou/` vs `social/`/`event/` → 删除逻辑仅允许本前缀 objectKey
- [微信合法域名未配] 真机请求失败 → 部署清单明确双域名；开发阶段可用开发者工具不校验

## Migration Plan

1. 创建库 `compane_profile`，跑表结构 + 公司文案种子
2. 部署服务到 `muhou.cuplay.top`，配置 `.env.prod` 与反向代理
3. 后台登录上传 logo / 作品 / 资费
4. 配置小程序合法域名，合并兄弟仓接线，删除 mock
5. 回滚：小程序可临时恢复 mock 分支；服务回退上一镜像 tag；DB 保留（勿 drop）

## Open Questions

- Nginx / TLS 证书由谁配置（若已有 cuplay.top 证书通配则可直接加 server 块）——实施时按现网惯例处理，不阻塞编码
- 作品对外 id：slug 字符串 vs 数字字符串——默认采用可编辑 `slug`，种子/新建时自动生成
