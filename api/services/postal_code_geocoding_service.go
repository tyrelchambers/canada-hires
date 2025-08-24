package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type PostalCodeGeocodingService interface {
	GeocodePostalCode(postalCode string, expectedProvince ...string) (latitude, longitude float64, err error)
	GeocodeMultiplePostalCodes(postalCodes []string) (map[string]models.PostalCodeCoordinates, error)
	GeocodeFullAddress(address string) (latitude, longitude float64, err error)
	GetAllPostalCodes() (map[string]models.PostalCodeCoordinates, error)
	UpsertPostalCode(postalCode *models.PostalCodeCoordinates) error
}

// Pelias geocoding response structure
type PeliasGeocodingInfo struct {
	Version     string                 `json:"version"`
	Attribution string                 `json:"attribution"`
	Query       map[string]interface{} `json:"query"`
	Warnings    []string               `json:"warnings,omitempty"`
	Errors      []string               `json:"errors,omitempty"`
	Engine      map[string]interface{} `json:"engine"`
	Timestamp   int64                  `json:"timestamp"`
}

type PeliasGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"` // [longitude, latitude]
}

type PeliasProperties struct {
	ID           string  `json:"id"`
	GID          string  `json:"gid"`
	Layer        string  `json:"layer"`
	Source       string  `json:"source"`
	SourceID     string  `json:"source_id"`
	CountryCode  string  `json:"country_code,omitempty"`
	Name         string  `json:"name"`
	PostalCode   string  `json:"postalcode,omitempty"`
	Confidence   float64 `json:"confidence"`
	MatchType    string  `json:"match_type,omitempty"`
	Distance     float64 `json:"distance,omitempty"`
	Accuracy     string  `json:"accuracy,omitempty"`
	Country      string  `json:"country,omitempty"`
	CountryGID   string  `json:"country_gid,omitempty"`
	CountryA     string  `json:"country_a,omitempty"`
	Region       string  `json:"region,omitempty"`
	RegionGID    string  `json:"region_gid,omitempty"`
	RegionA      string  `json:"region_a,omitempty"`
	Locality     string  `json:"locality,omitempty"`
	LocalityGID  string  `json:"locality_gid,omitempty"`
	Label        string  `json:"label,omitempty"`
}

type PeliasFeature struct {
	Type       string           `json:"type"`
	Geometry   PeliasGeometry   `json:"geometry"`
	Properties PeliasProperties `json:"properties"`
	Bbox       []float64        `json:"bbox,omitempty"`
}

type PeliasResponse struct {
	Geocoding PeliasGeocodingInfo `json:"geocoding"`
	Type      string              `json:"type"`
	Features  []PeliasFeature     `json:"features"`
	Bbox      []float64           `json:"bbox,omitempty"`
}

type postalCodeGeocodingService struct {
	peliasServerURL   string
	client            *http.Client
	postalCodeRepo    repos.PostalCodeRepository
	postalCodeService PostalCodeService
	addressCacheRepo  repos.AddressGeocodingCacheRepository
}

func NewPostalCodeGeocodingService(postalCodeRepo repos.PostalCodeRepository, postalCodeService PostalCodeService, addressCacheRepo repos.AddressGeocodingCacheRepository) PostalCodeGeocodingService {
	homeserverURL := os.Getenv("HOMESERVER_URL")
	if homeserverURL == "" {
		homeserverURL = "http://homeserver:4000"
		log.Info("Using default homeserver URL", "url", homeserverURL)
	} else {
		log.Info("Homeserver configured", "url", homeserverURL)
	}
	
	peliasServerURL := homeserverURL

	service := &postalCodeGeocodingService{
		peliasServerURL: peliasServerURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		postalCodeRepo:    postalCodeRepo,
		postalCodeService: postalCodeService,
		addressCacheRepo:  addressCacheRepo,
	}

	return service
}

// verifyLocationMatch checks if the geocoding result matches the expected province/territory
func (g *postalCodeGeocodingService) verifyLocationMatch(result PeliasFeature, expectedProvince string) bool {
	if expectedProvince == "" {
		return true // No verification needed if no expected province provided
	}
	
	// Normalize the expected province for comparison
	expectedNormalized := strings.ToLower(strings.TrimSpace(expectedProvince))
	
	// Check various location fields from Pelias response
	fieldsToCheck := []string{
		result.Properties.Region,        // Primary province/territory field
		result.Properties.RegionA,       // Abbreviated province (e.g., "AB", "ON")
		result.Properties.Label,         // Full address label
	}
	
	// Canadian province/territory mappings
	provinceMap := map[string][]string{
		"alberta":                    {"alberta", "ab"},
		"british columbia":           {"british columbia", "bc"},
		"manitoba":                   {"manitoba", "mb"},
		"new brunswick":              {"new brunswick", "nb"},
		"newfoundland and labrador":  {"newfoundland and labrador", "nl", "newfoundland"},
		"northwest territories":      {"northwest territories", "nt"},
		"nova scotia":                {"nova scotia", "ns"},
		"nunavut":                    {"nunavut", "nu"},
		"ontario":                    {"ontario", "on"},
		"prince edward island":       {"prince edward island", "pei", "pe"},
		"quebec":                     {"quebec", "qc", "québec"},
		"saskatchewan":               {"saskatchewan", "sk"},
		"yukon":                      {"yukon", "yt", "yukon territory"},
	}
	
	// Find which province the expected province matches
	var expectedVariants []string
	for province, variants := range provinceMap {
		for _, variant := range variants {
			if variant == expectedNormalized {
				expectedVariants = provinceMap[province]
				break
			}
		}
		if len(expectedVariants) > 0 {
			break
		}
	}
	
	// If we couldn't map the expected province, just do a simple string match
	if len(expectedVariants) == 0 {
		expectedVariants = []string{expectedNormalized}
	}
	
	// Check if any of the Pelias fields match any of the expected variants
	for _, field := range fieldsToCheck {
		if field == "" {
			continue
		}
		fieldNormalized := strings.ToLower(strings.TrimSpace(field))
		
		for _, variant := range expectedVariants {
			if strings.Contains(fieldNormalized, variant) {
				return true
			}
		}
	}
	
	return false
}

