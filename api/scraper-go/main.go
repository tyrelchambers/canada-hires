package main

import (
	"fmt"
	"log"

	"scraper/scraper"
)

func main() {
	fmt.Println("🇨🇦🍁 Job Bank Scraper - Starting...")
	fmt.Println("Scraping ALL LMIA jobs (fsrc=32)")

	numberOfPages := -1 // -1 means scrape all available pages

	if err := runScraper(numberOfPages); err != nil {
		log.Fatalf("Scraper failed: %v", err)
	}
}

func runScraper(numberOfPages int) error {
	// Initialize scraper
	scraper, err := scraper.NewScraper()
	if err != nil {
		return fmt.Errorf("failed to create scraper: %v", err)
	}
	defer scraper.Close()

	// Scrape jobs
	jobs, err := scraper.ScrapeLMIAJobs(numberOfPages)
	if err != nil {
		return fmt.Errorf("failed to scrape jobs: %v", err)
	}

	// Print summary statistics
	fmt.Println("\n=== SCRAPING SUMMARY ===")
	fmt.Printf("📊 Total jobs scraped: %d\n", len(jobs))
	
	// Count jobs by province/location
	locationCounts := make(map[string]int)
	for _, job := range jobs {
		locationCounts[job.Location]++
	}
	
	fmt.Printf("📍 Jobs by location:\n")
	for location, count := range locationCounts {
		if location != "" {
			fmt.Printf("   %s: %d jobs\n", location, count)
		}
	}
	
	// Count jobs by business (top 10)
	businessCounts := make(map[string]int)
	for _, job := range jobs {
		if job.Business != "" {
			businessCounts[job.Business]++
		}
	}
	
	fmt.Printf("🏢 Top employers:\n")
	count := 0
	for business, jobCount := range businessCounts {
		if count >= 10 {
			break
		}
		fmt.Printf("   %s: %d jobs\n", business, jobCount)
		count++
	}

	fmt.Println("\n✅ Scraping completed successfully!")
	fmt.Println("💾 Jobs are now ready to be saved to database")
	return nil
}