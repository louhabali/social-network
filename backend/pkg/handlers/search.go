package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type searchusers struct {
	Fname string
	Lname string
}

func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	searchvalue := r.FormValue("searchValue")
	searchq := `select first_name,last_name from users where first_name like CONCAT('%', ?, '%') OR last_name like CONCAT('%', ?, '%')  OR nickname like CONCAT('%', ?, '%') `
	
	rws, err := h.db.Query(searchq, searchvalue, searchvalue, searchvalue)
	if err != nil {
		fmt.Println("44", err)
		return
	}
	var Us []searchusers
	for rws.Next() {
		var U searchusers
		rws.Scan(&U.Fname, &U.Lname)
		
		Us = append(Us, U)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Us)
}
