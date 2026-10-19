# W013 — Reconnection & Anti-Griefing

Status: READY_FOR_REVIEW  
Type: Backend Feature / Security Hardening  
Owner: Project Owner assigns during execution  
Last updated: 2026-10-19

## Goals

- Restore an authorized Guest Session's active Match state after an unexpected disconnect within the confirmed grace period.
- Apply confirmed voluntary-leave and grace-expiration forfeits without affecting unrelated players or authoritative Room processing.
- Define safe Host continuity and abuse controls without letting a Host alter active Match outcomes.
- Enforce bounded, rate-limited, idempotent actions and clear recovery behavior for transient failures and server restart loss.

## Non-goals

- Persistent player identity, login, permanent bans, durable Match recovery, multi-instance routing, or perfect bot prevention.
- Changing Speed Battle scoring, card rules, Match timer, or the confirmed 15-second reconnection grace.
- Results ranking/rematch behavior (W012), core Room lifecycle/permissions (W009), Guest Session credential issuance (W008), or WebSocket transport implementation (W010), except through their approved interfaces.

## Scope

- Unexpected disconnect detection and 15-second grace timer using server time.
- Authorized same-session reconnection, connection replacement, and restoration of server-owned player state.
- Immediate voluntary forfeit, grace-expiry forfeit, one-eligible-player termination, all-forfeit no-winner outcome, and boundary handling when the normal Match deadline occurs during grace.
- Host continuity/migration according to an approved deterministic Room policy; preserve separation between Host permissions and Match authority.
- Domain abuse policies for gameplay and Room actions, with gameplay cooldown separated from security throttling. W010 owns transport/message-rate limits and malformed-message handling.
- Consume W010's bounded-queue, slow-client, connection-replacement, and transport rate-limit outcomes; apply W013 reconnect/eligibility policy when W010 closes a persistently slow or abusive connection.
- Safe cleanup of W013-owned grace timers and per-player eligibility state.
- Integrate with W009's Room lock, Room Code, Lobby-only kick, and Room/session abuse controls through explicit contracts; do not implement duplicate Room management.
- Safe logs and recoverable client messaging after server process loss.

## Requirement references and status

- `docs/REQUIREMENT.md` §§7.9–7.15 (Host permissions and migration), §§13.1–13.9: confirmed disconnect/reconnect/forfeit behavior, including 15-second grace, Match deadline behavior, voluntary leave, ownership proof, connection replacement, and in-memory restart limitation.
- §§14.1–14.9: confirmed server authority, action idempotency, suggested rate limits, gameplay/security cooldown separation, Host restrictions, Room abuse protections, bot limitation, safe logging.
- §15 EC-009–012, EC-030–036, EC-040–047; §16 FR-067–073; §17 INV-004, INV-010–017; §§18.8–18.11, 19.10–19.12, 20 SEC-001–006; §21 NFR-004–010/024/025, performance targets; §28 AC-017–019/023–025; §32 W013.
- Confirmed behavior must remain unchanged: unexpected loss is not immediate forfeit; grace is 15 seconds; Match continues; authorized reconnect restores server state and score without time refund; expired grace forfeits while Match active; normal deadline does not wait out remaining grace and uses server score while player remains eligible; explicit gameplay leave immediately forfeits and bars re-entry; username alone is not reconnect proof; new authorized connection invalidates old gameplay authority; all forfeit means no winner; one eligible player means Match ends with winner; in-memory server restart can lose state and UI reports recovery error.
- Rate limits in §14.4 are initial suggestions subject to testing, not fixed acceptance thresholds. Preserve their proposed values in configuration pending validation; do not silently treat them as confirmed production values.
- Voluntary Host departure is confirmed by §7.13: transfer Host ownership to the eligible remaining player with earliest join order. For unexpected Host disconnection, §7.13 confirms grace before transfer; successor selection is not explicit and remains a proposed Owner decision. Host powers never include Match scoring/result edits or active-Match kick/cancel.
- W008 Guest Session and W009 Room Management contracts are prerequisites. W010 owns transport framing/connection mechanics. W011 owns Match serialization, action execution, and terminal completion. W012 consumes terminal eligibility state. W013 must not duplicate these responsibilities.

## Expected behavior

