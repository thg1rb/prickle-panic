# Independent Specification Review — W009 Room Management

## Review 1

### Reviewed revision

- Pull Request: #11
- Commit SHA: `040608a5279d570b32b8332de635b9b537c728fc`
- Reviewer execution reference: Independent Sub-Agent
- Review scope: Specification and task plan only; no implementation or tests were inspected or run.

### Evidence inspected

- `docs/REQUIREMENT.md` relevant Room, session, leave/reconnect, security, architecture, acceptance, and W009 sections
- `docs/development/work-items/W009-room-management.md`
- `docs/development/work-items/W009-task-plan.md`
- `docs/protocol/v1/README.md` and `message.schema.json`
- `docs/development/work-items/W008-guest-session.md` pending dependency boundary
- Revision diff

### Findings

| ID | Severity | Evidence and impact | Required remediation |
|---|---|---|---|
| W009-SR-01 | MAJOR | §7.9 confirms Host may cancel Countdown. The design named a cancel hook, but expected behavior, acceptance, protocol and tests did not specify the operation. Implementation could omit a confirmed Host permission. | Define actor/state/result behavior, protocol event/error, acceptance criteria, and tests. |
| W009-SR-02 | MAJOR | The design did not define Room leave integration, including §13.6's immediate active-Match forfeit with no grace or re-entry. The domain interface also omitted leave. | Define Lobby/Countdown/active-Match leave behavior and W013 handoff, with acceptance criteria and tests. |
| W009-SR-03 | MAJOR | The proposed “single membership” rule imposed one-Room-at-a-time membership, while FR-008 only specifies one active connection per RoomPlayer. This added a user-visible restriction without approval. | Remove the unapproved restriction or label it as a Project Owner decision; keep connection arbitration with W010. |
| W009-SR-04 | MINOR | FR-005 Lobby username change is assigned to W009 by the W008 plan and is P1, but was not included or explicitly deferred in W009. | Add FR-005 scope/acceptance or record an explicit deferred owner. |

### Decision

FAIL

Reason: Three MAJOR coverage/policy gaps and one MINOR traceability gap require remediation. No implementation or tests were reviewed.

## Review 2

### Reviewed revision

- Pull Request: #11
- Commit SHA: `513708d28e46ecdf468c8302fbc89065d701642b`
- Reviewer execution reference: Independent Sub-Agent
- Review scope: Specification and task plan only; no implementation or tests were inspected or run.

### Re-review results

- W009-SR-01 — RESOLVED. Host countdown cancellation now specifies actor, transition, protocol scope, and test/acceptance coverage.
- W009-SR-02 — RESOLVED. Lobby/Countdown/active-Match leave behavior and the W013 forfeit seam are specified.
- W009-SR-03 — RESOLVED. The unsupported one-Room membership restriction was removed.
- W009-SR-04 — RESOLVED. FR-005 Lobby username change is assigned to W009 with acceptance coverage.

### Finding

| ID | Severity | Evidence and impact | Required remediation |
|---|---|---|---|
| W009-SR-05 | MAJOR | The spec permits memberships in multiple Rooms but only required W010 to enforce one active connection per RoomPlayer. `REQUIREMENT.md` INV-004 prohibits one Guest Session from controlling multiple active Room connections, so separate memberships could still enable concurrent active connections. | Distinguish stored membership from active authority; require global W010 enforcement of INV-004, define replacement/disconnect behavior, and add acceptance/test coverage. |

### Decision

FAIL

Reason: W009-SR-01 through W009-SR-04 are resolved, but W009-SR-05 remains a material conflict until the global Guest Session connection invariant and replacement behavior are explicit. No implementation or tests were inspected or run.

## Review 3

Pending independent re-review of the clarification to W009-SR-05.
