# Plan: Geocode Non-Compliant Businesses Without Postal Codes

**Objective**: Increase map coverage from 218 to all available non-compliant businesses by geocoding their addresses directly to coordinates.

## Problem Analysis

Currently, only 218 non-compliant businesses appear on the map because the system only shows businesses that have postal codes. Many businesses have addresses but lack properly formatted postal codes that can be extracted using the current regex pattern `[A-Za-z]\d[A-Za-z]\s*\d[A-Za-z]\d`.

## Solution Overview

Implement full address geocoding as a fallback for businesses without postal codes, using the existing Pelias geocoding infrastructure to convert full addresses directly to latitude/longitude coordinates.

## Phase 1: Backend Infrastructure

### 1. Extend Geocoding Service
- **File**: `api/services/postal_code_geocoding_service.go`
- **Add Method**: `GeocodeFullAddress(address string) (latitude, longitude float64, err error)`
- **Implementation**: Use Pelias `/v1/search?text=` endpoint for full address queries
- **Validation**: Ensure results are within Canada and have reasonable confidence scores

### 2. Update Non-Compliant Models
- **File**: `api/models/non_compliant_employer.go`
- **Add Fields** (if not already present):
  ```go
  Latitude    *float64   `json:"latitude" db:"latitude"`
  Longitude   *float64   `json:"longitude" db:"longitude"`  
  GeocodedAt  *time.Time `json:"geocoded_at" db:"geocoded_at"`
  ```

### 3. Database Migration
- **Create**: `api/migrations/XXX_add_coordinates_to_non_compliant_employers.up.sql`
- **Add Columns**: latitude, longitude, geocoded_at
- **Add Indexes**: For location-based queries

### 4. Extend Repository Layer
- **File**: `api/repos/non_compliant_repository.go`
- **Add Methods**:
  ```go
  GetEmployersWithoutCoordinates() ([]NonCompliantEmployer, error)
  UpdateEmployerCoordinates(employerID string, lat, lng float64) error
  GetLocationsByCoordinates(limit int) ([]NonCompliantPostalCodeLocation, error)
  ```

### 5. Batch Geocoding Service
- **File**: `api/services/non_compliant_service.go`
- **Add Method**: `BatchGeocodeEmployersWithoutCoordinates() error`
- **Logic**:
  1. Get all employers without coordinates but with addresses
  2. For each employer, attempt full address geocoding
  3. Update database with successful coordinates
  4. Log progress and failure rates

## Phase 2: Admin Interface

### 1. Add Geocoding Endpoint
- **File**: `api/controllers/non_compliant_controller.go`
- **Endpoint**: `POST /api/admin/non-compliant/batch-geocode`
- **Response**: Progress stats, success/failure counts

### 2. Update Admin Panel
- **File**: `web/src/components/admin/ScraperManager.tsx`
- **Add Controls**:
  - Geocoding statistics display
  - Manual batch geocoding trigger button
  - Progress indicator during geocoding

## Phase 3: Map Query Updates

### 1. Modify Location Queries
- **File**: `api/repos/non_compliant_repository.go`
- **Update**: `GetLocationsByPostalCode()` method
- **Logic**:
  1. Include employers with postal codes (current behavior)
  2. Include employers with direct coordinates but no postal codes
  3. Group coordinate-only employers by approximate location
  4. Combine results for comprehensive map coverage

### 2. Update Map Components
- **File**: `web/src/components/NonCompliantMapHeatmap.tsx`
- **Modifications**:
  - Handle mixed data sources (postal code clusters + individual coordinates)
  - Ensure popover data works for both data types
  - Maintain consistent visual styling

## Implementation Priority

1. **High Priority**: Backend geocoding service and database updates
2. **Medium Priority**: Admin interface for manual geocoding triggers
3. **Low Priority**: Advanced map clustering for coordinate-only employers

## Expected Outcomes

- **Coverage Increase**: From 218 businesses to potentially 500+ businesses on the map
- **Data Quality**: Better geographic representation of non-compliant employers
- **User Experience**: More comprehensive view of violations across Canada

## Technical Considerations

### Performance
- Implement rate limiting for Pelias requests
- Process geocoding in batches to avoid overwhelming the service
- Cache geocoding results to avoid re-processing

### Data Quality
- Validate geocoding results are within Canada
- Set minimum confidence thresholds for accepting coordinates
- Handle edge cases where addresses are incomplete or invalid

### Fallback Strategy
- If full address geocoding fails, attempt city/province geocoding
- Store geocoding errors for manual review
- Provide admin tools to manually fix problematic addresses

## Success Metrics

- **Coverage Rate**: Percentage of non-compliant employers with valid coordinates
- **Geocoding Success Rate**: Percentage of addresses successfully geocoded
- **Map Performance**: Load times with increased data volume
- **Data Accuracy**: Manual verification of sample geocoded locations