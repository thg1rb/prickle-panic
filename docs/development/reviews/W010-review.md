# Independent Specification Review — W010 WebSocket Foundation

## Review 1

### Reviewed revision

- Pull Request: #12
- Commit SHA: `d5c54bfb786d0a81e903336e3aba662515442052`
- Reviewer execution reference: Independent Sub-Agent
- Review scope: Specification and task plan only; no implementation or tests were inspected or run.

### Finding

| ID | Severity | Evidence and impact | Required remediation |
|---|---|---|---|
| W010-SR-01 | MINOR | W005 defines Room-scoped sequence values while private board events are owner-only. If hidden events consume the Room sequence, another player's normal private board updates appear as sequence gaps and cause unnecessary snapshot requests. | Clarify whether hidden events consume each recipient's sequence or use authorized streams; add sequence/privacy test coverage. |

### Decision

PASS

Reason: No BLOCKER or MAJOR findings. W010 is a Room-only foundation and this sequencing clarification is a minor follow-up before W011 private gameplay events. No implementation or tests were reviewed.

## Review 2

### Reviewed revision

- Pull Request: #12
- Commit SHA: `ec39f546875428c56cfe32419f7b9be006f863c8`
- Reviewer execution reference: Independent Sub-Agent
- Review scope: Specification and task plan only; no implementation or tests were inspected or run.

### Re-review results

- W010-SR-01 — RESOLVED. Transport sequence values are contiguous per authorized recipient stream; private events delivered to another player do not create artificial gaps or redacted placeholders.
- The sequence proposal is explicitly pending Project Owner approval and W005 contract review. T01 requires a recorded decision and coordinated contract/fixture updates before implementation.
- No new findings.

### Decision

PASS

Reason: The proposal resolves the sequencing/privacy gap while preserving the requirement's ordered-event invariant. Its change to W005's v1 contract is explicitly gated for approval and review. This is a specification-only review; no implementation or tests were inspected or run.

## Review 3 — Session-expiry reconciliation

### Reviewed revision

- Pull Request: #12
- Specification revision: `7bd1e72a88b5714ea0898a1117f4e2eefc28464b`
- Reviewer execution reference: Independent read-only Sub-Agent (`/root/w008_spec_review`)
- Review scope: W008 session-expiry ownership, established-WebSocket revocation/closure, typed expiry versus network-disconnect signaling, acceptance criteria, and test coverage.

### Findings

- No BLOCKER or MAJOR findings. The specification correctly enforces immediate authorization revocation on expiry, sends a typed expiry signal for W013, and reserves the 15-second grace for unexpected network loss.
- MINOR W010-SR-02: one dependency sentence said W008 approval was pending, conflicting with the approved baseline described elsewhere.

### Decision

PASS with the MINOR wording correction required. No implementation or tests were inspected or run.

## Re-review 4 — W010-SR-02

- Re-reviewed revision: `69e6e247f09f3b0e1afe804180663280a30608f2`
- Reviewer execution reference: Independent read-only Sub-Agent (`/root/w008_spec_review`)
- Finding: W010-SR-02 RESOLVED. The sentence now states W008 policy/specification approval is complete and PR #10 integration remains required before implementation.
- New findings: None.
- Decision: PASS. The independent review covers the substantive specification at `7bd1e72a88b5714ea0898a1117f4e2eefc28464b` and confirms its final wording correction at `69e6e247f09f3b0e1afe804180663280a30608f2`. No implementation or tests were inspected or run. This review does not authorize merge.
