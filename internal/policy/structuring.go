package policy

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sandeep/aml-triage/internal/alert"
)

/*
Implements # POL-STRUCT-01: Structuring and threshold avoidance
*/

type Structure struct {
	maxAmount float64
	txnWindow time.Duration
	txnCount  int
}

func Init() Structure {
	return Structure{
		maxAmount: 10000,
		txnWindow: 72 * time.Hour,
		txnCount:  3,
	}
}

func (s Structure) Name() string {
	return "POL-STRUCT-01: Structuring and threshold avoidance"
}

func (s Structure) Evaluate(a alert.Alert) (triggered bool, reasoning string, policyRef string) {
	//filter out transactions that are less than 10000
	subtransactions := make([]alert.Transaction, 0)

	for _, transaction := range a.Transactions {
		if transaction.AmountGBP <= s.maxAmount {
			subtransactions = append(subtransactions, transaction)
		}
	}

	sort.Slice(subtransactions, func(i, j int) bool {
		return subtransactions[i].Timestamp.Before(subtransactions[j].Timestamp)
	})

	lp, rp := 0, 0

	for ; lp < len(subtransactions)-s.txnCount+1; lp++ {
		txns := make([]alert.Transaction, 0)
		startTxn := subtransactions[lp].Timestamp
		endTxnWindow := startTxn.Add(s.txnWindow)

		for rp = lp; rp < len(subtransactions); rp++ {
			if subtransactions[rp].Timestamp.After(endTxnWindow) {
				break
			}
			txns = append(txns, subtransactions[rp])
		}

		if rp-lp >= s.txnCount {
			reasoning = s.GenerateStatement(txns)
			return true, reasoning, "POL-STRUCT-01"
		}
	}
	return false, "", ""
}

func (s Structure) GenerateStatement(txns []alert.Transaction) string {
	stmt := fmt.Sprintf("%d transactions at or under %.2f within %v hours:",
		len(txns), s.maxAmount, s.txnWindow.Hours())
	for i := 0; i < len(txns); i++ {
		stmt = fmt.Sprintf("%s %s %.2f (%s),", stmt, txns[i].ID,
			txns[i].AmountGBP, txns[i].Timestamp.Format("2006-01-02 15:04"))
	}
	stmt = strings.TrimSuffix(stmt, ",")
	stmt = fmt.Sprintf("%s. Span: %v hours", stmt,
		txns[len(txns)-1].Timestamp.Sub(txns[0].Timestamp).Hours())
	return stmt
}
