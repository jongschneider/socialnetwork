package main

import (
	"encoding/json"
	"net/http"
)

// LoginInput request body
type LoginInput struct {
	Email string `json:"email,omitempty"`
}

func login(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	email := input.Email

	// Find a user in the DB with the given email
	db.QueryRowContext(r.Context())
	// Issue JWT

	//Respond with the JWT
}
