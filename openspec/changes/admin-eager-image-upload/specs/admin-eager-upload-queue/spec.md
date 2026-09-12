## ADDED Requirements

### Requirement: Eager upload after image add
When an admin adds or replaces a gallery or cover image (including batch select), the client SHALL compress the file and then upload it via the existing admin upload API without waiting for form save. Browser-direct OSS upload SHALL NOT be required.

#### Scenario: Single add triggers background upload
- **WHEN** an admin selects an image for a gallery slot or cover
- **THEN** the client compresses it and starts an upload to `POST /admin/api/upload` in the background
- **AND** on success the slot stores CDN URLs and clears local pending blobs

### Requirement: Batch queue with limited concurrency
When multiple images are enqueued, the client SHALL process uploads with a maximum concurrency of 3 so that earlier items may be uploading while later items remain queued or still compressing.

#### Scenario: Front uploading back waiting
- **WHEN** an admin batch-selects more than three images
- **THEN** at most three upload jobs run at once
- **AND** remaining jobs wait in queue until a slot frees

### Requirement: Delete cancels queued or in-flight work
Removing a gallery slot (single remove or batch delete) SHALL cancel that slot's queued job or abort an in-flight compress/upload when possible. A late successful upload response for a removed `slotId` MUST NOT re-add the image to the list.

#### Scenario: Delete while queued
- **WHEN** an admin deletes a slot that is still queued
- **THEN** that job never uploads

#### Scenario: Delete while uploading
- **WHEN** an admin deletes a slot that is uploading
- **THEN** the client aborts or ignores the result for that slotId
- **AND** the slot does not reappear after a late response

### Requirement: Save requires uploads settled
Portfolio save SHALL NOT succeed while any gallery/cover slot is still compressing, queued, uploading, or in a failed upload state. Slots that already have remote URLs SHALL be submitted as URLs without re-uploading.

#### Scenario: Block save on failure
- **WHEN** any slot shows a failed upload status
- **THEN** save does not complete successfully until the admin retries successfully or removes the slot

#### Scenario: Block or wait while in flight
- **WHEN** uploads are still in progress
- **THEN** save is blocked or deferred with clear feedback until uploads finish or fail
