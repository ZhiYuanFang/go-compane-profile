# oss-media

## Purpose

oss-media 域行为规格（自变更 delta / 旧基线迁移）。

<!-- migrated-from: changes/cms-oss-orphan-cleanup/specs/cms-oss-cleanup, changes/company-profile-api-admin/specs/oss-dual-image -->

## Requirements

### Requirement: Diff-delete unused CMS image objects on save
When an authenticated admin successfully saves company logos, pricing image, portfolio cover, or portfolio gallery images, the system SHALL delete OSS objects for CDN URLs that were previously stored and are no longer present in the saved payload. Empty URL strings SHALL be ignored. Deletion SHALL occur only after the database write succeeds.

#### Scenario: Replace cover deletes previous pair
- **WHEN** an admin replaces a portfolio cover with a newly uploaded dual-URL pair and the save succeeds
- **THEN** the previous cover thumb and original objects under the configured `muhou/` prefix are deleted from OSS (best-effort)

#### Scenario: Unchanged gallery URL is retained
- **WHEN** an admin saves a gallery that still references an existing render URL and only adds or changes other items
- **THEN** OSS objects for the unchanged render URL are not deleted

#### Scenario: Cleared real slot deletes prior real objects
- **WHEN** an admin saves a gallery pair whose real dual-URL becomes empty after previously having non-empty URLs
- **THEN** the previous real thumb and original objects are deleted from OSS (best-effort)

### Requirement: Delete portfolio removes all related OSS images
When an authenticated admin successfully deletes a portfolio, the system SHALL best-effort delete OSS objects for that portfolio's cover and all gallery dual-URL images collected before the database delete.

#### Scenario: Portfolio delete cleans cover and gallery
- **WHEN** an admin deletes a portfolio that has a cover and gallery images
- **THEN** after the portfolio row is removed, OSS delete is attempted for each non-empty cover and gallery thumb/original URL

### Requirement: Company and pricing image replacement cleanup
Company logo / horizontal logo updates and pricing image updates SHALL use the same post-save URL diff cleanup as portfolios.

#### Scenario: Pricing image replaced
- **WHEN** an admin successfully saves a new pricing dual-URL pair that differs from the previous pair
- **THEN** the previous pricing thumb and original OSS objects are deleted (best-effort)

#### Scenario: Company logo cleared
- **WHEN** an admin successfully clears a company logo dual-URL that previously had URLs
- **THEN** the previous logo thumb and original OSS objects are deleted (best-effort)

### Requirement: Best-effort OSS delete must not fail CMS writes
OSS object deletion failures MUST NOT fail or roll back a successful CMS database write. The system SHALL log deletion failures and continue. Brief CDN continued serving of recently deleted objects is acceptable.

#### Scenario: OSS delete error does not fail save
- **WHEN** database save succeeds but one or more OSS deletes fail
- **THEN** the API still returns success for the CMS write and logs the delete failure

### Requirement: Only muhou-prefixed objects may be deleted
Cleanup SHALL use the existing server-side object delete path that resolves CDN URLs to keys and refuses keys outside the configured OSS prefix (default `muhou/`).

#### Scenario: Non-prefix URL skipped safely
- **WHEN** a stored URL cannot be resolved to an allowed `muhou/` object key
- **THEN** that URL is not deleted as an arbitrary bucket object and the overall CMS write remains successful

### Requirement: Browser-side dual-size preparation before upload
The admin console SHALL, after the operator selects an image, silently produce two local blobs: an original candidate at most 5MB and a thumbnail candidate at most 500KB. The operator MAY replace the selection before submit; only the current pair is eligible for upload.

#### Scenario: Replace selection before submit
- **WHEN** the operator selects image A, then selects image B, then submits
- **THEN** only B's processed original and thumbnail are uploaded, and A's blobs are discarded

#### Scenario: No upload until submit
- **WHEN** the operator selects and processes images but does not submit the form
- **THEN** no objects are written to OSS for those selections

### Requirement: Server-mediated paired OSS upload
The system SHALL accept authenticated admin uploads of original + thumbnail pairs and store them under the `muhou/` prefix in the shared Aliyun OSS bucket, using two distinct object keys. Uploaded objects MUST be readable via the configured CDN base URL. Access keys MUST come from environment configuration, never from the admin frontend bundle.

#### Scenario: Successful paired upload
- **WHEN** an authenticated admin uploads a valid original (≤5MB) and thumbnail (≤500KB) pair
- **THEN** the API returns CDN URLs for both objects under `https://resorce.cuplay.top/muhou/...`

#### Scenario: Reject oversized upload
- **WHEN** an upload exceeds the original 5MB or thumbnail 500KB limit
- **THEN** the API rejects the upload without creating OSS objects

### Requirement: Distinct keys to avoid CDN cache collision
Original and thumbnail for the same logical image MUST use different object keys (not the same key with processing query parameters).

#### Scenario: Separate object keys
- **WHEN** a dual-image pair is stored
- **THEN** the original and thumbnail object keys differ and each resolves independently on the CDN

### Requirement: Shared OSS configuration with muhou prefix
OSS connection SHALL reuse the same bucket, region, endpoint, and CDN base as go_ai_talk ucg (`pang-bao`, `cn-beijing`, `oss-cn-beijing.aliyuncs.com`, `https://resorce.cuplay.top`) with object key prefix `muhou/`.

#### Scenario: Prefix isolation
- **WHEN** the service uploads an object
- **THEN** the object key begins with `muhou/`
