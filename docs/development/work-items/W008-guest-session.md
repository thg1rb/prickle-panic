# W008 — Guest Session

Status: READY_FOR_REVIEW
Type: Backend and Frontend Feature

## Goal

Create and recover an anonymous player identity without login, validate the player's display username on the server, and authorize future protected requests with an unpredictable server-issued session credential.

## Scope

- Add a Go Guest Session module with opaque guest identity, cryptographically generated session credentials, explicit expiry, and in-memory active session storage for the MVP process lifetime.
- Add HTTP session creation and recovery endpoints; creation accepts the initial username and sets the credential in an HttpOnly cookie. Recovery reads the cookie and returns only the authorized guest identity and display username.
- Validate usernames on the server: trim leading/trailing whitespace; require 2–16 Unicode grapheme clusters after trimming; accept Thai, English, numbers, and internal spaces as explicitly required without treating this list as an exclusive allowlist; accept other safe Unicode text too; reject empty values, unsafe control characters, and values unsafe for text rendering. Duplicate-in-Room checks belong to W009 when Room membership exists.
- Keep guest identity separate from username. Do not use username as a credential or expose the credential in JSON, logs, browser storage, URLs, or errors.
- Let the client remember only the last display username in localStorage, create a session after valid submission, and restore an existing session after refresh through the recovery endpoint. Add Thai UI copy for validation, expiry, and service failures.
- Define testable interfaces so W009 and W010 can authorize Room and WebSocket actions from the same session manager.
- Document that an in-memory session is lost on process restart; active session durability across server restarts is not required by this Work Item.

## Non-goals

- Accounts, email, passwords, social login, durable database storage, or cross-device identity.
- Room membership, same-Room duplicate username checks, Lobby username changes, host authority, or Room actions (W009).
- WebSocket connection management, single-active-connection enforcement, reconnection handoff, or gameplay authorization (W010/W013).
- Broader abuse controls and later rate-limit tuning; W013 owns those. W008 must include protective session-creation rate limiting before public endpoint exposure.
- Changes to confirmed username/session requirements or production session expiration policy.

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

- Apply the authoritative username policy in §6.2 exactly: 2–16 visible characters using Unicode grapheme clusters, Thai/English/numbers and internal spaces explicitly permitted, outer whitespace trimmed, no empty/whitespace-only username, and no unsafe control characters or harmful text rendering. The listed categories are not an exclusive allowlist; other safe Unicode text must not be rejected merely because it is not listed. Render usernames as text, never markup.
- Guest ID is unique and distinct from display username; usernames can repeat globally.
- Raw credentials are unpredictable, server-issued, and not stored in localStorage.
- Production cookie attributes are HttpOnly, Secure, explicit expiry, and an appropriate SameSite setting. Cookie lifetime and SameSite value remain subject to the documented policy decision below.
- An invalid session cannot authorize protected operations. No endpoint returns raw credentials or verifier material.

## Edge cases and error handling

- Missing, malformed, overlong, or invalid Unicode input: reject with a stable error code and no session mutation.
- A username of one grapheme (including a multi-codepoint grapheme), two graphemes, sixteen graphemes, and seventeen graphemes must be tested explicitly.
- Whitespace-only, leading/trailing whitespace, internal spaces, Thai combining marks, emoji grapheme sequences (including ZWJ), Arabic text, C0/C1 controls, and each listed bidi control must be tested. Bidi controls are proposed to be rejected; ZWJ must remain accepted.
- Credential collision: retry generation a bounded number of times, then fail safely without issuing a duplicate identity.
- Missing, expired, malformed, unknown, or revoked cookie: return unauthenticated without revealing which credential check failed.
- Session store capacity or entropy-source failure: fail closed and do not issue a partial session.
- Cross-origin session creation/recovery requests: reject disallowed origins; production deployments must use HTTPS and an explicit cookie policy. Exact CSRF and Origin behavior must be reviewed with the HTTP API design.
- Process restart loses in-memory sessions; the client recovers as unauthenticated and can create a new guest identity. This limitation must be documented to the player through ordinary UX only if needed; it is not a server error detail.
- Session creation is an unauthenticated sensitive operation. Apply a configurable protective rate limit at the HTTP edge that does not rely exclusively on a guest session ID; reject excess attempts without allocating sessions. The initial threshold and source signal are proposed security/tuning decisions for technical design. Do not expose the endpoint publicly until this limit and its tests are in place.

## Acceptance criteria