1. On unexpected connection loss during an active Match, keep the player's authorized server-side Match state and mark the player temporarily disconnected. Start one server-time deadline exactly 15 seconds after loss; do not pause other players or refund time.
2. A reconnect must authenticate ownership of the original Guest Session through W008. A username, Room Code, or client-supplied Player ID alone is insufficient. The reconnect must identify the same Room and active Match through authorized server state.
3. If authenticated before the deadline and the player has not forfeited or the Match ended, replace the prior active connection, invalidate the old connection's action authority, and restore the latest authoritative score, reference, board, cooldown, and Room/Match state. Do not replay old client state as authority.
4. Deadline behavior is serialized with Match end and reconnect. If grace expires while the Match is active, mark the player forfeited once. If the normal deadline happens before grace expiration, finalize using that player's authoritative score and keep them eligible for that result; do not delay completion. If reconnect wins the ordering before deadline, restore connection normally. The exact equal-timestamp ordering must be owner-approved if not determined by W011's existing monotonic event ordering.
5. An explicit leave during gameplay forfeits immediately, receives no grace, and cannot re-enter that active Match even with a valid session. Lobby leave remains a W009 Room operation, not a gameplay forfeit.
6. If only one eligible competitor remains after resolved forfeits, terminate the Match and provide that terminal reason to W012; if none remain, terminate without a winner. W013 reports eligibility and reason; W012 calculates/display results.
7. If the Host voluntarily leaves, W009 transfers Host ownership to the eligible remaining player with earliest join order as required by §7.13; during active gameplay the Host also forfeits under §13.6. For unexpected Host disconnect, apply the confirmed reconnect grace before transferring Host ownership; successor selection for this unexpected-disconnect case remains an explicit decision. Host migration cannot pause, cancel, or alter Match state.
8. Apply security throttles independently from in-game mistake cooldowns. Rejected or throttled technical actions do not create gameplay penalties or extend an existing gameplay cooldown. Deduplicate logical actions by request ID and reject stale Room/Match/Wave versions.
9. W010 bounds message sizes, per-Room/outbound queues, and network writes. A slow client cannot block the Room event loop; after W010 applies its approved disconnect policy, W013 treats the resulting connection loss as unexpected and applies normal grace/eligibility behavior. W013 does not independently drop or reorder accepted authoritative mutations.
10. A process restart may destroy active in-memory Rooms/Matches. The client must leave stale gameplay and show a Thai recoverable unavailable/rejoin-room message. Do not fabricate a resumed Match or result.
11. Record structured security and lifecycle events with opaque IDs and categories. Never log raw credentials, cookies, unnecessary personal data, or full private board contents.

## Business rules and ownership boundaries

- W008 owns Guest Session credential creation, validation, and secure cookie/session proof. W013 consumes verified identity, never implements a parallel identity scheme.
- W009 owns Room membership, Host permissions, lock/kick/code lifecycle, and lobby state. W013 owns Match-time disconnect eligibility, reconnection and anti-griefing signals; reconcile shared behavior through explicit interface contracts.
- W010 owns WebSocket protocol, message-size/message-rate controls, request deduplication contract, bounded queues, slow-client handling, and connection lifecycle signals. W011 Room loop owns deterministic order of gameplay mutations and terminal state. W013's grace and domain-policy rate-limit timers enqueue versioned commands into the authoritative Room event loop, not mutate Match state from timer goroutines.
- W011 remains authoritative for scoring and completion. W013's eligibility decisions must be included in the same atomic/versioned terminal snapshot consumed by W012.
- Confirmed timing values stay configurable only as required by existing requirements; do not offer client or Host control over the 15-second confirmed baseline.
- For security rate limits, §14.4 proposed values (12 card picks/sec/player, 5 create-room attempts/min/session, 20 join attempts/min/session) are tuning defaults for measurement, not gameplay rules. Other operation thresholds need explicit proposals and approval if they affect user-visible behavior.
- Host migration must not imply identity ownership. A Host is a Room permission role, not an administrator of Match outcomes.

## Edge cases and error handling

