# W008 — Guest Session

Status: SPEC_READY
Type: Backend and Frontend Feature

## Goal

Create and recover an anonymous player identity without login, validate the player's display username on the server, and authorize future protected requests with an unpredictable server-issued session credential.

## Scope

- Add a Go Guest Session module with opaque guest identity, cryptographically generated session credentials, explicit expiry, and in-memory active session storage for the MVP process lifetime.
- Add HTTP session creation and recovery endpoints; creation accepts the initial username and sets the credential in an HttpOnly cookie. Recovery reads the cookie and returns only the authorized guest identity and display username.
- Validate usernames on the server using `docs/REQUIREMENT.md` §6.2: trim surrounding whitespace, normalize to NFC, enforce 2–16 Unicode extended grapheme clusters, a 256-byte raw-input limit, and the documented strict Thai/English/digit/single-ASCII-space character policy. Duplicate-in-Room comparison belongs to W009.
- Keep guest identity separate from username. Do not use username as a credential or expose the credential in JSON, logs, browser storage, URLs, or errors.
- Let the client remember only the last display username in localStorage, create a session after valid submission, and restore an existing session after refresh through the recovery endpoint. Add Thai UI copy for validation, expiry, and service failures.
- Define testable interfaces so W009 and W010 can authorize Room and WebSocket actions from the same session manager.
- Document that an in-memory session is lost on process restart; active session durability across server restarts is not required by this Work Item.

## Non-goals

- Accounts, email, passwords, social login, durable database storage, or cross-device identity.
- Room membership, same-Room duplicate username checks, Lobby username changes, host authority, or Room actions (W009).
- WebSocket connection management, single-active-connection enforcement, reconnection handoff, or gameplay authorization (W010/W013).
- Broader abuse controls and later rate-limit tuning; W013 owns those. W008 must include protective session-creation rate limiting before public endpoint exposure.
- Room membership, duplicate-name comparison implementation, Lobby username changes, host authority, or Room actions (W009).

## Requirement references

- `docs/REQUIREMENT.md` §§6.1–6.7, §16.1 (FR-001–008), §28 (AC-001), §32 (W008).
- NFR and security requirements: `docs/REQUIREMENT.md` §§17, 20, 21, 25, and 27 as applicable to sessions, privacy, HTTP security, validation, and logging.
- Protocol integration boundary: `docs/protocol/v1/README.md` states WebSocket credentials are not sent in JSON and the HTTP session contract is defined in W008.

## Expected behavior

1. A first-time visitor can enter a username without creating an account.
2. The client submits the entered username; the server trims leading/trailing Unicode whitespace, validates grapheme count, and returns the canonical display username. The client remembers only that returned value.
3. On valid creation, the server generates a unique opaque guest ID and a cryptographically unpredictable credential, stores only a non-reversible credential verifier, sets the raw credential in an HttpOnly cookie, and returns the guest ID, canonical display username, and expiry.
4. On invalid input, creation fails with a stable machine-readable error; no session or cookie is issued. Expired or invalid recovery clears the stale cookie and returns the same unauthenticated error.
5. On refresh, the client calls recovery; a valid credential restores the same guest identity and username, while missing, expired, or invalid credentials result in an unauthenticated state and Thai prompt to enter a username again.
6. The credential is sent only by cookie over HTTPS in production and is never exposed to JavaScript or included in protocol messages.

## Business rules

- Apply §6.2 exactly: trim surrounding whitespace, NFC-normalize, enforce 2–16 extended grapheme clusters and a 256-byte raw-input cap; permit Thai letters with valid associated marks, English letters, ASCII/Thai digits, and single internal ASCII spaces only. Reject controls, invisible/format/zero-width/bidirectional controls, standalone or malformed combining marks, emoji, unsupported scripts/punctuation, markup, and overlong input. Render as text, never markup. W009 performs same-Room duplicate comparison using an NFC key with ASCII-English case-insensitivity and preserves display casing.
- Guest ID is unique and distinct from display username; usernames can repeat globally.
- Raw credentials are unpredictable, server-issued, and not stored in localStorage.
- Session idle timeout is 24 hours after the last qualifying successful authenticated application activity or explicit recovery; absolute lifetime is 7 days from creation. The earlier deadline expires the session. Ping/pong, passive snapshots, background heartbeats, and failed requests do not qualify. Both values are configurable and enforced server-side across HTTP and WebSocket authorization. Session expiry during a Match revokes authorization and causes immediate forfeit, without the separate 15-second network reconnection grace.
- Production cookie is `__Host-prickle_guest; HttpOnly; Secure; SameSite=Lax; Path=/`, without Domain, with Max-Age no greater than the 7-day absolute lifetime. Local HTTP development uses a separate explicit development-only configuration. Production requires HTTPS/WSS, strict WebSocket Origin validation, and message-level authorization.
- Credentials are opaque CSPRNG tokens with at least 256 bits of randomness. Store a cryptographic hash server-side when practical; support invalidation; never log or expose raw credentials in URLs, localStorage, JSON, or error details; do not rotate on WebSocket reconnect.
- An invalid session cannot authorize protected operations. No endpoint returns raw credentials or verifier material.

