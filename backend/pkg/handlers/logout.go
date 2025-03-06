package handlers

import "net/http"

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no token found"))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "Token",
		Value:  "",     
		Path:   "/", 
		MaxAge: -1,
	})

	_, err = h.db.Exec(`UPDATE tokens SET token=NULL WHERE token=?`, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to update the token"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
