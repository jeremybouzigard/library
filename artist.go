package library

// Artist represents an artist resource object.
type Artist struct {
	Type       string           `json:"type,omitempty"`
	ID         string           `json:"id,omitempty"`
	Attributes ArtistAttributes `json:"attributes,omitempty"`
}

// ArtistAttributes represents information about the artist resource object.
type ArtistAttributes struct {
	Name string `json:"name,omitempty"`
	Sort string `json:"sort,omitempty"`
}

// ArtistService manages interactions with the artist data source.
type ArtistService interface {
	Artist(ID string) (*Artist, error)
	Artists(params map[string]string) ([]*Artist, error)
	CreateArtist(attributes *ArtistAttributes) error
}
