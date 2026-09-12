## ADDED Requirements

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
