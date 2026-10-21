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

### Reviewed revision

- Pull Request: #10
- Commit SHA: `97a6ff4f0c6b95420b76b45fb50ba9895dce2949`
- Reviewer execution reference: Independent Sub-Agent

### Evidence inspected

- Requirements: `docs/REQUIREMENT.md` §§6.1–6.7, 16.1, 20, and 32
- Specification: `docs/development/work-items/W008-guest-session.md`
- Task plan: `docs/development/work-items/W008-task-plan.md`
- Protocol boundary: `docs/protocol/v1/README.md`
- Revision diff from the prior review

### Findings and re-review

- W008-SR-03 — RESOLVED. The proposed policy names the Unicode Bidi_Control characters to reject, adds test vectors, and preserves Arabic text and ZWJ emoji.
- No new findings.

### Decision

PASS

Reason: The remaining stricter username filter is explicitly a proposal requiring Project Owner approval before implementation. No requirement conflict was found. This review covered specifications only; no implementation or tests were inspected or run.

## Review 5

### Reviewed revision

- Pull Request: #10
- Commit SHA: `cd37fb4096b1da295d7736985979b382f008bcf6`
- Reviewer execution reference: Independent Sub-Agent (`/root/w008_spec_review`)
- Review type: Specification and requirement source only

### Evidence inspected

- `docs/REQUIREMENT.md` §§6.2, 6.5, 13.1, and 14.4
- `docs/development/work-items/W008-guest-session.md`
- `docs/development/work-items/W008-task-plan.md`
- Exact revision and related changes in the PR

### Findings

No BLOCKER, MAJOR, MINOR, or SUGGESTION findings.

### Decision

PASS. The owner-approved username rules, 24-hour idle and 7-day absolute session lifetime, qualifying activity, credential/cookie policy, session creation limit, and immediate active-Match forfeit on session expiry agree between the requirement source and W008 specification. No implementation or tests were reviewed or run. This PASS applies to specification revision `cd37fb4096b1da295d7736985979b382f008bcf6` and does not authorize merge.

## Review 6 — PR scope reconciliation

- Pull Request: #10
- Reviewed PR head: `32bd3f45ae129ac0dfca9e0f9e6d41625629a90a`
- Reviewer execution reference: Independent read-only Sub-Agent (`/root/w008_spec_review`)
- Scope: Confirmed W008 content remains identical to the previously approved and reviewed policy revision; checked the complete PR diff for unrelated W007 edits.
- Decision: PASS. W008 policy/specification content is unchanged; the unrelated W007 file no longer appears in the PR diff. No new findings. No implementation or tests were reviewed or run. This review does not authorize merge.

## Review 7 — Specification readiness status

- Pull Request: #10
- Reviewed PR head: `daae4de3c4d84daae6264bf5d8e3c184f069df7d`
- Reviewer execution reference: Independent read-only Sub-Agent (`/root/w008_spec_review`)
- Scope: Consistency of W008 specification/task-plan state with exact-revision independent review evidence and pending integration status.
- Decision: PASS. Both documents identify the specification as `SPEC_READY`; the policy/specification PASS and scope review PASS support that state, and PR #10 integration remains explicitly pending. No new findings. This status-only review does not authorize merge or implementation.

## Review 8 — T01 technical design

- Pull Request: W008 T01 design PR (opened after this report is recorded)
- Final reviewed commit SHA: `3d9ed43a3215f339e19f10464aaedfa241f866f8`
- Reviewer execution reference: Independent read-only Sub-Agent (`/root/w008_spec_review`)
- Review type: Exact-revision W008 technical design and task-plan review, including two remediation re-reviews

### Evidence inspected

- `docs/REQUIREMENT.md` §§6.2, 6.5, 13.1–13.5, and 14.4
- `docs/agents/workflow.md`
- W008 specification and task-plan changes from the T01 branch
- W010, W011, and W013 session-expiry and Match-deadline contracts

### Review rounds and remediation