## Edge cases and error handling

- Missing, malformed, overlong, or invalid Unicode input: reject with a stable error code and no session mutation.
- Test one, two, sixteen, and seventeen grapheme boundaries; NFC-equivalent forms; Thai consonants, vowel marks and tone marks; ASCII and Thai digits; English case equivalence for Room duplicate keys; outer trimming, one internal ASCII space, and repeated-space rejection; 256-byte input boundary; controls, invisible/format/zero-width/bidirectional controls, malformed combining sequences, emoji, unsupported scripts and punctuation, and markup. Valid Thai combining sequences must remain accepted.
- Credential collision: retry generation a bounded number of times, then fail safely without issuing a duplicate identity.
- Missing, expired, malformed, unknown, or revoked cookie: return unauthenticated without revealing which credential check failed.
- Entropy-source failure or bounded credential-collision retries: fail closed and do not issue a partial session. No arbitrary active-session capacity is introduced; all stores and rate-limit state must remain bounded with cleanup.
- Cross-origin session creation/recovery requests: reject disallowed origins; production deployments must use HTTPS and an explicit cookie policy. Exact CSRF and Origin behavior must be reviewed with the HTTP API design.
- Process restart loses in-memory sessions; the client recovers as unauthenticated and can create a new guest identity. This limitation must be documented to the player through ordinary UX only if needed; it is not a server error detail.
- Session creation is limited server-side to 10 new sessions per 10 minutes per source IP as an initial configurable baseline. Reuse an existing valid session on refresh. Return an appropriate rate-limit response without allocating state. Configure trusted-proxy handling explicitly, bound limiter memory, and assess shared-network fairness and network-source bypass during abuse testing. W009 applies 5 room-create attempts per 10 minutes/session and 20 joins/minute/session; additional source controls may mitigate session churn.

## Acceptance criteria

- **W008-AC-01:** A visitor can create a guest session with a valid username without login; server issues a unique guest ID distinct from username.
- **W008-AC-02:** Server applies §6.2 exactly, including trimming, NFC, 2–16 extended grapheme clusters, 256-byte raw limit, strict allowed character categories and rejection rules. Tests cover Thai valid marks, normalization equivalence, English case-insensitive duplicate keys (W009), both digit sets, whitespace, grapheme boundaries, malformed/invisible/bidi/control input, emoji, unsupported characters and markup.
- **W008-AC-03:** Invalid creation returns a stable error, creates no session, and sets no credential cookie.
- **W008-AC-04:** Session credential is cryptographically unpredictable, carried in an HttpOnly cookie, never returned in JSON, never stored in localStorage, URL, or logs, and never placed in WebSocket JSON.
- **W008-AC-05:** Recovery after refresh returns the same guest ID and username for a valid session; invalid/expired/missing credentials are uniformly unauthenticated.
- **W008-AC-06:** Approved cookie, 24-hour idle/7-day absolute expiry, qualifying-activity definition, credential hashing/storage, endpoint schema, and in-memory restart limitation are documented; no durable database is introduced.
- **W008-AC-07:** Session manager APIs are isolated from HTTP/transport and testable; authorization lookup distinguishes unauthenticated from valid session without leaking secrets.
- **W008-AC-08:** React client remembers only display username, calls create/recovery, restores or prompts appropriately, and renders all player-visible content in Thai.
- **W008-AC-09:** Tests cover idle and absolute expiry boundaries, qualifying/non-qualifying activity, HTTP/WS authorization revocation, immediate active-Match forfeit on expiry, credential generation/hash/revocation/no rotation, cookie attributes and local-dev separation, 10-per-10-minute/IP session creation limits and reuse, Unicode cases listed in §6.2, refresh behavior, and no credential leakage; `pnpm check` and CI pass.
- **W008-AC-10:** Independent review PASS applies to the exact final PR commit; no BLOCKER/MAJOR findings remain, PR is merged to `develop`, feature branch retained, and evidence is recorded.

