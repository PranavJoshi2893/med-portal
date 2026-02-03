package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/PranavJoshi2893/med-portal/internal/config"
	"github.com/PranavJoshi2893/med-portal/pkg/auth"
	"github.com/PranavJoshi2893/med-portal/pkg/responses"
)

func AccessTokenMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				responses.WriteError(w, responses.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Status:  "UNAUTHORIZED",
					Message: "Unauthorized",
				})
				return
			}

			token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
			if token == "" {
				responses.WriteError(w, responses.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Status:  "UNAUTHORIZED",
					Message: "Unauthorized",
				})
				return
			}

			claims, err := auth.VerifyAccessToken(cfg.AccessTokenKey, token)
			if err != nil {
				responses.WriteError(w, responses.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Status:  "UNAUTHORIZED",
					Message: "Unauthorized",
				})
				return
			}
			role := claims.Role
			if role == "" {
				role = "user"
			}
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "role", role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RefreshTokenMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("refresh_token")
			if err != nil {
				responses.WriteError(w, responses.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Status:  "UNAUTHORIZED",
					Message: "Unauthorized",
				})
				return
			}

			claims, err := auth.VerifyRefreshToken(cfg.RefreshTokenKey, token.Value)
			if err != nil {
				responses.WriteError(w, responses.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Status:  "UNAUTHORIZED",
					Message: "Unauthorized",
				})
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "role", claims.Role)
			ctx = context.WithValue(ctx, "refresh_token", token.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
