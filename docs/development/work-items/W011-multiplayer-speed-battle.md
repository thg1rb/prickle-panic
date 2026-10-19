# W011 — Multiplayer Speed Battle

Status: READY_FOR_REVIEW
Type: Server-authoritative Gameplay Feature
Owner: Project Owner assigns during execution
Last updated: 2026-10-19

## Goals

- Replace the local-only W007 gameplay authority with a fully playable, server-authoritative Speed Battle for 2–6 Room participants.
- Run authoritative Match lifecycle, private boards, card selection validation, score/combo state, cooldown, and timer through Go game/match modules serialized by the owning Room.
- Synchronize responsive Thai gameplay UI through the approved W010 protocol while preserving private per-player state and deterministic testability.

## Non-goals

- Guest Session, Room/Lobby rules, WebSocket transport, or protocol envelope implementation (W008–W010).
- Final winner calculation, Results screen/lifecycle, return to Lobby, or Rematch (W012).
- Reconnection grace periods, forfeit policy, host migration after unexpected disconnect, or broader anti-griefing policy (W013); W011 exposes integration hooks and consumes authoritative eligibility/presence events.
- Classic Rush or Party Chaos; both remain post-MVP after Speed Battle and Closed Alpha readiness.
- New matching/scoring/tuning rules, client-authoritative score/timer, persistent Matches, or multi-server gameplay.

## Scope

- Coordinate Match preparation, 3-second initial countdown, IN_PROGRESS state, authoritative Match timer, completion/cancellation signals, and bounded state cleanup.
- Initialize 2–6 players with zero score/combo, identical starting reference and board configuration, and independent board states using W006's pure Go engine.
- Process authenticated `game.pick_card` intentions: validate active Match/RoomPlayer/card identity/versions/eligibility/cooldown and apply exactly W006/REQUIREMENT behavior.
- Synchronize private board/reference updates and public score updates through W010's approved recipient event streams; support snapshots and state resynchronization.
- Integrate W007's React–Phaser gameplay presentation with server-owned state, replacing or isolating its temporary local adapter from the online path.
- Expose typed completion/forfeit/disconnect hooks to W012/W013 without duplicating their Results or reconnection policy.

## Requirement references

- `docs/REQUIREMENT.md` §§3.1–3.2 (MVP/post-MVP mode scope), §§8.1–8.10 (card rules), §§9.1–9.12 (Speed Battle), §§12.1–12.6 (Match lifecycle), §§13.1–13.9 (disconnect/forfeit integration), §§14.1–14.5 and 14.8–14.9 (authority, idempotency, limits, bot/logging boundaries), §15 (EC-016–048), §§16.2–16.7 (FR-024–076), §17 (INV-001–018), §§18.6–18.11 (modules, authority, event loop, in-memory state), §§19.2–19.12 (protocol, snapshots, versions, timing, slow clients), §20 (SEC-001–015), §§21–25 (testing, operations, security, concurrency), §28 (AC-006–010, AC-023–024), §32 (W011).
- W005 protocol: `docs/protocol/v1/README.md` and schema/fixtures (Speed Battle action intent, per-player game snapshots, correct/incorrect action results, public scores, errors; recipient-stream sequencing proposal is under W010 review).
- W006 engine: `docs/development/work-items/W006-card-matching-engine.md` (DONE; Hamming distance ≤ 1, score/combo transitions, seeded solvable board generator).
- W007 prototype: `docs/development/work-items/W007-local-gameplay-prototype.md` (DONE; local-only UI adapter and React–Phaser bridge to replace/isolate).
- W008/W009/W010 contracts are currently proposed in open PRs #10–12 and await owner decisions/review/merge. W011 must consume only their approved interfaces; no policy stated in a pending spec is treated as confirmed.
- Status fidelity: Speed Battle behavior is BASELINE (§40); board size 4×4/16 cards, 60-second timer, +1 correct pick, zero incorrect score, 0.8-second cooldown, visual-only combo, continuous refill, and highest-score winner are configurable TUNING values (§9.2/§37). Countdown default is 3 seconds (§12.3; retain its documented status). Do not change them without playtesting/approval. One unresolved pick per player and the W006 random generation fallback are PROPOSED technical details, not product rules.

