# W008 — Guest Session

Status: SPEC_READY
Type: Backend and Frontend Feature

## Goal

Create and recover an anonymous player identity without login, validate the player's display username on the server, and authorize future protected requests with an unpredictable server-issued session credential.

## Scope

- Add a Go Guest Session module with opaque guest identity, cryptographically generated session credentials, explicit expiry, and in-memory active session storage for the MVP process lifetime.
- Add HTTP session creation and recovery endpoints; creation accepts the initial username and sets the credential in an HttpOnly cookie. Recovery reads the cookie and returns only the authorized guest identity and display username.
- Validate usernames on the server: trim leading/trailing whitespace; require 2–16 Unicode grapheme clusters after trimming; permit Thai, English, numbers, and internal spaces; reject empty values, control characters, and values unsafe for text rendering. Duplicate-in-Room checks belong to W009 when Room membership exists.
- Keep guest identity separate from username. Do not use username as a credential or expose the credential in JSON, logs, browser storage, URLs, or errors.
- Let the client remember only the last display username in localStorage, create a session after valid submission, and restore an existing session after refresh through the recovery endpoint. Add Thai UI copy for validation, expiry, and service failures.
- Define testable interfaces so W009 and W010 can authorize Room and WebSocket actions from the same session manager.
- Document that an in-memory session is lost on process restart; active session durability across server restarts is not required by this Work Item.

## Non-goals

- Accounts, email, passwords, social login, durable database storage, or cross-device identity.
- Room membership, same-Room duplicate username checks, Lobby username changes, host authority, or Room actions (W009).
- WebSocket connection management, single-active-connection enforcement, reconnection handoff, or gameplay authorization (W010/W013).
- Rate limiting and broader abuse controls beyond safe request bounds and session-creation resource limits; W013 owns policy and enforcement.
- Changes to confirmed username/session requirements or production session expiration policy.

## Requirement references

- `docs/REQUIREMENT.md` §§6.1–6.7, §16.1 (FR-001–008), §28 (AC-001), §32 (W008).
- NFR and security requirements: `docs/REQUIREMENT.md` §§17, 20, 21, 25, and 27 as applicable to sessions, privacy, HTTP security, validation, and logging.
- Protocol integration boundary: `docs/protocol/v1/README.md` states WebSocket credentials are not sent in JSON and the HTTP session contract is defined in W008.

## Expected behavior

1. A first-time visitor can enter a username without creating an account.
2. The client normalizes the submitted value only by trimming leading/trailing whitespace; server validation is authoritative and counts grapheme clusters.
3. On valid creation, the server generates a unique opaque guest ID and a cryptographically unpredictable credential, stores only a non-reversible credential verifier where practical, sets the raw credential in an HttpOnly cookie, and returns the guest ID and display username.
4. On invalid input, creation fails with a stable machine-readable error; no session or cookie is issued.
5. On refresh, the client calls recovery; a valid credential restores the same guest identity and username, while missing, expired, or invalid credentials result in an unauthenticated state and Thai prompt to enter a username again.
6. The credential is sent only by cookie over HTTPS in production and is never exposed to JavaScript or included in protocol messages.

## Business rules

- Apply the authoritative username policy in §6.2 exactly: 2–16 visible characters using Unicode grapheme clusters, Thai/English/numbers permitted, internal spaces permitted, outer whitespace trimmed, no empty/whitespace-only username, and no unsafe control characters or harmful text rendering.
- Guest ID is unique and distinct from display username; usernames can repeat globally.
- Raw credentials are unpredictable, server-issued, and not stored in localStorage.
- Production cookie attributes are HttpOnly, Secure, explicit expiry, and an appropriate SameSite setting. Cookie lifetime and SameSite value remain subject to the documented policy decision below.
- An invalid session cannot authorize protected operations. No endpoint returns raw credentials or verifier material.

## Edge cases and error handling

- Missing, malformed, overlong, or invalid Unicode input: reject with a stable error code and no session mutation.
- A username of one grapheme (including a multi-codepoint grapheme), two graphemes, sixteen graphemes, and seventeen graphemes must be tested explicitly.
- Whitespace-only, leading/trailing whitespace, internal spaces, Thai combining marks, emoji grapheme sequences, and control characters must be tested.
- Credential collision: retry generation a bounded number of times, then fail safely without issuing a duplicate identity.
- Missing, expired, malformed, unknown, or revoked cookie: return unauthenticated without revealing which credential check failed.
- Session store capacity or entropy-source failure: fail closed and do not issue a partial session.
- Cross-origin session creation/recovery requests: reject disallowed origins; production deployments must use HTTPS and an explicit cookie policy. Exact CSRF and Origin behavior must be reviewed with the HTTP API design.
- Process restart loses in-memory sessions; the client recovers as unauthenticated and can create a new guest identity. This limitation must be documented to the player through ordinary UX only if needed; it is not a server error detail.

## Acceptance criteria

