# wubzduh

https://wubzduh.grattafiori.dev

## A new release feed from hand-curated EDM artists

Wubzduh monitors Spotify for new music releases from a curated list of artists and serves a web feed. Background workers fetch new releases every 12 hours and purge entries older than 14 days.

## Architecture

- **Go** web server with cobra CLI, viper config, and PostgreSQL
- **Docker** deployment to [potatoserver](https://github.com/rossgrat/potatoserver) via Cloudflare Tunnel
- **CI/CD** via GitHub Actions — pushes to `main` build and publish a Docker image to GHCR

```
Internet → Cloudflare → cloudflared → Caddy → wubzduh:8080
                                                  ↕
                                              PostgreSQL
```

## Local Development

1. Copy and configure environment variables:
   ```bash
   cp .env.example .env
   nano .env
   ```

2. Start services:
   ```bash
   docker compose up
   ```

3. Visit http://localhost:8080

## CLI

```bash
wubzduh serve                    # Start the web server
wubzduh fetch                    # Manually fetch new releases
wubzduh fetch --no-date-check    # Fetch all latest releases regardless of date
wubzduh purge                    # Purge releases older than 14 days
wubzduh artists list             # List tracked artists
wubzduh artists add "Tycho"      # Search Spotify and add an artist
wubzduh artists search "ZHU"     # Search artists in the database
wubzduh migrate up               # Apply database migrations
wubzduh migrate down             # Roll back migrations
wubzduh migrate version          # Show current migration version
```

All commands accept `--config <path>` (default: `./config.yaml`).

## Deployment

```bash
make deploy    # Deploy to potatoserver
make stop      # Stop the service
make logs      # Tail logs
```

## Project Structure

```
cmd/           Cobra CLI commands (serve, fetch, purge, artists, migrate)
internal/      Business logic (config, db, service, worker)
plugins/       External integrations (spotify, logger)
web/           HTTP server, handlers, middleware, templates, static assets
migrations/    SQL migrations (embedded in binary via golang-migrate)
deploy/        Production docker-compose
```

## Configuration

- `config.yaml` — non-secret config (server port, DB host, log settings)
- `.env` — secrets (`DB_USER`, `DB_PASSWORD`, `SPOTIFY_CLIENT_ID`, `SPOTIFY_CLIENT_SECRET`)
