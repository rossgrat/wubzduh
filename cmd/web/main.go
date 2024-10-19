package main

import (
	"context"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/rossgrat/wubzduh/src/lib/db"
	"github.com/rossgrat/wubzduh/src/lib/threads"
	"github.com/rossgrat/wubzduh/src/lib/util"
	"github.com/zmb3/spotify/v2"
)

const (
	time1 = iota
	time2
	numTimes = 2
)

type ArtistPage struct {
	Title   string
	Artists []db.Artist
}
type FeedPage struct {
	Title  string
	Albums []db.Album
}
type PlaylistPage struct {
	Title    string
	Playlist string
}

var DB *sql.DB
var Client *spotify.Client
var Ctx context.Context
var Lock sync.RWMutex

var FetchTimer *time.Timer
var PurgeTimer *time.Timer

var Templates = template.Must(
	template.ParseFiles(
		"templates/feed.html",
		"templates/artists.html",
		"templates/navbar.html",
		"templates/album.html",
	),
)

func RequestLogger(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
		log.Printf("%s - %s %s\n", r.RemoteAddr, r.Method, r.RequestURI)
	})
}

// Redirects to Feed View
func FeedRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/feed/", http.StatusFound)
}

// Serves main feed fiew of all new albums and tracks
func FeedViewHandler(w http.ResponseWriter, r *http.Request) {
	var Page FeedPage
	albums, err := db.GetAlbums()
	if err != nil {
		log.Fatal(err.Error())
	}
	Page.Albums = albums

	Page.Title = "Feed"
	if err := Templates.ExecuteTemplate(w, "feed.html", &Page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Lists all artists in DB
func ArtistsViewHandler(w http.ResponseWriter, r *http.Request) {
	var Page ArtistPage
	artists, err := db.GetAllArtists()
	if err != nil {
		log.Fatal(err.Error())
	}
	Page.Artists = artists

	Page.Title = "Artists"
	if err := Templates.ExecuteTemplate(w, "artists.html", &Page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Serves favicon that shows up on tab
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/music.png")
}

func main() {
	db.Connect()
	Client, Ctx = util.ConnectToSpotify()

	go FetchThread(time2)
	go PurgeThread()

	mux := http.NewServeMux()
	mux.HandleFunc("/artists/", ArtistsViewHandler)
	mux.HandleFunc("/feed/", FeedViewHandler)
	mux.HandleFunc("/favicon.ico", FaviconHandler)
	mux.HandleFunc("/", FeedRedirect)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", RequestLogger(mux)))
}

// The Fetch thread is responsible for checking for new music at 12:30am and 11:30pm
// each day
func FetchThread(timeCase int) {
	currentTime := time.Now()
	var diff time.Duration
	switch timeCase {
	case time1:
		// 12:30am the next day
		currentTimePlusDay := currentTime.AddDate(0, 0, 1)
		targetTime := time.Date(currentTimePlusDay.Year(), currentTimePlusDay.Month(), currentTimePlusDay.Day(), 0, 30, 0, 0, currentTimePlusDay.Location())
		diff = targetTime.Sub(currentTime)
		log.Printf("Set Fetch Timer for 12:30am Tomorrow. ETE: %s\n", diff)
	case time2:
		//11:30pm the same day
		targetTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 30, 0, 0, currentTime.Location())
		diff = targetTime.Sub(currentTime)
		log.Printf("Set Fetch Timer for 11:30pm Today. ETE: %s\n", diff)
	}

	FetchTimer = time.NewTimer(diff)
	<-FetchTimer.C
	threads.Fetch(true)
	go FetchThread((timeCase + 1) % numTimes)
}

// The purge thread removes any music that has been released more than 14 days ago
func PurgeThread() {
	currentTime := time.Now()

	currentTimePlusDay := currentTime.AddDate(0, 0, 1)
	targetTime := time.Date(currentTimePlusDay.Year(), currentTimePlusDay.Month(), currentTimePlusDay.Day(), 12, 0, 0, 0, currentTimePlusDay.Location())
	diff := targetTime.Sub(currentTime)

	log.Printf("Set Purge Timer for 12:00pm Tomorrow. ETE: %s\n", diff)

	PurgeTimer = time.NewTimer(diff)
	<-PurgeTimer.C
	threads.Purge()
	go PurgeThread()
}