## Technical design

### Boundaries and HTTP contract

- `apps/server/internal/session/` owns guest identity, username validation, credential generation and verification, expiry, and a concurrency-safe in-memory store. It has no dependency on HTTP, WebSocket, Rooms, or Match packages. HTTP handlers translate transport input into session-manager operations; application composition wires lifecycle event sinks.
- `POST /api/v1/guest-session` accepts exactly `{ "username": "..." }`. A new session returns `201 Created`, sets the credential cookie, and returns `{ "guestId": "...", "username": "...", "expiresAt": "<RFC3339Nano UTC>" }`. `expiresAt` is the earlier idle/absolute deadline. If the cookie carries a valid session, reuse it without renaming, allocation, or quota use; return `200 OK` with the same response shape and cookie Max-Age limited to remaining absolute lifetime. `GET /api/v1/guest-session` recovers a valid cookie and returns `200 OK`; successful recovery updates idle activity. Missing, malformed, expired, or revoked credentials all return the same `401 SESSION_REQUIRED` and clear a stale cookie when present. There is no public DELETE endpoint; revocation is an internal operation.
- All responses on the `/api` path set `Cache-Control: no-store`. Cap request bodies at 4 KiB before decoding; require valid UTF-8, one JSON object with exactly one `username` field, no duplicate/unknown fields, and no trailing values. Map invalid username to `400 USERNAME_INVALID`; malformed structure to `400 REQUEST_INVALID`; oversized body to `413 REQUEST_TOO_LARGE`; disallowed Origin to `403 ORIGIN_NOT_ALLOWED`; quota exhaustion to `429 RATE_LIMITED` with `Retry-After`; unavailable session service or unusable trusted-proxy metadata to `503 SESSION_UNAVAILABLE`. Use `{ "error": { "code": "..." } }`; errors contain no credential/verifier details. Do not set a credential cookie on errors except to clear a stale recovery cookie.
- For POST require exact configured same-origin `Origin`. For both endpoints reject a present non-matching Origin or `Sec-Fetch-Site: cross-site`. Missing Fetch Metadata does not bypass POST Origin checks. Do not enable credentialed CORS. The client uses same-origin requests with credentials included. No endpoint returns raw credentials or verifier material.

### Credential, session, and cookie lifecycle

- Generate a separate 16-byte random guest ID and 32-byte (256-bit) session credential with Go `crypto/rand`; encode both as unpadded base64url with distinct type prefixes. Retry collisions a bounded number of times, then fail closed. Store `SHA-256(credential)` as the in-memory lookup key; never retain the raw token after setting its cookie. The 256-bit credential is uniform high-entropy data, so a fast digest is adequate; do not add password hashing or a pepper.
- Record creation and qualifying activity using one injected server `Clock` that returns paired UTC wall time and a monotonic process-local tick. Session expiry and all Match deadlines must use that same monotonic clock domain. Effective expiry is `min(lastQualifyingActivity + idleTTL, createdAt + absoluteTTL)`, initially configurable at 24 hours and 7 days respectively; add configured durations to monotonic ticks for enforcement and retain paired UTC values for cookie dates and diagnostics. Wall-clock adjustments must not move an already established enforcement deadline. Successful authenticated application actions and explicit recovery qualify, only after success; failed actions, ping/pong, passive snapshots, and background heartbeats do not. Keep authorization lookup separate from the explicit post-success activity update so a failed action cannot extend idle expiry. Serialize activity updates with expiry checks so concurrent work cannot revive an expired session or extend absolute lifetime. Check expiry on every authorization lookup.
- A session-manager-owned deadline scheduler proactively expires sessions even when the player has no connected socket and sends no further request. Maintain one indexed deadline entry per active session, updating it atomically when qualifying activity changes the idle deadline; scheduler memory therefore tracks the existing active-session store rather than an unbounded stale-timer queue. Wait until the earliest effective deadline, re-check it under session-manager synchronization, revoke at/after the deadline, and emit the expiry event after revocation. On shutdown cancel and join the scheduler and cleanup workers. Process restart loses all in-memory sessions.
- Production cookie: `__Host-prickle_guest`, HttpOnly, Secure, SameSite=Lax, Path `/`, no Domain. Max-Age/Expires is bounded by the remaining absolute lifetime (at most seven days); server-side idle expiry is authoritative. Production startup rejects any configuration that disables Secure and requires HTTPS/WSS. Local HTTP uses an explicitly selected development-only cookie name/configuration isolated from production and rejected by production startup validation. The T03 local configuration is `APP_ENV=development`, loopback-only `APP_ORIGIN`, and cookie name `prickle_guest_dev`; Secure is disabled only in this explicit local mode.