- Reconnect at just before/at/after the 15-second deadline: deterministic server timestamp and event ordering; no extension by repeated reconnect attempts.
- Reconnect races grace expiry, explicit leave, timer end, or Room close: Room event loop serializes state transition; exactly one terminal outcome wins.
- Old and new WebSocket send actions concurrently after replacement: old connection rejected, including already buffered actions unless accepted before authority replacement according to the sequenced contract.
- Reconnect with wrong Room, old Match, wrong session, expired/invalid cookie, or forfeited state: reject generically without disclosing another Room; do not rejoin active Match after forfeit.
- Reconnect with valid session after Match completed during grace: no active gameplay restoration; allow Results access only through approved W012 projection and current Room membership contract.
- All players disconnect, all forfeit, or only one remains: follow EC-035/036 and send one terminal reason.
- Host voluntarily leaves Lobby/Match: transfer to earliest-join eligible remaining player per §7.13; during Match the Host also forfeits. Unexpected Host disconnect: wait the confirmed reconnect grace before transfer; use the approved successor policy. No active Match cancellation or result manipulation.
- W010 rejects duplicate, reordered, oversized, malformed, or stale transport messages according to its contract; W013 ensures transport rejection never applies gameplay penalty.
- Rate-limit burst reaches threshold: deterministic bucket/key policy, safe rejection response, bounded memory; avoid permanent identity assumptions based only on IP.
- W010 outbound queue is full: W010 applies its documented resynchronization/disconnect policy; W013 handles any disconnect as unexpected while respecting Match completion and Results semantics.
- Server restart loses in-memory state: clear stale client scene and provide recovery action, without promising restoration.
- Cleanup after Room close: stop timers, release references, close connections, remove queue and limiter state within bounded time.

## Acceptance criteria

- **W013-AC-01:** Unexpected disconnect during active Match starts one 15-second server-authoritative grace period, leaves other players progressing, and stores the current authoritative player state without refunding elapsed Match time.
- **W013-AC-02:** Only proof of the original authorized Guest Session permits reconnect; username-only or another player's session cannot restore access.
- **W013-AC-03:** Reconnect before grace expiry restores latest authoritative state and score, replaces the previous connection, and makes the old connection unable to mutate Room/Match state.
- **W013-AC-04:** Grace expiry while Match active causes exactly one forfeit; explicit gameplay leave forfeits immediately and prevents re-entry to that Match.
- **W013-AC-05:** A normal Match deadline during active grace finalizes on time using server score and keeps the player eligible; grace expiration does not delay completion. Boundary races have deterministic serialized outcomes.
- **W013-AC-06:** One remaining eligible player ends the Match with that player eligible to win; all players forfeiting ends with no winner; W013 supplies version-correlated eligibility/terminal reason to W012 without calculating rankings.
- **W013-AC-07:** Host disconnect/migration follows the approved policy, preserves Room consistency, and cannot cancel active play, kick competitors during play, or alter score/results.
- **W013-AC-08:** W013 domain abuse policies remain separate from gameplay cooldown. W010 transport controls provide idempotency, stale-action checks, bounded payload/queues, rate limiting and slow-client isolation; technical rejection never applies or extends gameplay mistake cooldown.
- **W013-AC-09:** Process restart and Room closure produce the documented recoverable client state; lifecycle timers, queues, connections, and temporary identity state are cleaned up safely.
- **W013-AC-10:** Tests cover 15-second deadline boundaries, reconnect/forfeit/Match-end races, connection replacement, eligibility idempotency, Host behavior, domain rate-limit exhaustion, cleanup, and restart error behavior. W010 tests transport duplicate/stale messages, malformed/oversized payloads, queue bounds, and slow clients.
- **W013-AC-11:** Independent implementation review PASS applies to exact final PR SHA; required CI/security checks pass; no BLOCKER/MAJOR findings remain; docs synchronize; PR merges to `develop` by Merge Commit, branch is retained, and evidence is recorded.

## Technical design

- Extend the Room-owned event loop with typed, versioned events such as `ConnectionLost`, `ReconnectAuthorized`, `GraceExpired`, `VoluntaryLeave`, and `HostDisconnected`. W010 emits transport outcomes; W013 timers enqueue events and never mutate Room or Match state directly.
- Keep Guest Session validation behind W008 interface and room membership/Host operations behind W009. W010 maps protocol requests to validated commands and carries request IDs. W011 remains the serialization and score authority.
- Maintain one connection generation per Guest Session/Room. Increment generation on authorized replacement; bind queued actions to generation, request ID, Room revision and Match ID so stale connections cannot act.
- Consume W010's bounded request-deduplication and stale-action contract; W013 owns eligibility idempotency at the Room state transition and does not maintain a second transport deduplication cache.
- Schedule 15-second grace using monotonic server time and a Room revision/token; cancellation/reconnect invalidates the old timer token. At completion, atomically capture eligibility for W012 at the same authoritative revision.
- W013 domain abuse policies govern gameplay and Room operation semantics; W010 owns transport request limits, payload bounds, deduplication, queues, and backpressure. Review interfaces so throttles are not double-counted or reported as gameplay mistakes. Exact keys and thresholds beyond documented suggestions require design/owner approval.
- W013 consumes W010 slow-client disconnect signals and applies ordinary reconnect eligibility rules; W010 owns writer behavior, state snapshot recovery, and outbound queue policy.
- Host continuity uses W009's Room-owner state transition API and approved successor policy; no separate mutable Host field inside W013.
- Client shows Thai reconnecting, reconnected, forfeited, session-invalid, and server-state-lost messages. It must not imply reconnection or score restoration until authoritative state arrives.
- In-memory behavior only for MVP. No Redis, DB, multi-instance routing, or crash recovery introduced.

