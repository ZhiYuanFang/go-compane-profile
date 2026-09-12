## ADDED Requirements

### Requirement: Immediate local preview on file select
When an admin selects one or more images for gallery slots (single replace, empty pick, or batch), the client SHALL show a local preview for each selected file as soon as the file is associated with a slot, without waiting for compression or upload to start or finish.

#### Scenario: Batch select shows all previews
- **WHEN** an admin batch-selects multiple images
- **THEN** each new slot displays a local preview of its file promptly
- **AND** previews appear even for slots that are still waiting to compress or upload

#### Scenario: Replace updates preview immediately
- **WHEN** an admin replaces an image on a slot
- **THEN** the cell preview updates to the newly selected file before upload completes

### Requirement: Preview available while busy
While a slot is compressing, queued, or uploading, if a local or pending preview exists, the admin SHALL be able to open that preview (e.g. lightbox). Busy state MUST NOT blank or hide an available preview.

#### Scenario: Lightbox during upload
- **WHEN** a slot is uploading and has a local or pending preview
- **THEN** the admin can open the preview viewer for that slot
