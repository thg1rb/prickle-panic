# W004 — Development Standards & CI

Status: READY_FOR_REVIEW
Type: Foundation

## Goal

Establish repeatable formatting, linting, type-checking, Go validation, security checks, and pull-request CI gates for the existing monorepo.

## Scope

- Define and configure formatting and lint standards for current TypeScript/React and Go code.
- Add repository scripts for local validation.
- Add GitHub Actions validation for `feature/*` → `develop` and `develop` → `main` pull requests.
- Validate package lock reproducibility, frontend typecheck/build, Go formatting/test/vet, and dependency security advisories.
- Document checks, expected local commands, and CI limits.
- Mark W003 DONE with its exact review and merge evidence.

## Non-goals

- Feature code, new gameplay tests, end-to-end browser tests, deployment, release automation, or new external services.
- Required GitHub status-check enforcement until the checks have run successfully on GitHub and their exact names are verified.
- Changes to requirements, Git branch policy, or game rules.

## Requirement references

- `docs/REQUIREMENT.md` §§25–27, §29, §32, §40–41.
- AC-025, AC-026, SEC-012, RISK-006.
- Roadmap W004 — Development Standards & CI.

## Acceptance criteria

- **W004-AC-01:** Contributor-facing standards and commands are documented in English and linked from the agent entry point.
- **W004-AC-02:** Root scripts provide consistent formatting, lint, typecheck, Go test, vet, and build validation.
- **W004-AC-03:** Go formatting is checked without rewriting files in CI.
- **W004-AC-04:** CI runs on PRs targeting `develop` and `main`, and on updates to those branches, with least-privilege permissions and no secrets required for validation.
- **W004-AC-05:** CI uses frozen package installation, runs frontend checks/build and Go checks, and performs dependency advisory checks without claiming that audit tools prove the absence of vulnerabilities.
- **W004-AC-06:** Existing starter code passes the configured checks; no gameplay behavior is added.
- **W004-AC-07:** W003 status is synchronized to DONE with review/merge evidence.
- **W004-AC-08:** Required check names are reported for later branch-protection configuration only after observed on GitHub.

## Review focus

Inspect action permissions and trigger scope, dependency pinning, lockfile use, false confidence from security checks, local/CI parity, and whether required branch checks are withheld until their actual names are known.

## Technical notes

The TypeScript compiler baseline is TypeScript 7.0.2. The current `typescript-eslint` release rejects TypeScript 7.0, so the repository uses Oxlint for JS/TS linting rather than accepting an unsupported parser combination. The React `react-in-jsx-scope` rule is disabled because this Vite app uses the automatic JSX runtime. CI uses Go 1.27.2 after `govulncheck` identified fixed standard-library vulnerabilities in the prior local Go 1.25.1 toolchain; the official Go release history identifies 1.27.2 as the current security patch on 2026-10-19. Initial CI exposed an outdated `govulncheck` parser panic with Go 1.27, so the pinned scanner was updated from v1.1.4 to v1.9.0; the newer scanner completed successfully locally.

## Completion

Require all applicable local checks to pass, independent review PASS for the exact PR revision, blocking findings resolved, PR merged to `develop`, and feature branch retained.
