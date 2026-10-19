# W010 — WebSocket Foundation

Status: READY_FOR_REVIEW
Type: Backend and Frontend Infrastructure Feature
Owner: Project Owner assigns during execution
Last updated: 2026-10-19

## Goals

- Establish an authenticated, same-origin WSS connection for Room actions and real-time state synchronization.
- Route versioned protocol requests through validated Room and Match boundaries without giving the transport layer authority over game rules.
- Ensure one Guest Session cannot control multiple active Room connections; bound connection, request, and outbound resources.
- Provide a reconnectable protocol/client foundation that W011–W013 can extend for authoritative play, Match recovery, and anti-griefing.

## Non-goals

- Creating/recovering Guest Sessions (W008), implementing Room rules or Lobby UI/business behavior (W009), or changing the approved protocol contracts (W005) without review.
- Authoritative Speed Battle gameplay, card validation, scoring, Match timing, results, and rematch (W011/W012).
- Reconnection grace-period policy, gameplay restoration, forfeit, and delayed Host migration (W013); W010 owns transport replacement and presence signals only.
- Classic Rush, Party Chaos, public matchmaking, durable transport/session state, multi-server routing, Redis, databases, or Kubernetes.
- Treating client state, scores, or clocks as authoritative.

## Scope

- Add a Go WebSocket Gateway and lifecycle manager using the standard HTTP server and the requirement's preferred `github.com/coder/websocket` library after API/version/security review.
- Authenticate the handshake using the approved W008 Guest Session cookie and validate the exact browser Origin against configured allowed origins before upgrading.
- Enforce a global single-active-connection invariant per Guest Session (INV-004); atomically revoke a replaced connection's authority and notify Room/W013 as an unexpected disconnect.
- Strictly decode the W005 versioned JSON protocol, validate size and schema, dispatch Room actions to W009 and future gameplay actions to W011, and map failures to stable protocol errors.
- Implement bounded request deduplication, per-connection outbound queues, rate limits, heartbeat/deadline handling, safe close/shutdown, and authorized snapshots/incremental events.
- Add a typed TypeScript WebSocket client that owns connection lifecycle, request correlation, sequencing, snapshot recovery, and Thai connection/error UI; React owns product state and Phaser receives display state/intents only.
- Provide local development through the same-origin Vite `/api` or configured WebSocket proxy path; production uses HTTPS/WSS and configured Origin/public host.

## Requirement references

- `docs/REQUIREMENT.md` §§6.4–6.6 (Guest Session identity and active connection replacement), §§13.1–13.9 (disconnect/reconnect boundary), §§14.1–14.4 (authority, idempotency, rate limits), §§15 (EC-009, EC-031–045), §§16.1 and 16.7 (FR-008, FR-064–076), §§17 (INV-004, INV-014, INV-017), §§18.1–18.11 (architecture, Room event loop, in-memory single instance), §§19.1–19.12 (protocol, validation, synchronization, private state, ordering, timing, slow clients), §20 (SEC-001–012), §§21–25 (invariants, testing, security, operations), §28 (AC-009, AC-040–045), §32 (W010).
- W005 contract: `docs/protocol/v1/README.md` and `message.schema.json`; currently defines strict v1 envelopes, Room request/event subset, snapshots, errors, and proposed duplicate/frame behavior.
- W008 session boundary: `docs/development/work-items/W008-guest-session.md` in PR #10 contains the Project Owner-approved session policy baseline; its integration into `develop` remains pending. W010 consumes only that approved principal/cookie contract. W008 owns session creation, expiry deadlines, and authorization revocation; W010 owns enforcing revocation on established WebSocket connections.
- W009 Room boundary: `docs/development/work-items/W009-room-management.md` in pending PR #11; owner approval and merge are still pending. W010 must not duplicate Room state or treat stored Room membership as active connection authority.
- Status fidelity: WebSocket and architecture requirements are BASELINE (§40). The preferred library is a requirement recommendation but version and API compatibility must be verified. W005's 16 KiB frame ceiling and duplicate request cache behavior are PROPOSED and require technical review before use. The server must enforce message-size bounds, but the initial numeric limit is not confirmed.

## Expected behavior

