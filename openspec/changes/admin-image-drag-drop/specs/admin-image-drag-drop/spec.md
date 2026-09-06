## ADDED Requirements

### Requirement: Drag-and-drop image selection on dual-image fields
The admin dual-image field SHALL accept an image file dropped onto the field surface. The dropped file SHALL enter the same local dual-image processing pipeline used by the file picker (pending local pair until form submit). Non-image drops SHALL be rejected with an error message and SHALL NOT change the current value.

#### Scenario: Drop image onto empty field
- **WHEN** an operator drops a valid image file onto an empty dual-image field
- **THEN** the field processes that image and shows a local preview without requiring the system file dialog

#### Scenario: Reject non-image drop
- **WHEN** an operator drops a non-image file onto a dual-image field
- **THEN** the field shows an error and keeps the previous value unchanged

### Requirement: Confirm before replacing existing preview via drop
When the dual-image field already has a preview (saved or pending) and the operator drops a new image onto it, the system SHALL ask for confirmation before replacing. If the operator cancels, the current value SHALL remain unchanged.

#### Scenario: Confirm replace on drop
- **WHEN** an operator drops an image onto a field that already shows a preview and confirms the replace prompt
- **THEN** the field processes the new image and replaces the previous pending or displayed selection

#### Scenario: Cancel replace on drop
- **WHEN** an operator drops an image onto a field that already shows a preview and cancels the replace prompt
- **THEN** the field keeps the previous preview and value

### Requirement: Empty preview click opens file picker
When the dual-image preview is empty, activating the preview (click) SHALL open the same file picker used by the choose-image control. When the preview is non-empty, activating the preview SHALL open the existing enlarge/lightbox behavior and SHALL NOT open the file picker.

#### Scenario: Click empty preview to choose file
- **WHEN** an operator clicks an empty dual-image preview
- **THEN** the system file picker opens for image selection

#### Scenario: Click filled preview enlarges
- **WHEN** an operator clicks a non-empty dual-image preview
- **THEN** the enlarged lightbox opens and the file picker does not open

### Requirement: File button remains available
The dual-image field SHALL continue to provide an explicit control to choose or re-choose an image via the system file picker, in addition to drag-and-drop and empty-preview click.

#### Scenario: Button still works
- **WHEN** an operator uses the choose/re-choose button
- **THEN** the system file picker opens as before
