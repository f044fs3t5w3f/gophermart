package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/f044fs3t5w3f/gophermart/internal/service"
)

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func login(s *service.Service) http.HandlerFunc {
	h := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		loginRequest := &loginRequest{}
		err := decoder.Decode(loginRequest)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		token, err := s.Login(r.Context(), loginRequest.Login, loginRequest.Password)
		fmt.Println(token)
		switch err {
		case service.ErrInternalError:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		case service.ErrUserDoesNotExists, service.ErrIncorrectPassword:
			http.Error(w, "Incorrect login/password", http.StatusUnauthorized)
		case nil:
			http.SetCookie(w, &http.Cookie{
				Name:  "token",
				Value: token,
				Path:  "/",
			})
		}

	}
	return h
}
