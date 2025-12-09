package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/f044fs3t5w3f/gophermart/internal/service"
)

func listOrders(s *service.Service) http.HandlerFunc {
	h := func(w http.ResponseWriter, r *http.Request) {
		orders, err := s.ListOrders(r.Context())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		ordersResponse := make([]orderResponse, 0, len(orders))
		for _, order := range orders {
			ordersResponse = append(ordersResponse, orderResponse{
				Number:     order.Number,
				Status:     order.Status,
				Accrual:    9000, // TODO: calculate accrual
				UploadedAt: order.UploadedAt.Format(time.RFC3339),
			})
		}
		jsonEncoder := json.NewEncoder(w)
		err = jsonEncoder.Encode(ordersResponse)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	}
	return h
}
