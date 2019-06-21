package sqlite

import (
	"database/sql"

	"github.com/jeremybouzigard/library"
)

// AlbumDiscogService manages interactions that link artist to album in the
// data source.
type AlbumDiscogService struct {
	insert  *sql.Stmt
	session *Session
}

// NewAlbumDiscogService returns a new instance of an AlbumDiscogService that
// operates within the given session.
func NewAlbumDiscogService(s *Session) AlbumDiscogService {
	service := AlbumDiscogService{session: s}
	return service
}

// CreateTable creates the 'album_discographies' table and returns any errors.
func (service *AlbumDiscogService) CreateTable() (sql.Result, error) {
	create :=
		`CREATE TABLE IF NOT EXISTS album_discographies (
			artist_id INTEGER NOT NULL,
			album_id  INTEGER NOT NULL,
			PRIMARY KEY('artist_id','album_id'),
			FOREIGN KEY('artist_id') REFERENCES artists('artist_id'),
			FOREIGN KEY('album_id')  REFERENCES albums('album_id')
		)`
	return service.session.tx.Exec(create)
}

// DropTable drops the 'album_discographies' table and returns any errors.
func (service *AlbumDiscogService) DropTable() (sql.Result, error) {
	drop := `DROP TABLE IF EXISTS album_discographies`
	return service.session.tx.Exec(drop)
}

// CreateAlbumDiscog inserts a new record.
func (service *AlbumDiscogService) CreateAlbumDiscog(attributes *library.AlbumAttributes) error {
	if service.insert == nil {
		stmt, err := service.prepareInsert()
		if err != nil {
			service.session.Logger.Println(err)
			return err
		}
		service.insert = stmt
	}

	_, err := service.insert.Exec(
		attributes.ArtistName, attributes.ArtistSort,
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

// prepareInsert creates a prepared statement to insert a new record.
func (service *AlbumDiscogService) prepareInsert() (*sql.Stmt, error) {
	insert :=
		`INSERT OR IGNORE INTO album_discographies 
		                       (artist_id, 
		                        album_id) 
		                SELECT (SELECT artist_id 
		                         FROM artists 
		                        WHERE artist_name = ? 
		                          AND artist_sort = ?), 
		                       (SELECT album_id 
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

// Close closes all open statements.
func (service *AlbumDiscogService) Close() error {
	if service.insert != nil {
		err := service.insert.Close()
		if err != nil {
			service.session.Logger.Println(err)
			return err
		}
	}

	return nil
}
