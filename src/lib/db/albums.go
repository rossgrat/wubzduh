package db

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
)

func GetAlbums() ([]Album, error) {
	fn := "GetAlbums"
	albumRows, err := DB.Query(`
		SELECT 
			albums.id, 
			albums.title,
			albums.cover_art_url, 
			albums.release_date, 
			albums.type,
			albums.url, 
			artists.name
		FROM albums
		INNER JOIN artists 
			ON albums.artist_id=artists.id
		ORDER BY release_date DESC`)
	if err != nil {
		return []Album{}, errors.New(fn + ": failed to query albums - " + err.Error())
	}
	albums := []Album{}
	for albumRows.Next() {
		var a Album
		if err := albumRows.Scan(
			&a.ID,
			&a.Title,
			&a.CoverartURL,
			&a.ReleaseDate,
			&a.Type,
			&a.URL,
			&a.ArtistName,
		); err != nil {
			return albums, errors.New(fn + ": failed to scan album rows  - " + err.Error())
		}
		trackRows, err := DB.Query(`
			SELECT 
				id, 
				title, 
				length_ms, 
				number 
			FROM tracks 
			WHERE album_id=$1 
			ORDER BY number ASC`,
			a.ID)
		if err != nil {
			return albums, errors.New(fn + ": failed to query album tracks  - " + err.Error())
		}
		for trackRows.Next() {
			var t Track
			var durationMs int
			if err := trackRows.Scan(
				&t.ID,
				&t.Title,
				&durationMs,
				&t.Number,
			); err != nil {
				return albums, errors.New(fn + ": failed to query scan track rows  - " + err.Error())
			}
			t.DurationMinutes = fmt.Sprintf("%0d", (int)(math.Floor((float64)(durationMs)/(float64)(1000)/(float64)(60))))
			t.DurationSeconds = fmt.Sprintf("%02d", (durationMs/1000)%60)
			a.Tracks = append(a.Tracks, t)
		}
		albums = append(albums, a)
	}
	return albums, nil
}

// Only insert the album if no other albums with a matching spotify ID exist in the database
func InsertAlbum(album Album) (int, error) {
	fn := "InsertAlbum"
	albumID := -1
	row := DB.QueryRow(`
		INSERT INTO albums (
			title, 
			artist_id, 
			cover_art_url, 
			release_date, 
			type, 
			url, 
			spotify_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT
			DO NOTHING
		RETURNING id`,
		album.Title,
		album.ArtistID,
		album.CoverartURL,
		album.ReleaseDate,
		album.Type,
		album.URL,
		album.SpotifyID,
	)
	if err := row.Scan(
		&albumID,
	); err != nil && err != sql.ErrNoRows {
		return albumID, errors.New(fn + ": couldn't scan album ID  - " + err.Error())
	}
	return albumID, nil
}

func DeleteAlbum(albumID int) error {
	fn := "DeleteAlbum"
	if _, err := DB.Exec(`
		DELETE FROM albums 
			* 
		WHERE id=$1`,
		albumID,
	); err != nil {
		return errors.New(fn + ": failed to delete albums  - " + err.Error())
	}
	return nil
}
func GetAllAlbumReleaseDates() ([]Album, error) {
	fn := "GetAllAlbumReleaseDates"
	albums := []Album{}
	albumRows, err := DB.Query(`
	SELECT 
		id, 
		release_date 
	FROM albums 
	ORDER BY release_date ASC`)
	if err != nil {
		return albums, errors.New(fn + ": failed to query albums  - " + err.Error())
	}
	for albumRows.Next() {
		var album Album
		if err := albumRows.Scan(
			&album.ID,
			&album.ReleaseDate,
		); err != nil {
			return albums, errors.New(fn + ": failed to scan album row  - " + err.Error())
		}
		albums = append(albums, album)
	}
	return albums, nil
}