- **W008-AC-01:** A visitor can create a guest session with a valid username without login; server issues a unique guest ID distinct from username.
- **W008-AC-02:** Server applies the exact §6.2 trim, grapheme-count, and unsafe-control rules; Thai, English, numbers, and internal spaces are accepted, and other safe Unicode is not rejected solely by an allowlist. Boundary and Unicode tests cover valid/invalid examples, including proposed rejection of Bidi_Control characters while preserving Arabic text and ZWJ emoji.
- **W008-AC-03:** Invalid creation returns a stable error, creates no session, and sets no credential cookie.
- **W008-AC-04:** Session credential is cryptographically unpredictable, carried in an HttpOnly cookie, never returned in JSON, never stored in localStorage, URL, or logs, and never placed in WebSocket JSON.
- **W008-AC-05:** Recovery after refresh returns the same guest ID and username for a valid session; invalid/expired/missing credentials are uniformly unauthenticated.
- **W008-AC-06:** Cookie policy, expiry, credential hashing/storage, endpoint schema, and in-memory restart limitation are documented; no durable database is introduced.
- **W008-AC-07:** Session manager APIs are isolated from HTTP/transport and testable; authorization lookup distinguishes unauthenticated from valid session without leaking secrets.
- **W008-AC-08:** React client remembers only display username, calls create/recovery, restores or prompts appropriately, and renders all player-visible content in Thai.
- **W008-AC-09:** Tests cover session creation/recovery/expiry/revocation, collision/failure handling, cookie attributes, unauthenticated session-creation rate limiting, username Unicode boundaries, client refresh behavior, and no credential leakage; `pnpm check` and CI pass.
- **W008-AC-10:** Independent review PASS applies to the exact final PR commit; no BLOCKER/MAJOR findings remain, PR is merged to `develop`, feature branch retained, and evidence is recorded.

## Technical design

- Go package boundary: `apps/server/internal/session/` owns guest identity, username validation, credential generation/verifier, expiry, and a concurrency-safe in-memory session store. It does not depend on HTTP or Rooms.
- HTTP adapter proposal: `POST /api/v1/guest-session` creates a session with JSON `{ "username": "..." }`; `GET /api/v1/guest-session` recovers the session. Creation returns `201` with `{ "guestId", "username", "expiresAt" }`; the raw credential is cookie-only. Invalid usernames return `400 USERNAME_INVALID`; rate limits return `429 RATE_LIMITED` and `Retry-After`; missing/invalid sessions return `401 SESSION_REQUIRED`. Responses use `Cache-Control: no-store`. Reject unknown JSON fields, cap request bodies at 4 KiB, and reject malformed UTF-8. For POST require an exact configured `Origin`; for any request reject a present non-matching Origin or `Sec-Fetch-Site: cross-site`. Do not enable cross-origin credentials/CORS.
- React client: add a guest-entry bootstrap state and API client; on app startup recover session with credentials included; keep only display username in localStorage; keep no authentication secret in JS-readable storage.
- Integration seam: Room and WebSocket layers receive an authenticated `GuestSession` value from middleware/connection setup, never a caller-supplied guest ID.
- Proposed credential design: generate a separate 16-byte random guest ID and 32-byte random session credential with `crypto/rand`; encode both as unpadded base64url strings with distinct prefixes. Store `SHA-256(credential)` as the in-memory lookup key, never the raw credential. The 256-bit credential is high-entropy, so a fast one-way digest is adequate for verifier storage; do not add a password hashing function or pepper secret.
- Proposed expiry: fixed 24-hour absolute lifetime from creation; recovery does not extend it. Inject the clock into the session manager for deterministic expiry tests. Sweep expired sessions once per minute and stop the cleanup loop on server shutdown.
- Proposed production cookie: `__Host-pp_guest`, `Secure`, `HttpOnly`, `SameSite=Lax`, `Path=/`, no `Domain`, and `Max-Age=86400`/matching `Expires`. For local HTTP development only, use `pp_guest_dev` without `Secure`; keep it host-only, HttpOnly, SameSite=Lax, and Path `/`. Configure Vite to proxy `/api` to Go so the browser uses same-origin requests; never enable the development cookie in production.
- Proposed session-creation limiter: configurable in-memory rolling 60-second window, default 10 attempts per source IP, applied before request-body parsing and session allocation. Return `429` with `Retry-After`. Use `RemoteAddr` unless the immediate peer is in an explicit trusted-proxy list; only then parse forwarding headers from right to left, ignoring untrusted client-supplied prefixes. Bound limiter-key storage and prune stale entries. Add an active-session capacity of 10,000 by default; reject creation when full without allocating partial state. Initial thresholds/capacity require owner approval and load/abuse review.
- Username validation proposal: trim Go Unicode whitespace, reject invalid UTF-8, C0/C1 controls, and Unicode `Bidi_Control` characters (including U+061C, U+200E–U+200F, U+202A–U+202E, and U+2066–U+2069), then count extended grapheme clusters. Preserve all otherwise-safe Unicode, including Arabic text, ZWJ (U+200D) emoji sequences, and other safe scripts; render usernames as escaped text inside a bidirectional isolation element (`<bdi dir="auto">`) in the web UI. Do not treat markup characters as executable content or reject them solely by spelling. This is a proposed safety policy requiring Project Owner approval.
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

