package ProjectTapebucket

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	"github.com/gidoBOSSftw5731/log"
)

// AddToDB adds a block of text to the database and returns both the url and error
func AddToDB(text, identity *string, db *sql.DB) (string, error) {
	textHash := fmt.Sprintf("%x", sha256.Sum256([]byte(*text)))

	url := textHash[:8]

	test := db.QueryRow("SELECT hash FROM tapebucket WHERE url=?", url)
	switch {
	case test.Scan().Error() == fmt.Sprint(sql.ErrNoRows):
		log.Debug("New paste, adding..")
		_, err := db.Exec("INSERT INTO tapebucket VALUES(?, ?, ?, ?)",
			*text, textHash, url, *identity)
		if err != nil {
			return "", fmt.Errorf(test.Scan().Error())
		}

		//defer insert.Close()
		log.Debug("Added paste to table")
	case test.Scan().Error() != "":
		return "", fmt.Errorf(test.Scan().Error())		
	}
}
