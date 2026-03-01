package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rossgrat/wubzduh/src/lib/db"
	"github.com/rossgrat/wubzduh/src/lib/threads"
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
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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

	timer := time.NewTimer(diff)
	<-timer.C
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

	timer := time.NewTimer(diff)
	<-timer.C
	threads.Purge()
	go PurgeThread()
}

func main() {
	db.Connect()

	go FetchThread(time2)
	go PurgeThread()

	mux := http.NewServeMux()
	mux.HandleFunc("/artists/", ArtistsViewHandler)
	mux.HandleFunc("/feed/", FeedViewHandler)
	mux.HandleFunc("/favicon.ico", FaviconHandler)
	mux.HandleFunc("/", FeedRedirect)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := &http.Server{
		Addr:    ":8080",
		Handler: RequestLogger(mux),
	}

	// Listen for shutdown signals in the background
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-quit
		log.Printf("Received signal %s, shutting down...\n", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v\n", err)
		}
	}()

	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v\n", err)
	}
	log.Println("Server stopped")
}
