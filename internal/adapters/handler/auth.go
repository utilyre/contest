package handler

import (
	"encoding/json"
	"net/http"

	"github.com/utilyre/contest/internal/app/service"
)

type AuthHandler struct {
	accountSvc *service.AccountService
}

func NewAuthHandler(accountSvc *service.AccountService) *AuthHandler {
	return &AuthHandler{
		accountSvc: accountSvc,
	}
}

type RegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResp struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	out, err := ah.accountSvc.Register(
		r.Context(),
		&service.AccountRegisterParams{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&RegisterResp{
		ID:       out.Account.ID,
		Username: out.Account.Username.String(),
		Email:    out.Account.Email.String(),
	})
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token string `json:"token"`
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var params LoginReq
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	out, err := ah.accountSvc.Login(r.Context(), &service.AccountLoginParams{
		Username: params.Username,
		Password: params.Password,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&LoginResp{
		Token: out.Token,
	})
}
