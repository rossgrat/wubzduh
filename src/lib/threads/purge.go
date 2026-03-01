package threads

import (
	"log"
	"time"

	"github.com/rossgrat/wubzduh/src/lib/db"
)

// Delete all albums and their tracks from the database if the album was released more than a week ago
func Purge() {
	currentDate := time.Now()
	currentDate = currentDate.AddDate(0, 0, -14)

	albums, err := db.GetAllAlbumReleaseDates()
	if err != nil {
		log.Printf("Purge: failed to get album release dates: %s\n", err.Error())
		return
	}
	for _, album := range albums {
		if album.ReleaseDate.Before(currentDate) {
			if err := db.DeleteTracksForAlbum(album.ID); err != nil {
				log.Printf("Purge: failed to delete tracks for album %d: %s\n", album.ID, err.Error())
				continue
			}
			if err := db.DeleteAlbum(album.ID); err != nil {
				log.Printf("Purge: failed to delete album %d: %s\n", album.ID, err.Error())
			}
		}
	}
}
