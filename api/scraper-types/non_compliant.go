package types

// NonCompliantEmployerData represents a scraped non-compliant employer record
type NonCompliantEmployerData struct {
	BusinessOperatingName string   `json:"businessOperatingName"`
	BusinessLegalName     string   `json:"businessLegalName"`
	Address               string   `json:"address"`
	ReasonCodes           []string `json:"reasonCodes"` // Array of reason codes like ["5", "6", "15"]
	DateOfFinalDecision   string   `json:"dateOfFinalDecision"`
	PenaltyAmount         int      `json:"penaltyAmount"`
	PenaltyCurrency       string   `json:"penaltyCurrency"`
	Status                string   `json:"status"`
}