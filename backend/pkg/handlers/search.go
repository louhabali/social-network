package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type searchusers struct {
	Uid int
	Fname string
	Lname string
	Avatar *string
}

func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	searchvalue := r.FormValue("searchValue")
	searchq := `select id, first_name,last_name,avatar from users where first_name like CONCAT('%', ?, '%') OR last_name like CONCAT('%', ?, '%')  OR nickname like CONCAT('%', ?, '%') `
	
	rws, err := h.db.Query(searchq, searchvalue, searchvalue, searchvalue)
	if err != nil {
		fmt.Println("44", err)
		return
	}
	var Us []searchusers
	for rws.Next() {
		var U searchusers
		rws.Scan(&U.Uid ,&U.Fname, &U.Lname, &U.Avatar)
		if U.Avatar != nil {
			*U.Avatar,_  = EncodeImg(*U.Avatar)
		}
		
		Us = append(Us, U)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Us)
}
