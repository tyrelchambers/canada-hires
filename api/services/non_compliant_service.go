package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"canada-hires/scraper"
	scraper_types "canada-hires/scraper-types"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type NonCompliantService interface {
	ScrapeAndStoreNonCompliantEmployers() (*models.CronJob, error)
	GetNonCompliantEmployers(limit, offset int) ([]models.NonCompliantEmployerWithReasons, error)
	GetNonCompliantEmployersCount() (int, error)
	GetNonCompliantReasons() ([]models.NonCompliantReason, error)
	GetLatestScrapeInfo() (*models.CronJob, error)
	ExtractAndGeocodeEmployers() error
	GeocodeEmployersUsingAddresses() error
	ScrapeAndProcessAllGeocoding() (*models.CronJob, error)
	GetNonCompliantLocationsByPostalCode(limit int) (*models.NonCompliantLocationResponse, error)
	GetNonCompliantEmployersByPostalCode(postalCode string, limit, offset int) (*models.NonCompliantEmployersByPostalCodeResponse, error)
	GetNonCompliantEmployersByCoordinates(lat, lng float64, limit, offset int) (*models.NonCompliantEmployersByPostalCodeResponse, error)
	CleanExistingAddresses() error
}

type nonCompliantService struct {
	repo              repos.NonCompliantRepository
	logger            *log.Logger
	postalCodeService PostalCodeService
	geocodingService  PostalCodeGeocodingService
}

func NewNonCompliantService(repo repos.NonCompliantRepository, logger *log.Logger, postalCodeService PostalCodeService, geocodingService PostalCodeGeocodingService) NonCompliantService {
	return &nonCompliantService{
		repo:              repo,
		logger:            logger,
		postalCodeService: postalCodeService,
		geocodingService:  geocodingService,
	}
}

func (s *nonCompliantService) ScrapeAndStoreNonCompliantEmployers() (*models.CronJob, error) {
	s.logger.Info("Starting non-compliant employers scraper")

	// Create cron job record to track progress
	cronJob := &models.CronJob{
		JobName:   "non_compliant_scraper",
		Status:    "running",
		StartedAt: time.Now(),
	}

	// Note: We don't have a cron jobs table for non-compliant, but we'll create one similar to LMIA
	// For now, we'll return a mock cron job and focus on the scraping logic

	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("Panic occurred during scraping", "error", r)
			cronJob.Status = "failed"
			cronJob.ErrorMessage = func() *string { msg := fmt.Sprintf("Panic: %v", r); return &msg }()
		}
	}()

	// Initialize scraper
	scraperInstance, err := scraper.NewScraper()
	if err != nil {
		cronJob.Status = "failed"
		cronJob.ErrorMessage = func() *string { msg := fmt.Sprintf("Failed to initialize scraper: %v", err); return &msg }()
		return cronJob, fmt.Errorf("failed to create scraper: %w", err)
	}
	defer scraperInstance.Close()

	// Scrape non-compliant employers data with reason descriptions
	employersData, reasonDescriptions, err := scraperInstance.ScrapeNonCompliantEmployersWithReasons()
	if err != nil {
		cronJob.Status = "failed"
		cronJob.ErrorMessage = func() *string { msg := fmt.Sprintf("Scraping failed: %v", err); return &msg }()
		return cronJob, fmt.Errorf("failed to scrape non-compliant employers: %w", err)
	}

	s.logger.Info("Scraping completed", "employers_found", len(employersData), "reason_descriptions_found", len(reasonDescriptions))

	// Update reason descriptions in the database
	if len(reasonDescriptions) > 0 {
		s.logger.Info("Updating reason descriptions", "count", len(reasonDescriptions))
		for reasonCode, description := range reasonDescriptions {
			_, err := s.repo.UpsertReason(reasonCode, description)
			if err != nil {
				s.logger.Error("Failed to upsert reason", "code", reasonCode, "error", err)
				// Don't fail the whole operation for reason updates
			}
		}
		s.logger.Info("Reason descriptions updated successfully")
	}

	// Convert scraper data to models
	scraperData := make([]models.ScraperNonCompliantData, len(employersData))
	for i, employer := range employersData {
		scraperData[i] = convertToScraperNonCompliantData(employer)
	}

	// Use upsert approach to prevent duplicates while preserving data
	// Note: Uniqueness is based on business_name + address + date_of_final_decision
	// This allows the same business to have multiple violations over time
	s.logger.Info("Upserting scraped employers with reasons (prevents duplicates)", "count", len(scraperData))
	err = s.repo.UpsertEmployersWithReasons(scraperData)
	if err != nil {
		cronJob.Status = "failed"
		cronJob.ErrorMessage = func() *string { msg := fmt.Sprintf("Failed to upsert employers: %v", err); return &msg }()
		return cronJob, fmt.Errorf("failed to upsert scraped data: %w", err)
	}

	// Complete the job
	cronJob.Status = "completed"
	completedAt := time.Now()
	cronJob.CompletedAt = &completedAt
	cronJob.ResourcesProcessed = 1 // Number of scrape operations
	cronJob.RecordsProcessed = len(scraperData)

	s.logger.Info("Non-compliant employers scraping completed successfully",
		"employers_scraped", len(employersData),
		"records_processed", len(scraperData))

	return cronJob, nil
}

