## ADDED Requirements

### Requirement: Contain fit for dual-image preview
The admin dual-image field SHALL display the selected or saved image entirely within the existing preview box using contain-fit behavior (no cropping). The preview box dimensions SHALL remain unchanged from the current layout size.

#### Scenario: Selected image fits inside preview
- **WHEN** an operator selects or views an image in a dual-image field
- **THEN** the full image is visible inside the preview box without cropping, even if letterboxing appears inside the box

### Requirement: Click to enlarge dual-image preview
The admin dual-image field SHALL allow the operator to open an enlarged view of the current image by activating the preview (for example, clicking the preview). The enlarged view SHALL prefer the original URL when available, otherwise the current preview source (including a local pending preview). The operator SHALL be able to dismiss the enlarged view.

#### Scenario: Open and dismiss enlarge
- **WHEN** an operator activates a non-empty dual-image preview
- **THEN** an enlarged view of that image is shown and can be dismissed to return to the form

#### Scenario: Empty preview not enlargeable
- **WHEN** the dual-image preview has no image
- **THEN** activating the preview does not open an enlarged view

### Requirement: Enlarge applies to all dual-image fields
Every admin screen that uses the dual-image field (company logos, portfolio cover, pair gallery items, pricing) SHALL provide the same click-to-enlarge behavior. Portfolio list cover thumbnails are out of scope for enlarge.

#### Scenario: Pricing image enlarge
- **WHEN** an operator views a selected pricing image preview and activates it
- **THEN** the enlarged view opens for that image
