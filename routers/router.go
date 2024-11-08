package routers

import (
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/database"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/handlers"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/repositories"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/services"
	"github.com/go-chi/chi/v5"
)

func NewRouter() chi.Router {
	db := database.NewPostgresDB()
	r := chi.NewRouter()

	// Initialize handlers
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	AuthHandler := handlers.NewAuthHandler(authService)

	// Initialize router
	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", AuthHandler.RegisterHandler) // Register new user using POST /api/auth/register
			r.Post("/login", AuthHandler.LoginHandler)       // Login user using POST /api/auth/login
		})
	})

	return r
}
