# Task Plan — W011: Multiplayer Speed Battle

Specification: `docs/development/work-items/W011-multiplayer-speed-battle.md` (READY_FOR_REVIEW; implementation requires approved spec and W008–W010 dependencies)
Branch: `feature/W011-multiplayer-speed-battle`

| Task | Change | Requirement / acceptance refs | Deliverable | Validation |
|---|---|---|---|---|
| T01 | Approve lifecycle/deadline/input policy, tuning config, W008–W010 contracts, and W012/W013 interfaces | §§9, 12–13, 17–20; W011-AC-01–11 | Approved technical design and ADR if material architecture change | Owner decisions; independent design review; W008–W010 merged |
| T02 | Compose W006 pure game engine into per-Room Match aggregate with private player states and seeded solvable starts | FR-024–040, FR-065, INV-005/008, W011-AC-01–03 | `apps/server/internal/match/` domain | Go tests for start conditions, determinism, board privacy, solvability and settings |
| T03 | Implement authoritative action processor, timer/cooldown, idempotency integration, deadline ordering, and completion/forfeit hooks | FR-041–044, FR-065–073, INV-004–010/012–017, W011-AC-03–08 | Room-loop Match coordinator and typed events | Fake-clock, duplicate/stale, concurrency/race and security tests |
| T04 | Integrate W010 protocol projections and replace W007 local online authority with server-synchronized React/Phaser gameplay | FR-064–076, W011-AC-06,09 | Online Speed Battle UI and private/public event projections | Protocol fixtures; frontend tests; 2–6 client browser E2E, mouse/touch/keyboard/mobile checks |
| T05 | Validate, record evidence, request exact-head independent review, and prepare W012/W013 integration seam | W011-AC-01–11 | Implementation report and PR evidence | `pnpm check`, Go race tests, CI/security checks, independent review PASS and exact-SHA re-review after remediation |

## Dependencies

- W005 Protocol Contracts and W006 Card Matching Engine (DONE).
- W007 Local Gameplay Prototype (DONE; temporary online authority must be removed/isolated).
- W008 Guest Session, W009 Room Management, W010 WebSocket Foundation approved and merged.
- W012 consumes completion snapshot; W013 owns reconnection/forfeit eligibility. Implement against approved interfaces or test fakes until their integration work lands.

## Out of scope

- Results/winner UI and Rematch, reconnect grace/forfeit policy, Host migration timing, non-Speed-Battle modes, durable state, and multi-instance deployment.
- Changing BASELINE gameplay behavior or TUNING values without approved playtest evidence.

## Completion evidence

- [ ] Approved specification and W008–W010 contracts followed
- [ ] Actual validation outcomes recorded per acceptance criteria
- [ ] Independent review PASS for exact implementation commit
- [ ] Blocking findings resolved and re-reviewed
- [ ] Documentation synchronized; W012/W013 seams recorded
- [ ] PR merged to `develop` with Merge Commit; feature branch retained
