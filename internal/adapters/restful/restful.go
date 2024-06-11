package restful

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/utilyre/contest/internal/app/service"
)

func New(accountSvc *service.AccountService) *chi.Mux {
	mux := chi.NewMux()
	v1 := chi.NewRouter()

	v1.Route("/auth", func(r chi.Router) {
		handler := newAuthHandler(accountSvc)

		r.Post("/register", handler.register)
		r.Post("/login", handler.login)
	})

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.CleanPath)
	mux.Get("/health", newHealthHandler().check)
	mux.Mount("/api/v1", v1)

	return mux
}
