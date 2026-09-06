## Why

Docker Compose 服务键仍为 `compane-profile`，与镜像 repo、`container_name`（均为 `muhou`）以及产品对外命名不一致，运维命令与心智不统一。

## What Changes

- 将 `manifest/docker/docker-compose.yml` 中的服务名从 `compane-profile` 改为 `muhou`
- 同步注释 / `docs/deploy.md` 中凡引用该服务名的说明
- 修正 `.env.example` 的 `REGISTRY` 示例为 `…/compane-profile`（与 `${REGISTRY}/muhou` 拼出 `compane-profile/muhou` 一致）
- **不**修改 Compose project `name`（保持 `go-compane-profile`）
- **不**修改 Go module、Git 仓库名或 ACR 镜像路径公式（已是 `compane-profile/muhou`）

## Capabilities

### New Capabilities

- `compose-service-muhou`: Compose 单服务命名为 muhou，并与 ACR 镜像路径文档一致

### Modified Capabilities

- （无）

## Impact

- `manifest/docker/docker-compose.yml`、`docs/deploy.md`、`manifest/docker/env/.env.example`
- 已部署环境：服务键变更可能导致 compose 视为新服务；需停旧服务后再 `up`（见 design）
- CI workflow 镜像 tag 路径无需改动
