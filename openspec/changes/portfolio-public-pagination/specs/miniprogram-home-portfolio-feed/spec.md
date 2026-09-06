## ADDED Requirements

### Requirement: Home feed loads portfolios in pages of ten
The mini-program home page SHALL load public portfolios in pages of 10 via the paginated public list API, accumulate results across pages, and render the dual-column layout from the accumulated list using global index parity (even → left, odd → right).

#### Scenario: Initial load
- **WHEN** the user opens the home page
- **THEN** the app requests page 1 with pageSize 10 and displays returned items split into left and right columns

#### Scenario: Append preserves column parity
- **WHEN** a later page is appended to the accumulated list
- **THEN** columns are recomputed from the full accumulated list so global even/odd indices remain consistent

### Requirement: Load more on scroll to bottom with status text
When the user scrolls the home feed to the bottom and more items remain (`accumulated length < total`), the app SHALL request the next page and show a「加载中」status. While a load-more or refresh request is in flight, the app SHALL ignore duplicate load-more triggers. When no more items remain after a successful load (or `total` is already reached), the app SHALL show「没有更多」status text (or equivalent「没有更多了」copy).

#### Scenario: Reach bottom with more data
- **WHEN** the user scrolls to the bottom and `hasMore` is true and no request is in flight
- **THEN** the UI shows「加载中」, fetches the next page, appends items, and updates columns

#### Scenario: Reach bottom when exhausted
- **WHEN** the user scrolls to the bottom and all items for `total` are already loaded
- **THEN** the UI shows「没有更多」and does not issue another list request

#### Scenario: Concurrent scroll ignored
- **WHEN** a load-more request is already in flight and the user hits the bottom again
- **THEN** the app does not start a duplicate request

### Requirement: Pull-down refresh resets to first page
The home feed SHALL support pull-down refresh. On refresh, the app SHALL reload page 1 (pageSize 10), replace the accumulated list (not append), re-split columns, and end the refresher state. Given the home layout uses a vertical `scroll-view`, pull-down refresh MAY be implemented with `scroll-view` refresher events; behavior SHALL match user-visible pull-to-refresh reset semantics.

#### Scenario: Pull to refresh
- **WHEN** the user pulls down to refresh the home feed
- **THEN** the app fetches page 1, replaces prior list data with the new page, resets pagination state so subsequent load-more starts from page 2 when more exists, and completes the refresh indicator
