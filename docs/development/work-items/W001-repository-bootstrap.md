# W001 — Repository Bootstrap

Status: READY_FOR_REVIEW
Type: Foundation

## Goal

Create a minimal, reproducible web and Go server workspace that provides the verified starting point for later SDD Work Items.

## Scope

- Establish the pnpm workspace and React, TypeScript, Phaser, Tailwind, and Vite web shell.
- Establish a Go HTTP server shell with a health endpoint.
- Document local prerequisites and development, build, and type-check commands.
- Keep the visible MVP shell text in Thai.

## Non-goals

- Gameplay, card rules, guest sessions, Rooms, WebSocket protocol, assets, CI, production deployment, and Skill installation.
- Changing gameplay parameters or requirements.

## Requirements

- `docs/REQUIREMENT.md` §18.2, §18.5–18.7, §30.3, §31–32.
- NFR-019, NFR-026, NFR-030.
- Roadmap W001 — Repository Bootstrap.

## Expected behavior

- A fresh checkout can install the locked frontend dependencies and run the web shell.
- The web shell builds and type-checks with the documented commands.
- The Go server starts using the standard library and responds `200 OK` with `ok` at `GET /healthz`.
- The server uses `PORT` when set and defaults to port 8080.
- No authoritative gameplay or player session behavior is introduced.

## Edge cases and errors

- Missing `PORT` uses the documented default.
- An invalid or unavailable port causes the server process to report the standard listener error and exit.
- Missing frontend root element fails fast with a descriptive error.

## Acceptance criteria

- **W001-AC-01:** `pnpm install --frozen-lockfile` succeeds from a clean checkout.
- **W001-AC-02:** `pnpm build` and `pnpm typecheck` succeed.
- **W001-AC-03:** `go test ./...` and `go vet ./...` complete successfully in `apps/server` (no tests may be reported as no test files, not as executed test coverage).
- **W001-AC-04:** Running the documented server command and requesting `/healthz` returns HTTP 200 and `ok`.
- **W001-AC-05:** MVP-facing shell copy is Thai, excluding the approved product name.
- **W001-AC-06:** Dependency versions are recorded in manifests and the pnpm lockfile; Go uses only the standard library.

## Validation plan

Run the commands above on the feature branch. Record each result as PASS, FAIL, NOT_RUN, BLOCKED, or NOT_APPLICABLE. No test coverage is claimed for the starter shell.

## Review focus

Review the actual diff against the listed requirements and acceptance criteria; inspect dependency and lockfile changes, workspace boundaries, Thai-facing text, the health endpoint, and the recorded command results. Review evidence must identify the Work Item, reviewer execution, and exact commit SHA.

## Completion

W001 may be marked DONE only after the required validation evidence and independent review PASS are recorded, blocking findings are resolved, and its feature PR is merged into `develop` using a Merge Commit. Retain the feature branch.
