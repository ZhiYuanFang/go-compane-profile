# deploy-runtime

## Purpose

deploy-runtime 域行为规格（自变更 delta / 旧基线迁移）。

<!-- migrated-from: changes/company-profile-api-admin/specs/deploy-acr, changes/rename-compose-service-muhou/specs/compose-service-muhou -->

## Requirements

### Requirement: Secrets in environment files not in git-tracked config
Sensitive values (MySQL credentials, OSS access keys, admin password, registry credentials used at runtime) SHALL be supplied via `.env.prod` (or CI secrets for build push), not committed as plaintext in application source or tracked yaml.

#### Scenario: Runtime reads env
- **WHEN** the service starts in production
- **THEN** it reads database, OSS, CDN, and admin password settings from environment / `.env.prod` injection

### Requirement: Single-image Docker build for ACR
The repository SHALL provide a Dockerfile that builds a single image containing the Go API and the built admin static assets, suitable for push to Aliyun ACR under the configured registry namespace.

#### Scenario: Image contains admin and API
- **WHEN** the Docker image is built
- **THEN** the running container can serve public API routes and the `/admin/` SPA

### Requirement: GitHub Actions workflow simplified from go_ai_talk pattern
The repository SHALL include a GitHub Actions workflow that builds and pushes the single image on version tags (and optional workflow_dispatch), using ACR credentials from GitHub Secrets, deriving the public push registry from a possibly VPC pull `REGISTRY` by stripping `-vpc` when present.

#### Scenario: Tag push builds image
- **WHEN** a version tag is pushed and ACR secrets are configured
- **THEN** the workflow builds the image and pushes tags including the primary version tag

### Requirement: Compose service is named muhou
The Docker Compose file for this application SHALL define a single application service whose Compose service key is `muhou`. The service SHALL continue to use image `${REGISTRY}/muhou:${IMAGE_TAG}` and container name `muhou`.

#### Scenario: Service key is muhou
- **WHEN** an operator inspects `manifest/docker/docker-compose.yml`
- **THEN** the application service key is `muhou` and not `compane-profile`

#### Scenario: Image path remains namespace plus muhou
- **WHEN** `REGISTRY` ends with `/compane-profile` and `IMAGE_TAG` is `v0.0.1`
- **THEN** the resolved image reference is `{REGISTRY}/muhou:v0.0.1` (i.e. ACR path `compane-profile/muhou:v0.0.1`)

### Requirement: Example REGISTRY uses compane-profile namespace
The tracked `.env.example` SHALL set `REGISTRY` to an ACR host plus the `compane-profile` namespace (not a trailing `/muhou` namespace), so that compose interpolation `${REGISTRY}/muhou` matches the intended repository.

#### Scenario: Example registry namespace
- **WHEN** a new operator copies `.env.example`
- **THEN** the example `REGISTRY` value ends with `/compane-profile` (placeholder host allowed)

### Requirement: Deploy docs refer to muhou service
Deploy documentation SHALL refer to the Compose service as `muhou` when describing compose operations for this stack.

#### Scenario: Docs use muhou service name
- **WHEN** an operator reads the deploy runbook for compose usage
- **THEN** service-level references use `muhou` rather than `compane-profile` as the Compose service key
