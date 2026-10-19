# Task Plan — W008: Guest Session

Specification: `docs/development/work-items/W008-guest-session.md` (READY_FOR_REVIEW; Project Owner approval required before implementation)
Branch: `feature/W008-guest-session`

| Task | Change | Requirement / acceptance refs | Deliverable | Validation |
|---|---|---|---|---|
| T01 | Review proposed API, cookie, expiry, verifier, Origin/CSRF, Unicode dependency, and limiter details; obtain Project Owner approval | §6.1–6.7, §20 SEC-006, W008-AC-04–07, 09 | Approved technical design and ADR if a material architecture decision is made | Independent design review; documented owner approval for proposed policy values; endpoint exposure blocked until limiter and capacity behavior are tested |
| T02 | Implement pure guest-session domain, expiry, credential verifier, limiter/capacity interfaces, and username validator | FR-001–003, W008-AC-01–02, 06–07, 09 | `apps/server/internal/session/` | Unit, Unicode 17 boundary, limiter, capacity, concurrency/race tests; `go test`, `go vet` |
| T03 | Add same-origin HTTP session creation and recovery adapter with secure cookies | FR-003, FR-007, §20 SEC-006, W008-AC-03–07, 09 | Versioned HTTP endpoints, cookie policy, API errors, Vite `/api` proxy | HTTP tests for attributes, Origin, Fetch Metadata, limits, no credential leakage; security review |
| T04 | Add Thai client username entry and session bootstrap/recovery | FR-001–004, W008-AC-08 | React entry and recovery flow; display-name persistence only | Frontend component/state tests; browser refresh/session behavior |
| T05 | Integrate docs, validation and evidence; request exact-head independent review | W008-AC-09–10 | Implementation report and PR evidence | `pnpm check`, CI/security checks, independent review PASS and re-review after material fixes |

## Dependencies

- W004 Development Standards & CI (DONE).
- W005 Protocol Contracts (DONE).
- W009 and W010 depend on the authenticated guest-session boundary produced here.

## Out of scope

- Rooms, duplicate names within a Room, lobby name changes, WebSocket connections, active-connection arbitration, reconnect handoff, gameplay, durable account storage, or persistent session storage.
- Changing §6 username/session requirements or deciding policy through implementation defaults without approval.
- Broader abuse-control work and later rate-limit tuning owned by W013; W008 must provide an initial configurable session-creation limit before public exposure.

## Completion evidence

- [ ] Approved specification and architecture followed
- [ ] Actual validation outcomes recorded
- [ ] Independent review PASS for exact commit SHA
- [ ] Blocking findings resolved and re-reviewed
- [ ] Documentation synchronized
- [ ] PR merged to `develop`; feature branch retained
