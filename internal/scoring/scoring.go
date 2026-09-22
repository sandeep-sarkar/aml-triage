package scoring

import (
	"fmt"

	"github.com/sandeep/aml-triage/internal/alert"
	"github.com/sandeep/aml-triage/internal/policy"
)

type ScoreEngine struct {
	Policy policy.PolicyEngine
}

func (s ScoreEngine) Score(a alert.Alert) (alert.Disposition, error) {
	if len(a.Transactions) == 0 {
		return alert.Disposition{}, fmt.Errorf("alert %s has no transactions to score", a.ID)
	}

	matches := s.Policy.Assess(a)
	if len(matches) == 0 {
		return alert.Disposition{AlertID: a.ID, Decision: "close"}, nil
	}
	//Todo: combine matches into one Reasoning string and one PolicyRef
	// (e.g. join reasoning with "; ", join policy refs with ", ")
	// return Decision: "escalate"
}
