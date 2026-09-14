## ADDED Requirements

### Requirement: Compress on select without OSS upload

Across all admin dual-image controls, when the operator selects or drops an image file, the system SHALL show a local preview and compress the file into pending original and thumb blobs, and MUST NOT call the OSS upload API until the enclosing form save begins.

#### Scenario: Select image on portfolio gallery

- **WHEN** the operator selects an image for a portfolio gallery slot
- **THEN** the slot shows a local preview and reaches a ready-with-pending state without a network upload request

#### Scenario: Select image on company logo

- **WHEN** the operator selects an image for the company logo field
- **THEN** the system compresses locally into pending blobs and does not upload to OSS until company save

### Requirement: Upload pending images on save with progress

When the operator saves an admin form that contains dual-image slots with pending blobs, the system SHALL upload those pending slots via the existing streaming upload path, SHALL show per-slot upload progress, and SHALL persist CMS payloads only after all required uploads for that save succeed.

#### Scenario: Save portfolio with new pending cover and gallery images

- **WHEN** the operator clicks save and one or more slots have pending blobs
- **THEN** those slots upload with progress UI and the create/update APIs are called only after every pending upload for that save completes successfully

#### Scenario: Upload failure blocks persistence

- **WHEN** any pending upload fails during save
- **THEN** the system does not complete the CMS write for that submit with missing image URLs, and the failed slot shows an error state for retry

### Requirement: Unsaved removal creates no OSS objects

When the operator removes or clears a dual-image slot that has never been persisted with remote URLs for the current edit (pending-only or still uploading locally), the system SHALL discard local preview and pending blobs and MUST NOT issue an OSS delete or upload request for that removal.

#### Scenario: Delete gallery image before save

- **WHEN** the operator uploads/selects a gallery image (compress completes) and then removes that slot without saving
- **THEN** no OSS upload or delete network request is triggered for that removal

### Requirement: Ready pending does not block save initiation

Slots that are compressed and waiting with pending blobs (not currently compressing or uploading, and not in error) MUST allow the operator to initiate save; the save flow itself performs the deferred uploads.

#### Scenario: Save enabled when images are ready locally

- **WHEN** all non-empty dual slots are either already remote-DONE or ready-with-pending, and none are compressing, uploading, or in error
- **THEN** the save control is not blocked by a “still uploading” reason solely due to pending blobs

### Requirement: Site-wide deferred behavior

Deferred upload SHALL apply consistently to every admin surface that uses dual-image upload (including company, portfolio cover and gallery, activity, and pricing), not only portfolio gallery.

#### Scenario: Activity image follows same pipeline

- **WHEN** the operator selects an activity cover image and later saves the activity
- **THEN** OSS upload occurs during save with progress, not immediately on select
