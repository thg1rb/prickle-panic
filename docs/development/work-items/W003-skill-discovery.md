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

## Discovery evidence

Discovery was performed with `pnpm dlx skills@1.7.2 find` on 2026-10-19. Results below are untrusted candidate leads, not security or quality approvals. No candidate in this table was installed or approved. Search output is intentionally summarized to preserve useful evidence without turning the registry into a marketplace catalog.

| Category | Query | Representative result(s) observed | Disposition |
|---|---|---|---|
| Go practices | `Go best practices` | `samber/cc-skills-golang` testing, error-handling, security, concurrency, and lint Skills; `openai/skills@security-best-practices` | Not reviewed; unapproved |
| TypeScript / React | `React TypeScript performance` | `vercel-labs/agent-skills@vercel-react-best-practices`; `giuseppe-trisciuoglio/developer-kit@react-patterns` | Not reviewed; unapproved |
| Phaser / game development | `Phaser game development` | `playableintelligence/game-creator@phaser`; several Phaser 3/4 community results | Not reviewed; unapproved; Phaser version compatibility unknown |
| Accessibility / Playwright | `accessibility testing Playwright` | `currents-dev/playwright-best-practices-skill`; `pedronauck/skills@a11y-testing` | Not reviewed; unapproved |
| WebSocket / application security | `WebSocket security code review` and `WebSocket application security` | Search returned broad security skills, including `semgrep/skills@code-security`; no clearly scoped WebSocket result | No candidate selected; unapproved |
| Architecture / quality / review / refactoring | `software architecture code quality refactoring` | `wondelai/skills@improve-code-quality`, `@architecture-optimization`; `nelpos/code-review-plugin@code-review` | Not reviewed; unapproved |
| UI/UX | `UI UX design accessibility` | `anthropics/skills@frontend-design`; `vercel-labs/agent-skills@web-design-guidelines` | Not reviewed; unapproved |
| Docker / GitHub Actions | `Docker GitHub Actions` | `samber/cc-skills-golang@golang-continuous-integration`; `davila7/claude-code-templates@github-actions-creator` | Not reviewed; unapproved |
| Documentation / SDD | `technical documentation spec driven development` | `mattpocock/skills@grill-with-docs`, `@to-spec`, `@to-tickets` | Not reviewed; unapproved |

The initial `find-skills` source and installer received a deeper review because it is mandatory. Future candidates require source authenticity, maintainer, revision, license, scripts and referenced-command inspection, compatibility, security review, and explicit registry approval before trusted use. Search popularity alone is not sufficient.

## Review focus

Verify official-source authenticity, pinned revision integrity, project-only installation, absence of global mutations, correct risk characterization, and whether restricted approval is appropriate given the marketplace prompt-injection exposure.

## Completion

Require evidence-based independent review PASS for the exact PR revision, resolve BLOCKER/MAJOR findings, merge to `develop`, and retain the feature branch.
