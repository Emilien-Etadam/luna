# Luna

A self-hosted calendar aggregator with a single, fast UI for every calendar you care about.

![Luna in light and dark mode](./documentation/pictures/light-dark.png)

## What it does

- **One UI for every calendar.** Connect CalDAV servers (Nextcloud, Radicale, Baikal…), iCal feeds, and Google Calendar — Luna shows them side by side.
- **Multi-user.** Per-user accounts, sources, calendars, and preferences. Designed for self-hosting on a home server or VPS.
- **Public read-only mode.** Expose a curated view of your calendars to unauthenticated visitors without giving them write access.
- **Month and agenda views.** Navigate by month, week, day, or jump straight to a chronological agenda anchored on today.
- **VS Code Dark+ theme out of the box.** Tokens-driven design system in `frontend/static/theme/tokens.css`; a light theme is available with the same token names. Additional themes ship in `frontend/static/themes/`.
- **Sensible defaults, no fluff.** No shadows, no decorative animations, transitions capped at 100 ms, padding on a 4 px grid.

## Architecture

| Component | Stack | Notes |
|-----------|-------|-------|
| `frontend/` | SvelteKit 2 + Svelte 5 (runes), Vite, Bun runtime | Served by `svelte-adapter-bun` |
| `backend/`  | Go (`luna-backend`) | REST API, CalDAV/iCal/Google adapters, task scheduler |
| Database    | PostgreSQL 16 | Stores users, sources, calendars, cached events |

The frontend talks to the backend over HTTP. In public read-only mode, the SvelteKit hooks bypass `/login` redirects and route data fetching to `/api/public/*` endpoints exposed by the backend.

## Quick start (Docker)

Requirements: `make`, `docker`, `docker compose`.

```bash
git clone https://github.com/Emilien-Etadam/luna.git
cd luna
# Edit docker-compose.yml: set PUBLIC_URL, DB_PASSWORD, etc.
make run
```

`make` targets:

- `make run` — build and start the stack in detached mode (`docker compose up -d --build`)
- `make up` — start without rebuilding
- `make down` — stop the stack
- `make restart` — restart containers
- `make build` — rebuild images with `--no-cache`
- `make purge` — stop the stack and **wipe** the Postgres volume at `/srv/luna/postgres`

The default frontend port is `8080`. Persistent data lives in `/srv/luna/data` (backend) and `/srv/luna/postgres` (database) on the host.

A working `docker-compose.yml` is provided at the repository root; copy and adapt it before going to production.

## Local development

### Frontend

```bash
cd frontend
npm install
npm run dev      # vite dev server
npm run check    # svelte-kit sync + svelte-check
npm run build    # production build (svelte-adapter-bun)
```

### Backend

```bash
cd backend/src
go build ./...
go run ./luna-backend
```

The backend reads its configuration from environment variables (see `docker-compose.yml` for the full list).

## Theming

All colors, typography, spacing, radii, and scrollbar styling are centralised in [`frontend/static/theme/tokens.css`](./frontend/static/theme/tokens.css). Components consume CSS variables (`--bg-editor`, `--fg-primary`, `--border-default`, `--accent-blue`…) — no hardcoded hex values outside that file.

To create a new theme:

1. Duplicate `tokens.css` and override the variables under your own `html[data-theme="..."]` selector.
2. Drop additional theme stylesheets into `frontend/static/themes/`.

Further details: [`documentation/themes.md`](./documentation/themes.md).

## Project layout

```
.
├── backend/         Go API server (CalDAV, iCal, Google, auth, DB)
├── frontend/        SvelteKit app (UI, routing, theming)
├── documentation/   Deployment, security, API, theming, contribution
├── scripts/         Helper scripts
├── docker-compose.yml
└── Makefile
```

## Documentation

- [Deployment guide](./documentation/deployment.md)
- [API reference](./documentation/api.md)
- [Security & privacy](./documentation/security.md)
- [Theming](./documentation/themes.md)
- [Contribution guide](./documentation/contribution.md)

## Status

Luna is approaching a `1.0.0` release. Until then expect occasional breaking changes — back up your calendars before updating and be prepared to wipe the database between major versions. Track progress on the [development roadmap](https://todo.opisek.net/share/dvEazOyRLEYThqxohVosnqKskYLyoZ4nS8rQ63G1/auth?view=280).

## License

To be defined. Copyright © 2026 Kacper Darowski (Opisek) and contributors.
