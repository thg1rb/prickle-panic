# W015 — Closed Alpha Readiness

Status: READY_FOR_REVIEW  
Type: Integration Validation / Release Readiness  
Owner: Project Owner assigns during execution  
Last updated: 2026-10-19

## Goals

- Validate the integrated Speed Battle MVP against approved requirements and specifications before a Closed Alpha release candidate is declared ready.
- Exercise the complete private multiplayer path across multiple browser contexts and supported browser families, including reconnection, results, rematch, security, responsive layout, and documented failure cases.
- Produce a reproducible evidence package with actual test outcomes, known limitations, unresolved risk, and a clear readiness recommendation for the Project Owner.
- Separate technical readiness from authorization to merge a release, deploy, publish, or invite Closed Alpha participants.

## Non-goals

- Implementing new gameplay or silently waiving failed acceptance criteria.
- Expanding scope to Classic Rush or Party Chaos, public matchmaking, durable player identity, or persistent Match history.
- Deploying production, publishing a release, inviting external testers, changing repository protections, or granting access without explicit applicable authorization.
- Claiming cross-browser/device/performance/security coverage from a single local browser or CI unit test.

## Scope

- Confirm all required MVP Work Items and integration PRs are approved, merged to `develop`, traceable to requirements, and have exact-revision independent review evidence.
- Full functional E2E with 2–6 browser contexts: guest entry, create/join by code and invitation URL, Lobby/readiness, Host start, Speed Battle, scoring, end conditions, Results, Rematch, reconnect/forfeit, and recovery behavior.
- Browser compatibility for initial targets Chrome, Safari, Edge, Firefox on desktop and mobile/tablet layouts; record actual browser/OS/device/version and unavailable combinations.
- Concurrency and network testing for action ordering, duplicate/stale messages, simultaneous selections, Match-end races, disconnect/reconnect boundaries, slow clients, and rate limiting.
- Security review of session authorization, cookie handling, WebSocket Origin, request validation, permissions, payload limits, abuse controls, safe logging, dependency/asset/Skill provenance, and secrets handling.
- Performance/load measurements for server/Room isolation, 6-player Room capacity, client frame/asset loading, and 0/100/200/300 ms network scenarios; compare to documented targets and report measured values without inventing pass thresholds.
- Visual/UX and accessibility review of Thai rendering, Coco coverage across all §2.8 surfaces, 4×4 card readability, touch/keyboard interaction, focus/contrast, reduced motion, audio settings, and mobile viewport constraints.
- Release build verification, deployment configuration review, health checks, HTTPS/WSS routing review, rollback/recovery plan, and release documentation as read-only readiness evidence. Actual deploy/release remains separately authorized.
- Independent acceptance review of the integrated evidence and exact release candidate revision.

## Requirement references and status

- `docs/REQUIREMENT.md` §§2.5 (supported browsers/platforms), §3.1 (MVP), §§13–14 (reconnection/anti-griefing), §§18–22 (architecture, protocol, security, NFR, responsive UI), §§26–29 (DoD, testing, acceptance, CI), §§30 (deployment), §32 W015, §33 Work Item execution contract, §36 known risks, §37 tunings, §39 repository governance and release authorization, §41 Git/GitHub acceptance criteria; AC-001–026; NFR-001–030; SEC-001 onward; W015 roadmap entry.
- W015 validates approved requirements and feature specifications. It may report a requirement gap, but must not silently change confirmed behavior or lower its acceptance criteria.
- Initial browser targets are Chrome, Safari, Edge, and Firefox across desktop, mobile, and tablet web. Exact supported versions/device matrix is not specified and requires Owner approval before claiming compatibility.
- Initial measured targets include six players per Room, 60 FPS aim on capable devices, 320 CSS px mobile width, approximately 44×44 CSS px touch targets where practical, and latency scenarios of 0/100/200/300 ms. These are measurement targets, not guaranteed pass thresholds.
- Closed Alpha label/release name is a roadmap baseline; merge to `main`, production deployment, external invitations, and publication require separate authorized release approval under §29 and §41.

