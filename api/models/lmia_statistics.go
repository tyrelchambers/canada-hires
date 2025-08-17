package models

import (
	"encoding/json"
	"time"
)

type PeriodType string

const (
	PeriodTypeDaily   PeriodType = "daily"
	PeriodTypeMonthly PeriodType = "monthly"
)

type RegionData struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type LMIAStatistics struct {
	ID              string          `json:"id" db:"id"`
	Date            time.Time       `json:"date" db:"date"`
	PeriodType      PeriodType      `json:"period_type" db:"period_type"`
	TotalJobs       int             `json:"total_jobs" db:"total_jobs"`
	UniqueEmployers int             `json:"unique_employers" db:"unique_employers"`
	AvgSalaryMin    *float64        `json:"avg_salary_min" db:"avg_salary_min"`
	AvgSalaryMax    *float64        `json:"avg_salary_max" db:"avg_salary_max"`
	TopProvinces    json.RawMessage `json:"top_provinces" db:"top_provinces"`
	TopCities       json.RawMessage `json:"top_cities" db:"top_cities"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// GetTopProvinces unmarshals the TopProvinces JSON field
func (s *LMIAStatistics) GetTopProvinces() ([]RegionData, error) {
	var provinces []RegionData
	if s.TopProvinces != nil {
		err := json.Unmarshal(s.TopProvinces, &provinces)
		return provinces, err
	}
	return []RegionData{}, nil
}

// SetTopProvinces marshals the provinces data to JSON
func (s *LMIAStatistics) SetTopProvinces(provinces []RegionData) error {
	data, err := json.Marshal(provinces)
	if err != nil {
		return err
	}
	s.TopProvinces = data
	return nil
}

// GetTopCities unmarshals the TopCities JSON field
func (s *LMIAStatistics) GetTopCities() ([]RegionData, error) {
	var cities []RegionData
	if s.TopCities != nil {
		err := json.Unmarshal(s.TopCities, &cities)
		return cities, err
	}
	return []RegionData{}, nil
}

// SetTopCities marshals the cities data to JSON
func (s *LMIAStatistics) SetTopCities(cities []RegionData) error {
	data, err := json.Marshal(cities)
	if err != nil {
		return err
	}
	s.TopCities = data
	return nil
}

// JobStatisticsData represents raw aggregated data used to create statistics
type JobStatisticsData struct {
	TotalJobs       int
	UniqueEmployers int
	AvgSalaryMin    *float64
	AvgSalaryMax    *float64
	ProvincesCounts map[string]int
	CitiesCounts    map[string]int
}

// NewLMIAStatistics creates a new LMIAStatistics from aggregated data
func NewLMIAStatistics(date time.Time, periodType PeriodType, data JobStatisticsData) *LMIAStatistics {
	stats := &LMIAStatistics{
		Date:            date,
		PeriodType:      periodType,
		TotalJobs:       data.TotalJobs,
		UniqueEmployers: data.UniqueEmployers,
		AvgSalaryMin:    data.AvgSalaryMin,
		AvgSalaryMax:    data.AvgSalaryMax,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Convert province counts to sorted slice (top 10)
	provinces := make([]RegionData, 0, len(data.ProvincesCounts))
	for name, count := range data.ProvincesCounts {
		provinces = append(provinces, RegionData{Name: name, Count: count})
	}
	// Sort by count descending and take top 10
	provinces = sortAndLimitRegions(provinces, 10)
	stats.SetTopProvinces(provinces)

	// Convert city counts to sorted slice (top 10)
	cities := make([]RegionData, 0, len(data.CitiesCounts))
	for name, count := range data.CitiesCounts {
		cities = append(cities, RegionData{Name: name, Count: count})
	}
	// Sort by count descending and take top 10
	cities = sortAndLimitRegions(cities, 10)
	stats.SetTopCities(cities)

	return stats
}

// sortAndLimitRegions sorts regions by count (descending) and limits to specified count
func sortAndLimitRegions(regions []RegionData, limit int) []RegionData {
	// Simple bubble sort for small arrays
	for i := 0; i < len(regions)-1; i++ {
		for j := 0; j < len(regions)-i-1; j++ {
			if regions[j].Count < regions[j+1].Count {
				regions[j], regions[j+1] = regions[j+1], regions[j]
			}
		}
	}
	
	if len(regions) > limit {
		return regions[:limit]
	}
	return regions
}