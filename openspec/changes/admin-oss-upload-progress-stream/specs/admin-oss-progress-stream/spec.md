## ADDED Requirements

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
