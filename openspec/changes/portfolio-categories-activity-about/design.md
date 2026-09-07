## Context

本仓是公司名片的 Go API + Vue 后台；小程序在 sibling 仓 `wx-compane-profile-link`。现状：作品无类别、全局 `sort_order`；资费为单例双图；公司资料有地址/简介等，无电话/微信/奖项；无活动域。小程序首页为单列表瀑布流 + FAB「资费」。心流在详情用 `<view>` + `pre-wrap`，实机不换行；欢迎页用 `\n` split 可换行。

产品要一次做完：分类作品、活动位、关于我们替换资费、首页 Tab/ViewPager、多行文案修复。双边同发。

## Goals / Non-Goals

**Goals:**

- 作品四类（住宅/商业/办公/装置）+ 类内排序 + 公开可选 `category` 过滤
- 后台四条作品侧栏入口；允许改类
- 活动单例（封面+正文、可删）与小程序展示
- 公司资料扩展奖项/电话/微信；关于我们页；删除资费全链路
- 首页 Tab（全部+四类）+ ViewPager；活动条 80rpx
- 心流与奖项换行在小程序可见

**Non-Goals:**

- 存量作品自动迁移/默认类别（运营手改 DB）
- 多活动、活动排期、多正文图
- 简繁自动转换引擎
- Welcome 页展示奖项/电话/微信（仍只用现有公司简介字段）
- 后台作品列表全局「全部」视图

## Decisions

### 1. 类别枚举与存储

- DB / API 使用稳定英文：`residential` | `commercial` | `office` | `installation`
- 展示中文：住宅 / 商业 / 办公 / 装置
- `portfolio.category` NOT NULL；非法公开入参 → 400
- **Alternatives:** 中文直接入库 — 否决（排序/过滤/国际化脆弱）

### 2. 公开列表过滤与「全部」排序

- `GET /api/v1/portfolios?category=&page=&pageSize=`：省略 `category` = 全部
- 有 `category`：`WHERE category = ? ORDER BY sort_order ASC, id ASC`
- 全部：`ORDER BY sort_order ASC, id ASC`（同号仅用 `id` 稳定平局，不按类别再排）
- 分页契约保持现有 `page` / `pageSize` / `total`
- **Alternatives:** 「全部」按类分组拼接 — 否决（与「同号不分前后」冲突）

### 3. 类内排序与改类

- 后台每类独立列表；reorder 只接受本类 id 列表，写回 `0..n-1`
- 新建：入口所属类；`sort_order = max(该类)+1`
- 改类：允许；从旧类移除后，新类 `sort_order = max(新类)+1`；可选触发旧类稠密化（实现时：改类后旧类按当前顺序重写 `0..n-1`，避免空洞过大）
- **Alternatives:** 禁止改类 — 否决（产品明确允许）

### 4. 后台信息架构

- 侧栏四入口：`/portfolios/residential` 等（或 query `?category=`，路由需四条可点链接）
- 活动独立页；公司资料扩展字段；移除资费链接与路由
- 复用现有 `PortfolioListView` / `PortfolioEditView`，注入固定 `category`

### 5. 活动单例

- 新表 `activity`（或单行约定 `id=1`）：`cover_*`、`body_*` URL；无图即「空」
- 公开 `GET /api/v1/activity`：空则 `null`/空字段，小程序不渲染条
- 后台可上传封面/正文、删除（清空 URL）
- 首页封面容器高 **80rpx**，`mode=aspectFill`；正文页单图 `aspectFit` + 可滚动
- 镜像资费的 OSS dual-image 上传模式（category 如 `activity`）
- **Alternatives:** 复用/改造 `pricing` 表 — 否决（语义错误）

### 6. 关于我们与删资费

- 扩展 `company`：`awards`（TEXT，保留 `\n`）、`phone`、`wechat`；`address` 不变
- 公开 `GET /api/v1/company` 带回新字段
- 删除 `pricing` 表、公开/后台 API、Admin `PricingView`、小程序 `pages/pricing` 与 data
- 小程序新页 `pages/about`（或等价）：奖项可滚 + 固定底栏（地址复制、电话 `makePhoneCall`、微信复制）
- FAB 文案/图标改为关于我们

### 7. 多行文案（心流 / 奖项）

- 根因：微信里 `view` + `pre-wrap` 不可靠；欢迎页 `split('\n')` 可靠
- 详情心流、关于我们奖项：采用逐行渲染或 `<text>` 保证换行；不改 DB 存法
- API 原样返回含 `\n` 的字符串

### 8. 小程序首页 Tab + ViewPager

- Tab：全部 | 住宅 | 商业 | 办公 | 装置；默认全部；选中字重更大
- 横向 swiper 与 Tab 同步；每页独立拉列表（带对应 `category` 或省略）
- 活动条在 Tab **上方**；无活动数据则整块不占位或高度为 0
- 保留分页触底（每页独立 page 状态）

### 9. 仓与发版

- 实现 tasks 覆盖本仓与 `wx-compane-profile-link`
- 后端与小程序同窗口发布；先上 DB 列再发 API

## Risks / Trade-offs

- [「全部」同号交错] → 产品已接受；仅用 `id` 保证稳定顺序
- [改类后旧类顺序空洞] → 改类时重写旧类 sort_order
- [双边不同步] → 发版清单强制同窗口；旧小程序调已删 `/pricing` 会失败（可接受，因同步下线）
- [活动条 80rpx 偏矮] → 按产品定；外切可能裁切主体
- [公司资料页变重] → 可接受；不另开关于我们后台模块

## Migration Plan

1. 加 `portfolio.category`、`company` 新列、`activity` 表；部署 API
2. 运营为存量作品手写类别
3. 发后台（四入口、活动、公司字段、去资费）
4. 发小程序（首页/活动/关于我们/心流）
5. Rollback：保留列可空读旧客户端；pricing 删除后难回滚 → 发版前确认资费不再需要

## Open Questions

- （无阻塞项；活动正文 `aspectFit` 可滚、改类允许、单 change 已确认）
