## ADDED Requirements

### Requirement: Manage company profile
Authenticated admins SHALL be able to read and update company copy fields and logo URL pairs (square and horizontal).

#### Scenario: Update company copy
- **WHEN** an authenticated admin submits updated intro title, body, values, years label, and address
- **THEN** subsequent public company API responses reflect the new values

#### Scenario: Update logos
- **WHEN** an authenticated admin sets square and/or horizontal logo thumb and original URLs
- **THEN** subsequent public company API responses return those URLs

### Requirement: Manage portfolios with sort order
Authenticated admins SHALL be able to create, update, delete, and reorder portfolios. List order for the public API MUST follow the admin-defined sort order.

#### Scenario: Create portfolio
- **WHEN** an authenticated admin creates a portfolio with address, area, style, heartFlow, cover dual URLs, and optional render/real image lists
- **THEN** the portfolio appears in the public list at the configured sort position

#### Scenario: Reorder portfolios
- **WHEN** an authenticated admin changes portfolio sort order
- **THEN** `GET /api/v1/portfolios` returns items in the new order

#### Scenario: Delete portfolio
- **WHEN** an authenticated admin deletes a portfolio
- **THEN** the portfolio is no longer returned by public list or detail endpoints

### Requirement: Manage portfolio image galleries
Authenticated admins SHALL be able to manage ordered `renders` and `reals` image lists for a portfolio, each image stored as thumb + original URL pair.

#### Scenario: Update renders and reals
- **WHEN** an authenticated admin saves ordered render and real dual-URL lists for a portfolio
- **THEN** the public detail endpoint returns those arrays in the same order

### Requirement: Manage pricing image
Authenticated admins SHALL be able to set pricing sheet thumb and original URLs.

#### Scenario: Update pricing
- **WHEN** an authenticated admin sets pricing dual URLs
- **THEN** the public pricing endpoint returns those URLs

### Requirement: Seed company copy only
Initial database seed SHALL insert company copy text derived from the former mini-program mock and SHALL NOT insert mock portfolio items or pricing images.

#### Scenario: Fresh database seed
- **WHEN** migrations/seed run on an empty database
- **THEN** company copy fields are populated and portfolio count is zero
