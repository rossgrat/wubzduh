package db

import "time"

type Artist struct {
	ID        int
	Name      string
	SpotifyID string
}

type Track struct {
	ID              int
	Title           string
	SpotifyID       string
	Number          int
	DurationMS      int
	DurationMinutes string
	DurationSeconds string
	AlbumID         int
}

type Album struct {
	ID          int
	Title       string
	SpotifyID   string
	ArtistID    int
	ArtistName  string
	ReleaseDate time.Time
	CoverartURL string
	Type        string
	URL         string
	Tracks      []Track
}
