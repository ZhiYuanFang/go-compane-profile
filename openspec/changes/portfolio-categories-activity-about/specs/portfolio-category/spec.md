## ADDED Requirements

### Requirement: Portfolio has a category
The system SHALL store each portfolio with exactly one category among `residential`, `commercial`, `office`, and `installation` (displayed as 住宅、商业、办公、装置).

#### Scenario: Create portfolio in a category
- **WHEN** an admin creates a portfolio from a category-scoped entry
- **THEN** the portfolio is persisted with that category
- **AND** its `sort_order` is assigned as the next value within that category

#### Scenario: Change portfolio category
- **WHEN** an admin updates a portfolio’s category to a different valid category
- **THEN** the portfolio moves to the new category
- **AND** it receives `sort_order` at the end of the new category
- **AND** remaining portfolios in the old category keep a contiguous class-local order

### Requirement: Public portfolio list accepts optional category
The public portfolio list API SHALL accept an optional `category` query parameter. When omitted, the API SHALL return portfolios of all categories. When provided, the API SHALL return only portfolios of that category. Invalid category values SHALL result in a client error.

#### Scenario: List all portfolios
- **WHEN** a client calls `GET /api/v1/portfolios` without `category`
- **THEN** the response includes portfolios from every category
- **AND** items are ordered by `sort_order` ascending then `id` ascending

#### Scenario: List by category
- **WHEN** a client calls `GET /api/v1/portfolios?category=residential`
- **THEN** every returned item has category `residential`
- **AND** items are ordered by `sort_order` ascending then `id` ascending

#### Scenario: Invalid category
- **WHEN** a client passes an unsupported `category` value
- **THEN** the API responds with a 4xx error

### Requirement: Sort order is maintained per category in admin
Admin list and reorder operations SHALL be scoped to a single category. Reordering SHALL rewrite `sort_order` only for portfolios in that category.

#### Scenario: Reorder within category
- **WHEN** an admin reorders portfolios for `commercial`
- **THEN** only `commercial` portfolios receive updated `sort_order` values matching the submitted order
- **AND** portfolios in other categories are unchanged

### Requirement: Admin provides four category entry points
The admin console SHALL expose four navigation entries for 住宅、商业、办公、装置 portfolios, each opening a list/edit flow fixed to that category.

#### Scenario: Navigate to residential portfolios
- **WHEN** an admin opens the 住宅作品集 entry
- **THEN** the list shows only residential portfolios
- **AND** create/edit actions keep or default to the residential category
