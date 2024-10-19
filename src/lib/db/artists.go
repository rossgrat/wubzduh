package db

import (
	"errors"
)

func GetAllArtists() ([]Artist, error) {
	fn := "GetAllArtists"
	artists := []Artist{}
	rows, err := DB.Query(`
		SELECT 
			id,
			name,
			spotify_id
		FROM artists 
		ORDER BY name ASC`)
	if err != nil {
		return artists, errors.New(fn + " - failed to query artists" + err.Error())
	}
	for rows.Next() {
		var a Artist
		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.SpotifyID,
		); err != nil {
			return artists, errors.New(fn + " - failed to scan artist row" + err.Error())
		}
		artists = append(artists, a)
	}
	return artists, nil
}

func InsertArtist(artist Artist) error {
	fn := "InsertArtist"
	if _, err := DB.Exec(`
		INSERT INTO artists (
			artist_name,
			spotify_id
		)
		VALUES ($1, $2)`,
		artist.Name,
		artist.SpotifyID,
	); err != nil {
		return errors.New(fn + " - failed to insert artist" + err.Error())
	}
	return nil
}

func GetAllArtistsWithName(name string) ([]Artist, error) {
	fn := "GetAllArtists"
	artists := []Artist{}
	rows, err := DB.Query(`
		SELECT 
			id,
			name,
			spotify_id
		FROM artists 
		WHERE name=$1
		ORDER BY name ASC`,
		name)
	if err != nil {
		return artists, errors.New(fn + " - failed to query artists" + err.Error())
	}
	for rows.Next() {
		var a Artist
		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.SpotifyID,
		); err != nil {
			return artists, errors.New(fn + " - failed to scan artist row" + err.Error())
		}
		artists = append(artists, a)
	}
	return artists, nil
}
