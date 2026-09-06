## ADDED Requirements

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
