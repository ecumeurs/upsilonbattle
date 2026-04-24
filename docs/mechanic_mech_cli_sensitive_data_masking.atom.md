---
id: mechanic_mech_cli_sensitive_data_masking
status: STABLE
type: MECHANIC
parents:
  - [[api_auth_login]]
  - [[api_auth_register]]
  - [[api_auth_user]]
dependents: []
layer: IMPLEMENTATION
priority: 3
version: 1.0
human_name: CLI Sensitive Data Masking
---

# New Atom

## INTENT
To ensure user credentials and sensitive data are never exposed in plaintext during interactive CLI sessions or diagnostic logs.

## THE RULE / LOGIC
- **Input Masking:** If a parameter is marked as `Secret`, the CLI MUST use a non-echoing input method.
- **Output Masking:** Before printing any `[CURL]` command, the system MUST inspect the request body (if JSON) and replace any recognized sensitive top-level keys with a fixed mask string.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_mech_cli_sensitive_data_masking]]`
- **Masked Fields:** `password`, `password_confirmation`, `current_password`, `token`
- **Mechanism:** `readline.ReadPassword` for input, JSON manipulation for Curl output.

## EXPECTATION
- When prompting for a param marked as Secret, the input is obscured (masked).
- When printing a Curl command, any JSON body field named 'password', 'password_confirmation', 'current_password', or 'token' is replaced with '********'.
