# W012 — Results & Rematch

Status: READY_FOR_REVIEW
Type: Backend and Frontend Feature
Owner: Project Owner assigns during execution
Last updated: 2026-10-19

## Goals

- Publish final, server-calculated Speed Battle results consistent with authoritative Match scores and W013 eligibility/forfeit state.
- Display final scores, player names, rankings, joint winners/tie status, forfeited players, and Thai result copy.
- Return the existing Room to Lobby for a Rematch with a new Match ID and reset Match-specific state while preserving Room/member/Host/settings state.

## Non-goals

- Guest Sessions, Room state/Host permissions, WebSocket transport, or gameplay rules/scoring/timing (W008–W011).
- Reconnection grace, forfeit timing, host migration, or player eligibility decisions (W013).
- Classic Rush, Party Chaos, persistent Match history, public leaderboards, analytics, or cross-Room results.
- Changing winner/tie/rematch requirements or adding unapproved rematch voting/ready policies.

## Scope

- Consume a terminal authoritative Match snapshot from W011 with final scores and participant identities, plus eligibility/forfeit state from W013.
- Calculate eligible-player rankings and joint winners on the server; use shared placement for ties (for example, 1st, 1st, 3rd).
- Publish a stable result snapshot/event through W010 and render the Results screen in Thai with accessible desktop/mobile layout.
- Handle no-winner results for all-forfeit or cancelled Match outcomes, without inventing a winner from invalid state.
- Implement the approved Rematch transition from Results to the same Room's Lobby, preserving Room identity/code, connected members/current Host, and selected settings unless explicitly changed.
- Reset Match ID, score, combo, boards, references, timer, event schedule, and other Match-specific data; set every participant Not Ready for the next Match.
- Ensure duplicate/late completion or Rematch requests are idempotent and do not create multiple Matches or mutate completed results.

## Requirement references

