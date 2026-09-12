## ADDED Requirements

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
