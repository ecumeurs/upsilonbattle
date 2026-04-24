---
id: mech_web_catchall_router
human_name: "SPA Web Catch-all Router"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 2
tags: [routing, inertia, spa]
parents:
  - [[ui_tactical_infrastructure]]
dependents: []
---

# SPA Web Catch-all Router

## INTENT
To enable client-side routing for the Inertia entry points while explicitly preserving system routes required for health checks and API functionality.

## THE RULE / LOGIC
- All GET requests that do NOT match explicit routes must be captured by a catch-all route.
- **Exclusion Rule:** The catch-all MUST exclude:
    - Anything starting with `api/` (API logic).
    - The exact path `up` (System health check).
- When a match occurs, the router must render the `Welcome` view via Inertia as a fallback entry point.

## TECHNICAL INTERFACE (The Bridge)
- **Regex Filter:** `^(?!api\/|up$).*`
- **Code Tag:** `@spec-link [[mech_web_catchall_router]]`
- **Target Handler:** `Inertia::render('Welcome')`

## EXPECTATION (For Testing)
- `GET /random-non-existent-path` renders the Welcome page.
- `GET /api/v1/help` is NOT caught by this router (passes through to API).
- `GET /up` is NOT caught by this router (passes through to Framework Health).
