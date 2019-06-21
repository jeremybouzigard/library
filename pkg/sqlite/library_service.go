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
	session *Session
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
	if ls.session != nil {
		return fmt.Errorf("library session already opened")
	}
	ls.client.Open()
	ls.session = ls.client.Connect()
	return nil
}

// Close closes the current library session.
func (ls *Service) Close() error {
	if ls.session == nil {
		return fmt.Errorf("no open library session to close")
	}
	ls.client.Close()
	ls.session = nil
	return nil
}

// CreateLibrary creates library tables in the data source.
func (ls *Service) CreateLibrary() error {
	err := ls.session.BeginTx()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	_, err = ls.session.genreService.CreateTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	_, err = ls.session.artistService.CreateTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	_, err = ls.session.albumService.CreateTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	_, err = ls.session.songService.CreateTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	_, err = ls.session.AlbumDiscogService.CreateTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	_, err = ls.session.SongDiscogService.CreateTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	err = ls.session.CommitTx()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	return nil
}

// DeleteLibrary deletes all library data and drops tables from the data source.
func (ls *Service) DeleteLibrary() error {
	err := ls.session.BeginTx()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	defer ls.session.albumService.Close()
	defer ls.session.AlbumDiscogService.Close()
	defer ls.session.artistService.Close()
	defer ls.session.genreService.Close()
	defer ls.session.SongDiscogService.Close()

	_, err = ls.session.genreService.DropTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	_, err = ls.session.artistService.DropTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	_, err = ls.session.albumService.DropTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	_, err = ls.session.songService.DropTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	_, err = ls.session.AlbumDiscogService.DropTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	_, err = ls.session.SongDiscogService.DropTable()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	err = ls.session.CommitTx()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	return nil
}

// AddPath adds media data within the given path to the library.
func (ls *Service) AddPath(path string) error {
	err := ls.session.BeginTx()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}

	ms := metadata.Service{}
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		metadata, err := ms.Metadata(path)
		if err != nil {
			ls.session.Logger.Println(err)
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

			ls.session.genreService.CreateGenre(&genre)
			ls.session.artistService.CreateArtist(&artist)
			ls.session.albumService.CreateAlbum(&album)
			ls.session.songService.CreateSong(&song)
			ls.session.AlbumDiscogService.CreateAlbumDiscog(&album)
			ls.session.SongDiscogService.CreateSongDiscog(&song, &album)
		}

		return nil
	})

	err = ls.session.CommitTx()
	if err != nil {
		ls.session.Logger.Println(err)
		return err
	}
	return nil
}
