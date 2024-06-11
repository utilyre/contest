package restful

import (
	"encoding/json"
	"net/http"

	"github.com/utilyre/contest/internal/app/service"
)

type authHandler struct {
	accountSvc *service.AccountService
}

func newAuthHandler(accountSvc *service.AccountService) *authHandler {
	return &authHandler{
		accountSvc: accountSvc,
	}
}

type registerReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResp struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (ah *authHandler) register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
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
	_ = json.NewEncoder(w).Encode(&registerResp{
		ID:       out.Account.ID,
		Username: out.Account.Username.String(),
		Email:    out.Account.Email.String(),
	})
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResp struct {
	Token string `json:"token"`
}

func (ah *authHandler) login(w http.ResponseWriter, r *http.Request) {
	var params loginReq
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
	json.NewEncoder(w).Encode(&loginResp{
		Token: out.Token,
	})
}
