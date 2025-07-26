package services

import (
	"canada-hires/models"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/tealeg/xlsx/v3"
)

type LMIAParser interface {
	DownloadAndParseResource(resource *models.LMIAResource) ([]*models.LMIAEmployer, error)
	ParseCSV(filePath string, resourceID string, year int) ([]*models.LMIAEmployer, error)
	ParseXLSX(filePath string, resourceID string, year int) ([]*models.LMIAEmployer, error)
}

type lmiaParser struct {
	client  *http.Client
	tempDir string
}

func NewLMIAParser() LMIAParser {
	return &lmiaParser{
		client: &http.Client{
			Timeout: 5 * time.Minute, // Longer timeout for large files
		},
		tempDir: "/tmp/lmia_downloads",
	}
}

func (p *lmiaParser) DownloadAndParseResource(resource *models.LMIAResource) ([]*models.LMIAEmployer, error) {
	// Download the file first to get the actual filename
	fileName := fmt.Sprintf("%s.%s", resource.ResourceID, strings.ToLower(resource.Format))
	filePath := filepath.Join(p.tempDir, fileName)
	
	// Use the year from the resource's CreatedAt field.
	year := resource.CreatedAt.Year()
	
	log.Info("Downloading and parsing LMIA resource", 
		"resource_id", resource.ResourceID, 
		"resource_name", resource.Name,
		"url", resource.URL,
		"filename", fileName,
		"format", resource.Format,
		"extracted_year", year)

	// Create temp directory if it doesn't exist
	err := os.MkdirAll(p.tempDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	err = p.downloadFile(resource.URL, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	// Clean up the file after processing
	defer func() {
		if err := os.Remove(filePath); err != nil {
			log.Warn("Failed to remove temp file", "path", filePath, "error", err)
		}
	}()

	// Parse based on format
	var employers []*models.LMIAEmployer
	switch strings.ToUpper(resource.Format) {
	case "CSV":
		employers, err = p.ParseCSV(filePath, resource.ID, year)
	case "XLSX", "XLS":
		employers, err = p.ParseXLSX(filePath, resource.ID, year)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", resource.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	log.Info("Successfully parsed LMIA resource", 
		"resource_id", resource.ResourceID, 
		"year", year,
		"employers_count", len(employers))
	return employers, nil
}

func (p *lmiaParser) downloadFile(url, filePath string) error {
	resp, err := p.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}



func (p *lmiaParser) ParseCSV(filePath string, resourceID string, year int) ([]*models.LMIAEmployer, error) {
	log.Info("Parsing CSV file", "file_path", filePath, "year", year)
	
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Configure reader to be more lenient with malformed CSV
	reader.LazyQuotes = true       // Allow lazy quotes
	reader.TrimLeadingSpace = true // Trim leading spaces
	reader.FieldsPerRecord = -1    // Allow variable number of fields

	// Read records one by one to handle errors gracefully
	var records [][]string
	lineNumber := 0

	for {
		record, err := reader.Read()
		lineNumber++

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Warn("Skipping malformed CSV line", "line", lineNumber, "error", err.Error())
			continue
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty or all lines are malformed")
	}

	// Find the actual header row by looking for expected columns
	var headers []string
	var dataStartIndex int

	for i, record := range records {
		if len(record) == 0 {
			continue
		}

		// Check if this looks like a header row with expected LMIA columns
		hasExpectedColumns := false
		expectedColumnCount := 0
		for _, field := range record {
			lower := strings.ToLower(strings.TrimSpace(field))
			// Look for exact column names, not descriptions
			if lower == "employer" ||
				lower == "address" ||
				lower == "positions" ||
				lower == "occupation" ||
				strings.Contains(lower, "province") ||
				strings.Contains(lower, "program") ||
				strings.Contains(lower, "position") ||
				strings.Contains(lower, "noc") ||
				strings.Contains(lower, "stream") ||
				strings.Contains(lower, "incorporate") ||
				strings.Contains(lower, "approved") {
				hasExpectedColumns = true
				expectedColumnCount++
			}
		}

		if hasExpectedColumns && len(record) >= 2 && expectedColumnCount >= 2 { // Should have multiple expected columns
			headers = record
			dataStartIndex = i + 1
			log.Info("Found header row", "row_index", i, "headers", headers)
			break
		}
	}

	if len(headers) == 0 {
		return nil, fmt.Errorf("could not find valid header row in CSV")
	}

	// Find column indices
	columnMap := p.mapColumns(headers)

	// Log detected columns for debugging
	log.Info("CSV column mapping", "detected_columns", columnMap, "total_headers", len(headers))

	var employers []*models.LMIAEmployer
	skippedCount := 0

	for i, record := range records[dataStartIndex:] { // Skip to data rows
		if len(record) == 0 {
			continue
		}

		employer, reason := p.parseEmployerRecord(record, columnMap, resourceID, year)
		if employer != nil {
			employers = append(employers, employer)
		} else {
			skippedCount++
			if skippedCount <= 10 { // Only log first 10 failures to avoid spam
				log.Warn("Failed to parse employer record", "row", dataStartIndex+i+1, "reason", reason)
			}
		}
	}

	if skippedCount > 10 {
		log.Warn("Additional records skipped", "total_skipped", skippedCount)
	}

	log.Info("CSV parsing completed", "total_records", len(records)-1, "parsed_employers", len(employers), "skipped", skippedCount)
	return employers, nil
}

func (p *lmiaParser) ParseXLSX(filePath string, resourceID string, year int) ([]*models.LMIAEmployer, error) {
	log.Info("Parsing XLSX file", "file_path", filePath, "year", year)
	
	wb, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(wb.Sheets) == 0 {
		return nil, fmt.Errorf("XLSX file has no sheets")
	}

	sheet := wb.Sheets[0]

	var columnMap map[string]int
	var employers []*models.LMIAEmployer
	var headerFound bool
	skippedCount := 0
	rowIndex := 0

	err = sheet.ForEachRow(func(row *xlsx.Row) error {
		rowIndex++
		var record []string
		err := row.ForEachCell(func(cell *xlsx.Cell) error {
			val, err := cell.FormattedValue()
			if err != nil {
				val = cell.String()
			}
			record = append(record, strings.TrimSpace(val))
			return nil
		})
		if err != nil {
			log.Warn("Error reading cell, skipping row", "row", rowIndex, "error", err)
			return nil
		}

		if len(record) == 0 {
			return nil // Skip empty row
		}

		if !headerFound {
			// Check if this looks like a header row with expected LMIA columns
			hasExpectedColumns := false
			expectedColumnCount := 0
			for _, field := range record {
				lower := strings.ToLower(strings.TrimSpace(field))
				// Look for exact column names, not descriptions
				if lower == "employer" ||
					lower == "address" ||
					lower == "positions" ||
					lower == "occupation" ||
					strings.Contains(lower, "province") ||
					strings.Contains(lower, "program") ||
					strings.Contains(lower, "position") ||
					strings.Contains(lower, "noc") ||
					strings.Contains(lower, "stream") ||
					strings.Contains(lower, "incorporate") ||
					strings.Contains(lower, "approved") {
					hasExpectedColumns = true
					expectedColumnCount++
				}
			}

			if hasExpectedColumns && len(record) >= 2 && expectedColumnCount >= 2 { // Should have multiple expected columns
				headers := record
				columnMap = p.mapColumns(headers)
				headerFound = true
				log.Info("Found header row in XLSX", "row_index", rowIndex-1, "headers", headers)
				log.Info("XLSX column mapping", "detected_columns", columnMap)
				return nil // This was the header row, done with it.
			}

			// If we are here, it's not a header row, and we haven't found one yet.
			// It could be a note or empty line before the header. Skip it.
			return nil
		}

		// If we are here, header is found, so this is a data row
		employer, reason := p.parseEmployerRecord(record, columnMap, resourceID, year)
		if employer != nil {
			employers = append(employers, employer)
		} else {
			skippedCount++
			if skippedCount <= 10 {
				log.Warn("Failed to parse employer record", "row", rowIndex, "reason", reason)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if !headerFound {
		return nil, fmt.Errorf("could not find valid header row in XLSX")
	}

	if skippedCount > 10 {
		log.Warn("Additional records skipped", "total_skipped", skippedCount)
	}

	log.Info("XLSX parsing completed", "parsed_employers", len(employers), "skipped", skippedCount)
	return employers, nil
}

func (p *lmiaParser) mapColumns(headers []string) map[string]int {
	columnMap := make(map[string]int)

	for i, header := range headers {
		lower := strings.ToLower(strings.TrimSpace(header))

		// More specific rules first to avoid incorrect mapping by generic rules.
		if strings.Contains(lower, "program") && strings.Contains(lower, "stream") {
			columnMap["program_stream"] = i
		} else if strings.Contains(lower, "approved") && strings.Contains(lower, "lmia") {
			columnMap["approved_lmias"] = i
		} else if strings.Contains(lower, "approved") && strings.Contains(lower, "position") {
			columnMap["approved_positions"] = i
		} else if strings.Contains(lower, "incorporate") {
			columnMap["incorporate_status"] = i
		} else if strings.Contains(lower, "occupation") {
			columnMap["occupation"] = i
		} else if strings.Contains(lower, "address") {
			columnMap["address"] = i
		} else if strings.Contains(lower, "employer") && !strings.Contains(lower, "address") {
			columnMap["employer"] = i
		} else if strings.Contains(lower, "province") || strings.Contains(lower, "territory") {
			columnMap["province_territory"] = i
		} else if strings.Contains(lower, "position") && !strings.Contains(lower, "approved") {
			// Fallback for positions column if not explicitly "approved positions"
			columnMap["approved_positions"] = i
		}
	}

	return columnMap
}

func (p *lmiaParser) parseEmployerRecord(record []string, columnMap map[string]int, resourceID string, year int) (*models.LMIAEmployer, string) {
	// Helper function to safely get field value and clean UTF-8
	getField := func(fieldName string) string {
		if idx, exists := columnMap[fieldName]; exists && idx < len(record) {
			val := strings.TrimSpace(record[idx])
			return strings.ToValidUTF8(val, "")
		}
		return ""
	}

	// Get employer name (required)
	employerName := getField("employer")
	if employerName == "" {
		return nil, "missing employer"
	}

	employer := &models.LMIAEmployer{
		ResourceID: resourceID,
		Year:       year,
		Employer:   employerName,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Parse the ONLY 8 columns we expect
	if val := getField("province_territory"); val != "" {
		employer.ProvinceTerritory = &val
	}
	if val := getField("program_stream"); val != "" {
		employer.ProgramStream = &val
	}
	if val := getField("address"); val != "" {
		employer.Address = &val
	}
	if val := getField("occupation"); val != "" {
		employer.Occupation = &val
	}
	if val := getField("incorporate_status"); val != "" {
		employer.IncorporateStatus = &val
	}
	if val := getField("approved_lmias"); val != "" {
		if num, err := strconv.Atoi(val); err == nil {
			employer.ApprovedLMIAs = &num
		}
	}
	if val := getField("approved_positions"); val != "" {
		if num, err := strconv.Atoi(val); err == nil {
			employer.ApprovedPositions = &num
		}
	}

	return employer, ""
}

func (p *lmiaParser) parseDate(dateStr string) *time.Time {
	// Try common date formats
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"02/01/2006",
		"2006/01/02",
		"Jan 2, 2006",
		"January 2, 2006",
	}

	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return &date
		}
	}

	// Try to extract year if it's just a 4-digit year
	re := regexp.MustCompile(`\b(\d{4})\b`)
	if matches := re.FindStringSubmatch(dateStr); len(matches) > 1 {
		if year, err := strconv.Atoi(matches[1]); err == nil && year > 1900 && year < 2100 {
			date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			return &date
		}
	}

	return nil
}
