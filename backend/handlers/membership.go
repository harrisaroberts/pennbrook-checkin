package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"backend/db"
)

type FamilyMember struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Age              int    `json:"age"`
	MemberType       string `json:"member_type"`
	SwimTestPassed   bool   `json:"swim_test_passed"`
	ParentNoteOnFile bool   `json:"parent_note_on_file"`
    IsCheckedIn      bool   `json:"is_checked_in"`
}

type MembershipWithMembers struct {
	MembershipID   int            `json:"membership_id"`
	MembershipType string         `json:"membership_type"`
	Status         string         `json:"status"`
	Members        []FamilyMember `json:"members"`
}

func GetMembershipsByLastName(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if name == "" {
		http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
		return
	}

	query := `
		SELECT m.id, m.membership_type, m.status,
		       mem.id, mem.name, mem.age, mem.member_type, mem.swim_test_passed, mem.parent_note_on_file
		FROM memberships m
		JOIN members mem ON m.id = mem.membership_id
        WHERE LOWER(split_part(mem.name, ' ', -1)) LIKE LOWER($1) || '%'
		ORDER BY m.id, mem.id
	`

	rows, err := db.DB.Query(query, name)
	if err != nil {
		log.Println("DB query failed:", err)
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	families := make(map[int]*MembershipWithMembers)

	for rows.Next() {
		var (
			membershipID   int
			membershipType string
			status         string
			member         FamilyMember
		)

		if err := rows.Scan(
			&membershipID,
			&membershipType,
			&status,
			&member.ID,
			&member.Name,
			&member.Age,
			&member.MemberType,
			&member.SwimTestPassed,
			&member.ParentNoteOnFile,
		); err != nil {
			log.Println("Row scan error:", err)
			continue
		}

		if _, ok := families[membershipID]; !ok {
			families[membershipID] = &MembershipWithMembers{
				MembershipID:   membershipID,
				MembershipType: membershipType,
				Status:         status,
				Members:        []FamilyMember{},
			}
		}

		families[membershipID].Members = append(families[membershipID].Members, member)
	}

	// Flatten the map to a slice
	var result []MembershipWithMembers
	for _, fam := range families {
		result = append(result, *fam)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
