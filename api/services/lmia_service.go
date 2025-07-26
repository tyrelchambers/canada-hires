package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
)

type LMIAService interface {
	FetchAndStoreResources() error
	DownloadAndProcessResource(resource *models.LMIAResource) error
	ProcessAllUnprocessedResources() error
	RunFullUpdate() error
	GetLatestUpdateStatus() (*models.CronJob, error)
}

type lmiaService struct {
	repo   repos.LMIARepository
	parser LMIAParser
	client *http.Client
}

type OpenDataResponse struct {
	Result struct {
		Resources []struct {
			ID            string   `json:"id"`
			Name          string   `json:"name"`
			URL           string   `json:"url"`
			Format        string   `json:"format"`
			Language      []string `json:"language"`
			Size          *int64   `json:"size"`
			LastModified  *string  `json:"last_modified"`
			DatePublished *string  `json:"date_published"`
		} `json:"resources"`
	} `json:"result"`
}

func NewLMIAService(repo repos.LMIARepository) LMIAService {
	return &lmiaService{
		repo:   repo,
		parser: NewLMIAParser(),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *lmiaService) FetchAndStoreResources() error {
	log.Info("Fetching LMIA resources from Open Canada API")

	// Create cron job record
	cronJob := &models.CronJob{
		JobName:   "lmia_data_fetch",
		Status:    "running",
		StartedAt: time.Now(),
	}

	err := s.repo.CreateCronJob(cronJob)
	if err != nil {
		return fmt.Errorf("failed to create cron job record: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			errorMsg := fmt.Sprintf("Panic occurred: %v", r)
			s.repo.UpdateCronJobStatus(cronJob.ID, "failed", &errorMsg)
		}
	}()

	// Fetch data from Open Canada API
	apiURL := "https://open.canada.ca/data/api/3/action/package_show?id=90fed587-1364-4f33-a9ee-208181dc0b97"

	resp, err := s.client.Get(apiURL)
	if err != nil {
		errorMsg := fmt.Sprintf("failed to fetch API data: %v", err)
		s.repo.UpdateCronJobStatus(cronJob.ID, "failed", &errorMsg)
		return fmt.Errorf("failed to fetch API data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMsg := fmt.Sprintf("failed to read response body: %v", err)
		s.repo.UpdateCronJobStatus(cronJob.ID, "failed", &errorMsg)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse OpenDataResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		errorMsg := fmt.Sprintf("failed to parse JSON response: %v", err)
		s.repo.UpdateCronJobStatus(cronJob.ID, "failed", &errorMsg)
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	processedCount := 0

	// Process English resources only
	for _, resource := range apiResponse.Result.Resources {
		// Skip if not English
		isEnglish := false
		for _, lang := range resource.Language {
			if lang == "en" {
				isEnglish = true
				break
			}
		}
		if !isEnglish {
			continue
		}

		// Skip if not CSV or XLSX format
		if resource.Format != "CSV" && resource.Format != "XLSX" && resource.Format != "XLS" {
			continue
		}

		// Check if resource already exists
		existing, err := s.repo.GetResourceByResourceID(resource.ID)
		if err == nil && existing != nil {
			log.Info("Resource already exists, skipping", "resource_id", resource.ID)
			continue
		}

		// Parse quarter and year from name
		quarter, year := s.parseQuarterAndYear(resource.Name)
		if quarter == "" || year == 0 {
			log.Warn("Could not parse quarter/year from resource name", "name", resource.Name)
			continue
		}

		// Parse timestamps
		var lastModified, datePublished *time.Time
		if resource.LastModified != nil {
			if t, err := time.Parse(time.RFC3339, *resource.LastModified); err == nil {
				lastModified = &t
			}
		}
		if resource.DatePublished != nil {
			if t, err := time.Parse("2006-01-02", *resource.DatePublished); err == nil {
				datePublished = &t
			}
		}

		lmiaResource := &models.LMIAResource{
			ResourceID:    resource.ID,
			Name:          resource.Name,
			Quarter:       quarter,
			Year:          year,
			URL:           resource.URL,
			Format:        resource.Format,
			Language:      "en",
			SizeBytes:     resource.Size,
			LastModified:  lastModified,
			DatePublished: datePublished,
		}

		err = s.repo.CreateResource(lmiaResource)
		if err != nil {
			log.Error("Failed to create LMIA resource", "error", err, "resource_id", resource.ID)
			continue
		}

		processedCount++
		log.Info("Created LMIA resource", "resource_id", resource.ID, "quarter", quarter, "year", year)
	}

	// Update cron job as completed
	err = s.repo.UpdateCronJobCompleted(cronJob.ID, processedCount, 0)
	if err != nil {
		log.Error("Failed to update cron job as completed", "error", err)
	}

	log.Info("LMIA resources fetch completed", "processed_count", processedCount)
	return nil
}

func (s *lmiaService) parseQuarterAndYear(name string) (string, int) {
	// Match patterns like "2024Q1", "2023Q2", etc.
	re := regexp.MustCompile(`(\d{4})Q(\d)`)
	matches := re.FindStringSubmatch(name)
	if len(matches) == 3 {
		year, _ := strconv.Atoi(matches[1])
		quarter := "Q" + matches[2]
		return quarter, year
	}

	// Match patterns like "2024Q1Q2" for multi-quarter files
	re = regexp.MustCompile(`(\d{4})Q(\d)Q(\d)`)
	matches = re.FindStringSubmatch(name)
	if len(matches) == 4 {
		year, _ := strconv.Atoi(matches[1])
		quarter := "Q" + matches[2] + "Q" + matches[3]
		return quarter, year
	}

	// Match year patterns like "2015", "2016" etc
	re = regexp.MustCompile(`(\d{4})`)
	matches = re.FindStringSubmatch(name)
	if len(matches) >= 2 {
		year, _ := strconv.Atoi(matches[1])
		// If it's just a year, assume it's annual data
		return "ANNUAL", year
	}

	return "", 0
}

func (s *lmiaService) DownloadAndProcessResource(resource *models.LMIAResource) error {
	log.Info("Processing LMIA resource", "resource_id", resource.ResourceID)

	// Download and parse the file
	employers, err := s.parser.DownloadAndParseResource(resource)
	if err != nil {
		return fmt.Errorf("failed to download and parse resource: %w", err)
	}

	// Mark as downloaded
	err = s.repo.UpdateResourceDownloaded(resource.ID)
	if err != nil {
		return fmt.Errorf("failed to update resource as downloaded: %w", err)
	}

	// Store employers in batches
	batchSize := 1000
	for i := 0; i < len(employers); i += batchSize {
		end := i + batchSize
		if end > len(employers) {
			end = len(employers)
		}

		batch := employers[i:end]
		err = s.repo.CreateEmployersBatch(batch)
		if err != nil {
			return fmt.Errorf("failed to store employers batch: %w", err)
		}
	}

	// Mark as processed
	err = s.repo.UpdateResourceProcessed(resource.ID)
	if err != nil {
		return fmt.Errorf("failed to update resource as processed: %w", err)
	}

	log.Info("Resource processed successfully", "resource_id", resource.ResourceID, "employers_count", len(employers))
	return nil
}

func (s *lmiaService) ProcessAllUnprocessedResources() error {
	log.Info("Processing all unprocessed LMIA resources")

	// Create cron job record
	cronJob := &models.CronJob{
		JobName:   "lmia_data_process",
		Status:    "running",
		StartedAt: time.Now(),
	}

	err := s.repo.CreateCronJob(cronJob)
	if err != nil {
		return fmt.Errorf("failed to create cron job record: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			errorMsg := fmt.Sprintf("Panic occurred: %v", r)
			s.repo.UpdateCronJobStatus(cronJob.ID, "failed", &errorMsg)
		}
	}()

	// Get all unprocessed resources
	resources, err := s.repo.GetUnprocessedResources()
	if err != nil {
		errorMsg := fmt.Sprintf("failed to get unprocessed resources: %v", err)
		s.repo.UpdateCronJobStatus(cronJob.ID, "failed", &errorMsg)
		return fmt.Errorf("failed to get unprocessed resources: %w", err)
	}

	totalRecords := 0
	processedResources := 0

	for _, resource := range resources {
		err = s.DownloadAndProcessResource(resource)
		if err != nil {
			log.Error("Failed to process resource", "resource_id", resource.ResourceID, "error", err)
			continue
		}

		processedResources++

		// Get count of employers for this resource
		employers, err := s.repo.GetEmployersByResourceID(resource.ID)
		if err == nil {
			totalRecords += len(employers)
		}
	}

	// Update cron job as completed
	err = s.repo.UpdateCronJobCompleted(cronJob.ID, processedResources, totalRecords)
	if err != nil {
		log.Error("Failed to update cron job as completed", "error", err)
	}

	log.Info("Processing completed", "processed_resources", processedResources, "total_records", totalRecords)
	return nil
}

func (s *lmiaService) RunFullUpdate() error {
	log.Info("Running full LMIA data update")

	// First, fetch and store new resources
	err := s.FetchAndStoreResources()
	if err != nil {
		return fmt.Errorf("failed to fetch resources: %w", err)
	}

	// Then process all unprocessed resources
	err = s.ProcessAllUnprocessedResources()
	if err != nil {
		return fmt.Errorf("failed to process resources: %w", err)
	}

	log.Info("Full LMIA data update completed")
	return nil
}

func (s *lmiaService) GetLatestUpdateStatus() (*models.CronJob, error) {
	return s.repo.GetLatestCronJob("lmia_data_fetch")
}