## Security and privacy

- Enforce SEC-001/002 through W008-authenticated sessions, secure cookies and no secret logging. W013 never accepts usernames or client IDs as authentication.
- Validate every reconnect, Room action, Match/Wave reference, request ID and payload size server-side. Authorization is checked per action, including replacement-connection generation.
- Validate WebSocket Origin in W010. W013 must not weaken Origin or cookie controls.
- Rate-limit resource-creating and gameplay operations separately from penalties. Bound concurrent reconnect attempts, queue memory, dedup cache, and per-session timers.
- Do not return detailed auth/Room existence errors that enable enumeration; use stable error categories and safe Thai UI messages.
- Host cannot kick during Match, modify scores or outcomes, or cancel to avoid defeat. User-visible Room moderation remains W009-owned.
- Log only stable categories and opaque IDs. Avoid raw credentials, IP retention beyond approved operational need, private game state, and personal data.
- Document that server authority does not prevent bots from analyzing visible cards; do not claim complete bot prevention.

## Test plan

- Unit tests with a controllable monotonic clock for grace scheduling, expiry, cancellation, and exact boundary policy.
- Room-loop integration tests for reconnect vs expiry vs timer end vs leave, single terminal transition, eligibility snapshot version, and no time refund.
- Security tests for wrong/expired session, same session replacement, stale connection generation, cross-Room access, post-forfeit re-entry, and safe error/log outputs.
- W013 idempotency/domain tests for duplicate terminal eligibility transitions, stale Room/Match state and burst operation limits. W010 tests duplicate request IDs, reused IDs with changed payload, stale protocol messages, malformed/oversized payloads, queue bounds, and limiter-cardinality exhaustion.
- Concurrency/race tests for simultaneous disconnect/reconnect/leave/completion and the W010 disconnect-signal integration; use Go race detector where available. W010 owns direct backpressure/slow-writer tests.
- Integration tests for Host disconnect and migration under Lobby and active Match, verifying authority boundaries and W009 ownership.
- Browser multi-context tests: refresh/reconnect within grace, reconnect after expiry, two tabs/session replacement, voluntary leave, network loss, timer end during grace, and process/server unavailable recovery UI.
- Run documented Go tests/race/vet, protocol fixtures, `pnpm check`, security checks and Playwright E2E as applicable; report actual PASS/FAIL/NOT_RUN/BLOCKED outcomes and evidence.

## Dependencies and owner decisions

- W008 Guest Session, W009 Room Management, W010 WebSocket Foundation, W011 Multiplayer Speed Battle are implementation dependencies and must be approved/integrated before implementation.
- W012 Results & Rematch consumes W013 eligibility. W012 and W013 specs may be reviewed in parallel; implementation sequence should put W013 before W012 integration.
- Owner decisions required before implementation where not already settled by requirement: exact equal-timestamp ordering at grace/deadline boundary if not defined by W011; successor selection for unexpected Host disconnection; domain abuse limiter thresholds/key policy beyond §14.4 suggested values. Voluntary Host departure successor order is fixed by §7.13 and requires no approval. W010 owns malformed-message and slow-client thresholds/queue recovery.
- Preserve 15-second grace, immediate voluntary leave, Match end without waiting, and eligibility rules. Do not seek approval to alter confirmed behavior.
- Any overlap with pending W008/W009/W010/W011 contracts must be reconciled and independently re-reviewed before implementation.

## Approval

- Project Owner / authorized approver: Pending
- Independent specification review: Pending
- Decision and date: Pending
- W008–W011 dependency approval/merge: Pending
