## ADDED Requirements

### Requirement: Activity is a deletable singleton with cover and body images
The system SHALL store at most one activity consisting of one cover image pair and one body image pair (original + thumb, following existing dual-image upload conventions). Admins SHALL be able to set or clear either image. When both cover and body are cleared (or no activity row is present), the activity SHALL be treated as empty.

#### Scenario: Save activity images
- **WHEN** an admin uploads a cover and a body image for the activity
- **THEN** the activity stores both image pairs
- **AND** subsequent reads return those images

#### Scenario: Delete activity
- **WHEN** an admin clears/deletes the activity images
- **THEN** the public activity read indicates empty/absent activity

### Requirement: Public activity API
The system SHALL expose `GET /api/v1/activity` that returns the current activity cover and body image URLs when present, or an empty/null payload when absent.

#### Scenario: Activity present
- **WHEN** a client calls `GET /api/v1/activity` and images are configured
- **THEN** the response includes cover and body image references usable by the mini-program

#### Scenario: Activity absent
- **WHEN** a client calls `GET /api/v1/activity` and no activity is configured
- **THEN** the response indicates no activity (null or empty fields)
- **AND** does not return a server error

### Requirement: Admin activity management UI
The admin console SHALL provide an activity settings page to upload/replace cover and body images and to delete the activity.

#### Scenario: Admin opens activity page
- **WHEN** an admin navigates to the activity entry
- **THEN** they can view current images (if any), upload cover and body, and delete

### Requirement: Mini-program activity banner and detail
When the home page receives a non-empty activity, it SHALL show the cover above the category tabs in a container of height 80rpx with the image filling via aspect fill (外切). Tapping the cover SHALL open a page that shows the body image with `aspectFit` inside a vertically scrollable area. When activity is empty, the home page SHALL not show the banner.

#### Scenario: Show and open activity
- **WHEN** home loads and activity has a cover
- **THEN** a 80rpx-tall cover banner appears above the tabs
- **AND** tapping it opens the activity body page with a scrollable `aspectFit` image

#### Scenario: Hide when empty
- **WHEN** home loads and activity is empty
- **THEN** no activity banner is shown above the tabs
