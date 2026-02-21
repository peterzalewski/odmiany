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

type pastCorpusEntry struct {
	Infinitive string `json:"infinitive"`
	Sg1M       string `json:"sg1m"`
	Sg1F       string `json:"sg1f"`
	Sg2M       string `json:"sg2m"`
	Sg2F       string `json:"sg2f"`
	Sg3M       string `json:"sg3m"`
	Sg3F       string `json:"sg3f"`
	Sg3N       string `json:"sg3n"`
	Pl1V       string `json:"pl1v"`
	Pl1NV      string `json:"pl1nv"`
	Pl2V       string `json:"pl2v"`
	Pl2NV      string `json:"pl2nv"`
	Pl3V       string `json:"pl3v"`
	Pl3NV      string `json:"pl3nv"`
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

func loadPastCorpus(t *testing.T) []pastCorpusEntry {
	t.Helper()
	data, err := os.ReadFile("testdata/verbs_past.json")
	if err != nil {
		t.Fatalf("failed to load past corpus: %v", err)
	}
	var entries []pastCorpusEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("failed to parse past corpus: %v", err)
	}
	return entries
}

func TestCorpusAccuracy(t *testing.T) {
	entries := loadCorpus(t)

	// Group corpus entries by infinitive to handle variants/homographs
	byInfinitive := make(map[string][]PresentTense)
	for _, e := range entries {
		pt := PresentTense{
			Sg1: e.Sg1, Sg2: e.Sg2, Sg3: e.Sg3,
			Pl1: e.Pl1, Pl2: e.Pl2, Pl3: e.Pl3,
		}
		byInfinitive[e.Infinitive] = append(byInfinitive[e.Infinitive], pt)
	}

	var passed, failed, noMatch int
	failures := make(map[string]int) // pattern -> count

	for infinitive, corpusParadigms := range byInfinitive {
		paradigms, err := ConjugatePresent(infinitive)
		if err != nil {
			noMatch++
			pattern := classifyFailure(infinitive, "no_match")
			failures[pattern]++
			continue
		}

		// Check if ANY of our paradigms matches ANY corpus paradigm
		anyMatch := false
		for _, ourP := range paradigms {
			for _, corpusP := range corpusParadigms {
				if ourP.PresentTense.Equals(corpusP) {
					anyMatch = true
					break
				}
			}
			if anyMatch {
				break
			}
		}

		if anyMatch {
			passed++
		} else {
			failed++
			pattern := classifyFailure(infinitive, describeError(infinitive, corpusParadigms[0], paradigms[0].PresentTense))
			failures[pattern]++
		}
	}

	total := len(byInfinitive)
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

		paradigms, err := ConjugatePresent(inf)
		if err != nil {
			t.Logf("%s: no match (expected: %s, %s, %s...)", inf, e.Sg1, e.Sg2, e.Sg3)
			continue
		}

		// Check if any paradigm matches
		anyMatch := false
		for _, p := range paradigms {
			if p.PresentTense.Equals(expected) {
				anyMatch = true
				break
			}
		}

		got := paradigms[0].PresentTense
		if anyMatch {
			t.Logf("%s: ✓ %s, %s, %s...", inf, got.Sg1, got.Sg2, got.Sg3)
		} else {
			t.Logf("%s: ✗ got %s, %s, %s; want %s, %s, %s",
				inf, got.Sg1, got.Sg2, got.Sg3, e.Sg1, e.Sg2, e.Sg3)
		}
	}
}

func TestCorpusPastAccuracy(t *testing.T) {
	entries := loadPastCorpus(t)

	// Group corpus entries by infinitive to handle variants
	byInfinitive := make(map[string][]PastTense)
	for _, e := range entries {
		pt := PastTense{
			Sg1M: e.Sg1M, Sg1F: e.Sg1F,
			Sg2M: e.Sg2M, Sg2F: e.Sg2F,
			Sg3M: e.Sg3M, Sg3F: e.Sg3F, Sg3N: e.Sg3N,
			Pl1V: e.Pl1V, Pl1NV: e.Pl1NV,
			Pl2V: e.Pl2V, Pl2NV: e.Pl2NV,
			Pl3V: e.Pl3V, Pl3NV: e.Pl3NV,
		}
		byInfinitive[e.Infinitive] = append(byInfinitive[e.Infinitive], pt)
	}

	var passed, failed, noMatch int
	failures := make(map[string]int) // pattern -> count

	for infinitive, corpusParadigms := range byInfinitive {
		paradigms, err := ConjugatePast(infinitive)
		if err != nil {
			noMatch++
			pattern := classifyFailure(infinitive, "no_match")
			failures[pattern]++
			continue
		}

		// Check if ANY of our paradigms matches ANY corpus paradigm
		anyMatch := false
		for _, ourP := range paradigms {
			for _, corpusP := range corpusParadigms {
				if ourP.PastTense.Equals(corpusP) {
					anyMatch = true
					break
				}
			}
			if anyMatch {
				break
			}
		}

		if anyMatch {
			passed++
		} else {
			failed++
			pattern := classifyFailure(infinitive, describePastError(infinitive, corpusParadigms[0], paradigms[0].PastTense))
			failures[pattern]++
		}
	}

	total := len(byInfinitive)
	accuracy := float64(passed) / float64(total) * 100

	t.Logf("Past tense corpus accuracy: %.2f%% (%d/%d passed, %d failed, %d no match)",
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
}

type verbalNounCorpusEntry struct {
	Infinitive string `json:"infinitive"`
	VerbalNoun string `json:"verbal_noun"`
}

func loadVerbalNounCorpus(t *testing.T) []verbalNounCorpusEntry {
	t.Helper()
	data, err := os.ReadFile("testdata/verbs_verbal_noun.json")
	if err != nil {
		t.Fatalf("failed to load verbal noun corpus: %v", err)
	}
	var entries []verbalNounCorpusEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("failed to parse verbal noun corpus: %v", err)
	}
	return entries
}

