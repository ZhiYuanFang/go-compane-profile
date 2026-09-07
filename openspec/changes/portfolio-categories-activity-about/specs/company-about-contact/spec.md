## ADDED Requirements

### Requirement: Company profile includes awards and contact fields
The company profile SHALL include `awards` (multiline text), `phone`, and `wechat` in addition to existing fields. `address` SHALL continue to be the company address field used by the about page. Public `GET /api/v1/company` and admin company read/update SHALL include these fields. Newlines in `awards` SHALL be preserved in storage and API responses.

#### Scenario: Admin updates about fields
- **WHEN** an admin saves awards, phone, and wechat on the company profile page
- **THEN** subsequent public company reads return those values
- **AND** awards retains newline characters as entered

#### Scenario: Address reuse
- **WHEN** the about page displays address
- **THEN** it uses the existing `company.address` value

### Requirement: Pricing capability is removed
The system SHALL remove pricing end-to-end: database table usage, public and admin pricing APIs, admin pricing UI, and mini-program pricing page and entry. Clients MUST NOT rely on `GET /api/v1/pricing`.

#### Scenario: Pricing endpoints gone
- **WHEN** a client calls the former public or admin pricing endpoints after this change
- **THEN** the endpoints are unavailable (not found or removed from routing)

#### Scenario: Admin nav has no pricing
- **WHEN** an admin views the sidebar
- **THEN** there is no 资费 entry

### Requirement: Mini-program About Us page replaces pricing entry
The mini-program SHALL replace the former pricing FAB/entry with an About Us entry that opens an About page. The About page SHALL show scrollable awards text (honoring newlines) and a fixed bottom section with: address plus copy action, phone that initiates a dial action, and wechat plus copy action.

#### Scenario: Open About Us
- **WHEN** the user taps the About Us entry on home
- **THEN** the About page opens with awards content and the fixed contact footer

#### Scenario: Contact actions
- **WHEN** the user taps copy on address or wechat
- **THEN** the corresponding value is copied to the clipboard
- **WHEN** the user taps the phone row
- **THEN** the mini-program initiates a phone call to `phone`
