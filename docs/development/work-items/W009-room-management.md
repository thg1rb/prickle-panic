# W009 — Room Management

Status: READY_FOR_REVIEW
Type: Backend and Frontend Feature
Owner: Project Owner assigns during execution
Last updated: 2026-10-19

## Goals

- Let an authenticated guest create a private Room, invite others by code or URL, and manage the Lobby for an MVP Speed Battle.
- Keep Room membership, host authority, readiness, and lifecycle state authoritative on the Go server.
- Provide transport-independent, concurrency-safe Room domain APIs for the WebSocket gateway and game coordinator.

## Non-goals

- Guest identity/session issuance and recovery (W008).
- WebSocket transport, active-connection arbitration, and event delivery (W010).
- Match engine, authoritative gameplay, results, and rematch (W011/W012), except for Room lifecycle seams and validated Match-start requests.
- Reconnection grace-period behavior, forfeit processing, and unexpected-disconnect host migration timing (W013); W009 must provide explicit lifecycle hooks.
- Classic Rush and Party Chaos implementation, room browsing, public matchmaking, spectators, waiting queues, durable storage, multi-server coordination, or permanent identity bans.
- Changing requirement-defined Room rules or TUNING values without the required approval/playtest evidence.

## Scope

- Create an in-memory private Room with a unique internal ID, active shareable code, initial Host, and Host RoomPlayer.
- Join by Room Code or invitation URL; enforce state, lock, membership, capacity, reservation, duplicate-name, and rate-limit rules atomically.
- Provide authorized Lobby snapshots and mutations: leave, Lobby-only username change (P1), ready, validated settings, game-mode selection, Lobby kick, Room lock, and P1 code regeneration. Lock reversibility is an unresolved product behavior.
- Enforce Room state transitions and invariants while exposing hooks for W010 connection state and W011/W012 Match lifecycle.
- Expire and clean up eligible Rooms, codes, memberships, reservations, indexes, and event-loop resources.
- Implement MVP-facing create/join/Lobby flows with Thai player-visible copy and responsive desktop/mobile layout, integrated through the W010 protocol when available. Until W010 is implemented, domain and UI work must use a clearly isolated adapter seam.

## Requirement references

- `docs/REQUIREMENT.md` §§3.1–3.2 (MVP and post-MVP scope), §§6.4–6.7 (Guest Session, RoomPlayer identity, username changes), §§7.1–7.15 (Room rules), §§12.2–12.3 (start validation/countdown), §§13.1–13.8 (reconnection/leave boundaries), §§14.2, 14.4, 14.6–14.7 (authority, limits, Host and Room abuse), §§15 (EC-003–015 and related cases), §§16.1–16.2 (FR-005–023), §19.5 (stable errors), §20 (security), §21 (invariants), §§22–25 (architecture/concurrency/security), §28 (AC-002–005), §32 (W009).
- W005 contract: `docs/protocol/v1/README.md` and `message.schema.json` (strict v1 envelope and Room snapshot contract; currently lacks several W009 commands and Lobby fields).
- W008 boundary: `docs/development/work-items/W008-guest-session.md` (currently reviewed but still awaiting Project Owner approval and merge). W009 consumes an authenticated session principal; it does not accept caller-supplied Guest IDs.
- Status fidelity: Room Management is BASELINE (§40). Six-character Room Code is an INITIAL PROPOSAL (§7.2), Room expiry values are TUNING (§7.15), Room create/join limits are initial suggestions subject to testing (§14.4), and protocol extensions/capabilities described below are design recommendations pending review.

## Expected behavior