func (s *nonCompliantService) GetNonCompliantEmployers(limit, offset int) ([]models.NonCompliantEmployerWithReasons, error) {
	return s.repo.GetEmployersWithReasons(limit, offset)
}

func (s *nonCompliantService) GetNonCompliantEmployersCount() (int, error) {
	return s.repo.GetTotalEmployersCount()
}

func (s *nonCompliantService) GetNonCompliantReasons() ([]models.NonCompliantReason, error) {
	return s.repo.GetAllReasons()
}

func (s *nonCompliantService) GetLatestScrapeInfo() (*models.CronJob, error) {
	// For now, return a mock job with the latest scraped date
	latestDate, err := s.repo.GetLatestScrapedDate()
	if err != nil {
		return nil, err
	}

	if latestDate == nil {
		return nil, nil // No scrapes yet
	}

	// Get total count for the mock job
	totalCount, err := s.repo.GetTotalEmployersCount()
	if err != nil {
		totalCount = 0
	}

	mockJob := &models.CronJob{
		JobName:            "non_compliant_scraper",
		Status:             "completed",
		StartedAt:          *latestDate,
		CompletedAt:        latestDate,
		ResourcesProcessed: 1,
		RecordsProcessed:   totalCount,
		CreatedAt:          *latestDate,
	}

	return mockJob, nil
}

// convertToScraperNonCompliantData converts scraper output to model input format
func convertToScraperNonCompliantData(employer scraper_types.NonCompliantEmployerData) models.ScraperNonCompliantData {
	return models.ScraperNonCompliantData{
		BusinessOperatingName: employer.BusinessOperatingName,
		BusinessLegalName:     employer.BusinessLegalName,
		Address:               employer.Address,
		ReasonCodes:           employer.ReasonCodes,
		DateOfFinalDecision:   employer.DateOfFinalDecision,
		PenaltyAmount:         employer.PenaltyAmount,
		PenaltyCurrency:       employer.PenaltyCurrency,
		Status:                employer.Status,
	}
}

