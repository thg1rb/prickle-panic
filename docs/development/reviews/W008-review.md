# W008 Independent Specification Review

## Review 1

### Reviewed revision

- Pull Request: #10
- Commit SHA: `8391d663f7176e1fddd64829f00ca18ba69e062d`
- Reviewer execution reference: Independent Sub-Agent

### Evidence inspected

- Requirements: `docs/REQUIREMENT.md` §§6.1–6.7, 16.1, 20, 28, and 32
- Specification: `docs/development/work-items/W008-guest-session.md`
- Task plan: `docs/development/work-items/W008-task-plan.md`
- Protocol boundary: `docs/protocol/v1/README.md`
- Bundled W007 completion status update and exact merged review/CI evidence

### Findings

| ID | Severity | Evidence and impact | Required remediation |
|---|---|---|---|
| W008-SR-01 | MAJOR | §6.2 explicitly permits Thai, English, and numbers but does not define them as an exclusive allowlist. W008 did not state whether other safe Unicode, including emoji, is valid, risking an unapproved restriction. | State that listed categories are guaranteed accepted but not an exclusive allowlist, or obtain Project Owner approval before introducing a narrower policy. |
| W008-SR-02 | MINOR | §20 SEC-006 requires rate limits for sensitive operations. W008 introduces unauthenticated session creation but assigned rate-limit enforcement solely to W013, leaving an exposed endpoint without an interim guard. | Add a configurable session-creation limit and tests, or gate public endpoint exposure on W013. Record ownership and validation evidence. |

### Bundled W007 update

The reviewer confirmed changing W007 to DONE is supported by exact-head review PASS, CI pass, merged PR #9, and retained branch. No scope issue was found.

### Decision

FAIL

Reason: W008-SR-01 is a material behavior ambiguity and W008-SR-02 leaves an unauthenticated sensitive endpoint without an explicit release protection. No implementation or tests were run; this was a specification-only review.

## Review 2

### Reviewed revision

- Pull Request: #10
- Commit SHA: `675f770f25583edd4bf9d149e53986b689172573`
- Reviewer execution reference: Independent Sub-Agent

### Findings and re-review

- W008-SR-01 — RESOLVED. The named character categories are accepted but not an exclusive allowlist; other safe Unicode is permitted.
- W008-SR-02 — RESOLVED at specification level. A configurable pre-authentication limiter, non-session-only protection, no-allocation rejection tests, and a public-exposure gate are required. Threshold/source remain T01 design proposals.
- No new BLOCKER, MAJOR, or MINOR findings. The reviewer recommends defining concrete unsafe Unicode/rendering test vectors in T01 while preserving safe scripts.
- Bundled W007 completion update remains supported by its merged PR, exact-head PASS, CI, and retained branch.

### Decision

PASS

Reason: Both prior findings are resolved in the specification, and no blocking or major findings remain. This PASS applies to the exact specification revision above; no implementation or tests were run.

## Review 3

### Reviewed revision

- Pull Request: #10
- Commit SHA: `41be8e8f2d984151b2fbb185cd18271a5e4fb21d`
- Reviewer execution reference: Independent Sub-Agent

### Evidence inspected

- Requirements: `docs/REQUIREMENT.md` §§6.1–6.7, 16.1, 20, and 32
- Specification: `docs/development/work-items/W008-guest-session.md`
- Task plan: `docs/development/work-items/W008-task-plan.md`
- Protocol boundary: `docs/protocol/v1/README.md`

### Finding

| ID | Severity | Evidence and impact | Required remediation |
|---|---|---|---|
| W008-SR-03 | MINOR | The design covers bidi isolation in the UI but does not explicitly assess embedded Unicode bidirectional override/isolate controls. Those controls can affect the displayed order of username text despite surrounding isolation. | Explicitly assess and test bidi controls while preserving ZWJ emoji support. Keep any policy change subject to Project Owner approval. |

### Decision

PASS

Reason: No BLOCKER or MAJOR findings remain. The MINOR follow-up is addressed in the next specification revision with a proposed rejection policy and test vectors for Bidi_Control characters, while preserving Arabic text and ZWJ emoji; that revision requires a fresh review. This review covered specifications only; no implementation or tests were inspected or run.

## Review 4

Pending independent review of the exact specification revision that proposes rejecting Unicode Bidi_Control characters and adds corresponding test vectors.