func (g *postalCodeGeocodingService) GeocodePostalCode(postalCode string, expectedProvince ...string) (latitude, longitude float64, err error) {
	if postalCode == "" {
		return 0, 0, fmt.Errorf("postal code is empty")
	}

	// Clean and format the postal code using the postal code service
	cleanedPostalCode := g.postalCodeService.FormatPostalCode(postalCode)
	if cleanedPostalCode == "" {
		return 0, 0, fmt.Errorf("invalid postal code format: %s", postalCode)
	}

	// First, check if we already have this postal code in our database
	cachedCoords, err := g.postalCodeRepo.GetByPostalCode(cleanedPostalCode)
	if err == nil && cachedCoords != nil && cachedCoords.Latitude.Valid && cachedCoords.Longitude.Valid {
		return cachedCoords.Latitude.Float64, cachedCoords.Longitude.Float64, nil
	}

	// Build Pelias search URL - search by postal code text
	apiURL := fmt.Sprintf("%s/v1/search?text=%s",
		g.peliasServerURL,
		url.QueryEscape(cleanedPostalCode))

	// Make the request
	resp, err := g.client.Get(apiURL)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to make Pelias geocoding request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("Pelias API returned status %d", resp.StatusCode)
	}

	// Parse the response
	var peliasResult PeliasResponse
	if err := json.NewDecoder(resp.Body).Decode(&peliasResult); err != nil {
		return 0, 0, fmt.Errorf("failed to decode Pelias response: %w", err)
	}

	// Check if we got any results
	if len(peliasResult.Features) == 0 {
		return 0, 0, fmt.Errorf("no geocoding results found for postal code: %s", cleanedPostalCode)
	}

	// Get the coordinates from the first result
	result := peliasResult.Features[0]
	if len(result.Geometry.Coordinates) < 2 {
		return 0, 0, fmt.Errorf("invalid coordinates in Pelias response")
	}
	
	// Verify postal code match - the returned postal code should match our request
	if result.Properties.PostalCode != "" && result.Properties.PostalCode != cleanedPostalCode {
		return 0, 0, fmt.Errorf("geocoded postal code %s does not match requested postal code %s", 
			result.Properties.PostalCode, cleanedPostalCode)
	}
	
	// Verify location match if expected province is provided
	var expectedProvinceStr string
	if len(expectedProvince) > 0 {
		expectedProvinceStr = expectedProvince[0]
	}
	
	if !g.verifyLocationMatch(result, expectedProvinceStr) {
		return 0, 0, fmt.Errorf("geocoded location for postal code %s does not match expected province %s (got: %s)", 
			cleanedPostalCode, expectedProvinceStr, result.Properties.Region)
	}
	
	// Pelias returns coordinates as [longitude, latitude]
	longitude = result.Geometry.Coordinates[0]
	latitude = result.Geometry.Coordinates[1]

	// Save the coordinates to the database for future use
	coordsToSave := &models.PostalCodeCoordinates{
		PostalCode: cleanedPostalCode,
		Latitude:   sql.NullFloat64{Float64: latitude, Valid: true},
		Longitude:  sql.NullFloat64{Float64: longitude, Valid: true},
	}

	if err := g.postalCodeRepo.Upsert(coordsToSave); err != nil {
		// Don't fail the request, just continue silently
	}

	return latitude, longitude, nil
}

