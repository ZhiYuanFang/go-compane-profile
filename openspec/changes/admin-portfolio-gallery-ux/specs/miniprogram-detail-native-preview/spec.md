## ADDED Requirements

### Requirement: Native preview for non-cover detail images
On the mini-program portfolio detail page, activating a non-cover gallery main image SHALL open WeChat native image preview (`wx.previewImage`) so the user can view fullscreen with platform gesture zoom. The cover intro image SHALL NOT open native preview via this behavior.

#### Scenario: Tap gallery image opens native preview
- **WHEN** the user taps the main image on a non-cover detail page
- **THEN** native image preview opens for that image with gesture zoom available

#### Scenario: Cover does not open native preview
- **WHEN** the user views the cover intro page
- **THEN** tapping the cover does not open native image preview through the gallery preview handler

### Requirement: Preview URL set excludes cover and empty reals
The native preview URL list SHALL include non-cover gallery images for the portfolio (render originals preferred, plus non-empty real originals) in a stable browse order, and SHALL set `current` to the activated image. Empty real placeholders SHALL be omitted from the URL list. Cover SHALL be omitted.

#### Scenario: Swipe across gallery in native preview
- **WHEN** native preview opens from a gallery image and multiple non-empty gallery URLs exist
- **THEN** the user can swipe among those gallery URLs inside native preview

### Requirement: Per-pair real presence without fallback to first real
When building detail image pages from `renders` and `reals`, the client SHALL pair by index. An empty or missing real at an index SHALL mean that page has no real, and the client MUST NOT fall back to another index's real (including `reals[0]`). The render/real toggle control SHALL appear only when that page has a non-empty real.

#### Scenario: Middle pair without real
- **WHEN** `reals[1]` is empty and `reals[0]` and `reals[2]` are non-empty
- **THEN** the second image page shows no real toggle and does not display `reals[0]` as its real

### Requirement: Page-local pinch scale may be removed
The detail page MAY remove in-page `movable-view` pinch-zoom and swiper zoom-lock logic once native preview provides zoom, provided vertical paging among cover and gallery pages remains usable.

#### Scenario: Vertical paging still works after removing page zoom
- **WHEN** page-local scale handling has been removed in favor of native preview
- **THEN** the user can still vertically page between the cover and gallery items
