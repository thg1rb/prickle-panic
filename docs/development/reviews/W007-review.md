# W007 Independent Review

## Review 1

- Decision: FAIL
- Reviewed commit: `fbad25b330a04a1adbfada5a5a73777546832074`
- Pull Request: #9
- CI run: `39209046674` — PASS
- Reviewer: Independent Sub-Agent

### Findings

#### W007-IR-01 — MAJOR — Player-facing copy is not fully centralized

`describeCard` kept Thai attribute descriptions in the game rules module, and the home footer appended Thai copy directly in the React component. This conflicts with W007-AC-01's centralized Thai copy requirement. Move all displayed strings to `apps/web/src/i18n/th.ts`; keep card bit/attribute logic language-neutral.

#### W007-IR-02 — MINOR — Pointer/touch and Phaser unmount QA is incomplete

The React component tests mock Phaser, and synthetic browser pointer events did not exercise the Phaser input path. Verify pointer/touch selection and canvas removal after unmount in a browser when an input method that can target card coordinates is available. Do not report these behaviors as browser-verified until then.

### Review observations

- W006 rule parity was consistent: valid picks score, advance the reference, replace a slot, increment combo, and preserve solvability; invalid picks preserve board/reference/score, reset combo, and apply a non-extending cooldown.
- Local frontend tests and exact-SHA CI passed.
- At 320×568 CSS pixels, the browser document fit the viewport and card controls were about 65×65 CSS pixels.

### Remediation

Pending.
