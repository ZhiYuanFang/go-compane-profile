## ADDED Requirements

### Requirement: Vue 3 admin console with glassmorphism gray theme
The system SHALL provide a Vue 3 admin SPA with a glassmorphism visual language and high-end gray theme (deep gray atmosphere, translucent panels, cool silver accents). The SPA SHALL be served from the application host under `/admin/`.

#### Scenario: Admin entry
- **WHEN** an operator opens `/admin/` on the public host
- **THEN** the Vue 3 admin application loads and presents a login gate before content management

### Requirement: Content management screens
The admin console SHALL provide screens to edit company profile, manage ordered portfolios (including cover and render/real galleries), and manage the pricing dual image.

#### Scenario: Navigate modules
- **WHEN** an authenticated operator uses the admin navigation
- **THEN** they can reach company, portfolio, and pricing management views

### Requirement: Local dual-image UX on image fields
For every image field (logos, covers, gallery items, pricing), the console SHALL show local preview of the processed pair, allow re-selection before submit, and upload pairs only on form submit together with the rest of the save operation (or explicitly as part of the save flow such that unused drafts are not persisted to OSS).

#### Scenario: Preview before submit
- **WHEN** the operator selects an image for a field
- **THEN** the UI shows a local preview without requiring a prior OSS upload
