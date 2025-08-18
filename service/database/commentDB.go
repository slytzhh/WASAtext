package database

import (
	"strings"

	"github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) InsertComment(messageid int, userid int, emoji string) error {

	// check if the user has already commented his message
	alreadycommented := false
	err := db.c.QueryRow(`SELECT EXISTS(SELECT * FROM Comment WHERE MessageId=? AND UserId=?)`, messageid, userid).Scan(&alreadycommented)
	if err != nil {
		return err
	}

	// if the user has already commented update the comment,
	// otherwise inserts the comment
	if alreadycommented {
		_, err = db.c.Exec(`UPDATE Comment SET Emoji=? WHERE MessageId=? AND UserId=?`, emoji, messageid, userid)
		if err != nil {
			return err
		}
	} else {
		_, err = db.c.Exec(`INSERT INTO Comment(MessageId,UserId,Emoji) VALUES(?,?,?)`, messageid, userid, emoji)
		if err != nil {
			if strings.Contains(err.Error(), sqlite3.ErrConstraintForeignKey.Error()) {
				return ErrMessNotFound
			}
			return err
		}
	}

	return err
}

func (db *appdbimpl) DeleteComment(messageid int, userid int) error {

	_, err := db.c.Exec(`DELETE FROM Comment WHERE MessageId=? AND UserId=?`, messageid, userid)
	if err != nil {
		return err
	}

	return err
}
