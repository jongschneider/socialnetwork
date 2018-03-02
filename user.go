package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/lib/pq"
)

// CreateUserInput from request body
type CreateUserInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

// User model
type User struct {
	ID        string  `json:"-"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatarUrl"` //nullable
}

// Profile model
type Profile struct {
	Email           string    `json:"email,omitempty"`
	Username        string    `json:"username"`
	AvatarURL       *string   `json:"avatarUrl"`
	FollowersCount  int       `json:"followersCount"`
	FollowingCount  int       `json:"followingCount"`
	CreatedAt       time.Time `json:"createdAt"`
	Me              bool      `json:"me"`
	FollowersOfMine bool      `json:"followersOfMine"`
	FollowingOfMine bool      `json:"followingOfMine"`
}

const queryGetUser = "SELECT username, avatar_url FROM users WHERE is = $1"

func createUser(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	email := input.Email
	username := input.Username
	// TODO: validate input
	// Insert the user into the db
	var user Profile
	err := db.QueryRowContext(r.Context(), `
		INSERT INTO users (email, username) VALUES ($1, $2)
		RETURNING created_at
	`, email, username).Scan(&user.CreatedAt)
	if errPq, ok := err.(*pq.Error); ok && errPq.Code.Name() == "unique_violation" {
		if strings.Contains(errPq.Error(), "users_email_key") {
			respondJSON(w, map[string]string{
				"email": "Email taken",
			}, http.StatusUnprocessableEntity)
			return
		}
		if strings.Contains(errPq.Error(), "users_email_key") {
			respondJSON(w, map[string]string{
				"username": "Username taken",
			}, http.StatusUnprocessableEntity)
			return
		} else if err != nil {
			respondError(w, err)
			return
		}
		user.Email = email
		user.Username = username
		user.Me = true
	}
	// Respond with the created user
	respondJSON(w, user, http.StatusCreated)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUserID, authenticated := ctx.Value(keyAuthUserID).(string)
	username := chi.URLParam(r, "username")

	// Find user on the DB
	var userID string
	var user Profile
	if err := db.QueryRowContext(ctx, `
		SELECT
			id, 
			email,
			username, 
			avatar_url,
			followers_count,
			following_count,
			created_at
		FROM users
		WHERE username = $1
		`, username).Scan(
		&userID,
		&user.Email,
		&user.Username,
		&user.AvatarURL,
		&user.FollowersCount,
		&user.FollowingCount,
		&user.CreatedAt,
	); err == sql.ErrNoRows {
		http.Error(w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound)
		return
	} else if err != nil {
		respondError(w, err)
		return
	}

	// Hide details when the account isn't mine
	if !authenticated || authUserID != userID {
		user.Email = ""
	}

	// Respond with the found user profile
	respondJSON(w, user, http.StatusOK)
}
