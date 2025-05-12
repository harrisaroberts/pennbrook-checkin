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
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.ListenAndServe(":8080", middleware.EnableCORS(r))

	r.HandleFunc("/members", handlers.SearchMembers).Methods("GET")
	r.HandleFunc("/memberships", handlers.GetMemberships).Methods("GET")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

