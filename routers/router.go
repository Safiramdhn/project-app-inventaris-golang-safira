package routers

import (
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/database"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/handlers"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/middlewares"
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

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(*categoryRepo)
	CategoryHandler := handlers.NewCategoryHandler(categoryService)

	// Initialize router
	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", AuthHandler.RegisterHandler) // Register new user using POST /api/auth/register
			r.Post("/login", AuthHandler.LoginHandler)       // Login user using POST /api/auth/login
		})

		r.Route("/categories", func(r chi.Router) {
			r.With(middlewares.AuthMiddleware).Post("/", CategoryHandler.CreateCategoryHandler)
			r.With(middlewares.AuthMiddleware).Put("/{id}", CategoryHandler.UpdateCategoryHandler)
			r.With(middlewares.AuthMiddleware).Delete("/{id}", CategoryHandler.DeleteCategoryHandler)
			r.With(middlewares.AuthMiddleware).Get("/", CategoryHandler.GetCategoriesHandler)
			r.With(middlewares.AuthMiddleware).Get("/{id}", CategoryHandler.GetCategoryByIDHandler)
		})

	})

	return r
}
