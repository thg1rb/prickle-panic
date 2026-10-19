# Task Plan — W008: Guest Session

Specification: `docs/development/work-items/W008-guest-session.md` (READY_FOR_REVIEW; approved revision required before implementation)
Branch: `feature/W008-guest-session`

| Task | Change | Requirement / acceptance refs | Deliverable | Validation |
|---|---|---|---|---|
| T01 | Finalize session API, cookie, expiry, verifier, Origin/CSRF, Unicode dependency, and unauthenticated session-creation rate-limit decisions | §6.1–6.7, §20 SEC-006, W008-AC-04–07, 09 | Approved technical design and ADR if a material architecture decision is made | Independent design review; documented owner approval for material user-facing/security policy decisions; endpoint exposure blocked until rate-limit policy and tests are specified |
| T02 | Implement pure guest-session domain and username validator | FR-001–003, W008-AC-01–02, 06–07 | `apps/server/internal/session/` | Unit, Unicode boundary, concurrency/race tests; `go test`, `go vet` |
| T03 | Add session HTTP creation and recovery adapter | FR-003, FR-007, W008-AC-03–07 | Versioned HTTP endpoints, cookie policy, API errors | HTTP tests for attributes, Origin, limits, no credential leakage; security review |
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
