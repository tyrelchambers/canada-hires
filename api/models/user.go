package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type VerificationTier string

const (
	VerificationBasic    VerificationTier = "basic"
	VerificationEnhanced VerificationTier = "enhanced"
	VerificationTrusted  VerificationTier = "trusted"
)

type IPAddresses []string

func (ips IPAddresses) Value() (driver.Value, error) {
	return json.Marshal(ips)
}

func (ips *IPAddresses) Scan(value interface{}) error {
	if value == nil {
		*ips = IPAddresses{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, ips)
	case string:
		return json.Unmarshal([]byte(v), ips)
	default:
		*ips = IPAddresses{}
		return nil
	}
}

type User struct {
	ID               string           `json:"id" db:"id"`
	Email            string           `json:"email" db:"email"`
	VerificationTier VerificationTier `json:"verification_tier" db:"verification_tier"`
	EmailDomain      *string          `json:"email_domain" db:"email_domain"`
	IPAddresses      IPAddresses      `json:"ip_addresses" db:"ip_addresses"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	LastActive       time.Time        `json:"last_active" db:"last_active"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
}
