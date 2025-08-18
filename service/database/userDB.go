package database

import (
	"database/sql"
	"errors"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) InsertUser(username string) (int, string, error) {

	// default photo
	default_photo := "data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4NCg0KPCFET0NUWVBFIHN2ZyBQVUJMSUMgIi0vL1czQy8vRFREIFNWRyAxLjEvL0VOIiAiaHR0cDovL3d3dy53My5vcmcvR3JhcGhpY3MvU1ZHLzEuMS9EVEQvc3ZnMTEuZHRkIj4NCjwhLS0gVXBsb2FkZWQgdG86IFNWRyBSZXBvLCB3d3cuc3ZncmVwby5jb20sIEdlbmVyYXRvcjogU1ZHIFJlcG8gTWl4ZXIgVG9vbHMgLS0+CjxzdmcgZmlsbD0iIzAwMDAwMCIgdmVyc2lvbj0iMS4xIiBpZD0iTGF5ZXJfMSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgDQoJIHdpZHRoPSI4MDBweCIgaGVpZ2h0PSI4MDBweCIgdmlld0JveD0iNzk2IDc5NiAyMDAgMjAwIiBlbmFibGUtYmFja2dyb3VuZD0ibmV3IDc5NiA3OTYgMjAwIDIwMCIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSI+DQo8cGF0aCBkPSJNODk2LDc5NmMtNTUuMTQsMC05OS45OTksNDQuODYtOTkuOTk5LDEwMGMwLDU1LjE0MSw0NC44NTksMTAwLDk5Ljk5OSwxMDBjNTUuMTQxLDAsOTkuOTk5LTQ0Ljg1OSw5OS45OTktMTAwDQoJQzk5NS45OTksODQwLjg2LDk1MS4xNDEsNzk2LDg5Niw3OTZ6IE04OTYuNjM5LDgyNy40MjVjMjAuNTM4LDAsMzcuMTg5LDE5LjY2LDM3LjE4OSw0My45MjFjMCwyNC4yNTctMTYuNjUxLDQzLjkyNC0zNy4xODksNDMuOTI0DQoJcy0zNy4xODctMTkuNjY3LTM3LjE4Ny00My45MjRDODU5LjQ1Miw4NDcuMDg1LDg3Ni4xMDEsODI3LjQyNSw4OTYuNjM5LDgyNy40MjV6IE04OTYsOTgzLjg2DQoJYy0yNC42OTIsMC00Ny4wMzgtMTAuMjM5LTYzLjAxNi0yNi42OTVjLTIuMjY2LTIuMzM1LTIuOTg0LTUuNzc1LTEuODQtOC44MmM1LjQ3LTE0LjU1NiwxNS43MTgtMjYuNzYyLDI4LjgxNy0zNC43NjENCgljMi44MjgtMS43MjgsNi40NDktMS4zOTMsOC45MSwwLjgyOGM3LjcwNiw2Ljk1OCwxNy4zMTYsMTEuMTE0LDI3Ljc2NywxMS4xMTRjMTAuMjQ5LDAsMTkuNjktNC4wMDEsMjcuMzE4LTEwLjcxOQ0KCWMyLjQ4OC0yLjE5MSw2LjEyOC0yLjQ3OSw4LjkzMi0wLjcxMWMxMi42OTcsOC4wMDQsMjIuNjE4LDIwLjAwNSwyNy45NjcsMzQuMjUzYzEuMTQ0LDMuMDQ3LDAuNDI1LDYuNDgyLTEuODQyLDguODE3DQoJQzk0My4wMzcsOTczLjYyMSw5MjAuNjkxLDk4My44Niw4OTYsOTgzLjg2eiIvPg0KPC9zdmc+"

	// check if the user already exists
	var user_exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT * FROM User WHERE Username=?)`, username).Scan(&user_exists)
	if err != nil {
		return 0, "", err
	}

	var id int
	var photo string
	if !user_exists {
		// insert the user in db if not exists, returning the id
		err = db.c.QueryRow(`INSERT INTO User(Username,Photo) VALUES(?,?) RETURNING UserId`, username, default_photo).Scan(&id)
		photo = default_photo
	} else {
		// take the id of the already existing user
		err = db.c.QueryRow(`SELECT UserId, Photo FROM User WHERE Username=?`, username).Scan(&id, &photo)
	}
	if err != nil {
		return 0, "", err
	}

	return id, photo, err
}

func (db *appdbimpl) GetIdFromUsername(username string) (int, error) {

	var userid int
	err := db.c.QueryRow(`SELECT UserId FROM User WHERE Username=?`, username).Scan(&userid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, err
	}
	return userid, err
}

func (db *appdbimpl) GetUsernameFromId(userid int) (string, error) {

	var username string
	err := db.c.QueryRow(`SELECT Username FROM User WHERE UserId=?`, userid).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return username, err
}

func (db *appdbimpl) ChangeUsername(userid int, username string) error {

	res, err := db.c.Exec("UPDATE User SET Username=? WHERE UserId=?", username, userid)
	if err != nil {
		if strings.Contains(err.Error(), sqlite3.ErrConstraintUnique.Error()) {
			return ErrUsernameAlreadyExists
		}
		return err
	}

	// check if the row effected are 0 which mean the user don't exists
	eff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if eff == 0 {
		return ErrUserNotFound
	}

	return err
}

func (db *appdbimpl) ChangeUserPhoto(userid int, photo string) error {

	res, err := db.c.Exec("UPDATE User SET Photo=? WHERE UserId=?", photo, userid)
	if err != nil {
		return err
	}

	// check if the row effected are 0 which mean the user don't exists
	eff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if eff == 0 {
		return ErrUserNotFound
	}

	return err
}

func (db *appdbimpl) IsUserInChat(chatid int, userid int) (bool, error) {

	userinchat := false
	err := db.c.QueryRow(`SELECT EXISTS(SELECT * FROM ChatUser WHERE ChatId=? AND UserId=?)`, chatid, userid).Scan(&userinchat)
	if err != nil {
		return false, err
	}

	return userinchat, err
}

func (db *appdbimpl) SetLastAccess(userid int) error {

	_, err := db.c.Exec(`UPDATE User SET LastAccess=CURRENT_TIMESTAMP WHERE UserId=?`, userid)
	if err != nil {
		return err
	}

	return err
}

func (db *appdbimpl) SetLastRead(userid int, chatid int) error {

	_, err := db.c.Exec(`UPDATE ChatUser SET LastRead=CURRENT_TIMESTAMP WHERE UserId=? AND ChatId=?`, userid, chatid)
	if err != nil {
		return err
	}

	return err
}

func (db *appdbimpl) SearchUsers(username string) ([]components.User, error) {

	userlist := []components.User{}
	userrows, err := db.c.Query(`SELECT * FROM User WHERE Username LIKE ?||'%'`, username)
	if err != nil {
		return userlist, err
	}

	defer userrows.Close()

	// cicle for all the users
	for userrows.Next() {
		var user components.User
		err = userrows.Scan(&user.UserId, &user.Username, &user.Photo, &user.LastAccess)
		if err != nil {
			return userlist, err
		}

		// append the user to userlist
		userlist = append(userlist, user)
	}

	if userrows.Err() != nil {
		return userlist, err
	}

	return userlist, err
}
