package wubzduh

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/zmb3/spotify/v2"
)

//For each artist in the database, get the latest release by each artist, may be either a full length album or an EP / Single, we make the following assumptions here:
//1) An arist will not release both a new single and a new album in one day
//2) An artist will not release more than one new single or more than one new album in one day
func spotifyGetArtistsLatestAlbums(client *spotify.Client, ctx context.Context, artists []Artist, albumType spotify.AlbumType) (albums []Album) {
	//For all artists, check if latest album released has release date equal to todays date
	for _, a := range artists {

		fmt.Printf("%s %s\n", a.ArtistName, a.SpotifyID)
		albumResults, err := client.GetArtistAlbums(ctx, (spotify.ID)(a.SpotifyID), []spotify.AlbumType{albumType}, spotify.Limit(1))
		if err != nil {
			log.Fatal(err)
		}

		//Create new album from album retrieved from spotify API
		newAlbum := Album{
			AlbumTitle:  albumResults.Albums[0].Name,
			ArtistID:    a.ID,
			ReleaseDate: albumResults.Albums[0].ReleaseDate,
			CoverartURL: albumResults.Albums[0].Images[1].URL,
			AlbumType:   albumResults.Albums[0].AlbumType,
			AlbumURL:    (string)(albumResults.Albums[0].ExternalURLs["spotify"]),
			SpotifyID:   (string)(albumResults.Albums[0].ID),
		}

		albums = append(albums, newAlbum)
	}
	return albums
}

func spotifyGetAlbumTracks(client *spotify.Client, ctx context.Context, album Album) (tracks []Track) {
	trackResults, err := client.GetAlbumTracks(ctx, (spotify.ID)(album.SpotifyID))
	if err != nil {
		log.Fatal(err)
	}
	for _, tr := range trackResults.Tracks {
		track := Track{
			TrackTitle:  tr.Name,
			SpotifyID:   (string)(tr.ID),
			TrackNumber: tr.TrackNumber,
			DurationMS:  tr.Duration,
		}
		tracks = append(tracks, track)
	}
	return tracks
}

func checkAlbumReleasedToday(album Album) (releasedToday bool) {
	currentDate := time.Now().Format("2006-01-02")
	if album.ReleaseDate == currentDate {
		return true
	} else {
		return false
	}
}

//Fetch - Main function to update database with newly released albums and their tracks
//For all artists, get latest releases. If any artist has an album released today, add that album and its tracks to the database
func Fetch(db *sql.DB, client *spotify.Client, ctx context.Context) {
	//Get slice of latest albums from all artists
	artists := GetAllArtists(db)
	newAlbums := spotifyGetArtistsLatestAlbums(client, ctx, artists, spotify.AlbumTypeAlbum)
	newSingles := spotifyGetArtistsLatestAlbums(client, ctx, artists, spotify.AlbumTypeSingle)

	newAlbums = append(newAlbums, newSingles...)

	for _, a := range newAlbums {
		if checkAlbumReleasedToday(a) {
			//Create new album entry in database, if we could not insert the album (album already exists in db), continue
			if InsertAlbum(db, a) == false {
				continue
			}
			//Album does not already exist in database
			//Query spotify for all tracks on album
			tracks := spotifyGetAlbumTracks(client, ctx, a)
			//Get ID of album created in database
			a.ID = GetAlbumIDWithSpotifyID(db, a.SpotifyID)
			//Insert all album tracks into database, associate with album ID
			InsertTracksWithAlbumID(db, tracks, a.ID)
		}
	}
}

//Fetches all latest albums and their tracks and inserts them into databases. Used for populating the database for testing
func FetchAllLatest(db *sql.DB, client *spotify.Client, ctx context.Context) {
	//Get slice of latest albums from all artists
	artists := GetAllArtists(db)
	newAlbums := spotifyGetArtistsLatestAlbums(client, ctx, artists, spotify.AlbumTypeAlbum)
	newSingles := spotifyGetArtistsLatestAlbums(client, ctx, artists, spotify.AlbumTypeSingle)

	newAlbums = append(newAlbums, newSingles...)

	for _, a := range newAlbums {
		//Create new album entry in database, if we could not insert the album (album already exists in db), continue
		if InsertAlbum(db, a) == false {
			continue
		}
		//Album does not already exist in database
		//Query spotify for all tracks on album
		tracks := spotifyGetAlbumTracks(client, ctx, a)
		//Get ID of album created in database
		a.ID = GetAlbumIDWithSpotifyID(db, a.SpotifyID)
		//Insert all album tracks into database, associate with album ID
		InsertTracksWithAlbumID(db, tracks, a.ID)
	}
}