### Trusted client source and rate limits

- A boundary `ClientIPResolver` uses parsed `RemoteAddr` and ignores forwarded headers if no trusted proxy CIDRs are configured. If the direct peer is in an explicitly configured trusted CIDR, parse `X-Forwarded-For` and walk right-to-left through trusted hops, selecting the first untrusted address. The configurable hop limit defaults to 32. Never trust forwarded headers from an untrusted peer, arbitrary `Forwarded`, or an unchecked leftmost XFF value. Invalid/ambiguous metadata from a trusted peer fails closed for creation with `503 SESSION_UNAVAILABLE`; validate proxy CIDRs at startup against the actual topology.
- A server-side rolling 10-minute window admits at most 10 successfully allocated sessions per canonical source IP. Only successful new allocations count. Invalid requests and username failures do not consume quota; valid-session reuse consumes no new-session quota. Quota check and allocation are atomic. HTTP may use a non-consuming Manager preflight before reading a bounded body; allocation repeats the quota check atomically. `429 RATE_LIMITED` includes `Retry-After`, the ceiling seconds until the oldest successful allocation leaves the window. In-memory counters expire and cleanup; tracked source keys have an explicit configurable bound (implementation default 65,536), and capacity exhaustion fails closed with 503 rather than evicting live counters. Shared-network fairness is a tradeoff; IP is only a quota source, never identity. The initial threshold is configurable but is not changed without owner approval.
- W008 owns only session-creation limiting. W009 owns 5 Room-create attempts per 10 minutes and 20 Room joins per minute per Guest Session. W008 does not implement Room operations or their counters.

### Username validation and dependency selection

