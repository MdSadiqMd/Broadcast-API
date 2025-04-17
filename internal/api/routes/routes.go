package api

import (
	"net/http"
	"time"

	"github.com/MdSadiqMd/Broadcast-API/internal/api/handlers"
	appMiddleware "github.com/MdSadiqMd/Broadcast-API/internal/api/middleware"
	"github.com/MdSadiqMd/Broadcast-API/internal/services"
	"github.com/MdSadiqMd/Broadcast-API/pkg/config"
	"github.com/MdSadiqMd/Broadcast-API/pkg/email"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func Setup(r *chi.Mux, db *gorm.DB, jwtSecret string) {
	userService := services.NewUserService(db)
	compaignService := services.NewCampaignService(db)
	contactService := services.NewContactService(db)
	broadcastService := services.NewBroadcastService(db)

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	smtpClient := email.NewSMTPClient(email.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		FromName: cfg.SMTP.FromName,
		FromAddr: cfg.SMTP.FromAddr,
		UseTLS:   false,
	})

	mailService := services.NewMailService(db, smtpClient)
	auth := appMiddleware.NewAuth(appMiddleware.AuthConfig{
		JWTSecret:     jwtSecret,
		TokenDuration: 24 * time.Hour,
		UserService:   userService,
	})

	authHandler := handlers.NewAuthHandler(userService, auth)
	compaignHandler := handlers.NewCampaignHandler(compaignService, auth)
	contactHandler := handlers.NewContactHandler(contactService, auth)
	broadcastHandler := handlers.NewBroadcastHandler(broadcastService, auth)
	mailHandler := handlers.NewMailHandler(mailService, auth)

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

			r.Post("/contact", contactHandler.CreateContact)
			r.Post("/contacts", contactHandler.GetAllContacts)
			r.Get("/contact/{id}", contactHandler.GetContactByID)
			r.Put("/contact/{id}", contactHandler.UpdateContact)
			r.Delete("/contact/{id}", contactHandler.DeleteContact)

			r.Post("/broadcast", broadcastHandler.CreateBroadcast)
			r.Get("/broadcast/{id}", broadcastHandler.GetBroadcastByID)
			r.Put("/broadcast/{id}", broadcastHandler.UpdateBroadcast)
			r.Get("/broadcasts", broadcastHandler.ListBroadcasts)
			r.Post("/broadcast/{id}/send", broadcastHandler.SendBroadcast)
			r.Delete("/broadcast/{id}", broadcastHandler.DeleteBroadcast)

			r.Post("/mail/test", mailHandler.SendTestEmail)
			r.Post("/mail/transactional", mailHandler.SendTransactionalEmail)
			r.Post("/mail/job/{id}/process", mailHandler.ProcessEmailJob)
			r.Post("/mail/campaign/send", mailHandler.ProcessCampaignEmail)
			r.Post("/mail/campaign/bulk", mailHandler.BulkSendCampaign)

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
