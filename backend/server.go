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

	http.HandleFunc("/api/login", handleCORS(handler.LoginHandler))
	http.HandleFunc("/api/register", handleCORS(handler.RegisterHandler))
	http.HandleFunc("/api/posts", handleCORS(handler.Midd(handler.ShowPosts)))
	http.HandleFunc("/api/posting", handleCORS(handler.Midd(handler.PostCreationHandler)))
	http.HandleFunc("/api/profile", handleCORS(handler.Midd(handler.GetProfile)))
	http.HandleFunc("/api/searchData", handleCORS(handler.HandleSearch))
	http.HandleFunc("/api/userdata", handleCORS(handler.Midd(handler.UserData)))
	http.HandleFunc("/api/images/", handleCORS(handler.ServeImage))
	http.HandleFunc("/api/like", handleCORS(handler.Midd(handler.LikeHandler)))
	http.HandleFunc("/api/logout", handleCORS(handler.Midd(handler.LogoutHandler)))
	
	http.ListenAndServe(":8080", nil)
}

func handleCORS(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

