package library

// Service manages interactions with the media library.
type Service interface {
	CreateLibrary() error
	DeleteLibrary() error
}
