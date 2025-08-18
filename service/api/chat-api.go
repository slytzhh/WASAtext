package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// post /chats/newchat
func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the chat to create from the request body
	var chat components.ChatCreation
	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// check if the are at least 2 person in chat
	if len(chat.UsernameList) < 2 {
		http.Error(w, database.ErrLessTwoUserInChat.Error(), http.StatusBadRequest) // 400
		return
	}

	// if there are more than 2 users check that the group name is not empty
	if len(chat.UsernameList) > 2 && len(chat.GroupName) == 0 {
		http.Error(w, database.ErrGroupNameLength.Error(), http.StatusBadRequest) // 400
		return
	}

	// check if is a chat and the message is not empty, unless is forwarded
	if len(chat.GroupName) == 0 && len(chat.FirstMessage.Text) == 0 && len(chat.FirstMessage.Photo) == 0 && chat.ForwardedId == 0 {
		http.Error(w, database.ErrMessageEmpty.Error(), http.StatusBadRequest) // 400
		return
	}

	// Inserts the chat in the database
	var chatid, messageid int
	chatid, messageid, err = rt.db.InsertChat(chat, userperformingid)
	if errors.Is(err, database.ErrTransaction) {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	var ids components.ChatMessId
	ids.ChatId = chatid
	ids.MessageId = messageid

	// set the header of the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(ids)
}

// put /chats/{chat_id}/users
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// take the username list from the request body
	var usernamelist components.UsernameList
	err = json.NewDecoder(r.Body).Decode(&usernamelist)
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

	// check if the chat is a group
	isgroup, err := rt.db.IsGroup(chatid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}
	if !isgroup {
		http.Error(w, database.ErrOperationGroupOnly.Error(), http.StatusBadRequest) // 400
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

	// insert the usernames in the db
	err = rt.db.AddUsersToGroup(usernamelist.UsernameList, chatid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // 204
	_ = json.NewEncoder(w).Encode(usernamelist)
}

// get /chats/{chat_id}/users
func (rt *_router) getChatUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// check if the chat is a group
	isgroup, err := rt.db.IsGroup(chatid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}
	if !isgroup {
		http.Error(w, database.ErrOperationGroupOnly.Error(), http.StatusBadRequest) // 400
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

	// get the users in the chat
	var usersList []components.Username
	usersList, err = rt.db.GetUsersInChat(chatid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	_ = json.NewEncoder(w).Encode(usersList)
}

// put /chats/{chat_id}/users/{user_id}
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// take the user id from the URL
	userid, err := strconv.Atoi(ps.ByName("user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// check if the chat is a group
	isgroup, err := rt.db.IsGroup(chatid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}
	if !isgroup {
		http.Error(w, database.ErrOperationGroupOnly.Error(), http.StatusBadRequest) // 400
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

	// check if the user performing the action is the same ad the one leaving the group
	if userperformingid != userid {
		http.Error(w, database.ErrDifferentUser.Error(), http.StatusUnauthorized) // 401
		return
	}

	// deletes the user from the group
	err = rt.db.DeleteUserFromGroup(userid, chatid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.WriteHeader(http.StatusNoContent) // 204
}

// put /chats/{chat_id}/name
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// Take the groupname from the request body
	var newgroupname components.GroupName
	err = json.NewDecoder(r.Body).Decode(&newgroupname)
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

	// check if the group name is empty
	if len(newgroupname.GroupName) == 0 {
		http.Error(w, database.ErrGroupNameLength.Error(), http.StatusBadRequest) // 400
		return
	}

	// check if the chat is a group
	isgroup, err := rt.db.IsGroup(chatid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}
	if !isgroup {
		http.Error(w, database.ErrOperationGroupOnly.Error(), http.StatusBadRequest) // 400
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

	// change the group name
	err = rt.db.ChangeGroupName(chatid, newgroupname.GroupName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.WriteHeader(http.StatusNoContent) // 204
}

// put /users/{chat_id}/photo
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// Take the group photo from the request body
	var newphoto components.Photo
	err = json.NewDecoder(r.Body).Decode(&newphoto)
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

	// check if the chat is a group
	isgroup, err := rt.db.IsGroup(chatid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}
	if !isgroup {
		http.Error(w, database.ErrOperationGroupOnly.Error(), http.StatusBadRequest) // 400
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

	// if the photo is empty sets the dafult photo
	if len(newphoto.Photo) == 0 {
		newphoto.Photo = "photo.png"
	}

	// change the group photo
	err = rt.db.ChangeGroupPhoto(chatid, newphoto.Photo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.WriteHeader(http.StatusNoContent) // 204
}

// get /chats
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// takes all the chat of the user from db
	var chats []components.ChatPreview
	chats, err = rt.db.GetUserChats(userperformingid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// order the chats based on timestamp of the last message
	sort.Slice(chats, func(i, j int) bool {
		if chats[i].LastMessage.MessageId != 0 && chats[j].LastMessage.MessageId != 0 {
			return chats[i].LastMessage.TimeStamp > chats[j].LastMessage.TimeStamp
		} else if chats[i].LastMessage.MessageId == 0 && chats[j].LastMessage.MessageId != 0 {
			return chats[i].TimeCreated > chats[j].LastMessage.TimeStamp
		} else if chats[i].LastMessage.MessageId != 0 && chats[j].LastMessage.MessageId == 0 {
			return chats[i].LastMessage.TimeStamp > chats[j].TimeCreated
		} else {
			return chats[i].TimeCreated > chats[j].TimeCreated
		}
	})

	// sets the LastAccess for the user
	err = rt.db.SetLastAccess(userperformingid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	_ = json.NewEncoder(w).Encode(chats)

}

// get /chats/{chat_id}
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// takes the chat from the db
	var chat components.Chat
	chat, err = rt.db.GetChat(chatid, userperformingid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// sets the LastRead of the chat for the user
	err = rt.db.SetLastRead(userperformingid, chatid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// set the header of the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	_ = json.NewEncoder(w).Encode(chat)

}
