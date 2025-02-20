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
	http.HandleFunc("/api/posts", handleCORS(handler.ShowPosts))
	http.HandleFunc("/api/posting", handleCORS(handler.ShowPosts))
	http.HandleFunc("/api/profile", handleCORS(handler.GetProfile))

	http.ListenAndServe(":8080", nil)
}

func handleCORS(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
