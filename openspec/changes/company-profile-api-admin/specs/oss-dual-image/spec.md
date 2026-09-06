## ADDED Requirements

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
