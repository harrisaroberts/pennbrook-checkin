package handlers

import (
	"encoding/json"
	"net/http"
	"log"

	"backend/db"
)

type Membership struct {
	ID                int    `json:"id"`
	MembershipType    string `json:"membership_type"`
	Status            string `json:"status"`
	QuickBooksID      string `json:"quickbooks_id"`
	GuestLimitMonthly int    `json:"guest_limit_monthly"`
	GuestLimitTotal   int    `json:"guest_limit_total"`
	Notes             string `json:"notes"`
}

func GetMemberships(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, membership_type, status, quickbooks_id, guest_limit_monthly, guest_limit_total, notes FROM memberships")
	if err != nil {
		log.Printf("Error querying memberships: %v", err)
		http.Error(w, "Failed to query memberships", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var memberships []Membership
	for rows.Next() {
		var m Membership
		if err := rows.Scan(&m.ID, &m.MembershipType, &m.Status, &m.QuickBooksID, &m.GuestLimitMonthly, &m.GuestLimitTotal, &m.Notes); err != nil {
			http.Error(w, "Failed to parse memberships", http.StatusInternalServerError)
			return
		}
		memberships = append(memberships, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(memberships)
}
