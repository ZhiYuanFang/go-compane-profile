## Why

小程序首页仍是单一作品瀑布流，无法按住宅/商业/办公/装置浏览；缺少活动运营位与「关于我们」联络页，却仍保留已不再需要的资费入口。后台作品排序也是全局一把梭，无法按类维护。需要一次产品刷新：分类作品、活动、关于我们替换资费，并修好心流换行在详情页不生效的问题。

## What Changes

- 作品增加类别（住宅 / 商业 / 办公 / 装置）；公开列表 API 增加可选 `category` 入参，省略则返回全部
- 作品 `sort_order` 在**分类内**维护；「全部」列表仍按 `sort_order` 排序，同号不强制二次打破平局
- 后台侧栏改为四条作品入口（每类独立列表/排序/编辑）；允许编辑时改类（改后进入新类末尾）
- 存量仅一条作品：不写迁移默认值，运营手动改 DB 类别即可
- 新增活动单例：封面图 + 正文图，可删到空；公开 `GET` 活动接口；小程序首页 Tab 上方展示封面（高 80rpx、外切填满），点击进正文页（单图 `aspectFit`、可滚动）
- **BREAKING**：删除资费整域（表/API/后台页/小程序资费页与入口）
- 「公司资料」扩展奖项描述、电话、微信号；地址继续用 `company.address`；小程序原资费 FAB 改为「关于我们」联络页（奖项可滚 + 底栏地址/电话/微信）
- 心流与奖项等多行文案：后台换行在小程序侧须可见换行（心流详情当前为一行，需修复）
- 小程序首页改为顶部 Tab + ViewPager：默认「全部」，选中 Tab 字变粗变大
- 范围：`go-compane-profile`（API + admin）与 `wx-compane-profile-link` 同步发版

## Capabilities

### New Capabilities

- `portfolio-category`: 作品类别字段、公开/后台按类过滤与类内排序、改类规则、后台四入口
- `activity`: 活动单例（封面+正文、可删）、公开/后台 API 与小程序展示
- `company-about-contact`: 公司资料扩展联络/奖项字段、关于我们页、删除资费
- `miniprogram-home-category-tabs`: 首页 Tab + ViewPager 分类浏览、活动条、关于我们入口
- `miniprogram-multiline-copy`: 心流与奖项等换行在小程序侧正确展示

### Modified Capabilities

- （无现有主规格目录；本 change 以新增能力规格描述行为）

## Impact

- **Backend** (`go-compane-profile`): `portfolio` 表/服务/公开与后台 API；新建 `activity`；扩展 `company`；删除 `pricing` 全链路；`schema.sql` / seed
- **Admin**: 侧栏四作品入口 + 活动页；公司资料表单字段；移除资费页
- **Miniprogram** (`wx-compane-profile-link`): `pages/home` 重构；活动页；关于我们页替换资费；详情心流换行；`app.json` / data 层
- **Deploy**: 后端与小程序需同窗口发布；DB 需手工为存量作品设类别
