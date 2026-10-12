package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sandeep/aml-triage/internal/alert"
	"github.com/sandeep/aml-triage/internal/scoring"
)

func main() {
	fmt.Println("Starting eval harness")
	err := start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Encountered Error: ", err)
		os.Exit(1)
	}
}

func start() error {
	freader, err := os.Open("data/alerts.json")
	if err != nil {
		return fmt.Errorf("Error in opening file: %w", err.Error())
	}
	defer freader.Close()

	alertDecoder := json.NewDecoder(freader)
	if _, err := alertDecoder.Token(); err != nil { // read [
		return fmt.Errorf("Error in parsing alerts.json: %w", err)
	}

	scoreEngine := scoring.NewScoreEngine()

	for alertDecoder.More() {
		var a alert.Alert
		if err := alertDecoder.Decode(&a); err != nil {
			return fmt.Errorf("Error in parsing alert json: %w", err)
		}
		disposition, err := scoreEngine.Score(a)
		if err != nil {
			return fmt.Errorf("Error in alert disposition: %w", err)
		}
		fmt.Printf("Alert_id: %s, Decision: %s, Reasoning: %s, PolicyRef: %s\n",
			disposition.AlertID, disposition.Decision, disposition.Reasoning,
			disposition.PolicyRef)
	}
	return nil
}
