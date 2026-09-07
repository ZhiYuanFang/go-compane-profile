## ADDED Requirements

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
