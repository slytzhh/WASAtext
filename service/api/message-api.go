package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// post /chats/{chat_id}/messages
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the message to send from the request body
	var newmessage components.MessageToSend
	err = json.NewDecoder(r.Body).Decode(&newmessage)
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

	// check if the user performing the action is in the chat
	userinchat, err := rt.db.IsUserInChat(chatid, userperformingid)
	if err != nil {
		if errors.Is(err, database.ErrChatNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest) // 400
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}
	if !userinchat {
		http.Error(w, database.ErrUserNotInChat.Error(), http.StatusUnauthorized) // 401
		return
	}

	// check if the message is empty
	if len(newmessage.Text) == 0 && len(newmessage.Photo) == 0 {
		http.Error(w, database.ErrMessageEmpty.Error(), http.StatusBadRequest) // 400
		return
	}

	// Inserts the message in the database
	var messageid int
	messageid, err = rt.db.InsertMessage(newmessage, false, chatid, userperformingid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	var id components.MessageId
	id.MessageId = messageid

	// set the header of the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(id)
}

// post /chats/{chat_id}/forwardedmessages
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the message id of the message to forward from the request body
	var messageid components.MessageId
	err = json.NewDecoder(r.Body).Decode(&messageid)
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

	// takes the info of the message to forward
	text, photo, err := rt.db.GetMessage(messageid.MessageId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// creates the new message
	var newmessage components.MessageToSend
	newmessage.Text = text
	newmessage.Photo = photo

	// Inserts the message in the database
	var newmessageid int
	newmessageid, err = rt.db.InsertMessage(newmessage, true, chatid, userperformingid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// creates the message id to return
	var id components.MessageId
	id.MessageId = newmessageid

	// set the header of the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(id)
}

// delete /chats/{chat_id}/messages/{message_id}
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// check if the user is the same that sent the message
	userid, err := rt.db.GetUserFromMessage(messageid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}
	if userid != userperformingid {
		http.Error(w, database.ErrMessNotSent.Error(), http.StatusUnauthorized) // 401
		return
	}

	// deletes the message if exists
	err = rt.db.DeleteMessage(messageid, chatid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.WriteHeader(http.StatusNoContent) // 204
}

// post /chats/{chat_id}/repliedmessages
func (rt *_router) replyMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the message from the request body
	var reply components.MessageToSend
	err = json.NewDecoder(r.Body).Decode(&reply)
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

	// Inserts the message in the database
	var newmessageid int
	newmessageid, err = rt.db.InsertMessage(reply, false, chatid, userperformingid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// creates the message id to return
	var id components.MessageId
	id.MessageId = newmessageid

	// set the header of the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(id)
}
