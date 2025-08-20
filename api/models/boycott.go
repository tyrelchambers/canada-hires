package models

import "time"

type Boycott struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	BusinessName    string    `json:"business_name" db:"business_name"`
	BusinessAddress *string   `json:"business_address" db:"business_address"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type BoycottStats struct {
	BusinessName    string `json:"business_name" db:"business_name"`
	BusinessAddress string `json:"business_address" db:"business_address"`
	BoycottCount    int    `json:"boycott_count" db:"boycott_count"`
}