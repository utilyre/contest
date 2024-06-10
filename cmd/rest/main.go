package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/utilyre/contest/internal/adapters/handler"
	"github.com/utilyre/contest/internal/adapters/postgres"
	"github.com/utilyre/contest/internal/app/service"
)

func main() {
	db, err := postgres.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	accountRepo := postgres.NewAccountRepo(db)
	accountSvc := service.NewAccountService(accountRepo)

	r := chi.NewRouter()
	v1 := chi.NewRouter()

	v1.Route("/auth", func(r chi.Router) {
		handler := handler.NewAuthHandler(accountSvc)

		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})

	r.Get("/health", handler.NewHealthHandler().Check)
	r.Mount("/api/v1", v1)

	log.Println("Listening on http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
