package handler

import (
	"github.com/f044fs3t5w3f/gophermart/internal/compress"
	"github.com/f044fs3t5w3f/gophermart/internal/service"
	"github.com/go-chi/chi/v5"
)

func GetRouter(s *service.Service) *chi.Mux {
	r := chi.NewRouter()
	r.Use(compress.Middleware)
	r.Post("/api/user/register", register(s))
	r.Post("/api/user/login", login(s))
	return r
}
