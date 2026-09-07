## ADDED Requirements

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
