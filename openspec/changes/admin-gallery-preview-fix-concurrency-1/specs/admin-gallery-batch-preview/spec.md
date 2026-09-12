## ADDED Requirements

### Requirement: Batch select keeps all local previews
When an admin batch-selects multiple images, the client SHALL associate a local preview with every new slot in the same update cycle (or equivalent) such that no selected slot loses its preview due to concurrent status patches on other slots.

#### Scenario: Ten files all show previews immediately
- **WHEN** an admin batch-selects ten images
- **THEN** all ten gallery cells show a local preview without waiting for compression or upload to finish
- **AND** subsequent status updates on one slot MUST NOT clear another slot’s local preview

### Requirement: Async patches do not clobber sibling slots
Per-slot pipeline patches (queued / compressing / uploading / done / error) SHALL update only the targeted `slotId` while preserving other slots’ fields including `localPreview`, `pending`, and remote URLs.

#### Scenario: Compress patch preserves neighbors
- **WHEN** slot A’s status is patched to compressing while slot B already has a localPreview
- **THEN** slot B’s localPreview remains visible