// ExtractAndGeocodeEmployers extracts postal codes from addresses and geocodes them
func (s *nonCompliantService) ExtractAndGeocodeEmployers() error {
	s.logger.Info("Starting postal code extraction and geocoding for non-compliant employers")

	// Get all employers without postal codes
	employers, err := s.repo.GetEmployersWithoutPostalCodes()
	if err != nil {
		return fmt.Errorf("failed to get employers without postal codes: %w", err)
	}

	if len(employers) == 0 {
		s.logger.Info("No employers found without postal codes")
		return nil
	}

	s.logger.Info("Found employers to geocode", "count", len(employers))

	successCount := 0
	errorCount := 0

	for i, employer := range employers {
		if i%500 == 0 && i > 0 {
			s.logger.Info("Postal code extraction progress", "processed", i, "total", len(employers))
		}

		// Skip if no address
		if employer.Address == nil || *employer.Address == "" {
			errorCount++
			continue
		}

		// Extract postal code from address
		postalCode := s.postalCodeService.ExtractPostalCode(*employer.Address)
		if postalCode == "" {
			errorCount++
			continue
		}

		// Validate the postal code exists in our postal_codes table
		_, _, err := s.geocodingService.GeocodePostalCode(postalCode)
		if err != nil {
			errorCount++
			continue
		}

		// Update the employer record with postal code
		err = s.repo.UpdateEmployerPostalCode(employer.ID, postalCode)
		if err != nil {
			// Skip individual update error logs to reduce noise - will log summary at the end
			errorCount++
			continue
		}

		successCount++
	}

	s.logger.Info("Geocoding completed",
		"total_processed", len(employers),
		"successful", successCount,
		"failed", errorCount)

	return nil
}

// GetNonCompliantLocationsByPostalCode returns aggregated location data for the map
func (s *nonCompliantService) GetNonCompliantLocationsByPostalCode(limit int) (*models.NonCompliantLocationResponse, error) {
	locations, err := s.repo.GetLocationsByPostalCode(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations by postal code: %w", err)
	}

	return &models.NonCompliantLocationResponse{
		Locations: locations,
		Count:     len(locations),
		Limit:     limit,
	}, nil
}

// GetNonCompliantEmployersByPostalCode returns all employers for a specific postal code
func (s *nonCompliantService) GetNonCompliantEmployersByPostalCode(postalCode string, limit, offset int) (*models.NonCompliantEmployersByPostalCodeResponse, error) {
	employers, err := s.repo.GetEmployersByPostalCode(postalCode, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get employers by postal code: %w", err)
	}

	// Calculate total penalty for this postal code
	totalPenalty := 0
	for _, employer := range employers {
		if employer.PenaltyAmount != nil {
			totalPenalty += *employer.PenaltyAmount
		}
	}

	return &models.NonCompliantEmployersByPostalCodeResponse{
		Employers:    employers,
		PostalCode:   postalCode,
		Count:        len(employers),
		TotalPenalty: totalPenalty,
	}, nil
}

// GeocodeEmployersUsingAddresses geocodes employers using full address geocoding (clean lookup approach)
func (s *nonCompliantService) GeocodeEmployersUsingAddresses() error {
	s.logger.Info("Starting address geocoding for employers without extractable postal codes")

	// Get employers that don't have extractable postal codes
	employers, err := s.repo.GetEmployersWithoutExtractablePostalCodes()
	if err != nil {
		return fmt.Errorf("failed to get employers without extractable postal codes: %w", err)
	}

	if len(employers) == 0 {
		s.logger.Info("No employers found without extractable postal codes")
		return nil
	}

	s.logger.Info("Found employers to geocode via address lookup", "count", len(employers))

	successCount := 0
	errorCount := 0
	cacheHits := 0

	for i, employer := range employers {
		if i%500 == 0 && i > 0 {
			s.logger.Info("Address geocoding progress", "processed", i, "total", len(employers), 
				"successful", successCount, "failed", errorCount)
		}

		// Skip if no address
		if employer.Address == nil || *employer.Address == "" {
			errorCount++
			continue
		}

		// Geocode the full address (this will check cache first, then geocode if needed)
		_, _, err := s.geocodingService.GeocodeFullAddress(*employer.Address)
		if err != nil {
			// Skip individual error logs to reduce noise - will log summary at the end
			errorCount++
			continue
		}

		// Removed individual geocoding success logs to reduce noise

		successCount++

		// Add a small delay between requests to avoid overwhelming the geocoding service
		time.Sleep(100 * time.Millisecond)
	}

	successRate := float64(successCount) / float64(len(employers)) * 100
	s.logger.Info("Address geocoding completed",
		"total_processed", len(employers),
		"successful", successCount,
		"failed", errorCount,
		"cache_hits", cacheHits,
		"success_rate", fmt.Sprintf("%.1f%%", successRate))

	return nil
}

