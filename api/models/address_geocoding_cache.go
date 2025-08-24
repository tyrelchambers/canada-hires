package models

import (
	"time"
)

type AddressGeocodingCache struct {
	ID                int       `json:"id" db:"id"`
	Address           string    `json:"address" db:"address"`
	NormalizedAddress string    `json:"normalized_address" db:"normalized_address"`
	Latitude          float64   `json:"latitude" db:"latitude"`
	Longitude         float64   `json:"longitude" db:"longitude"`
	Confidence        *float64  `json:"confidence" db:"confidence"`
	GeocodedAt        time.Time `json:"geocoded_at" db:"geocoded_at"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}