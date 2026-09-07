## ADDED Requirements

### Requirement: Either side of a gallery pair may be empty
In the admin portfolio pair gallery editor, each pair SHALL allow an empty render (effect) image and an empty real (site) image. A pair with at least one non-empty side SHALL be valid to save. Cover image requirements are unchanged by this capability and SHALL remain required as today.

#### Scenario: Save pair with only real image
- **WHEN** an operator fills only the real slot of a pair and leaves render empty
- **THEN** the pair is accepted on save
- **AND** the submitted `reals` entry at that index is non-empty
- **AND** the submitted `renders` entry at that index is an empty dual placeholder

#### Scenario: Save pair with only render image
- **WHEN** an operator fills only the render slot and leaves real empty
- **THEN** the pair is accepted on save with an empty real placeholder at the same index

#### Scenario: Cover still required
- **WHEN** an operator saves a portfolio
- **THEN** cover required behavior is unchanged by this gallery-pair rule change

### Requirement: Silently drop pairs with both sides empty
On save, pairs where both render and real are empty SHALL be omitted from the submitted `renders` and `reals` arrays without blocking the save or showing a validation error for that empty pair.

#### Scenario: Empty pair skipped
- **WHEN** the editor contains a pair with both sides empty among other valid pairs
- **THEN** save succeeds
- **AND** that empty pair does not appear as an index in the submitted arrays

### Requirement: Persist equal-length aligned arrays
Valid pairs SHALL still be submitted as equal-length `renders` and `reals` arrays in pair order, using empty dual placeholders for the missing side.

#### Scenario: Mixed optional sides keep alignment
- **WHEN** the operator saves pairs that are render-only, real-only, and both-filled in that order
- **THEN** `renders.length` equals `reals.length`
- **AND** index 0 has empty real placeholder, index 1 has empty render placeholder, index 2 has both non-empty

### Requirement: Editor copy reflects optional sides
The pair editor UI SHALL label both render and real slots as optional (or otherwise clearly indicate neither side is mandatory within a pair).

#### Scenario: Labels show optional
- **WHEN** an operator views a gallery pair row
- **THEN** neither slot is presented as the sole required field of the pair
