# Task Plan — W010: WebSocket Foundation

Specification: `docs/development/work-items/W010-websocket-foundation.md` (READY_FOR_REVIEW; implementation requires approved spec and W008/W009 contracts)
Branch: `feature/W010-websocket-foundation`

| Task | Change | Requirement / acceptance refs | Deliverable | Validation |
|---|---|---|---|---|
| T01 | Approve session/Room contracts; verify `coder/websocket` version/API/security; decide frame, queue, cache, rate, timeout, Origin and replacement settings; approve and resolve W005 recipient-stream sequence semantics | §§6, 13–20; W005; W010-AC-01–11 | Approved technical design and ADR if material architecture change | Recorded owner decision; independent design review; W008 and W009 merged |
| T02 | Add authenticated WebSocket upgrade, strict Origin/session policy, bounded parsing, and stable error/close mapping | SEC-001–006, SEC-013–015, W010-AC-01,03,08 | Gateway handshake and request decoder | Go HTTP/WebSocket integration tests; `go test`, `go vet` |
| T03 | Add global connection registry, generation-based replacement, Room/W013 signals, dispatcher, request dedupe, and bounded outbound writer | FR-008, FR-064, FR-071–073, INV-004/014/017, W010-AC-02,04–08 | Connection manager and typed Room/Match adapters | Race, duplicate, stale-generation, slow-client, cross-Room, cleanup tests |
| T04 | Implement typed React WebSocket client and Room protocol synchronization/recovery integration | FR-064, W010-AC-04–05,09 | Client service and application state bridge | Frontend tests; multi-context browser E2E; Thai errors/status; credential storage audit |
| T05 | Validate protocol coverage, security/resource bounds, document evidence, and request exact-head independent review | W010-AC-01–11 | Implementation report and PR evidence | `pnpm check`, `pnpm protocol:check`, Go race tests, CI/security checks, independent review PASS/re-review |

## Dependencies

- W005 Protocol Contracts (DONE).
- W008 Guest Session approved and merged; W010 does not implement a substitute auth mechanism.
- W009 Room Management approved and merged; W010 dispatches to its domain boundary and does not duplicate Room authority.
- W011 consumes the gameplay dispatch seam; W013 consumes unexpected-disconnect/replacement notifications.

## Out of scope

- Room business behavior, game rules/scoring/timing, Results/rematch, reconnect grace/forfeit policy, additional modes, durable state, and multi-instance routing.
- Changing confirmed session, Room, or protocol requirements through transport implementation.

## Completion evidence

- [ ] Approved specification and W008/W009 API boundaries followed
- [ ] Actual validation outcomes recorded
- [ ] INV-004 and connection-replacement evidence tied to exact SHA
- [ ] Independent review PASS for exact implementation commit
- [ ] Blocking findings resolved and re-reviewed
- [ ] Documentation synchronized
- [ ] PR merged to `develop` with Merge Commit; feature branch retained
