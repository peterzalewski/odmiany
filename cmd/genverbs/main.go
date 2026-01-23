package main

import (
	"bufio"
	"compress/bzip2"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

// VerbParadigm holds the present tense forms for a verb.
type VerbParadigm struct {
	Infinitive string `json:"infinitive"`
	Sg1        string `json:"sg1,omitempty"` // ja
	Sg2        string `json:"sg2,omitempty"` // ty
	Sg3        string `json:"sg3,omitempty"` // on/ona/ono
	Pl1        string `json:"pl1,omitempty"` // my
	Pl2        string `json:"pl2,omitempty"` // wy
	Pl3        string `json:"pl3,omitempty"` // oni/one
	Aspect     string `json:"aspect,omitempty"`
}

func main() {
	inputPath := flag.String("input", "data/polish.txt.bz2", "path to polish.txt.bz2")
	flag.Parse()

	f, err := os.Open(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	reader := bzip2.NewReader(f)
	scanner := bufio.NewScanner(reader)

	// Map infinitive -> paradigm
	verbs := make(map[string]*VerbParadigm)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) != 3 {
			continue
		}
		lemma, form, tags := parts[0], parts[1], parts[2]

		// Only interested in finite verbs (present tense)
		if !strings.Contains(tags, "verb:fin:") {
			continue
		}

		// Skip if not imperfective (present tense is imperfective)
		// Actually, we want both for now to see the data
		// Perfective verbs have future meaning in "present" forms

		paradigm, ok := verbs[lemma]
		if !ok {
			paradigm = &VerbParadigm{Infinitive: lemma}
			verbs[lemma] = paradigm
		}

		// Extract aspect
		if strings.Contains(tags, ":imperf") {
			paradigm.Aspect = "imperf"
		} else if strings.Contains(tags, ":perf") && paradigm.Aspect == "" {
			paradigm.Aspect = "perf"
		}

		// Parse person/number from tags like "verb:fin:sg:pri:imperf:refl.nonrefl"
		// Format: verb:fin:NUMBER:PERSON:ASPECT:REFL
		tagParts := strings.Split(tags, ":")
		if len(tagParts) < 4 {
			continue
		}

		number := tagParts[2] // sg or pl
		person := tagParts[3] // pri, sec, ter

		switch {
		case number == "sg" && person == "pri":
			if paradigm.Sg1 == "" {
				paradigm.Sg1 = form
			}
		case number == "sg" && person == "sec":
			if paradigm.Sg2 == "" {
				paradigm.Sg2 = form
			}
		case number == "sg" && person == "ter":
			if paradigm.Sg3 == "" {
				paradigm.Sg3 = form
			}
		case number == "pl" && person == "pri":
			if paradigm.Pl1 == "" {
				paradigm.Pl1 = form
			}
		case number == "pl" && person == "sec":
			if paradigm.Pl2 == "" {
				paradigm.Pl2 = form
			}
		case number == "pl" && person == "ter":
			if paradigm.Pl3 == "" {
				paradigm.Pl3 = form
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scan: %v\n", err)
		os.Exit(1)
	}

	// Output as JSON array
	var paradigms []VerbParadigm
	for _, p := range verbs {
		// Only include verbs with complete paradigms
		if p.Sg1 != "" && p.Sg2 != "" && p.Sg3 != "" &&
			p.Pl1 != "" && p.Pl2 != "" && p.Pl3 != "" {
			paradigms = append(paradigms, *p)
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(paradigms); err != nil {
		fmt.Fprintf(os.Stderr, "encode: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Extracted %d complete verb paradigms\n", len(paradigms))
}
