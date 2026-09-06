## ADDED Requirements

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
