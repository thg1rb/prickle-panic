# Prickle Panic!

Prickle Panic! is a Thai-first browser party game. The MVP target is a private multiplayer Speed Battle for 2–6 players.

The product requirements and gameplay rules are maintained in [`docs/REQUIREMENT.md`](docs/REQUIREMENT.md). Development and project documentation use English; player-facing MVP text uses Thai.

## Prerequisites

- Node.js 22.12 or newer
- pnpm 10.29.2 or newer
- Go 1.25 or newer

The frontend dependencies are pinned by `pnpm-lock.yaml`. The Go server currently uses only the standard library.

## Development

Install frontend dependencies:

```sh
pnpm install
```

Run the web development server:

```sh
pnpm dev
```

Build and type-check the web application:

```sh
pnpm build
pnpm typecheck
```

Run the Go health-check server:

```sh
pnpm server:run
```

The server listens on port `8080` by default. Set `PORT` to use another port. Its health endpoint is `GET /healthz`.

## Repository layout

- `apps/web/` — React, TypeScript, Phaser, and Tailwind web client
- `apps/server/` — Go HTTP server foundation
- `docs/REQUIREMENT.md` — authoritative product and development requirements