1. An authenticated guest can create a Room; the creator becomes the Host and first participant. The server returns the Room snapshot and a share link based on the configured public origin, not a hard-coded production domain.
2. Another authenticated guest can join via a valid active Room Code or `/room/<code>` invitation. The client preserves the invitation destination while it obtains a guest username/session.
3. A join is rejected consistently for unknown/expired code, locked Room, full Room including reserved disconnected slots, duplicate same-Room username, or a non-joinable Room state. A new participant cannot join COUNTDOWN, IN_PROGRESS, or RESULTS.
4. The Lobby shows Room Code/share link, player list, Host, readiness, player count, selected MVP mode, validated settings, and relevant controls. Player-visible copy is Thai.
5. Every player, including Host, explicitly sets Ready. The server accepts Match start only with 2–6 eligible connected participants, all Ready, valid settings, and an idle Lobby; failed or duplicate starts do not create duplicate Matches. During COUNTDOWN, only the Host can explicitly cancel; the pending Match is cancelled and Room returns to LOBBY. If a player leaves during COUNTDOWN, revalidate conditions and cancel when they no longer hold.
6. Any gameplay-setting change in Lobby resets every participant's Ready state. Match settings cannot change after countdown begins.
7. Only the Host can perform Host operations, and only in the permitted Room state. Kick is Lobby-only; a kicked Guest Session cannot rejoin that Room. Room lock rejects new joins but permits authorized reconnection. The requirement does not specify whether/how a Host can unlock a Room; defer unlock behavior pending owner decision. Code regeneration invalidates old invitation codes while preserving current members.
8. A player may explicitly leave a Room. In LOBBY, remove that RoomPlayer and transfer Host to the eligible remaining participant with the earliest join order when the Host leaves. During active gameplay, route voluntary leave to W013's forfeit coordinator immediately, with no reconnect grace and no re-entry into that active Match; Room membership/state is retained as required to finalize the Match. If leave occurs during COUNTDOWN, cancel/revalidate as required by EC-015. Unexpected disconnect is represented as disconnected/pending-host state; W013 later applies the required grace period before migration. Host changes cannot alter Match outcomes.
9. Room operations are serialized by the owning Room context. Different Rooms can progress concurrently. Slow or failed event consumers cannot block authoritative Room mutation or cause unbounded memory growth.
10. Eligible empty/inactive Lobby Rooms expire according to configured initial values; an active Match is not expired by the Lobby timeout. Cleanup removes all indexes and stops timers/loops exactly once.
11. A Guest Session may keep Room memberships in multiple Rooms, but W010 must authorize commands through exactly one active Room connection for that session across the server. Replacing a connection revokes the older connection's command/gameplay authority and reports an unexpected disconnect to W013; replacement does not silently leave or forfeit the prior Room/Match.

## Business rules

- Room starts in LOBBY. Supported states are LOBBY, COUNTDOWN, IN_PROGRESS, RESULTS, and CLOSED (§7.5). Primary flow is LOBBY → COUNTDOWN → IN_PROGRESS → RESULTS → LOBBY; eligible termination/expiration can transition to CLOSED. State changes are server-owned.
- At most six participants/active reservations exist; minimum Match size is two. A disconnected participant with an active reservation counts toward the cap. An active reservation is maintained by W013 policy; until then Room Manager accepts explicit reserve/release signals rather than inventing its own timeout.
- Duplicate usernames are disallowed within one Room; names may repeat across Rooms. The requirement does not define comparison/canonicalization semantics. The recommendation is compare the W008-canonical display names after Unicode NFC normalization and full Unicode case folding, without changing displayed spelling. This is a material player-visible policy and needs Project Owner approval before implementation.
- Host settings changes reset readiness. For the MVP, Speed Battle is the only selectable mode; Classic Rush and Party Chaos remain deferred until after Closed Alpha. The precise set of Host-editable settings is unspecified even though FR-017 is P0. Recommendation: implement a versioned, allow-listed settings object with §9.2 defaults and decide which values are editable before implementation; do not silently expose arbitrary tuning or imply settings are immutable.
- Readiness is idempotently set to the requested boolean; it is not a toggle dependent on stale client state. A settings version change clears all readiness. Every accepted mutation increments a Room revision and publishes a consistent snapshot/event sequence through W010.
- Create/join operations must enforce configured limits. §14.4 suggests 5 create attempts/minute/session and 20 join attempts/minute/session, subject to testing. Add a bounded supplementary source limiter so fresh Guest Sessions cannot bypass all protection; source identity is abuse mitigation, not permanent identity.
- Six-character uppercase Room Codes using a confusion-reduced alphabet follow §7.2's initial proposal only if approved. Secure random generation and active-code uniqueness are required. The current protocol accepts 1–12 characters and does not confirm six. Regeneration retires the old code atomically.
- §7.15 initial expiration values are approximately 5 minutes for empty Rooms and 30 minutes for inactive Lobbies. Keep them configurable and preserve as TUNING values; active Matches are exempt from Lobby expiry.
- A RoomPlayer ID is an opaque Room-scoped protocol identity and must not expose Guest Session credentials or the Guest ID. Repeated actions must be idempotent or rejected without partial state mutation.

