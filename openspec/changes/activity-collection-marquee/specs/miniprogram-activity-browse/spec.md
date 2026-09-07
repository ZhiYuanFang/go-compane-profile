## ADDED Requirements

### Requirement: Home vertical marquee of all activity titles
When at least one activity exists, the mini-program home SHALL show a notice bar above the portfolio category tabs that vertically cycles through all activity titles in sort order. Dwell time SHALL be 1 second per title. The bar SHALL support both automatic playback and manual swipe. Each title SHALL display on one line with end ellipsis when overflowing. Tapping the bar SHALL open the detail page for the currently shown activity. When no activities exist, the bar SHALL be hidden entirely.

#### Scenario: Auto and manual cycle
- **WHEN** multiple activities exist on home
- **THEN** titles advance automatically about every 1 second
- **AND** the user can manually swipe between titles
- **AND** tapping opens the corresponding activity detail

#### Scenario: Hide when empty
- **WHEN** the activities list is empty
- **THEN** no activity notice bar is shown on home

### Requirement: Activity detail layout and interactions
The activity detail page SHALL lay out content in order: image, then bold title, then HTML body via rich-text. The image SHALL span the screen width with height adapting to aspect ratio. Tapping the image SHALL open native preview with zoom. A control labeled 「查看全部活动」 SHALL appear at the bottom-right of the image area (visual style consistent with the home portfolio address badge) and SHALL navigate to the all-activities list page. User-visible copy SHALL be Simplified Chinese.

#### Scenario: Detail order and preview
- **WHEN** the user opens an activity detail
- **THEN** they see image above bold title above rich-text body
- **AND** tapping the image opens zoomable preview
- **AND** tapping 「查看全部活动」 opens the activity list page

### Requirement: All-activities waterfall list
The mini-program SHALL provide an activity list page that shows all activities in sort order as a two-column waterfall of rounded cards. Each card SHALL show the activity image on top and the title below, with the title limited to at most two lines and ellipsis for overflow. Tapping a card SHALL open that activity’s detail page.

#### Scenario: Open list from detail
- **WHEN** the user enters the activity list page
- **THEN** all activities appear as two-column cards with image and up-to-two-line titles
- **AND** tapping a card opens the matching detail