// ScrapeAndProcessAllGeocoding performs complete workflow: scrape -> postal code geocoding -> address geocoding
func (s *nonCompliantService) ScrapeAndProcessAllGeocoding() (*models.CronJob, error) {
	s.logger.Info("Starting complete non-compliant processing workflow")

	// Step 1: Scrape non-compliant employers data
	s.logger.Info("Step 1: Scraping non-compliant employers")
	cronJob, err := s.ScrapeAndStoreNonCompliantEmployers()
	if err != nil {
		s.logger.Error("Step 1 failed: scraping", "error", err)
		return cronJob, err
	}
	s.logger.Info("Step 1 completed: scraping", "employers_scraped", cronJob.RecordsProcessed)

	// Step 2: Extract and geocode postal codes from addresses
	s.logger.Info("Step 2: Extracting and geocoding postal codes")
	err = s.ExtractAndGeocodeEmployers()
	if err != nil {
		s.logger.Error("Step 2 failed: postal code geocoding", "error", err)
		// Don't fail the whole job, just log and continue
	} else {
		s.logger.Info("Step 2 completed: postal code geocoding")
	}

	// Step 3: Geocode full addresses for employers without postal codes
	s.logger.Info("Step 3: Geocoding full addresses for employers without postal codes")
	err = s.GeocodeEmployersUsingAddresses()
	if err != nil {
		s.logger.Error("Step 3 failed: address geocoding", "error", err)
		// Don't fail the whole job, just log and continue
	} else {
		s.logger.Info("Step 3 completed: address geocoding")
	}

	s.logger.Info("Complete non-compliant processing workflow finished successfully")
	
	// Update the cron job to reflect the complete workflow
	cronJob.JobName = "non_compliant_complete_workflow"
	
	return cronJob, nil
}

// GetNonCompliantEmployersByCoordinates returns all employers for specific lat/lng coordinates
func (s *nonCompliantService) GetNonCompliantEmployersByCoordinates(lat, lng float64, limit, offset int) (*models.NonCompliantEmployersByPostalCodeResponse, error) {
	employers, err := s.repo.GetEmployersByCoordinates(lat, lng, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get employers by coordinates: %w", err)
	}

	// Calculate total penalty for these coordinates
	totalPenalty := 0
	for _, employer := range employers {
		if employer.PenaltyAmount != nil {
			totalPenalty += *employer.PenaltyAmount
		}
	}

	return &models.NonCompliantEmployersByPostalCodeResponse{
		Employers:    employers,
		PostalCode:   fmt.Sprintf("%.6f,%.6f", lat, lng), // Use coordinates as identifier
		Count:        len(employers),
		TotalPenalty: totalPenalty,
	}, nil
}

// CleanExistingAddresses fixes malformed addresses that were previously scraped
func (s *nonCompliantService) CleanExistingAddresses() error {
	s.logger.Info("Starting address cleaning for existing non-compliant employers")
	
	// Get all employers with addresses that need cleaning
	employers, err := s.repo.GetAllEmployers()
	if err != nil {
		return fmt.Errorf("failed to get employers: %w", err)
	}
	
	if len(employers) == 0 {
		s.logger.Info("No employers found to clean")
		return nil
	}
	
	s.logger.Info("Found employers to check for address cleaning", "count", len(employers))
	
	successCount := 0
	errorCount := 0
	unchangedCount := 0
	
	for i, employer := range employers {
		if i%100 == 0 && i > 0 {
			s.logger.Info("Address cleaning progress", "processed", i, "total", len(employers))
		}
		
		// Skip if no address
		if employer.Address == nil || *employer.Address == "" {
			unchangedCount++
			continue
		}
		
		// Clean the address
		originalAddress := *employer.Address
		cleanedAddress := s.cleanAddress(originalAddress)
		
		// Skip if no change needed
		if cleanedAddress == originalAddress {
			unchangedCount++
			continue
		}
		
		// Update the employer record with cleaned address
		err = s.repo.UpdateEmployerAddress(employer.ID, cleanedAddress)
		if err != nil {
			s.logger.Error("Failed to update employer address", "id", employer.ID, "error", err)
			errorCount++
			continue
		}
		
		s.logger.Debug("Cleaned address", "id", employer.ID, "original", originalAddress, "cleaned", cleanedAddress)
		successCount++
	}
	
	s.logger.Info("Address cleaning completed",
		"total_processed", len(employers),
		"successful", successCount,
		"unchanged", unchangedCount,
		"failed", errorCount)
	
	return nil
}

