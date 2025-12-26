package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/f044fs3t5w3f/gophermart/internal/service"
)

type withdrawResponse struct {
	Order       string  `json:"withdraw"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func ListWithdraws(s *service.Service) http.HandlerFunc {
	h := func(w http.ResponseWriter, r *http.Request) {
		withdraws, err := s.ListWithdraws(r.Context())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if len(withdraws) == 0 {
			http.Error(w, "[]", http.StatusNoContent)
		}
		withdrawsResponse := make([]withdrawResponse, 0, len(withdraws))
		for _, withdraw := range withdraws {
			withdrawsResponse = append(withdrawsResponse, withdrawResponse{
				Order:       withdraw.Order,
				Sum:         withdraw.Sum,
				ProcessedAt: withdraw.ProcessedAt.Format(time.RFC3339),
			})
		}
		jsonEncoder := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		err = jsonEncoder.Encode(withdrawsResponse)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	}
	return h
}
