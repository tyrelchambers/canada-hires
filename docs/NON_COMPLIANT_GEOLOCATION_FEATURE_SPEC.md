# Non-Compliant Employers Geolocation Feature Specification

## Overview

This document outlines the implementation plan for adding postal code extraction and geolocation support to the non-compliant employers system. This will enable map-based visualization of non-compliant businesses similar to the existing LMIA heatmap functionality.

## Current State Analysis

### Existing Non-Compliant System
- **Model**: `NonCompliantEmployer` with address field but no geolocation data
- **API Endpoints**: Basic CRUD operations for non-compliant employers and reasons
- **Frontend**: Admin scraper management only, no public-facing map visualization
- **Data Source**: Daily scraping from government non-compliant employer listings

### Existing Geolocation Infrastructure
- **Postal Codes Table**: Already exists with latitude/longitude data
- **Geocoding Service**: `postal_code_geocoding_service.go` for address-to-postal-code extraction
- **LMIA Integration**: LMIA employers already have postal code extraction and mapping

## Feature Requirements

### 1. Database Schema Changes

#### Add Geolocation Fields to NonCompliantEmployer
```sql
-- Migration: 032_add_geolocation_to_non_compliant_employers.up.sql
ALTER TABLE non_compliant_employers 
ADD COLUMN postal_code VARCHAR(7),
ADD COLUMN latitude DECIMAL(10, 8),
ADD COLUMN longitude DECIMAL(11, 8),
ADD COLUMN geocoded_at TIMESTAMP;

CREATE INDEX idx_non_compliant_employers_postal_code ON non_compliant_employers(postal_code);
CREATE INDEX idx_non_compliant_employers_location ON non_compliant_employers(latitude, longitude);
```

#### Update Model Structure
```go
type NonCompliantEmployer struct {
    // ... existing fields ...
    PostalCode  *string    `json:"postal_code" db:"postal_code"`
    Latitude    *float64   `json:"latitude" db:"latitude"`
    Longitude   *float64   `json:"longitude" db:"longitude"`
    GeocodedAt  *time.Time `json:"geocoded_at" db:"geocoded_at"`
}
```

### 2. Backend Implementation

#### Update NonCompliantService
- **Method**: `ExtractAndGeocode()` - Extract postal codes from addresses
- **Method**: `GetNonCompliantLocationsByPostalCode()` - Group by postal code with coordinates
- **Method**: `GetNonCompliantEmployersByPostalCode()` - Get employers for specific postal code
- **Integration**: Use existing `postal_code_geocoding_service.go` for address parsing

#### New API Endpoints
```go
// GET /api/non-compliant/locations
// Returns postal codes with non-compliant employer counts and coordinates
type NonCompliantLocationResponse struct {
    Locations []NonCompliantPostalCodeLocation `json:"locations"`
    Count     int                              `json:"count"`
    Limit     int                              `json:"limit"`
}

type NonCompliantPostalCodeLocation struct {
    PostalCode          string  `json:"postal_code"`
    Latitude           float64 `json:"latitude"`
    Longitude          float64 `json:"longitude"`
    EmployerCount      int     `json:"employer_count"`
    TotalPenaltyAmount int     `json:"total_penalty_amount"`
    ViolationCount     int     `json:"violation_count"`
    MostRecentViolation *time.Time `json:"most_recent_violation"`
}

// GET /api/non-compliant/employers/postal-code/{postal_code}
// Returns all non-compliant employers for a specific postal code
type NonCompliantEmployersByPostalCodeResponse struct {
    Employers   []NonCompliantEmployerWithDetails `json:"employers"`
    PostalCode  string                           `json:"postal_code"`
    Count       int                              `json:"count"`
    TotalPenalty int                             `json:"total_penalty"`
}
```

#### Geocoding Integration
- **Process**: Extract postal codes during scraping process
- **Fallback**: Batch geocoding job for existing records without postal codes
- **Validation**: Ensure postal codes exist in `postal_codes` table
- **Performance**: Cache geocoding results, avoid re-geocoding same addresses

### 3. Data Processing Pipeline

#### During Scraping Process
1. **Address Parsing**: Use `postal_code_geocoding_service` to extract postal code from address
2. **Postal Code Validation**: Check if postal code exists in `postal_codes` table
3. **Coordinate Lookup**: Get latitude/longitude from `postal_codes` table
4. **Record Update**: Store postal code and coordinates in `non_compliant_employers`

#### Batch Geocoding for Existing Data
```go
// Admin endpoint: POST /api/admin/non-compliant/geocode
func (s *NonCompliantService) BatchGeocodeExistingEmployers() error {
    // Get all non-compliant employers without postal codes
    // Extract postal codes from addresses
    // Update records with coordinates
    // Track success/failure rates
}
```

### 4. Frontend Implementation

#### TypeScript Types
```typescript
export interface NonCompliantPostalCodeLocation {
  postal_code: string;
  latitude: number;
  longitude: number;
  employer_count: number;
  total_penalty_amount: number;
  violation_count: number;
  most_recent_violation?: string;
}

export interface NonCompliantLocationResponse {
  locations: NonCompliantPostalCodeLocation[];
  count: number;
  limit: number;
}

export interface NonCompliantEmployerDetails {
  id: string;
  business_operating_name: string;
  business_legal_name?: string;
  address?: string;
  date_of_final_decision?: string;
  penalty_amount?: number;
  penalty_currency: string;
  status?: string;
  reason_codes: string[];
  reasons: NonCompliantReason[];
  scraped_at: string;
}
```

