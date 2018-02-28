package main

import (
	"database/sql"
	"fmt"
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
	mux.Route("/", func(api chi.Router) {
		jsonRequired := middleware.AllowContentType("application/json")
		api.With(jsonRequired).Post("/login", login)
	})

	// Server
	port := env("PORT", "80")
	fmt.Printf("Now serving on port: %s", port)
	log.Fatalln(http.ListenAndServe(":"+port, mux))

}

func env(key, fallbackValue string) string {
	val, present := os.LookupEnv(key)
	if present {
		return val
	}

	return fallbackValue
}
