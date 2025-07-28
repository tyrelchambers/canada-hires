package types

// JobData represents a scraped job listing
type JobData struct {
	JobTitle  string `json:"jobTitle"`
	Business  string `json:"business"`
	Salary    string `json:"salary"`
	Location  string `json:"location"`
	JobURL    string `json:"jobUrl"`
	Date      string `json:"date"`
	JobBankID string `json:"jobBankId,omitempty"`
}