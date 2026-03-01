package db

import "fmt"

func (s *Store) InsertTrack(track Track) error {
	_, err := s.db.Exec(`
		INSERT INTO tracks (title, album_id, length_ms, number, spotify_id)
		VALUES ($1, $2, $3, $4, $5)`,
		track.Title, track.AlbumID, track.DurationMS, track.Number, track.SpotifyID,
	)
	if err != nil {
		return fmt.Errorf("inserting track: %w", err)
	}
	return nil
}

func (s *Store) DeleteTracksForAlbum(albumID int) error {
	_, err := s.db.Exec(`DELETE FROM tracks WHERE album_id = $1`, albumID)
	if err != nil {
		return fmt.Errorf("deleting tracks for album %d: %w", albumID, err)
	}
	return nil
}
