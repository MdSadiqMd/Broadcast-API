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
	"github.com/MdSadiqMd/Broadcast-API/pkg/config"
	"github.com/go-chi/chi/v5"

	// "github.com/go-co-op/gocron"
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

	err = runMigrations(db)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	r := chi.NewRouter()
	api.Setup(r, db, cfg.JWT.Secret)

	/* s := gocron.NewScheduler(time.UTC)
	err = scheduler.Setup(s, db)
	if err != nil {
		log.Fatalf("Failed to setup scheduler: %v", err)
	}
	s.StartAsync() */

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		log.Println("Shutting down server...")
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("Error during shutdown: %v\n", err)
		}
		// s.Stop()

		serverStopCtx()
	}()

	log.Printf("Server is running on port %d\n", cfg.Server.Port)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v\n", err)
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
	)
}
