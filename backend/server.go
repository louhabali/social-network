package main

import (
	"net/http"

	"main/pkg/app/handlers"
	"main/pkg/app/utils"
	database "main/pkg/db/sqlite"
)

func main() {
	db := database.CreateTables()
	utils.CreateDirectories()
	defer db.Close()
	handler := handlers.NewHandler(db)

	http.HandleFunc("/api/login", handlers.HandleCORS(handler.LoginHandler))
	http.HandleFunc("/api/register", handlers.HandleCORS(handler.RegisterHandler))
	http.HandleFunc("/api/posts", handlers.HandleCORS(handler.Midd(handler.ShowPosts)))
	http.HandleFunc("/api/posting", handlers.HandleCORS(handler.Midd(handler.PostCreationHandler)))
	http.HandleFunc("/api/searchData", handlers.HandleCORS(handler.HandleSearch))
	http.HandleFunc("/api/suggestions", handlers.HandleCORS(handler.Midd(handler.Suggestions)))

	http.HandleFunc("/api/userdata", handlers.HandleCORS(handler.Midd(handler.UserData)))
	http.HandleFunc("/api/profile/{userId}", handlers.HandleCORS(handler.Midd(handler.GetProfile)))
	http.HandleFunc("/api/profile/{userId}/privacy", handlers.HandleCORS(handler.Midd(handler.TogglePrivacy)))
	http.HandleFunc("/api/profile/{userId}/follow", handlers.HandleCORS(handler.Midd(handler.FollowUser)))

	http.HandleFunc("/api/posts/{postId}/comments", handlers.HandleCORS(handler.Midd(handler.CommentCreationHandler)))
	http.HandleFunc("/api/posts/{postId}", handlers.HandleCORS(handler.Midd(handler.ShowPost)))

	http.HandleFunc("/api/like", handlers.HandleCORS(handler.Midd(handler.LikeHandler)))
	http.HandleFunc("/api/logout", handlers.HandleCORS(handler.Midd(handler.LogoutHandler)))

	http.HandleFunc("/api/groups/{groupId}", handlers.HandleCORS(handler.Midd(handler.GetGroupByIdHandler)))
	http.HandleFunc("/api/groups", handlers.HandleCORS(handler.Midd(handler.CreateGroupHandler)))
	http.HandleFunc("/api/groups/{groupId}/join", handlers.HandleCORS(handler.Midd(handler.JoinGroupHandler)))
	http.HandleFunc("/api/group-join/reject", handlers.HandleCORS(handler.Midd(handler.JoinReject)))
	http.HandleFunc("/api/group-join/accept", handlers.HandleCORS(handler.Midd(handler.JoinAccept)))
	http.HandleFunc("/api/groups/{groupId}/events", handlers.HandleCORS(handler.Midd(handler.GroupEventsHandler)))
	http.HandleFunc("/api/groups/{groupId}/events/going", handlers.HandleCORS(handler.Midd(handler.CreateEventResponse)))
	http.HandleFunc("/api/groups/{groupId}/invite", handlers.HandleCORS(handler.Midd(handler.InviteUsersGroup)))

	http.HandleFunc("/api/images/", handlers.HandleCORS(handler.ServeImage))
	// http.HandleFunc("/api/getallusers", handlers.HandleCORS(handler.Midd(handler.GetAllusers)))
	http.HandleFunc("/api/getfollowedusers", handlers.HandleCORS(handler.Midd(handler.GetAllusers)))

	http.HandleFunc("/api/ws", handlers.HandleCORS(handler.Midd(handler.HandleWebSocket)))
	http.HandleFunc("/api/notifications", handlers.HandleCORS(handler.Midd(handler.GetNotifications)))
	http.HandleFunc("/api/notifications/read", handlers.HandleCORS(handler.Midd(handler.MarkNotifAsRead)))
	http.HandleFunc("/api/notifications/count", handlers.HandleCORS(handler.Midd(handler.GetUnreadCount)))
	http.HandleFunc("/api/conversations", handlers.HandleCORS(handler.Midd(handler.GetConversations)))
	http.HandleFunc("/api/messages", handlers.HandleCORS(handler.Midd(handler.GetMessages)))
	http.HandleFunc("/api/messages/read", handlers.HandleCORS(handler.Midd(handler.MarkMsgAsRead)))
	http.HandleFunc("/api/follow/accept", handlers.HandleCORS(handler.Midd(handler.AcceptFollowReq)))
	http.HandleFunc("/api/follow/reject", handlers.HandleCORS(handler.Midd(handler.RejectFollowReq)))

	http.HandleFunc("/api/validate-token", handlers.HandleCORS(handler.ValidateToken))

	http.ListenAndServe(":8080", nil)
}
