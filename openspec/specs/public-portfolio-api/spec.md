# public-portfolio-api

## Purpose

public-portfolio-api 域行为规格（自变更 delta / 旧基线迁移）。

<!-- migrated-from: changes/company-profile-api-admin/specs/public-portfolio-api, changes/portfolio-categories-activity-about/specs/company-about-contact, changes/portfolio-public-pagination/specs/public-portfolio-pagination -->

## Requirements

### Requirement: Public company profile endpoint
The system SHALL expose an unauthenticated `GET` endpoint that returns the company profile fields needed by the mini-program welcome and branding UI: intro title, intro body, values, years label, address, and logo pairs (square and horizontal), each as thumb and original CDN URLs (URLs MAY be empty when not yet uploaded).

#### Scenario: Company profile returned
- **WHEN** a client calls `GET /api/v1/company`
- **THEN** the response includes `introTitle`, `introBody`, `values`, `yearsLabel`, `address`, and logo fields with `thumb` and `original` URLs for both square and horizontal logos

### Requirement: Public portfolio list uses thumb covers
The system SHALL expose an unauthenticated portfolio list ordered by server-defined sort order (no client-side re-sort required). Each item's list cover SHALL be the thumbnail CDN URL.

#### Scenario: Ordered list with thumb cover
- **WHEN** a client calls `GET /api/v1/portfolios`
- **THEN** items are returned in ascending sort order and each item's `cover` value is the thumbnail URL

#### Scenario: Empty portfolio list
- **WHEN** no portfolios exist
- **THEN** the endpoint returns an empty list successfully

### Requirement: Public portfolio detail with dual image URLs
The system SHALL expose an unauthenticated detail endpoint by portfolio id. Cover, each render, and each real image SHALL be represented as `{ thumb, original }` CDN URL pairs. Render and real arrays SHALL preserve server sort order.

#### Scenario: Detail found
- **WHEN** a client requests an existing portfolio id
- **THEN** the response includes project fields (`address`, `area`, `style`, `heartFlow`) and dual-URL image structures for cover, `renders`, and `reals`

#### Scenario: Detail missing
- **WHEN** a client requests an unknown portfolio id
- **THEN** the API returns a not-found error suitable for the mini-program「作品不存在」handling

### Requirement: Public pricing dual image
The system SHALL expose an unauthenticated pricing endpoint returning thumb and original CDN URLs for the fee sheet image (URLs MAY be empty when not uploaded).

#### Scenario: Pricing returned
- **WHEN** a client calls `GET /api/v1/pricing`
- **THEN** the response includes `thumb` and `original` URLs

### Requirement: Company profile includes awards and contact fields
The company profile SHALL include `awards` (multiline text), `phone`, and `wechat` in addition to existing fields. `address` SHALL continue to be the company address field used by the about page. Public `GET /api/v1/company` and admin company read/update SHALL include these fields. Newlines in `awards` SHALL be preserved in storage and API responses.

#### Scenario: Admin updates about fields
- **WHEN** an admin saves awards, phone, and wechat on the company profile page
- **THEN** subsequent public company reads return those values
- **AND** awards retains newline characters as entered

#### Scenario: Address reuse
- **WHEN** the about page displays address
- **THEN** it uses the existing `company.address` value

### Requirement: Pricing capability is removed
The system SHALL remove pricing end-to-end: database table usage, public and admin pricing APIs, admin pricing UI, and mini-program pricing page and entry. Clients MUST NOT rely on `GET /api/v1/pricing`.

#### Scenario: Pricing endpoints gone
- **WHEN** a client calls the former public or admin pricing endpoints after this change
- **THEN** the endpoints are unavailable (not found or removed from routing)

#### Scenario: Admin nav has no pricing
- **WHEN** an admin views the sidebar
- **THEN** there is no 资费 entry

### Requirement: Public portfolio list pagination
The system SHALL support offset pagination on unauthenticated `GET /api/v1/portfolios` via `page` and `pageSize`, returning `list`, `total`, `page`, and `pageSize`. Items SHALL stay ordered by ascending `sort_order` then `id`, and list `cover` SHALL remain the thumbnail CDN URL.

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