// cleanAddress parses and reformats malformed addresses
func (s *nonCompliantService) cleanAddress(rawAddress string) string {
	if rawAddress == "" {
		return rawAddress
	}
	
	// Canadian provinces/territories (both English and French)
	provinces := []string{
		"Alberta", "AB", "Colombie-Britannique", "British Columbia", "BC",
		"Manitoba", "MB", "Nouveau-Brunswick", "New Brunswick", "NB",
		"Terre-Neuve-et-Labrador", "Newfoundland and Labrador", "NL",
		"Territoires du Nord-Ouest", "Northwest Territories", "NT",
		"Nouvelle-Écosse", "Nova Scotia", "NS", "Nunavut", "NU",
		"Ontario", "ON", "Île-du-Prince-Édouard", "Prince Edward Island", "PE",
		"Québec", "Quebec", "QC", "Saskatchewan", "SK", "Yukon", "YT",
	}
	
	// Create regex pattern for provinces
	provincePattern := "(" + strings.Join(provinces, "|") + ")"
	regex := regexp.MustCompile(provincePattern + `\s*$`)
	
	// Check if address ends with a province
	provinceMatch := regex.FindStringSubmatch(rawAddress)
	if len(provinceMatch) == 0 {
		return rawAddress // No province found, return as-is
	}
	
	province := provinceMatch[1]
	addressWithoutProvince := strings.TrimSpace(regex.ReplaceAllString(rawAddress, ""))
	
	// Look for common patterns where city runs into street address
	// Pattern 1: Street name directly followed by city name (no comma)
	streetPatterns := []*regexp.Regexp{
		regexp.MustCompile(`^(.+(?:rue|street|ave|avenue|blvd|boulevard|rd|road|dr|drive|place|pl|way)\s+\w+)([A-Z][a-z]+(?:\s+[A-Z][a-z]+)*)\s*,?\s*$`),
		regexp.MustCompile(`^(.+\s+)([A-Z][a-z]+(?:\s+[A-Z][a-z]+)*)\s*,?\s*$`),
	}
	
	for _, pattern := range streetPatterns {
		match := pattern.FindStringSubmatch(addressWithoutProvince)
		if len(match) > 2 && match[2] != "" {
			streetAddress := strings.TrimSpace(match[1])
			cityName := strings.TrimSpace(match[2])
			
			// Additional validation: city name should be reasonable length (2+ chars, not all caps)
			if len(cityName) > 1 && !regexp.MustCompile(`^[A-Z]+$`).MatchString(cityName) {
				return streetAddress + ", " + cityName + ", " + province
			}
		}
	}
	
	// Fallback: if we have a province but couldn't parse the city, just add commas where needed
	if strings.Contains(addressWithoutProvince, ",") {
		return rawAddress // Already has commas, probably formatted correctly
	} else {
		// Try to add comma before last word before province (assuming it's the city)
		parts := strings.Fields(strings.TrimSpace(addressWithoutProvince))
		if len(parts) > 1 {
			cityPart := parts[len(parts)-1]
			streetPart := strings.Join(parts[:len(parts)-1], " ")
			return streetPart + ", " + cityPart + ", " + province
		}
	}
	
	return rawAddress // Return original if can't parse
}