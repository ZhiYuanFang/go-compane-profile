# 项目上下文

本文件是仓库 **唯一** AI / OpenSpec 项目入口（Cursor 规则与 `openspec/config.yaml` 均指向此处）。**不再维护** 根目录 `AGENTS.md`。行为级 MUST 以 `openspec/specs/<domain>/spec.md` 为准；下文只提供名片、布局与短指针。

## 项目目标

`go-compane-profile`：公司形象站后端与管理端（公开作品/公司资料 API + Admin CMS/控制台）。对外域名见 `docs/deploy.md`（如 `muhou.cuplay.top`）。微信小程序客户端规格在兄弟仓维护，本仓 **不** 将 `miniprogram-*` 等纯客户端能力写入 Go 域真相。

## 技术栈

- Go / GoFrame
- MySQL（库名 `compane_profile`）
- 阿里云 OSS + CDN（前缀 `muhou/`）
- Docker Compose（服务键 `muhou`）+ 阿里云 ACR
- Admin 控制台（SPA，随本仓部署）
- OpenSpec CLI 2.1.x（`spec-driven`：proposal → specs → design → tasks）

## OpenSpec 布局

- **真相**：`openspec/specs/<domain>/spec.md`（十域，见下）。
- **变更**：`openspec/changes/<change-id>/`；delta 位于 `changes/<id>/specs/<domain>/spec.md`。
- **归档**：合并 delta 进主规格（若实现阶段已写妥主规格可 `--skip-specs`），然后 **删除**该 change 目录；默认 **不保留** `openspec/changes/archive/`。
- **制品语言**：简体中文；每条 `### Requirement` 须含英文 **SHALL** 或 **MUST**，并含 `#### Scenario`。
- **禁止** 版本巨石 `openspec/specs/v*.md` 与收版巨石脚本。

### 十域

| 域 | 摘要 |
|---|---|
| `public-portfolio-api` | 公开公司资料、作品列表/详情/分页等未认证 API |
| `admin-auth` | 管理端登录、会话与鉴权 |
| `admin-cms` | 管理端内容 CRUD（作品、公司、活动等写路径） |
| `admin-console-ui` | Admin SPA 交互与展示（画廊网格、预览、拖拽、槽位状态等） |
| `admin-media-upload` | 管理端图片压缩/队列/并发/进度与上传管线 |
| `oss-media` | OSS 双图（thumb/original）、孤儿清理与 CDN 约定 |
| `activity` | 活动集合公开 API 与活动 CMS |
| `portfolio-category` | 作品分类模型与 API/CMS |
| `view-counts` | 浏览量计数 API 与管理端展示 |
| `deploy-runtime` | Compose/ACR、服务名 `muhou`、运行端口与部署约定 |

## 工程约束指针

实现代码前对照对应域规格（勿在本文件重复长文 MUST）：

- 公开 `GET /api/v1/*` → `public-portfolio-api` / `activity` / `portfolio-category` / `view-counts`
- Admin 会话与密码 → `admin-auth`
- Admin 写模型与 CMS → `admin-cms`
- 控制台 UX（画廊、槽位色、拖拽）→ `admin-console-ui`
- 选图后上传队列/并发/进度 → `admin-media-upload`
- OSS 双 URL、清理 → `oss-media`
- Compose / ACR / `muhou` 服务键 → `deploy-runtime`
- 小程序 UI/接线 → 兄弟仓，勿并入本仓十域

## 重要约束（评审检查项）

- 行为变更是否落到某域 `### Requirement` / `#### Scenario`，禁止规格悬空。
- 是否误用已删除的 `v*.md` 巨石。
- 归档是否合并（或 skip-specs）并删除 change，而非堆 `changes/archive/`。
- 是否把纯小程序能力误写入 Go 域真相。
- Compose 是否仍以 `muhou` 服务部署且配置无明文密钥。
