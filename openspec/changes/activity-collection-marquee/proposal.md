## Why

活动目前是单例双图（封面+正文图），无法运营多条活动，首页也只能展示一张封面。需要改为可增删改查的活动集合：首页纵向轮播标题、详情图文分离并用 HTML 正文、另提供全部活动列表入口。

## What Changes

- **BREAKING**：废止单例活动契约（`GET /api/v1/activity`、cover/body 双图模型、首页封面条）；旧单例数据清空重录，不做迁移
- 活动改为多行集合：标题、必填活动图（dual）、正文 HTML（字号/加粗/颜色；图文分离，正文内不插图）
- 公开 API：全量有序列表 + 详情（不分页）；admin：列表、创建、更新、删除、上下移排序
- 详情页顺序：图（满宽、高度自适应、点击放大）→ 标题加粗 → HTML 正文；图右下角「查看全部活动」进入列表
- 首页：纵向跑马灯遍历全部活动标题（停 1s，可手滑+自动播，一行尾省略）；无活动隐藏整条；点击进入对应详情
- 活动列表页：双列瀑布流圆角卡片（上图下标题，标题最多两行省略）
- 面向用户文案使用简体中文
- 公司简介、作品心流本 change **不做**富文本

## Capabilities

### New Capabilities

- `activity-collection-api`: 活动集合数据模型与公开/后台 API（CRUD、排序、全量列表）
- `activity-admin-cms`: 后台活动列表与编辑（标题、必填图、HTML 正文、排序）
- `miniprogram-activity-browse`: 首页纵向跑马灯、详情、全部活动列表瀑布流

### Modified Capabilities

- （无现有主规格目录；本 change 以新增能力规格描述行为，并显式废止既有单例活动行为）

## Impact

- **Backend** (`go-compane-profile`): `activity` 表重建/替换字段；service/controller/API；移除单例 cover/body 接口
- **Admin**: 活动由单页表单改为列表+编辑；引入正文富文本编辑器（HTML 子集）
- **Miniprogram** (`wx-compane-profile-link`): 首页活动条、详情页重做、新增活动列表页；`app.json` / data 层
- **Deploy**: 与小程序同窗口；DB 清空或按新 schema 重建 activity
