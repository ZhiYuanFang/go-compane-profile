## ADDED Requirements

### Requirement: Public portfolio list pagination
The system SHALL support offset pagination on unauthenticated `GET /api/v1/portfolios` via query parameters `page` and `pageSize`. The response SHALL include `list` (current page items), `total` (total matching rows), `page` (effective page), and `pageSize` (effective page size). Items SHALL remain ordered by ascending `sort_order` then ascending `id`. Each item's list `cover` SHALL remain the thumbnail CDN URL. When `page` or `pageSize` are omitted, the system SHALL default to `page=1` and `pageSize=10` (not the full collection).

#### Scenario: First page with defaults
- **WHEN** a client calls `GET /api/v1/portfolios` without query parameters
- **THEN** the response contains at most 10 items in `list`, `page` is 1, `pageSize` is 10, and `total` equals the full portfolio count

#### Scenario: Explicit page request
- **WHEN** a client calls `GET /api/v1/portfolios?page=2&pageSize=10` and at least 11 portfolios exist
- **THEN** `list` contains the second page of items in global sort order, `page` is 2, `pageSize` is 10, and `total` is the full count

#### Scenario: Empty portfolio list
- **WHEN** no portfolios exist
- **THEN** the endpoint returns `list` as an empty array, `total` as 0, and successful status

#### Scenario: Page beyond last
- **WHEN** a client requests a `page` greater than the last page for the given `pageSize`
- **THEN** `list` is empty, and `total` still reflects the full count

### Requirement: Pagination parameter clamping
The system SHALL treat `page` values less than 1 as page 1. The system SHALL clamp `pageSize` to a positive value with a maximum of 50 (values above 50 become 50).

#### Scenario: Invalid page and oversized pageSize
- **WHEN** a client calls `GET /api/v1/portfolios?page=0&pageSize=999`
- **THEN** the effective `page` is 1 and the effective `pageSize` is 50 in the response metadata and query window
