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
		log.Fatal(err.Error())
	}
	for _, album := range albums {
		if album.ReleaseDate.Before(currentDate) {
			if err := db.DeleteTracksForAlbum(album.ID); err != nil {
				log.Fatal(err.Error())
			}
			if err := db.DeleteAlbum(album.ID); err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}
