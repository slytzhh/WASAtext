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

// put /chats/{chat_id}/messages/{message_id}/comments
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the emoji to comment from the request body
	var emoji components.Emoji
	err = json.NewDecoder(r.Body).Decode(&emoji)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the chat id from the URL
	chatid, err := strconv.Atoi(ps.ByName("chat_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the message id from the URL
	messageid, err := strconv.Atoi(ps.ByName("message_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// check if the message is in the chat
	messageinchat, err := rt.db.IsMessageInChat(chatid, messageid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}
	if !messageinchat {
		http.Error(w, database.ErrMessNotInChat.Error(), http.StatusUnauthorized) // 401
		return
	}

	// check if the user performing the action is in the chat
	userinchat, err := rt.db.IsUserInChat(chatid, userperformingid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}
	if !userinchat {
		http.Error(w, database.ErrUserNotInChat.Error(), http.StatusUnauthorized) // 401
		return
	}

	// check the comment is a single character
	if len(emoji.Emoji) < 1 && len(emoji.Emoji) > 4 {
		http.Error(w, database.ErrCommentLength.Error(), http.StatusBadRequest) // 400
		return
	}

	// inserts the comment in the db
	err = rt.db.InsertComment(messageid, userperformingid, emoji.Emoji)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.WriteHeader(http.StatusNoContent) // 204
}

// delete /chats/{chat_id}/messages/{message_id}/comments/{user_id}
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the chat id from the URL
	chatid, err := strconv.Atoi(ps.ByName("chat_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the message id from the URL
	messageid, err := strconv.Atoi(ps.ByName("message_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// check if the message is in the chat
	messageinchat, err := rt.db.IsMessageInChat(chatid, messageid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}
	if !messageinchat {
		http.Error(w, database.ErrMessNotInChat.Error(), http.StatusUnauthorized) // 401
		return
	}

	// check if the user performing the action is in the chat
	userinchat, err := rt.db.IsUserInChat(chatid, userperformingid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}
	if !userinchat {
		http.Error(w, database.ErrUserNotInChat.Error(), http.StatusUnauthorized) // 401
		return
	}

	// deletes the comment if exists
	err = rt.db.DeleteComment(messageid, userperformingid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.WriteHeader(http.StatusNoContent) // 204
}
