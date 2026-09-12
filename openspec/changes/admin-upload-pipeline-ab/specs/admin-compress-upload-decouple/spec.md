## ADDED Requirements

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
