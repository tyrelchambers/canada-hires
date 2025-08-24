package models

import (
	"time"

	"github.com/lib/pq"
)

type NonCompliantEmployer struct {
	ID                    string     `json:"id" db:"id"`
	BusinessOperatingName string     `json:"business_operating_name" db:"business_operating_name"`
	BusinessLegalName     *string    `json:"business_legal_name" db:"business_legal_name"`
	Address               *string    `json:"address" db:"address"`
	DateOfFinalDecision   *time.Time `json:"date_of_final_decision" db:"date_of_final_decision"`
	PenaltyAmount         *int       `json:"penalty_amount" db:"penalty_amount"`
	PenaltyCurrency       string     `json:"penalty_currency" db:"penalty_currency"`
	Status                *string    `json:"status" db:"status"`
	ReasonCodes           []string   `json:"reason_codes" db:"reason_codes"`
	PostalCode            *string    `json:"postal_code" db:"postal_code"`
	ScrapedAt             time.Time  `json:"scraped_at" db:"scraped_at"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`

	// Related reasons (populated by joins)
	Reasons []NonCompliantReason `json:"reasons,omitempty"`
}

type NonCompliantReason struct {
	ID          int       `json:"id" db:"id"`
	ReasonCode  string    `json:"reason_code" db:"reason_code"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ScraperNonCompliantData represents the raw scraped data before database insertion
type ScraperNonCompliantData struct {
	BusinessOperatingName string
	BusinessLegalName     string
	Address               string
	ReasonCodes           []string // Array of reason codes like ["5", "6", "15"]
	DateOfFinalDecision   string   // Will be parsed to time.Time
	PenaltyAmount         int
	PenaltyCurrency       string
	Status                string
}

// NonCompliantEmployerWithReasons is a flattened view for API responses
type NonCompliantEmployerWithReasons struct {
	ID                    string               `json:"id" db:"id"`
	BusinessOperatingName string               `json:"business_operating_name" db:"business_operating_name"`
	BusinessLegalName     *string              `json:"business_legal_name" db:"business_legal_name"`
	Address               *string              `json:"address" db:"address"`
	DateOfFinalDecision   *time.Time           `json:"date_of_final_decision" db:"date_of_final_decision"`
	PenaltyAmount         *int                 `json:"penalty_amount" db:"penalty_amount"`
	PenaltyCurrency       string               `json:"penalty_currency" db:"penalty_currency"`
	Status                *string              `json:"status" db:"status"`
	ReasonCodes           pq.StringArray       `json:"reason_codes" db:"reason_codes"` // Array of reason codes
	Reasons               []NonCompliantReason `json:"reasons,omitempty" db:"reasons"` // Full reason objects with descriptions
	PostalCode            *string              `json:"postal_code" db:"postal_code"`
	ScrapedAt             time.Time            `json:"scraped_at" db:"scraped_at"`
	CreatedAt             time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at" db:"updated_at"`
}

// NonCompliantPostalCodeLocation represents aggregated location data for the map
type NonCompliantPostalCodeLocation struct {
	PostalCode          string     `json:"postal_code" db:"postal_code"`
	Latitude            float64    `json:"latitude" db:"latitude"`
	Longitude           float64    `json:"longitude" db:"longitude"`
	EmployerCount       int        `json:"employer_count" db:"employer_count"`
	TotalPenaltyAmount  int        `json:"total_penalty_amount" db:"total_penalty_amount"`
	ViolationCount      int        `json:"violation_count" db:"violation_count"`
	MostRecentViolation *time.Time `json:"most_recent_violation" db:"most_recent_violation"`
}

// NonCompliantLocationResponse is the response for location endpoint
type NonCompliantLocationResponse struct {
	Locations []NonCompliantPostalCodeLocation `json:"locations"`
	Count     int                              `json:"count"`
	Limit     int                              `json:"limit"`
}

// NonCompliantEmployersByPostalCodeResponse is the response for employers by postal code
type NonCompliantEmployersByPostalCodeResponse struct {
	Employers    []NonCompliantEmployerWithReasons `json:"employers"`
	PostalCode   string                            `json:"postal_code"`
	Count        int                               `json:"count"`
	TotalPenalty int                               `json:"total_penalty"`
}
