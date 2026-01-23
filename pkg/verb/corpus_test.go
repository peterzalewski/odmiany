package verb

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

type corpusEntry struct {
	Infinitive string `json:"infinitive"`
	Sg1        string `json:"sg1"`
	Sg2        string `json:"sg2"`
	Sg3        string `json:"sg3"`
	Pl1        string `json:"pl1"`
	Pl2        string `json:"pl2"`
	Pl3        string `json:"pl3"`
	Aspect     string `json:"aspect"`
}

func loadCorpus(t *testing.T) []corpusEntry {
	t.Helper()
	data, err := os.ReadFile("testdata/verbs.json")
	if err != nil {
		t.Fatalf("failed to load corpus: %v", err)
	}
	var entries []corpusEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("failed to parse corpus: %v", err)
	}
	return entries
}

func TestCorpusAccuracy(t *testing.T) {
	entries := loadCorpus(t)

	var passed, failed, noMatch int
	failures := make(map[string]int) // pattern -> count

	for _, e := range entries {
		expected := PresentTense{
			Sg1: e.Sg1, Sg2: e.Sg2, Sg3: e.Sg3,
			Pl1: e.Pl1, Pl2: e.Pl2, Pl3: e.Pl3,
		}

		got, err := ConjugatePresent(e.Infinitive)
		if err != nil {
			noMatch++
			pattern := classifyFailure(e.Infinitive, "no_match")
			failures[pattern]++
			continue
		}

		if got.Equals(expected) {
			passed++
		} else {
			failed++
			pattern := classifyFailure(e.Infinitive, describeError(e.Infinitive, expected, got))
			failures[pattern]++
		}
	}

	total := len(entries)
	accuracy := float64(passed) / float64(total) * 100

	t.Logf("Corpus accuracy: %.2f%% (%d/%d passed, %d failed, %d no match)",
		accuracy, passed, total, failed, noMatch)

	// Print top failure patterns
	type failurePattern struct {
		pattern string
		count   int
	}
	var patterns []failurePattern
	for p, c := range failures {
		patterns = append(patterns, failurePattern{p, c})
	}
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].count > patterns[j].count
	})

	t.Log("\nTop failure patterns:")
	for i, p := range patterns {
		if i >= 20 {
			break
		}
		t.Logf("  %4d: %s", p.count, p.pattern)
	}

	// For now, don't fail the test - we're iterating on heuristics
	// Uncomment this when we want to enforce a threshold:
	// if accuracy < 95.0 {
	// 	t.Errorf("accuracy %.2f%% below threshold 95%%", accuracy)
	// }
}

// classifyFailure returns a pattern string for grouping similar failures.
func classifyFailure(infinitive, errorType string) string {
	// Get the last few characters of the infinitive
	suffix := infinitive
	if len(suffix) > 6 {
		suffix = suffix[len(suffix)-6:]
	}
	return fmt.Sprintf("[%s] %s", suffix, errorType)
}

// describeError returns a short description of how the conjugation differs.
func describeError(infinitive string, expected, got PresentTense) string {
	var diffs []string

	if expected.Sg1 != got.Sg1 {
		// Extract the ending pattern
		diffs = append(diffs, fmt.Sprintf("1sg: want %s got %s", ending(expected.Sg1), ending(got.Sg1)))
	}
	if expected.Sg2 != got.Sg2 {
		diffs = append(diffs, fmt.Sprintf("2sg: want %s got %s", ending(expected.Sg2), ending(got.Sg2)))
	}
	if expected.Sg3 != got.Sg3 {
		diffs = append(diffs, fmt.Sprintf("3sg: want %s got %s", ending(expected.Sg3), ending(got.Sg3)))
	}

	if len(diffs) > 2 {
		return diffs[0] // Just show the first difference
	}
	return strings.Join(diffs, "; ")
}

func ending(form string) string {
	if len(form) <= 3 {
		return form
	}
	return "..." + form[len(form)-3:]
}

// TestSampleVerbs tests specific verbs to understand patterns.
func TestSampleVerbs(t *testing.T) {
	entries := loadCorpus(t)

	// Build a map for quick lookup
	corpus := make(map[string]corpusEntry)
	for _, e := range entries {
		corpus[e.Infinitive] = e
	}

	samples := []string{
		"czytać", "pisać", "robić", "brać", "być", "mieć",
		"pracować", "umieć", "myć", "żyć", "pić",
		"jechać", "nieść", "móc", "chcieć",
	}

	for _, inf := range samples {
		e, ok := corpus[inf]
		if !ok {
			t.Logf("%s: not in corpus", inf)
			continue
		}

		expected := PresentTense{
			Sg1: e.Sg1, Sg2: e.Sg2, Sg3: e.Sg3,
			Pl1: e.Pl1, Pl2: e.Pl2, Pl3: e.Pl3,
		}

		got, err := ConjugatePresent(inf)
		if err != nil {
			t.Logf("%s: no match (expected: %s, %s, %s...)", inf, e.Sg1, e.Sg2, e.Sg3)
			continue
		}

		if got.Equals(expected) {
			t.Logf("%s: ✓ %s, %s, %s...", inf, got.Sg1, got.Sg2, got.Sg3)
		} else {
			t.Logf("%s: ✗ got %s, %s, %s; want %s, %s, %s",
				inf, got.Sg1, got.Sg2, got.Sg3, e.Sg1, e.Sg2, e.Sg3)
		}
	}
}
