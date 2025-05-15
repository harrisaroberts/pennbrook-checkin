package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"backend/db"
	"backend/handlers"
	"backend/middleware"
)

func main() {
	_ = godotenv.Load()
	db.Connect()

	r := mux.NewRouter()

	// All routes BEFORE ListenAndServe
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
    r.HandleFunc("/checkins", handlers.CheckInMember).Methods("POST")
    r.HandleFunc("/checkins/{id}", handlers.UndoCheckIn).Methods("DELETE")
    r.HandleFunc("/checkins/today", handlers.GetTodayCheckins).Methods("GET")
	r.HandleFunc("/members", handlers.SearchMembers).Methods("GET")
    r.HandleFunc("/memberships/by-last-name", handlers.GetMembershipsByLastName).Methods("GET")
    r.HandleFunc("/members/family", handlers.GetFamilyByMemberID).Methods("GET")
    r.HandleFunc("/guests", handlers.AddGuests).Methods("POST")
    r.HandleFunc("/guests/monthly-total", handlers.GetMonthlyGuestCount).Methods("GET")
    r.HandleFunc("/guests/today", handlers.GetTodayGuestsByMembership).Methods("GET")
    r.HandleFunc("/guests/{id}", handlers.DeleteGuest).Methods("DELETE")



	// Enable CORS and start the server
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.EnableCORS(r)))
}

