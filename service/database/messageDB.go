package database

import (
	"database/sql"
	"errors"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
)

func (db *appdbimpl) InsertMessage(message components.MessageToSend, isforwarded bool, chatid int, userperformingid int) (int, error) {

	// check if there is a text in message
	var text sql.NullString
	text.Valid = true
	if len(message.Text) == 0 {
		text.Valid = false
	} else {
		text.String = message.Text
	}

	// check if there is a photo in message
	var photo sql.NullString
	photo.Valid = true
	if len(message.Photo) == 0 {
		photo.Valid = false
	} else {
		photo.String = message.Photo
	}

	// check if the message is a reply
	var replyid sql.NullInt32
	replyid.Valid = true
	if message.ReplyId == 0 {
		replyid.Valid = false
	} else {
		replyid.Int32 = int32(message.ReplyId)
	}

	var messageid int
	err := db.c.QueryRow(`INSERT INTO Message(ChatId,UserId,Text,Photo,IsForwarded,RepliedId) VALUES(?,?,?,?,?,?) RETURNING MessageId`, chatid, userperformingid, text, photo, isforwarded, replyid).Scan(&messageid)
	if err != nil {
		return 0, err
	}

	return messageid, err
}

func (db *appdbimpl) GetMessage(messageid int) (string, string, error) {

	var textnull sql.NullString
	var photonull sql.NullString
	err := db.c.QueryRow(`SELECT Text,Photo FROM Message WHERE MessageId=?`, messageid).Scan(&textnull, &photonull)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrMessNotFound
		}
		return "", "", err
	}

	var text string
	var photo string

	// copies the values that can be NULL
	if textnull.Valid {
		text = textnull.String
	}
	if photonull.Valid {
		photo = photonull.String
	}

	return text, photo, err
}

func (db *appdbimpl) IsMessageInChat(chatid int, messageid int) (bool, error) {

	messageinchat := false
	err := db.c.QueryRow(`SELECT EXISTS(SELECT * FROM Message WHERE ChatId=? AND MessageId=?)`, chatid, messageid).Scan(&messageinchat)
	if err != nil {
		return false, err
	}

	return messageinchat, err
}

func (db *appdbimpl) DeleteMessage(messageid int, chatid int) error {

	// set to NULL the RepliedId on all the message that replayed to this
	_, err := db.c.Exec(`UPDATE Message SET RepliedId=NULL WHERE RepliedId=?`, messageid)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`DELETE FROM Message WHERE MessageId=?`, messageid, chatid)
	if err != nil {
		return err
	}

	// check if the chat is a group
	isgroup, err := db.IsGroup(chatid)
	if err != nil {
		return err
	}

	// if the aren't other message and is a chat, the chat is deleted
	var n_messages int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM Message WHERE ChatId=?`, chatid).Scan(&n_messages)
	if err != nil {
		return err
	}

	if n_messages == 0 && !isgroup {
		_, err = db.c.Exec(`DELETE FROM Chat WHERE ChatId=?`, chatid)
		if err != nil {
			return err
		}
	}

	return err
}

func (db *appdbimpl) GetUserFromMessage(messageid int) (int, error) {

	var userid int
	err := db.c.QueryRow(`SELECT UserId FROM Message WHERE MessageId=?`, messageid).Scan(&userid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrMessNotFound
		}
		return 0, err
	}

	return userid, err
}

func (db *appdbimpl) IsAllReceived(messageid int, userid int) (bool, error) {

	allreceived := false
	err := db.c.QueryRow(`SELECT NOT EXISTS(SELECT *
	FROM Message m JOIN ChatUser cu ON m.ChatId=cu.ChatId JOIN User u ON cu.UserId=u.UserId
	WHERE m.MessageId=? AND m.TimeStamp>u.LastAccess AND u.UserId<>?)`, messageid, userid).Scan(&allreceived)
	if err != nil {
		return false, err
	}

	return allreceived, err
}

func (db *appdbimpl) IsAllRead(messageid int, userid int) (bool, error) {

	allread := false
	err := db.c.QueryRow(`SELECT NOT EXISTS(SELECT *
	FROM Message m JOIN ChatUser cu ON m.ChatId=cu.ChatId
	WHERE m.MessageId=? AND cu.UserId<>? AND (m.TimeStamp>cu.LastRead OR cu.LastRead IS NULL))`, messageid, userid).Scan(&allread)
	if err != nil {
		return false, err
	}

	return allread, err
}