## Readiness gates

1. **Requirement and scope gate:** every MVP requirement has a linked approved spec, implementation, test/evidence, independent review, and status. Deviations are classified and authorized; no unresolved BLOCKER or MAJOR defect is silently accepted.
2. **Integration gate:** all prerequisite implementation PRs are merged to `develop` through the prescribed PR workflow, required CI passes on exact merge revisions, feature branches remain retained, and documentation/ADRs match integrated behavior.
3. **Functional gate:** required 2–6 player flows complete using real browser clients and server-authoritative state. No client-side score or eligibility is accepted as truth.
4. **Security gate:** review checklist is complete; actionable findings are severity-classified and fixed/re-reviewed or have explicit authorized risk acceptance where policy permits. No raw secrets or sensitive payloads appear in logs or shipped bundles.
5. **Compatibility/accessibility gate:** actual browser/device matrix and UX evidence are recorded. Unsupported/unavailable combinations are explicit; no blanket compatibility claim based on untested environments.
6. **Performance gate:** measurement reports include environment, method, raw/summary observations, and relation to initial targets. Owner reviews any material failure or target change.
7. **Operational readiness gate:** release configuration and recovery/rollback procedure are reviewed; health checks, HTTPS/WSS, asset availability, and in-memory restart limitation are understood. No deployment is performed by this Work Item without explicit authorization.
8. **Independent acceptance gate:** separate Reviewer inspects source revision, specs, test reports, defects, security/performance/browser evidence, and proposed limitations; PASS/BLOCKED/FAIL applies to exact candidate SHA.
9. **Owner decision gate:** Project Owner receives the readiness report and explicitly authorizes or declines release actions. A technical PASS alone is not release authorization.

## Required scenarios and evidence

### Functional multiplayer

- Two through six distinct browser contexts create/join private Rooms by code and invitation URL, show the same Lobby membership, exercise Ready/Host/start restrictions, and complete Speed Battle.
- Verify private boards remain private, accepted picks update authoritative scores once, incorrect/technical rejection behavior is correct, and shared protocol state is consistent.
- Exercise timer expiry, invalid-board recovery/cancellation, ties, all-forfeit, one eligible player, voluntary leave, unexpected disconnect/reconnect within/outside grace, Match end during grace, and authorized connection replacement.
- Verify Results/winner projection, Room-preserving Rematch, Ready reset, new Match ID, stale-action rejection, and Room cleanup.
- Validate server restart behavior as documented: active in-memory Rooms may be lost and clients receive a clear Thai recovery path rather than stale gameplay.

### Security and abuse

- Attempt unauthorized/cross-Room/session actions, session misuse, forged scores/timers, wrong Origin, malformed/oversized payloads, duplicate/replayed/stale requests, rate-limit bursts, Host privilege violations, and resource-exhaustion patterns.
- Inspect cookie flags and secret handling in the actual configured deployment path; verify no credentials/private boards/internal errors are exposed to peers or logs.
- Review third-party dependencies, approved project-local Skills, and all shipped asset source/license/provenance records; global Skill installation remains prohibited.

### Compatibility and UX

- Capture environment metadata and evidence for each tested browser/device class; include Thai glyph rendering, layout, input, sound/autoplay, reduced motion, errors/reconnect, and safe asset loading.
- Verify the full 4×4 board at 320 CSS px and representative larger sizes, with no unintended scrolling and practical touch targets; report actual dimensions and any limitation.
- Record keyboard-only and screen-reader-oriented semantic interaction checks; verify five card attributes are distinguishable without color alone.

### Performance and concurrency

- Measure six concurrent players in a Room and activity across multiple Rooms; verify one Room's slow client does not block other Rooms.
- Simulate documented latency values 0, 100, 200, and 300 ms, disconnects, duplicate/delayed messages, reconnects, and slow-client behavior.
- Record frame rate/long tasks, initial asset load, server action/update timing, CPU/memory/queue behavior, and environment. Do not assert universal 60 FPS or capacity beyond measured setup.

## Acceptance criteria

