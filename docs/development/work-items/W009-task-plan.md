# Task Plan — W009: Room Management

Specification: `docs/development/work-items/W009-room-management.md` (READY_FOR_REVIEW; implementation requires approved spec and W008 boundary)
Branch: `feature/W009-room-management`

| Task | Change | Requirement / acceptance refs | Deliverable | Validation |
|---|---|---|---|---|
| T01 | Resolve owner decisions; review Room code and username policy, settings allow-list, empty-Room close behavior, protocol extension, and W008 session contract | §§6–7, W009-AC-01–09 | Approved technical design and ADR if architecture materially changes | Recorded owner decisions; independent spec review PASS; W008 approval/merge |
| T02 | Implement pure Room aggregate, states/transitions, membership, Host permissions, readiness/settings, and immutable snapshots | FR-006, FR-009, FR-013–020, W009-AC-01,03–07 | `apps/server/internal/room/` domain | Unit tests/table tests; `go test`, `go vet` |
| T03 | Implement Room manager, code/index uniqueness, bounded per-Room processing, rate limits, expiry and cleanup | FR-010–013, FR-022–023, SEC-006/009, W009-AC-01,03,06–08 | In-memory RoomManager and lifecycle | Concurrency/race tests; collisions, last-slot join, saturation, timer and shutdown tests |
| T04 | Extend strict v1 protocol schemas, Room commands/events/errors and fixtures; bind transport through W010-compatible interface | FR-011–021, W005 contract, W009-AC-02,04,09 | Updated schema, fixtures, typed server/client boundary | `pnpm protocol:check`; compatibility/security review |
| T05 | Implement Thai Room create/join/invitation/Lobby UI and responsive interactions | FR-009, FR-011–020, W009-AC-02,04,05 | React Room feature and typed client projection | Frontend tests; two-context browser E2E; desktop/mobile and accessibility checks |
| T06 | Integrate, validate, record evidence, and request exact-head independent implementation review | W009-AC-01–11 | Implementation report and PR evidence | `pnpm check`, CI/security checks, Go race tests, E2E where available; independent review PASS and exact-SHA re-review after fixes |

## Dependencies

- W001–W007 (DONE); W005 protocol contract (DONE).
- W008 approved and merged for authenticated Guest Session principal, username policy, and membership identity.
- W010 protocol/runtime integration is coordinated; Room domain can be built before WebSocket integration, but end-to-end acceptance cannot be claimed until W010 supplies its approved adapter.
- W011/W012 provide the authoritative Match coordinator and callbacks for countdown, in-progress, results, and return to Lobby.
- W013 supplies unexpected-disconnect reservations, grace period, forfeit, and delayed Host migration.

## Out of scope

- Guest session creation/recovery, WebSocket connection lifecycle, gameplay/match engine, results/rematch, unexpected disconnect policy, and additional game modes.
- Changing BASELINE Room behavior, FR priorities, or TUNING values without approved design/playtest evidence.
- Public matchmaking, room directory, durable persistence, or multiple server instances.

## Completion evidence

- [ ] Owner decisions and approved specification recorded
- [ ] W008 contract dependency approved and merged
- [ ] Actual validation outcomes recorded
- [ ] Independent review PASS for exact implementation SHA
- [ ] Blocking findings resolved and re-reviewed
- [ ] Documentation synchronized
- [ ] PR merged to `develop` with Merge Commit; feature branch retained
