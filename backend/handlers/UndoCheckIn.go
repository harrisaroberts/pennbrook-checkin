package handlers

import (
	"backend/db"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func UndoCheckIn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	memberID, err := strconv.Atoi(idStr)
	if err != nil || memberID == 0 {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`
		DELETE FROM checkins WHERE member_id = $1 AND checkin_date = CURRENT_DATE
	`, memberID)
	if err != nil {
		log.Println("Failed to undo check-in:", err)
		http.Error(w, "Failed to undo check-in", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
