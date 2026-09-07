## ADDED Requirements

### Requirement: Home uses category tabs and a viewpager
The mini-program home page SHALL present tabs for 全部、住宅、商业、办公、装置. The default selected tab SHALL be 全部. Content SHALL be shown in a horizontal pager (swiper/viewpager) synced with the selected tab. Selecting a tab SHALL enlarge and bold the selected tab label relative to unselected tabs.

#### Scenario: Default tab is all
- **WHEN** the user opens home
- **THEN** the 全部 tab is selected
- **AND** the pager shows the all-categories portfolio list

#### Scenario: Switch category via tab
- **WHEN** the user selects 办公
- **THEN** the 办公 label is visually emphasized (larger and bolder)
- **AND** the pager shows the office-category portfolio list

#### Scenario: Swipe pager updates tab
- **WHEN** the user swipes the pager to another category page
- **THEN** the corresponding tab becomes selected and emphasized

### Requirement: Each pager page loads portfolios for its filter
The 全部 page SHALL request portfolios without a category filter. Each category page SHALL request portfolios with the matching category parameter. Existing pagination / load-more behavior SHALL apply per page independently.

#### Scenario: Category page requests filter
- **WHEN** the residential pager page loads or refreshes
- **THEN** it requests the public portfolio list with `category=residential`

#### Scenario: All page requests unfiltered list
- **WHEN** the 全部 pager page loads or refreshes
- **THEN** it requests the public portfolio list without a category parameter

### Requirement: Activity banner sits above tabs
When an activity cover is available, the home layout SHALL place the activity banner above the category tabs (not between tabs and list).

#### Scenario: Banner above tabs
- **WHEN** activity cover is present
- **THEN** the banner is rendered above the tab row
