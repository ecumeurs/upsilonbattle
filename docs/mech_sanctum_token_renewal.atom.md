---
id: mech_sanctum_token_renewal
human_name: Sliding Token Renewal Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: STABLE
priority: 5
tags: [auth, tokens, middleware]
parents:
  - [[shared:req_security_token_ttl]]
dependents: []
---
# Sliding Token Renewal Mechanic

## INTENT
To implement a proactive token renewal system that extends user sessions without requiring re-authentication, while ensuring security through short TTLs and grace periods.

## THE RULE / LOGIC
The renewal process runs in the hub's `TokenRenewal` middleware, after bearer authentication:

1. **Authentication:** The request must carry a valid opaque bearer token (personal access token).
2. **Age Calculation:**
   - `CurrentAge = Now - Token.CreatedAt`
3. **Trigger Condition:**
   - If `CurrentAge >= RenewAfter (10 minutes)` AND `CurrentAge < TokenTTL (15 minutes)`:
     - **Check for Active Grace Period:** If `Token.ExpiresAt` is set and near (`< GraceTTL = 20 seconds`), skip renewal (already in progress).
     - **Issue New Token:** Create a new personal access token for the user with `ExpiresAt = Now + TokenTTL`.
     - **Set Grace Period:** Update the *current* token's `ExpiresAt` to `Now + GraceTTL`.
     - **Store for Response:** Stash the new plaintext token in the request context.
4. **Response Modification:**
   - The envelope writer injects, when a new token was issued:
     - `meta.token = <NewToken>`
     - `meta.message = "Token renewed"`

## TECHNICAL INTERFACE (The Bridge)
- **Middleware:** `upsilonhub/internal/gateway/middleware/auth.go` `TokenRenewal` (constants in `internal/platform/identity`)
- **Code Tag:** `@spec-link [[mech_sanctum_token_renewal]]`
- **Tests:** `upsilonhub/internal/gateway/token_renewal_test.go`

## EXPECTATION (For Testing)
- Requests at T+9m -> No renewal.
- Requests at T+11m -> New token in meta, old token persists for 20s.
- Requests at T+11m05s (using old token) -> Normal response, no new renewal triggered.
- Requests at T+11m25s (using old token) -> 401 Unauthorized.
- A client idle past TokenTTL between requests is expired by design (see ISS-105 for the CLI keepalive direction).
