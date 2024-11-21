# wubzduh.com

https://wubzduh.grattafiori.com

## A New Release feed from a list of hand-curated EDM artists
Wubzduh is a multi-threaded Golang web server using a postgres database that integrates with the Spotify API to provide viewers with a feed of new music from a list of artists as soon as that new music is relased (12:05 UTC).

Entries that are older than two weeks are cleared out by a cleanup thread at 18:00 UTC.

## Local Development Setup
1. Create an env.txt with the following fields:
```
DB_USERNAME=<username>
DB_PASSWORD=<password>
DB_NAME=<databaseName>
SPOTIFY_CLIENT_ID=<your-client-id>
SPOTIFY_CLIENT_SECRET=<your-client-secret>
```
2. Use the `src/config/setup.sh` to install postgres if necessary, as well as create and initialize the database and neccessary roles

# Server

## Caddy
Caddy is used for both SSL, HTTP to HTTPS redirection, and as a reverse proxy so I can have more than one application running per machine. The CaddyFile is present in `cmd/web/deploy`

Caddy can be installed with the instructions [here](https://caddyserver.com/docs/install)

### Quick Reference
May have to run any of these commands with `sudo` for the correct permissions.
- `caddy start` - Starts Caddy in the background
- `caddy run` - Starts Caddy in the foreground
- `caddy adapt` - Reloads the Caddyfile, must be in the directory of the Caddy file, else need to use `caddy adapt --config /path/to/file`
- `ss -tulnp` - Sometimes it is neccessary to see what ports are already in use, as Caddy may throw an error if you have already started it, use the `ss` socket investigator for this
    - `-t` displays TCP info
    - `-u` displays UDP info
    - `-l` displays listening sockets
    - `-n` shows ports numerically instead of by name
    - `-p` shows processes using the sockets
- Caddy has a systemd unit that should be made use of, the default config for Caddy is in `/etc/caddy/Caddyfile`. In order to prevent duplicate Caddyfiles, it is neccessary to use the above commands with `--config /etc/caddy/Caddyfile`

## Initial Deployment
- Initial deployment for `web` is done with `remote-setup.sh`
- Initial deployment for `cli` is the same as continous deployment, just use the `build-and-deploy.sh`

## Continous Deployment
- Continuous deployment for `web` is done with `build-and-deploy.sh`, to change the daemon, use `deploy-daemon.sh`
- Continuous deployment for `cli` is the same as initial deployment, just use the `build-and-deploy.sh`


# TODO
- Server is crashing for some reason, figure out why
    - Setup logging for Caddy
- Add CLI functionality to add no artists if the search does not populate any valid artists
- Add throttling for bots
- Get rid of all of the ugly relative directory strings `../../, etc` in favor of variables
- Allow specification of HOST in config files
- Do something intelligent with page and thread errors instead of just crashing
- Add logger to record visits, use zlogger package
    - Use lnav to examine logs
- Use JSON for config files, not text and environment variables
- Postgres is too heavyweight for this application, consider using sqlite or leveldb
- Add some real styling, move feed to the middle of the screen for wider screens

## Future Ideas
- Write Playlist module
    For all new albums
        Get track and add it to playlists
    For all albums older than one week  
        Remove entry from playlist
- CLI Module Additions
    Print contents of Album table
    Print contents of Track table
- Allow sorting of albums by genres
    Add logic in FeedHandler to grab tracks based on genre parameter/feed/?genre=house
    https://golangbyexample.com/net-http-package-get-query-params-golang/
- Tag albums release day of with NEW RELEASE when displaying on website?
- Add total duration to album displayed information
- Add last release date field in artists tab
- Make this whole thing into a single-page web application and a progressive web app, save on server costs
- Do push notifications as a PWA
- Use Spotify OAUTH to allow people to login to their spotify profiles, retrieve their liked artists, push notifications based on those

