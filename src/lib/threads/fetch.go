package threads

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/rossgrat/wubzduh/src/lib/db"
	"github.com/rossgrat/wubzduh/src/lib/util"
	"github.com/zmb3/spotify/v2"
)

// TOOD: Does it make sense to put the spotify stuff in it's own package, and combine the
// Fetch and Purge functions into their own package?

// For each artist in the database, get the latest release by each artist, may
// be either a full length album or an EP / Single, we make the following
// assumptions here:
// 1. An arist will not release both a new single and a new album in one day
// 2. An artist will not release more than one new single or more than one
// new album in one day
func spotifyGetArtistsLatestAlbums(
	client *spotify.Client,
	ctx context.Context,
	artists []db.Artist,
	albumType spotify.AlbumType,
) (
	albums []db.Album,
) {
	for _, a := range artists {
		albumResults, err := client.GetArtistAlbums(
			ctx,
			(spotify.ID)(a.SpotifyID),
			[]spotify.AlbumType{albumType},
			spotify.Limit(1),
		)
		if err != nil {
			log.Printf("Failed to get albums for artist %s: %v\n", a.Name, err)
			continue
		}
		if len(albumResults.Albums) < 1 {
			continue
		}
		coverartURL := ""
		if len(albumResults.Albums[0].Images) > 1 {
			coverartURL = albumResults.Albums[0].Images[1].URL
		} else if len(albumResults.Albums[0].Images) > 0 {
			coverartURL = albumResults.Albums[0].Images[0].URL
		}
		newAlbum := db.Album{
			Title:       albumResults.Albums[0].Name,
			ArtistID:    a.ID,
			ReleaseDate: albumResults.Albums[0].ReleaseDateTime(),
			CoverartURL: coverartURL,
			Type:        strings.ToTitle(albumResults.Albums[0].AlbumType),
			URL:         albumResults.Albums[0].ExternalURLs["spotify"],
			SpotifyID:   albumResults.Albums[0].ID.String(),
		}
		albums = append(albums, newAlbum)
	}
	return albums
}

func spotifyGetAlbumTracks(
	client *spotify.Client,
	ctx context.Context,
	album db.Album,
) (
	tracks []db.Track,
) {
	trackResults, err := client.GetAlbumTracks(
		ctx,
		(spotify.ID)(album.SpotifyID),
	)
	if err != nil {
		log.Printf("Failed to get tracks for album %s: %v\n", album.Title, err)
		return tracks
	}
	for _, tr := range trackResults.Tracks {
		track := db.Track{
			Title:      tr.Name,
			SpotifyID:  tr.ID.String(),
			Number:     tr.TrackNumber,
			DurationMS: tr.Duration,
			AlbumID:    album.ID,
		}
		tracks = append(tracks, track)
	}
	return tracks
}

func checkAlbumReleasedToday(album db.Album) (releasedToday bool) {
	currentDate := time.Now().Format("2006-01-02")
	if album.ReleaseDate.Format("2006-01-02") == currentDate {
		return true
	} else {
		return false
	}
}

// Fetch - Main function to update database with newly released albums and
// their tracks
// For all artists, get latest releases. If any artist has an album released
// today, add that album and its tracks to the database if it is not already
// present in the database
func Fetch(isReleasedTodayCheck bool) {
	client, ctx := util.ConnectToSpotify()
	artists, err := db.GetAllArtists()
	if err != nil {
		log.Printf("Fetch: failed to get artists: %s\n", err.Error())
		return
	}
	newAlbums := spotifyGetArtistsLatestAlbums(
		client,
		ctx,
		artists,
		spotify.AlbumTypeAlbum,
	)
	newSingles := spotifyGetArtistsLatestAlbums(
		client,
		ctx,
		artists,
		spotify.AlbumTypeSingle,
	)
	newAlbums = append(newAlbums, newSingles...)

	for _, a := range newAlbums {
		if isReleasedTodayCheck {
			if !checkAlbumReleasedToday(a) {
				continue
			}
		}
		albumID, err := db.InsertAlbum(a)
		if err != nil {
			log.Printf("Fetch: failed to insert album %s: %s\n", a.Title, err.Error())
			continue
		}
		if albumID == -1 {
			log.Printf("Album %s already exists.", a.Title)
			continue
		}
		a.ID = albumID
		tracks := spotifyGetAlbumTracks(client, ctx, a)
		for _, track := range tracks {
			if err := db.InsertTrack(track); err != nil {
				log.Printf("Fetch: failed to insert track %s: %s\n", track.Title, err.Error())
			}
		}

	}
}
