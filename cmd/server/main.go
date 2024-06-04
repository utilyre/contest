package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"github.com/utilyre/contest/internal/database"
	"golang.org/x/crypto/bcrypt"
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

type RegisterParams struct {
	Username string
	Email    string
	Password string
}

func (ah AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	var params RegisterParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to generate hash from password:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := ah.queries.CreateAccount(r.Context(), database.CreateAccountParams{
		Username: params.Username,
		Email:    params.Email,
		Password: hashedPassword,
	}); err != nil {
		log.Println("failed to create account:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Account Created")
}

type LoginParams struct {
	Username string
	Password string
}

func (ah AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	var params LoginParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	account, err := ah.queries.GetAccountEmailPassword(r.Context(), params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Account Not Found", http.StatusNotFound)
			return
		}

		log.Println("failed to query database for account password:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword(account.Password, []byte(params.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			http.Error(w, "Account Not Found", http.StatusNotFound)
			return
		}

		log.Println("failed to query database for account password:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, AccountClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiryTime)),
		},
		Email: account.Email,
	}).SignedString([]byte("secret"))
	if err != nil {
		log.Println("failed to generate jwt token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"token": token})
}

const AccessTokenExpiryTime = 5 * time.Minute

type AccountClaims struct {
	jwt.RegisteredClaims
	Email string
}
