package models

import (
	"database/sql"
)

// PostalCodeCoordinates represents the geographic coordinates of a postal code
type PostalCodeCoordinates struct {
	PostalCode string         `json:"postal_code" db:"postal_code"`
	Latitude   sql.NullFloat64 `json:"latitude" db:"latitude"`
	Longitude  sql.NullFloat64 `json:"longitude" db:"longitude"`
	Error      string         `json:"error,omitempty" db:"-"`
}

// PostalCodeLocation represents a grouping of businesses by postal code
type PostalCodeLocation struct {
	PostalCode    string      `json:"postal_code"`
	Latitude      float64     `json:"latitude"`
	Longitude     float64     `json:"longitude"`
	Businesses    []Business  `json:"businesses"`
	TotalLMIAs    int         `json:"total_lmias"`
	BusinessCount int         `json:"business_count"`
}

// Business represents a single business within a postal code location
type Business struct {
	Employer          string `json:"employer"`
	Occupation        string `json:"occupation"`
	ApprovedLMIAs     int    `json:"approved_lmias"`
	ApprovedPositions int    `json:"approved_positions"`
}
