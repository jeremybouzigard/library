package sqlite

import (
	"database/sql"

	"github.com/jeremybouzigard/library"
)

// SongDiscogService manages interactions that link artist to song in the
// data source.
type SongDiscogService struct {
	session *Session
	insert  *sql.Stmt
}

// NewSongDiscogService returns a new instance of a SongDiscogService that
// operates within the given session.
func NewSongDiscogService(s *Session) SongDiscogService {
	sds := SongDiscogService{session: s}
	return sds
}

// CreateTable creates the 'song_discographies' table and returns any errors.
func (sds *SongDiscogService) CreateTable() (sql.Result, error) {
	create :=
		`CREATE TABLE IF NOT EXISTS song_discographies (
			artist_id INTEGER NOT NULL,
			song_id   INTEGER NOT NULL,
			album_id  INTEGER,
			PRIMARY KEY('artist_id','song_id'),
			FOREIGN KEY('artist_id') REFERENCES artists('artist_id'),
			FOREIGN KEY('song_id')   REFERENCES songs('song_id'),
			FOREIGN KEY('album_id')  REFERENCES albums('album_id')
		)`
	return sds.session.tx.Exec(create)
}

// DropTable drops the 'song_discographies' table and returns any errors.
func (sds *SongDiscogService) DropTable() (sql.Result, error) {
	drop := `DROP TABLE IF EXISTS song_discographies`
	return sds.session.tx.Exec(drop)
}

// PrepareInsert creates a prepared statement to insert a new record.
func (sds *SongDiscogService) PrepareInsert() (*sql.Stmt, error) {
	insert :=
		`INSERT OR IGNORE INTO song_discographies 
		                       (artist_id, 
								song_id,
								album_id) 
		                SELECT (SELECT artist_id 
		                          FROM artists 
		                         WHERE artist_name = ? 
		                           AND artist_sort = ?), 
		                       (SELECT song_id 
		                          FROM songs 
								 WHERE file_path = ?),
		                       (SELECT album_id 
								 FROM  albums 
								WHERE  album_name = ? 
								  AND  album_sort = ?)`
	return sds.session.tx.Prepare(insert)
}

// CreateSongDiscog inserts a new record.
func (sds *SongDiscogService) CreateSongDiscog(sa *library.SongAttributes, aa *library.AlbumAttributes) error {
	if sds.insert == nil {
		stmt, err := sds.PrepareInsert()
		if err != nil {
			sds.session.Logger.Println(err)
			return err
		}
		sds.insert = stmt
	}

	_, err := sds.insert.Exec(
		sa.ArtistName, sa.ArtistSort,
		sa.FilePath,
		aa.Name, aa.Sort)

	if err != nil {
		sds.session.Logger.Println(err)
		return err
	}
	return nil
}

// Close closes all open statements.
func (sds *SongDiscogService) Close() error {
	if sds.insert != nil {
		err := sds.insert.Close()
		if err != nil {
			sds.session.Logger.Println(err)
			return err
		}
	}

	return nil
}
