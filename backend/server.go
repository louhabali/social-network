package main

import (
	"net/http"

	database "main/pkg/db/sqlite"
	"main/pkg/handlers"
)

func main() {
	db := database.CreateTables()
	defer db.Close()
	handler := handlers.NewHandler(db)


	/*server := http.Server {
		handler: handleCORS(mux),
		addr: 8080,
	}
*/
	http.HandleFunc("/api/login", handleCORS(handler.LoginHandler))
	http.HandleFunc("/api/register", handleCORS(handler.RegisterHandler))
	http.HandleFunc("/api/posts", handleCORS(handler.Midd(handler.ShowPosts)))
	http.HandleFunc("/api/posting", handleCORS(handler.Midd(handler.PostCreationHandler)))
	http.HandleFunc("/api/searchData", handleCORS(handler.HandleSearch))
	http.HandleFunc("/api/suggestions", handleCORS(handler.Midd(handler.Suggestions)))

	http.HandleFunc("/api/userdata", handleCORS(handler.Midd(handler.UserData)))
	http.HandleFunc("/api/profile/{userId}", handleCORS(handler.Midd(handler.GetProfile)))
	http.HandleFunc("/api/profile/{userId}/privacy", handleCORS(handler.Midd(handler.TogglePrivacy)))
	http.HandleFunc("/api/profile/{userId}/follow", handleCORS(handler.Midd(handler.FollowUser)))

	http.HandleFunc("/api/posts/{postId}/comments", handleCORS(handler.Midd(handler.CommentCreationHandler)))
	http.HandleFunc("/api/posts/{postId}", handleCORS(handler.Midd(handler.ShowPost)))

	http.HandleFunc("/api/like", handleCORS(handler.Midd(handler.LikeHandler)))
	http.HandleFunc("/api/logout", handleCORS(handler.Midd(handler.LogoutHandler)))

	http.HandleFunc("/api/groups/{groupId}", handleCORS(handler.Midd(handler.GetGroupByIdHandler)))
	http.HandleFunc("/api/groups", handleCORS(handler.Midd(handler.CreateGroupHandler)))
	http.HandleFunc("/api/groups/{groupId}/join", handleCORS(handler.Midd(handler.JoinGroupHandler)))
	http.HandleFunc("/api/groups/{groupId}/events", handleCORS(handler.Midd(handler.GroupEventsHandler)))
	http.HandleFunc("/api/groups/{groupId}/events/going", handleCORS(handler.Midd(handler.CreateEventResponse)))
	http.HandleFunc("/api/groups/{groupId}/invite", handleCORS(handler.Midd(handler.InviteUsersGroup)))

	http.HandleFunc("/api/images/", handleCORS(handler.ServeImage))
	http.HandleFunc("/api/getallusers", handleCORS(handler.Midd(handler.GetAllusers)))
	http.HandleFunc("/api/getfollowedusers", handleCORS(handler.Midd(handler.GetAllusers)))

	http.HandleFunc("/api/ws", handleCORS(handler.Midd(handler.HandleWebSocket)))
	http.HandleFunc("/api/notifications", handleCORS(handler.Midd(handler.GetNotifications)))
	http.HandleFunc("/api/conversations", handleCORS(handler.Midd(handler.GetConversations)))
	http.HandleFunc("/api/messages", handleCORS(handler.Midd(handler.GetMessages)))

	http.ListenAndServe(":8080", nil)
}

func handleCORS(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
