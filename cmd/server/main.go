package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/MdSadiqMd/Broadcast-API/internal/api/routes"
	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"github.com/MdSadiqMd/Broadcast-API/internal/scheduler"
	"github.com/MdSadiqMd/Broadcast-API/pkg/config"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := setupDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	r := chi.NewRouter()
	api.Setup(r, db, cfg.JWT.Secret)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		sched := scheduler.NewScheduler(db, *cfg)
		if err := sched.Start(); err != nil {
			log.Printf("Failed to start scheduler: %v", err)
			return
		}
		defer sched.Stop()
		select {}
	}()

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Println("Graceful shutdown timed out")
			}
		}()

		log.Println("Shutting down server gracefully...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		serverStopCtx()
	}()

	log.Printf("Server starting on port %d", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	<-serverCtx.Done()
	log.Println("Server stopped")
}

func setupDatabase(config config.DatabaseConfig) (*gorm.DB, error) {
	dsn := config.URL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.Name,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Campaign{},
		&models.Contact{},
		&models.Broadcast{},
		&models.CampaignAudience{},
		&models.JWTClaims{},
		&models.EmailJob{},
		&models.EmailLog{},
		&models.Template{},
		&models.Message{},
		&models.Subscriber{},
		&models.List{},
	)
}
