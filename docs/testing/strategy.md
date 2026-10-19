# Testing and Validation Strategy

This strategy implements `docs/REQUIREMENT.md` §27–29. The currently implemented baseline checks and local commands are defined in [`../development/standards.md`](../development/standards.md). It does not claim coverage until checks exist and results are recorded.

## Test levels

- **Unit:** Pure card/match/scoring rules and frontend behavior; map tests to requirement IDs or acceptance criteria.
- **Integration:** Go Room, session, and WebSocket boundaries; verify authorization, ordering, idempotency, state privacy, and lifecycle.
- **End-to-end:** Playwright with multiple browser contexts for create/join, ready/start, gameplay, score sync, end/results, rematch, and reconnect.
- **Concurrency:** Go race detector and deterministic action scheduling for same-card contests, duplicate starts/actions, timer boundaries, reconnection, and event/state transitions.
- **Security:** Invalid sessions/Origins, unauthorized actions, malformed/oversized payloads, rate limits, duplicate actions, and resource bounds.
- **Performance and network:** Six players per Room, slow clients, disconnections, and 0/100/200/300 ms latency scenarios; record measurements before making claims.
- **Browser and responsive:** Chrome, Safari, Edge, Firefox; mobile widths starting at 320 CSS px; touch and mouse; Thai font and layout review.
- **Visual/UX:** Card attribute distinction, Coco consistency, reduced motion, contrast, and no-scroll 4×4 readability. Use visual regression only where stable and useful.

## Core engine coverage

Test all 32 × 32 = 1,024 card-variant pairs against Hamming distance ≤ 1. Cover card-instance identity, valid/invalid selections, deterministic generation, solvability after refill, cooldown versus technical rejection, Combo, score, and Match termination.

## Evidence rules

Record exact commands/scenarios and actual results as PASS, FAIL, NOT_RUN, BLOCKED, or NOT_APPLICABLE with reasons. Distinguish developer, CI, and independent reviewer execution. No test result may be inferred from a build, and no unexecuted test may be reported as passing. Reviewer evidence must identify the exact reviewed commit SHA.

## Quality gates

Feature PRs target `develop` and run formatting, lint, type-check, unit, applicable integration, build, protocol, and relevant security checks. Release PRs from `develop` to `main` add full regression, multiplayer integration, E2E, browser compatibility, release-build, and release review. Go concurrency-sensitive packages use `go test -race ./...` when available. W004 establishes the actual CI workflows and required status checks.