func (g *postalCodeGeocodingService) GeocodeMultiplePostalCodes(postalCodes []string) (map[string]models.PostalCodeCoordinates, error) {
	results := make(map[string]models.PostalCodeCoordinates)
	
	if len(postalCodes) == 0 {
		return results, nil
	}
	
	log.Info("Starting sequential geocoding", "total_postal_codes", len(postalCodes))
	
	for i, postalCode := range postalCodes {
		// Log progress every 100 postal codes
		if i%100 == 0 && i > 0 {
			log.Info("Geocoding progress", "processed", i, "total", len(postalCodes), "percent", fmt.Sprintf("%.1f%%", float64(i)/float64(len(postalCodes))*100))
		}
		
		cleanedPostalCode := g.postalCodeService.FormatPostalCode(postalCode)
		if cleanedPostalCode == "" || len(cleanedPostalCode) < 3 || len(cleanedPostalCode) > 10 {
			results[postalCode] = models.PostalCodeCoordinates{
				PostalCode: postalCode,
				Error:      "invalid postal code format",
			}
			continue
		}
		
		latitude, longitude, err := g.GeocodePostalCode(postalCode)
		if err != nil {
			results[postalCode] = models.PostalCodeCoordinates{
				PostalCode: postalCode,
				Error:      err.Error(),
			}
		} else {
			results[postalCode] = models.PostalCodeCoordinates{
				PostalCode: postalCode,
				Latitude:   sql.NullFloat64{Float64: latitude, Valid: true},
				Longitude:  sql.NullFloat64{Float64: longitude, Valid: true},
			}
		}
	}
	
	// Final summary
	successCount := 0
	errorCount := 0
	for _, coords := range results {
		if coords.Error == "" {
			successCount++
		} else {
			errorCount++
		}
	}
	
	log.Info("Sequential geocoding completed",
		"total_processed", len(postalCodes),
		"successful", successCount,
		"failed", errorCount,
		"success_rate", fmt.Sprintf("%.1f%%", float64(successCount)/float64(len(postalCodes))*100))
	
	return results, nil
}

func (g *postalCodeGeocodingService) GetAllPostalCodes() (map[string]models.PostalCodeCoordinates, error) {
	return g.postalCodeRepo.GetAllPostalCodes()
}

func (g *postalCodeGeocodingService) UpsertPostalCode(postalCode *models.PostalCodeCoordinates) error {
	return g.postalCodeRepo.Upsert(postalCode)
}

// GeocodeFullAddress geocodes a full address string using Pelias search API
func (g *postalCodeGeocodingService) GeocodeFullAddress(address string) (latitude, longitude float64, err error) {
	if address == "" {
		return 0, 0, fmt.Errorf("address is empty")
	}

	// Clean the address by trimming whitespace
	cleanedAddress := strings.TrimSpace(address)
	if cleanedAddress == "" {
		return 0, 0, fmt.Errorf("address is empty after cleaning")
	}

	// Check cache first
	normalizedAddress := repos.NormalizeAddress(cleanedAddress)
	if cached, err := g.addressCacheRepo.GetByNormalizedAddress(normalizedAddress); err == nil && cached != nil {
		return cached.Latitude, cached.Longitude, nil
	}

	// Build Pelias search URL for full address search
	// Use boundary.country=CAN to restrict results to Canada
	apiURL := fmt.Sprintf("%s/v1/search?text=%s&boundary.country=CAN&size=1",
		g.peliasServerURL,
		url.QueryEscape(cleanedAddress))

	// Make the request
	resp, err := g.client.Get(apiURL)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to make Pelias geocoding request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("Pelias API returned status %d", resp.StatusCode)
	}

	// Parse the response
	var peliasResult PeliasResponse
	if err := json.NewDecoder(resp.Body).Decode(&peliasResult); err != nil {
		return 0, 0, fmt.Errorf("failed to decode Pelias response: %w", err)
	}

	// Check if we got any results
	if len(peliasResult.Features) == 0 {
		return 0, 0, fmt.Errorf("no geocoding results found for address: %s", cleanedAddress)
	}

	// Get the first result
	result := peliasResult.Features[0]
	if len(result.Geometry.Coordinates) < 2 {
		return 0, 0, fmt.Errorf("invalid coordinates in Pelias response")
	}


	// Ensure the result is in Canada
	if result.Properties.CountryCode != "CA" && result.Properties.CountryA != "CAN" {
		return 0, 0, fmt.Errorf("geocoding result is not in Canada: %s for address: %s", 
			result.Properties.Country, cleanedAddress)
	}

	// Pelias returns coordinates as [longitude, latitude]
	longitude = result.Geometry.Coordinates[0]
	latitude = result.Geometry.Coordinates[1]

	// Validate coordinates are within reasonable bounds for Canada
	// Canada latitude: approximately 41.7 to 83.1
	// Canada longitude: approximately -141.0 to -52.6
	if latitude < 41.0 || latitude > 84.0 || longitude < -142.0 || longitude > -52.0 {
		return 0, 0, fmt.Errorf("coordinates outside Canada bounds: lat=%.6f, lng=%.6f for address: %s", 
			latitude, longitude, cleanedAddress)
	}

	// Save to cache for future use
	cacheEntry := &models.AddressGeocodingCache{
		Address:           cleanedAddress,
		NormalizedAddress: normalizedAddress,
		Latitude:          latitude,
		Longitude:         longitude,
		Confidence:        &result.Properties.Confidence,
		GeocodedAt:        time.Now(),
	}

	// Don't fail the whole request if cache save fails
	if err := g.addressCacheRepo.Upsert(cacheEntry); err != nil {
		// Just log the error and continue
	}

	return latitude, longitude, nil
}


