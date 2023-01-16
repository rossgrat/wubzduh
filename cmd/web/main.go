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

func RequestLogger(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)

		log.Printf("%s - %s %s\n", r.RemoteAddr, r.Method, r.RequestURI)
	})

}

func main() {

	//Initialize global variables
	DB = wubzduh.ConnectToDB()
	Client, Ctx = wubzduh.ConnectToSpotify()

	//Start the first purge and fetch threads
	go FetchThread(time2)
	go PurgeThread()

	//Set handler functions and start web server
	mux := http.NewServeMux()
	mux.HandleFunc("/artists/", ArtistsViewHandler)
	mux.HandleFunc("/feed/", FeedViewHandler)
	mux.HandleFunc("/playlist/", PlaylistViewHandler)
	mux.HandleFunc("/favicon.ico", FaviconHandler)
	mux.HandleFunc("/", FeedRedirect)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", RequestLogger(mux)))
}
