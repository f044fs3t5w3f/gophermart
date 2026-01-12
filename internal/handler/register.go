package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/f044fs3t5w3f/gophermart/internal/auth"
	"github.com/f044fs3t5w3f/gophermart/internal/service"
)

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func register(s *service.Service) http.HandlerFunc {
	h := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		registerRequest := &registerRequest{}
		err := decoder.Decode(registerRequest)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		token, err := s.Register(r.Context(), registerRequest.Login, registerRequest.Password)
		if err == nil {
			http.SetCookie(w, &http.Cookie{
				Name:  auth.TokenCookieName,
				Value: token,
				Path:  "/",
			})
			return
		}
		if errors.Is(err, service.ErrUserExists) {
			http.Error(w, "User exists", http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	return h
}
