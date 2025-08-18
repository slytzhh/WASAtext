/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// User
	InsertUser(username string) (int, string, error)
	ChangeUsername(userid int, username string) error
	ChangeUserPhoto(userid int, photo string) error
	IsUserInChat(chatid int, userid int) (bool, error)
	GetUsernameFromId(userid int) (string, error)
	GetIdFromUsername(username string) (int, error)
	SetLastAccess(userid int) error
	SetLastRead(userid int, chatid int) error
	SearchUsers(username string) ([]components.User, error)

	// Chat
	InsertChat(chat components.ChatCreation, userperformingid int) (int, int, error)
	AddUsersToGroup(usernamelist []string, chatid int) error
	DeleteUserFromGroup(userid int, chatid int) error
	IsGroup(chatid int) (bool, error)
	ChangeGroupName(chatid int, groupname string) error
	ChangeGroupPhoto(chatid int, photo string) error
	GetUserChats(userid int) ([]components.ChatPreview, error)
	GetChat(chatid int, userid int) (components.Chat, error)
	GetUsersInChat(chatid int) ([]components.Username, error)

	// Message
	InsertMessage(message components.MessageToSend, isforwarded bool, chatid int, userperformingid int) (int, error)
	GetMessage(messageid int) (string, string, error)
	IsMessageInChat(chatid int, messageid int) (bool, error)
	DeleteMessage(messageid int, chatid int) error
	GetUserFromMessage(messageid int) (int, error)
	IsAllReceived(messageid int, userid int) (bool, error)
	IsAllRead(messageid int, userid int) (bool, error)

	// Comment
	InsertComment(messageid int, userid int, emoji string) error
	DeleteComment(messageid int, userid int) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	User := `CREATE TABLE IF NOT EXISTS User(
				UserId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				Username TEXT NOT NULL UNIQUE,
				Photo TEXT NOT NULL,
				LastAccess DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				CHECK(LENGTH(Username)>=3 AND LENGTH(Username)<=16)
				);`

	Chat := `CREATE TABLE IF NOT EXISTS Chat(
				ChatId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				ChatName TEXT,
				ChatPhoto TEXT,
				TimeCreated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				CHECK ((ChatName IS NULL AND ChatPhoto IS NULL) OR (ChatName IS NOT NULL AND ChatPhoto IS NOT NULL))
				);`

	Message := `CREATE TABLE IF NOT EXISTS Message(
					MessageId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					ChatId INTEGER NOT NULL,
					UserId INTEGER NOT NULL,
					Text TEXT,
					Photo TEXT,
					IsForwarded BOOLEAN NOT NULL,
					RepliedId INTEGER,
					TimeStamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					CHECK (Text IS NOT NULL OR Photo IS NOT NULL),
					FOREIGN KEY(ChatId) REFERENCES Chat(ChatId) ON DELETE CASCADE,
					FOREIGN KEY(UserId) REFERENCES User(UserId),
					FOREIGN KEY(RepliedId) REFERENCES Message(MessageId)
					);`

	Comment := `CREATE TABLE IF NOT EXISTS Comment(
					MessageId INTEGER NOT NULL,
					UserId INTEGER NOT NULL,
					Emoji TEXT,
					PRIMARY KEY(MessageId,UserId),
					FOREIGN KEY(MessageId) REFERENCES Message(MessageId) ON DELETE CASCADE,
					FOREIGN KEY(UserId) REFERENCES User(UserId)
					);`

	ChatUser := `CREATE TABLE IF NOT EXISTS ChatUser(
					UserId INTEGER NOT NULL,
					ChatId INTEGER NOT NULL,
					TimeAdded DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					LastRead DATETIME,
					PRIMARY KEY(UserId,ChatId),
					CHECK((LastRead>=TimeAdded) OR (LastRead IS NULL)),
					FOREIGN KEY(ChatId) REFERENCES Chat(ChatId) ON DELETE CASCADE,
					FOREIGN KEY(UserId) REFERENCES User(UserId)
					);`

	var err error
	_, err = db.Exec("PRAGMA foreign_keys=ON")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(User)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(Chat)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(Message)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(Message)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(Comment)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(ChatUser)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
