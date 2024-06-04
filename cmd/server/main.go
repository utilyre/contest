package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/utilyre/contest/internal/database"
)

func main() {
	db, err := sql.Open("postgres", "postgresql://admin:admin@localhost:5432/contest?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := database.New(db)

	r := chi.NewRouter()

	r.Get("/helloworld", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello world")
	})

	r.Mount("/api/v1", NewV1Router(queries))

	log.Println("Listening on http://localhost:3000")
	http.ListenAndServe(":3000", r)
}

func NewV1Router(queries *database.Queries) chi.Router {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		handler := AuthHandler{queries: queries}

		r.Post("/register", handler.register)
		r.Post("/login", handler.login)
	})

	return r
}

type AuthHandler struct {
	queries *database.Queries
}

func (ah AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	log.Println("bro got hit")

	account, err := ah.queries.GetAccount(r.Context(), "utilyre@gmail.com")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(account)
}

func (ah AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	log.Println("other bro")
}
