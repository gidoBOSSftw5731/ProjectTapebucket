package ProjectTapebucket

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	"github.com/gidoBOSSftw5731/log"
)

const (
	URLLen = 8
)

// AddToDB adds a block of text to the database and returns both the url and error
// database should be mysql
func AddToDB(text, identity *string, db *sql.DB) (string, error) {
	textHash := fmt.Sprintf("%x", sha256.Sum256([]byte(*text)))

	url := textHash[:URLLen]
	log.Traceln(textHash)

	var out bool
	err := db.QueryRow("SELECT 1 FROM tapebucket WHERE url=?", url).Scan(&out)
	switch err {
	case sql.ErrNoRows:
		log.Debug("New paste, adding..")
		_, err := db.Exec("INSERT INTO tapebucket VALUES(?, ?, ?, ?)",
			*text, textHash, url, *identity)
		if err != nil {
			return "", err
		}

		//defer insert.Close()
		log.Debug("Added paste to table")
	case nil:
	default:
		return "", err
	}

	return url, nil
}

// ReturnFromDB returns a pointer to the paste from a database and an error.
func ReturnFromDB(url string, db *sql.DB) (*string, error) {
	if len(url) != URLLen {
		return nil, fmt.Errorf("badurlformat")
	}
	// so now we know the URL could be valid

	var paste string
	err := db.QueryRow("SELECT paste FROM tapebucket WHERE url=?", url).Scan(&paste)
	if err != nil {
		return nil, err
	}
	return &paste, nil
}

/*
CREATE TABLE tapebucket (
	paste LONGTEXT,
	hash char(64),
	url char(8),
	identity tinytext);
*/
