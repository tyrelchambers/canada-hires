package services

import (
	"regexp"
	"strings"
)

// Canadian postal code regex pattern: A1A 1A1 or A1A1A1
var postalCodeRegex = regexp.MustCompile(`([A-Za-z]\d[A-Za-z])\s*(\d[A-Za-z]\d)`)

// PostalCodeService handles postal code extraction and validation
type PostalCodeService interface {
	ExtractPostalCode(address string) string
	ValidatePostalCode(postalCode string) bool
	FormatPostalCode(postalCode string) string
}

type postalCodeService struct{}

func NewPostalCodeService() PostalCodeService {
	return &postalCodeService{}
}

// ExtractPostalCode extracts a Canadian postal code from an address string
func (p *postalCodeService) ExtractPostalCode(address string) string {
	if address == "" {
		return ""
	}

	// Clean the address - remove common prefixes and extra spaces
	cleaned := strings.TrimSpace(address)
	
	// Find postal code pattern in the address
	matches := postalCodeRegex.FindStringSubmatch(cleaned)
	if len(matches) >= 3 {
		// Format as A1A 1A1
		return strings.ToUpper(matches[1] + " " + matches[2])
	}
	
	return ""
}

// ValidatePostalCode checks if a postal code follows Canadian format
func (p *postalCodeService) ValidatePostalCode(postalCode string) bool {
	if postalCode == "" {
		return false
	}
	
	// Remove spaces and check format
	cleaned := strings.ReplaceAll(strings.ToUpper(postalCode), " ", "")
	if len(cleaned) != 6 {
		return false
	}
	
	// Check pattern: Letter-Digit-Letter-Digit-Letter-Digit
	pattern := regexp.MustCompile(`^[A-Z]\d[A-Z]\d[A-Z]\d$`)
	return pattern.MatchString(cleaned)
}

// FormatPostalCode ensures consistent formatting (A1A 1A1)
func (p *postalCodeService) FormatPostalCode(postalCode string) string {
	if postalCode == "" {
		return ""
	}
	
	// Remove all spaces and convert to uppercase
	cleaned := strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))
	
	if len(cleaned) == 6 && p.ValidatePostalCode(cleaned) {
		// Insert space in the middle: A1A1A1 -> A1A 1A1
		return cleaned[:3] + " " + cleaned[3:]
	}
	
	return postalCode // Return original if invalid
}