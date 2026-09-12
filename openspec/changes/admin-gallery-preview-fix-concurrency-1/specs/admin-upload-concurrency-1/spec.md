## ADDED Requirements

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
