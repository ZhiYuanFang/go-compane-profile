## ADDED Requirements

### Requirement: Pending images show local preview
When a gallery dual has `pending` blobs and no remote URLs yet, the admin UI SHALL show a local object-URL preview of the pending thumb (or equivalent) instead of a blank preview area.

#### Scenario: Batch-added pending preview
- **WHEN** batch select finishes processing a file into a pending dual
- **THEN** that slot's thumbnail/preview shows the local image
- **AND** status may still indicate pending upload until form submit

### Requirement: Progressive batch select with progress
Batch file selection SHALL process images and append each completed slot to the list progressively, showing progress such as processed/total while work is in progress. The UI MUST remain understandable during long batches (not a blank list until the entire batch completes).

#### Scenario: Progressive append
- **WHEN** an admin batch-selects multiple images
- **THEN** slots appear as each file finishes processing
- **AND** a progress indicator shows how many of the batch are done

#### Scenario: Partial failure
- **WHEN** one file in a batch fails processing after some succeeded
- **THEN** successfully processed slots remain in the list
- **AND** an error message explains the failure

### Requirement: Batch controls disabled while processing
While batch processing is running, batch select and destructive batch actions SHALL be disabled or otherwise blocked to avoid conflicting edits.

#### Scenario: Processing lock
- **WHEN** batch processing is in progress
- **THEN** starting another batch select is disabled until processing finishes
