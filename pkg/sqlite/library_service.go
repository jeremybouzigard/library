package sqlite

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeremybouzigard/library"
	"github.com/jeremybouzigard/metadata/pkg/metadata"
)

// Service manages interactions with the media library.
type Service struct {
	client  *Client
	Session *Session
}

// NewService returns a new instance of a Service that operates on the library
// at the given path.
func NewService(path string) Service {
	client := NewClient(path)
	ls := Service{client: client}
	return ls
}

// Open opens and initializes a new library session.
func (ls *Service) Open() error {
	if ls.Session != nil {
		return fmt.Errorf("library session already opened")
	}
	ls.client.Open()
	ls.Session = ls.client.Connect()
	return nil
}

// Close closes the current library session.
func (ls *Service) Close() error {
	if ls.Session == nil {
		return fmt.Errorf("no open library session to close")
	}
	ls.client.Close()
	ls.Session = nil
	return nil
}

// CreateLibrary creates library tables in the data source.
func (ls *Service) CreateLibrary() error {
	err := ls.Session.BeginTx()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	_, err = ls.Session.genreService.CreateTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	_, err = ls.Session.artistService.CreateTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	_, err = ls.Session.albumService.CreateTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	_, err = ls.Session.songService.CreateTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	_, err = ls.Session.AlbumDiscogService.CreateTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	_, err = ls.Session.SongDiscogService.CreateTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	err = ls.Session.CommitTx()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	return nil
}

// DeleteLibrary deletes all library data and drops tables from the data source.
func (ls *Service) DeleteLibrary() error {
	err := ls.Session.BeginTx()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	defer ls.Session.albumService.Close()
	defer ls.Session.AlbumDiscogService.Close()
	defer ls.Session.artistService.Close()
	defer ls.Session.genreService.Close()
	defer ls.Session.SongDiscogService.Close()

	_, err = ls.Session.genreService.DropTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	_, err = ls.Session.artistService.DropTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	_, err = ls.Session.albumService.DropTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	_, err = ls.Session.songService.DropTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	_, err = ls.Session.AlbumDiscogService.DropTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	_, err = ls.Session.SongDiscogService.DropTable()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	err = ls.Session.CommitTx()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	return nil
}

// AddPath adds media data within the given path to the library.
func (ls *Service) AddPath(path string) error {
	err := ls.Session.BeginTx()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}

	ms := metadata.Service{}
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		metadata, err := ms.Metadata(path)
		if err != nil {
			ls.Session.Logger.Println(err)
		} else {
			genre := library.GenreAttributes{
				Name: metadata.Genre}

			artist := library.ArtistAttributes{
				Name: metadata.Artist,
				Sort: metadata.ArtistSort}

			album := library.AlbumAttributes{
				Name:        metadata.Album,
				Sort:        metadata.AlbumSort,
				ArtistName:  metadata.Artist,
				ArtistSort:  metadata.ArtistSort,
				GenreName:   metadata.Genre,
				ReleaseDate: metadata.Year}

			song := library.SongAttributes{
				FilePath:    path,
				FileBase:    filepath.Base(path),
				FileDir:     filepath.Dir(path),
				ArtistName:  metadata.Artist,
				ArtistSort:  metadata.ArtistSort,
				Name:        metadata.Title,
				GenreName:   metadata.Genre,
				TrackNumber: metadata.Track,
				ReleaseDate: metadata.Year,
				Lyrics:      metadata.Lyrics,
				Comments:    metadata.Comment}

			ls.Session.genreService.CreateGenre(&genre)
			ls.Session.artistService.CreateArtist(&artist)
			ls.Session.albumService.CreateAlbum(&album)
			ls.Session.songService.CreateSong(&song)
			ls.Session.AlbumDiscogService.CreateAlbumDiscog(&album)
			ls.Session.SongDiscogService.CreateSongDiscog(&song, &album)
		}

		return nil
	})

	err = ls.Session.CommitTx()
	if err != nil {
		ls.Session.Logger.Println(err)
		return err
	}
	return nil
}
