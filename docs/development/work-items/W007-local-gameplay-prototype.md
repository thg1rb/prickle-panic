# W007 — Local Gameplay Prototype

Status: PLANNED
Type: Frontend Feature

## Goal

Replace the starter screen with a polished, Thai-language, single-player local Speed Battle prototype that demonstrates the 4×4 Coco card interaction, scoring, incorrect-pick cooldown, timer, and results on desktop and mobile browsers.

## Scope

- Build a small React shell for a Thai title/home state, local game, and local results/replay flow.
- Integrate Phaser 4 as a React-owned board view: Phaser draws the cards and reports pointer selections; React owns navigation, HUD, accessible control semantics, and the temporary local game state.
- Add a disposable TypeScript local game adapter that mirrors W006's Go engine for browser-only interaction. Keep it explicitly local and replace it with server-authoritative state in W011.
- Display a reference card and 16 selectable Coco variants in a visible 4×4 board.
- Implement the baseline local rules: distance 0–1 is correct, +1 score, update reference, replace same slot, preserve a valid selection, increment combo, incorrect pick awards zero, preserves board/reference/score, resets combo, and applies the 0.8-second cooldown without extension.
- Run the configurable baseline 60-second countdown; show local results and allow a replay with a fresh deterministic/random seed.
- Keep player-facing copy in Thai and centralized in one translation module; no server messages, Rooms, accounts, or English localization.
- Support touch, mouse, keyboard, visible focus, and screen-reader access while keeping the Phaser canvas decorative and providing accessible per-card controls.
- Respect `prefers-reduced-motion` and avoid requiring hover or scrolling to play.
- Use temporary procedural vector card visuals. Do not create or claim finished Coco artwork; W014 owns production assets and licensing records.
- Add focused frontend tests for local rules, accessible interactions, timer/cooldown, and end/replay flow using the project's standard TypeScript test tooling.

## Non-goals

- Guest Session, username persistence, Room creation/joining, Lobby, invitations, WebSocket, multiplayer, reconnection, or server authority. W008–W011 own those features.
- Claiming the local adapter as production game authority. W011 must remove or isolate it from the online build path.
- Classic Rush, Party Chaos, matchmaking, or shared boards.
- Production mascot art, licensed audio, or asset pipeline; W014 owns that work.
- Changing the approved BASELINE matching/scoring behavior or configurable TUNING values.

## Requirement references

- `docs/REQUIREMENT.md` §§1.3, 3.1, 4.1–4.10, 5.1–5.5, 8.1–8.10, 9.1–9.7, 18.2–18.4, 22.1–22.4, §32 (W007).
- LOC-001–003; FR-024–035, FR-039–040, FR-043; INV-005, INV-008; NFR-013–018, NFR-020, NFR-022–023, NFR-029–030; AC-006–007.
- Go game domain source: `apps/server/internal/game/` (W006).

## Design direction

- Preserve “Coco's Berry Picnic”: warm cream field, cocoa text and card outlines, sage accents, strawberry coral for the primary action, honey yellow for small game highlights; use Mitr for primary Thai UI and Mali sparingly for decorative headings.
- Keep the home and results states calm and characterful. During play, the 4×4 grid is the visual anchor; the reference card and score/timer remain immediately scannable.
- In the provisional vector cards, distinguish all five binary attributes with shape/placement cues as well as color. Provide descriptive Thai accessible labels for each card.
- At 320 CSS px, target a 4×4 grid without page scrolling and touch targets around 44×44 CSS px or larger where practical. Validate actual text and card readability rather than claiming compliance from CSS dimensions alone.
- Motion should explain a result (correct/incorrect) and remain subtle; reduced-motion mode disables nonessential movement.

## Proposed technical boundary

- `apps/web/src/game/localSpeedBattle.ts`: pure local state transitions, deterministic injectable random source, timer/cooldown state, and testable result types. Copy BASELINE behavior from W006; document this as a temporary adapter to be removed/replaced at W011.
- `apps/web/src/game/cocoCardVariant.ts`: explicit mapping to the Go/protocol MSB-first attribute order. Keep TypeScript values in 0–31 and test all 1,024 ordered pairs against the same rule.
- `apps/web/src/game/PhaserBoard.tsx` or equivalent: create/destroy the Phaser game with React lifecycle; pass immutable board snapshots into the scene; forward selected card IDs through a callback/ref without recreating the Phaser instance on every HUD update.
- React owns screen state, timer display, status messages, and accessible controls; Phaser owns canvas drawing and pointer hit feedback. A synchronized React button grid supplies keyboard/screen-reader operation and invokes the same action callback.
- Centralize Thai copy in `apps/web/src/i18n/th.ts` (or an equivalent translation map); do not put display strings in the game adapter.
- Use only project dependencies already present where practical; if adding test tooling, pin compatible versions in the pnpm lockfile and validate them in CI.

## Acceptance criteria

- **W007-AC-01:** Home, game, error/feedback, timeout, and result UI text is Thai and comes from centralized translation keys; no English player-visible copy is introduced.
- **W007-AC-02:** The local game displays a reference card and exactly 16 selectable cards in a visible 4×4 board; card variation mapping follows W005/W006 bit order.
- **W007-AC-03:** TypeScript matching is exhaustively checked across all 1,024 ordered variation pairs and agrees with the documented `distance <= 1` rule.
- **W007-AC-04:** Correct picks score +1, increment combo, update reference, replace the selected position, and preserve board solvability. Incorrect picks score zero, preserve score/reference/board, reset combo, show Thai feedback, and apply a non-extending 0.8-second cooldown.
- **W007-AC-05:** Configurable 60-second baseline timer ends the local match, shows result score, and replay creates fresh state with score/combo/timer reset.
- **W007-AC-06:** React owns app state and Phaser lifecycle is cleaned up on unmount; card pointer selection and accessible button selection invoke the same transition once.
- **W007-AC-07:** All 16 cards are usable by touch and keyboard; visible focus, Thai accessible labels, reduced-motion preference, readable contrast, and no hover-only action are supported.
- **W007-AC-08:** At 320 CSS px the full board and primary gameplay controls fit without document-level horizontal or vertical scrolling; browser validation records actual viewport and observed behavior.
- **W007-AC-09:** The UI clearly represents a local prototype. The local adapter has no network authority and no Room/session code; W011 replacement is documented.
- **W007-AC-10:** Focused unit/component tests and CI pass; no test result is claimed unless it ran.

## Dependencies

- W001–W004 (DONE).
- W006 — Card Matching Engine (DONE).
- W005 — Protocol Contracts (DONE; W011 will bind runtime events).

## Validation

- Type check, formatting, lint, frontend production build, and focused component/state tests.
- Browser QA at a desktop viewport and 320 CSS px mobile viewport, including pointer and keyboard use, reduced motion, correct/incorrect pick, cooldown edge, timeout, and replay.
- Inspect Phaser cleanup under screen transitions and verify one action per card activation.
- Run `pnpm check` and CI security checks; report browser and test outcomes separately.

## Independent review focus

- Thai-only visible text and translation centralization.
- Matching/score parity with W006 and no client-local rule leaking into future authority.
- Phaser lifecycle and React/Phaser state synchronization.
- Touch dimensions, 4×4 fit, focus/keyboard/screen reader labels, reduced motion, and accessible feedback.
- Timer/cooldown edge cases, replay resets, and test evidence.

## Completion

Require the approved specification, focused tests and browser QA evidence, CI pass, independent review PASS for the exact final PR revision, no unresolved BLOCKER/MAJOR findings, documentation synchronized, PR merge to `develop`, and retained feature branches.
