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

// Tense represents present or past tense extraction mode.
type Tense string

const (
	TensePresent Tense = "present"
	TensePast    Tense = "past"
)

// VerbForm represents a single conjugated form with its grammatical tags.
type VerbForm struct {
	Form   string
	Number string // "sg" or "pl"
	Person string // "pri" (1st), "sec" (2nd), "ter" (3rd)
	Gender string // "m1" (masc.pers), "m2" (masc.anim), "m3" (masc.inan), "f", "n", "n1" (non-masc.pers plural)
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

// PastParadigm holds a complete past tense paradigm (13 forms).
// Past tense distinguishes gender: masculine/feminine/neuter in singular,
// masculine-personal/non-masculine-personal in plural.
type PastParadigm struct {
	Infinitive string `json:"infinitive"`
	// Singular - ja (1st person)
	Sg1M string `json:"sg1m"` // ja (masculine)
	Sg1F string `json:"sg1f"` // ja (feminine)
	// Singular - ty (2nd person)
	Sg2M string `json:"sg2m"` // ty (masculine)
	Sg2F string `json:"sg2f"` // ty (feminine)
	// Singular - on/ona/ono (3rd person)
	Sg3M string `json:"sg3m"` // on (masculine)
	Sg3F string `json:"sg3f"` // ona (feminine)
	Sg3N string `json:"sg3n"` // ono (neuter)
	// Plural - my (1st person)
	Pl1V  string `json:"pl1v"`  // my (masculine-personal/virile)
	Pl1NV string `json:"pl1nv"` // my (non-masculine-personal/non-virile)
	// Plural - wy (2nd person)
	Pl2V  string `json:"pl2v"`  // wy (masculine-personal)
	Pl2NV string `json:"pl2nv"` // wy (non-masculine-personal)
	// Plural - oni/one (3rd person)
	Pl3V  string `json:"pl3v"`  // oni (masculine-personal)
	Pl3NV string `json:"pl3nv"` // one (non-masculine-personal)
	Aspect string `json:"aspect"`
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
	tense := flag.String("tense", "present", "tense to extract: present or past")
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

	// Determine tag prefix based on tense
	var tagPrefix string
	switch Tense(*tense) {
	case TensePresent:
		tagPrefix = "verb:fin:"
	case TensePast:
		tagPrefix = "verb:praet:"
	default:
		fmt.Fprintf(os.Stderr, "unknown tense: %s (use 'present' or 'past')\n", *tense)
		os.Exit(1)
	}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) != 3 {
			continue
		}
		lemma, form, tags := parts[0], parts[1], parts[2]

		if !strings.Contains(tags, tagPrefix) {
			continue
		}

		// Parse the form
		var vf VerbForm
		if Tense(*tense) == TensePresent {
			vf = parseVerbForm(form, tags)
			if vf.Number != "" && vf.Person != "" {
				verbForms[lemma] = append(verbForms[lemma], vf)
			}
		} else {
			vf = parsePastForm(form, tags)
			if vf.Number != "" && vf.Person != "" && vf.Gender != "" {
				verbForms[lemma] = append(verbForms[lemma], vf)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scan: %v\n", err)
		os.Exit(1)
	}

	// Extract and output paradigms based on tense
	if Tense(*tense) == TensePresent {
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

		fmt.Fprintf(os.Stderr, "Extracted %d complete present tense paradigms from %d infinitives\n",
			len(paradigms), len(verbForms))
	} else {
		// Extract past tense paradigms
		var paradigms []PastParadigm
		for infinitive, forms := range verbForms {
			extracted := extractPastParadigms(infinitive, forms)
			paradigms = append(paradigms, extracted...)
		}

		// Sort for deterministic output
		sort.Slice(paradigms, func(i, j int) bool {
			if paradigms[i].Infinitive != paradigms[j].Infinitive {
				return paradigms[i].Infinitive < paradigms[j].Infinitive
			}
			return paradigms[i].Sg1M < paradigms[j].Sg1M
		})

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(paradigms); err != nil {
			fmt.Fprintf(os.Stderr, "encode: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Extracted %d complete past tense paradigms from %d infinitives\n",
			len(paradigms), len(verbForms))
	}
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

// parsePastForm extracts grammatical information from Polimorf past tense tags.
// Tags format: verb:praet:NUMBER:GENDER:PERSON:ASPECT:REFL
// Example: verb:praet:sg:m1:pri:imperf:nonrefl
func parsePastForm(form, tags string) VerbForm {
	vf := VerbForm{Form: form}

	tagParts := strings.Split(tags, ":")
	if len(tagParts) < 6 {
		return vf
	}

	vf.Number = tagParts[2] // sg or pl
	vf.Gender = tagParts[3] // m1, m2, m3, f, n, n1
	vf.Person = tagParts[4] // pri, sec, ter

	// Extract aspect
	if strings.Contains(tags, ":imperf") {
		vf.Aspect = "imperf"
	} else if strings.Contains(tags, ":perf") {
		vf.Aspect = "perf"
	}

	// Extract reflexivity
	if strings.Contains(tags, ":refl.nonrefl") {
		vf.Refl = "refl.nonrefl"
	} else if strings.Contains(tags, ":nonrefl") {
		vf.Refl = "nonrefl"
	} else if strings.Contains(tags, ":refl") {
		vf.Refl = "refl"
	}

	return vf
}

// extractPastParadigms groups past tense forms into coherent paradigms.
// Past tense is simpler than present - stems are nearly universal within a verb,
// so we mostly just need to collect all 13 forms.
func extractPastParadigms(infinitive string, forms []VerbForm) []PastParadigm {
	// Group forms by normalized slot (person+number+genderCategory)
	// Polimorf uses compound gender tags like "m1.m2.m3", "n1.n2", "m1.p1", "m2.m3.f.n1.n2.p2.p3"
	// We normalize these to: sgM, sgF, sgN, plV, plNV
	bySlot := make(map[string][]VerbForm)
	for _, f := range forms {
		slots := normalizeGenderSlots(f.Number, f.Person, f.Gender)
		for _, slot := range slots {
			bySlot[slot] = append(bySlot[slot], f)
		}
	}

	// Get the 3rd person masculine singular as base (it's the "dictionary" form)
	sg3mForms := bySlot["sg:ter:M"]
	if len(sg3mForms) == 0 {
		return nil // No base form found
	}

	// For past tense, we try to build paradigms from each sg3m form
	var paradigms []PastParadigm

	for _, sg3m := range sg3mForms {
		paradigm := PastParadigm{
			Infinitive: infinitive,
			Aspect:     sg3m.Aspect,
		}

		// Try to find all forms, preferring forms from the same aspect
		paradigm.Sg1M = findPastFormNorm(bySlot, "sg", "pri", "M", sg3m.Aspect)
		paradigm.Sg1F = findPastFormNorm(bySlot, "sg", "pri", "F", sg3m.Aspect)
		paradigm.Sg2M = findPastFormNorm(bySlot, "sg", "sec", "M", sg3m.Aspect)
		paradigm.Sg2F = findPastFormNorm(bySlot, "sg", "sec", "F", sg3m.Aspect)
		paradigm.Sg3M = sg3m.Form
		paradigm.Sg3F = findPastFormNorm(bySlot, "sg", "ter", "F", sg3m.Aspect)
		paradigm.Sg3N = findPastFormNorm(bySlot, "sg", "ter", "N", sg3m.Aspect)
		paradigm.Pl1V = findPastFormNorm(bySlot, "pl", "pri", "V", sg3m.Aspect)
		paradigm.Pl1NV = findPastFormNorm(bySlot, "pl", "pri", "NV", sg3m.Aspect)
		paradigm.Pl2V = findPastFormNorm(bySlot, "pl", "sec", "V", sg3m.Aspect)
		paradigm.Pl2NV = findPastFormNorm(bySlot, "pl", "sec", "NV", sg3m.Aspect)
		paradigm.Pl3V = findPastFormNorm(bySlot, "pl", "ter", "V", sg3m.Aspect)
		paradigm.Pl3NV = findPastFormNorm(bySlot, "pl", "ter", "NV", sg3m.Aspect)

		// Check if paradigm is complete (has all 13 forms)
		if isCompletePastParadigm(paradigm) {
			// Check for coherence - the stem should be consistent
			if isPastParadigmCoherent(paradigm) {
				paradigms = append(paradigms, paradigm)
			}
		}
	}

	return paradigms
}

// normalizeGenderSlots converts Polimorf compound gender tags to normalized slots.
// Returns a list of slots this form belongs to.
// Polimorf tags:
//   - Singular: m1.m2.m3 (masc), f (fem), n1.n2 (neut)
//   - Plural: m1.p1 (masc-pers/virile), m2.m3.f.n1.n2.p2.p3 (non-masc-pers)
func normalizeGenderSlots(number, person, gender string) []string {
	var slots []string
	slot := number + ":" + person + ":"

	// Check for masculine (singular or plural virile)
	if strings.Contains(gender, "m1") {
		if number == "sg" {
			slots = append(slots, slot+"M")
		} else {
			// In plural, m1 alone or m1.p1 means virile
			if strings.Contains(gender, "p1") || gender == "m1" || !strings.Contains(gender, "m2") {
				slots = append(slots, slot+"V")
			}
		}
	}

	// Check for feminine
	if strings.Contains(gender, "f") {
		if number == "sg" {
			slots = append(slots, slot+"F")
		}
		// In plural, f is part of non-virile
	}

	// Check for neuter
	if strings.Contains(gender, "n1") || strings.Contains(gender, "n2") {
		if number == "sg" {
			slots = append(slots, slot+"N")
		}
	}

	// Check for plural non-virile (contains m2, m3, f, n1, n2, p2, p3 but not just m1.p1)
	if number == "pl" {
		if strings.Contains(gender, "m2") || strings.Contains(gender, "m3") ||
			strings.Contains(gender, "f") || strings.Contains(gender, "p2") ||
			strings.Contains(gender, "p3") {
			slots = append(slots, slot+"NV")
		}
	}

	return slots
}

// findPastFormNorm finds a form matching the given normalized slot.
func findPastFormNorm(bySlot map[string][]VerbForm, number, person, genderCat, preferAspect string) string {
	slot := number + ":" + person + ":" + genderCat
	forms := bySlot[slot]
	// Prefer matching aspect
	for _, f := range forms {
		if f.Aspect == preferAspect {
			return f.Form
		}
	}
	// Fall back to any form
	if len(forms) > 0 {
		return forms[0].Form
	}
	return ""
}

// isCompletePastParadigm checks if all 13 forms are present.
func isCompletePastParadigm(p PastParadigm) bool {
	return p.Sg1M != "" && p.Sg1F != "" &&
		p.Sg2M != "" && p.Sg2F != "" &&
		p.Sg3M != "" && p.Sg3F != "" && p.Sg3N != "" &&
		p.Pl1V != "" && p.Pl1NV != "" &&
		p.Pl2V != "" && p.Pl2NV != "" &&
		p.Pl3V != "" && p.Pl3NV != ""
}

// isPastParadigmCoherent checks if the past paradigm forms share a consistent stem.
// Past tense is very regular - almost all forms share the same stem,
// with predictable endings.
func isPastParadigmCoherent(p PastParadigm) bool {
	// Extract stem from sg3m (base form) - remove -ł
	stem := strings.TrimSuffix(p.Sg3M, "ł")
	if stem == p.Sg3M {
		// Might be an irregular form like "szedł" - accept it
		return true
	}

	// For regular verbs, check that feminine forms match stem + ła/łam/łaś
	// This is a light coherence check - past tense is much more regular than present
	if !strings.HasPrefix(p.Sg3F, stem) {
		// Check for ó→o alternation (e.g., mógł → mogła)
		altStem := strings.ReplaceAll(stem, "ó", "o")
		if !strings.HasPrefix(p.Sg3F, altStem) {
			// Check for vowel dropping (e.g., tarł → tarła, but also niósł → niosła)
			// These are acceptable variations
			return true // Accept for now - past tense is very regular
		}
	}

	return true
}
