# Development Standards

## Supported toolchain

- Node.js 24 and pnpm 10.29.2, matching the root package manager declaration and CI.
- Go 1.27.2 is pinned in CI and `apps/server/go.mod`; update the patch version after reviewing official security releases and rerunning vulnerability scans.
- Dependency versions are recorded in manifests and lockfiles; CI uses frozen installation.

## Local validation

Run `pnpm check` before opening or updating a pull request. It runs formatting checks, lint, TypeScript checking, protocol schema fixtures, Go tests, Go vet, and the frontend production build. Security advisory checks are separately run with `pnpm audit --prod --audit-level high` and `cd apps/server && go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...`.

Oxlint is used for JS/TS linting because the current `typescript-eslint` parser does not support the project's TypeScript 7.0.2 compiler. React plugin rules are enabled; `react-in-jsx-scope` is allowed because the application uses React's automatic JSX runtime.

Formatting changes may be applied with `pnpm format:write`. Go formatting is performed with `gofmt`; `pnpm go:format:check` reports failure without changing files.

## Pull request CI

`.github/workflows/ci.yml` runs on pull requests targeting `develop` or `main` and on pushes to those integration/release branches. The `validate` job checks frozen dependency installation, formatting, lint, type checking, Go tests/vet, frontend build, JavaScript production dependency advisories, and Go known vulnerabilities.

The workflow grants only `contents: read`, uses no secrets, applies a job timeout, and cancels superseded runs. Action references currently use maintained major-version tags; update them deliberately and review upstream changes. Dependency audit results are advisory signals and do not prove that the application is secure. A new vulnerability finding must be triaged and recorded rather than hidden by weakening the threshold.

## Required GitHub checks

The check is named `validate` in workflow `CI`. Confirm the exact status-check context on a successful GitHub run before adding it as a required branch rule for `main` and `develop`. This Work Item does not change branch rules based on an unobserved check name.

## Test growth

The current starter has no Go test files and no browser test suite. Add focused unit and integration tests with the Work Item that introduces the behavior, then make those suites mandatory in CI. Do not report absent suites as passing tests; report them as not applicable or not yet introduced where appropriate.
