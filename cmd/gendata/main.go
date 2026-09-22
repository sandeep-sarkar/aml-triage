// Command gendata writes synthetic labelled alerts to data/alerts.json.
//
// Usage:
//
//	go run ./cmd/gendata -n 200 -seed 42 -out data/alerts.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sandeep/aml-triage/internal/gendata"
)

func main() {
	n := flag.Int("n", 200, "number of alerts to generate")
	seed := flag.Int64("seed", 42, "random seed, for reproducibility")
	out := flag.String("out", "data/alerts.json", "output file path")
	flag.Parse()

	alerts := gendata.GenerateAlerts(*seed, *n)

	f, err := os.Create(*out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating output file:", err)
		os.Exit(1)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(alerts); err != nil {
		fmt.Fprintln(os.Stderr, "error encoding alerts:", err)
		os.Exit(1)
	}

	positive := 0
	for _, a := range alerts {
		if a.TrueDisposition == "escalate" {
			positive++
		}
	}
	fmt.Printf("wrote %d alerts to %s (%d escalate, %d close)\n", len(alerts), *out, positive, len(alerts)-positive)
}
