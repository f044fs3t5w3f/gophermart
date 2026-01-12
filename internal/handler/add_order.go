package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/f044fs3t5w3f/gophermart/internal/service"
)

func addOrder(s *service.Service) http.HandlerFunc {
	h := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		orderNumber := strings.Trim(string(body), "\n")
		ctx := r.Context()
		created, err := s.AddOrder(ctx, "", orderNumber)
		switch err {
		case service.ErrInternalError:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		case service.ErrOrderAlreadyExists:
			http.Error(w, "Order already exists", http.StatusConflict)
		case service.ErrInvalidNumber:
			http.Error(w, "Invalid number", http.StatusUnprocessableEntity)
		case nil:
			if created {
				w.WriteHeader(http.StatusAccepted)
			}
		}

	}
	return h
}
