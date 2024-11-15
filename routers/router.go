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
	categoryService := services.NewCategoryService(categoryRepo)
	CategoryHandler := handlers.NewCategoryHandler(categoryService)

	itemRepo := repositories.NewItemRepository(db)
	itemService := services.NewItemService(itemRepo)
	itemHandler := handlers.NewItemHandler(itemService)

	itemInvesmentRepo := repositories.NewItemInvestmentRepository(db)
	itemInvesmentService := services.NewItemInvestmentService(itemInvesmentRepo)
	itemInvesmentHandler := handlers.NewItemInvestmentHandler(*itemInvesmentService)

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

		r.Route("/items", func(r chi.Router) {
			r.With(middlewares.AuthMiddleware).Post("/", itemHandler.CreateItemHandler)
			r.With(middlewares.AuthMiddleware).Get("/{id}", itemHandler.GetItemByIDHandler)
			r.With(middlewares.AuthMiddleware).Put("/{id}", itemHandler.UpdateItemHandler)
			r.With(middlewares.AuthMiddleware).Delete("/{id}", itemHandler.DeleteItemHandler)
			r.With(middlewares.AuthMiddleware).Get("/", itemHandler.GetAllItemsHandler)
			r.With(middlewares.AuthMiddleware).Get("/need-replacement", itemHandler.GetReplacementItemsHandler)

			r.Route("/investment", func(r chi.Router) {
				r.With(middlewares.AuthMiddleware).Get("/", itemInvesmentHandler.CountAllItemInvestmentsHandler)
				r.With(middlewares.AuthMiddleware).Get("/{id}", itemInvesmentHandler.GetItemInvesmentByItemIdHandler)
			})
		})

	})

	return r
}