## Expected behavior

1. A Host starts Speed Battle only through the authorized Room/Match path after W009 readiness validation. Match preparation validates 2–6 eligible participants, approved settings, Room/Match state, and board solvability before countdown.
2. The server creates a unique Match ID, assigns zero score and combo to every participant, generates one shared starting reference/board configuration, and copies independent board state to each player. Clients receive only their own board/reference plus public scores.
3. During the 3-second countdown, no gameplay action is accepted. At the authoritative start instant the Match enters IN_PROGRESS and server time begins; client clocks are display-only.
4. A correct pick is exactly W006's `HammingDistance(reference, selected) <= 1`, including identical cards. The server awards the configured +1, increments visual combo, changes that player's reference, removes/replaces only that player's selected instance in the same slot, preserves a 16-card board and at least one valid choice, then sends a private board update and appropriate public score update.
5. An actual incorrect pick awards zero, resets that player's combo, leaves board/reference/score unchanged, and applies the configured non-extending 0.8-second mistake cooldown. Technical rejection (stale state, invalid ID, cooldown, unauthorized, duplicate) does not cause an incorrect-pick penalty or reset/extend cooldown.
6. A player's correct pick never changes another player's private board/reference. Cards may have the same variation across players; Card Instance IDs are player-board scoped and cannot be confused across owners.
7. The server processes simultaneous actions deterministically through the owning Room loop; one logical request cannot mutate twice. Actions after Match end, before start, during countdown, from an ineligible/forfeited player, or with stale/invalid state are rejected safely.
8. An AFK participant does not pause the timer or other boards. A player may resume while still eligible and the Match remains active.
9. Match ends at configured duration, when W013 reports only one eligible competitor remains after forfeit/grace resolution, or when authoritative state cannot be safely maintained. W011 freezes further gameplay mutations and emits final authoritative Match state or a no-winner cancellation reason; W012 owns winner/Results presentation and return/Rematch behavior.
10. A player in W013's reconnect grace remains eligible until the agreed timeout or normal Match-end policy resolves them; do not declare an early forfeit victory while an opponent remains eligible to reconnect.
11. An impossible board triggers safe server correction without player penalty; if recovery fails, cancel the Match without winner. Anomalies are recorded without sensitive payloads.
12. The online UI renders server state and Thai feedback. Phaser emits selection intentions only; React/application state and the server protocol determine displayed score, board, cooldown, and remaining time.

## Business rules

- Match state values are PREPARING, COUNTDOWN, IN_PROGRESS, COMPLETED, and CANCELLED (§12.1); server owns all transitions. Room state is coordinated with W009 and may be COUNTDOWN/IN_PROGRESS/RESULTS as documented.
- Match preparation follows §12.2 order: validate Room/player/readiness/settings, create Match ID, generate/validate initial state, schedule authoritative start. No action is accepted unless Match is IN_PROGRESS.
- Preserve §9 BASELINE rules and all §9.2 TUNING defaults. Combo remains visual only; it never multiplies or changes score. Incorrect pick is a processed gameplay outcome; technical errors are separate per FR-076/W005.
- Each player's board/reference/score/combo/cooldown is private state owned by the server. Public scores are shared. The client cannot choose board seed, card replacement, score, timer, eligibility, or outcome.
- Use W006 pure transitions for matching/scoring/generation. W011 owns clock/deadline, cooldown window, eligibility, request orchestration, event projection, and Match lifecycle integration. Do not copy matching rules into React/Phaser.
- Server obtains an unpredictable match seed from a cryptographic source and keeps it private; pass deterministic seed/state to W006. Do not reveal PRNG state/future card sequence. Same initial conditions are required; difficulty need not remain equal after choices.
- Tuning values come from approved configuration, validated at preparation and frozen when countdown starts. Host cannot change settings after countdown begins. A settings change before start follows W009 and resets readiness.
- Action ordering is authoritative through the Room event loop. **PROPOSED fairness rule for T01:** accept a pick only if the Gateway's trusted server receipt time is strictly before the monotonic deadline; at/after deadline reject. Never trust client timestamp. Record how queue delay and timer event ordering are handled and confirm with owner/design review because the requirement does not define the exact deadline-boundary rule.
- **PROPOSED input policy:** allow at most one unresolved pick per player initially, consistent with §19.9's initial option, until latency tests justify bounded queuing. This limit is not a gameplay penalty and does not extend cooldown.
- Duplicate request IDs use W010's approved bounded idempotency semantics (same request/payload cannot score twice; same ID/different payload returns technical duplicate error). Card state versions are checked against the player's board/reference; stale or removed card is a technical rejection, never an incorrect-pick result.
- W013 owns reconnection grace, forfeit, and player eligibility transitions. W011 accepts idempotent eligibility updates and keeps Match state while a disconnected player remains eligible. Voluntary gameplay leave is immediate forfeit under §13.6; unexpected disconnect is not immediate forfeit. W012 owns winner/tie calculation and Results, using W011's authoritative scores plus W013 eligibility.
- Match completion is atomic with respect to queued picks: an action is either processed before the authoritative end and included in final state, or rejected after completion; it cannot mutate a completed Match (INV-010).

