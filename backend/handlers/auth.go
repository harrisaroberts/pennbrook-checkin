package handlers

import (
	"backend/db"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
    "strings"
    "os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt")

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) 

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var creds LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var (
		id           int
		passwordHash string
		role         string
	)

	err := db.DB.QueryRow(
		"SELECT id, password_hash, role FROM users WHERE username = $1",
		creds.Username,
	).Scan(&id, &passwordHash, &role)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"username": creds.Username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Could not sign token", http.StatusInternalServerError)
		return
	}

	// Return token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
        "role": role,
	})
}

func GetUserRoleFromJWT(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", http.ErrNoCookie
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", http.ErrNoCookie
	}

	tokenStr := parts[1]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Optional: validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrInvalidKey
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", jwt.ErrInvalidKey
	}

	return role, nil
}

func isGuardOrAdmin(r *http.Request) bool {
	role, err := GetUserRoleFromJWT(r)
	if err != nil {
		return false
	}
	return role == "guard" || role == "admin"
}

