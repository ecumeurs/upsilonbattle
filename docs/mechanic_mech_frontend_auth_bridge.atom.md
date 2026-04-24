---
id: mechanic_mech_frontend_auth_bridge
status: STABLE
human_name: Frontend Auth Bridge
type: MECHANIC
priority: 3
tags: auth,axios,security
dependents: []
layer: IMPLEMENTATION
parents:
  - [[req_security_token_ttl]]
version: 1.0
---

# New Atom

## INTENT
Provide a centralized Axios instance that automatically manages JWT tokens, standardized API envelopes, and automated token renewal.

## THE RULE / LOGIC
- **Authorization:** Every request MUST include the `Authorization: Bearer <token>` header, derived from the `upsilon_token` key in `localStorage`.
- **Request ID:** Every request MUST include a unique UUIDv7 in the `X-Request-ID` header.
- **Envelope Handling:** All responses MUST be unwrapped from the standardized Upsilon envelope `{ success, data, message, meta }`, returning `data` to the caller.
- **Token Renewal:** If a response contains `meta.token`, the bridge MUST update the `upsilon_token` in `localStorage` immediately to ensure subsequent requests use the new token.
- **Error Propagation:** 4xx and 5xx errors from the backend must be normalized and rejected with the error details from the envelope.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_frontend_auth_bridge]]`
- **Stored Tokens:** `localStorage['upsilon_token']`, `localStorage['upsilon_user']`
- **Request ID Header:** `X-Request-ID`

## EXPECTATION
- Requests to protected routes fail with 401 if `upsilon_token` is missing or invalid.
- Requests succeed and include the bearer token if `upsilon_token` is present.
- If a response includes `meta.token`, `localStorage['upsilon_token']` is updated with the new value.
