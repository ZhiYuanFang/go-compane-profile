## ADDED Requirements

### Requirement: Record portfolio open onLoad
When the portfolio detail page loads (`onLoad`), the mini-program SHALL call `POST /api/v1/portfolios/{id}/view` once for that page load, using the opened portfolio id. Tracking failures MUST NOT block rendering of the detail content.

#### Scenario: Portfolio detail onLoad tracks
- **WHEN** the user opens a portfolio detail page and `onLoad` runs with a valid id
- **THEN** the mini-program issues one view-increment request for that id
- **AND** detail content still loads if the increment request fails

### Requirement: Record activity open onLoad
When the activity detail page loads (`onLoad`), the mini-program SHALL call `POST /api/v1/activities/{id}/view` once for that page load. Tracking failures MUST NOT block rendering.

#### Scenario: Activity detail onLoad tracks
- **WHEN** the user opens an activity detail page and `onLoad` runs with a valid id
- **THEN** the mini-program issues one view-increment request for that id
- **AND** activity content still loads if the increment request fails

### Requirement: Record about-us open onLoad
When the about-us page loads (`onLoad`), the mini-program SHALL call `POST /api/v1/company/about/view` once for that page load. Tracking failures MUST NOT block rendering.

#### Scenario: About page onLoad tracks
- **WHEN** the user opens the about-us page and `onLoad` runs
- **THEN** the mini-program issues one about-view increment request
- **AND** about content still loads if the increment request fails

### Requirement: No mini-program heat UI
The mini-program SHALL NOT display portfolio, activity, or about view counts in end-user UI as part of this change.

#### Scenario: End user does not see counts
- **WHEN** a user browses detail or about pages
- **THEN** no heat/view-count label is shown on those pages for this feature
