package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"backend/db"
	"log"
)

type CheckinRequest struct {
	MemberID int `json:"member_id"`
}

func CheckInMember(w http.ResponseWriter, r *http.Request) {
	var req CheckinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.MemberID == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec(`
		INSERT INTO checkins (member_id, checkin_date)
		VALUES ($1, CURRENT_DATE)
		ON CONFLICT (member_id, checkin_date) DO NOTHING
	`, req.MemberID)

	if err != nil {
		log.Println("Failed to insert checkin:", err)
		http.Error(w, "Failed to check in member", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"status":     "checked_in",
		"member_id":  req.MemberID,
		"date":       time.Now().Format("2006-01-02"),
	})
}

func GetTodayCheckins(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT m.id, m.name, m.age, m.member_type
		FROM checkins c
		JOIN members m ON c.member_id = m.id
		WHERE c.checkin_date = CURRENT_DATE
		ORDER BY c.created_at
	`)
	if err != nil {
		http.Error(w, "Failed to fetch check-ins", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Member struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Age        int    `json:"age"`
		MemberType string `json:"member_type"`
	}

	var results []Member
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.ID, &m.Name, &m.Age, &m.MemberType); err != nil {
			http.Error(w, "Failed to parse check-ins", http.StatusInternalServerError)
			return
		}
		results = append(results, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

