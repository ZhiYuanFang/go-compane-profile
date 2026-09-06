## ADDED Requirements

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
