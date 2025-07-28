package main

import (
	"canada-hires/container"
	"canada-hires/services"
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found", "error", err)
	}

	// Define command line flags
	var (
		jobTitle     = flag.String("title", "", "Job title to search for (empty for all)")
		province     = flag.String("province", "", "Province to search in (empty for all)")
		pages        = flag.Int("pages", -1, "Number of pages to scrape (-1 for all)")
		saveToAPI    = flag.Bool("api", true, "Save results to API database")
		dryRun       = flag.Bool("dry-run", false, "Run without saving data")
		help         = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		fmt.Println("Job Scraper CLI")
		fmt.Println("Usage: go run cmd/job_scrape/main.go [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  go run cmd/job_scrape/main.go                           # Scrape all jobs")
		fmt.Println("  go run cmd/job_scrape/main.go -title='Software Engineer' # Scrape specific job title")
		fmt.Println("  go run cmd/job_scrape/main.go -province='ON' -pages=5    # Scrape Ontario jobs, 5 pages")
		fmt.Println("  go run cmd/job_scrape/main.go -dry-run                   # Test run without saving")
		return
	}

	log.Info("🇨🇦🍁 Job Bank Scraper - Starting CLI...")
	
	// Create dependency injection container
	cn, err := container.New()
	if err != nil {
		log.Fatal("Failed to create container", "error", err)
	}

	// Create context with timeout for the scraping operation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Execute scraping via the scraper service
	err = cn.Invoke(func(scraperService services.ScraperService) {
		log.Info("Starting job scraping", 
			"title", *jobTitle, 
			"province", *province, 
			"pages", *pages,
			"save_api", *saveToAPI,
			"dry_run", *dryRun)

		if *dryRun {
			log.Info("DRY RUN MODE - No data will be saved")
			// For dry run, just log the parameters
			return
		}

		// Trigger the scraper
		if _, err := scraperService.RunScraperWithConfig(ctx, services.ScraperConfig{
			JobTitle:    *jobTitle,
			Province:    *province,
			Pages:       *pages,
			SaveToAPI:   *saveToAPI,
		}); err != nil {
			log.Fatal("Scraping failed", "error", err)
		}

		log.Info("Job scraping completed successfully")
	})

	if err != nil {
		log.Fatal("Failed to execute scraper", "error", err)
	}
}