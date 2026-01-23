package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

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

// Approximate frequency rankings for common Polish verbs
var freq = map[string]int{
	// Top tier (extremely common)
	"być": 100, "mieć": 99, "móc": 98, "musieć": 97, "chcieć": 96,
	"wiedzieć": 95, "mówić": 94, "robić": 93, "iść": 92, "dać": 91,
	"widzieć": 90, "stać": 89, "myśleć": 88, "brać": 87, "jechać": 86,
	"pisać": 85, "czytać": 84, "jeść": 83, "pić": 82, "spać": 81,

	// High frequency
	"słyszeć": 80, "patrzeć": 79, "wziąć": 78, "rozumieć": 77, "leżeć": 76,
	"siedzieć": 75, "prosić": 74, "pytać": 73, "trzymać": 72, "czekać": 71,
	"znać": 70, "żyć": 69, "grać": 68, "pracować": 67, "mieszkać": 66,
	"kochać": 65, "wierzyć": 64, "bać": 63, "nieść": 62, "płynąć": 61,

	// Common prefixed forms
	"przejść": 60, "wejść": 59, "zejść": 58, "dojść": 57, "wyjść": 56,
	"przyjść": 55, "dostać": 54, "zostać": 53, "przestać": 52, "powstać": 51,
	"zobaczyć": 50, "powiedzieć": 49, "zrobić": 48, "dowiedzieć": 47,
	"poznać": 46, "zacząć": 45, "skończyć": 44, "wrócić": 43,

	// Moderately common
	"biec": 42, "lecieć": 41, "rosnąć": 40, "kroić": 39, "kleić": 38,
	"grzać": 37, "rwać": 36, "prać": 35, "słać": 34, "trzeć": 33,
	"drzeć": 32, "ciąć": 31, "giąć": 30, "paść": 29, "siąść": 28,

	// Less common but still used
	"trzęść": 27, "gnieść": 26, "mieść": 25, "wieść": 24, "pleść": 23,
	"kraść": 22, "kłaść": 21, "żreć": 20, "przeć": 19, "mrzeć": 18,

	// Additional common verbs
	"zdobyć": 50, "przybyć": 48, "nabyć": 46,
}

type failure struct {
	Infinitive string
	Freq       int
	Got        string
	Want       string
	NoMatch    bool
}

func main() {
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

	var failures []failure

	for _, e := range entries {
		expected := verb.PresentTense{
			Sg1: e.Sg1, Sg2: e.Sg2, Sg3: e.Sg3,
			Pl1: e.Pl1, Pl2: e.Pl2, Pl3: e.Pl3,
		}

		got, err := verb.ConjugatePresent(e.Infinitive)

		if err != nil {
			f := getFreq(e.Infinitive)
			failures = append(failures, failure{
				Infinitive: e.Infinitive,
				Freq:       f,
				Got:        "",
				Want:       e.Sg1,
				NoMatch:    true,
			})
		} else if !got.Equals(expected) {
			f := getFreq(e.Infinitive)
			failures = append(failures, failure{
				Infinitive: e.Infinitive,
				Freq:       f,
				Got:        got.Sg1,
				Want:       e.Sg1,
				NoMatch:    false,
			})
		}
	}

	// Sort by frequency (descending)
	sort.Slice(failures, func(i, j int) bool {
		return failures[i].Freq > failures[j].Freq
	})

	// Print results
	for _, f := range failures {
		status := "WRONG"
		if f.NoMatch {
			status = "NO_MATCH"
		}
		fmt.Printf("%-20s freq=%3d  %-10s got=%-15s want=%s\n",
			f.Infinitive, f.Freq, status, f.Got, f.Want)
	}

	fmt.Fprintf(os.Stderr, "\nTotal failures: %d\n", len(failures))
}

func getFreq(infinitive string) int {
	if f, ok := freq[infinitive]; ok {
		return f
	}
	// Check if it's a prefixed form of a known verb
	for base, baseFreq := range freq {
		if len(infinitive) > len(base) &&
			infinitive[len(infinitive)-len(base):] == base {
			return baseFreq - 5
		}
	}
	return 0
}
