package main

import (
	"html/template"
	"net/http"

	wubzduh "github.com/rossgrat/wubzduh/src"
)

type ArtistPage struct {
	Title   string
	Artists []wubzduh.Artist
}

type FeedPage struct {
	Title  string
	Albums []wubzduh.Album
}

type PlaylistPage struct {
	Title    string
	Playlist string
}

var Templates = template.Must(template.ParseFiles("feed.html", "artists.html", "navbar.html", "album.html"))

// Redirects to Feed View
func FeedRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/feed/", http.StatusFound)
}

// Serves main feed fiew of all new albums and tracks
func FeedViewHandler(w http.ResponseWriter, r *http.Request) {
	var Page FeedPage
	//Get all albums to be shown on page
	Page.Albums = wubzduh.GetAlbumsAndTracks(DB)

	Page.Title = "Feed"
	err := Templates.ExecuteTemplate(w, "feed.html", &Page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Lists all artists in DB
func ArtistsViewHandler(w http.ResponseWriter, r *http.Request) {
	var Page ArtistPage
	//Get all artists to be shown on the page
	Page.Artists = wubzduh.GetAllArtists(DB)

	Page.Title = "Artists"
	err := Templates.ExecuteTemplate(w, "artists.html", &Page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Serves favicon that shows up on tab
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/music.png")
}
