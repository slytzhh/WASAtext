package database

import (
	"database/sql"
	"errors"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
)

func (db *appdbimpl) InsertChat(chat components.ChatCreation, userperformingid int) (int, int, error) {

	// default photo
	default_photo := "data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4NCg0KPCFET0NUWVBFIHN2ZyBQVUJMSUMgIi0vL1czQy8vRFREIFNWRyAxLjEvL0VOIiAiaHR0cDovL3d3dy53My5vcmcvR3JhcGhpY3MvU1ZHLzEuMS9EVEQvc3ZnMTEuZHRkIj4NCjwhLS0gVXBsb2FkZWQgdG86IFNWRyBSZXBvLCB3d3cuc3ZncmVwby5jb20sIEdlbmVyYXRvcjogU1ZHIFJlcG8gTWl4ZXIgVG9vbHMgLS0+CjxzdmcgZmlsbD0iIzAwMDAwMCIgdmVyc2lvbj0iMS4xIiBpZD0iTGF5ZXJfMSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgDQoJIHdpZHRoPSI4MDBweCIgaGVpZ2h0PSI4MDBweCIgdmlld0JveD0iNzk2IDc5NiAyMDAgMjAwIiBlbmFibGUtYmFja2dyb3VuZD0ibmV3IDc5NiA3OTYgMjAwIDIwMCIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSI+DQo8cGF0aCBkPSJNODk2LDc5NmMtNTUuMTQsMC05OS45OTksNDQuODYtOTkuOTk5LDEwMGMwLDU1LjE0MSw0NC44NTksMTAwLDk5Ljk5OSwxMDBjNTUuMTQxLDAsOTkuOTk5LTQ0Ljg1OSw5OS45OTktMTAwDQoJQzk5NS45OTksODQwLjg2LDk1MS4xNDEsNzk2LDg5Niw3OTZ6IE04OTYuNjM5LDgyNy40MjVjMjAuNTM4LDAsMzcuMTg5LDE5LjY2LDM3LjE4OSw0My45MjFjMCwyNC4yNTctMTYuNjUxLDQzLjkyNC0zNy4xODksNDMuOTI0DQoJcy0zNy4xODctMTkuNjY3LTM3LjE4Ny00My45MjRDODU5LjQ1Miw4NDcuMDg1LDg3Ni4xMDEsODI3LjQyNSw4OTYuNjM5LDgyNy40MjV6IE04OTYsOTgzLjg2DQoJYy0yNC42OTIsMC00Ny4wMzgtMTAuMjM5LTYzLjAxNi0yNi42OTVjLTIuMjY2LTIuMzM1LTIuOTg0LTUuNzc1LTEuODQtOC44MmM1LjQ3LTE0LjU1NiwxNS43MTgtMjYuNzYyLDI4LjgxNy0zNC43NjENCgljMi44MjgtMS43MjgsNi40NDktMS4zOTMsOC45MSwwLjgyOGM3LjcwNiw2Ljk1OCwxNy4zMTYsMTEuMTE0LDI3Ljc2NywxMS4xMTRjMTAuMjQ5LDAsMTkuNjktNC4wMDEsMjcuMzE4LTEwLjcxOQ0KCWMyLjQ4OC0yLjE5MSw2LjEyOC0yLjQ3OSw4LjkzMi0wLjcxMWMxMi42OTcsOC4wMDQsMjIuNjE4LDIwLjAwNSwyNy45NjcsMzQuMjUzYzEuMTQ0LDMuMDQ3LDAuNDI1LDYuNDgyLTEuODQyLDguODE3DQoJQzk0My4wMzcsOTczLjYyMSw5MjAuNjkxLDk4My44Niw4OTYsOTgzLjg2eiIvPg0KPC9zdmc+"

	// check if is a group
	var groupname sql.NullString
	var groupphoto sql.NullString
	groupname.Valid = true
	groupphoto.Valid = true
	if len(chat.GroupName) == 0 {
		groupname.Valid = false
		groupphoto.Valid = false
	} else {
		groupname.String = chat.GroupName
		// set the default image if not specified
		if len(chat.GroupPhoto) == 0 {
			groupphoto.String = default_photo
		} else {
			groupphoto.String = chat.GroupPhoto
		}
	}

	// start a transaction
	tx, err := db.c.Begin()
	if err != nil {
		return 0, 0, ErrTransaction
	}

	// insert the chat in db
	var chatid int
	err = tx.QueryRow(`INSERT INTO Chat(ChatName,ChatPhoto) VALUES(?,?) RETURNING ChatId`, groupname, groupphoto).Scan(&chatid)
	if err != nil {
		errtx := tx.Rollback()
		if errtx != nil {
			return 0, 0, ErrTransaction
		}
		return 0, 0, err
	}

	// insert all the user in ChatUser and also check if the user performing the action is in the list
	userinchat := false
	for i := 0; i < len(chat.UsernameList); i++ {

		// takes the id of the user
		var userid int
		userid, err = db.GetIdFromUsername(chat.UsernameList[i])
		if err != nil {
			return 0, 0, ErrUserNotFound
		}

		// create a row in ChatUser
		if userid == userperformingid {
			userinchat = true
			// sets the LastRead for the user creating the chat
			_, err = tx.Exec("INSERT INTO ChatUser(UserId,ChatId,LastRead) VALUES(?,?,CURRENT_TIMESTAMP)", userid, chatid)
			if err != nil {
				return 0, 0, err
			}
		} else {
			// for the other user LastRead is not set
			_, err = tx.Exec("INSERT INTO ChatUser(UserId,ChatId) VALUES(?,?)", userid, chatid)
			if err != nil {
				return 0, 0, err
			}
		}
	}

	// checks userinchat
	if !userinchat {
		errtx := tx.Rollback()
		if errtx != nil {
			return 0, 0, ErrTransaction
		}
		return 0, 0, ErrUserNotInChat
	}

	// if is a group commit the transaction
	if len(chat.GroupName) != 0 {
		err = tx.Commit()
		if err != nil {
			return 0, 0, ErrTransaction
		}

		return chatid, 0, err
	}

	// if the first message is forwarded, take the info of the message to forward
	if chat.ForwardedId != 0 {
		chat.FirstMessage.Text, chat.FirstMessage.Photo, err = db.GetMessage(chat.ForwardedId)
		if err != nil {
			errtx := tx.Rollback()
			if errtx != nil {
				return 0, 0, ErrTransaction
			}
			return 0, 0, err
		}
	}

	// check if there is a text in message
	var text sql.NullString
	text.Valid = true
	if len(chat.FirstMessage.Text) == 0 {
		text.Valid = false
	} else {
		text.String = chat.FirstMessage.Text
	}

	// check if there is a photo in message
	var photo sql.NullString
	photo.Valid = true
	if len(chat.FirstMessage.Photo) == 0 {
		photo.Valid = false
	} else {
		photo.String = chat.FirstMessage.Photo
	}

	// set the isForwarded flag
	isforwarded := false
	if chat.ForwardedId != 0 {
		isforwarded = true
	}

	// insert the first message
	var messageid int
	err = tx.QueryRow("INSERT INTO Message(ChatId,UserId,Text,Photo,IsForwarded) VALUES(?,?,?,?,?) RETURNING MessageId", chatid, userperformingid, text, photo, isforwarded).Scan(&messageid)
	if err != nil {
		errtx := tx.Rollback()
		if errtx != nil {
			return 0, 0, ErrTransaction
		}
		return 0, 0, err
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, 0, ErrTransaction
	}

	return chatid, messageid, err
}

