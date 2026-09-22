package alert

import "time"

type Transaction struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	AmountGBP float64   `json:"amount_gbp"`
	Direction string    `json:"direction"` // "inbound" | "outbound"
	Counterparty string `json:"counterparty"`
	Channel   string    `json:"channel"` // "wire" | "card" | "cash" | "crypto"
}

type Customer struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AccountType string `json:"account_type"`
	RiskRating  string `json:"risk_rating"` // "low" | "medium" | "high"
	KYCDate     string `json:"kyc_date"`
}

type Alert struct {
	ID           string        `json:"id"`
	Typology     string        `json:"typology"` // e.g. "structuring"
	Customer     Customer      `json:"customer"`
	Transactions []Transaction `json:"transactions"`

	// Ground truth for eval, not shown to the model
	TrueDisposition string `json:"true_disposition,omitempty"` // "escalate" | "close"
	TruePolicyRef   string `json:"true_policy_ref,omitempty"`
}

// Disposition is what both the baseline scorer and the LLM must produce.
type Disposition struct {
	AlertID    string `json:"alert_id"`
	Decision   string `json:"decision"` // "escalate" | "close"
	Reasoning  string `json:"reasoning"`
	PolicyRef  string `json:"policy_ref,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}