1. The client opens a same-origin WSS connection with its HttpOnly Guest Session cookie. The Gateway validates Origin, session, protocol version/handshake, limits, and server availability before accepting an upgraded connection.
2. The registry allows at most one active Room connection per Guest Session across all Rooms. A newly authorized connection atomically supersedes the older one; the older connection can no longer submit authorized mutations even during the close race. The old Room receives an unexpected-disconnect signal, not an explicit leave/forfeit signal.
3. Every request is versioned JSON with request ID and a schema-bounded payload. Malformed JSON, unknown fields, wrong version, oversized frames, invalid request IDs, or invalid session/Room authority fail safely with stable errors or a documented protocol close.
4. Valid Room requests are dispatched with server-resolved Guest Session, connection generation, Room membership, and RoomPlayer context. The server ignores caller-supplied identity/authority fields and checks state/Host permissions through W009. Gameplay requests are not executed until the W011 authoritative engine is integrated.
5. The client receives an authorized Room snapshot on connect/join/recovery and ordered incremental Room events. Each Room's event sequence advances monotonically. A sequence gap or incompatible state triggers an authorized snapshot request; another player's private board is never included in a Room broadcast.
6. Identical duplicate request IDs with identical action context do not repeat mutations while retained; reuse with different content returns a stable duplicate error. Cache size/lifetime are bounded and documented; cache expiration and reconnect boundaries are explicit.
7. A full/slow outbound queue never blocks a Room's authoritative event loop. The connection is closed or resynchronized by a bounded policy; W013 receives a disconnect signal and applies its separate grace/forfeit policy.
8. Connection close, heartbeat timeout, server shutdown, invalid session, replacement, and write failure release registry slots, queues, subscriptions, timers, and goroutines once. A later authenticated reconnect can request a fresh snapshot.
9. The React client displays connection/loading/reconnecting/failure state using Thai text, reconnects under a bounded backoff policy, and never replays non-idempotent actions blindly. Server authority remains final; prediction is presentation only.

## Business rules

- Authentication happens before upgrade. WSS is mandatory in production. Cookie semantics and session expiry come from the approved W008 specification; no new credential mechanism is introduced in W010.
- Validate the exact `Origin` against an explicit allowlist, require an Origin for browser clients, and reject unapproved/missing origins. SameSite cookies are defense in depth and do not replace Origin validation. No permissive cross-origin credential mode is enabled.
- INV-004 is global per Guest Session, not per RoomPlayer: only one active Room WebSocket connection may hold mutation authority for a session across the server. Replacing a connection atomically revokes prior authority before the replacement can issue a Room mutation. The previous Room membership remains stored; W013 treats the connection loss as unexpected and applies the 15-second grace policy if in Match.
- Each active connection uses one bounded inbound reader and one serialized outbound writer. Room event processing remains owned by W009's Room loop; WebSocket write latency cannot hold Room locks or block the loop.
- Protocol JSON follows the strict, versioned W005 contract. Unknown fields are rejected. Unsupported versions use the reviewed stable error/close behavior. Additive schema evolution requires an explicit v1 compatibility decision and coordinated deployment; breaking changes use a new protocol version.
- **PROPOSED protocol resolution pending Project Owner approval and W005 contract review:** transport `seq` is scoped to one authorized recipient stream for a RoomPlayer, is contiguous and monotonic for events delivered to that recipient, and is assigned by the authoritative Room owner. A private event delivered to another player's stream does not consume this recipient's `seq`; no redacted placeholder is sent. The Room may maintain a separate internal global revision for serialization/diagnostics. Snapshots include the recipient-stream synchronization watermark, and a reconnecting connection starts from a fresh authorized snapshot/watermark. Room snapshots never contain gameplay boards; W011's game snapshots are per-player. This resolves the mismatch between W005's current Room-scoped wording and private event projections without changing the requirement to order authoritative events.
- Request IDs are scoped to the authenticated Guest Session (and where necessary Room/match action context); payload fingerprints prevent same-ID/different-action mutation. Deduplication retention is bounded and should outlive expected retransmission/reconnect windows without retaining secrets or unbounded data. Exact cache capacity/lifetime are PROPOSED technical values to be validated in T01.
- Enforce frame-size and message-rate limits before decoding/dispatch, plus connection caps and bounded queues. W005's proposed 16 KiB frame size must be checked against the largest Room and per-player Match snapshots before adoption. No per-IP value or heartbeat interval is promoted without measurement/review.
- A replacement WebSocket is transport replacement, not voluntary Room leave. Close the old connection and signal unexpected disconnect. W008 owns session expiry and revokes authorization immediately; W010 invalidates and closes an established WebSocket on expiry and emits a typed `session_expired` signal, not an unexpected network-disconnect signal. W013 consumes that signal and applies the approved immediate active-Match forfeit through the Room event loop. The 15-second grace applies only to unexpected network loss.
- Development and production deployment must keep browser HTTP and WebSocket traffic same-origin where practical. Production reverse proxy must pass Upgrade/Connection headers, enforce TLS, and preserve the configured Host/Origin policy.

