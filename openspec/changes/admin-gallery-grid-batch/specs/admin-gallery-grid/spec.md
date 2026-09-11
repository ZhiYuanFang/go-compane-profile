## ADDED Requirements

### Requirement: Two-column thumbnail gallery grid
The admin portfolio gallery editor SHALL display gallery slots in a two-column grid of compact thumbnails (not full-width DualImageField rows as the primary list surface). On narrow viewports the grid MAY collapse to one column.

#### Scenario: Dense grid layout
- **WHEN** an admin opens a portfolio with multiple gallery images
- **THEN** thumbnails appear in a two-column grid
- **AND** each cell shows an index and a small preview when an image is present

### Requirement: Click thumbnail to edit
Clicking a gallery thumbnail SHALL open an edit surface (dialog or equivalent) that exposes DualImageField controls for that slot (replace, clear, full preview). Checkbox clicks SHALL NOT open the editor.

#### Scenario: Open slot editor
- **WHEN** an admin clicks a gallery thumbnail (not its checkbox)
- **THEN** an editor opens for that slot with dual-image controls
- **AND** closing the editor returns to the grid with any changes reflected

#### Scenario: Empty slot still editable
- **WHEN** an admin clicks an empty gallery thumbnail/placeholder
- **THEN** the editor opens so they can choose an image for that slot

### Requirement: Compact per-slot affordances
Each grid cell SHALL support selection via checkbox. Reorder (上移/下移) and single remove MAY live in the edit dialog and/or as compact cell actions; the primary list SHALL NOT require a full-width DualImageField per row.

#### Scenario: Cell shows selection control
- **WHEN** gallery items are shown in the grid
- **THEN** each cell has a checkbox for multi-select
