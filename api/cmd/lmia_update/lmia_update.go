package main

import (
	"canada-hires/container"
	"canada-hires/db"
	"canada-hires/services"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Warn("Could not load .env file", "error", err)
	}

	// Initialize database
	database := db.InitDB()
	defer database.Close()

	// Test database connection
	if err := database.Ping(); err != nil {
		log.Fatal("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Create container and get LMIA service
	cn, err := container.New()
	if err != nil {
		log.Fatal("Failed to create container", "error", err)
		os.Exit(1)
	}

	var lmiaService services.LMIAService
	if err := cn.Invoke(func(s services.LMIAService) {
		lmiaService = s
	}); err != nil {
		log.Fatal("Failed to get LMIA service", "error", err)
		os.Exit(1)
	}

	log.Info("Starting LMIA data update job...")

	// Run the full update
	if err := lmiaService.RunFullUpdate(); err != nil {
		log.Error("LMIA update failed", "error", err)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	log.Info("LMIA data update completed successfully")
	fmt.Println("LMIA data update completed successfully")
}