## Edge cases and error handling

- Concurrent create/join/code-regeneration operations must not duplicate active codes, exceed capacity, admit two memberships for one Guest Session, or create partial memberships. Collision handling retries a bounded number of times and fails closed.
- Two simultaneous join attempts for the final slot: exactly one succeeds; the other receives the documented full-room error.
- Concurrent ready/settings/start commands are processed in order by the owning Room loop. A stale ready request or settings revision cannot start a Match using outdated readiness.
- A Guest Session may retain RoomPlayer membership in more than one Room; this does not authorize simultaneous active Room connections. INV-004 requires at most one active Room connection per Guest Session across all Rooms. W010 owns a global connection registry keyed by Guest Session; a replacement connection invalidates the older connection's Room/gameplay authority and reports it as an unexpected disconnect to W013, without implicitly leaving or forfeiting its Room/Match. W009 maintains a Guest Session → set of Room memberships and checks that commands come from the currently authorized connection context supplied by W010.
- Repeated room.create or room.join request IDs do not create duplicate Rooms or membership mutations. Protocol-level request deduplication is currently PROPOSED in W005/W010; W009 must define the Room Manager's idempotency boundary and failure behavior.
- Joining while a Room transitions to COUNTDOWN is serialized: it is either accepted before transition and included in start validation, or rejected after transition.
- Host kick of self, unknown player, already departed player, or any player outside LOBBY fails without mutation. A kicked principal remains denied until that Room closes; the block is Room-local and does not claim a durable account ban.
- Room code regeneration while joins are in flight is atomic: exactly the current code resolves; previous code becomes unavailable, while existing members remain.
- Empty Room after Host leave closes immediately or follows the configured empty-room lifecycle; recommended behavior is immediate close and cleanup because no Host/member remains. Confirm this with owner before implementation if the requirement does not resolve it.
- Duplicate username comparison policy, editable settings list, unlock behavior, and immediate empty-Room close behavior are approval blockers for the detailed spec. Unknown production origin must fail link construction safely; use configured public origin or a relative invitation path.
- Process restart destroys in-memory Rooms and active Matches, consistent with MVP requirements. The client must receive a recoverable unavailable/recreate path; no stale membership or invitation is treated as valid.

## Acceptance criteria

- **W009-AC-01:** Authenticated Guest Session creates one private Room; creator is Host and a RoomPlayer with opaque IDs; state and active-code indexes are consistent under concurrent creation.
- **W009-AC-02:** Joining by code and invitation route produces the same authorized membership result and preserves the intended invitation destination through username/session entry.
- **W009-AC-03:** Server enforces Room state, lock, capacity (including reservations), same-Room duplicate usernames, supported stored memberships across Rooms, and create/join rate limits atomically, with stable errors and no partial mutations. Room Manager rejects missing/stale W010 connection context; W010 must enforce INV-004 globally across the Guest Session before a command reaches a Room.
- **W009-AC-04:** Lobby projections include required public Room/player/Host/Ready/mode/settings information and never expose credentials, Guest IDs, private game boards, or internal-only data; all player-facing strings are Thai.
- **W009-AC-05:** Ready, settings, mode, start, countdown cancellation, kick, lock, and code-regeneration operations enforce actor permissions and state. Only Host may cancel Countdown; cancellation returns Room to LOBBY and closes/cancels the pending Match once. Settings changes reset all Ready states; unauthorized or invalid actions leave state unchanged.
- **W009-AC-06:** Lobby leave removes the RoomPlayer and Host departure follows join-order migration; active-Match leave emits exactly one typed voluntary-leave event for the W013 forfeit coordinator, grants no grace, and marks the RoomPlayer ineligible to re-enter that Match. W009 retains the Match participant record for finalization; W013 owns applying the forfeit outcome. COUNTDOWN leave revalidates and cancels if conditions fail. Disconnect/reconnect hooks defer unexpected-disconnect timing to W013; authorized reconnect can remain possible while locked.
- **W009-AC-07:** A player can change username only in LOBBY. The server applies the approved W008 username policy, enforces same-Room duplicate rules atomically, preserves Guest Session identity, and changes no score or Match ownership. During COUNTDOWN/IN_PROGRESS/RESULTS the change is rejected.
- **W009-AC-08:** Room transitions, Countdown/start/cancel guards, expiry, close, and cleanup are bounded, deterministic under an injected clock, and safe under concurrent access; no active Match is removed by Lobby expiry.
- **W009-AC-09:** Room ID/code generation uses the reviewed secure source; collisions, entropy failures, capacity exhaustion, and limiter exhaustion fail safely without leaked secrets or partial state.
- **W009-AC-10:** The v1 protocol/schema and fixtures cover every W009 request, event, public Room snapshot field, and stable error; compatibility impact is reviewed before adding or changing fields.
- **W009-AC-11:** Tests cover state-transition tables, all capacity/state/permission boundaries, concurrent last-slot joins and code rotation, duplicate request behavior, settings/readiness races, countdown cancellation and cancellation races, Lobby rename, voluntary leaves in Lobby/Countdown/Match, expiration cleanup, host departure, and Room isolation. Go race tests and frontend/browser flows pass; exact outcomes are recorded.
- **W009-AC-12:** Independent review PASS applies to exact final PR SHA, required CI passes, no BLOCKER/MAJOR findings remain, PR merges to `develop` with Merge Commit, branch is retained, and evidence is recorded.

