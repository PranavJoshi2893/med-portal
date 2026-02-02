package server

import (
	"net/http"
	"time"

	"github.com/PranavJoshi2893/med-portal/internal/config"
	"github.com/PranavJoshi2893/med-portal/internal/handler"
	appMiddleware "github.com/PranavJoshi2893/med-portal/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routes(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, cfg *config.Config) http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(appMiddleware.AccessTokenMiddleware(cfg))
			r.Get("/", userHandler.GetAll)
			r.Delete("/{id}", userHandler.DeleteByID)
			r.Get("/{id}", userHandler.GetByID)
			r.Patch("/{id}", nil)
		})
	})

	return r
}
