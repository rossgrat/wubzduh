package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sync"

	wubzduh "github.com/rossgrat/wubzduh/src"
	"github.com/zmb3/spotify/v2"
)

var DB *sql.DB
var Client *spotify.Client
var Ctx context.Context
var Lock sync.RWMutex

func main() {

	//Initialize global variables
	DB = wubzduh.ConnectToDB()
	Client, Ctx = wubzduh.ConnectToSpotify()

	//Start the first purge and fetch threads
	go FetchThread(time2)
	go PurgeThread()

	//Set handler functions and start web server
	http.HandleFunc("/artists/", ArtistsViewHandler)
	http.HandleFunc("/feed/", FeedViewHandler)
	http.HandleFunc("/playlist/", PlaylistViewHandler)
	http.HandleFunc("/favicon.icoa", FaviconHandler)
	http.HandleFunc("/", FeedRedirect)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
