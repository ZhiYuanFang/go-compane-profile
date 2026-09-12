## ADDED Requirements

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
