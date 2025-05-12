package handlers

import (
	"encoding/json"
	"net/http"
	"backend/db"
)

type Member struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Age          int    `json:"age"`
	MemberType   string `json:"member_type"`
	MembershipID int    `json:"membership_id"`
}

func SearchMembers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("name")
	if query == "" {
		http.Error(w, "Missing 'name' query param", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
        SELECT id, name, age, member_type, membership_id
        FROM members
        WHERE LOWER(name) LIKE '%' || LOWER($1) || '%'
    `, query)
	if err != nil {
		http.Error(w, "Failed to search members", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []Member
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.ID, &m.Name, &m.Age, &m.MemberType, &m.MembershipID); err != nil {
			http.Error(w, "Error parsing result", http.StatusInternalServerError)
			return
		}
		results = append(results, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

