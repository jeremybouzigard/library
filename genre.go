package library

// Genre represents a genre resource object.
type Genre struct {
	Type       string          `json:"type,omitempty"`
	ID         string          `json:"id,omitempty"`
	Attributes GenreAttributes `json:"attributes,omitempty"`
}

// GenreAttributes represents information about the genre resource object.
type GenreAttributes struct {
	Name string `json:"name,omitempty"`
}

// GenreService manages interactions with the genres data source.
type GenreService interface {
	Genre(ID string) (*Genre, error)
	Genres() ([]*Genre, error)
	CreateGenre(attributes *GenreAttributes) error
}
