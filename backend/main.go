package main

import (
	"backend/db"
	"backend/handlers"
	"backend/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	r.HandleFunc("/guests", handlers.AddGuests).Methods("POST")
	r.HandleFunc("/guests/monthly-total", handlers.GetMonthlyGuestCount).Methods("GET")
	r.HandleFunc("/guests/today", handlers.GetTodayGuestsByMembership).Methods("GET")
	r.HandleFunc("/guests/{id}", handlers.DeleteGuest).Methods("DELETE")
	r.HandleFunc("/login", handlers.HandleLogin).Methods("POST")

	// 👇 Register static route BEFORE dynamic member ID routes
	r.HandleFunc("/members/by-membership", handlers.GetMembersByMembershipID).Methods("GET")

	// Dynamic member routes
	r.HandleFunc("/members/{id}/swim-test", handlers.UpdateSwimTestStatus).Methods("PATCH")
	r.HandleFunc("/members/{id}", handlers.GetMemberByID).Methods("GET")
	r.HandleFunc("/members/{id}", handlers.UpdateMemberByID).Methods("PUT")

	// Enable CORS and start the server
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.EnableCORS(r)))
}