## Edge cases and error handling

- Match start with 1 or 7 participants, any non-Ready/disconnected participant, invalid setting, duplicate RoomPlayer, existing Match, or invalid initial board: reject/cancel preparation without partial Match state.
- Card ID belongs to another player's board, is unknown, removed/replaced, or has stale board/reference version: return stable technical error; no score/combo/cooldown/board mutation.
- Identical variation/card is valid under distance-zero baseline; 2+ attribute difference is an incorrect gameplay result only when the selected instance is valid/current/owned and the request is otherwise eligible.
- Duplicate request while pending/completed cannot score twice. Reused request ID with different payload is rejected. If dedupe entry expires, retention must still comply with the approved W010 policy; tests cover retries around retention boundary.
- Two requests for same card/player arrive concurrently: Room loop processes at most one current instance; the later stale/duplicate one cannot award again. Different players' independent boards can process concurrently only through deterministic Room serialization.
- Exact timer boundary uses the approved trusted receipt/deadline rule. Timer callback and pick race produces one serial outcome; no score is awarded after completion.
- Wrong pick just before cooldown expiry, repeated wrong pick during cooldown, cooldown state after disconnect/reconnect, and Match end during cooldown are handled by authoritative clock; no repeated error extends cooldown.
- AFK player's clock continues. All-player forfeit produces no winner; one remaining eligible competitor can end Match under §9.10; a player in grace prevents premature victory. W013 provides the authoritative eligibility transition.
- Recipient sees only own board updates. Private event sequencing follows the owner-approved W010/W005 recipient stream; private events to opponents do not leak content or create artificial sequence gaps.
- Initial board/replacement cannot produce no-valid-choice state. If an invariant check fails, attempt deterministic/secure correction and publish a private correction; if impossible, cancel all participants consistently with no winner.
- Server restart destroys in-memory Match state; client receives recoverable recovery error and cannot resume stale IDs. No false claim of durable Match recovery.
- Malformed, oversized, unauthorized, stale, wrong-match, or after-end messages never panic or mutate state and never receive a gameplay penalty.

## Acceptance criteria

