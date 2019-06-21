package sqlite

import (
	"database/sql"
	"github.com/jeremybouzigard/library"
	"log"
	"os"
)

// Session represents an open connection to the database.
type Session struct {
	db     *sql.DB
	tx     *sql.Tx
	Logger *log.Logger

	// Services
	albumService       AlbumService
	artistService      ArtistService
	genreService       GenreService
	songService        SongService
	LibraryService     library.Service
	AlbumDiscogService AlbumDiscogService
	SongDiscogService  SongDiscogService
}

// newSession returns a new instance of a Session attached to the database.
func newSession(db *sql.DB) *Session {
	s := &Session{
		db:     db,
		Logger: log.New(os.Stderr, "", log.LstdFlags)}
	s.genreService = NewGenreService(s)
	s.artistService = NewArtistService(s)
	s.albumService = NewAlbumService(s)
	s.songService = NewSongService(s)
	s.AlbumDiscogService = NewAlbumDiscogService(s)
	s.SongDiscogService = NewSongDiscogService(s)
	return s
}

// BeginTx starts a transaction within a Session.
func (s *Session) BeginTx() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	s.tx = tx
	return nil
}

// CommitTx commits the transaction within a Session.
func (s *Session) CommitTx() error {
	err := s.tx.Commit()
	if err != nil {
		return err
	}
	s.tx = nil
	return nil
}

// GenreService returns a genre service associated with this session.
func (s *Session) GenreService() library.GenreService {
	return &s.genreService
}

// AlbumService returns an album service associated with this session.
func (s *Session) AlbumService() library.AlbumService {
	return &s.albumService
}

// ArtistService returns an album service associated with this session.
func (s *Session) ArtistService() library.ArtistService {
	return &s.artistService
}

// SongService returns an album service associated with this session.
func (s *Session) SongService() library.SongService {
	return &s.songService
}
