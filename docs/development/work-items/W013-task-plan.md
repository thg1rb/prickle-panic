# Task Plan — W013: Reconnection & Anti-Griefing

Specification: `docs/development/work-items/W013-reconnection-anti-griefing.md` (READY_FOR_REVIEW; implementation requires approved policy decisions and W008–W011 dependencies)
Branch: `feature/W013-reconnection-anti-griefing-spec`

| Task | Change | Requirement / acceptance refs | Deliverable | Validation |
|---|---|---|---|---|
| T01 | Approve Host continuity, deadline tie ordering, operational rate/queue thresholds, and reconcile W008–W011 interfaces | §§13–14, 18–20; W013-AC-01–11 | Approved design and dependency contracts | Owner decisions; independent design review; dependency PRs approved |
| T02 | Implement disconnect state, 15-second grace, authorized reconnect, connection replacement, and voluntary/expired forfeit | §§13.1–13.8, INV-016/017; W013-AC-01–05 | Room-loop reconnect/forfeit commands | Clock-controlled tests; boundary/concurrency tests |
| T03 | Integrate terminal eligibility/Host continuity and lifecycle cleanup with Room/Match owners | EC-010/011/035/036, INV-015/018; W013-AC-06–07/09 | Versioned eligibility and Host handoff | Integration tests across W009/W011/W012 |
| T04 | Add bounded rate limiting, idempotency/stale checks, slow-client handling, and safe restart recovery | §§14, 19.10–12, 20; W013-AC-08–09 | Security/resource controls and recovery UI | Abuse, load-bound, security, and UI tests |
| T05 | Validate full workflows, record evidence, and request exact-head independent implementation review | W013-AC-01–11 | Implementation report and PR evidence | Go tests/race/vet, `pnpm check`, CI/security, multi-context E2E, exact-SHA review |

## Dependencies

- W008 Guest Session, W009 Room Management, W010 WebSocket Foundation, and W011 Multiplayer Speed Battle approved and merged.
- W005 Protocol Contracts (DONE); W006 Card Matching Engine and W007 Local Gameplay Prototype (DONE).
- W012 Results & Rematch must consume the approved atomic/versioned terminal eligibility handoff; W012 integration follows W013.

## Out of scope

- Durable sessions/Match recovery, multi-instance deployment, permanent identity/bans, complete bot prevention, scoring/card rule changes, Results ranking/rematch implementation, and additional modes.
- Changing confirmed 15-second grace, voluntary forfeit, or Match-end-during-grace behavior.

## Completion evidence

- [ ] Owner decisions and approved specification recorded
- [ ] Actual validation outcomes recorded per acceptance criteria
- [ ] Independent review PASS for exact implementation commit
- [ ] Blocking findings resolved and re-reviewed
- [ ] Documentation synchronized
- [ ] PR merged to `develop` with Merge Commit; feature branch retained
