## ADDED Requirements

### Requirement: Hide portfolio slug in admin editor
The portfolio create and edit forms SHALL NOT present a slug input. On create, the client SHALL omit a custom slug (empty) so the server may auto-generate one. On edit, the client SHALL NOT require the operator to manage slug.

#### Scenario: Create without slug field
- **WHEN** an operator opens the new portfolio form
- **THEN** no slug field is shown and a successful create still persists a server-generated slug

### Requirement: Paired render and real gallery editor
The portfolio editor SHALL present gallery images as ordered pairs, each pair containing one render (effect) image field and one optional real (site photo) image field. Adding a pair SHALL add both slots together. Reordering or removing SHALL apply to the whole pair. Real MAY be empty.

#### Scenario: Add pair with optional real
- **WHEN** an operator adds a gallery pair and fills only the render image
- **THEN** the pair is valid to keep in the editor and the real slot may remain empty

#### Scenario: Reorder pairs
- **WHEN** an operator moves a pair up or down
- **THEN** both the render and real for that pair move together

### Requirement: Persist pairs as aligned render and real arrays
On save, the editor SHALL submit `renders` and `reals` arrays of equal length corresponding to pair order. An empty real SHALL be represented as an empty dual-URL placeholder at the same index so public clients can align by index.

#### Scenario: Save with a middle empty real
- **WHEN** an operator saves three pairs where the second real is empty
- **THEN** the saved `reals` list has length three and the second entry has empty thumb and original URLs

#### Scenario: Load unequal historical lists
- **WHEN** the editor loads a portfolio whose stored render and real list lengths differ
- **THEN** the UI zips them into pairs by index and pads the shorter side with empty dual slots so no images are dropped

### Requirement: Navigate to portfolio list after successful save
After a successful portfolio create or update, the admin console SHALL navigate the operator to the portfolio list page.

#### Scenario: Create returns to list
- **WHEN** an operator successfully creates a portfolio
- **THEN** the console navigates to the portfolio list

#### Scenario: Update returns to list
- **WHEN** an operator successfully updates a portfolio
- **THEN** the console navigates to the portfolio list
