// Package main extracts verb paradigms from Polimorf morphological dictionary data.
//
// # Problem This Solves
//
// Polimorf (the source data) contains multiple valid conjugation forms for many verbs:
//   - Homographs: "stać" can mean "to stand" (stoję) or "to become" (stanę)
//   - Variants: "brzmieć" has both standard (brzmię) and colloquial (brzmieję) forms
//
// A naive extraction that takes "first form seen" for each slot creates mixed paradigms
// like: {sg1: brzmieję, sg2: brzmiejesz, sg3: brzmi} - mixing two different patterns!
//
// # Solution
//
// We extract ALL forms for each infinitive, then group them into coherent paradigms
// based on Polish conjugation patterns. Forms belong to the same paradigm if their
// endings are consistent with each other.
//
// # Polish Conjugation Pattern Primer
//
// Polish verbs conjugate in predictable patterns. The 1sg (ja) form determines
// the pattern for all other forms:
//
//	Pattern A (-ę/-esz): piszę → piszesz, pisze, piszemy, piszecie, piszą
//	Pattern B (-ę/-isz): robię → robisz, robi, robimy, robicie, robią
//	Pattern C (-ę/-ysz): uczę → uczysz, uczy, uczymy, uczycie, uczą
//	Pattern D (-am/-asz): czytam → czytasz, czyta, czytamy, czytacie, czytają
//	Pattern E (-em/-esz): umiem → umiesz, umie, umiemy, umiecie, umieją
//	Pattern F (-eję/-ejesz): starzeję → starzejesz, starzeje, starzejemy, starzejecie, starzeją
//
// By checking if forms have compatible endings, we can group them correctly.
package main