1. Review 8 at `d062b87fb88ad7e695af2013b3989d725c87dfd3`: FAIL with two MAJOR findings and one MINOR finding. The equal-time Match-deadline policy needed traceable Owner approval; expiry-event delivery needed a proactive trigger when no request or socket exists; and the grapheme module/package paths needed to be distinguished.
2. Review 9 at `b286294803d8287b0d5b8fceb016ec26af940bd2`: FAIL with one MAJOR finding. UTC session-expiry timestamps and monotonic Match deadlines lacked a common ordering basis. The reviewer confirmed the three prior findings were resolved.
3. Review 10 at `3d9ed43a3215f339e19f10464aaedfa241f866f8`: PASS with no findings. The design now uses one injected monotonic clock domain for expiry and Match deadlines. Expiry events carry a monotonic `DeadlineTick` for ordering and UTC time for diagnostics. Delayed delivery, wall-clock adjustments, no-request/no-socket expiry, and exact-deadline ties are included in the test plan.

Other remediation details: W008 records the Owner-approved Match-deadline precedence and requires W013 specification synchronization before its implementation; specifies a proactive indexed scheduler with one deadline entry per active session; and distinguishes Go module and package paths.

### Final decision

PASS at commit `3d9ed43a3215f339e19f10464aaedfa241f866f8`. The final review confirmed the approved W008 policy, API, cookie, proxy, quota, Unicode, and W008/W010/W013 integration boundaries. This is a design review only; no guest-session implementation or implementation tests were reviewed or run. T02 has not started. PASS does not authorize T01 PR merge.

## Review 9 — T02 Guest Session implementation

- Work Item: W008 T02
- Branch: `feature/W008-guest-session-implementation`
- Reviewed implementation SHA: `b24338a1d3d987e6fe23ff3dc2c39c9bc29ef615`
- Reviewer execution reference: Independent read-only Sub-Agent (`/root/w008_spec_review`)
- Review type: Exact-revision code, test, requirement, and integration-boundary review; one implementation re-review and one test-only re-review

### Evidence inspected

- `docs/REQUIREMENT.md` §§6.2, 6.5, 13.1, 14.4, and W008 security/functional requirement IDs
- W008 specification, task plan, T01 design, and implementation report
- Full T02 implementation diff and tests in `apps/server/internal/session/`, `internal/clientip/`, and `internal/clock/`
- W009, W010, and W013 boundaries and the approved session-expiry policy
- Actual developer test results and the current Git revision

### Findings and remediation

1. PASS review at `ddc8a46a956a56a124d8f3a523adebe1df5fac8a` found no substantive issues and two MINOR test gaps: the entropy failure test failed during guest-ID generation rather than credential generation, and the absolute-lifetime test could fail at the earlier idle deadline. Both were fixed in `b24338a`: the test reader now returns guest-ID bytes before failing the credential read, and the absolute-lifetime test records successful activity through tick 24 before verifying rejection at the fixed tick-25 deadline.
2. The reviewer re-reviewed `b24338a1d3d987e6fe23ff3dc2c39c9bc29ef615` and confirmed both MINOR findings resolved, with no new findings.

### Validation independently run by reviewer

- `pnpm check` — PASS
- `go test ./...` — PASS
- `go test -race ./...` — PASS
- `go vet ./...` — PASS
- `git diff --check` — PASS

The reviewer verified the current implementation report's dependency-audit and `govulncheck` evidence but did not rerun those commands. GitHub CI had not run because the implementation PR was not yet open.

### Final decision

PASS at `b24338a1d3d987e6fe23ff3dc2c39c9bc29ef615`. No unresolved BLOCKER, MAJOR, MINOR, or SUGGESTION findings remain. The implementation matches the approved session lifetime, credential storage, activity qualification, creation quota, username policy, trusted proxy boundary, and expiry-event contract. HTTP/cookies, WebSocket handling, Room logic, Match forfeit, and frontend behavior remain out of T02 scope. The PASS validates implementation quality; it does not authorize PR merge.
