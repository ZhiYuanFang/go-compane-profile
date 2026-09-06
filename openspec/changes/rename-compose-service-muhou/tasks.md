## 1. Compose and env

- [x] 1.1 Rename Compose service key `compane-profile` → `muhou` in `manifest/docker/docker-compose.yml` (keep project `name`, image, container_name)
- [x] 1.2 Fix `manifest/docker/env/.env.example` `REGISTRY` to use `/compane-profile` namespace
- [x] 1.3 Fix misleading namespace comments in `manifest/docker/env/.env.prod` if present (do not change secrets)

## 2. Docs

- [x] 2.1 Update `docs/deploy.md` compose/service wording so the service is referred to as `muhou`
