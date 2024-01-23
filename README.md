# wubzduh.com
## A New Release feed from a list of hand-curated EDM artists
Wubzduh is a multi-threaded Golang web server using a postgres database that integrates with the Spotify API to provide viewers with a feed of new music from a list of artists as soon as that new music is relased (12:05 UTC).

Entries that are older than two weeks are cleared out by a cleanup thread at 18:00 UTC.
