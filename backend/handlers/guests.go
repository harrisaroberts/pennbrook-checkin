package handlers

import (
	"backend/db"
	"encoding/json"
	"net/http"
	"strings"
	"time"
    "github.com/gorilla/mux"
)

type AddGuestsRequest struct {
	MembershipID int    `json:"membership_id"`
	GuestNames   string `json:"guest_names"` // comma-separated
}

func AddGuests(w http.ResponseWriter, r *http.Request) {
	var req AddGuestsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	names := strings.Split(req.GuestNames, ",")
	visitDate := time.Now().Format("2006-01-02")

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}

    stmt, err := tx.Prepare("INSERT INTO guests (membership_id, guest_name, visit_date) VALUES ($1, $2, $3)")
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to prepare insert", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, rawName := range names {
		name := strings.TrimSpace(rawName)
		if name == "" {
			continue
		}
		if _, err := stmt.Exec(req.MembershipID, name, visitDate); err != nil {
			tx.Rollback()
			http.Error(w, "Failed to insert guest", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetMonthlyGuestCount(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing membership ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT COUNT(*) FROM guests
		WHERE membership_id = $1
		AND date_trunc('month', visit_date) = date_trunc('month', CURRENT_DATE)
	`

	var count int
	if err := db.DB.QueryRow(query, id).Scan(&count); err != nil {
		http.Error(w, "Failed to fetch guest count", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"guest_count": count})
}

func GetTodayGuestsByMembership(w http.ResponseWriter, r *http.Request) {
	membershipID := r.URL.Query().Get("membership_id")
	if membershipID == "" {
		http.Error(w, "Missing membership ID", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, guest_name, visit_date
		FROM guests
		WHERE membership_id = $1 AND visit_date = CURRENT_DATE
		ORDER BY id DESC
	`, membershipID)
	if err != nil {
		http.Error(w, "Failed to fetch guests", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Guest struct {
		ID        int    `json:"id"`
		GuestName string `json:"guest_name"`
		VisitDate string `json:"visit_date"`
	}

	var guests []Guest
	for rows.Next() {
		var g Guest
		if err := rows.Scan(&g.ID, &g.GuestName, &g.VisitDate); err != nil {
			http.Error(w, "Failed to scan guest", http.StatusInternalServerError)
			return
		}
		guests = append(guests, g)
	}

	json.NewEncoder(w).Encode(guests)
}

func DeleteGuest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.DB.Exec("DELETE FROM guests WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete guest", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

