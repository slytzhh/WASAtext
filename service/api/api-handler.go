package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Login route
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User routes
	rt.router.PUT("/users/:user_id/name", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/:user_id/photo", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/users", rt.wrap(rt.searchUsers))

	// Chats routes
	rt.router.POST("/newchat", rt.wrap(rt.createChat))
	rt.router.PUT("/chats/:chat_id/users", rt.wrap(rt.addToGroup))
	rt.router.GET("/chats/:chat_id/users", rt.wrap(rt.getChatUsers))
	rt.router.DELETE("/chats/:chat_id/users/:user_id", rt.wrap(rt.leaveGroup))
	rt.router.PUT("/chats/:chat_id/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/chats/:chat_id/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.GET("/chats", rt.wrap(rt.getMyConversations))
	rt.router.GET("/chats/:chat_id", rt.wrap(rt.getConversation))

	// Message routes
	rt.router.POST("/chats/:chat_id/messages", rt.wrap(rt.sendMessage))
	rt.router.POST("/chats/:chat_id/forwardedmessages", rt.wrap(rt.forwardMessage))
	rt.router.POST("/chats/:chat_id/repliedmessages", rt.wrap(rt.replyMessage))
	rt.router.DELETE("/chats/:chat_id/messages/:message_id", rt.wrap(rt.deleteMessage))

	// Comment routes
	rt.router.PUT("/chats/:chat_id/messages/:message_id/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/chats/:chat_id/messages/:message_id/comments", rt.wrap(rt.uncommentMessage))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
