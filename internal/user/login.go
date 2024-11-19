package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Serve the login HTML page for GET requests
			http.ServeFile(w, r, "static/login.html")
			return
		}

		if r.Method == http.MethodPost {
			var credentials struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}

			if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			user, err := GetUser(db, credentials.Username)
			if err != nil {
				log.Printf("Error getting user: %v", err)
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}

			if err := VerifyPassword(user.Password, credentials.Password); err != nil {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "token": "some-jwt-or-session-token"})
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