- **W011-AC-01:** Host-authorized Match preparation enforces Room/participant/readiness/settings preconditions and creates one unique Match ID with PREPARING → COUNTDOWN → IN_PROGRESS transitions; duplicate start does not create a second Match.
- **W011-AC-02:** 2–6 participants receive zero score/combo, identical initial reference/board configuration, and separate independent private board state; the server uses W006 and keeps seed/future sequence private.
- **W011-AC-03:** All matching/scoring/refill behavior exactly matches §9 and W006, including distance ≤1, +1 correct, zero incorrect, visual-only combo, same-slot replacement, board size and solvability, with no cross-player board mutation.
- **W011-AC-04:** Actual incorrect pick applies the configured non-extending cooldown and changes combo while preserving score, board, and reference; technical errors (authorization, duplicate, stale, cooldown, invalid card, wrong state) have no gameplay penalty or partial mutation.
- **W011-AC-05:** Server owns Match timer/deadline and all score/card outcomes. Client-supplied score/clock/seed is ignored or rejected. At timer boundary, exactly one authoritative serial result is applied; completed/cancelled Matches accept no mutations.
- **W011-AC-06:** W010 protocol emits private per-player snapshots/board updates and public scores/results-to-W012 projections with approved sequence/correlation semantics; no opponent board or future generation state leaks.
- **W011-AC-07:** AFK behavior, W013 disconnect/grace/eligibility and voluntary leave/forfeit seams, and W012 completion/winner handoff are contract-tested using fakes without duplicating their ownership or changing policy; full cross-Work-Item E2E is reserved for W015.
- **W011-AC-08:** Invalid board recovery is safe and unpenalized; unrecoverable state cancels without a winner and records a diagnosable anomaly.
- **W011-AC-09:** React/Phaser online gameplay displays authoritative state in Thai, renders immediate feedback without claiming unconfirmed score, and the W007 local adapter is removed from or explicitly isolated from the online path.
- **W011-AC-10:** Tests cover gameplay boundaries, all match lifecycle states, request idempotency, stale versions, cooldown, timer boundary races, concurrent players/picks, private projection, deterministic random reproduction, disconnect eligibility and board recovery. Go race tests and frontend/browser E2E with 2–6 clients pass and actual results are recorded.
- **W011-AC-11:** Independent review PASS applies to exact final implementation SHA, required CI/security checks pass, no BLOCKER/MAJOR findings remain, documentation is synchronized, PR merges to `develop` with Merge Commit, branch retained, and evidence recorded.

## Technical design

- Domain packages: keep W006's pure `apps/server/internal/game/` engine unchanged for rules. Add a Match coordinator package (candidate `apps/server/internal/match/`) that composes W006, Match state/timer/eligibility, and emits typed transitions. No HTTP/WebSocket/browser imports in the core coordinator.
- State ownership: the owning W009 Room event loop serializes Match start, pick actions, timer/deadline events, presence/eligibility updates, and completion. Match state is owned by one Room and is not concurrently mutated outside that loop. Different Rooms run independently.
- Match aggregate: Match ID, mode `speed_battle`, phase, prepared/start/deadline timestamps, immutable approved settings, per-player eligibility/status, each player's independent W006 game state, cooldown deadline, request outcome associations, and final authoritative snapshot for W012. Keep player board fields private in projections.
- Randomness: generate a cryptographically random seed per Match (or secure per-board sub-seeds from it) using the standard library; keep it internal. Feed deterministic RNG/state into W006. Tests inject fixed seeds/clocks. Review the W006 fallback's distribution behavior; do not alter its matching or TUNING behavior.
- Clock: inject a monotonic clock/timer interface. Schedule countdown and end events; do not tick per frame. On any operation, compare authoritative server time against deadlines. Persist only in process memory, clean timers on completion/cancel/shutdown.
- Action pipeline: W010 authenticates connection/session and validates frame/request ID. Match coordinator then checks active connection generation, Room membership/MatchPlayer eligibility, state, rate/cooldown, Match ID, card ID ownership/currentness, and versions. It applies a W006 transition exactly once and emits typed private/public event projections.
- Cooldown: server stores per-player monotonic `cooldownUntil`; actual incorrect picks set to `processedAt + configuredCooldown` only when not already cooling down. Technical cooldown rejection does not extend it. A disconnect/reconnect per W013 policy preserves remaining cooldown if Match remains active.
- Timer/deadline race: Room loop assigns order using trusted server arrival/receipt time plus serial processing. The strict `< deadline` rule remains a proposal pending T01 approval. At finalization, drain or reject queued pre-deadline actions according to that approved policy; never process after completion. Avoid wall-clock jumps by monotonic elapsed-time source.
- Event projections: private board/reference updates and initial game snapshot go only to their owner; score updates contain only public per-player scores. Use W010's approved recipient stream. W012 receives an internal authoritative final snapshot (scores and eligibility) and controls public result/winner state. No client can change eligibility or scores.
- Frontend: replace W007's local state transition as online authority with `game.pick_card` intentions and server snapshot/events. Keep local pressed/hover animation as presentation only. React store owns connection and authoritative display projection; Phaser canvas receives snapshots/deltas and emits card instance IDs. Support keyboard/accessible React fallback and reduced motion.
- W011 consumes W008/W009/W010 APIs as dependencies; if those approved interfaces differ from pending specifications, update the design and obtain review before implementation. No new protocol field is assumed until W005/W010 contracts include matching fixtures.

