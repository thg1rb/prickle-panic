# W002 — SDD Documentation Foundation

Status: READY_FOR_REVIEW
Type: Foundation

## Goal

Provide a compact and executable documentation system for Agent-first SDD without duplicating the authoritative requirements.

## Scope

- Add a root `AGENTS.md` entry point and documentation index/ownership policy.
- Document Work Item states, specification-first workflow, traceability, review, and Definition of Done.
- Add reusable specification, task, implementation report, review, and ADR templates.
- Record the baseline architecture boundaries, roadmap, testing/evidence strategy, and Skills registry state.
- Mark W001 DONE now that its review passed and PR #1 merged.

## Non-goals

- Implementation of application features, CI, Skill installation/approval, or changes to `docs/REQUIREMENT.md`.
- Rewriting the requirements document into duplicate product/feature documents.

## Requirement references

- `docs/REQUIREMENT.md` §§1, 18–21, 23–27, 29–33, 40–41.
- NFR-030, AC-025, AC-026.
- Roadmap W002 — SDD Documentation Foundation.

## Acceptance criteria

- **W002-AC-01:** Root `AGENTS.md` directs agents to authoritative docs, language policy, project-local Skills policy, and scope/change control without copying the requirements.
- **W002-AC-02:** Documentation index identifies source of truth, ownership, structure, and maintenance/review policy.
- **W002-AC-03:** Workflow defines spec/design/task/implementation/verification/review/remediation/merge/completion sequence, states, evidence labels, traceability, and DoD.
- **W002-AC-04:** Templates cover feature specs, task plans, implementation reports, review reports, and ADRs.
- **W002-AC-05:** Architecture/testing/roadmap documents reflect the baseline and distinguish unresolved technical design from approved requirements.
- **W002-AC-06:** Skills registry explicitly shows that no project Skills are yet approved or installed.
- **W002-AC-07:** W001 is marked DONE with its PR #1, independent review and merge evidence available on GitHub.
- **W002-AC-08:** All documentation links resolve and `docs/REQUIREMENT.md` is unchanged.

## Review focus

Check for contradictions, copied or silently promoted requirements, missing owners, traceability gaps, overstated review/CI evidence, and incorrect Work Item status.

## Completion

Require the documentation checks and independent review PASS for the exact PR revision, resolve BLOCKER/MAJOR findings, merge to `develop`, and retain the feature branch.
