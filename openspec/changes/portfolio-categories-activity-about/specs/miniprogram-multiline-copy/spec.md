## ADDED Requirements

### Requirement: Heart flow newlines render on detail
The portfolio detail page SHALL display `heartFlow` such that newline characters from the API appear as line breaks on screen (not collapsed into a single line).

#### Scenario: Multiline heart flow
- **WHEN** the detail page receives a `heartFlow` value containing `\n` separators
- **THEN** the user sees distinct lines corresponding to those separators

### Requirement: Awards newlines render on About page
The About Us page SHALL display `awards` such that newline characters from the company API appear as line breaks on screen.

#### Scenario: Multiline awards
- **WHEN** the About page receives an `awards` value containing `\n` separators
- **THEN** the user sees distinct lines corresponding to those separators

### Requirement: Prefer reliable mini-program newline strategy
Multiline CMS text on the mini-program SHALL use a rendering approach known to preserve newlines on WeChat (for example splitting on `\n` into multiple nodes, or an equivalent `<text>` strategy), rather than relying solely on `white-space: pre-wrap` on a `view`.

#### Scenario: Implementation does not depend on view pre-wrap alone
- **WHEN** multiline copy is rendered for heart flow or awards
- **THEN** the implementation uses split-lines or `<text>`-based preservation
- **AND** line breaks remain visible on device/simulator
