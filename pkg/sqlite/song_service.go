package sqlite

import (
	"bytes"
	"database/sql"

	"github.com/jeremybouzigard/library"
)

// SongService manages interactions with the song data source.
type SongService struct {
	session *Session
	insert  *sql.Stmt
}

// NewSongService returns a new instance of a SongService that operates within
// the given session.
func NewSongService(s *Session) SongService {
	ss := SongService{session: s}
	return ss
}

// CreateTable creates the 'songs' table and returns any errors.
func (ss *SongService) CreateTable() (sql.Result, error) {
	create :=
		`CREATE TABLE IF NOT EXISTS songs (
			song_id            INTEGER PRIMARY KEY,
			file_path          TEXT    NOT NULL,
			file_base          TEXT    NOT NULL,
			file_dir           TEXT    NOT NULL,
			artist_id          INTEGER NOT NULL,
			song_name          TEXT,
			genre_id           INTEGER,
			release_date       TEXT,
			track_number       INTEGER,
			disc_number        INTEGER,
			duration_in_millis INTEGER,
			artist_sort        TEXT,
			composer_name      TEXT,
			composer_sort      TEXT,
			conductor          TEXT,
			song_name_sort     TEXT,
			lyrics             TEXT,
			FOREIGN KEY('artist_id') REFERENCES artists('artist_id'),
			FOREIGN KEY('genre_id')  REFERENCES genres('genre_id')
		)`
	return ss.session.tx.Exec(create)
}

// DropTable drops the 'songs' table and returns any errors.
func (ss *SongService) DropTable() (sql.Result, error) {
	drop := `DROP TABLE IF EXISTS songs`
	return ss.session.tx.Exec(drop)
}

// CreateSong inserts a new song.
func (ss *SongService) CreateSong(sa *library.SongAttributes) error {
	if ss.insert == nil {
		stmt, err := ss.PrepareInsert()
		if err != nil {
			ss.session.Logger.Println(err)
			return err
		}
		ss.insert = stmt
	}

	_, err := ss.insert.Exec(
		sa.FilePath,
		sa.FileBase,
		sa.FileDir,
		sa.ArtistName, sa.ArtistSort,
		sa.Name,
		sa.GenreName,
		sa.ReleaseDate,
		sa.TrackNumber,
		sa.TrackNumber,
		sa.TrackNumber,
		sa.Lyrics,
		sa.FilePath)

	if err != nil {
		ss.session.Logger.Println(err)
		return err
	}
	return nil
}

// PrepareInsert creates a prepared statement to insert a new song.
func (ss *SongService) PrepareInsert() (*sql.Stmt, error) {
	insert :=
		`INSERT INTO songs 
		             (file_path, 
		              file_base, 
		              file_dir, 
		              artist_id, 
		              song_name, 
		              genre_id, 
		              release_date, 
		              track_number, 
		              disc_number, 
		              duration_in_millis, 
		              lyrics) 
		                       SELECT ?, 
		                              ?, 
		                              ?, 
		                              (SELECT artist_id 
		                                 FROM artists 
		                                WHERE artist_name = ? 
		                                  AND artist_sort = ?), 
		                              ?, 
		                              (SELECT genre_id 
		                                 FROM genres 
		                                WHERE genre_name = ?), 
		                              ?, 
		                              ?, 
		                              ?, 
		                              ?, 
		                              ? 
		             WHERE NOT EXISTS (SELECT 1 
		                                FROM songs 
		                               WHERE file_path = ?)`
	return ss.session.tx.Prepare(insert)
}

// Song queries the 'songs' table for a song  with the given ID and returns the
// result along with any error.
func (ss *SongService) Song(ID string) (*library.Song, error) {
	var s library.Song
	query :=
		`SELECT
		  songs.song_id,
		  songs.file_path,
		  songs.file_base,
		  songs.file_dir,
		  artists.artist_name,
		  artists.artist_sort,
		  genres.genre_name,
		  songs.song_name,
		  songs.release_date,
		  songs.track_number,
		  songs.lyrics
		FROM
		  song_discographies
		  INNER JOIN songs ON song_discographies.song_id = songs.song_id
		  INNER JOIN artists ON song_discographies.artist_id = artists.artist_id
		  INNER JOIN genres ON songs.genre_id = genres.genre_id
		WHERE 
		  songs.song_id = ?`
	err := ss.session.db.QueryRow(query, ID).Scan(
		&s.ID,
		&s.Attributes.FilePath,
		&s.Attributes.FileBase,
		&s.Attributes.FileDir,
		&s.Attributes.ArtistName,
		&s.Attributes.ArtistSort,
		&s.Attributes.GenreName,
		&s.Attributes.Name,
		&s.Attributes.ReleaseDate,
		&s.Attributes.TrackNumber,
		&s.Attributes.Lyrics)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		ss.session.Logger.Println(err)
		return &s, err
	}
	return &s, nil
}

// Songs queries the 'songs' table for all songs that meet the given
// criteria and returns the result along with any error.
func (ss *SongService) Songs(queries map[string]string) ([]*library.Song, error) {
	var results []*library.Song

	rows, err := ss.Query(queries)
	if err != nil {
		ss.session.Logger.Println(err)
		return results, err
	}

	for rows.Next() {
		var s library.Song
		err := rows.Scan(
			&s.ID,
			&s.Attributes.FilePath,
			&s.Attributes.FileBase,
			&s.Attributes.FileDir,
			&s.Attributes.ArtistName,
			&s.Attributes.ArtistSort,
			&s.Attributes.GenreName,
			&s.Attributes.Name,
			&s.Attributes.ReleaseDate,
			&s.Attributes.TrackNumber,
			&s.Attributes.Lyrics)
		s.Type = "songs"
		if err != nil {
			ss.session.Logger.Println(err)
			return results, err
		}
		results = append(results, &s)
	}

	if len(results) < 1 {
		return nil, nil
	}

	return results, nil
}

// Query executes a query for artists that meet the given predicate criteria and
// returns the results along with any error.
func (ss *SongService) Query(predicates map[string]string) (*sql.Rows, error) {
	query := bytes.NewBufferString(
		`SELECT
		  songs.song_id,
		  songs.file_path,
		  songs.file_base,
		  songs.file_dir,
		  artists.artist_name,
		  artists.artist_sort,
		  genres.genre_name,
		  songs.song_name,
		  songs.release_date,
		  songs.track_number,
		  songs.lyrics
		FROM
		  song_discographies
		  INNER JOIN songs ON song_discographies.song_id = songs.song_id
		  INNER JOIN artists ON song_discographies.artist_id = artists.artist_id
		  INNER JOIN genres ON songs.genre_id = genres.genre_id`)
	query, args := Where(query, predicates)
	return ss.session.db.Query(query.String(), args...)
}

// Close closes all open statements.
func (ss *SongService) Close() error {
	if ss.insert != nil {
		err := ss.insert.Close()
		if err != nil {
			ss.session.Logger.Println(err)
			return err
		}
	}

	return nil
}