## Technical design

- Domain boundary: `apps/server/internal/room/` owns Room aggregate state, RoomPlayer membership, Host permissions, codes, settings, readiness, revisions, lifecycle, and pure Room transitions. It consumes an authenticated principal (`GuestSession` from W008) plus a W010-validated current connection context; it does not parse HTTP/WebSocket frames or trust client-supplied actor IDs. Connection replacement/presence is signaled by W010 and W013, and Room Manager never treats stored membership as proof of an active connection.
- Manager boundary: process-local `RoomManager` indexes Room ID, active code, and Guest Session membership. The manager atomically creates the aggregate and unique code, routes commands to the Room-owned bounded event loop, and performs lifecycle cleanup. No database, Redis, or cross-process coordination is introduced.
- Concurrency: one bounded command queue/event loop per Room serializes mutations. Manager indexes use short critical sections and never hold a global lock while waiting on Room processing. Each command has a context/deadline; queue full returns a retryable stable error. Notifications are snapshots/events delivered non-blockingly; a slow client cannot stall the Room loop. W010 owns network write queues/backpressure.
- Domain interfaces should include create, join, leave, change-username, set-ready, update-allowed-settings, select-mode, start/cancel-countdown, lobby kick, lock, regenerate-code, presence/reservation updates, read-authorized snapshot, expire/close, and shutdown. Each command validates against current state when dequeued; caller validation is not authoritative. An active-Match leave emits an idempotent typed voluntary-leave event for W013 before releasing active connection authority; the RoomPlayer is ineligible to re-enter that Match while the Match participant/result record remains until finalization. W013 owns applying the forfeit outcome; W009 exposes and tests the integration seam.
- Room state snapshots are immutable copies with a monotonic Room revision. Public snapshot projection includes only protocol-approved Room/player fields. Guest Session ID maps to private membership indexes; use Room-scoped opaque Player IDs externally.
- Host migration uses monotonic join order, with a deterministic tie-breaker for any implementation-level equal order. Voluntary leave is immediate. Unexpected disconnect marks the Host unavailable and signals W013; W013's grace-period completion authorizes migration. If no eligible participant remains, follow the approved empty-Room lifecycle.
- Protocol v1 currently has `room.create`, `room.join`, `room.leave`, `room.set_ready`, and `match.start`, but lacks Host operations, Room state/mode/settings/lock metadata, invitation URL output, `room.change_username`, and `match.cancel_countdown`/countdown-cancelled event. W009 should add explicit strict schema branches and fixtures for required commands/events/errors before transport implementation; don't overload existing messages or accept arbitrary extra fields. W010 will bind them to WebSocket routing. Review compatibility under W005's strict v1 policy; deployment is coordinated, and no live users are assumed without verification.
- Frontend organization should add a Room feature for create/join/invitation route/Lobby, consume a typed Room API/store, and project it into React components. React owns Lobby UI state; Phaser remains isolated to gameplay and receives no Room authority. Keep route destination across guest bootstrap. Render usernames as text in bidi isolation, provide accessible Thai labels, focus states, and touch-sized actions.
- W008 remains pending owner approval at the time of this specification. Its proposed session expiration, cookie, username safety filter, and credential design are not treated as confirmed here. W009 requires only the approved authenticated session principal and stable session identity contract; if W008 changes, re-review this boundary before implementation.