## Edge cases and error handling

- Session expires between HTTP bootstrap and WebSocket handshake: reject upgrade with a non-sensitive HTTP/auth response; client returns to the Thai guest-entry flow.
- Session expires after a WebSocket is established: reject all subsequent actions immediately, invalidate and close the connection, release its resources once, and emit a typed session-expired signal for W013; do not classify expiry as network loss or start the 15-second grace. No reconnect may restore the expired session's authority.
- Two simultaneous connections for one Guest Session race to register: one generation wins atomically; at no point can both mutate. Stale generation requests are rejected even if the old socket has not completed close.
- A replacement connection arrives while the previous connection is in a Room/Match: retain membership and notify W013 via unexpected-disconnect. It does not apply voluntary-leave forfeit; a reconnection within the policy window restores authorized state.
- Origin missing, malformed, mismatched scheme/host/port, or untrusted proxy metadata: reject before upgrade. Trust forwarded host/proto only from configured reverse proxies.
- Unknown protocol version, malformed envelope, oversized/fragmented message, invalid UTF-8, unknown fields, missing request ID, or invalid payload: bounded parse/reject/close behavior; never panic or log raw credentials/payload secrets.
- Duplicate request while in-flight returns/correlates to a single logical operation; duplicate after completion uses cached result if retained. Same request ID with different payload/context is rejected; cache saturation fails closed or evicts only entries past approved retention, with no duplicate mutation.
- Client disconnects while outbound queue has pending events: release queue and subscription; reconnect snapshot repairs state. A slow consumer is isolated from other players and other Rooms.
- Room closes during a queued request or broadcast: return stable unavailable error; no mutation occurs after closure. Server shutdown stops accepting handshakes, revokes connections, closes queues, and joins goroutines under a bounded deadline.
- If snapshots exceed the configured frame limit, fail safely and surface a diagnosable server-side metric; do not truncate JSON or silently omit private/public state. Revisit snapshot projection/limit before release.
- Simultaneous connection replacement plus stale Room command: generation validation rejects stale command regardless of client timing; integration tests cover the race.
- Multiple stored Room memberships are allowed by W009's pending proposal, but the connection registry still permits only one active Room connection per Guest Session. The design must not infer membership deletion or forfeit from replacement.

## Acceptance criteria

- **W010-AC-01:** Production upgrade is available only over TLS/WSS, validates configured browser Origin and an authorized W008 Guest Session before upgrade, and rejects invalid/missing session or disallowed Origin without allocating a connection slot.
- **W010-AC-02:** At most one active Room connection per Guest Session can issue commands globally across Rooms. Replacement invalidates old authority atomically, rejects stale-generation requests, and sends an unexpected-disconnect event to the old Room/W013 without voluntary leave/forfeit.
- **W010-AC-03:** Requests conform to the approved strict v1 schema, have bounded UTF-8 JSON frame sizes, request IDs, rate limits, and stable error mapping; malformed, unknown, oversized, and unsupported inputs do not mutate Room/Match state.
- **W010-AC-04:** The Gateway dispatches actions using server-resolved identity and W009 Room membership, enforces authorization and state before domain commands, and exposes no gameplay authority before W011.
- **W010-AC-05:** Under the approved W005 sequence design, events use a contiguous sequence per authorized recipient stream, and snapshots include that stream's watermark. Clients detect missing events in their own stream and recover with an authorized snapshot. Private events sent to another player's stream do not create gaps, and no redacted placeholder reveals their activity. Public Room state never leaks private Match boards.
- **W010-AC-06:** Duplicate request protection is bounded, session-scoped, payload-aware, tested for in-flight/completed/same-ID-different-payload cases, and cannot repeat a Room mutation within the documented retention window.
- **W010-AC-07:** Inbound processing, per-connection queues, connection count, timers, and shutdown are bounded. Slow or disconnected clients cannot block the Room event loop or other Rooms; queue overflow follows the documented close/resync policy.
- **W010-AC-08:** Connection replacement, close, heartbeat timeout, session expiry, and shutdown release all registrations, subscriptions, queues, timers, and goroutines exactly once. Expiry on an established socket immediately revokes action authorization, closes the connection, and emits a typed session-expired signal to W013; it never starts network grace. Reconnect restores transport and an authorized Room snapshot only with a currently valid W008 session; W013 owns Match network-grace/recovery semantics and immediate active-Match forfeit handling.
- **W010-AC-09:** The TypeScript client has tested lifecycle, request correlation, event sequencing, snapshot recovery, bounded retry, and Thai error/status UI; it never stores the session credential or trusts client scores/timers.
- **W010-AC-10:** Integration/concurrency/security tests verify Origin/session handling, including expiry after an established WebSocket and rejection of all later actions, typed expiry signaling (distinct from unexpected network loss), INV-004 races, stale generation rejection, Room isolation, duplicate actions, event order, private-state projection, slow-client isolation, resource cleanup, and protocol fixtures. Actual developer/CI/reviewer outcomes are recorded.
- **W010-AC-11:** Independent implementation review PASS applies to the exact PR SHA; no BLOCKER/MAJOR findings remain; required CI/security checks pass; docs are synchronized; PR merges to `develop` with Merge Commit and branch retained.

