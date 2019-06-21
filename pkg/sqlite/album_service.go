package sqlite

import (
	"bytes"
	"database/sql"

	"github.com/jeremybouzigard/library"
)

// AlbumService manages interactions with the album data source.
type AlbumService struct {
	session *Session
	insert  *sql.Stmt
}

// NewAlbumService returns a new instance of an AlbumService that operates
// within the given session.
func NewAlbumService(s *Session) AlbumService {
	service := AlbumService{session: s}
	return service
}

// CreateTable creates the 'albums' table and returns any errors.
func (service *AlbumService) CreateTable() (sql.Result, error) {
	create :=
		`CREATE TABLE IF NOT EXISTS albums (
			album_id          INTEGER PRIMARY KEY,
			album_name        TEXT    NOT NULL,
			artist_id         INTEGER NOT NULL,
			genre_id          INTEGER,
			release_date      TEXT,
			track_total       INTEGER,
			album_sort        TEXT,
			album_artist      TEXT,
			album_artist_sort TEXT,
			FOREIGN KEY('artist_id') REFERENCES artists('artist_id'),
			FOREIGN KEY('genre_id')  REFERENCES genres('genre_id')
		)`
	return service.session.tx.Exec(create)
}

// DropTable drops the 'albums' table and returns any errors.
func (service *AlbumService) DropTable() (sql.Result, error) {
	drop := `DROP TABLE IF EXISTS albums`
	return service.session.tx.Exec(drop)
}

// CreateAlbum inserts a new album.
func (service *AlbumService) CreateAlbum(attributes *library.AlbumAttributes) error {
	if service.insert == nil {
		stmt, err := service.prepareInsert()
		if err != nil {
			service.session.Logger.Println(err)
			return err
		}
		service.insert = stmt
	}

	_, err := service.insert.Exec(
		attributes.Name,
		attributes.ArtistName, attributes.ArtistSort,
		attributes.GenreName,
		attributes.ReleaseDate,
		attributes.Sort,
		attributes.AlbumArtist,
		attributes.AlbumArtistSort,
		attributes.Name,
		attributes.Sort,
		attributes.ReleaseDate,
		attributes.ArtistName, attributes.ArtistSort,
		attributes.GenreName)

	if err != nil {
		service.session.Logger.Println(err)
		return err
	}
	return nil
}

// prepareInsert creates a prepared statement to insert a new album.
func (service *AlbumService) prepareInsert() (*sql.Stmt, error) {
	insert :=
		`INSERT INTO albums 
		             (album_name, 
		              artist_id, 
		              genre_id, 
		              release_date, 
		              album_sort, 
		              album_artist, 
		              album_artist_sort)
		                      SELECT ?, 
		                             (SELECT artist_id 
		                                FROM artists 
		                               WHERE artist_name = ? 
		                                 AND artist_sort = ?), 
		                             (SELECT genre_id 
		                                FROM genres 
		                               WHERE genre_name = ?), 
		                             ?, 
		                             ?, 
		                             ?, 
		                             ? 
		            WHERE NOT EXISTS (SELECT 1 
		                                FROM albums 
		                               WHERE album_name = ? 
		                                 AND album_sort = ? 
		                                 AND release_date = ? 
		                                 AND artist_id = (SELECT artist_id 
		                                                    FROM artists 
		                                                   WHERE artist_name = ? 
		                                                     AND artist_sort = ?) 
		                                 AND genre_id = (SELECT genre_id 
		                                                   FROM genres 
		                                                  WHERE genre_name = ?))`
	return service.session.tx.Prepare(insert)
}

// Album queries the 'albums' table for an album with the given ID and returns
// the result along with any error.
func (service *AlbumService) Album(ID string) (*library.Album, error) {
	var a library.Album
	query :=
		`SELECT
			albums.album_id,
			albums.album_name,
			albums.album_sort,
			artists.artist_name,
			artists.artist_sort,
			genres.genre_name,
			albums.release_date
		FROM
			album_discographies
			INNER JOIN artists ON album_discographies.artist_id = artists.artist_id
			INNER JOIN albums ON album_discographies.album_id = albums.album_id
			INNER JOIN genres ON albums.genre_id = genres.genre_id
		WHERE 
			albums.album_id = ?`
	err := service.session.db.QueryRow(query, ID).Scan(
		&a.ID,
		&a.Attributes.Name,
		&a.Attributes.Sort,
		&a.Attributes.ArtistName,
		&a.Attributes.ArtistSort,
		&a.Attributes.GenreName,
		&a.Attributes.ReleaseDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return &a, nil
		}
		service.session.Logger.Println(err)
		return &a, err
	}
	return &a, nil
}

// Albums queries the 'albums' table for all albums that meet the given criteria
// and returns the result along with any error.
func (service *AlbumService) Albums(queries map[string]string) ([]*library.Album, error) {
	var results []*library.Album

	rows, err := service.Query(queries)
	if err != nil {
		service.session.Logger.Println(err)
		return results, err
	}

	for rows.Next() {
		var a library.Album
		err := rows.Scan(
			&a.ID,
			&a.Attributes.Name,
			&a.Attributes.Sort,
			&a.Attributes.ArtistName,
			&a.Attributes.ArtistSort,
			&a.Attributes.GenreName,
			&a.Attributes.ReleaseDate)
		a.Type = "albums"
		if err != nil {
			service.session.Logger.Println(err)
			return results, err
		}
		results = append(results, &a)
	}

	return results, nil
}

// Query executes a query for albums that meet the given predicate criteria and
// returns the results along with any error.
func (service *AlbumService) Query(predicates map[string]string) (*sql.Rows, error) {
	query := bytes.NewBufferString(
		`SELECT
		  albums.album_id,
		  albums.album_name,
		  albums.album_sort,
		  artists.artist_name,
		  artists.artist_sort,
		  genres.genre_name,
		  albums.release_date
		FROM
		  album_discographies
		  INNER JOIN artists ON album_discographies.artist_id = artists.artist_id
		  INNER JOIN albums ON album_discographies.album_id = albums.album_id
		  INNER JOIN genres ON albums.genre_id = genres.genre_id`)
	query, args := Where(query, predicates)
	return service.session.db.Query(query.String(), args...)
}

// Close closes all open statements.
func (service *AlbumService) Close() error {
	if service.insert != nil {
		err := service.insert.Close()
		if err != nil {
			service.session.Logger.Println(err)
			return err
		}
	}
	return nil
}
