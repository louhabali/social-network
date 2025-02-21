package handlers

import (
	"net/http"
)

func (h *Handler) ServeImage(w http.ResponseWriter, r *http.Request) {
	imageID := r.URL.Path[len("/api/images/"):]
	http.ServeFile(w, r, imageID)
}
