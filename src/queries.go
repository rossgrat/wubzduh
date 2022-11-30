package wubzduh

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

type Artist struct {
	ID         int
	ArtistName string
	SpotifyID  string
}

type Track struct {
	ID              int
	TrackTitle      string
	SpotifyID       string
	TrackNumber     int
	DurationMS      int
	DurationMinutes string
	DurationSeconds string
}

type Album struct {
	ID          int
	AlbumTitle  string
	SpotifyID   string
	ArtistID    int
	ArtistName  string
	ReleaseDate string
	CoverartURL string
	AlbumType   string
	AlbumURL    string
	Tracks      []Track
}

//Database queries
const (
	getAlbumsJoinArtistName = `SELECT Albums.ID, Albums.AlbumTitle, Albums.CoverartURL, Albums.ReleaseDate, Albums.AlbumType, Albums.AlbumURL, Artists.ArtistName FROM Albums INNER JOIN Artists ON Albums.ArtistID = Artists.ID ORDER BY ReleaseDate DESC`
	getAlbumTracks          = `SELECT ID, TrackTitle, LengthMS, TrackNumber FROM Tracks WHERE AlbumID=$1 ORDER BY TrackNumber ASC`
	getArtists              = `SELECT * FROM Artists ORDER BY ArtistName ASC;`
	getAlbumIDBySpotifyID   = `SELECT ID FROM Albums WHERE SpotifyID=$1 ORDER BY ID DESC`
	insertAlbum             = `INSERT INTO Albums (AlbumTitle, ArtistID, CoverartURL, ReleaseDate, AlbumType, AlbumURL, SpotifyID) values ($1, $2, $3, $4, $5, $6, $7)`
	insertTrack             = `INSERT INTO Tracks (TrackTitle, AlbumID, LengthMS, TrackNumber, SpotifyID, AddedToPlaylist) values ($1, $2, $3, $4, $5, 'false')`
)

//All Database access must first lock this Reader Writer Mutex
//Database reads must use the read lock
//Database writes must use the write lock (INSERTs)
//Even if different database operators are active, our RWMutexes implementation still works, i.e,
// any operator can read so long as no operators are writing, only one operator may write so long
// as no other operators are reading
var DBLock sync.RWMutex

//Get all albums and the tracks for each album
func GetAlbumsAndTracks(db *sql.DB) (albums []Album) {
	DBLock.RLock()
	//Get all albums and join with artist name
	albumRows, err := db.Query(getAlbumsJoinArtistName)
	if err != nil {
		log.Fatalf("Error - failed to query albums: %v", err)
	}
	for albumRows.Next() {
		var a Album
		var releaseDateTime time.Time
		err := albumRows.Scan(&a.ID, &a.AlbumTitle, &a.CoverartURL, &releaseDateTime, &a.AlbumType, &a.AlbumURL, &a.ArtistName)
		if err != nil {
			log.Fatalf("Error - failed to scan album row: %v", err)
		}
		a.ReleaseDate = releaseDateTime.Format("2006-01-02")

		//Get all album tracks
		trackRows, err := db.Query(getAlbumTracks, a.ID)
		if err != nil {
			log.Fatalf("Error - failed to query album tracks: %v", err)
		}
		for trackRows.Next() {
			var t Track
			var durationMs int
			err := trackRows.Scan(&t.ID, &t.TrackTitle, &durationMs, &t.TrackNumber)
			if err != nil {
				log.Fatalf("Error - failed to scan track row: %v", err)
			}
			t.DurationMinutes = fmt.Sprintf("%0d", (int)(math.Floor((float64)(durationMs)/(float64)(1000)/(float64)(60))))
			t.DurationSeconds = fmt.Sprintf("%02d", (durationMs/1000)%60)
			//Append track to album tracks
			a.Tracks = append(a.Tracks, t)
		}
		//Append album to albums slice
		albums = append(albums, a)
	}
	DBLock.RUnlock()
	return albums
}

//Get all artists in database
func GetAllArtists(db *sql.DB) (artists []Artist) {
	DBLock.RLock()
	rows, err := db.Query(getArtists)
	if err != nil {
		log.Fatalf("Error - failed to query artists: %v", err)
	}

	for rows.Next() {
		var a Artist
		err = rows.Scan(&a.ID, &a.ArtistName, &a.SpotifyID)
		if err != nil {
			log.Fatalf("Error - failed to scan artist row: %v", err)
		}
		artists = append(artists, a)
	}
	DBLock.RUnlock()
	return artists
}

//Get the ID for the album with Spotify ID
//We enforce that only one album with a given spotify ID may be inserted into the database during InsertAlbum
//In the case that we have two albums with the same title and different spotify IDs, this query returns only the album
// with the highest (latest) ID
func GetAlbumIDWithSpotifyID(db *sql.DB, SpotifyID string) (albumID int) {
	DBLock.RLock()
	//Get album ID from db
	rows, err := db.Query(getAlbumIDBySpotifyID, SpotifyID)
	if err != nil {
		log.Fatalf("Error - failed to query album: %v", err)
	}
	rows.Next()
	err = rows.Scan(&albumID)
	if err != nil {
		log.Fatalf("Error - failed to scan row: %v", err)
	}
	DBLock.RUnlock()
	return albumID
}

//Insert an album into the database
//Only insert the album if no other albums with a matching spotify ID exist in the database
func InsertAlbum(db *sql.DB, album Album) (success bool) {
	DBLock.Lock()
	//Check if an album with this albums spotify ID exists in the database
	albumQuery, err := db.Exec(getAlbumIDBySpotifyID, album.SpotifyID)
	if err != nil {
		log.Fatalf("Error - couldn't query Album: %v", err)
	}
	count, err := albumQuery.RowsAffected()
	if err != nil {
		log.Fatalf("Error - couldn't check Album query rows: %v", err)
	}
	if count != 0 {
		log.Printf("Could not insert album, matching Spotify ID already exists in DB.")
		DBLock.Unlock()
		return false
	}
	//Insert the album
	_, err = db.Exec(insertAlbum, album.AlbumTitle, album.ArtistID, album.CoverartURL, album.ReleaseDate, album.AlbumType, album.AlbumURL, album.SpotifyID)
	if err != nil {
		log.Fatalf("Error - couldn't insert Album: %v", err)
	}
	DBLock.Unlock()
	return true
}

//Insert a list of tracks into the database and associate them with an album that is in the database
func InsertTracksWithAlbumID(db *sql.DB, tracks []Track, albumID int) (success bool) {
	DBLock.Lock()
	for _, t := range tracks {
		_, err := db.Exec(insertTrack, t.TrackTitle, albumID, t.DurationMS, t.TrackNumber, t.ID)
		if err != nil {
			log.Fatalf("Error - couldn't insert Track: %v", err)
		}
	}
	DBLock.Unlock()
	return true
}
