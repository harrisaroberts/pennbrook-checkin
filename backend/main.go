package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/harrisaroberts/pennbrook-checkin/backend/db"
)

func main() {
	_ = godotenv.Load()

	db.Connect()

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