- Validation order is normative: (1) require valid UTF-8 in the bounded 4 KiB JSON body; (2) reject username raw UTF-8 input above 256 bytes before normalization; (3) reject control, format, zero-width, and bidirectional control characters before trimming, so tabs/newlines cannot be silently removed; (4) trim surrounding Unicode whitespace; (5) normalize to NFC; (6) reject emoji, unsupported scripts/punctuation, markup, and malformed or standalone combining marks while permitting valid Thai consonant/vowel/tone sequences; (7) allow only Thai/English letters, ASCII/Thai digits, and single internal ASCII spaces; (8) require 2–16 Unicode extended grapheme clusters. Preserve normalized display casing. Validation is server-authoritative. W009 derives its duplicate-in-Room key from that display value with ASCII-English case folding only.
- T02 pins module `golang.org/x/text` v0.43.0, package `golang.org/x/text/unicode/norm` (BSD-3-Clause; Go project source; Go 1.26 minimum) for NFC; and module `github.com/clipperhouse/uax29/v2` v2.7.0, package `github.com/clipperhouse/uax29/v2/graphemes` (MIT; Unicode 17 UAX #29 rules and conformance suite) for graphemes. T02 records the module graph, license notices, and vulnerability check in the implementation report. T03 adds no backend framework or third-party dependency.

### Expiry integration contract with W010 and W013

- W008 owns the authoritative validity decision and defines an internal `ExpiryEvent` sink contract at the session/application boundary. It carries stable event ID, opaque guest ID, reason (`IDLE_TIMEOUT` or `ABSOLUTE_TIMEOUT`), `ExpiredAtUTC` for diagnostics, and `DeadlineTick` from the shared monotonic clock domain for ordering; no credential, hash, username, or cookie. The tick is process-local, is not persisted, and is meaningful only alongside Match deadline ticks produced by the same injected clock. Session code depends only on the sink interface and imports neither W010 nor W013.
- Application composition fans each event independently to W010 and W013. W010 invalidates/closes every established socket for the guest with internal reason `SESSION_EXPIRED`; it must not classify this as network loss or begin reconnection grace. W013 enqueues the event into the owning Room event loop and immediately applies active-Match forfeit. This works without a currently connected socket. The 15-second grace is only for unexpected network disconnection.
- Revoke authorization before notification. Delivery is at-least-once; consumers deduplicate by event ID and handle idempotently. The session manager retains failed deliveries for retry in a bounded in-memory outbox with configurable capacity (initial default 1,024) and capped exponential backoff from 100 ms to 30 seconds. Do not drop an event when the outbox is full: retain its pending-delivery state on the expired session record, reject new session creation and authenticated lookups with a typed backpressure error until capacity is available, and promote pending events as delivery drains. Expired sessions are always unauthorized, regardless of sink health. Retry the same stable event ID because a sink can fail after accepting an event; consumers deduplicate. Outbox/session state is process-local and is not recovered after restart.
- Expiry strictly before Match deadline forfeits; equal timestamps are won by the Match deadline, as explicitly approved by the Project Owner on 2026-10-19; already-terminal matches remain unchanged. Delayed delivery compares `DeadlineTick`, not event arrival time or UTC wall-clock values. The monotonic tick is the common ordering basis for session expiry and Match deadlines; UTC is diagnostic only, so wall-clock jumps do not alter the winner. Repeated notifications cannot duplicate forfeit/results. Public HTTP continues to return uniform `401 SESSION_REQUIRED`; `SESSION_EXPIRED` is an internal transport/Room reason, not a new credential detail. W013's open specification currently records equal-time ordering as requiring a decision; this approval evidence should accompany W008 T01 integration, and W013's integration spec must be synchronized before its implementation begins.
- W008 tests event shape, revoke-before-notify ordering, proactive expiry with no socket or authorization request, duplicate delivery, monotonic ordering, wall-clock adjustments, and deadline tie behavior. W010 tests socket closure and no network-grace classification. W013 tests Room-loop forfeit ordering and deadline ties. This design does not rewrite their specifications or wire protocol.

### Test plan

- T02 unit tests: fake-clock idle/absolute boundaries; activity qualification and immutable absolute deadline; revoke/lookup races; credential generation/encoding/hash/revocation/no rotation, collision and entropy failures; rolling limiter admission, Retry-After, concurrent atomicity, cleanup and bounded-key exhaustion; trusted/untrusted/direct/malformed proxy chains and spoofing; Thai/English Unicode vectors, NFC equivalence, grapheme boundaries, malformed marks, controls, bidi, zero-width, emoji, unsupported scripts, markup and excess length.
- T03 HTTP tests: creation/reuse status and quota behavior; GET recovery and uniform unauthorized mapping; strict/oversized JSON; cookie attributes/absolute Max-Age/local-only separation; Origin and Fetch Metadata; no cookie on errors except stale-cookie clearing; 429 Retry-After; absence of secrets from response, URL, logs and errors; unavailable-service mapping.
- T01 validation is manual traceability and independent design review against W008 acceptance criteria and requirement IDs. These are design checks only; no implementation or implementation-test result is claimed.

## Security and privacy

- Require TLS for production; set `Secure`, `HttpOnly`, explicit `Max-Age`/expiry, and a reviewed `SameSite` value.
- Validate Origin on session-mutating endpoints and use same-origin deployment; do not rely on SameSite alone as CSRF protection.
- Bound request body, username bytes, session allocation, credential generation retries, and unauthenticated session-creation attempts. The limiter must not rely exclusively on guest session IDs; document its source signal, proxy-trust assumptions, and spoofing/availability tradeoffs.
- Exclude cookies, credential verifiers, and raw secrets from application logs and error responses; redact sensitive request headers in observability.
- Do not interpret usernames as HTML. Store/display as text and reject unsafe control characters according to the source requirement.
- Avoid claiming sessions survive process restart; no account recovery is implied.

## Test plan

- Go unit tests for username grapheme policy, unique guest IDs, CSPRNG credential generation/hash/revocation, idle and absolute expiry, activity qualification, authorization invalidation, limiter boundary/source parsing, and race-safe concurrent access. Include UAX #29 tests for Thai vowel/tone marks and normalization-equivalent forms; ASCII/Thai digits; one internal space and repeated-space rejection; malformed/invisible/bidi controls, emoji, unsupported scripts/punctuation, markup, and byte/grapheme boundaries. Verify safe text rendering. Do not include active-session capacity tests absent an approved cap.
- HTTP tests for endpoint status/error mapping, 4 KiB request-body limit, cookie flags, no-cookie-on-error, Origin and Fetch Metadata behavior, limiter before body parsing/session allocation, expired-cookie clearing, and credential absence from response bodies/log hooks.
- Frontend tests for username input, Thai errors, create/recovery bootstrap, refresh restoration, and localStorage containing display name only.
- Security checks for dependency audit and static analysis; `pnpm check` and exact-head CI. Verify the session-creation limiter before authentication and confirm excess attempts allocate no sessions.
- Independent review inspects the implementation diff, specs, tests, actual results, and final SHA.

## Owner decisions and implementation gate

The Project Owner approved the session lifetime, qualifying activity, credential, cookie, session/Room rate limits, username policy, successful-allocation-only session quota accounting, and Match-deadline tie behavior (the Match deadline wins at equal timestamps) on 2026-10-19. These are the initial configurable W008 baseline; they are no longer pending decisions. The API, dependency revisions, error codes, trusted proxy algorithm, expiry scheduler, and session-expiry event interface are finalized here as T01 design choices, subject to independent design review. Public session creation remains unavailable until its limiter and tests are implemented.

Security references reviewed 2026-10-19: [OWASP Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html), [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html), and the official [`crypto/rand` documentation](https://pkg.go.dev/crypto/rand). The grapheme dependency candidate is documented at [pkg.go.dev](https://pkg.go.dev/github.com/clipperhouse/uax29/v2/graphemes) and its [version-pinned repository](https://github.com/clipperhouse/uax29/tree/v2.7.0).

## Dependencies and risks

- Depends on W004 standards/CI and W005 versioned protocol contract (DONE).
- W009 consumes authenticated guest identity and introduces per-Room username uniqueness. W010 consumes session authorization for WebSocket establishment.
- Username grapheme counting is not equivalent to Go rune count or UTF-16 length; dependency quality and test vectors matter.
- In-memory sessions disappear on restart and are single-server only, consistent with MVP architecture but a deliberate limitation.
- W009 must implement the approved Room rate limits and duplicate key policy; W013 must implement immediate Match forfeit on session expiry without network grace. T02 pins the reviewed NFC and grapheme dependency revisions in `apps/server/go.mod`; the implementation report records their license and vulnerability-check evidence.

## Specification review history

Review 1 at `8391d663f7176e1fddd64829f00ca18ba69e062d`: FAIL with W008-SR-01 (MAJOR, username allowlist ambiguity) and W008-SR-02 (MINOR, no interim unauthenticated session-creation limit). Remediation clarifies acceptance of other safe Unicode text and adds a configurable rate-limit requirement and public-exposure gate.
- Review 2: PASS at specification revision `675f770f25583edd4bf9d149e53986b689172573`; both findings resolved, no new BLOCKER/MAJOR/MINOR findings. Reviewer recommends that T01 define unsafe Unicode/rendering test vectors without rejecting safe scripts. Full report: `docs/development/reviews/W008-review.md`.
- Review 3: PASS at specification revision `41be8e8f2d984151b2fbb185cd18271a5e4fb21d`; the expanded design had no BLOCKER/MAJOR findings and one MINOR follow-up to assess embedded bidi controls. Review 4 resolved the follow-up.
- Review 4: PASS at specification revision `97a6ff4f0c6b95420b76b45fb50ba9895dce2949`; this review predates the owner decisions below and does not apply to this revision.
- Review 5: PASS at specification and requirement revision `cd37fb4096b1da295d7736985979b382f008bcf6`; the reviewer confirmed that the approved lifetime, activity definition, credentials, cookie, session limit, username policy, and expiry-forfeit behavior are reflected consistently. No findings. Report: `docs/development/reviews/W008-review.md`. This does not authorize merge; implementation has not started.
- Project Owner approved the session and username policy baseline on 2026-10-19.

## Implementation validation evidence

T02 was independently reviewed and integrated by PR #19 Merge Commit `c319fb78a6a9c2438738c2007ba0d9e11603e436`; its feature branch is retained. T03 implementation and exact-revision validation/review evidence are recorded in `docs/development/implementation-reports/W008.md`. W008 remains incomplete until T03, T04, and T05 are completed and integrated.

## Completion

Require approved specification, implementation and tests satisfying W008-AC-01–10, actual CI pass, independent exact-SHA review PASS, blocking findings resolved, documentation synchronized, PR merged to `develop`, and feature branch retained.
