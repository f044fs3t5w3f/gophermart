package handler

import (
	"encoding/json"
	"net/http"

	"github.com/f044fs3t5w3f/gophermart/internal/service"
)

type balanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func balance(s *service.Service) http.HandlerFunc {
	h := func(w http.ResponseWriter, r *http.Request) {
		current, withdrawn, err := s.Balance(r.Context())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		response := &balanceResponse{
			Current:   current,
			Withdrawn: withdrawn,
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	return h
}
