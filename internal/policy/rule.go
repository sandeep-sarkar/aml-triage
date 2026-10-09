package policy

import "github.com/sandeep/aml-triage/internal/alert"

type Match struct {
	RuleName  string
	Reasoning string
	PolicyRef string
}

type Rule interface {
	Name() string
	Evaluate(a alert.Alert) (triggered bool, reasoning string, policyRef string)
}

type PolicyEngine struct {
	Rules []Rule
}

func (e PolicyEngine) Assess(a alert.Alert) []Match {
	// TODO: loop over e.Rules, call Evaluate, collect Matches where triggered == true

	matches := make([]Match, 0)

	for _, rule := range e.Rules {
		triggered, reasoning, policyRef := rule.Evaluate(a)
		if triggered {
			match := Match{
				RuleName:  rule.Name(),
				Reasoning: reasoning,
				PolicyRef: policyRef,
			}
			matches = append(matches, match)
		}
	}

	return matches
}
