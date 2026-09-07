## ADDED Requirements

### Requirement: Admin activity list with reorder
The admin console SHALL provide an activity list showing all activities with controls to move items up/down (or equivalent order editing), open edit, create new, and delete.

#### Scenario: Open activity list
- **WHEN** an admin opens the activity management entry
- **THEN** they see all activities in current sort order
- **AND** can reorder, edit, create, and delete

### Requirement: Admin activity editor
The create/edit form SHALL include title, a required dual-image field for the activity image, and a rich-text editor for body limited to font size, bold, and color. Saving without an image SHALL be blocked in the UI and/or by API error.

#### Scenario: Edit rich text body
- **WHEN** an admin applies bold, color, or font size in the body editor and saves
- **THEN** the stored body HTML reflects those styles

#### Scenario: Image required in form
- **WHEN** an admin tries to save with an empty image
- **THEN** save does not succeed as a valid activity

### Requirement: Simplified Chinese admin copy
Admin-facing labels and messages for this feature SHALL use Simplified Chinese.