## Technical design

- Backend package proposal: `apps/server/internal/transport/websocket/` owns upgrade, connection reader/writer, schema dispatch, close/retry mapping, and connection cleanup. `apps/server/internal/connection/` (or an equivalent small session module) owns the global active Guest Session → generation/connection registry. Final package layout is reviewed against repository conventions in T01.
- HTTP integration: register `/api/v1/ws` on the existing Go `net/http` server. Authenticate using W008's middleware/session manager and authorize exact Origin before calling the reviewed WebSocket library's Accept/upgrade function. Configure trusted proxy headers only through explicit deployment config. Keep HTTP session APIs and WebSocket upgrade same-origin.
- Connection lifecycle: `Connection` holds immutable session/connection IDs, monotonic generation, subscriptions, bounded inbound state, outbound queue, close signal, and cancellation context. A registry lock/CAS serializes replacement. Before activating a new generation, revoke old command authority; then close old socket and publish unexpected-disconnect to its Room/W013. Validate generation again at dispatch, not only during handshake.
- Reader/dispatcher: one reader enforces read limit/deadline, UTF-8 and JSON envelope decode, schema validation, rate limit, correlation, and dispatch. It does not mutate domain state directly. Each request carries server-resolved `GuestSession`, connection generation, and Room membership to the W009 command interface. W011 registers typed gameplay handlers later.
- Writer: one writer serializes all outbound frames from a bounded channel; no concurrent writes on one socket. Broadcast producers enqueue non-blockingly. Queue saturation invokes the documented slow-client close/resync action and publishes disconnect without blocking the Room event loop.
- Event synchronization proposal: W009's Room loop owns serialization and internal Room revision. Gateway maintains a transport sequence for each authorized RoomPlayer stream and assigns the next number only when an event is delivered to that recipient. The same public event receives each recipient's next stream number; a private event advances only its owner's stream. No-op/redacted events are not emitted. Snapshot request returns the current authorized projection with that stream's watermark. Detect local delivery gaps within that stream and request/send a snapshot; never patch state with an out-of-order event. Project Owner approval and W005's wording/fixtures must be updated as a coordinated v1 contract clarification before implementation.
- Idempotency: bounded cache keyed by Guest Session plus request ID and relevant action scope, storing payload hash and terminal response/event correlation. Same key/hash returns prior outcome; same key/different hash returns `ACTION_DUPLICATE`. Cache design must coordinate with W005's proposed semantics and avoid caching secrets/private state; tune retention and capacity based on retry expectations.
- Frontend structure: a typed `WebSocketClient` service owns connect, close, correlation, event sequence, bounded retry, and snapshot request; feature stores/components consume typed actions/events. React owns user-facing connection state and all Thai UI. Phaser gets render data and emits intentions via bridge; it never opens an independent socket.
- Local/test configuration: permit `ws://` only for explicitly configured loopback development/test origins; do not weaken production WSS/origin checks. Vite proxy uses the same browser origin. Production deploys one Game Server behind HTTPS reverse proxy with WebSocket Upgrade support and health checks; no sticky multi-instance routing is needed until architecture is revised.
- W010 must not decide unresolved W008 policies. If W008 changes the session cookie/identity/API shape, update this design and review the exact dependency before implementation. W010 also consumes only owner-approved W009 membership/command contracts.

## Security and privacy