#### React Hooks
```typescript
// hooks/useNonCompliant.ts
export function useNonCompliantLocations(limit: number = 1000) {
  return useQuery<NonCompliantLocationResponse>({
    queryKey: ["non-compliant", "locations", limit],
    queryFn: async () => {
      const response = await apiClient.get(`/non-compliant/locations?limit=${limit}`);
      return response.data;
    },
  });
}

export function useNonCompliantByPostalCode(postalCode: string, limit: number = 100) {
  return useQuery<NonCompliantEmployersByPostalCodeResponse>({
    queryKey: ["non-compliant", "postal-code", postalCode, limit],
    queryFn: async () => {
      const response = await apiClient.get(`/non-compliant/employers/postal-code/${postalCode}?limit=${limit}`);
      return response.data;
    },
    enabled: !!postalCode,
  });
}
```

#### Map Component
- **Base**: Copy `LMIAMapHeatmap.tsx` structure
- **Component**: `NonCompliantMapHeatmap.tsx`
- **Color Scheme**: Orange/red theme to differentiate from LMIA (red) data
- **Markers**: Show violation count and total penalty amounts
- **Popups**: Display comprehensive violation details

#### Route Integration
- **Route**: `/non-compliant-map`
- **Navigation**: Add to main navigation menu
- **SEO**: Meta tags for non-compliant business map

### 5. Data Visualization Features

#### Map Display
- **Markers**: Postal code locations with non-compliant employers
- **Circle Overlays**: Size based on total penalty amounts or violation counts
- **Color Coding**: Intensity based on severity/recent violations
- **Clustering**: Group nearby violations for better performance

#### Sidebar Information
```typescript
// For each postal code area:
- Employer count
- Total penalty amount
- Most recent violation date
- Breakdown by violation types

// For individual employers:
- Business name (operating and legal)
- Full address
- Penalty amount and currency
- Decision date and status
- All violation reasons with descriptions
- Reason code explanations
```

#### Search and Filtering
- **Location Search**: Reuse existing `MapSearch` component
- **Filters**: By penalty amount range, violation date range, reason codes
- **Sorting**: By penalty amount, recent violations, employer name

### 6. Implementation Phases

#### Phase 1: Backend Infrastructure
1. Create database migration for geolocation fields
2. Update `NonCompliantEmployer` model with new fields
3. Integrate postal code extraction in scraping service
4. Add location-based API endpoints
5. Create batch geocoding job for existing data

#### Phase 2: Frontend Foundation
1. Add TypeScript types for location data
2. Create React hooks for location-based queries
3. Copy and modify LMIA map component structure
4. Add new route and navigation

#### Phase 3: Enhanced Visualization
1. Implement comprehensive data display in popups/sidebar
2. Add violation reason descriptions and explanations
3. Implement search and filtering functionality
4. Performance optimization and testing

#### Phase 4: Admin Tools
1. Add geocoding management to admin panel
2. Create monitoring for geocoding success rates
3. Add tools for manual geocoding fixes
4. Performance monitoring and optimization

### 7. Technical Considerations

#### Performance
- **Database Indexing**: Postal code and location indexes for fast queries
- **API Caching**: Cache location data with appropriate TTL
- **Frontend Optimization**: Memoize map markers and circles
- **Lazy Loading**: Load employer details on demand

#### Data Quality
- **Address Parsing**: Handle various Canadian address formats
- **Postal Code Validation**: Ensure extracted codes are valid Canadian postal codes
- **Coordinate Accuracy**: Use existing validated postal code coordinates
- **Error Handling**: Graceful degradation when geocoding fails

#### Security and Privacy
- **Public Data**: All non-compliant data is already public government information
- **Rate Limiting**: Protect geocoding endpoints from abuse
- **Input Validation**: Sanitize postal code and coordinate inputs

### 8. Success Metrics

#### Data Coverage
- **Geocoding Rate**: Percentage of non-compliant employers with valid postal codes
- **Coordinate Accuracy**: Match rate with existing postal code database
- **Data Freshness**: Time between scraping and geocoding completion

#### User Experience
- **Map Performance**: Load times for different data volumes
- **Search Functionality**: Accuracy and speed of location searches
- **Mobile Responsiveness**: Usability on different screen sizes

#### Business Value
- **Public Transparency**: Easy visualization of violation patterns
- **Geographic Insights**: Regional compliance patterns
- **Penalty Awareness**: Clear display of financial consequences

## Dependencies

### Existing Systems
- **Postal Codes Table**: Contains validated Canadian postal codes with coordinates
- **Geocoding Service**: Address parsing and postal code extraction
- **Scraping Infrastructure**: Daily data collection system
- **Map Components**: Leaflet integration and existing UI patterns

### External Services
- **Pelias Geocoding**: For address search functionality
- **OpenStreetMap**: Base map tiles
- **Government Data Source**: Non-compliant employer listings

## Risks and Mitigation

### Technical Risks
- **Geocoding Accuracy**: Some addresses may not parse correctly
  - *Mitigation*: Manual review process, fallback to city-level geocoding
- **Performance**: Large datasets may slow map rendering
  - *Mitigation*: Implement clustering, pagination, and caching
- **Data Synchronization**: Postal codes and employer data may get out of sync
  - *Mitigation*: Add data validation and sync monitoring

### Business Risks
- **Data Privacy**: Ensure only public information is displayed
  - *Mitigation*: All data is already public, no additional privacy concerns
- **Accuracy Concerns**: Incorrect penalty or violation information
  - *Mitigation*: Clear data source attribution, regular data validation

## Conclusion

This feature will provide valuable geographic visualization of non-compliant employers, matching the functionality and user experience of the existing LMIA heatmap. The implementation leverages existing infrastructure while adding comprehensive violation details and penalty information to support public transparency and informed decision-making.