package handler

import (
	"github.com/f044fs3t5w3f/gophermart/internal/auth"
	"github.com/f044fs3t5w3f/gophermart/internal/repository"
	"github.com/f044fs3t5w3f/gophermart/internal/service"
	"github.com/f044fs3t5w3f/gophermart/pkg/compress"
	"github.com/go-chi/chi/v5"
)

func GetRouter(s *service.Service, repo repository.Repository) *chi.Mux {
	r := chi.NewRouter()
	r.Use(compress.Middleware)
	r.Post("/api/user/register", register(s))
	r.Post("/api/user/login", login(s))

	authMiddleware := auth.GetAuthMiddleware(repo)
	authRequiredRoutes := chi.NewRouter()
	authRequiredRoutes.Use(authMiddleware)
	authRequiredRoutes.Post("/api/user/orders", addOrder(s))
	authRequiredRoutes.Get("/api/user/orders", listOrders(s))
	r.Mount("/", authRequiredRoutes)
	return r
}
