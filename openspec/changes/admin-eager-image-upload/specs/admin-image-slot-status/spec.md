## ADDED Requirements

### Requirement: Per-slot upload status badge
Each gallery thumbnail (and cover dual field) SHALL show a Simplified Chinese status reflecting the slot's upload pipeline state.

#### Scenario: Status vocabulary
- **WHEN** a slot is compressing, queued, uploading, successfully uploaded, or failed
- **THEN** the UI shows a corresponding label such as 处理中…、排队中、上传中、已上传、失败

### Requirement: Remote URL shows uploaded
A slot that already has remote thumb/original URLs and no pending upload work SHALL display 已上传 (not 待上传).

#### Scenario: Loaded portfolio images
- **WHEN** an admin opens an existing portfolio whose gallery images already have CDN URLs
- **THEN** those slots show 已上传

### Requirement: Failed state allows retry
A failed upload status SHALL remain visible until the admin retries or removes the slot, and a retry control SHALL be available.

#### Scenario: Retry after failure
- **WHEN** an admin retries a failed slot
- **THEN** the slot re-enters the upload pipeline and status updates accordingly
