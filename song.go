package library

// Song represents a song resource object.
type Song struct {
	Type       string         `json:"type,omitempty"`
	ID         string         `json:"id,omitempty"`
	Attributes SongAttributes `json:"attributes,omitempty"`
}

// SongAttributes represents information about the song resource object.
type SongAttributes struct {
	FilePath    string `json:"filePath,omitempty"`
	FileBase    string `json:"fileBase,omitempty"`
	FileDir     string `json:"fileDir,omitempty"`
	ArtistName  string `json:"artistName,omitempty"`
	ArtistSort  string `json:"artistSort,omitempty"`
	Name        string `json:"name,omitempty"`
	GenreName   string `json:"genreName,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	TrackNumber string `json:"trackNumber,omitempty"`
	Lyrics      string `json:"lyrics,omitempty"`
	Comments    string `json:"comments,omitempty"`
}

// SongService manages interactions with the song data source.
type SongService interface {
	Song(ID string) (*Song, error)
	Songs(params map[string]string) ([]*Song, error)
	CreateSong(attributes *SongAttributes) error
}
