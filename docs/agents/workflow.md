# Agent-first Spec-Driven Development Workflow

## Required sequence

1. Read `docs/REQUIREMENT.md` and identify relevant IDs and their statuses.
2. Analyze the requested scope, edge cases, dependencies, risks, and conflicts.
3. Write a Feature Specification and Technical Design using the templates in this directory.
4. Break approved work into a bounded Task Plan with acceptance criteria and validation steps.
5. Obtain the approval required by the project policy before implementation. Confirmed product decisions and material architecture changes require Project Owner approval.
6. Create a feature branch from the current `develop`; implement only the assigned scope.
7. Run applicable validation and record actual outcomes. Never report unexecuted checks as passing.
8. Open a PR to `develop`; provide the diff, Work Item, requirements, acceptance criteria, test evidence, and exact commit SHA.
9. Launch a separate read-only Reviewer Agent. The implementer cannot approve their work. The reviewer independently inspects source, diff, specifications, tests, and evidence.
10. Remediate FAIL findings and request re-review. Material changes after PASS require review of the new revision.
11. Merge only after independent PASS, required CI, documentation synchronization, and applicable authorization. Use a Merge Commit and retain the feature branch.
12. Record the final revision, PR, validation, review, limitations, and follow-ups before marking the Work Item DONE.

## Work Item states

Use the following states:

`PLANNED` → `SPEC_READY` → `IN_PROGRESS` → `READY_FOR_REVIEW` → `CHANGES_REQUESTED` / `REVIEW_BLOCKED` → `VERIFIED` → `DONE`.

`CANCELLED` is reserved for explicitly cancelled work. `BLOCKED` must identify the missing evidence, access, tool, or decision. A blocked item is not verified or done. A review PASS alone does not authorize merge, release, or deployment.

## Traceability

Every specification, task plan, implementation report, and PR references requirement IDs. Tests should reference acceptance criteria or requirement IDs where practical. Review evidence names the exact commit SHA and tests inspected or run. Keep a trace chain from requirement → specification → task → implementation → test → independent review.

Use test outcomes `PASS`, `FAIL`, `NOT_RUN`, `BLOCKED`, or `NOT_APPLICABLE`; provide a reason for `NOT_APPLICABLE`. Test reports must distinguish developer-run checks from CI and reviewer-run checks.

## Scope and change control

Do not add unrelated work or alter confirmed requirements. Record unrelated discoveries as follow-up Work Items. A material architecture change requires justification, impact analysis, an ADR update, and Project Owner approval when applicable. Tuning values start at their documented values; changes require playtesting evidence or approved design adjustment.

## Independent review

Decisions are `PASS`, `FAIL`, or `BLOCKED`. PASS requires all acceptance criteria met, sufficient revision-specific evidence, and no unresolved BLOCKER or MAJOR findings. Record findings as BLOCKER, MAJOR, MINOR, or SUGGESTION. The review must identify the reviewer execution and exact SHA. Do not treat an implementer-authored PASS file as independent evidence.

If independent review is unavailable, mark the Work Item `REVIEW_BLOCKED`; do not merge or mark DONE. A previous PASS does not apply to material changes on a later commit.

## Definition of Done

A Work Item is DONE only when the approved specification is satisfied, applicable checks actually pass, acceptance criteria are verified, independent review passes, blocking findings are resolved, documentation is synchronized, required PR merge is complete, and evidence is recorded. Preserve the feature branch after merge.
