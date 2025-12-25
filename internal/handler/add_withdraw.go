package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/f044fs3t5w3f/gophermart/internal/repository"
	"github.com/f044fs3t5w3f/gophermart/internal/service"
)

type addWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func addWithdraw(s *service.Service) http.HandlerFunc {
	h := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		addWithdrawRequest := &addWithdrawRequest{}
		err := decoder.Decode(addWithdrawRequest)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = s.AddWithdraw(r.Context(), addWithdrawRequest.Order, addWithdrawRequest.Sum)
		if errors.Is(err, repository.ErrWithdrawAllreadyExistsForAnotherUser) {
			http.Error(w, "Withdraw exists", http.StatusConflict)
		} else if errors.Is(err, repository.ErrWithdrawAllreadyExists) {
			return
		} else if errors.Is(err, repository.ErrNotEnough) {
			http.Error(w, "No enough points", http.StatusPaymentRequired)
		} else if errors.Is(err, service.ErrInvalidNumber) {
			http.Error(w, "Invalid number", http.StatusUnprocessableEntity)
		} else if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	return h
}
