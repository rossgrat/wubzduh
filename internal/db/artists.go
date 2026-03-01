package db

import "fmt"

func (s *Store) GetAllArtists() ([]Artist, error) {
	rows, err := s.db.Query(`
		SELECT id, name, spotify_id
		FROM artists
		ORDER BY name ASC`)
	if err != nil {
		return nil, fmt.Errorf("querying artists: %w", err)
	}
	defer rows.Close()

	var artists []Artist
	for rows.Next() {
		var a Artist
		if err := rows.Scan(&a.ID, &a.Name, &a.SpotifyID); err != nil {
			return nil, fmt.Errorf("scanning artist row: %w", err)
		}
		artists = append(artists, a)
	}
	return artists, nil
}

func (s *Store) InsertArtist(artist Artist) error {
	_, err := s.db.Exec(`
		INSERT INTO artists (name, spotify_id)
		VALUES ($1, $2)`,
		artist.Name, artist.SpotifyID,
	)
	if err != nil {
		return fmt.Errorf("inserting artist: %w", err)
	}
	return nil
}

func (s *Store) GetArtistsByName(name string) ([]Artist, error) {
	rows, err := s.db.Query(`
		SELECT id, name, spotify_id
		FROM artists
		WHERE name = $1
		ORDER BY name ASC`,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("querying artists by name: %w", err)
	}
	defer rows.Close()

	var artists []Artist
	for rows.Next() {
		var a Artist
		if err := rows.Scan(&a.ID, &a.Name, &a.SpotifyID); err != nil {
			return nil, fmt.Errorf("scanning artist row: %w", err)
		}
		artists = append(artists, a)
	}
	return artists, nil
}