- Go unit tests for username grapheme policy, unique guest IDs, credential generation, expiry/revocation, capacity/failure behavior, limiter boundary/source parsing, and race-safe concurrent access. Include known UAX #29 vectors for Thai combining sequences, emoji joined by ZWJ, regional-indicator flags, Arabic, and Unicode 17 test data where useful; invalid UTF-8, C0/C1 controls, and listed Bidi_Control code points reject under the proposed policy. Confirm U+200D ZWJ remains accepted. Literal `<script>`-like text remains accepted as a username and must render inertly as text.
- HTTP tests for endpoint status/error mapping, 4 KiB request-body limit, cookie flags, no-cookie-on-error, Origin and Fetch Metadata behavior, limiter before body parsing/session allocation, expired-cookie clearing, and credential absence from response bodies/log hooks.
- Frontend tests for username input, Thai errors, create/recovery bootstrap, refresh restoration, and localStorage containing display name only.
- Security checks for dependency audit and static analysis; `pnpm check` and exact-head CI. Verify the session-creation limiter before authentication and confirm excess attempts allocate no sessions.
- Independent review inspects the implementation diff, specs, tests, actual results, and final SHA.

## T01 design proposal and approval

The following defaults are proposals, not confirmed product requirements. Project Owner approval is required before implementation: 24-hour absolute session expiry; 256-bit cookie credential and SHA-256 verifier; production `__Host-pp_guest` Secure/HttpOnly/SameSite=Lax cookie; 10 session-creation attempts per minute per trusted source IP; 10,000 active sessions per server process; and the exact control-character policy. The API shapes, Unicode 17 grapheme implementation candidate, error codes, and local development proxy are technical recommendations subject to independent design review. Public session creation remains unavailable until the limiter and capacity behavior are implemented and tested.

Security references reviewed 2026-10-19: [OWASP Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html), [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html), and the official [`crypto/rand` documentation](https://pkg.go.dev/crypto/rand). The grapheme dependency candidate is documented at [pkg.go.dev](https://pkg.go.dev/github.com/clipperhouse/uax29/v2/graphemes) and its [version-pinned repository](https://github.com/clipperhouse/uax29/tree/v2.7.0).

## Dependencies and risks

- Depends on W004 standards/CI and W005 versioned protocol contract (DONE).
- W009 consumes authenticated guest identity and introduces per-Room username uniqueness. W010 consumes session authorization for WebSocket establishment.
- Username grapheme counting is not equivalent to Go rune count or UTF-16 length; dependency quality and test vectors matter.
- In-memory sessions disappear on restart and are single-server only, consistent with MVP architecture but a deliberate limitation.
- Approval is still required for the proposed cookie lifetime, SameSite, credential format, rate limit/capacity defaults, and Unicode control policy before implementation. The package candidate may be replaced after version/license/security review.

## Specification review history

Review 1 at `8391d663f7176e1fddd64829f00ca18ba69e062d`: FAIL with W008-SR-01 (MAJOR, username allowlist ambiguity) and W008-SR-02 (MINOR, no interim unauthenticated session-creation limit). Remediation clarifies acceptance of other safe Unicode text and adds a configurable rate-limit requirement and public-exposure gate.
- Review 2: PASS at specification revision `675f770f25583edd4bf9d149e53986b689172573`; both findings resolved, no new BLOCKER/MAJOR/MINOR findings. Reviewer recommends that T01 define unsafe Unicode/rendering test vectors without rejecting safe scripts. Full report: `docs/development/reviews/W008-review.md`.
- Review 3: PASS at specification revision `41be8e8f2d984151b2fbb185cd18271a5e4fb21d`; the expanded design had no BLOCKER/MAJOR findings and one MINOR follow-up to assess embedded bidi controls. Review 4 is pending after adding that proposed policy and its test vectors.
- Project Owner approval of the specification and proposed session policy is pending; implementation has not started.

## Implementation validation evidence

Not run; implementation has not started.

## Completion

Require approved specification, implementation and tests satisfying W008-AC-01–10, actual CI pass, independent exact-SHA review PASS, blocking findings resolved, documentation synchronized, PR merged to `develop`, and feature branch retained.
