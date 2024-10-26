# wubzduh.com

https://wubzduh.grattafiori.com

## A New Release feed from a list of hand-curated EDM artists
Wubzduh is a multi-threaded Golang web server using a postgres database that integrates with the Spotify API to provide viewers with a feed of new music from a list of artists as soon as that new music is relased (12:05 UTC).

Entries that are older than two weeks are cleared out by a cleanup thread at 18:00 UTC.

## Local Development Setup
1. Create a postgres database with `CREATE DATABASE databaseName;`
2. Create a postgres role and database for the server to use, with prileges necessary to edit the database.
```
CREATE ROLE serverUser;
CREATE DATABASE serverUser;
\c databaseName
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO serverUser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO serverUser;
```
3. Create an env.txt with the following fields:
```
DB_USERNAME=<username>
DB_PASSWORD=<password>
DB_NAME=<databaseName>
SPOTIFY_CLIENT_ID=<your-client-id>
SPOTIFY_CLIENT_SECRET=<your-client-secret>
```

# Server

## Caddy
Caddy is used for both SSL, HTTP to HTTPS redirection, and as a reverse proxy so I can have more than one application running per machine. The CaddyFile is present in `cmd/web/deploy`

Caddy can be installed with the instructions [here](https://caddyserver.com/docs/install)

### Quick Reference
- `caddy start` - Starts Caddy in the background
- `caddy run` - Starts Caddy in the foreground
- `caddy adapt` - Reloads the Caddyfile, must be in the directory of the Caddy file, else need to use `caddy adapt --config /path/to/file`


## Initial Deployment to Server
1. Copy over wubzduh.service to correct directory
2. Load wubzdub.service into systemctl daemon
3. Create database and initialize tables

## Deployment to Server
1. Build web executable
2. Copy web executable, templates, and enviroment file to server
3. Restart service daemon

# TODO
- Do something intelligent with page and thread errors instead of just crashing
- Add logger to record visits, use zlogger package
    - Use lnav to examine logs
- Use JSON for config files, not text and environment variables
- Postgres is too heavyweight for this application, consider using sqlite or leveldb
- Add some real styling 
-

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