- `docs/REQUIREMENT.md` §§3.1 (MVP Results/Rematch), §§7.5–7.9 (Room state, Host/readiness), §§9.10–9.11 (end conditions/winners), §§12.1 and 12.5–12.8 (Match state, completion, Results, Rematch), §§13.2–13.7 (eligibility/forfeit interactions), §§14.2–14.3 (server authority/idempotency), §15 EC-034–040 and EC-050, §§16.7 (FR-074–076), §17 INV-010/011/014/015/016/018, §§18–20 (server authority, transport, security), §§21–25 (invariants, tests, UX/security), §28 AC-020–022, §32 W012.
- W011 Multiplayer Speed Battle (specification in PR #13): provides final authoritative score/state and completion/cancellation reason; not yet owner-approved/merged.
- W013 Reconnection & Anti-Griefing: provides final eligibility/forfeit status; not yet specified/approved/implemented.
- W009/W010 provide Room Lobby transitions and transport events; proposals remain in open PRs #11–12 and are not treated as approved contracts.
- Status fidelity: highest eligible score wins; ties are allowed and use shared placement (§9.11, BASELINE Speed Battle). Rematch reset/preserve rules in §12.8 are confirmed behavior within the BASELINE lifecycle. Which button/request triggers a Room-wide return to Lobby and how forfeited players are ranked/displayed are not fully specified and remain explicit decisions below.

## Expected behavior

1. When W011 marks a Match COMPLETED, W012 consumes one immutable final snapshot with Match ID, Room ID, each player's final score/name/opaque Player ID, and W013 eligibility/forfeit status. The server validates that it matches the active Room/Match and finalizes exactly once.
2. W012 ranks eligible players by descending server-recorded score. Players tied at the highest score are joint winners; shared placement uses competition ranking (e.g. 1st, 1st, 3rd). The client never calculates or submits winners.
3. Results display final scores, player names, rankings, winners, tie status when applicable, forfeited participants, and an available Rematch action, all in Thai. Usernames render as safe text. Results projections contain no credentials, private boards, PRNG state, or internal errors.
4. A CANCELLED Match caused by unrecoverable consistency error has no winner. All players forfeiting yields no winner; one remaining eligible competitor can be the winner after W013 resolves the grace/forfeit policy. W012 follows the server-provided terminal reason and eligibility rather than recomputing disconnect policy.
5. A Rematch returns the same participants/Room to LOBBY (subject to the approved trigger semantics), creates a new Match ID only when the next Match is actually started, resets all Match state, preserves Room identity/code/members/current Host/settings, and sets every player Not Ready.
6. The next start follows W009's normal Host/readiness/start checks and W011's preparation; no shortcut bypasses Ready, capacity, eligibility, or settings validation.
7. Duplicate terminal events do not republish divergent Results; a Rematch request repeated with the same logical request does not reset or create duplicate state. A stale action referencing a previous Match cannot affect the next one.
8. Results remain available while the Room is in RESULTS until the approved Rematch/Room-close path occurs. Process restart may lose the in-memory Room/Match/Results; no persistence is implied.

## Business rules

- Ranking source is immutable server-finalized score plus eligibility. Client score claims are ignored. W011 owns gameplay score; W013 owns eligibility/forfeit; W012 owns ranking/winner calculation and result projection.
- Highest score among eligible participants wins. Ties at the highest eligible score produce joint winners. Shared placement follows the documented example (1,1,3), not dense ranking (1,1,2).
- Forfeited players cannot win. Results must indicate forfeiture (§12.7). Whether they receive a numeric ranking or appear in a separate forfeited section is unspecified. Recommendation: show final score with a Thai forfeited label in a separate non-ranked section; owner approval required before implementation.
- A Match cancelled for consistency failure has no winner even if scores exist. No winner may be inferred from incomplete/corrupted score data.
- Results entry is idempotent and terminal. No score/result mutations after completion (INV-010/018); Host has no result-edit permissions (INV-011).
- Rematch preserves the existing Room, Room Code, connected members, current Host, and selected settings unless the Host changes permitted settings. It resets scores, combo, boards, references, Match timer, event schedule, Match ID, cooldowns, and all Match-specific state. All participants must become Not Ready and ready again before the next start.
- §12.8 says a Rematch returns players to Lobby, but does not specify who initiates it or whether every player must consent. **PROPOSED recommendation:** any current Room participant may request a single Room-wide return to Lobby; the request is not a vote to start. The Host still starts the next Match after all readiness checks. Require owner approval before implementation.
- The result screen's player-visible copy is Thai. Shared placement and tie copy must be clear for Thai-first players and use accessible semantic labels.

## Edge cases and error handling

- Duplicate, stale, or mismatched completion event (wrong Room/Match ID, already finalized, partial participant list): ignore/reject as technical error and preserve existing terminal snapshot; do not calculate alternate winners.
- Two Match completion triggers race (timer and last eligible player): W011/Room serialization emits one terminal state; W012 uses compare-and-transition/idempotency so result finalization occurs once.
- Equal top score: all eligible tied players are joint winners with equal rank; following ranks skip tied positions (1,1,3).
- All players forfeited: display no winner and all forfeited statuses. One eligible player after W013 resolution: that player wins, even if Match ended early for eligibility reasons.
- Player in reconnect grace remains eligible until W013 resolves; do not declare early winner while a competitor can reconnect under the applicable policy.
- Cancelled Match: no winner and a Thai recoverable result/cancellation message; no rematch state may reuse the cancelled Match ID.
- Rematch request from an unauthorized session, nonmember, or invalid Room state is rejected; no partial reset. Exact caller eligibility for the Room-wide transition follows the owner-approved trigger policy.
- Player disconnects/leaves during Results: Result snapshot remains immutable; connected members and Host preservation follow W009/W013 contracts. Do not change final eligibility retroactively after Results.
- Settings changed in Lobby after Rematch reset readiness again as required by W009. Ready actions during RESULTS are rejected.
- Older Match action/result event arrives after new Match starts: reject/ignore by Match ID and Room state; no stale score leaks into the new Match.
- Room expires/closes during Results: W009 lifecycle owns cleanup; W012 releases result references on Room close.
- Username missing or malformed in final projection: use the server-approved safe display fallback or fail the projection safely; never render untrusted markup. Exact fallback should follow W008/W009 UI contract.

## Acceptance criteria

- **W012-AC-01:** A valid terminal W011 Match snapshot plus W013 eligibility state is finalized once by the server and bound to the correct Room/Match ID; duplicate/stale/partial snapshots cannot change the result.
- **W012-AC-02:** Server ranks eligible players by final score, declares all top-score ties joint winners, uses shared placement (1,1,3), excludes forfeited players from winning, and produces no winner for cancelled/all-forfeit results.
- **W012-AC-03:** Results projection contains final scores, player names, rankings, winners, tie status, forfeited status, and no private boards, credentials, PRNG state, or internal diagnostics.
- **W012-AC-04:** Thai Results UI renders all required data safely and accessibly on desktop/mobile, including no-winner, tie, forfeited, loading, and error states.
- **W012-AC-05:** Approved Rematch trigger returns the existing Room to LOBBY without creating a Match; it preserves Room ID/code/members/current Host/settings, clears all Match-specific data, and marks every participant Not Ready.
- **W012-AC-06:** A subsequent Match requires a new Match ID and normal W009/W011 Host, readiness, settings, and preparation gates; repeated/stale Rematch commands cannot create duplicate Match state.
- **W012-AC-07:** W011/W013/W009/W010 integration contracts are followed: gameplay scores remain immutable after completion, eligibility is not recalculated by the UI, Room ownership stays with W009, and Room closure releases Results.
- **W012-AC-08:** Tests cover score ordering, shared ties, 1,1,3 ranking, forfeits, all-forfeit, one eligible, cancelled Match, reconnect grace boundary, duplicate/stale result events, rematch reset/preservation, stale previous-Match actions, Thai rendering, accessibility, and Room-close cleanup.
- **W012-AC-09:** Independent implementation review PASS applies to exact final PR SHA, required CI/security checks pass, no BLOCKER/MAJOR findings remain, docs synchronize, PR merges to `develop` with Merge Commit, branch retained, and evidence recorded.

## Technical design

- Domain boundary: `apps/server/internal/results/` (or an equivalent small package) consumes immutable W011 final Match snapshot plus W013 eligibility and produces deterministic result projection. It must not reach into mutable game state or decide reconnect/forfeit policy.
- Match/Room orchestration: Room remains in RESULTS after W011 terminal transition. W009 owns Room membership/Host/settings and the Room state transition back to LOBBY. W012 calls an explicit approved Room command and does not mutate Room fields directly.
- Ranking: sort eligible participants by score descending with stable deterministic tie order for display only; assign competition rank `1 + number of eligible players with a strictly greater score`. All with max eligible score are joint winners. Forfeited players are represented separately per the approved display/ranking policy.
- Result snapshot includes result version, Room ID, Match ID, completed/cancelled reason, all participant scores and eligibility/forfeit labels, eligible ranking, winners, tie flag, and finalization time. It is immutable once published. Include only approved public data.
- Rematch transition is a state command to W009; W012 validates the approved initiator rule and current RESULTS state, then requests `returnToLobbyAndResetReady`. W009 preserves Room state and members while W011 creates fresh Match-specific state only on next valid start. Idempotency uses W010 request IDs and Match/Room state versions.
- Protocol: extend W005/W010 strict v1 contract with `match.results` and the approved Rematch request/result/event. Include fixtures for complete, tied, forfeited, all-forfeit, cancelled, and reset states; don't send winner claims from clients. Review whether state fields can be added compatibly before implementation.
- Frontend: React Results feature consumes typed result projection; Phaser gameplay scene unmounts or freezes on terminal event and clears old Match display state after transition. Semantic rank/winner/tie/forfeit labels and Thai text are accessible. A Rematch request shows pending/error state and returns to Lobby only after authoritative Room update.
- In-memory only: result state lifetime equals Room/Match lifetime. On process restart, show the documented recoverable state-loss path rather than reconstructing false results.
- W011/W013/W009/W010 remain pending dependencies. If their approved shapes differ, revise this spec and re-review before implementation.

## Security and privacy

- Authenticate each Results/Rematch action through W008 and active W010 connection; verify Room membership and approved permission via W009.
- Results and rankings are generated from server-owned final state. Ignore client scores, rankings, winner selection, eligibility, and Match IDs unless verified against current authority.
- Never expose session credentials, private boards, random seeds, internal errors, or another Room's data in result messages/logs.
- Render usernames as escaped text under the approved W008 Unicode safety policy. Do not insert username/Room Code content as HTML.
- Ensure Match IDs and Room revisions prevent replay of old result/rematch requests. Idempotency cannot change a finalized result or reset a different/newer Match.
- Logging records result finalization/cancellation and stable categories with opaque IDs; avoid cookies, private content, unnecessary personal data, and high-cardinality metrics.

## Test plan

- Go unit tests for deterministic ranking, top ties, shared placement, forfeited exclusion, no eligible winner, cancellation/no-winner, immutable projections, and malformed/inconsistent terminal input.
- Integration tests for W011 completion event exactly once, W013 eligibility changes before terminal state, active grace preventing premature winner, Results state transition, W009 Room reset preservation, W010 duplicate/rematch request handling, and Room closure cleanup.
- Frontend component tests for Thai rank/tie/winner/forfeit/no-winner/error/loading UI, safe username rendering, accessible semantics, and authoritative navigation only after server response.
- Browser E2E with 2–6 clients covers timer end, score/result consistency, ties, forfeit, Results display, Room state return, readiness reset, starting a new Match with new ID and reset state, plus mobile viewport.
- Run `pnpm check`, protocol fixtures, Go tests/race/vet, CI/security checks, and record actual results per criterion. Browser/E2E environment limitations must be NOT_RUN/BLOCKED with evidence, not implied pass.

## Dependencies and risks

- W011 authoritative Match completion/score snapshot; W013 final eligibility/forfeit outcomes; W009 Room Results→Lobby/reset command; W010 transport/idempotency/protocol. All are pending or proposed at this point.
- Risks: result event races; eligibility changes near timer end; inconsistent or partial snapshots; tie/forfeit interpretation; user consent/trigger ambiguity for Room-wide Rematch; protocol v1 evolution; preserving selected TUNING settings; and lost in-memory Results on process restart.
- Decisions requiring Project Owner approval before implementation: Rematch trigger/consent model; ranking/display for forfeited participants; and any protocol or result field changes that materially alter player-visible behavior. Preserve confirmed tie/winner/reset semantics.

## Approval

- Project Owner / authorized approver: Pending
- Independent specification review: Pending
- Decision and date: Pending
- W009/W010/W011/W013 dependency approval/merge: Pending
