# activity

## Purpose

activity 域行为规格（自变更 delta / 旧基线迁移）。

<!-- migrated-from: changes/activity-collection-marquee/specs/activity-admin-cms, changes/activity-collection-marquee/specs/activity-collection-api, changes/portfolio-categories-activity-about/specs/activity -->

## Requirements

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

#### Scenario: 默认验收

- **WHEN** 该需求适用的条件成立
- **THEN** 系统行为符合本 Requirement

### Requirement: Activity is a multi-record collection
The system SHALL store zero or more activities. Each activity SHALL have a title, a required dual image (thumb + original), a body stored as HTML, and a sort order. Cover-image and body-image singleton fields SHALL NOT be used.

#### Scenario: Create activity with required image
- **WHEN** an admin creates an activity with title, non-empty image dual, and optional HTML body
- **THEN** the activity is persisted with those fields and a sort order at the end of the list

#### Scenario: Reject missing image
- **WHEN** an admin attempts to create or update an activity without an image
- **THEN** the API rejects the request with a client error

### Requirement: Public list returns all activities in sort order
The system SHALL expose `GET /api/v1/activities` that returns all activities ordered by `sort_order` ascending then `id` ascending, without pagination.

#### Scenario: Full ordered list
- **WHEN** a client requests the public activities list
- **THEN** every activity is included
- **AND** items appear in admin-defined sort order

#### Scenario: Empty list
- **WHEN** no activities exist
- **THEN** the response is an empty list without server error

### Requirement: Public activity detail
The system SHALL expose `GET /api/v1/activities/{id}` returning title, image dual URLs, and body HTML for one activity.

#### Scenario: Get by id
- **WHEN** a client requests a valid activity id
- **THEN** the response includes title, image, and bodyHtml

#### Scenario: Missing activity
- **WHEN** the id does not exist
- **THEN** the API returns not found

### Requirement: Admin CRUD and reorder
Admin APIs SHALL support list, get, create, update, delete, and reorder (ordered id list rewriting `sort_order`). Reorder SHALL only affect activities.

#### Scenario: Reorder
- **WHEN** an admin submits an ordered list of activity ids
- **THEN** `sort_order` values match that order for subsequent public and admin lists

### Requirement: Body HTML is a sanitized subset
Activity body HTML SHALL allow only a safe subset suitable for font size, bold, and color (and structural line breaks/paragraphs). Disallowed tags and attributes SHALL be stripped or rejected on save.

#### Scenario: Styled paragraph saved
- **WHEN** an admin saves body HTML using allowed bold/color/font-size markup
- **THEN** subsequent reads return that content for mini-program rich-text display

### Requirement: Singleton activity API removed
**Reason**: Replaced by the activity collection APIs. The system SHALL satisfy this requirement.
**Migration**: Clients must use `GET /api/v1/activities` and `GET /api/v1/activities/{id}`; operators re-enter data after schema cutover (no automated migration of old singleton cover/body rows).

#### Scenario: Old singleton path unavailable
- **WHEN** a client calls the former singleton activity endpoint after this change
- **THEN** the endpoint is unavailable (removed or not found)

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
