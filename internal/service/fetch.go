package service

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/rossgrat/wubzduh/internal/db"
	spotifyplugin "github.com/rossgrat/wubzduh/plugins/spotify"
	spotifylib "github.com/zmb3/spotify/v2"
)

type FetchService struct {
	store   *db.Store
	spotify *spotifyplugin.Client
	logger  *slog.Logger
}

type FetchOption func(*FetchService)

func FetchWithLogger(log *slog.Logger) FetchOption {
	return func(f *FetchService) { f.logger = log }
}

func NewFetchService(store *db.Store, spotify *spotifyplugin.Client, opts ...FetchOption) *FetchService {
	f := &FetchService{
		store:   store,
		spotify: spotify,
		logger:  slog.Default(),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *FetchService) Fetch(checkReleaseDate bool) error {
	artists, err := f.store.GetAllArtists()
	if err != nil {
		return fmt.Errorf("getting artists: %w", err)
	}

	newAlbums := f.getArtistsLatestAlbums(artists, spotifylib.AlbumTypeAlbum)
	newSingles := f.getArtistsLatestAlbums(artists, spotifylib.AlbumTypeSingle)
	newAlbums = append(newAlbums, newSingles...)

	for _, a := range newAlbums {
		if checkReleaseDate && !isReleasedToday(a) {
			continue
		}

		albumID, err := f.store.InsertAlbum(a)
		if err != nil {
			f.logger.Error("failed to insert album", "title", a.Title, "error", err)
			continue
		}
		if albumID == -1 {
			f.logger.Debug("album already exists", "title", a.Title)
			continue
		}

		a.ID = albumID
		tracks := f.getAlbumTracks(a)
		for _, track := range tracks {
			if err := f.store.InsertTrack(track); err != nil {
				f.logger.Error("failed to insert track", "title", track.Title, "error", err)
			}
		}
		f.logger.Info("added new album", "title", a.Title, "artist_id", a.ArtistID, "tracks", len(tracks))
	}

	return nil
}

func (f *FetchService) getArtistsLatestAlbums(artists []db.Artist, albumType spotifylib.AlbumType) []db.Album {
	var albums []db.Album
	for _, a := range artists {
		albumResults, err := f.spotify.API().GetArtistAlbums(
			f.spotify.Context(),
			spotifylib.ID(a.SpotifyID),
			[]spotifylib.AlbumType{albumType},
			spotifylib.Limit(1),
		)
		if err != nil {
			f.logger.Error("failed to get albums for artist", "artist", a.Name, "error", err)
			continue
		}
		if len(albumResults.Albums) < 1 {
			continue
		}

		result := albumResults.Albums[0]
		coverartURL := ""
		if len(result.Images) > 1 {
			coverartURL = result.Images[1].URL
		} else if len(result.Images) > 0 {
			coverartURL = result.Images[0].URL
		}

		albums = append(albums, db.Album{
			Title:       result.Name,
			ArtistID:    a.ID,
			ReleaseDate: result.ReleaseDateTime(),
			CoverartURL: coverartURL,
			Type:        strings.ToTitle(string(result.AlbumType)),
			URL:         result.ExternalURLs["spotify"],
			SpotifyID:   result.ID.String(),
		})
	}
	return albums
}

func (f *FetchService) getAlbumTracks(album db.Album) []db.Track {
	trackResults, err := f.spotify.API().GetAlbumTracks(
		f.spotify.Context(),
		spotifylib.ID(album.SpotifyID),
	)
	if err != nil {
		f.logger.Error("failed to get tracks for album", "album", album.Title, "error", err)
		return nil
	}

	var tracks []db.Track
	for _, tr := range trackResults.Tracks {
		tracks = append(tracks, db.Track{
			Title:      tr.Name,
			SpotifyID:  tr.ID.String(),
			Number:     int(tr.TrackNumber),
			DurationMS: int(tr.Duration),
			AlbumID:    album.ID,
		})
	}
	return tracks
}

func isReleasedToday(album db.Album) bool {
	today := time.Now().Format("2006-01-02")
	return album.ReleaseDate.Format("2006-01-02") == today
}
