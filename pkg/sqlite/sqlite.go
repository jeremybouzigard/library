package sqlite

import (
	"bytes"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Registers database driver.
)

// Where appends WHERE clauses to the query using the given predicates.
func Where(query *bytes.Buffer, predicates map[string]string) (*bytes.Buffer, []interface{}) {
	args := []interface{}{}

	artistID := predicates["artistID"]
	if len(artistID) > 0 {
		if len(args) > 0 {
			query.WriteString(` AND artists.artist_id = ?`)
		} else {
			query.WriteString(` WHERE artists.artist_id = ?`)
		}
		args = append(args, artistID)
	}

	albumID := predicates["albumID"]
	if len(albumID) > 0 {
		if len(args) > 0 {
			query.WriteString(` AND albums.album_id = ?`)
		} else {
			query.WriteString(` WHERE albums.album_id = ?`)
		}
		args = append(args, albumID)
	}

	genreID := predicates["genreID"]
	if len(genreID) > 0 {
		if len(args) > 0 {
			query.WriteString(` AND genres.genre_id = ?`)
		} else {
			query.WriteString(` WHERE genres.genre_id = ?`)
		}
		args = append(args, genreID)
	}

	fmt.Println(query)
	return query, args
}

// Select executes the given query and returns the results represented by a
// slice of strings.
// func (tx *Tx) Select(query string) ([]string, error) {
// 	var results []string

// 	rows, err := tx.Query(query)
// 	if err != nil {
// 		return results, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var res string
// 		err := rows.Scan(&res)
// 		if err != nil {
// 			return results, err
// 		}

// 		results = append(results, res)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		return results, err
// 	}
// 	return results, nil
// }
