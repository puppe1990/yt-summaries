# YT Summaries

Super­-summaries of YouTube videos: a transcript recap **plus** a full read of the comment section, published as an editorial reader page.

**Live:** https://resumos.apps.gestaobem.com · site copy is pt-BR

## How it works

1. The video's transcript (`data/raw/<videoId>/transcricao.md`) and comment export (`comentarios.json`) are archived in the repo.
2. A structured super-summary is written to `data/videos/<videoId>.json`: kicker, editorial headline, thesis plates, numbered sections, quotes with PT translations, and the audience read (top comments, themes, insights, consensus vs. open fight, verdict).
3. On boot the app upserts every payload into SQLite — idempotent, safe on every deploy.
4. `/` renders the index as cards (the newest item, the "next" one, is always first, flagged ▶); `/resumos/<slug>` renders the full super-summary.
5. New summaries can be published at any time via `POST /api/resumos` (Bearer `ADMIN_TOKEN`) — upserted straight into SQLite, no deploy needed. Payloads are source-agnostic (`source_id`/`source`/`source_type`): YouTube, podcasts, reels, articles.

## Stack

- Go 1.26 (`net/http`) + [Amarra Views + Drive](https://github.com/puppe1990/amarra-cais)
- SQLite (`modernc.org/sqlite`, no CGO) + embedded migrations
- Tailwind CSS · Fraunces / Source Serif 4 / IBM Plex Mono
- Deployed on Cleat, auto-deploy on push to `main`

## Development

```bash
amarra-cais install   # npm + go mod tidy
amarra-cais dev       # or: make build && PORT=:8091 ./bin/server
make test             # go test ./... -race
make ci               # test + lint + format-check
```

Agent conventions live in [`AGENTS.md`](AGENTS.md).

## Brand (icon + Open Graph)

The scaffold placeholders in `web/static/` are replaced by the darkroom/paper mark
(film rail + thesis plate + play) and a 1200×630 social card.

```bash
scripts/generate-brand.sh
```

Sources: `web/static/icons/icon.svg`, `web/static/icons/icon-maskable.svg`,
`design/og.html`. The maskable icon keeps the mark inside the inner 80% so
Android circle crops do not clip the sprockets or the play tip.

## Layout

| Path                  | Responsibility                                                    |
| --------------------- | ----------------------------------------------------------------- |
| `data/videos/*.json`  | super-summary payloads (imported on boot, embedded in the binary) |
| `data/raw/<videoId>/` | original transcript + comments export                             |
| `internal/store/`     | SQLite, migrations, videos upsert/query                           |
| `internal/handlers/`  | index cards, video page, designed 404                             |
| `web/templates/`      | `pages/` (home, video, notfound) + `partials/resumo.html`         |
