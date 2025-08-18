package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// put /users/{user_id}/name
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// Take the username from the request body
	var newusername components.Username
	err = json.NewDecoder(r.Body).Decode(&newusername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the user id from the URL
	userid, err := strconv.Atoi(ps.ByName("user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// check if the user performing the action is the same of the user changing username
	if userperformingid != userid {
		http.Error(w, database.ErrDifferentUser.Error(), http.StatusUnauthorized) // 401
		return
	}

	// check if the username has the correct length
	if len(newusername.Username) < 3 || len(newusername.Username) > 18 {
		http.Error(w, database.ErrUsernameLength.Error(), http.StatusBadRequest) // 400
		return
	}

	// change the user username
	err = rt.db.ChangeUsername(userid, newusername.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.WriteHeader(http.StatusNoContent) // 204
}

// put /users/{user_id}/photo
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// Take the user photo from the request body
	var newphoto components.Photo
	err = json.NewDecoder(r.Body).Decode(&newphoto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the user id from the URL
	userid, err := strconv.Atoi(ps.ByName("user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// check if the user performing the action is the same of the user changing photo
	if userperformingid != userid {
		http.Error(w, database.ErrDifferentUser.Error(), http.StatusUnauthorized) // 401
		return
	}

	// if the photo is empty sets the dafult photo
	if len(newphoto.Photo) == 0 {
		newphoto.Photo = "photo.png"
	}

	// change the user photo
	err = rt.db.ChangeUserPhoto(userid, newphoto.Photo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.WriteHeader(http.StatusNoContent) // 204
}

// get /users
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Take the username to search from the query
	username := r.URL.Query().Get("username")

	// gets the users with username that starts equal to the username in input
	var userlist components.UserList
	var err error
	userlist.UserList, err = rt.db.SearchUsers(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	_ = json.NewEncoder(w).Encode(userlist)
}
