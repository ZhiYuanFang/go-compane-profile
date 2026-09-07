## ADDED Requirements

### Requirement: Detail pages support render-only, real-only, and both
The mini-program portfolio detail gallery SHALL build one page per aligned pair index. A page MAY have only a render, only a real, or both. Pairs where both sides are empty SHALL NOT produce a gallery page.

#### Scenario: Real-only pair shows as a page
- **WHEN** `renders[i]` is empty and `reals[i]` is non-empty
- **THEN** the detail swiper includes a page for index `i`
- **AND** that page displays the real image by default

#### Scenario: Render-only pair shows as a page
- **WHEN** `renders[i]` is non-empty and `reals[i]` is empty
- **THEN** the page displays the render image
- **AND** no render/real toggle is shown

#### Scenario: Both present
- **WHEN** both `renders[i]` and `reals[i]` are non-empty
- **THEN** the page defaults to the render image
- **AND** a toggle control is available to switch to the real image

### Requirement: Toggle only when both sides exist
The render/real toggle control SHALL appear only when that page has both a non-empty render and a non-empty real.

#### Scenario: No toggle for single-sided page
- **WHEN** a page has only one non-empty side
- **THEN** the toggle control is not shown

### Requirement: Native preview skips empty sides
When opening native image preview from a gallery page, the URL list SHALL include only non-empty render and real originals for gallery pages, and `current` SHALL be the currently displayed side’s URL.

#### Scenario: Preview from real-only page
- **WHEN** the user opens native preview on a real-only page
- **THEN** preview opens with that real image as current
- **AND** empty render placeholders are not included in the URL list
