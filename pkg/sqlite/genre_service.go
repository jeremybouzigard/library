package sqlite

import (
	"database/sql"

	"github.com/jeremybouzigard/library"
)

// GenreService manages interactions with the genres data source.
type GenreService struct {
	insert  *sql.Stmt
	session *Session
	Select  *sql.Stmt
}

// NewGenreService returns a new instance of an GenreService that operates
// within the given session.
func NewGenreService(s *Session) GenreService {
	gs := GenreService{session: s}
	return gs
}

// CreateTable creates the 'genres' table and returns any errors.
func (service *GenreService) CreateTable() (sql.Result, error) {
	create :=
		`CREATE TABLE IF NOT EXISTS genres (
			genre_id   INTEGER PRIMARY KEY,
			genre_name TEXT    UNIQUE NOT NULL
		)`
	return service.session.tx.Exec(create)
}

// DropTable drops the 'genres' table and returns any errors.
func (service *GenreService) DropTable() (sql.Result, error) {
	drop := `DROP TABLE IF EXISTS genres`
	return service.session.tx.Exec(drop)
}

// CreateGenre inserts a new genre.
func (service *GenreService) CreateGenre(attributes *library.GenreAttributes) error {
	if service.insert == nil {
		stmt, err := service.prepareInsert()
		if err != nil {
			service.session.Logger.Println(err)
			return err
		}
		service.insert = stmt
	}

	_, err := service.insert.Exec(attributes.Name, attributes.Name)
	if err != nil {
		service.session.Logger.Println(err)
		return err
	}
	return nil
}

// prepareInsert creates a prepared statement to insert a new genre.
func (service *GenreService) prepareInsert() (*sql.Stmt, error) {
	insert :=
		`     INSERT INTO genres (genre_name) 
		           SELECT ?
		 WHERE NOT EXISTS (SELECT 1 
		                     FROM genres 
							WHERE genre_name = ?)`
	return service.session.tx.Prepare(insert)
}

// Genre queries the 'genres' table for a genre with the given ID and returns
// the result along with any error.
func (service *GenreService) Genre(ID string) (*library.Genre, error) {
	var g library.Genre
	query := `SELECT genre_id, genre_name FROM genres WHERE genre_id = ?`

	err := service.session.db.QueryRow(query, ID).Scan(
		&g.ID,
		&g.Attributes.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return &g, nil
		}
		service.session.Logger.Println(err)
		return &g, err
	}
	return &g, nil
}

// Genres queries the 'genres' table for all distinct genres and returns the
// result along with any error.
func (service *GenreService) Genres() ([]*library.Genre, error) {
	var results []*library.Genre

	query := `SELECT genre_id, genre_name FROM genres`
	rows, err := service.session.db.Query(query)
	if err != nil {
		service.session.Logger.Println(err)
		return results, err
	}

	for rows.Next() {
		var res library.Genre
		res.Type = "genres"
		err := rows.Scan(
			&res.ID,
			&res.Attributes.Name)
		if err != nil {
			service.session.Logger.Println(err)
			return results, err
		}
		results = append(results, &res)
	}
	return results, nil
}

// Close closes all open statements.
func (service *GenreService) Close() error {
	if service.insert != nil {
		err := service.insert.Close()
		if err != nil {
			service.session.Logger.Println(err)
			return err
		}
	}

	if service.Select != nil {
		err := service.insert.Close()
		if err != nil {
			service.session.Logger.Println(err)
			return err
		}
	}

	return nil
}
