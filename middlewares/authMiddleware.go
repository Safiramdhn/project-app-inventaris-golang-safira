package middlewares

import (
	"net/http"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/database"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/repositories"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/services"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/utils"
)

var JsonResp = &utils.JSONResponse{}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract and validate token from request header or cookie
		cookie, err := r.Cookie("token")
		if err != nil || cookie == nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
		}

		db := database.NewPostgresDB()
		authRepo := repositories.NewAuthRepository(db)
		authService := services.NewAuthService(authRepo)
		_, err = authService.GetSession(cookie.Value)
		if err != nil {
			JsonResp.SendError(w, http.StatusUnauthorized, "Invalid token", err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
