# W005 — Protocol Contracts

Status: DONE
Type: Core Engine

## Goal

Define a versioned, machine-readable WebSocket contract for guest room participation and Speed Battle actions/events, preserving server authority and private player state.

## Scope

- Define JSON envelopes, client request and server event categories, request correlation, state versioning, ordering, error shape, and compatibility policy.
- Define payload schemas for the MVP Room/Lobby/Speed Battle interactions required by the approved requirements.
- Preserve the specified v1 example fields and error codes; distinguish protocol failures from gameplay-invalid selections.
- Specify snapshot and incremental-event synchronization, authorized projection rules, and stale/duplicate request handling.
- Add representative valid/invalid contract fixtures and an automated schema validation command wired into CI.
- Document a path for future game modes without defining their entire event systems.

## Non-goals

- Implementing a WebSocket server/client, Guest Sessions, Room Manager, or Match Engine.
- Changing gameplay scoring, matching, cooldown, timeout, forfeit, or victory behavior.
- Implementing Classic Rush or Party Chaos event catalogs.

## Requirement references

- `docs/REQUIREMENT.md` §§18.4–18.10, 19.1–19.12, §§20.1–20.8, §27.
- FR-064, NFR-011, NFR-027, SEC-001, SEC-003–005, AC-025.
- Roadmap W005 — Protocol Contracts.

## Acceptance criteria

- **W005-AC-01:** The protocol defines v1 JSON client request, server response/event, snapshot, and structured error envelopes with stable request correlation.
- **W005-AC-02:** Every message schema sets explicit required fields, types, bounds, unknown-field policy, and version behavior; malformed payloads fail closed.
- **W005-AC-03:** Contract coverage includes guest initialization boundary, room create/join/leave, ready state, start, card pick, synchronization, score/result, and predictable errors without implementing handlers.
- **W005-AC-04:** Server event sequence/state versions and client action versions prevent stale or cross-Match/Wave actions while retaining mode-aware validation semantics.
- **W005-AC-05:** Snapshot and event schemas prevent another player's private board from appearing in public room messages; authorized self-state is explicit.
- **W005-AC-06:** Error codes include all required documented values, define safe retry/resync behavior where possible, and do not expose session secrets or internals.
- **W005-AC-07:** Valid/invalid fixtures are checked by a repeatable command included in CI; tests map to protocol acceptance criteria.
- **W005-AC-08:** Contract documentation marks proposed technical details as such and leaves non-MVP modes deferred.

## Review focus

Check all message fields against §19, identity and authorization assumptions, private-state disclosure, ordering/idempotency semantics, bounded payloads, error-code preservation, schema tests, and accidental policy changes to gameplay behavior.

## Completion

Completed: schema fixtures/validation and CI passed, independent review PASS was recorded for the implementation and final review-evidence revision, PR #5 merged to `develop` as `1a0ac6ba0c10fa1b027ba989985dc3ef74d55e95`, and the feature branch is retained. Final PR-head review evidence is documented on PR #5; protocol implementation was reviewed at `0c9f050d6a5e5514b23598b6a646bc1f8b17084c`, with review-evidence-only final commit `f13423071461150486e78d6a12b5a9443e2a6fa5` also independently verified.