func (db *appdbimpl) AddUsersToGroup(usernamelist []string, chatid int) error {

	// start a transaction
	tx, err := db.c.Begin()
	if err != nil {
		return ErrTransaction
	}

	// insert all the user in ChatUser
	for i := 0; i < len(usernamelist); i++ {

		// takes the id of the user
		var userid int
		userid, err = db.GetIdFromUsername(usernamelist[i])
		if err != nil {
			errtx := tx.Rollback()
			if errtx != nil {
				return ErrTransaction
			}
			return ErrUserNotFound
		}

		// create a row in ChatUser
		_, err = tx.Exec("INSERT OR IGNORE INTO ChatUser(UserId,ChatId) VALUES(?,?)", userid, chatid)
		if err != nil {
			errtx := tx.Rollback()
			if errtx != nil {
				return ErrTransaction
			}
			return err
		}
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return ErrTransaction
	}

	return err
}

func (db *appdbimpl) DeleteUserFromGroup(userid int, chatid int) error {

	_, err := db.c.Exec(`DELETE FROM ChatUser WHERE UserId=? AND ChatId=?`, userid, chatid)
	if err != nil {
		return err
	}

	// if the aren't other users in the chat, the chat is deleted
	var n_users int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM ChatUser WHERE ChatId=?`, chatid).Scan(&n_users)
	if err != nil {
		return err
	}

	if n_users == 0 {
		_, err = db.c.Exec(`DELETE FROM Chat WHERE ChatId=?`, chatid)
		if err != nil {
			return err
		}
	}

	return err
}

func (db *appdbimpl) IsGroup(chatid int) (bool, error) {

	var name sql.NullString
	err := db.c.QueryRow("SELECT ChatName FROM Chat WHERE ChatId=?", chatid).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrChatNotFound
		}
		return false, err
	}

	if !name.Valid {
		return false, err
	}

	return true, err
}

func (db *appdbimpl) ChangeGroupName(chatid int, groupname string) error {

	res, err := db.c.Exec("UPDATE Chat SET ChatName=? WHERE ChatId=?", groupname, chatid)
	if err != nil {
		return err
	}

	// check if the row effected are 0 which mean the chat don't exists
	eff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if eff == 0 {
		return ErrChatNotFound
	}

	return err
}

func (db *appdbimpl) ChangeGroupPhoto(chatid int, photo string) error {

	res, err := db.c.Exec("UPDATE Chat SET ChatPhoto=? WHERE ChatId=?", photo, chatid)
	if err != nil {
		return err
	}

	// check if the row effected are 0 which mean the chat don't exists
	eff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if eff == 0 {
		return ErrChatNotFound
	}

	return err
}

func (db *appdbimpl) GetUserChats(userid int) ([]components.ChatPreview, error) {

	chats := []components.ChatPreview{}

	// search all the groups where the user is
	rowsgroup, err := db.c.Query(`SELECT c.ChatId, c.ChatName, c.ChatPhoto, c.TimeCreated
	FROM ChatUser cu JOIN Chat c ON cu.ChatId = c.ChatId 
	WHERE cu.UserId = ? AND c.ChatName IS NOT NULL`, userid)
	if err != nil {
		return chats, err
	}

	defer rowsgroup.Close()

	// add all the group to the list
	for rowsgroup.Next() {

		var currentchat components.ChatPreview
		err = rowsgroup.Scan(&currentchat.ChatId, &currentchat.GroupName, &currentchat.GroupPhoto, &currentchat.TimeCreated)
		if err != nil {
			return chats, err
		}

		// takes the last message of the chat
		var text sql.NullString
		var photo sql.NullString
		err = db.c.QueryRow(`SELECT m.MessageId, m.ChatId, m.UserId, u.Username, m.Text, m.Photo, m.TimeStamp
		FROM Message m JOIN User u ON m.UserId = u.UserId
		WHERE m.ChatId = ? AND 
		m.MessageId = (SELECT MessageId FROM Message WHERE ChatId = m.ChatId ORDER BY Timestamp DESC, MessageId DESC LIMIT 1)`, currentchat.ChatId).Scan(&currentchat.LastMessage.MessageId, &currentchat.LastMessage.ChatId, &currentchat.LastMessage.UserId, &currentchat.LastMessage.Username, &text, &photo, &currentchat.LastMessage.TimeStamp)

		// if the last message exists
		if err == nil {

			// check if text and photo of the message are not null
			if text.Valid {
				currentchat.LastMessage.Text = text.String
			}
			if photo.Valid {
				currentchat.LastMessage.Photo = photo.String
			}

			// check if the message is received by all the users in the chat
			var err2 error
			currentchat.LastMessage.IsAllReceived, err2 = db.IsAllReceived(currentchat.LastMessage.MessageId, userid)
			if err2 != nil {
				return chats, err2
			}

			// check if the message is read by all the users in the chat
			currentchat.LastMessage.IsAllRead, err2 = db.IsAllRead(currentchat.LastMessage.MessageId, userid)
			if err2 != nil {
				return chats, err2
			}
		} else if !errors.Is(err, sql.ErrNoRows) {
			return chats, err
		}

		chats = append(chats, currentchat)
	}

	if rowsgroup.Err() != nil {
		return chats, err
	}

	// search all the chats where the user is with the last message, setting the name and photo equal to the other user of the chat
	rowschat, err := db.c.Query(`SELECT c.ChatId, u.Username, u.Photo, c.TimeCreated, m.MessageId, m.ChatId, m.UserId, u2.Username, m.Text, m.Photo, m.Timestamp
	FROM ChatUser cu JOIN Chat c ON cu.ChatId = c.ChatId JOIN Message m ON m.ChatId = c.ChatId JOIN ChatUser cu2 ON cu2.ChatId=cu.ChatId JOIN User u ON cu2.UserId = u.UserId JOIN User u2 ON u2.UserId = m.UserId
	WHERE cu.UserId = ? AND c.ChatName IS NULL AND u.UserId<>cu.UserId
	AND m.MessageId = (SELECT MessageId FROM Message WHERE ChatId = c.ChatId ORDER BY Timestamp DESC, MessageId DESC LIMIT 1)`, userid)
	if err != nil {
		return chats, err
	}

	defer rowschat.Close()

	// add all the chat to the list
	for rowschat.Next() {
		var currentchat components.ChatPreview
		var text sql.NullString
		var photo sql.NullString
		err = rowschat.Scan(&currentchat.ChatId, &currentchat.GroupName, &currentchat.GroupPhoto, &currentchat.TimeCreated, &currentchat.LastMessage.MessageId, &currentchat.LastMessage.ChatId, &currentchat.LastMessage.UserId, &currentchat.LastMessage.Username, &text, &photo, &currentchat.LastMessage.TimeStamp)
		if err != nil {
			return chats, err
		}

		// check if text and photo of the message are not null
		if text.Valid {
			currentchat.LastMessage.Text = text.String
		}
		if photo.Valid {
			currentchat.LastMessage.Photo = photo.String
		}

		// check if the message is received by all the users in the chat
		currentchat.LastMessage.IsAllReceived, err = db.IsAllReceived(currentchat.LastMessage.MessageId, userid)
		if err != nil {
			return chats, err
		}

		// check if the message is read by all the users in the chat
		currentchat.LastMessage.IsAllRead, err = db.IsAllRead(currentchat.LastMessage.MessageId, userid)
		if err != nil {
			return chats, err
		}

		chats = append(chats, currentchat)
	}

	if rowschat.Err() != nil {
		return chats, err
	}

	return chats, err
}

func (db *appdbimpl) GetChat(chatid int, userid int) (components.Chat, error) {

	var chat components.Chat

	isgroup, err := db.IsGroup(chatid)
	if err != nil {
		return chat, err
	}

	chat.IsGroup = isgroup

	// if the chat is a group takes the info, if not the name and group are equal to the other user of the chat
	if isgroup {
		err = db.c.QueryRow(`SELECT * FROM Chat WHERE ChatId=?`, chatid).Scan(&chat.ChatId, &chat.GroupName, &chat.GroupPhoto, &chat.TimeCreated)
	} else {
		err = db.c.QueryRow(`SELECT c.ChatId, u.Username, u.Photo, c.TimeCreated 
		FROM Chat c JOIN ChatUser cu ON c.ChatId=cu.ChatId JOIN User u ON cu.UserId=u.UserId
		WHERE c.ChatId=? AND u.UserId<>?`, chatid, userid).Scan(&chat.ChatId, &chat.GroupName, &chat.GroupPhoto, &chat.TimeCreated)
	}

	if err != nil {
		return chat, err
	}

	// gets all the users in the chat
	usersrows, err := db.c.Query(`SELECT u.Username FROM User u JOIN ChatUser cu ON u.UserId = cu.UserId
	WHERE cu.ChatId = ?`, chatid)
	if err != nil {
		return chat, err
	}

	defer usersrows.Close()

	// cicle for all the users
	for usersrows.Next() {
		var username string
		err = usersrows.Scan(&username)
		if err != nil {
			return chat, err
		}

		chat.UsernameList = append(chat.UsernameList, username)
	}

	if usersrows.Err() != nil {
		return chat, err
	}

	chat.MessageList = []components.Message{}

	// gets all the message in the chat
	messagerows, err := db.c.Query(`SELECT m.MessageId,m.ChatId,m.UserId,u.Username,m.Text,m.Photo,m.IsForwarded,m.Timestamp
	FROM Message m JOIN User u ON m.UserId = u.UserId
	WHERE ChatId=?`, chatid)
	if err != nil {
		return chat, err
	}

	defer messagerows.Close()

	// cicle for all the message
	for messagerows.Next() {
		var message components.Message
		var text sql.NullString
		var photo sql.NullString
		err = messagerows.Scan(&message.MessageId, &message.ChatId, &message.UserId, &message.Username, &text, &photo, &message.IsForwarded, &message.TimeStamp)
		if err != nil {
			return chat, err
		}

		if text.Valid {
			message.Text = text.String
		}
		if photo.Valid {
			message.Photo = photo.String
		}

		// check if the message is received by all the users in the chat
		message.IsAllReceived, err = db.IsAllReceived(message.MessageId, userid)
		if err != nil {
			return chat, err
		}

		// check if the message is read by all the users in the chat
		message.IsAllRead, err = db.IsAllRead(message.MessageId, userid)
		if err != nil {
			return chat, err
		}

		// check if the message is a reply to another message
		var replyid sql.NullInt32
		err = db.c.QueryRow(`SELECT RepliedId FROM Message WHERE MessageId=?`, message.MessageId).Scan(&replyid)
		if err != nil {
			return chat, err
		}

		if replyid.Valid {
			message.ReplyMessage = components.MessagePreview{}
			var replytext sql.NullString
			var replyphoto sql.NullString
			err = db.c.QueryRow(`SELECT m.MessageId,m.ChatId,m.UserId,u.Username,m.Text,m.Photo,m.Timestamp
						FROM Message m JOIN User u ON m.UserId = u.UserId
						WHERE MessageId=?`, replyid).Scan(&message.ReplyMessage.MessageId, &message.ReplyMessage.ChatId, &message.ReplyMessage.UserId, &message.ReplyMessage.Username, &replytext, &replyphoto, &message.ReplyMessage.TimeStamp)
			if err != nil {
				return chat, err
			}

			if replytext.Valid {
				message.ReplyMessage.Text = replytext.String
			}
			if replyphoto.Valid {
				message.ReplyMessage.Photo = replyphoto.String
			}

		}

		message.CommentList = []components.Comment{}

		// gets all the comment of the message
		commentrows, err := db.c.Query(`SELECT c.MessageId, c.UserId, c.Emoji, u.Username
		FROM Comment c JOIN User u ON c.UserId=u.UserId
		WHERE MessageId=?`, message.MessageId)
		if err != nil {
			return chat, err
		}

		defer commentrows.Close()

		// cicle for all the comments
		for commentrows.Next() {
			var comment components.Comment
			err = commentrows.Scan(&comment.MessageId, &comment.UserId, &comment.Emoji, &comment.Username)
			if err != nil {
				return chat, err
			}

			// append the comment to the message
			message.CommentList = append(message.CommentList, comment)
		}

		if commentrows.Err() != nil {
			return chat, err
		}

		// append the message to the chat
		chat.MessageList = append(chat.MessageList, message)
	}

	if messagerows.Err() != nil {
		return chat, err
	}

	return chat, err
}

func (db *appdbimpl) GetUsersInChat(chatid int) ([]components.Username, error) {

	// gets all the users in the chat
	userList := []components.Username{}
	userrows, err := db.c.Query(`SELECT u.Username FROM User u JOIN ChatUser cu ON u.UserId = cu.UserId
	WHERE cu.ChatId = ?`, chatid)
	if err != nil {
		return userList, err
	}

	defer userrows.Close()

	// cicle for all the users
	for userrows.Next() {
		var user components.Username
		err = userrows.Scan(&user.Username)
		if err != nil {
			return userList, err
		}

		userList = append(userList, user)
	}

	if userrows.Err() != nil {
		return userList, err
	}

	return userList, err

}
