package library

// Album represents an album resource object.
type Album struct {
	Type       string          `json:"type,omitempty"`
	ID         string          `json:"id,omitempty"`
	Attributes AlbumAttributes `json:"attributes,omitempty"`
}

// AlbumAttributes represents information about the album resource object.
type AlbumAttributes struct {
	Name            string `json:"name,omitempty"`
	Sort            string `json:"sort,omitempty"`
	ArtistName      string `json:"artistName,omitempty"`
	ArtistSort      string `json:"artistSort,omitempty"`
	GenreName       string `json:"genreName,omitempty"`
	ReleaseDate     string `json:"releaseDate,omitempty"`
	AlbumArtist     string `json:"albumArtist,omitempty"`
	AlbumArtistSort string `json:"albumArtistSort,omitempty"`
}

// AlbumService manages interactions with the album data source.
type AlbumService interface {
	Album(ID string) (*Album, error)
	Albums(params map[string]string) ([]*Album, error)
	CreateAlbum(attributes *AlbumAttributes) error
}
