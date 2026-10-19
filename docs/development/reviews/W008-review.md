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

Pending independent re-review of the remediation revision.
