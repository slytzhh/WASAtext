package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// post /session
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Take the username from the request body
	var username components.Username
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// check if the username has the correct length
	if len(username.Username) < 3 || len(username.Username) > 18 {
		http.Error(w, database.ErrUsernameLength.Error(), http.StatusBadRequest) // 400
		return
	}
	// Insert the user in the database
	var id int
	var photo string
	id, photo, err = rt.db.InsertUser(username.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// create the user info
	var userinfo components.UserIdPhoto
	userinfo.UserId = id
	userinfo.Photo = photo

	// set the header of the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(userinfo)
}
