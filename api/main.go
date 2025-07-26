package main

import (
	"canada-hires/container"
	"canada-hires/db"
	"canada-hires/router"
	"canada-hires/services"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	database := db.InitDB()
	defer database.Close()

	// Run migrations on startup
	if err := db.RunMigrations(database.GetDB().DB); err != nil {
		log.Fatal("Failed to run migrations", "error", err)
	}
	// Create router
	r := chi.NewRouter()
	// Start the application with dependency injection
	// Test database connection
	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Info("Successfully connected to database")

	// Setup router with CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://canadahires.info"},
		AllowedMethods:   []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Origin", "Content-Type", "file-type", "Authorization", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	cn, err := container.New()
	if err != nil {
		log.Fatal("Failed to create container", "error", err)
	}

	router.InitRoutes(cn, r)

	// Start cron service
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = cn.Invoke(func(cronService services.CronService) {
		go cronService.Start(ctx)
	})
	if err != nil {
		log.Fatal("Failed to start cron service", "error", err)
	}

	// Setup graceful shutdown
	server := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	// Create a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Info("Server starting on port 8000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", "error", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Info("Shutting down server...")

	// Cancel context to stop cron service
	cancel()

	// Shutdown server gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	} else {
		log.Info("Server exited gracefully")
	}

}
