## ADDED Requirements

### Requirement: Replace local mock with API data modules
The WeChat mini-program SHALL obtain company, portfolio list/detail, and pricing data from `https://muhou.cuplay.top` public APIs. Local mock modules `data/company.js` and `data/portfolio.js` SHALL be removed or no longer used as the data source.

#### Scenario: Welcome uses API company
- **WHEN** the welcome page loads
- **THEN** it renders company copy and logos from the public company API

#### Scenario: Home uses API portfolio list
- **WHEN** the home page loads
- **THEN** it renders the waterfall from the public portfolio list API in API order

### Requirement: List cover uses thumbnail
Portfolio cards on the home page SHALL use the thumbnail cover URL from the API.

#### Scenario: Card cover is thumb
- **WHEN** a portfolio card is rendered
- **THEN** its cover image src is the thumbnail URL

### Requirement: Detail progressive image swap
On the detail page, cover and gallery images SHALL first display the thumbnail URL, then silently replace with the original URL after the original has loaded successfully.

#### Scenario: Progressive swap
- **WHEN** a detail image begins loading
- **THEN** the UI shows the thumb first and switches to original after original load success without requiring user action

### Requirement: Pricing page uses API dual image
The pricing page SHALL load the fee sheet from the public pricing API using the same thumb-then-original progressive behavior.

#### Scenario: Pricing from API
- **WHEN** the pricing page opens
- **THEN** it displays the API pricing image ( Progressive thumb then original ) rather than a packaged local `cash.png` mock dependency as the source of truth
