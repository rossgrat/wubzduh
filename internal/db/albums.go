package db

import (
	"database/sql"
	"fmt"
	"math"
)

func (s *Store) GetAlbums() ([]Album, error) {
	albumRows, err := s.db.Query(`
		SELECT
			albums.id,
			albums.title,
			albums.cover_art_url,
			albums.release_date,
			albums.type,
			albums.url,
			artists.name
		FROM albums
		INNER JOIN artists ON albums.artist_id = artists.id
		ORDER BY release_date DESC`)
	if err != nil {
		return nil, fmt.Errorf("querying albums: %w", err)
	}
	defer albumRows.Close()

	var albums []Album
	for albumRows.Next() {
		var a Album
		if err := albumRows.Scan(
			&a.ID, &a.Title, &a.CoverartURL,
			&a.ReleaseDate, &a.Type, &a.URL, &a.ArtistName,
		); err != nil {
			return nil, fmt.Errorf("scanning album row: %w", err)
		}

		trackRows, err := s.db.Query(`
			SELECT id, title, length_ms, number
			FROM tracks
			WHERE album_id = $1
			ORDER BY number ASC`,
			a.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("querying tracks for album %d: %w", a.ID, err)
		}

		for trackRows.Next() {
			var t Track
			var durationMs int
			if err := trackRows.Scan(&t.ID, &t.Title, &durationMs, &t.Number); err != nil {
				trackRows.Close()
				return nil, fmt.Errorf("scanning track row: %w", err)
			}
			t.DurationMinutes = fmt.Sprintf("%d", int(math.Floor(float64(durationMs)/1000/60)))
			t.DurationSeconds = fmt.Sprintf("%02d", (durationMs/1000)%60)
			a.Tracks = append(a.Tracks, t)
		}
		trackRows.Close()

		albums = append(albums, a)
	}
	return albums, nil
}

func (s *Store) InsertAlbum(album Album) (int, error) {
	var albumID int
	err := s.db.QueryRow(`
		INSERT INTO albums (title, artist_id, cover_art_url, release_date, type, url, spotify_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT DO NOTHING
		RETURNING id`,
		album.Title, album.ArtistID, album.CoverartURL,
		album.ReleaseDate, album.Type, album.URL, album.SpotifyID,
	).Scan(&albumID)

	if err == sql.ErrNoRows {
		return -1, nil
	}
	if err != nil {
		return -1, fmt.Errorf("inserting album: %w", err)
	}
	return albumID, nil
}

func (s *Store) DeleteAlbum(albumID int) error {
	_, err := s.db.Exec(`DELETE FROM albums WHERE id = $1`, albumID)
	if err != nil {
		return fmt.Errorf("deleting album %d: %w", albumID, err)
	}
	return nil
}

func (s *Store) GetAllAlbumReleaseDates() ([]Album, error) {
	rows, err := s.db.Query(`
		SELECT id, release_date
		FROM albums
		ORDER BY release_date ASC`)
	if err != nil {
		return nil, fmt.Errorf("querying album release dates: %w", err)
	}
	defer rows.Close()

	var albums []Album
	for rows.Next() {
		var a Album
		if err := rows.Scan(&a.ID, &a.ReleaseDate); err != nil {
			return nil, fmt.Errorf("scanning album release date row: %w", err)
		}
		albums = append(albums, a)
	}
	return albums, nil
}
