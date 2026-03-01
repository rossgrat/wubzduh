package web

import (
	"net/http"

	"github.com/rossgrat/wubzduh/internal/db"
)

type feedPage struct {
	Title  string
	Albums []db.Album
}

type artistPage struct {
	Title   string
	Artists []db.Artist
}

func (s *Server) feedRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/feed/", http.StatusFound)
}

func (s *Server) feedHandler(w http.ResponseWriter, r *http.Request) {
	albums, err := s.store.GetAlbums()
	if err != nil {
		s.logger.Error("failed to get albums", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	page := feedPage{Title: "Feed", Albums: albums}
	if err := s.templates.ExecuteTemplate(w, "feed.html", &page); err != nil {
		s.logger.Error("failed to render feed template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) artistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := s.store.GetAllArtists()
	if err != nil {
		s.logger.Error("failed to get artists", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	page := artistPage{Title: "Artists", Artists: artists}
	if err := s.templates.ExecuteTemplate(w, "artists.html", &page); err != nil {
		s.logger.Error("failed to render artists template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) faviconHandler(w http.ResponseWriter, r *http.Request) {
	data, err := embeddedStatic.ReadFile("static/music.png")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(data)
}
