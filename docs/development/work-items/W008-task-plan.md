# Task Plan — W008: Guest Session

Specification: `docs/development/work-items/W008-guest-session.md` (SPEC_READY; approved policy integrated by PR #10; T01 design PASS at `3d9ed43a3215f339e19f10464aaedfa241f866f8`; integrated by PR #18 merge commit `774ec869a2ac2c96914cfb60a0bfe9cb7b70a08d`)
Branch: `feature/W008-guest-session-implementation` (T02 implementation; based on updated develop at `774ec869a2ac2c96914cfb60a0bfe9cb7b70a08d`)
Current execution: T02 — IN_PROGRESS. T03 HTTP/cookie work and T04 frontend work have not started.

| Task | Change | Requirement / acceptance refs | Deliverable | Validation |
|---|---|---|---|---|
| T01 | Finalize API, expiry activity semantics, expiry-to-W010/W013 event contract, cookie dev/prod split, trusted-proxy policy, grapheme/NFC dependencies, Origin/CSRF behavior, and limiter details | §§6.1–6.7, §20 SEC-006, W008-AC-04–07, 09 | Reviewed technical design in W008 specification; ADR only if a material architecture decision requires one | Documentation consistency, requirement traceability, `git diff --check`, `pnpm check`, independent exact-revision design review; integrated by PR #18 |
| T02 | Implement pure guest-session domain, configurable idle/absolute expiry, credential verifier/revocation, creation limiter, trusted source resolution, and strict username validator | FR-001–003, W008-AC-01–02, 06–07, 09 | `apps/server/internal/session/`, `apps/server/internal/clientip/`, and `apps/server/internal/clock/` | Deterministic expiry/activity tests; Unicode grapheme and approved-input tests; limiter boundary, proxy trust, concurrency and cleanup tests; bounded expiry outbox/retry/backpressure and race tests; `go test`, `go test -race`, `go vet` |
| T03 | Add same-origin HTTP session creation and recovery adapter with secure production cookie and isolated development configuration | FR-003, FR-007, §20 SEC-006, W008-AC-03–07, 09 | Versioned HTTP endpoints, cookie policy, API errors, Vite `/api` proxy | HTTP tests for attributes, Origin, Fetch Metadata, expiry, limits, reuse, no credential leakage; security review |
| T04 | Add Thai client username entry and session bootstrap/recovery | FR-001–004, W008-AC-08 | React entry and recovery flow; display-name persistence only | Frontend component/state tests; browser refresh/session behavior |
| T05 | Integrate docs, validation and evidence; request exact-head independent implementation review | W008-AC-09–10 | Implementation report and PR evidence | `pnpm check`, CI/security checks, independent review PASS and re-review after material fixes |

## Dependencies

- W004 Development Standards & CI (DONE).
- W005 Protocol Contracts (DONE).
- W009 and W010 depend on the authenticated guest-session boundary produced here; W013 consumes the session-expiry forfeit signal.

## Out of scope

- Rooms, duplicate names within a Room, lobby name changes, WebSocket connections, active-connection arbitration, reconnect handoff, gameplay, durable account storage, or persistent session storage.
- Changing §6 username/session policy without Project Owner approval.
- Room-operation and gameplay abuse controls owned by W009/W013; W008 must provide the approved configurable session-creation limit before public exposure.

## Completion evidence

- [x] Exact-revision specification review PASS and approved architecture followed
- [ ] Actual validation outcomes recorded in `docs/development/implementation-reports/W008.md`
- [ ] Independent review PASS for exact commit SHA
- [ ] Blocking findings resolved and re-reviewed
- [ ] Documentation synchronized
- [ ] PR merged to `develop`; feature branch retained

## T01 design acceptance

- [x] HTTP request/response, status/error, cookie, reuse, and Origin contract is explicit.
- [x] Credential lifecycle, effective expiration, qualifying activity, invalidation, and cleanup are explicit and preserve §6.5.
- [x] Trusted-proxy source resolution and successful-allocation-only rolling session quota are explicit and testable.
- [x] Username validation order, Unicode cases, and W009 duplicate-key boundary are explicit and preserve §6.2.
- [x] W008-to-W010/W013 expiry event ownership, ordering, idempotency, and deadline tie behavior are explicit without circular dependencies.
- [x] Session expiry and Match deadlines share one monotonic clock domain; UTC remains diagnostic, with wall-clock adjustment and delayed-delivery boundary tests specified.
- [x] Proactive session expiry emits the event at the authoritative deadline without requiring a later request or active socket; scheduler state is bounded by active sessions and shuts down cleanly.
- [x] Selected NFC/grapheme dependency versions, licenses, source, and implementation-time checks are recorded; no dependency is installed during T01.
- [x] Owner approval of equal-time Match-deadline precedence is recorded and flagged for synchronization into W013 before its implementation.
- [x] Independent reviewer PASS covers the final T01 design revision; that review was completed before T02 implementation began.
