package auth

import (
	"context"
	"net/http"

	"github.com/f044fs3t5w3f/gophermart/internal/models"
)

type authRepo interface {
	GetUserByToken(ctx context.Context, token string) (*models.User, error)
}

const TokenCookieName = "token"

type key int

const ContextUserKey key = iota

func GetAuthMiddleware(authRepo authRepo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			cookie, err := r.Cookie(TokenCookieName)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token := cookie.Value
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			user, err := authRepo.GetUserByToken(ctx, token)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, ContextUserKey, user)
			newRequest := r.WithContext(ctx)
			next.ServeHTTP(w, newRequest)
		}
		return http.HandlerFunc(fn)
	}
}