import (
	"bufio"
	"compress/bzip2"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

// VerbForm represents a single conjugated form with its grammatical tags.
type VerbForm struct {
	Form   string
	Number string // "sg" or "pl"
	Person string // "pri" (1st), "sec" (2nd), "ter" (3rd)
	Aspect string // "imperf" or "perf"
	Refl   string // reflexivity tag
}

// VerbParadigm holds a complete present tense paradigm.
type VerbParadigm struct {
	Infinitive string `json:"infinitive"`
	Sg1        string `json:"sg1"` // ja
	Sg2        string `json:"sg2"` // ty
	Sg3        string `json:"sg3"` // on/ona/ono
	Pl1        string `json:"pl1"` // my
	Pl2        string `json:"pl2"` // wy
	Pl3        string `json:"pl3"` // oni/one
	Aspect     string `json:"aspect"`
}

// conjugationPattern defines expected ending patterns for a conjugation class.
// If sg1 ends with Sg1Suffix, we expect sg2 to end with Sg2Suffix, etc.
type conjugationPattern struct {
	Name      string
	Sg1Suffix string
	Sg2Suffix string
	Sg3Suffix string
	Pl1Suffix string
	Pl2Suffix string
	Pl3Suffix string
}

// knownPatterns lists the main Polish conjugation patterns.
// Order matters: more specific patterns should come first.
var knownPatterns = []conjugationPattern{
	// -eję/-ejesz pattern (inchoative verbs like starzeć)
	{"eję", "eję", "ejesz", "eje", "ejemy", "ejecie", "eją"},
	// -ję/-jesz pattern (dawać → daję)
	{"ję", "ję", "jesz", "je", "jemy", "jecie", "ją"},
	// -uję/-ujesz pattern (pracować → pracuję)
	{"uję", "uję", "ujesz", "uje", "ujemy", "ujecie", "ują"},
	// -am/-asz pattern (czytać → czytam)
	{"am", "am", "asz", "a", "amy", "acie", "ają"},
	// -em/-esz pattern (umieć → umiem)
	{"em", "em", "esz", "e", "emy", "ecie", "eją"},
	// -ę/-isz pattern (robić → robię)
	{"ę/isz", "ę", "isz", "i", "imy", "icie", "ą"},
	// -ę/-ysz pattern (uczyć → uczę)
	{"ę/ysz", "ę", "ysz", "y", "ymy", "ycie", "ą"},
	// -ę/-esz pattern (pisać → piszę)
	{"ę/esz", "ę", "esz", "e", "emy", "ecie", "ą"},
	// -ę/-iesz pattern (nieść → niosę type verbs)
	{"ę/iesz", "ę", "iesz", "ie", "iemy", "iecie", "ą"},
	// -nę/-niesz pattern (ciągnąć → ciągnę)
	{"nę", "nę", "niesz", "nie", "niemy", "niecie", "ną"},
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

	// Collect ALL forms for each infinitive
	verbForms := make(map[string][]VerbForm)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) != 3 {
			continue
		}
		lemma, form, tags := parts[0], parts[1], parts[2]

		// Only interested in finite verbs (present/future tense)
		if !strings.Contains(tags, "verb:fin:") {
			continue
		}

		// Parse the form
		vf := parseVerbForm(form, tags)
		if vf.Number != "" && vf.Person != "" {
			verbForms[lemma] = append(verbForms[lemma], vf)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scan: %v\n", err)
		os.Exit(1)
	}

	// Extract coherent paradigms from collected forms
	var paradigms []VerbParadigm
	for infinitive, forms := range verbForms {
		extracted := extractCoherentParadigms(infinitive, forms)
		paradigms = append(paradigms, extracted...)
	}

	// Sort for deterministic output
	sort.Slice(paradigms, func(i, j int) bool {
		if paradigms[i].Infinitive != paradigms[j].Infinitive {
			return paradigms[i].Infinitive < paradigms[j].Infinitive
		}
		return paradigms[i].Sg1 < paradigms[j].Sg1
	})

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(paradigms); err != nil {
		fmt.Fprintf(os.Stderr, "encode: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Extracted %d complete verb paradigms from %d infinitives\n",
		len(paradigms), len(verbForms))
}

// parseVerbForm extracts grammatical information from Polimorf tags.
func parseVerbForm(form, tags string) VerbForm {
	vf := VerbForm{Form: form}

	// Tags format: verb:fin:NUMBER:PERSON:ASPECT:REFL
	// Example: verb:fin:sg:pri:imperf:nonrefl
	tagParts := strings.Split(tags, ":")
	if len(tagParts) < 4 {
		return vf
	}

	vf.Number = tagParts[2] // sg or pl
	vf.Person = tagParts[3] // pri, sec, ter

	// Extract aspect
	if strings.Contains(tags, ":imperf") {
		vf.Aspect = "imperf"
	} else if strings.Contains(tags, ":perf") {
		vf.Aspect = "perf"
	}

	// Extract reflexivity (useful for distinguishing some paradigms)
	if strings.Contains(tags, ":refl.nonrefl") {
		vf.Refl = "refl.nonrefl"
	} else if strings.Contains(tags, ":nonrefl") {
		vf.Refl = "nonrefl"
	} else if strings.Contains(tags, ":refl") {
		vf.Refl = "refl"
	}

	return vf
}

// extractCoherentParadigms groups forms into coherent paradigms based on ending patterns.
func extractCoherentParadigms(infinitive string, forms []VerbForm) []VerbParadigm {
	// Group forms by slot (person+number)
	bySlot := make(map[string][]VerbForm)
	for _, f := range forms {
		slot := f.Number + ":" + f.Person
		bySlot[slot] = append(bySlot[slot], f)
	}

	// Get all sg1 forms - these determine the paradigms
	sg1Forms := bySlot["sg:pri"]
	if len(sg1Forms) == 0 {
		return nil
	}

	// For each sg1 form, try to build a complete paradigm with compatible forms
	var paradigms []VerbParadigm
	usedForms := make(map[string]bool) // track which forms we've used

	for _, sg1 := range sg1Forms {
		if usedForms[sg1.Form] {
			continue
		}

		// Find the conjugation pattern for this sg1
		pattern := findPattern(sg1.Form)
		if pattern == nil {
			// Unknown pattern - skip for now
			continue
		}

		// Try to find compatible forms for each slot
		paradigm := VerbParadigm{
			Infinitive: infinitive,
			Sg1:        sg1.Form,
			Aspect:     sg1.Aspect,
		}

		// Find sg2
		if sg2 := findCompatibleForm(bySlot["sg:sec"], sg1, pattern.Sg1Suffix, pattern.Sg2Suffix); sg2 != "" {
			paradigm.Sg2 = sg2
		} else {
			continue // incomplete paradigm
		}

		// Find sg3
		if sg3 := findCompatibleForm(bySlot["sg:ter"], sg1, pattern.Sg1Suffix, pattern.Sg3Suffix); sg3 != "" {
			paradigm.Sg3 = sg3
		} else {
			continue
		}

		// Find pl1
		if pl1 := findCompatibleForm(bySlot["pl:pri"], sg1, pattern.Sg1Suffix, pattern.Pl1Suffix); pl1 != "" {
			paradigm.Pl1 = pl1
		} else {
			continue
		}

		// Find pl2
		if pl2 := findCompatibleForm(bySlot["pl:sec"], sg1, pattern.Sg1Suffix, pattern.Pl2Suffix); pl2 != "" {
			paradigm.Pl2 = pl2
		} else {
			continue
		}

		// Find pl3
		if pl3 := findCompatibleForm(bySlot["pl:ter"], sg1, pattern.Sg1Suffix, pattern.Pl3Suffix); pl3 != "" {
			paradigm.Pl3 = pl3
		} else {
			continue
		}

		// Skip archaic forms
		if isArchaicParadigm(paradigm) {
			continue
		}

		// Mark forms as used
		usedForms[sg1.Form] = true
		usedForms[paradigm.Sg2] = true
		usedForms[paradigm.Sg3] = true
		usedForms[paradigm.Pl1] = true
		usedForms[paradigm.Pl2] = true
		usedForms[paradigm.Pl3] = true

		paradigms = append(paradigms, paradigm)
	}

	return paradigms
}

// isArchaicParadigm returns true if the paradigm uses archaic conjugation patterns.
// These are forms that were standard in older Polish but have been replaced in modern usage.
func isArchaicParadigm(p VerbParadigm) bool {
	inf := p.Infinitive
	sg1 := p.Sg1

	// Pattern 1: -tać verbs with -tam instead of modern -czę
	// Archaic: szeptać → szeptam, mamrotać → mamrotam
	// Modern: szeptać → szepczę, mamrotać → mamroczę
	// Exception: regular -ać verbs like czytać → czytam are NOT archaic
	if strings.HasSuffix(inf, "tać") && !strings.HasSuffix(inf, "ytać") {
		// Check for -otać, -etać, -ptać patterns that should use -czę
		if strings.HasSuffix(inf, "otać") || strings.HasSuffix(inf, "etać") ||
			strings.HasSuffix(inf, "ptać") {
			if strings.HasSuffix(sg1, "tam") {
				return true // archaic -tam form
			}
		}
	}

	// Pattern 2: -ywać verbs with -wam instead of modern -uję
	// Archaic: wykonywać → wykonywam, pokazywać → pokazywam
	// Modern: wykonywać → wykonuję, pokazywać → pokazuję
	// Exception: bywać family (from być) correctly uses -wam
	if strings.HasSuffix(inf, "ywać") && strings.HasSuffix(sg1, "wam") {
		// Check if it's NOT a bywać derivative
		if !isBywacDerivative(inf) {
			return true // archaic -wam form
		}
	}

	// Pattern 3: -iwać verbs with -wam instead of modern -uję
	if strings.HasSuffix(inf, "iwać") && strings.HasSuffix(sg1, "wam") {
		return true // archaic -wam form
	}

	// Pattern 4: -awać verbs with -wam instead of modern -ję
	// Archaic: stawać → stawam, napawać → napawam
	// Modern: stawać → staję, napawać → napaję
	// Exception: -ywać handled above, -iwać handled above
	if strings.HasSuffix(inf, "awać") && strings.HasSuffix(sg1, "wam") {
		// Skip if already handled by -ywać or -iwać
		if !strings.HasSuffix(inf, "ywać") && !strings.HasSuffix(inf, "iwać") {
			return true // archaic -wam form
		}
	}

	// Pattern 5: -ować verbs with -wam instead of modern -uję
	// Archaic: kować → kowam, knować → knowam
	// Modern: kować → kuję, knować → knuję
	if strings.HasSuffix(inf, "ować") && strings.HasSuffix(sg1, "wam") {
		return true // archaic -wam form
	}

	// Pattern 6: -przeć verbs with -eję instead of standard -ę
	// Standard: oprzeć → oprę, przeć → prę
	// Variant: oprzeć → oprzeję (less common, treat as archaic for consistency)
	if strings.HasSuffix(inf, "przeć") && strings.HasSuffix(sg1, "eję") {
		return true
	}

	return false
}

// isBywacDerivative checks if a verb is derived from bywać (być + -wać).
// These correctly use -wam: bywać → bywam, przebywać → przebywam
func isBywacDerivative(inf string) bool {
	// Must end in -bywać
	if !strings.HasSuffix(inf, "bywać") {
		return false
	}
	// The part before -bywać should be empty or a valid prefix
	prefix := strings.TrimSuffix(inf, "bywać")
	if prefix == "" {
		return true // bywać itself
	}
	// Check for common verbal prefixes
	prefixes := []string{
		"do", "na", "o", "ob", "od", "po", "pod", "prze", "przy",
		"roz", "u", "w", "wy", "z", "za",
	}
	for _, p := range prefixes {
		if prefix == p {
			return true
		}
	}
	return false
}

// findPattern returns the conjugation pattern that matches the given sg1 form.
func findPattern(sg1 string) *conjugationPattern {
	for i := range knownPatterns {
		if strings.HasSuffix(sg1, knownPatterns[i].Sg1Suffix) {
			return &knownPatterns[i]
		}
	}
	return nil
}

// findCompatibleForm finds a form whose ending is consistent with the pattern.
// Given sg1 ending and expected ending, it looks for a form with the same stem.
func findCompatibleForm(candidates []VerbForm, sg1 VerbForm, sg1Suffix, expectedSuffix string) string {
	// Calculate the stem from sg1
	stem := strings.TrimSuffix(sg1.Form, sg1Suffix)
	expectedForm := stem + expectedSuffix

	// Look for exact match first
	for _, c := range candidates {
		if c.Form == expectedForm {
			// Prefer forms with matching aspect/reflexivity
			if c.Aspect == sg1.Aspect || c.Aspect == "" || sg1.Aspect == "" {
				return c.Form
			}
		}
	}

	// If no exact match, look for any form with matching aspect
	for _, c := range candidates {
		if c.Form == expectedForm {
			return c.Form
		}
	}

	return ""
}
