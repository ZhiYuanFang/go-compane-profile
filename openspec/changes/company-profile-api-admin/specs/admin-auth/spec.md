## ADDED Requirements

### Requirement: Password login for admin
The system SHALL allow admin access only after successful login with the configured admin password from environment configuration (`ADMIN_PASSWORD`). The password MUST NOT be hardcoded in source.

#### Scenario: Correct password
- **WHEN** a client posts the correct password to the admin login endpoint
- **THEN** the system establishes an authenticated admin session (cookie or token) and allows subsequent admin API calls

#### Scenario: Incorrect password
- **WHEN** a client posts an incorrect password
- **THEN** the system rejects the login and does not grant admin access

### Requirement: Admin API protection
The system SHALL reject unauthenticated requests to admin API routes other than login.

#### Scenario: Unauthenticated admin call
- **WHEN** a client calls a protected admin API without a valid session
- **THEN** the system responds with an authentication error

### Requirement: Logout
The system SHALL invalidate the admin session on logout.

#### Scenario: Logout clears access
- **WHEN** an authenticated admin logs out
- **THEN** subsequent admin API calls with the old session fail authentication
