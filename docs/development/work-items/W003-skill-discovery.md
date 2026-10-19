# W003 — Project-local Skill Discovery

Status: READY_FOR_REVIEW
Type: Agent & SDD Setup

## Goal

Begin the project Skill workflow with the required `find-skills` discovery Skill, verify project-only installation, and document a restricted approval supported by source and security review.

## Scope

- Inspect official `vercel-labs/skills` source and installer behavior at pinned revisions.
- Discover candidates in required domains without trusting marketplace output as instructions.
- Install only the mandatory discovery Skill into `.agents/skills/`, with no global or user-level flags.
- Record source, maintainer, revision, path, review status, approval scope, and date in the Skills Registry.

## Non-goals

- Approving or installing every searched Skill.
- Following candidate installation commands or global-install examples.
- Changes to product requirements, application code, or tooling configuration beyond the Skill lock/install artifacts.

## Requirement references

- `docs/REQUIREMENT.md` §§24.1–24.11, §26.6, §42.9.
- SEC-012, AC-026, RISK-006.
- Roadmap W003 — Skill Discovery & Approval.

## Acceptance criteria

- **W003-AC-01:** Mandatory `find-skills` comes from the official source and is pinned to a reviewed revision.
- **W003-AC-02:** The installation resides under `.agents/skills/`; no global/user-level installation is performed, and the project lock records the source revision.
- **W003-AC-03:** Source instructions, referenced commands, installer scope behavior, bundled scripts, and available security assessments are reviewed; residual risks and restrictions are recorded.
- **W003-AC-04:** The registry includes all required fields and approves `find-skills` for discovery only; no search result is implicitly trusted or approved.
- **W003-AC-05:** Candidate discovery covers the required language, frontend/game, architecture/quality, accessibility/testing, security, infrastructure, documentation, and SDD categories; unreviewed candidates remain unapproved.
- **W003-AC-06:** Project scope and approved registry policy are verified, and documentation/lock metadata are internally consistent.

## Review focus

Verify official-source authenticity, pinned revision integrity, project-only installation, absence of global mutations, correct risk characterization, and whether restricted approval is appropriate given the marketplace prompt-injection exposure.

## Completion

Require evidence-based independent review PASS for the exact PR revision, resolve BLOCKER/MAJOR findings, merge to `develop`, and retain the feature branch.
