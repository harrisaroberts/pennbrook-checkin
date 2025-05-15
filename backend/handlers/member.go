package handlers

import (
	"backend/db"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type Member struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Age          int    `json:"age"`
	MemberType   string `json:"member_type"`
	MembershipID int    `json:"membership_id"`
}

type MemberDetail struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Age              int    `json:"age"`
	MemberType       string `json:"member_type"`
	SwimTestPassed   bool   `json:"swim_test_passed"`
	ParentNoteOnFile bool   `json:"parent_note_on_file"`
	IsCheckedIn      bool   `json:"is_checked_in"`
    MembershipID     int    `json:"membership_id"`
}

// SearchMembers returns all members or filters by name
func SearchMembers(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("name"))

	var rows *sql.Rows
	var err error

	if query == "" {
		rows, err = db.DB.Query(`SELECT id, name, age, member_type, membership_id FROM members`)
	} else {
		rows, err = db.DB.Query(`
			SELECT id, name, age, member_type, membership_id
			FROM members
			WHERE LOWER(name) LIKE '%' || LOWER($1) || '%'`, query)
	}

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

// GetFamilyByMemberID returns all members in the same membership
func GetFamilyByMemberID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing member ID", http.StatusBadRequest)
		return
	}

	var membershipID int
	err := db.DB.QueryRow(`
		SELECT membership_id FROM members WHERE id = $1
	`, idStr).Scan(&membershipID)
	if err != nil {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	rows, err := db.DB.Query(`
		SELECT 
			id, 
			name, 
			age, 
			member_type, 
			swim_test_passed, 
			parent_note_on_file,
            membership_id,
			EXISTS (
				SELECT 1 
				FROM checkins 
				WHERE checkins.member_id = members.id 
				AND checkin_date = CURRENT_DATE
			) AS is_checked_in
		FROM members
		WHERE membership_id = $1
		ORDER BY id
	`, membershipID)
	if err != nil {
		http.Error(w, "Failed to query family", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var family []MemberDetail
	for rows.Next() {
		var m MemberDetail
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Age,
			&m.MemberType,
			&m.SwimTestPassed,
			&m.ParentNoteOnFile,
            &m.MembershipID,
			&m.IsCheckedIn,
		); err != nil {
			http.Error(w, "Failed to scan member", http.StatusInternalServerError)
			return
		}
		family = append(family, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(family)
}