- **W008-AC-01:** A visitor can create a guest session with a valid username without login; server issues a unique guest ID distinct from username.
- **W008-AC-02:** Server applies the exact §6.2 trim, grapheme-count, permitted-character, and unsafe-control rules; boundary and Unicode tests cover valid/invalid examples.
- **W008-AC-03:** Invalid creation returns a stable error, creates no session, and sets no credential cookie.
- **W008-AC-04:** Session credential is cryptographically unpredictable, carried in an HttpOnly cookie, never returned in JSON, never stored in localStorage, URL, or logs, and never placed in WebSocket JSON.
- **W008-AC-05:** Recovery after refresh returns the same guest ID and username for a valid session; invalid/expired/missing credentials are uniformly unauthenticated.
- **W008-AC-06:** Cookie policy, expiry, credential hashing/storage, endpoint schema, and in-memory restart limitation are documented; no durable database is introduced.
- **W008-AC-07:** Session manager APIs are isolated from HTTP/transport and testable; authorization lookup distinguishes unauthenticated from valid session without leaking secrets.
- **W008-AC-08:** React client remembers only display username, calls create/recovery, restores or prompts appropriately, and renders all player-visible content in Thai.
- **W008-AC-09:** Tests cover session creation/recovery/expiry/revocation, collision/failure handling, cookie attributes, username Unicode boundaries, client refresh behavior, and no credential leakage; `pnpm check` and CI pass.
- **W008-AC-10:** Independent review PASS applies to the exact final PR commit; no BLOCKER/MAJOR findings remain, PR is merged to `develop`, feature branch retained, and evidence is recorded.

## Technical design

- Go package boundary: `apps/server/internal/session/` owns guest identity, username validation, credential generation/verifier, expiry, and a concurrency-safe in-memory session store. It does not depend on HTTP or Rooms.
- HTTP adapter: versioned endpoints under `/api/v1/guest-session` create and recover sessions, set/clear cookies, enforce request-size limits, and map domain errors to stable API codes. Exact endpoint/schema remains a proposal pending API design and review.
- React client: add a guest-entry bootstrap state and API client; on app startup recover session with credentials included; keep only display username in localStorage; keep no authentication secret in JS-readable storage.
- Integration seam: Room and WebSocket layers receive an authenticated `GuestSession` value from middleware/connection setup, never a caller-supplied guest ID.
- Use `crypto/rand`; store a keyed or one-way verifier rather than raw cookie credentials if compatible with operational revocation and expiry requirements. Keep comparison constant-time. Choose the concrete verifier and cookie name in technical design.
- Use a well-maintained Go Unicode grapheme segmentation dependency only if the standard library cannot meet §6.2; inspect licensing, maintenance, and compatibility and pin it through Go module metadata.
- Start with process-local storage; no persistence or distributed session dependency.

## Security and privacy

- Require TLS for production; set `Secure`, `HttpOnly`, explicit `Max-Age`/expiry, and a reviewed `SameSite` value.
- Validate Origin on session-mutating endpoints and use same-origin deployment; do not rely on SameSite alone as CSRF protection.
- Bound request body, username bytes, session allocation, and credential generation retries.
- Exclude cookies, credential verifiers, and raw secrets from application logs and error responses; redact sensitive request headers in observability.
- Do not interpret usernames as HTML. Store/display as text and reject unsafe control characters according to the source requirement.
- Avoid claiming sessions survive process restart; no account recovery is implied.

## Test plan

- Go unit tests for username grapheme policy, uniqueness, credential generation, expiry, revocation, constant-time verifier use where testable, capacity/failure behavior, and race-safe concurrent access.
- HTTP tests for endpoint status/error mapping, request limits, cookie flags, no-cookie-on-error, Origin policy, and credential absence from response bodies/log hooks.
- Frontend tests for username input, Thai errors, create/recovery bootstrap, refresh restoration, and localStorage containing display name only.
- Security checks for dependency audit and static analysis; `pnpm check` and exact-head CI.
- Independent review inspects the implementation diff, specs, tests, actual results, and final SHA.

## Dependencies and risks

- Depends on W004 standards/CI and W005 versioned protocol contract (DONE).
- W009 consumes authenticated guest identity and introduces per-Room username uniqueness. W010 consumes session authorization for WebSocket establishment.
- Username grapheme counting is not equivalent to Go rune count or UTF-16 length; dependency quality and test vectors matter.
- In-memory sessions disappear on restart and are single-server only, consistent with MVP architecture but a deliberate limitation.
- Cookie lifetime, SameSite, credential format/length, revocation behavior, and endpoint schema are technical decisions; the cookie lifetime and production policy need Project Owner approval before implementation if they change user-visible recovery expectations.

## Specification review history

Pending independent review.

## Implementation validation evidence

Not run; implementation has not started.

## Completion

Require approved specification, implementation and tests satisfying W008-AC-01–10, actual CI pass, independent exact-SHA review PASS, blocking findings resolved, documentation synchronized, PR merged to `develop`, and feature branch retained.
