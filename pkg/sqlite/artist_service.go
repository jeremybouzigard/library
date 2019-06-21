package sqlite

import (
	"bytes"
	"database/sql"

	"github.com/jeremybouzigard/library"
)

// ArtistService manages interactions with the artist data source.
type ArtistService struct {
	session *Session
	insert  *sql.Stmt
}

// NewArtistService returns a new instance of an ArtistService that operates
// within the given session.
func NewArtistService(s *Session) ArtistService {
	service := ArtistService{session: s}
	return service
}

// CreateTable creates the 'artists' table and returns any errors.
func (service *ArtistService) CreateTable() (sql.Result, error) {
	create :=
		`CREATE TABLE IF NOT EXISTS artists (
			artist_id   INTEGER PRIMARY KEY,
			artist_name TEXT    NOT NULL,
			artist_sort TEXT
		)`
	return service.session.tx.Exec(create)
}

// DropTable drops the 'artists' table and returns any errors.
func (service *ArtistService) DropTable() (sql.Result, error) {
	drop := `DROP TABLE IF EXISTS artists`
	return service.session.tx.Exec(drop)
}

// CreateArtist inserts a new artist.
func (service *ArtistService) CreateArtist(attributes *library.ArtistAttributes) error {
	if service.insert == nil {
		stmt, err := service.prepareInsert()
		if err != nil {
			service.session.Logger.Println(err)
			return err
		}
		service.insert = stmt
	}

	_, err := service.insert.Exec(
		attributes.Name, attributes.Sort,
		attributes.Name, attributes.Sort)

	if err != nil {
		service.session.Logger.Println(err)
		return err
	}
	return nil
}

// prepareInsert creates a prepared statement to insert a new artist.
func (service *ArtistService) prepareInsert() (*sql.Stmt, error) {
	insert :=
		`     INSERT INTO artists 
		                  (artist_name, 
		                   artist_sort) 
		           SELECT ?, 
		                  ? 
		 WHERE NOT EXISTS (SELECT 1 
		                    FROM artists 
		                   WHERE artist_name = ? 
		                     AND artist_sort = ?)`
	return service.session.tx.Prepare(insert)
}

// Artist queries the 'artists' table for an artist with the given ID and
// returns the result along with any error.
func (service *ArtistService) Artist(ID string) (*library.Artist, error) {
	var a library.Artist
	query :=
		`SELECT
			artist_id,
			artist_name,
			artist_sort
		FROM
			artists
		WHERE 
			artist_id = ?`
	err := service.session.db.QueryRow(query, ID).Scan(
		&a.ID,
		&a.Attributes.Name,
		&a.Attributes.Sort)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		service.session.Logger.Println(err)
		return &a, err
	}
	return &a, nil
}

// Artists queries the 'artists' table for all artists that meet the given
// criteria and returns the result along with any error.
func (service *ArtistService) Artists(queries map[string]string) ([]*library.Artist, error) {
	var results []*library.Artist

	rows, err := service.Query(queries)
	if err != nil {
		service.session.Logger.Println(err)
		return results, err
	}

	for rows.Next() {
		var a library.Artist
		err := rows.Scan(
			&a.ID,
			&a.Attributes.Name,
			&a.Attributes.Sort)
		a.Type = "artists"
		if err != nil {
			service.session.Logger.Println(err)
			return results, err
		}
		results = append(results, &a)
	}

	if len(results) < 1 {
		return nil, nil
	}

	return results, nil
}

// Query executes a query for artists that meet the given predicate criteria and
// returns the results along with any error.
func (service *ArtistService) Query(predicates map[string]string) (*sql.Rows, error) {
	query := bytes.NewBufferString(
		`SELECT
		  artist_id,
		  artist_name,
		  artist_sort
		FROM
		  artists`)
	query, args := Where(query, predicates)
	return service.session.db.Query(query.String(), args...)
}

// Close closes all open statements.
func (service *ArtistService) Close() error {
	if service.insert != nil {
		err := service.insert.Close()
		if err != nil {
			service.session.Logger.Println(err)
			return err
		}
	}

	return nil
}
