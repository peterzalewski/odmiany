package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
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

type failure struct {
	Infinitive string
	Freq       int
	Got        string
	Want       string
	NoMatch    bool
}

func main() {
	// Load frequency data from OpenSubtitles (hermitdave/FrequencyWords)
	freqMap := loadFrequency("pkg/verb/testdata/pl_freq.txt")

	// Load verb corpus
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

		// Get frequency - check infinitive and all conjugated forms
		freq := getVerbFrequency(freqMap, e)

		if err != nil {
			failures = append(failures, failure{
				Infinitive: e.Infinitive,
				Freq:       freq,
				Got:        "",
				Want:       e.Sg1,
				NoMatch:    true,
			})
		} else if !got.Equals(expected) {
			failures = append(failures, failure{
				Infinitive: e.Infinitive,
				Freq:       freq,
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
		fmt.Printf("%-20s freq=%9d  %-10s got=%-15s want=%s\n",
			f.Infinitive, f.Freq, status, f.Got, f.Want)
	}

	fmt.Fprintf(os.Stderr, "\nTotal failures: %d\n", len(failures))
	fmt.Fprintf(os.Stderr, "Frequency source: OpenSubtitles 2018 (hermitdave/FrequencyWords)\n")
}

// loadFrequency loads word frequency data from hermitdave format: "word count"
func loadFrequency(path string) map[string]int {
	freq := make(map[string]int)

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load frequency data: %v\n", err)
		return freq
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			word := parts[0]
			count, _ := strconv.Atoi(parts[1])
			freq[word] = count
		}
	}

	return freq
}

// getVerbFrequency returns the highest frequency among the infinitive and all conjugated forms
func getVerbFrequency(freqMap map[string]int, e corpusEntry) int {
	maxFreq := 0

	// Check infinitive
	if f, ok := freqMap[e.Infinitive]; ok && f > maxFreq {
		maxFreq = f
	}

	// Check all conjugated forms (these appear more often in subtitles)
	forms := []string{e.Sg1, e.Sg2, e.Sg3, e.Pl1, e.Pl2, e.Pl3}
	for _, form := range forms {
		if f, ok := freqMap[form]; ok && f > maxFreq {
			maxFreq = f
		}
	}

	return maxFreq
}
