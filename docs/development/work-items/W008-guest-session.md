# W008 — Guest Session

Status: READY_FOR_REVIEW
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

- Go package boundary: `apps/server/internal/session/` owns guest identity, username validation, credential generation/verifier, expiry, and a concurrency-safe in-memory session store. It does not depend on HTTP or Rooms.
- HTTP adapter proposal: `POST /api/v1/guest-session` creates a session with JSON `{ "username": "..." }`; `GET /api/v1/guest-session` recovers the session. Creation returns `201` with `{ "guestId", "username", "expiresAt" }`; the raw credential is cookie-only. Invalid usernames return `400 USERNAME_INVALID`; rate limits return `429 RATE_LIMITED` and `Retry-After`; missing/invalid sessions return `401 SESSION_REQUIRED`. Responses use `Cache-Control: no-store`. Reject unknown JSON fields, cap request bodies at 4 KiB, and reject malformed UTF-8. For POST require an exact configured `Origin`; for any request reject a present non-matching Origin or `Sec-Fetch-Site: cross-site`. Do not enable cross-origin credentials/CORS.
- React client: add a guest-entry bootstrap state and API client; on app startup recover session with credentials included; keep only display username in localStorage; keep no authentication secret in JS-readable storage.
- Integration seam: Room and WebSocket layers receive an authenticated `GuestSession` value from middleware/connection setup, never a caller-supplied guest ID.
- Proposed credential design: generate a separate 16-byte random guest ID and 32-byte random session credential with `crypto/rand`; encode both as unpadded base64url strings with distinct prefixes. Store `SHA-256(credential)` as the in-memory lookup key, never the raw credential. The 256-bit credential is high-entropy, so a fast one-way digest is adequate for verifier storage; do not add a password hashing function or pepper secret.
- Approved expiry: enforce configurable 24-hour idle and 7-day absolute deadlines; use the earlier. Count successful authenticated application actions and explicit session recovery as qualifying activity, but exclude WebSocket ping/pong, passive snapshots, background heartbeats, and failed requests. Inject a clock for deterministic tests. Expiration invalidates HTTP and WebSocket authorization consistently. If it happens during an active Match, revoke/close and immediately forfeit, without the distinct network-loss grace period. Cleanup must stop on shutdown.
- Approved production cookie: `__Host-prickle_guest`, Secure, HttpOnly, SameSite=Lax, Path `/`, no Domain, with Max-Age/Expires capped by the absolute expiration. Define a separate development-only HTTP cookie configuration; never weaken production settings. Use HTTPS/WSS in production, strict Origin validation, and per-message authorization.
- Approved session-creation limit: no more than 10 newly allocated sessions per 10 minutes per source IP, configurable and server-side. Reuse valid sessions on refresh. Return 429/Retry-After or the approved stable rate-limit mapping. Use only explicitly trusted proxy metadata; bound and clean limiter-key storage. Do not add an active-session cap without measured need and a separately approved policy.
- Approved username policy is specified in `docs/REQUIREMENT.md` §6.2. Enforce its 256-byte raw input bound before normalization. Grapheme counting must use extended grapheme clusters and preserve valid Thai combining sequences. W009 owns duplicate keys (NFC plus ASCII-English case-insensitive comparison).
- Grapheme dependency candidate: `github.com/clipperhouse/uax29/v2/graphemes` at v2.7.0. The package documents Unicode 17 grapheme rules, official Unicode test-suite use, and an MIT license; confirm the exact source/revision and dependency graph at implementation time. Do not install it until the specification and dependency review are approved.
- Start with process-local storage; no persistence or distributed session dependency.

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

The Project Owner approved the session lifetime, qualifying activity, credential, cookie, session/Room rate limits, and username policy recorded above on 2026-10-19. These are the initial configurable W008 baseline; they are no longer pending decisions. API shapes, dependency revisions, error codes, proxy trust configuration, and exact session-expiry-to-W013 event interface remain technical design details and require specification-level review before implementation. Public session creation remains unavailable until its limiter and tests are implemented.

Security references reviewed 2026-10-19: [OWASP Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html), [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html), and the official [`crypto/rand` documentation](https://pkg.go.dev/crypto/rand). The grapheme dependency candidate is documented at [pkg.go.dev](https://pkg.go.dev/github.com/clipperhouse/uax29/v2/graphemes) and its [version-pinned repository](https://github.com/clipperhouse/uax29/tree/v2.7.0).

## Dependencies and risks

- Depends on W004 standards/CI and W005 versioned protocol contract (DONE).
- W009 consumes authenticated guest identity and introduces per-Room username uniqueness. W010 consumes session authorization for WebSocket establishment.
- Username grapheme counting is not equivalent to Go rune count or UTF-16 length; dependency quality and test vectors matter.
- In-memory sessions disappear on restart and are single-server only, consistent with MVP architecture but a deliberate limitation.
- W009 must implement the approved Room rate limits and duplicate key policy; W013 must implement immediate Match forfeit on session expiry without network grace. The grapheme package candidate may be replaced after version/license/security review; the chosen dependency revision requires review before installation.

## Specification review history

Review 1 at `8391d663f7176e1fddd64829f00ca18ba69e062d`: FAIL with W008-SR-01 (MAJOR, username allowlist ambiguity) and W008-SR-02 (MINOR, no interim unauthenticated session-creation limit). Remediation clarifies acceptance of other safe Unicode text and adds a configurable rate-limit requirement and public-exposure gate.
- Review 2: PASS at specification revision `675f770f25583edd4bf9d149e53986b689172573`; both findings resolved, no new BLOCKER/MAJOR/MINOR findings. Reviewer recommends that T01 define unsafe Unicode/rendering test vectors without rejecting safe scripts. Full report: `docs/development/reviews/W008-review.md`.
- Review 3: PASS at specification revision `41be8e8f2d984151b2fbb185cd18271a5e4fb21d`; the expanded design had no BLOCKER/MAJOR findings and one MINOR follow-up to assess embedded bidi controls. Review 4 resolved the follow-up.
- Review 4: PASS at specification revision `97a6ff4f0c6b95420b76b45fb50ba9895dce2949`; this review predates the owner decisions below and does not apply to this revision.
- Review 5: PASS at specification and requirement revision `cd37fb4096b1da295d7736985979b382f008bcf6`; the reviewer confirmed that the approved lifetime, activity definition, credentials, cookie, session limit, username policy, and expiry-forfeit behavior are reflected consistently. No findings. Report: `docs/development/reviews/W008-review.md`. This does not authorize merge; implementation has not started.
- Project Owner approved the session and username policy baseline on 2026-10-19.

## Implementation validation evidence

Not run; implementation has not started.

## Completion

Require approved specification, implementation and tests satisfying W008-AC-01–10, actual CI pass, independent exact-SHA review PASS, blocking findings resolved, documentation synchronized, PR merged to `develop`, and feature branch retained.
