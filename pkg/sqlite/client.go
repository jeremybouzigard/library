package sqlite

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Registers database driver.
)

// Client represents a client to the underlying SQLite3 database.
type Client struct {
	Path   string
	db     *sql.DB
	Logger *log.Logger
}

// NewClient returns a new instance of a Client initialized with the given path.
func NewClient(path string) *Client {
	c := &Client{
		Path:   path,
		Logger: log.New(os.Stderr, "", log.LstdFlags)}
	return c
}

// Connect returns a new session to the SQLite3 database.
func (c *Client) Connect() *Session {
	s := newSession(c.db)
	return s
}

// Open opens and initializes the SQLite3 database.
func (c *Client) Open() error {
	db, err := sql.Open("sqlite3", c.Path)
	if err != nil {
		c.Logger.Println(err)
		return err
	}
	c.db = db
	return nil
}

// Close closes the underlying SQLite3 database.
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
