# Task Plan — W012: Results & Rematch

Specification: `docs/development/work-items/W012-results-rematch.md` (READY_FOR_REVIEW; implementation requires approved spec and W009–W011/W013 dependencies)
Branch: `feature/W012-results-rematch-spec`

| Task | Change | Requirement / acceptance refs | Deliverable | Validation |
|---|---|---|---|---|
| T01 | Resolve owner decisions on Rematch initiator/consent and forfeited ranking/display; reconcile approved W009–W011/W013 contracts and protocol fields | §§7, 9, 12–13, 17–20; W012-AC-01–09 | Approved technical design and protocol decision | Owner decisions; independent design review; required dependencies merged |
| T02 | Implement immutable result snapshot, eligibility-aware score ranking, shared ties, and no-winner cases | FR-044, FR-074, INV-010/011/018, W012-AC-01–03 | Results domain projection | Unit tests for ranks, ties, forfeits, cancellation, immutability |
| T03 | Implement Room Results→Lobby/Rematch command and Match reset/preservation integration | FR-075, INV-014/015, W012-AC-05–07 | Typed W009/W011 rematch transition | Integration/idempotency tests, Room state and readiness checks |
| T04 | Extend protocol and build Thai accessible Results/Rematch UI | FR-074–075, W012-AC-03–06 | Strict schema/events, Results screen and Lobby navigation | Protocol fixtures; frontend tests; 2–6 context browser E2E; mobile/a11y checks |
| T05 | Validate full lifecycle, record evidence, and request exact-head independent review | W012-AC-01–09 | Implementation report and PR evidence | `pnpm check`, Go race tests, CI/security and E2E; independent PASS/re-review |

## Dependencies

- W009 Room Management, W010 WebSocket Foundation, W011 Multiplayer Speed Battle and W013 Reconnection & Anti-Griefing approved and integrated before implementation. Specification review may occur in parallel, but implementation/integration order places W013 before W012 because W012 consumes W013's final eligibility/forfeit state.
- W005 Protocol Contracts (DONE); W006 Card Matching Engine (DONE); W007 local prototype (DONE); W008 Guest Session approved and merged.

## Out of scope

- Gameplay scoring/timing, Room management, reconnection/forfeit policy, additional modes, durable Match history, and public leaderboards.
- Changing confirmed tie/winner/Match reset behavior without Project Owner approval.

## Completion evidence

- [ ] Owner decisions and approved specification recorded
- [ ] Actual validation outcomes recorded per acceptance criteria
- [ ] Independent review PASS for exact implementation commit
- [ ] Blocking findings resolved and re-reviewed
- [ ] Documentation synchronized
- [ ] PR merged to `develop` with Merge Commit; feature branch retained
