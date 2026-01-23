package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"petezalew.ski/odmiany/pkg/verb"
)

type corpusEntry struct {
	Infinitive string `json:"infinitive"`
	Sg1        string `json:"sg1"`
	Sg2        string `json:"sg2"`
	Sg3        string `json:"sg3"`
	Pl1        string `json:"pl1"`
	Pl2        string `json:"pl2"`
	Pl3        string `json:"pl3"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: conjugate <prefix|infinitive>")
		fmt.Println("  Search corpus for verbs matching prefix and show conjugations")
		fmt.Println("  If exact infinitive given, shows detailed comparison")
		os.Exit(1)
	}

	query := os.Args[1]

	// Load corpus for comparison
	data, err := os.ReadFile("pkg/verb/testdata/verbs.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading corpus: %v\n", err)
		os.Exit(1)
	}

	var entries []corpusEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing corpus: %v\n", err)
		os.Exit(1)
	}

	// Build corpus map
	corpus := make(map[string]corpusEntry)
	for _, e := range entries {
		corpus[e.Infinitive] = e
	}

	// Check for exact match first
	if e, ok := corpus[query]; ok {
		showDetailed(query, e)
		return
	}

	// Search by prefix or suffix
	var matches []corpusEntry
	for _, e := range entries {
		if strings.HasPrefix(e.Infinitive, query) || strings.HasSuffix(e.Infinitive, query) {
			matches = append(matches, e)
		}
	}

	if len(matches) == 0 {
		// Try conjugating anyway (might not be in corpus)
		fmt.Printf("No corpus matches for %q, attempting conjugation:\n\n", query)
		got, err := verb.ConjugatePresent(query)
		if err != nil {
			fmt.Printf("  %s: NO MATCH (%v)\n", query, err)
		} else {
			printParadigm(query, got)
		}
		return
	}

	fmt.Printf("Found %d matches for %q:\n\n", len(matches), query)
	for _, e := range matches {
		showComparison(e)
	}
}

func showDetailed(infinitive string, e corpusEntry) {
	fmt.Printf("=== %s ===\n\n", infinitive)

	expected := verb.PresentTense{
		Sg1: e.Sg1, Sg2: e.Sg2, Sg3: e.Sg3,
		Pl1: e.Pl1, Pl2: e.Pl2, Pl3: e.Pl3,
	}

	got, err := verb.ConjugatePresent(infinitive)

	fmt.Println("Expected (corpus):")
	printParadigm("", expected)

	fmt.Println("\nGot (heuristic):")
	if err != nil {
		fmt.Printf("  NO MATCH: %v\n", err)
	} else {
		printParadigm("", got)
	}

	if err == nil {
		fmt.Println("\nComparison:")
		compare("Sg1", expected.Sg1, got.Sg1)
		compare("Sg2", expected.Sg2, got.Sg2)
		compare("Sg3", expected.Sg3, got.Sg3)
		compare("Pl1", expected.Pl1, got.Pl1)
		compare("Pl2", expected.Pl2, got.Pl2)
		compare("Pl3", expected.Pl3, got.Pl3)
	}
}

func showComparison(e corpusEntry) {
	got, err := verb.ConjugatePresent(e.Infinitive)

	status := "✓"
	if err != nil {
		status = "✗ NO_MATCH"
	} else {
		expected := verb.PresentTense{
			Sg1: e.Sg1, Sg2: e.Sg2, Sg3: e.Sg3,
			Pl1: e.Pl1, Pl2: e.Pl2, Pl3: e.Pl3,
		}
		if !got.Equals(expected) {
			status = "✗ WRONG"
		}
	}

	if err != nil {
		fmt.Printf("%-20s %s (want: %s)\n", e.Infinitive, status, e.Sg1)
	} else {
		fmt.Printf("%-20s %s got=%-15s want=%s\n", e.Infinitive, status, got.Sg1, e.Sg1)
	}
}

func printParadigm(label string, p verb.PresentTense) {
	if label != "" {
		fmt.Printf("%s:\n", label)
	}
	fmt.Printf("  Sg: %s, %s, %s\n", p.Sg1, p.Sg2, p.Sg3)
	fmt.Printf("  Pl: %s, %s, %s\n", p.Pl1, p.Pl2, p.Pl3)
}

func compare(form, expected, got string) {
	if expected == got {
		fmt.Printf("  %s: ✓ %s\n", form, got)
	} else {
		fmt.Printf("  %s: ✗ got %q, want %q\n", form, got, expected)
	}
}
