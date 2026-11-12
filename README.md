# wubzduh

[New music releases](https://wubzduh.grattafiori.dev) from a curated list of EDM artists on Spotify.
The Go service stores results in Postgres, fetches every 12 hours, and removes releases older than 14 days.

## Run locally

Copy [.env.example](.env.example) to `.env` and supply Spotify credentials and database credentials.
Non-secret configuration lives in [config.yaml](config.yaml).

The [Compose file](docker-compose.yml) requires the external `caddy-network`.
If absent, create it once with `docker network create caddy-network`.
Initialize the database, then expose the app only on localhost:

```sh
docker compose run --rm wubzduh ./wubzduh migrate up
docker compose run --rm -p 127.0.0.1:8080:8080 wubzduh
```

Open `http://localhost:8080`.
Plain `docker compose up` does not publish an application port.
Use the same container command form for `artists list`, `artists add "Tycho"`, `fetch`, or `purge`.
The CLI's `--help` lists commands and flags.

Production runs on K3s. Follow [potatoserver operations](https://github.com/rossgrat/potatoserver/blob/main/docs/operations.md) rather than the retired Docker deployment targets.
