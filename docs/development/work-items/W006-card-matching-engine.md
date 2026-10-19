# W006 — Card Matching Engine

Status: PLANNED
Type: Core Engine

## Goal

Implement a pure, deterministic Go domain package for Coco's five-bit card variations, authoritative matching, Speed Battle scoring/combo transitions, and generation of a 16-card board that always contains a valid selection. The package must have no dependency on HTTP, WebSocket, storage, or frontend code.

## Scope

- Represent 32 card variations as values 0–31 with the W005 MSB-first attribute mapping.
- Represent card instances separately from variations and keep instance IDs unique within a board/game-state scope.
- Implement `HammingDistance` and the confirmed `distance <= 1` match rule.
- Implement correct-pick score/combo transition semantics and incorrect-pick outcome semantics without embedding transport or cooldown timing.
- Implement seeded deterministic board/reference generation suitable for repeatable tests; maintain at least one valid selection after initial generation and every replacement.
- Provide immutable or copy-safe operations so callers cannot mutate engine-owned input unexpectedly.
- Add table-driven edge tests and exhaustive matching over all 1,024 ordered variation pairs.
- Document the generation policy as a technical proposal and make no tuning parameter change.

## Non-goals

- Match timers, network/abuse rate limits, mistake-cooldown clocks, Rooms, sessions, WebSocket handlers, and persistence.
- UI, card artwork, Phaser integration, or generated assets.
- Full Match lifecycle or winner calculation; W011/W012 own those integrations.
- Classic Rush or Party Chaos behavior.
- Changing authoritative matching, +1 score, incorrect-pick effects, or the configurable 4×4/16-card/60-second/0.8-second tuning values.

## Requirement references

- `docs/REQUIREMENT.md` §§1.3, 8.1–8.10, 9.2–9.7, §32 (W006).
- FR-024–040, NFR-020–023, INV-005, INV-008, AC-006, AC-007, AC-024.
- Protocol mapping: `docs/protocol/v1/README.md` and `docs/protocol/v1/message.schema.json`.

## Proposed generation design

- A caller supplies a seed or `math/rand` source; no package-global randomness is permitted. The server retains random state privately and never exposes future sequences.
- Generate a reference variation, then independently sample 16 variation values for board instances. Duplicate variations are allowed because requirements distinguish variation from instance.
- If an initial board contains no valid match, replace its last slot with a uniformly selected valid variation from the six values at Hamming distance 0 or 1. After a correct pick, first place the replacement in the selected slot; if the resulting board has no valid choice relative to the new reference, replace that same slot with a uniformly selected valid variation. This is a proposed deterministic fallback; it preserves the confirmed solvability invariant but slightly changes distribution only for a board that would otherwise be unsolvable.
- Every generated instance receives a unique opaque ID scoped to the generated board. Replacement generation must not reuse any active card-instance ID.
- Treat randomness distribution, seed plumbing, and replacement policy as implementation details; do not expose PRNG state to clients or turn the current tuning table into constants that cannot be configured later.

## Acceptance criteria

- **W006-AC-01:** All 32 variations map to the five attributes in the protocol's documented MSB-first order; conversions reject values outside 0–31.
- **W006-AC-02:** Matching is exactly `HammingDistance(reference, selected) <= 1`, including identical variations; all 1,024 ordered pairs are exhaustively verified.
- **W006-AC-03:** A variation and a card instance are distinct types/concepts; generated card IDs are unique in scope, including replacements, while duplicate variations remain allowed.
- **W006-AC-04:** Correct pick transition awards the configured +1 point, advances reference to the selected variation, removes/replaces the selected slot, increments combo, and preserves a 16-card board with a valid choice. Incorrect pick awards zero, preserves board/reference/score, and resets combo. Cooldown timing and technical rejection are caller concerns.
- **W006-AC-05:** Seeded generation is reproducible, does not use package-global random state, and guarantees an available valid choice on every produced interactive board.
- **W006-AC-06:** Public engine APIs are deterministic and do not mutate caller-owned input; invalid indices/identifiers return explicit errors without partial mutation.
- **W006-AC-07:** Unit tests cover bit/attribute mapping, both body-color values, all single-bit toggles, Hamming distance boundaries, duplicate variations with unique instances, scoring/combo transitions, board cardinality, solvability, replacement identity, invalid actions, and multiple seeds.
- **W006-AC-08:** The independent package has no transport or frontend imports and passes `gofmt`, `go test`, and `go vet` in CI.
- **W006-AC-09:** The generation fallback and any chosen random-source/ID strategy are documented as proposed technical choices; confirmed gameplay rules and TUNING values remain unchanged.

## Implementation boundary

Proposed location: `apps/server/internal/game/` with small card, board, and scoring domain files. Keep APIs limited to operations needed by W007 local play and W011 authoritative server gameplay. The implementer should inspect current package conventions before adding subpackages.

## Dependencies

- W004 — Development Standards & CI (DONE).
- W005 — Protocol Contracts (DONE).

## Validation

- Run `gofmt` and `go test ./...` in `apps/server`.
- Run Go coverage reporting for the pure game package; report the actual result without setting an unapproved coverage threshold.
- Run `go vet ./...` and repository `pnpm check`.
- Confirm exhaustive pair verification covers all 1,024 ordered pairs, not merely unique unordered pairs.

## Independent review focus

- Exact matching and bit order, inclusive zero-distance rule, exhaustive pair count.
- No accidental change to confirmed score, combo, cooldown, or board rules.
- Deterministic solvability guarantee and documented sampling bias.
- Card-instance identity uniqueness and replacement behavior.
- Pure boundaries, error atomicity, and test evidence.

## Completion

Require an approved specification, implementation and required validation, independent review PASS tied to the final PR SHA, no unresolved BLOCKER/MAJOR findings, documentation synchronized, merge to `develop` through a PR, retained feature branch, and actual validation/review/merge evidence recorded.