## Security and privacy

- Authorize every mutation from server-resolved Guest Session → RoomPlayer membership; verify Host ownership server-side for privileged operations. Never accept Guest ID, Player ID, Host ID, score, or connection state as proof from clients.
- Validate and normalize Room Code at the boundary (uppercase ASCII alphabet only as specified by chosen code policy), cap payload sizes, reject malformed/unknown fields, and return stable errors without disclosing whether sensitive internal IDs exist.
- Enforce exact-origin/CSRF/session authorization at W008/W010 boundaries, plus separate create/join rate limits and bounded Room, code, queue, and timer resources.
- Use cryptographically secure Room codes and opaque identifiers; avoid logging session cookies, credentials, full invitation query secrets, or unnecessary user data. Room invitation code is a bearer capability and should not be copied to analytics/logs.
- Use escaped text rendering for usernames and reject unsafe input according to the approved W008 policy. Do not generate HTML from username or Room Code.
- Bound Room count and per-Room queue/event resources. Define shutdown semantics so all loops and timers terminate and no goroutine survives a closed Room.

## Test plan

- Go unit tests for state transition table, membership and Host rules, readiness/settings reset, eligible mode/settings validation, duplicate-name comparison, kick denylist, lock/recode behavior, expiration, immutable snapshots, and error atomicity.
- Go integration tests for concurrent create/code collisions, simultaneous final-slot joins, stored cross-Room membership without active authority bypass, Room isolation, event ordering, bounded queue saturation/slow consumers, clock-driven expiry, close/shutdown, and race detector. W010 must add a global connection-registry integration test proving two Room connections for one Guest Session cannot both act, and that replacement revokes the old connection and signals unexpected disconnect to W013.
- Protocol schema fixtures for all new commands, projections, bounds, errors, unknown fields, and any compatibility decision; run `pnpm protocol:check`.
- Frontend component tests for Thai create/join/errors/Lobby/readiness, invitation preservation through guest entry, accessible interaction, and state updates. Browser E2E with at least two contexts covers create, invite/join, ready, Host controls, capacity, lock, and recovery path. Test desktop/mobile viewport and Thai text overflow.
- Test each acceptance criterion and record `PASS`, `FAIL`, `NOT_RUN`, `BLOCKED`, or reasoned `NOT_APPLICABLE`; developer and reviewer results remain separate.

## Dependencies and risks

- W005 Protocol Contracts (DONE) provides strict envelope/schema rules; W008 Guest Session is reviewed but owner approval and merge are still pending.
- W010 consumes Room APIs and provides active connection/presence, event delivery, network backpressure, FR-008 one-active-connection enforcement per RoomPlayer, and INV-004 one-active-Room-connection enforcement across the Guest Session. Replacing a connection must revoke the old connection's authority and notify W013 using the unexpected-disconnect path.
- W011/W012 own Match coordination and state transitions; W013 owns unexpected-disconnect grace periods, forfeit, and delayed Host migration.
- Major design risks: room code entropy/collision handling under high load; multi-index atomicity; Room event queue lifecycle/leaks; undefined username equality; unspecified Host-editable settings; protocol v1 field evolution; interactions between kick, rejoin, and guest identity; and in-memory loss at restart.
- Decisions requiring Project Owner approval before implementation: duplicate-username equality/canonicalization; whether the six-character code proposal is approved; the permitted Host-editable Speed Battle setting list; whether a Host can unlock a locked Room; and whether a Room with no remaining members closes immediately. Keep all such recommendations marked PROPOSED until decided.

## Approval

- Project Owner / authorized approver: Pending
- Independent specification review: PASS at `48d7d76f0fafa545d72938942c2d5e400a0ad4ed`, recorded in `docs/development/reviews/W009-review.md`
- Decision and date: Pending
- W008 dependency approval/merge: Pending
