# admin-media-upload

## Purpose

admin-media-upload 域行为规格（自变更 delta / 旧基线迁移）。

<!-- migrated-from: changes/admin-deferred-oss-upload/specs/admin-deferred-upload, changes/admin-eager-image-upload/specs/admin-eager-upload-queue, changes/admin-eager-image-upload/specs/admin-image-size-3m, changes/admin-gallery-preview-fix-concurrency-1/specs/admin-upload-concurrency-1, changes/admin-image-upload-parallel/specs/admin-image-size-align, changes/admin-image-upload-parallel/specs/admin-upload-concurrency, changes/admin-oss-upload-progress-stream/specs/admin-oss-progress-stream, changes/admin-upload-pipeline-ab/specs/admin-compress-upload-decouple -->

## Requirements

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

### Requirement: Original max 3MiB
Client compression and server upload validation SHALL enforce an original image maximum of 3 MiB (`3 * 1024 * 1024` bytes).

#### Scenario: Client compress target
- **WHEN** a selected original exceeds 3MiB
- **THEN** the client attempts to compress it under 3MiB before upload

#### Scenario: Server rejects oversized original
- **WHEN** an upload presents an original larger than 3MiB
- **THEN** the server rejects it with an invalid-parameter style error

### Requirement: Thumb max aligned at 0.5MiB
Client and server SHALL both treat the thumbnail maximum as `0.5 * 1024 * 1024` bytes so a thumb accepted by client compression is not rejected solely for exceeding `500 * 1024`.

#### Scenario: Client and server agree
- **WHEN** a thumb size is between `500*1024` and `0.5*1024*1024` bytes and within the compress success bound
- **THEN** both client acceptance and server validation allow it

### Requirement: Upload and compress concurrency is one
The shared admin image compress pool and upload pool SHALL each run with a maximum concurrency of 1 (at most one compress job and one HTTP upload in flight at a time across the shared queues).

#### Scenario: Second upload waits for first
- **WHEN** two slots are ready to upload
- **THEN** only one HTTP upload request runs at a time
- **AND** the second starts only after the first completes, fails, or is cancelled

#### Scenario: Second compress waits for first
- **WHEN** two slots need compression
- **THEN** only one compression runs at a time
- **AND** slots waiting for compress MUST still show their local previews if set

### Requirement: Aligned thumbnail size limit
Client-side image processing SHALL use the same byte ceiling for thumbnail success checks as the compression target of 0.5 MiB (`0.5 * 1024 * 1024` bytes). A thumbnail that successfully finishes compression under that target MUST NOT be rejected solely for exceeding `500 * 1024` bytes.

#### Scenario: Thumb in former dead zone succeeds
- **WHEN** compression produces a thumbnail whose size is between `500*1024` and `0.5*1024*1024` bytes inclusive of the compress success bound
- **THEN** `processDualImage` accepts it as a valid thumb
- **AND** does not throw「缩略图超过 500KB…」for that reason alone

### Requirement: Consistent limit constants
Original and thumb max sizes SHALL be derived from the same unit convention as `compressUnder`'s `maxSizeMB * 1024 * 1024` checks so post-conditions cannot disagree with compress success.

#### Scenario: Single source of truth
- **WHEN** thumb max is defined for processing
- **THEN** it equals `0.5 * 1024 * 1024` (or an equivalent expression of 0.5 MiB)

### Requirement: Parallel pending uploads with concurrency limit
When resolving multiple pending dual images for portfolio save, the admin client SHALL upload those pending items concurrently with a maximum concurrency of 3 (not strictly one-after-another). Items that already have remote URLs and no pending blobs SHALL not be uploaded again.

#### Scenario: Multiple pending uploads overlap
- **WHEN** an admin saves a portfolio with more than three pending gallery images
- **THEN** at most three upload requests are in flight at once
- **AND** all pending images are uploaded before the portfolio/gallery write continues

#### Scenario: Preserve gallery order
- **WHEN** pending images at positions 1..N are resolved via the pool
- **THEN** the resolved URL list retains the same order as the input gallery list

