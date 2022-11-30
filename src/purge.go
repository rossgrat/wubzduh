package wubzduh

import (
	"database/sql"
	"log"
	"time"
)

const (
	getAllAlbums               = `SELECT ID, ReleaseDate FROM Albums ORDER BY ReleaseDate ASC`
	deleteAllTracksWithAlbumID = `DELETE FROM Tracks * WHERE AlbumID=$1`
	deleteAllAlbumsWithID      = `DELETE FROM Albums * WHERE ID=$1`
)

//Delete all albums and their tracks from the database if the album was released more than a week ago
func Purge(db *sql.DB) {
	DBLock.Lock()

	//Get current date and subtract 7 days from it
	currentDate := time.Now()
	currentDate = currentDate.AddDate(0, 0, -14)

	albumRows, err := db.Query(getAllAlbums)
	if err != nil {
		log.Fatalf("Failed to query albums: %v", err)
	}

	var id int
	var date time.Time
	for albumRows.Next() {
		err := albumRows.Scan(&id, &date)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		if date.Before(currentDate) {
			_, err := db.Exec(deleteAllTracksWithAlbumID, id)
			if err != nil {
				log.Fatalf("Failed to delete tracks: %v", err)
			}
			_, err = db.Exec(deleteAllAlbumsWithID, id)
			if err != nil {
				log.Fatalf("Failed to delete albums: %v", err)
			}
		}
	}
	DBLock.Unlock()
}
