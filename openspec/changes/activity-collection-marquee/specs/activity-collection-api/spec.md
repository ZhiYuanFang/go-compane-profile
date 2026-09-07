## ADDED Requirements

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
**Reason**: Replaced by the activity collection APIs.  
**Migration**: Clients must use `GET /api/v1/activities` and `GET /api/v1/activities/{id}`; operators re-enter data after schema cutover (no automated migration of old singleton cover/body rows).

#### Scenario: Old singleton path unavailable
- **WHEN** a client calls the former singleton activity endpoint after this change
- **THEN** the endpoint is unavailable (removed or not found)
