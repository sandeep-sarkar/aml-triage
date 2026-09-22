// Package gendata creates synthetic customers, transactions, and labelled
// alerts for the structuring typology, plus a handful of fake policy
// documents. It uses only the standard library so it runs with no extra
// dependencies.
package gendata

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sandeep/aml-triage/internal/alert"
)

var firstNames = []string{
	"James", "Priya", "Mohammed", "Sarah", "Chen", "Fatima", "Liam",
	"Aisha", "Daniel", "Ngozi", "Tom", "Elena", "Raj", "Emily", "Yuki",
}

var lastNames = []string{
	"Smith", "Patel", "Khan", "Johnson", "Wang", "Brown", "Okafor",
	"Garcia", "Kim", "Müller", "Singh", "Davies", "Rossi", "Novak",
}

var counterparties = []string{
	"ABC Trading Ltd", "Fast Cash Exchange", "Riverside Motors",
	"Global Freight Co", "Green Leaf Holdings", "Unity Retail",
	"North Star Logistics", "Harbor View Consulting", "Private individual",
}

var accountTypes = []string{"personal current", "business current", "savings"}
var riskRatings = []string{"low", "medium", "high"}

func randChoice(r *rand.Rand, s []string) string {
	return s[r.Intn(len(s))]
}

func randCustomer(r *rand.Rand, id int) alert.Customer {
	return alert.Customer{
		ID:          fmt.Sprintf("CUST-%04d", id),
		Name:        randChoice(r, firstNames) + " " + randChoice(r, lastNames),
		AccountType: randChoice(r, accountTypes),
		RiskRating:  randChoice(r, riskRatings),
		KYCDate:     time.Date(2021+r.Intn(4), time.Month(1+r.Intn(12)), 1+r.Intn(28), 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
	}
}

// genStructuringPositive creates a transaction set that IS structuring:
// several deposits just under a £10,000 reporting threshold, in a short
// window.
func genStructuringPositive(r *rand.Rand, start time.Time) []alert.Transaction {
	n := 3 + r.Intn(3) // 3-5 transactions
	txns := make([]alert.Transaction, 0, n)
	for i := 0; i < n; i++ {
		amount := 8500 + r.Float64()*1400 // 8500-9900, just under 10k
		txns = append(txns, alert.Transaction{
			ID:           fmt.Sprintf("TXN-%08d", r.Int31()),
			Timestamp:    start.Add(time.Duration(i) * time.Duration(4+r.Intn(20)) * time.Hour),
			AmountGBP:    round2(amount),
			Direction:    "inbound",
			Counterparty: randChoice(r, counterparties),
			Channel:      "cash",
		})
	}
	return txns
}

// genNormal creates an ordinary transaction pattern: a few transactions,
// varied amounts, no threshold-avoidance pattern.
func genNormal(r *rand.Rand, start time.Time) []alert.Transaction {
	n := 1 + r.Intn(4)
	txns := make([]alert.Transaction, 0, n)
	for i := 0; i < n; i++ {
		amount := 50 + r.Float64()*4000
		dir := "outbound"
		if r.Intn(2) == 0 {
			dir = "inbound"
		}
		txns = append(txns, alert.Transaction{
			ID:           fmt.Sprintf("TXN-%08d", r.Int31()),
			Timestamp:    start.Add(time.Duration(i) * time.Duration(1+r.Intn(96)) * time.Hour),
			AmountGBP:    round2(amount),
			Direction:    dir,
			Counterparty: randChoice(r, counterparties),
			Channel:      []string{"wire", "card", "cash"}[r.Intn(3)],
		})
	}
	return txns
}

func round2(f float64) float64 {
	return float64(int(f*100)) / 100
}

// GenerateAlerts produces n labelled alerts, roughly balanced between
// structuring-positive and normal cases.
func GenerateAlerts(seed int64, n int) []alert.Alert {
	r := rand.New(rand.NewSource(seed))
	alerts := make([]alert.Alert, 0, n)

	for i := 0; i < n; i++ {
		cust := randCustomer(r, i+1)
		start := time.Date(2026, time.Month(1+r.Intn(9)), 1+r.Intn(27), 9, 0, 0, 0, time.UTC)

		isPositive := r.Intn(2) == 0
		var txns []alert.Transaction
		var trueDisposition, truePolicyRef string

		if isPositive {
			txns = genStructuringPositive(r, start)
			trueDisposition = "escalate"
			truePolicyRef = "POL-STRUCT-01"
		} else {
			txns = genNormal(r, start)
			trueDisposition = "close"
			truePolicyRef = ""
		}

		alerts = append(alerts, alert.Alert{
			ID:              fmt.Sprintf("ALERT-%04d", i+1),
			Typology:        "structuring",
			Customer:        cust,
			Transactions:    txns,
			TrueDisposition: trueDisposition,
			TruePolicyRef:   truePolicyRef,
		})
	}
	return alerts
}
