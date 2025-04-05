package api

import (
	"net/http"
	"time"

	"github.com/MdSadiqMd/Broadcast-API/internal/api/handlers"
	appMiddleware "github.com/MdSadiqMd/Broadcast-API/internal/api/middleware"
	"github.com/MdSadiqMd/Broadcast-API/internal/services"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func Setup(r *chi.Mux, db *gorm.DB, jwtSecret string) {
	userService := services.NewUserService(db)
	compaignService := services.NewCampaignService(db)

	auth := appMiddleware.NewAuth(appMiddleware.AuthConfig{
		JWTSecret:     jwtSecret,
		TokenDuration: 24 * time.Hour,
		UserService:   userService,
	})

	authHandler := handlers.NewAuthHandler(userService, auth)
	compaignHandler := handlers.NewCampaignHandler(compaignService, auth)

	r.Use(auth.Middleware())

	r.Route("/api", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware())
			r.Post("/campaign", compaignHandler.CreateCampaign)
			r.Get("/campaigns", compaignHandler.GetAllCampaigns)
			r.Get("/campaign/{id}", compaignHandler.GetCampaignByID)
			r.Delete("/campaign/{id}", compaignHandler.DeleteCampaign)
			// Admin Routes
			r.Group(func(r chi.Router) {
				r.Use(appMiddleware.RequireRole("admin"))
				r.Get("/admin/healthz", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("OK"))
				})
			})
		})
	})
}