- Authenticate before upgrade; session cookie is HttpOnly/Secure per approved W008 policy and is never echoed in protocol, browser storage, URL, logs, metrics, or error content.
- Validate browser Origin exactly (scheme/host/port) against configured allowlist; reject null, missing, or untrusted origin for browser clients. Ensure reverse-proxy metadata is trusted only from configured immediate peers.
- Enforce authorization on every action using current Guest Session, active connection generation, Room membership, Host permission, Match membership, and state. Client-provided IDs are selectors only, never credentials.
- Bound frame size, JSON nesting/complexity, message rate, concurrent connections, read/write deadlines, inbound work, outbound queue, idempotency cache, and event subscriptions. Close abusive sessions with stable, non-sensitive codes.
- Revoke old connection authority atomically on replacement. Prevent close/dispatch races with generation checks. Ensure connection replacement cannot create two valid active paths due to parallel handshakes.
- Separate transport errors, protocol errors, and gameplay outcomes. A transport/auth rejection never becomes an incorrect-pick penalty. W011/W012 own gameplay result mapping.
- Redact Cookie and Authorization headers; do not log raw frames or username payloads by default. Log stable error category, request ID only if non-sensitive, Room-scoped opaque ID, and coarse metrics. Avoid high-cardinality Guest Session labels.
- Do not enable permissive CORS with credentials for WebSocket. Apply WSS/TLS in production and securely configure origins, proxy trust, limits, and secrets.

## Test plan

- Go unit tests for strict envelope decoding, message-size/UTF-8 bounds, error mapping, Origin parser/policy, generation registry replacement, stale-generation rejection, idempotency collision/fingerprint behavior, queue saturation, and cleanup idempotence.
- Go integration tests use httptest and real WebSocket clients to validate upgrade/auth/origin ordering, multiple parallel connections for one Guest Session, cross-Room INV-004 enforcement, connection replacement with in-flight requests, Room/W013 disconnect signal, snapshots/sequence gaps, authorization, shutdown, and slow clients.
- Protocol contract tests add valid/invalid fixtures for all W009/W010 messages, strict unknown fields, bounds, errors, supported/unsupported versions, recipient-stream sequencing (including private-event non-delivery without artificial gaps), snapshot watermarks, and private-state leakage. Run `pnpm protocol:check`.
- Frontend tests for connect/reconnect/status transitions, correlation, sequence gap recovery, timeout/close/error mapping, bounded backoff, and no credential persistence. Browser E2E with multiple contexts verifies private Room Lobby synchronization and same-session tab replacement; gameplay actions are tested in W011.
- Race detector and load tests exercise replacement races, concurrent Room operations, queue saturation, and cleanup. Benchmark handshake/message/snapshot sizes before confirming frame limit or rate thresholds.
- Record per-criterion evidence as `PASS`, `FAIL`, `NOT_RUN`, `BLOCKED`, or reasoned `NOT_APPLICABLE`; distinguish developer, CI, and independent reviewer results.

## Dependencies and risks

- W005 protocol contract (DONE) is the schema/error baseline.
- W008 approved/merged Guest Session is required for cookie-authenticated upgrade and stable session principal; current approval is pending.
- W009 approved/merged Room manager is required for membership-checked routing, Room event sequences, state snapshots, and disconnect/presence APIs; its reviewed spec is pending owner decisions/merge.
- W011 consumes the Gateway dispatch seam for authoritative gameplay. W013 consumes connection loss/replacement signals for grace, forfeit, and recovery.
- Key risks: handshake/replacement races; server shutdown leaks; stale authority during socket close; event sequence gaps; strict protocol incompatibility; private snapshot leaks; idempotency cache memory/expiry; proxy Origin/host spoofing; slow clients; and unverified library behavior/limits.
- Technical decisions requiring T01 review: exact `coder/websocket` version/API; numeric frame/queue/connection limits; idempotency retention and key scope; heartbeat/read/write deadlines; cross-tab/Room replacement UX; Origin configuration; close/retry codes; and recipient-stream sequencing plus the related W005 v1 contract clarification. Any behavior that changes INV-004, authenticated session semantics, or the protocol contract requires Project Owner approval.

## Approval

- Project Owner / authorized approver: Pending
- Independent specification review: PASS at `ec39f546875428c56cfe32419f7b9be006f863c8`, recorded in `docs/development/reviews/W010-review.md`
- Decision and date: Pending
- W008/W009 dependency approval/merge: Pending
