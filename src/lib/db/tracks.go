package db

import (
	"errors"
)

func InsertTrack(track Track) error {
	fn := "InsertTrack"
	_, err := DB.Exec(`
		INSERT INTO tracks (
			title, 
			album_id, 
			length_ms, 
			number, 
			spotify_id
		) VALUES ($1, $2, $3, $4, $5)`,
		track.Title,
		track.AlbumID,
		track.DurationMS,
		track.Number,
		track.SpotifyID,
	)
	if err != nil {
		return errors.New(fn + ": failed to query tracks - " + err.Error())
	}
	return nil
}

func DeleteTracksForAlbum(albumID int) error {
	fn := "DeleteTracksForAlbum"
	if _, err := DB.Exec(`
		DELETE FROM tracks 
			* 
		WHERE album_id=$1`,
		albumID,
	); err != nil {
		return errors.New(fn + ": failed to delete tracks - " + err.Error())
	}
	return nil
}
