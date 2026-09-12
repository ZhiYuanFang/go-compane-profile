## ADDED Requirements

### Requirement: Four status tone tiers
The admin UI SHALL classify each image slot status into one of four visual tones: wait (排队中), busy (处理中 or 上传中), done (已上传), or error (失败). Compressing and uploading SHALL share the busy tone.

#### Scenario: Queued is wait
- **WHEN** a slot status is queued
- **THEN** its tone is wait

#### Scenario: Compressing and uploading are busy
- **WHEN** a slot status is compressing or uploading
- **THEN** its tone is busy

#### Scenario: Done and error tones
- **WHEN** a slot status is done
- **THEN** its tone is done
- **WHEN** a slot status is error
- **THEN** its tone is error

### Requirement: Status text color by tone
Gallery status badges and cover DualImageField status text SHALL use distinct colors per tone so wait, busy, done, and error are visually distinguishable at a glance.

#### Scenario: Badge colors differ by tone
- **WHEN** several gallery slots show wait, busy, done, and error statuses
- **THEN** their status label colors differ according to those four tones

### Requirement: Border color by tone
Gallery cells and the cover dual container SHALL apply a border (outline) color matching the same tone as the status text when a tone is present.

#### Scenario: Cell border matches status tone
- **WHEN** a gallery slot is uploading (busy)
- **THEN** the cell border uses the busy tone color
- **AND** the status text also uses the busy tone color