func TestCorpusVerbalNounAccuracy(t *testing.T) {
	entries := loadVerbalNounCorpus(t)

	// Group corpus entries by infinitive (some verbs have multiple valid forms)
	byInfinitive := make(map[string][]string)
	for _, e := range entries {
		byInfinitive[e.Infinitive] = append(byInfinitive[e.Infinitive], e.VerbalNoun)
	}

	var passed, failed, noMatch int
	failures := make(map[string]int)

	for infinitive, corpusForms := range byInfinitive {
		predicted, err := VerbalNoun(infinitive)
		if err != nil {
			noMatch++
			pattern := classifyFailure(infinitive, "no_match")
			failures[pattern]++
			continue
		}

		// Check if ANY corpus form appears in our predicted forms
		anyMatch := false
		for _, corpusForm := range corpusForms {
			for _, predForm := range predicted {
				if corpusForm == predForm {
					anyMatch = true
					break
				}
			}
			if anyMatch {
				break
			}
		}

		if anyMatch {
			passed++
		} else {
			failed++
			desc := fmt.Sprintf("want %s got %s", corpusForms[0], predicted[0])
			pattern := classifyFailure(infinitive, desc)
			failures[pattern]++
		}
	}

	total := len(byInfinitive)
	accuracy := float64(passed) / float64(total) * 100

	t.Logf("Verbal noun corpus accuracy: %.2f%% (%d/%d passed, %d failed, %d no match)",
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
}

// describePastError returns a short description of how the past conjugation differs.
func describePastError(infinitive string, expected, got PastTense) string {
	var diffs []string

	if expected.Sg3M != got.Sg3M {
		diffs = append(diffs, fmt.Sprintf("3sgM: want %s got %s", ending(expected.Sg3M), ending(got.Sg3M)))
	}
	if expected.Sg3F != got.Sg3F {
		diffs = append(diffs, fmt.Sprintf("3sgF: want %s got %s", ending(expected.Sg3F), ending(got.Sg3F)))
	}
	if expected.Pl3V != got.Pl3V {
		diffs = append(diffs, fmt.Sprintf("3plV: want %s got %s", ending(expected.Pl3V), ending(got.Pl3V)))
	}

	if len(diffs) > 2 {
		return diffs[0] // Just show the first difference
	}
	return strings.Join(diffs, "; ")
}

// TestSamplePastVerbs tests specific verbs to understand past tense patterns.
func TestSamplePastVerbs(t *testing.T) {
	entries := loadPastCorpus(t)

	// Build a map for quick lookup
	corpus := make(map[string]pastCorpusEntry)
	for _, e := range entries {
		corpus[e.Infinitive] = e
	}

	samples := []string{
		"czytać", "pisać", "robić", "brać", "być", "mieć",
		"pracować", "umieć", "myć", "żyć", "pić",
		"nieść", "móc", "chcieć", "iść", "wziąć",
		"ciągnąć", "kopnąć",
	}

	for _, inf := range samples {
		e, ok := corpus[inf]
		if !ok {
			t.Logf("%s: not in corpus", inf)
			continue
		}

		expected := PastTense{
			Sg1M: e.Sg1M, Sg1F: e.Sg1F,
			Sg2M: e.Sg2M, Sg2F: e.Sg2F,
			Sg3M: e.Sg3M, Sg3F: e.Sg3F, Sg3N: e.Sg3N,
			Pl1V: e.Pl1V, Pl1NV: e.Pl1NV,
			Pl2V: e.Pl2V, Pl2NV: e.Pl2NV,
			Pl3V: e.Pl3V, Pl3NV: e.Pl3NV,
		}

		paradigms, err := ConjugatePast(inf)
		if err != nil {
			t.Logf("%s: no match (expected: %s, %s, %s...)", inf, e.Sg3M, e.Sg3F, e.Sg3N)
			continue
		}

		// Check if any paradigm matches
		anyMatch := false
		for _, p := range paradigms {
			if p.PastTense.Equals(expected) {
				anyMatch = true
				break
			}
		}

		got := paradigms[0].PastTense
		if anyMatch {
			t.Logf("%s: ✓ %s, %s, %s, %s, %s", inf, got.Sg3M, got.Sg3F, got.Sg3N, got.Pl3V, got.Pl3NV)
		} else {
			t.Logf("%s: ✗ got %s/%s/%s want %s/%s/%s",
				inf, got.Sg3M, got.Sg3F, got.Sg3N, e.Sg3M, e.Sg3F, e.Sg3N)
		}
	}
}