## Security and privacy

- Require authenticated W008 Guest Session and current W010 connection generation for every action; validate Room and Match membership server-side. Never trust client player, score, reference, timer, seed, eligibility, or card replacement claims.
- Enforce server-side message/frame/queue/rate limits from W010; separately enforce gameplay mistake cooldown. Malformed and technical failures do not become incorrect gameplay penalties.
- Use W006 matching exactly; validate card instance ownership and current board/reference versions to prevent cross-player selection, stale action, or duplicate scoring.
- Private boards, reference cards, seeds, PRNG state, cooldown internals, and other players' action details remain private. Public score projection includes only approved public fields.
- Prevent timing races and duplicate mutations via Room-owned event loop, idempotency, monotonic deadline, and atomic Match terminal state.
- Logs contain no session cookies, credentials, raw websocket frames, private board data, or unnecessary usernames. Record stable error category, Match/Room opaque IDs, and anomaly metrics with bounded cardinality.
- Fail closed on unrecoverable Match state. No client is allowed to choose a winner or modify results.

## Test plan

- Pure/domain tests integrate W006: correct/incorrect transitions, all attribute difference boundaries, identical card validity, board identity/refill/solvability, score/combo, invalid action atomicity, and multiple deterministic seeds.
- Match lifecycle tests for prepare validation, unique ID, countdown, exact start/end, completion/cancel, no input outside IN_PROGRESS, idempotent start, all-eligible/one-eligible/all-forfeit states, and shutdown cleanup using fake clock.
- Concurrency/race tests for simultaneous actions on same/different players, timer-versus-pick boundary, replacement plus stale request, duplicate IDs and payload mismatch, state snapshot under serialization, and independent Room progress.
- Security tests for forged score/clock/player/seed, unauthorized session/generation, wrong Match/Room, cross-player card IDs, oversized/malformed payloads, rate limits, and logs/projections without secrets/private state.
- Protocol fixtures/browser E2E: 2–6 clients create/join, ready/start, countdown, correct/incorrect picks, score synchronization, independent boards, timer completion and state resync. Verify invitation/Thai UI and physical mouse/touch/keyboard through browser where available. Results/Rematch continue in W012.
- Run Go race tests, `go test`, `go vet`, `pnpm protocol:check`, `pnpm check`, and CI. Record actual results per acceptance criterion; do not treat CI from a different SHA as evidence.

## Dependencies and risks

- W006 Card Matching Engine (DONE) and W005 Protocol Contracts (DONE).
- W007 Local Gameplay Prototype (DONE); online path must replace/isolate local authority.
- W008 Guest Session, W009 Room Management, and W010 WebSocket Foundation are required and currently under separate review/owner gates (PRs #10–12); W011 must not begin before approved, merged interfaces exist.
- W012 Results/Rematch consumes final state; W013 provides grace/forfeit/eligibility/host migration policy. W011 can implement adapters/fakes but final MVP integration requires those contracts.
- Risks: timer boundary fairness and event-loop queue delay; identical initial difficulty but evolving boards differ; random sequence confidentiality; duplicate action retention; exact Room event protocol/private recipient sequencing; reconnect eligibility races; server restart erases state; Phaser local prototype may leak into online authority; full 2–6 player mobile performance.
- Decisions needing Project Owner approval during T01: deadline boundary/queue fairness rule, initial one-unresolved-pick policy if adopted, any change to TUNING defaults, random seed derivation if materially changing W006 design, and approved W008/W009/W010 policy/contracts.

## Approval

- Project Owner / authorized approver: Pending
- Independent specification review: PASS at `af31e5ec939946edcafeb97275fd8b2f99ddc366`, recorded in `docs/development/reviews/W011-review.md`
- Decision and date: Pending
- W008/W009/W010 dependency approval/merge: Pending
