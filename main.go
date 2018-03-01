package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	// Database Connection
	databaseURL := env("DATABASE_URL", "postgresql://root@127.0.0.1:26257/socialnetwork?sslmode=disable")
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}

	// Router
	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Route("/api", func(api chi.Router) {
		jsonRequired := middleware.AllowContentType("application/json")
		api.With(jsonRequired).Post("/login", login)
		api.With(jsonRequired).Post("/users", createUser)
	})

	// Server
	port := env("PORT", "8080")
	log.Printf("Now serving on port: %s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, mux))

}

func env(key, fallbackValue string) string {
	val, present := os.LookupEnv(key)
	if present {
		return val
	}

	return fallbackValue
}

func respondError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func respondJSON(w http.ResponseWriter, v interface{}, code int) {
	b, err := json.Marshal(v)
	if err != nil {
		respondError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

// Middleware -
type Middleware func(http.Handler) http.Handler

func pipe(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
