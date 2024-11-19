package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Serve the registration page
			http.ServeFile(w, r, "static/register.html")
		case http.MethodPost:
			// Handle user registration
			var user User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			if err := CreateUser(db, &user); err != nil {
				log.Printf("Error creating user: %v", err)
				http.Error(w, "Error creating user", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