- **W015-AC-01:** Every MVP requirement has a traceable implementation and acceptance-evidence status; missing or deferred coverage is explicit, with no requirement silently treated as passed.
- **W015-AC-02:** End-to-end 2–6 player scenarios exercise create/join, Lobby/Ready/start, Speed Battle, score synchronization, terminal cases, Results, Rematch, reconnection, and restart recovery across independent browser contexts.
- **W015-AC-03:** Required unit/integration/E2E/concurrency/security/build checks pass for the exact release candidate or are explicitly NOT_RUN/BLOCKED/N/A with reason; no result is fabricated or inferred from unrelated CI.
- **W015-AC-04:** Browser/device matrix covers initial Chrome, Safari, Edge, Firefox and desktop/mobile/tablet classes to the extent available; exact version/environment and gaps are recorded.
- **W015-AC-05:** Security review covers SEC requirements and adversarial scenarios; no unresolved BLOCKER/MAJOR security findings remain unless explicitly handled through an authorized policy decision.
- **W015-AC-06:** Performance/network measurements are recorded against §21.1 initial targets and documented scenarios, including environment, data, method, observed limitations, and owner decisions where target changes are proposed.
- **W015-AC-07:** Visual/accessibility evidence includes all required Coco surfaces, Thai UI, 4×4 mobile readability, keyboard/touch, non-color cues, reduced motion, and audio controls; gaps are explicit and assigned.
- **W015-AC-08:** Production/release configuration, HTTPS/WSS, health checks, secrets, asset delivery, in-memory limitation, and rollback/recovery instructions are reviewed; no deployment or publication occurs without separate authorization.
- **W015-AC-09:** Independent release acceptance review PASS identifies the exact candidate SHA, evidence reviewed, limitations, and unresolved findings; material later changes trigger affected validation/re-review.
- **W015-AC-10:** Project Owner receives a concise readiness recommendation and explicitly approves or rejects release authorization; completion is not declared DONE before required review, documentation, PR integration, and evidence are complete.

## Technical and operational plan

- Build a requirement-to-evidence matrix from the source requirement IDs and approved W001–W014 specs. Each row records verification artifact/command, environment, run date, status (PASS/FAIL/NOT_RUN/BLOCKED/N/A), owner, and linked report.
- Use the existing pnpm/Go/Playwright/GitHub Actions toolchain rather than introducing a new test framework without a demonstrated gap.
- Use isolated temporary test Rooms and browser contexts; ensure cleanup and avoid using real player credentials or external personal data.
- Keep release candidate immutable while collecting evidence. Record commit SHA and CI run IDs for every result; rebuild/test after material source changes.
- File discovered defects as separate traceable Work Items; block readiness when a finding violates a confirmed requirement or required security gate.
- Produce a release readiness report with summary, tested revision, requirement coverage, commands/results, browser matrix, security/performance evidence, known risks, unresolved issues, recommendation, and owner decision record.
- Release process follows `develop → main` Pull Request with Merge Commit, retained branches per policy, independent release review, required CI, and explicit authorized merge/deploy approval. W015 does not perform unauthorized operations.

## Dependencies and approval requirements

- All required MVP foundations and implementation Work Items W001–W014, their approved specs, and integration changes must be complete enough for full-system validation. Missing feature contracts or unmerged implementation PRs block the associated validation claim.
- Approved W008–W013 sessions/Room/protocol/multiplayer/results/reconnection contracts; approved W014 Coco/assets/UI integrations; W004 CI and W005 protocol validation.
- Project Owner approval is required for exact browser support versions, representative performance environments/interpretation of targets, any acceptance exception/risk acceptance, final Closed Alpha go/no-go, release promotion, production deployment, or external tester invitations.
- External dependencies: access to representative Safari/mobile devices, production-like HTTPS/WSS environment, security review capability, and test accounts/rooms. Missing access is recorded as BLOCKED/NOT_RUN with mitigation; it is not represented as passing.

## Approval

- Project Owner / authorized approver: Pending
- Independent specification review: Pending
- Release authorization: Not requested by this specification review; separate approval required
- Decision and date: Pending
