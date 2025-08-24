package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"canada-hires/scraper"
	scraper_types "canada-hires/scraper-types"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

type NonCompliantService interface {
	ScrapeAndStoreNonCompliantEmployers() (*models.CronJob, error)
	GetNonCompliantEmployers(limit, offset int) ([]models.NonCompliantEmployerWithReasonCodes, error)
	GetNonCompliantEmployersCount() (int, error)
	GetNonCompliantReasons() ([]models.NonCompliantReason, error)
	GetLatestScrapeInfo() (*models.CronJob, error)
	ExtractAndGeocodeEmployers() error
	GetNonCompliantLocationsByPostalCode(limit int) (*models.NonCompliantLocationResponse, error)
	GetNonCompliantEmployersByPostalCode(postalCode string, limit, offset int) (*models.NonCompliantEmployersByPostalCodeResponse, error)
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

	// Scrape non-compliant employers data
	employersData, err := scraperInstance.ScrapeNonCompliantEmployers()
	if err != nil {
		cronJob.Status = "failed"
		cronJob.ErrorMessage = func() *string { msg := fmt.Sprintf("Scraping failed: %v", err); return &msg }()
		return cronJob, fmt.Errorf("failed to scrape non-compliant employers: %w", err)
	}

	s.logger.Info("Scraping completed", "employers_found", len(employersData))

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

func (s *nonCompliantService) GetNonCompliantEmployers(limit, offset int) ([]models.NonCompliantEmployerWithReasonCodes, error) {
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
		if i%100 == 0 && i > 0 {
			s.logger.Info("Geocoding progress", "processed", i, "total", len(employers))
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
			s.logger.Error("Failed to update employer geolocation", "employer_id", employer.ID, "error", err)
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