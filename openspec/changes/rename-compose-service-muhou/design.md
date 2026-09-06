## Context

ACR 镜像已为 `{host}/compane-profile/muhou:{tag}`；compose 中 `image` 与 `container_name` 已用 `muhou`。服务键仍为 `compane-profile`，与命名不一致。产品决定仅改服务名，不改 Compose project name。

## Goals / Non-Goals

**Goals:**

- Compose 服务键改为 `muhou`
- 文档与示例 env 中 REGISTRY 命名空间与 `compane-profile/muhou` 一致

**Non-Goals:**

- 不改 `name: go-compane-profile`（project）
- 不改 Go module / 仓库目录名
- 不改 CI 镜像路径（已正确）
- 不改端口、env 业务变量

## Decisions

### 1. 仅重命名 service key

```yaml
services:
  muhou:
    image: ${REGISTRY}/muhou:${IMAGE_TAG}
    container_name: muhou
```

- **备选：** 同时改 project `name` 为 `muhou` — 否决（产品选仅服务名）

### 2. REGISTRY 示例修正

- `.env.example`：`…/muhou` → `…/compane-profile`
- `.env.prod` 注释若写「namespace muhou」改为「namespace compane-profile / image muhou」

### 3. 服务器迁移

改服务键后 compose 可能新建 `muhou` 服务并留下旧 `compane-profile` 服务定义痕迹；因 `container_name: muhou` 固定，需先：

```bash
docker compose --env-file env/.env.prod -f docker-compose.yml down
# 然后 pull && up
```

避免名称冲突。

## Risks / Trade-offs

- [旧服务键残留 / 容器名冲突] → 部署前 `down` 再 `up`
- [脚本写死 `compane-profile` 服务名] → 全文检索并更新文档

## Migration Plan

1. 合并 compose / docs / example 变更
2. 服务器：`down` → `pull` → `up -d --no-build`
3. 回滚：恢复旧 compose 服务键即可

## Open Questions

- 无。