### Requirement: Upload progress during save
While save-time uploads are running, the admin UI SHALL show progress indicating completed uploads over total uploads needed (e.g.「上传中 3/12」).

#### Scenario: Progress updates
- **WHEN** pending dual uploads are in progress during save
- **THEN** the UI shows a fraction of finished uploads versus total pending uploads

### Requirement: Fail save on any upload error
If any pending dual upload fails, the save operation SHALL fail and MUST NOT treat the gallery as successfully updated for that submit.

#### Scenario: One upload fails
- **WHEN** one of several parallel uploads returns an error
- **THEN** the save surfaces an error
- **AND** does not complete as a successful save

### Requirement: Streaming upload with OSS progress events
The system SHALL provide an authenticated admin upload endpoint that accepts the existing dual-image multipart fields (`original`, `thumb`, `category`), uploads both objects to OSS, and **streams progress events** to the client while OSS puts are in progress—not only after both puts complete.

#### Scenario: Progress events before done
- **WHEN** an admin client uploads a dual image pair via the streaming endpoint
- **THEN** the client receives one or more progress events with an overall percent before the final success event
- **AND** the final success event includes CDN URLs for original and thumb

#### Scenario: Error event on failure
- **WHEN** OSS upload fails after the multipart body was accepted
- **THEN** the stream ends with an error event containing a message
- **AND** the client MUST NOT treat the upload as successful

### Requirement: Progress reflects OSS body read
Progress percent SHALL be derived primarily from bytes consumed while uploading object bodies to OSS (original and thumb), weighted by their known sizes, so that progress continues to advance after the browser has finished sending the HTTP request body.

#### Scenario: Percent advances during OSS put
- **WHEN** the server is putting the original object to OSS
- **THEN** progress percent increases as the original body is read for upload
- **AND** thumb upload continues to advance overall percent until completion

### Requirement: Admin UI shows stream progress
The admin eager-upload pipeline SHALL use the streaming endpoint and display the stream’s overall percent on the slot status (e.g.「上传中 N%」) until success or failure. Cancel/abort SHALL abort the in-flight stream request.

#### Scenario: Badge tracks OSS-weighted percent
- **WHEN** the stream reports progress pct=42
- **THEN** the gallery/cover status shows uploading progress at 42%

#### Scenario: Done clears progress
- **WHEN** the stream reports done with URLs
- **THEN** the slot stores CDN URLs, status becomes uploaded, and progress is cleared

### Requirement: Compress outside upload queue
Gallery (and cover when using the shared pipeline) image compression SHALL NOT occupy a slot in the shared upload concurrency pool. Compression SHALL run in a separate limited pool; only the HTTP upload of already-compressed pending blobs SHALL use the upload queue (maximum concurrency 3).

#### Scenario: Later images compress while earlier upload
- **WHEN** more than three images are selected and the first uploads are in progress
- **THEN** additional images MAY begin or complete compression without waiting for an upload pool slot
- **AND** at most three HTTP uploads run concurrently

#### Scenario: Upload queue receives pending blobs only
- **WHEN** a slot finishes compression successfully
- **THEN** the client enqueues an upload of the compressed original+thumb pair
- **AND** the upload job does not re-run full dual compression

### Requirement: Parallel compress with concurrency limit
The client SHALL compress multiple pending files with a maximum concurrency of 3 (same order of magnitude as upload). Slots waiting for a compress slot SHALL keep their local preview and a waiting/processing status as appropriate.

#### Scenario: Compress pool capped
- **WHEN** more than three slots need compression at once
- **THEN** at most three compressions run simultaneously
- **AND** remaining slots wait without blocking UI display of their previews

### Requirement: Cancel aborts compress and upload
Deleting or clearing a slot SHALL cancel that slot’s compress work and upload work. A late compress or upload completion for a cancelled or removed `slotId` MUST NOT restore the slot into the list.

#### Scenario: Delete during compress
- **WHEN** an admin deletes a slot that is compressing
- **THEN** compression for that slot is cancelled or its result discarded
- **AND** no upload is started for that slotId
