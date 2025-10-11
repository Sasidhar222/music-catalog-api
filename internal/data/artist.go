package data

import (
	"database/sql"
	"errors"

	"github.com/Sasidhar222/music-api/models"
	"github.com/google/uuid"
)

type ArtistModel struct {
	DB *sql.DB
}

var ErrNotFound = errors.New("NOT FOUND")

func (m ArtistModel) GetAll() ([]*models.Artist, error) {
	// Create a slice to hold the results
	var artists []*models.Artist

	rows, err := m.DB.Query("SELECT id, name FROM artists")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows
	for rows.Next() {
		var artist models.Artist
		// Scan the row data into the artist struct fields
		if err := rows.Scan(&artist.ID, &artist.Name); err != nil {
			return nil, err
		}
		artists = append(artists, &artist)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return artists, nil
}

func (m ArtistModel) AddArtist(newArtist *models.Artist) (*models.Artist, error) {
	// 1. Generate a new unique ID
	newArtist.ID = uuid.New().String()
	// 2. Prepare the SQL INSERT statement with placeholders to prevent SQL injection
	sqlStatement := `INSERT INTO artists (id, name) VALUES ($1, $2)`
	_, err := m.DB.Exec(sqlStatement, newArtist.ID, newArtist.Name)

	if err != nil {
		return nil, err
	}
	return newArtist, nil
}

func (m ArtistModel) UpdateArtist(id string, artist *models.Artist) (*models.Artist, error) {
	sqlStatement := `UPDATE artists set name = $1 where id = $2`
	res, err := m.DB.Exec(sqlStatement, artist.Name, id)

	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrNotFound
	}
	artist.ID = id
	return artist, nil
}

func (m ArtistModel) DeleteArtist(id string) (bool, error) {
	sqlStatement := `DELETE FROM artists where id = $1`
	res, err := m.DB.Exec(sqlStatement, id)

	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, ErrNotFound
	}

	return true, nil
}
