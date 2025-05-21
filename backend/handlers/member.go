package handlers

import (
	"backend/db"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Member struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Age              int    `json:"age"`
	MemberType       string `json:"member_type"`
	MembershipID     int    `json:"membership_id"`
	SwimTestPassed   bool   `json:"swim_test_passed"`
	ParentNoteOnFile bool   `json:"parent_note_on_file"`
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

func UpdateSwimTestStatus(w http.ResponseWriter, r *http.Request) {
	if !isGuardOrAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var input struct {
		SwimTestPassed bool `json:"swim_test_passed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec(
		"UPDATE members SET swim_test_passed = $1 WHERE id = $2",
		input.SwimTestPassed,
		id,
	)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetMemberByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var m Member
	err := db.DB.QueryRow(`
		SELECT id, name, age, member_type, swim_test_passed, parent_note_on_file, membership_id
		FROM members
		WHERE id = $1
	`, id).Scan(
		&m.ID,
		&m.Name,
		&m.Age,
		&m.MemberType,
		&m.SwimTestPassed,
		&m.ParentNoteOnFile,
		&m.MembershipID,
	)
	if err != nil {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func UpdateMemberByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var m Member
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec(`UPDATE members
		SET name = $1,
		    age = $2,
		    member_type = $3,
		    swim_test_passed = $4,
		    parent_note_on_file = $5
		WHERE id = $6`,
		m.Name, m.Age, m.MemberType, m.SwimTestPassed, m.ParentNoteOnFile, id)
	if err != nil {
		http.Error(w, "Failed to update member", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetMembersByMembershipID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing membership ID", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT 
			id, name, age, member_type, swim_test_passed, parent_note_on_file, membership_id,
			EXISTS (
				SELECT 1 FROM checkins 
				WHERE checkins.member_id = members.id 
				AND checkin_date = CURRENT_DATE
			) AS is_checked_in
		FROM members
		WHERE membership_id = $1
		ORDER BY id
	`, idStr)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []MemberDetail
	for rows.Next() {
		var m MemberDetail
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Age,
			&m.MemberType,
			&m.SwimTestPassed,
			&m.ParentNoteOnFile,
			&m.MembershipID,
			&m.IsCheckedIn,
		)
		if err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		members = append(members, m)
	}

	if len(members) == 0 {
		http.Error(w, "No members found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}
