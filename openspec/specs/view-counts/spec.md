# view-counts

## Purpose

view-counts 域行为规格（自变更 delta / 旧基线迁移）。

<!-- migrated-from: changes/content-view-counts/specs/view-count-admin, changes/content-view-counts/specs/view-count-api -->

## Requirements

### Requirement: Portfolio list shows heat
Admin portfolio list responses SHALL include each item's `viewCount`, and the admin portfolio list UI SHALL display it as「热度 {n}」(e.g. `热度 12`).

#### Scenario: List row heat label
- **WHEN** an admin opens a category portfolio list that has items with view counts
- **THEN** each row shows `热度` followed by that item's count

### Requirement: Activity list shows heat
Admin activity list responses SHALL include each item's `viewCount`, and the admin activity list UI SHALL display it as「热度 {n}」.

#### Scenario: Activity list row heat label
- **WHEN** an admin opens the activity list
- **THEN** each row shows `热度` followed by that item's count

### Requirement: Company page shows about-us views
The admin company profile payload SHALL include `aboutViewCount`, and the company edit page SHALL show a separate read-only module with the exact copy pattern「查阅关于我们次数为 {n}」.

#### Scenario: About views module
- **WHEN** an admin opens the company profile page
- **THEN** they see a dedicated module stating 查阅关于我们次数为 followed by the current count
- **AND** the count is not editable in that module

### Requirement: Simplified Chinese admin copy
Admin-facing labels for this feature SHALL use Simplified Chinese as specified (`热度`, `查阅关于我们次数为`).

#### Scenario: Copy language
- **WHEN** an admin views portfolio list, activity list, or company about-views module
- **THEN** the heat/about labels use the specified Simplified Chinese phrases

### Requirement: Increment portfolio view count
The system SHALL expose `POST /api/v1/portfolios/{id}/view` that atomically increments `view_count` for an existing portfolio identified by numeric id or slug (same resolution rules as portfolio detail GET).

#### Scenario: Successful increment
- **WHEN** a client POSTs to `/api/v1/portfolios/{id}/view` for an existing portfolio
- **THEN** that portfolio's `view_count` increases by 1
- **AND** the response indicates success

#### Scenario: Missing portfolio
- **WHEN** a client POSTs for a non-existent portfolio id/slug
- **THEN** the system returns not found and does not change any counts

### Requirement: Increment activity view count
The system SHALL expose `POST /api/v1/activities/{id}/view` that atomically increments `view_count` for an existing activity by numeric id.

#### Scenario: Successful increment
- **WHEN** a client POSTs to `/api/v1/activities/{id}/view` for an existing activity
- **THEN** that activity's `view_count` increases by 1
- **AND** the response indicates success

#### Scenario: Missing activity
- **WHEN** a client POSTs for a non-existent activity id
- **THEN** the system returns not found and does not change any counts

### Requirement: Increment about-us view count
The system SHALL expose `POST /api/v1/company/about/view` that atomically increments `about_view_count` on the company singleton row.

#### Scenario: Successful about increment
- **WHEN** a client POSTs to `/api/v1/company/about/view`
- **THEN** the company `about_view_count` increases by 1
- **AND** the response indicates success

### Requirement: Counts not tied to content GET
Content GET endpoints for portfolio detail, activity detail, and company SHALL NOT increment view counters as a side effect.

#### Scenario: GET detail does not count
- **WHEN** a client GETs portfolio detail, activity detail, or company profile
- **THEN** the corresponding view counter is unchanged by that request alone
