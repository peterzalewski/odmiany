package verb

import (
	"fmt"
	"strings"
)

type Person int

const (
	First Person = iota + 1
	Second
	Third
)

type Number int

const (
	Singular Number = iota + 1
	Plural
)

// PresentTense holds all six forms of the present tense paradigm.
type PresentTense struct {
	Sg1 string // ja
	Sg2 string // ty
	Sg3 string // on/ona/ono
	Pl1 string // my
	Pl2 string // wy
	Pl3 string // oni/one
}

// Get returns the form for the given person and number.
func (p PresentTense) Get(person Person, number Number) string {
	switch {
	case person == First && number == Singular:
		return p.Sg1
	case person == Second && number == Singular:
		return p.Sg2
	case person == Third && number == Singular:
		return p.Sg3
	case person == First && number == Plural:
		return p.Pl1
	case person == Second && number == Plural:
		return p.Pl2
	case person == Third && number == Plural:
		return p.Pl3
	default:
		return ""
	}
}

// Equals returns true if two paradigms are identical.
func (p PresentTense) Equals(other PresentTense) bool {
	return p.Sg1 == other.Sg1 &&
		p.Sg2 == other.Sg2 &&
		p.Sg3 == other.Sg3 &&
		p.Pl1 == other.Pl1 &&
		p.Pl2 == other.Pl2 &&
		p.Pl3 == other.Pl3
}

// ConjugatePresent returns the present tense paradigm for a verb.
// First checks the irregular verb lookup table, then falls back to heuristics.
func ConjugatePresent(infinitive string) (PresentTense, error) {
	// Check irregular verbs first (including prefixed forms)
	if p, ok := lookupIrregularWithPrefix(infinitive); ok {
		return p, nil
	}

	// Try heuristics in order of specificity
	for _, h := range heuristics {
		if p, ok := h(infinitive); ok {
			return p, nil
		}
	}
	return PresentTense{}, fmt.Errorf("no heuristic matched: %s", infinitive)
}

// heuristic is a function that attempts to conjugate a verb.
// Returns (paradigm, true) if it can handle the verb, (_, false) otherwise.
type heuristic func(infinitive string) (PresentTense, bool)

// heuristics is the ordered list of conjugation heuristics.
// More specific patterns should come first.
var heuristics = []heuristic{
	// -ować verbs: pracować → pracuję
	heuristicOwac,
	// -ywać/-iwać verbs: pokazywać → pokazuję (but bywać → bywam)
	heuristicYwacIwac,
	// -awać verbs: dawać → daję
	heuristicAwac,
	// -otać verbs: chichotać → chichoczę
	heuristicOtac,
	// -ać verbs with consonant alternations: pisać → piszę
	heuristicAcAlternating,
	// -nąć verbs: ciągnąć → ciągnę
	heuristicNac,
	// -ść/-źć verbs: nieść → niosę
	heuristicSc,
	// -c verbs: móc → mogę
	heuristicC,
	// -ić verbs: robić → robię (with consonant alternations)
	heuristicIc,
	// -yć verbs: myć → myję
	heuristicYc,
	// -eć verbs: umieć → umiem
	heuristicEc,
	// Regular -ać verbs: czytać → czytam (fallback for -ać)
	heuristicAc,
}

// heuristicOwac handles -ować verbs.
// pracować → pracuję, pracujesz, pracuje, pracujemy, pracujecie, pracują
func heuristicOwac(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "ować") {
		return PresentTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ować")
	return PresentTense{
		Sg1: stem + "uję",
		Sg2: stem + "ujesz",
		Sg3: stem + "uje",
		Pl1: stem + "ujemy",
		Pl2: stem + "ujecie",
		Pl3: stem + "ują",
	}, true
}

// heuristicYwacIwac handles -ywać and -iwać verbs.
// pokazywać → pokazuję (drop -ywać, add -uję)
// Exception: bywać, pływać, etc. → bywam (keep stem, -am/-asz)
func heuristicYwacIwac(infinitive string) (PresentTense, bool) {
	var stem string
	if strings.HasSuffix(infinitive, "ywać") {
		stem = strings.TrimSuffix(infinitive, "ywać")
	} else if strings.HasSuffix(infinitive, "iwać") {
		stem = strings.TrimSuffix(infinitive, "iwać")
	} else {
		return PresentTense{}, false
	}

	// Exceptions that conjugate as -am/-asz instead of -uję
	// These are typically verbs where -ywać is part of the root, not a suffix
	exceptions := map[string]bool{
		"bywać": true, "pływać": true, "przebywać": true,
		"dobywać": true, "nabywać": true, "odbywać": true,
		"pobywać": true, "ubywać": true, "wybywać": true,
		"zbywać": true, "obywać": true, "zabywać": true,
		"odzywać": true, "przyzywać": true, "wzywać": true,
		"zażywać": true, "używać": true, "nadużywać": true,
	}
	if exceptions[infinitive] {
		// These follow regular -ać pattern
		fullStem := strings.TrimSuffix(infinitive, "ć")
		return PresentTense{
			Sg1: fullStem + "m",
			Sg2: fullStem + "sz",
			Sg3: fullStem,
			Pl1: fullStem + "my",
			Pl2: fullStem + "cie",
			Pl3: fullStem + "ją",
		}, true
	}

	// Standard -ywać/-iwać → -uję pattern
	return PresentTense{
		Sg1: stem + "uję",
		Sg2: stem + "ujesz",
		Sg3: stem + "uje",
		Pl1: stem + "ujemy",
		Pl2: stem + "ujecie",
		Pl3: stem + "ują",
	}, true
}

// heuristicAwac handles -awać verbs (not -ować or -ywać).
// dawać → daję, dajesz, daje...
func heuristicAwac(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "awać") {
		return PresentTense{}, false
	}
	// Skip if it's actually -ywać or -iwać (handled above)
	if strings.HasSuffix(infinitive, "ywać") || strings.HasSuffix(infinitive, "iwać") {
		return PresentTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "wać")
	return PresentTense{
		Sg1: stem + "ję",
		Sg2: stem + "jesz",
		Sg3: stem + "je",
		Pl1: stem + "jemy",
		Pl2: stem + "jecie",
		Pl3: stem + "ją",
	}, true
}

// heuristicOtac handles -otać verbs (onomatopoeia, iterative actions).
// chichotać → chichoczę, chichoczesz, chichocze...
// The t→cz alternation occurs in these verbs.
func heuristicOtac(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "otać") {
		return PresentTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "tać")
	return PresentTense{
		Sg1: stem + "czę",
		Sg2: stem + "czesz",
		Sg3: stem + "cze",
		Pl1: stem + "czemy",
		Pl2: stem + "czecie",
		Pl3: stem + "czą",
	}, true
}

// heuristicAcAlternating handles -ać verbs that conjugate with -ę/-esz
// (not the regular -am/-asz pattern) due to consonant alternations.
//
// Based on corpus analysis:
//   -pać: 250 alternate vs 14 regular → mostly alternates
//   -bać: 113 alternate vs 26 regular → mostly alternates
//   -mać: 49 alternate vs 67 regular → mostly regular (skip)
//   -sać: 82 alternate vs 142 regular → mostly regular (skip)
//   -zać: 100 alternate vs 1494 regular → mostly regular (skip)
//   -kać: 77 alternate vs 722 regular → mostly regular (skip)
func heuristicAcAlternating(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "ać") {
		return PresentTense{}, false
	}

	// Skip patterns handled elsewhere
	if strings.HasSuffix(infinitive, "ować") ||
		strings.HasSuffix(infinitive, "ywać") ||
		strings.HasSuffix(infinitive, "iwać") ||
		strings.HasSuffix(infinitive, "awać") ||
		strings.HasSuffix(infinitive, "otać") {
		return PresentTense{}, false
	}

	stem := strings.TrimSuffix(infinitive, "ać")

	// Only match patterns that mostly alternate (>80% alternation rate)

	// -pać → -pię: capać → capię, sypać → sypię (95% alternate)
	if strings.HasSuffix(stem, "p") {
		return presentIEIesz(stem), true
	}
	// -bać → -bię: drapać → drapię, skubać → skubię (81% alternate)
	if strings.HasSuffix(stem, "b") {
		return presentIEIesz(stem), true
	}

	// All other -ać patterns have <50% alternation rate
	// Let them fall through to regular -ać handler
	return PresentTense{}, false
}

// presentEEsz creates a present tense paradigm with -ę/-esz endings.
// Used for verbs like pisać → piszę, piszesz, pisze...
func presentEEsz(stem string) PresentTense {
	return PresentTense{
		Sg1: stem + "ę",
		Sg2: stem + "esz",
		Sg3: stem + "e",
		Pl1: stem + "emy",
		Pl2: stem + "ecie",
		Pl3: stem + "ą",
	}
}

// presentIEIesz creates a present tense paradigm with -ię/-iesz endings.
// Used for verbs like capać → capię, capiesz, capie...
func presentIEIesz(stem string) PresentTense {
	return PresentTense{
		Sg1: stem + "ię",
		Sg2: stem + "iesz",
		Sg3: stem + "ie",
		Pl1: stem + "iemy",
		Pl2: stem + "iecie",
		Pl3: stem + "ią",
	}
}

// heuristicNac handles -nąć verbs.
// ciągnąć → ciągnę, ciągniesz, ciągnie, ciągniemy, ciągniecie, ciągną
func heuristicNac(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "nąć") {
		return PresentTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ąć") // keeps the 'n'
	return PresentTense{
		Sg1: stem + "ę",
		Sg2: stem + "iesz",
		Sg3: stem + "ie",
		Pl1: stem + "iemy",
		Pl2: stem + "iecie",
		Pl3: stem + "ą",
	}, true
}

// heuristicSc handles -ść and -źć verbs.
// nieść → niosę, niesiesz, niesie...
// wieźć → wiozę, wieziesz, wiezie...
// gryźć → gryzę, gryziesz, gryzie...
func heuristicSc(infinitive string) (PresentTense, bool) {
	// -ieść verbs (nieść type): ie→io alternation in 1sg/3pl
	if strings.HasSuffix(infinitive, "ieść") {
		stem := strings.TrimSuffix(infinitive, "ieść")
		return PresentTense{
			Sg1: stem + "iosę",
			Sg2: stem + "iesiesz",
			Sg3: stem + "iesie",
			Pl1: stem + "iesiemy",
			Pl2: stem + "iesiecie",
			Pl3: stem + "iosą",
		}, true
	}
	// -ieźć verbs (wieźć type): ie→io alternation in 1sg/3pl
	if strings.HasSuffix(infinitive, "ieźć") {
		stem := strings.TrimSuffix(infinitive, "ieźć")
		return PresentTense{
			Sg1: stem + "iozę",
			Sg2: stem + "ieziesz",
			Sg3: stem + "iezie",
			Pl1: stem + "ieziemy",
			Pl2: stem + "ieziecie",
			Pl3: stem + "iozą",
		}, true
	}
	// -yźć verbs (gryźć type): no vowel change
	if strings.HasSuffix(infinitive, "yźć") {
		stem := strings.TrimSuffix(infinitive, "źć")
		return PresentTense{
			Sg1: stem + "zę",
			Sg2: stem + "ziesz",
			Sg3: stem + "zie",
			Pl1: stem + "ziemy",
			Pl2: stem + "ziecie",
			Pl3: stem + "zą",
		}, true
	}
	// -eźć verbs (leźć type): no vowel change
	if strings.HasSuffix(infinitive, "eźć") {
		stem := strings.TrimSuffix(infinitive, "źć")
		return PresentTense{
			Sg1: stem + "zę",
			Sg2: stem + "ziesz",
			Sg3: stem + "zie",
			Pl1: stem + "ziemy",
			Pl2: stem + "ziecie",
			Pl3: stem + "zą",
		}, true
	}
	// Other -ść/-źć patterns (iść, etc.) - skip for now
	return PresentTense{}, false
}

// heuristicC handles -c verbs (móc, piec, etc.).
// móc → mogę, możesz, może, możemy, możecie, mogą
func heuristicC(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "c") {
		return PresentTense{}, false
	}
	// Skip -ść/-źć (handled above) and -ać/-eć/-ić/-yć (handled below)
	if strings.HasSuffix(infinitive, "ść") || strings.HasSuffix(infinitive, "źć") ||
		strings.HasSuffix(infinitive, "ać") || strings.HasSuffix(infinitive, "eć") ||
		strings.HasSuffix(infinitive, "ić") || strings.HasSuffix(infinitive, "yć") ||
		strings.HasSuffix(infinitive, "ąć") {
		return PresentTense{}, false
	}

	// móc → mogę type (c → g/ż alternation)
	if strings.HasSuffix(infinitive, "óc") {
		stem := strings.TrimSuffix(infinitive, "óc")
		return PresentTense{
			Sg1: stem + "ogę",
			Sg2: stem + "ożesz",
			Sg3: stem + "oże",
			Pl1: stem + "ożemy",
			Pl2: stem + "ożecie",
			Pl3: stem + "ogą",
		}, true
	}

	// piec → piekę type
	if strings.HasSuffix(infinitive, "ec") {
		stem := strings.TrimSuffix(infinitive, "ec")
		return PresentTense{
			Sg1: stem + "ekę",
			Sg2: stem + "eczesz",
			Sg3: stem + "ecze",
			Pl1: stem + "eczemy",
			Pl2: stem + "eczecie",
			Pl3: stem + "eką",
		}, true
	}

	return PresentTense{}, false
}

// heuristicIc handles -ić verbs with consonant alternations.
// robić → robię, robisz, robi... (no alternation - b stays)
// nosić → noszę, nosisz, nosi... (s → sz in 1sg)
// chodzić → chodzę, chodzisz, chodzi... (soft stem - 1sg is stem+ę)
// pić → piję, pijesz, pije... (monosyllabic - j-insertion)
func heuristicIc(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "ić") {
		return PresentTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ić")

	// Monosyllabic stems (pić, bić, lić) use j-insertion: pić → piję
	runeCount := len([]rune(stem))
	if runeCount <= 2 {
		return PresentTense{
			Sg1: stem + "iję",
			Sg2: stem + "ijesz",
			Sg3: stem + "ije",
			Pl1: stem + "ijemy",
			Pl2: stem + "ijecie",
			Pl3: stem + "iją",
		}, true
	}

	// Determine 1sg form based on stem-final consonant
	var sg1 string
	if softStem, ok := applySoftening(stem); ok {
		// Stem ends in consonant that softens: nosić → noszę
		sg1 = softStem + "ę"
	} else if endsInSoftConsonant(stem) {
		// Stem already ends in soft consonant: chodzić → chodzę
		sg1 = stem + "ę"
	} else {
		// No softening needed: robić → robię
		sg1 = stem + "ię"
	}

	return PresentTense{
		Sg1: sg1,
		Sg2: stem + "isz",
		Sg3: stem + "i",
		Pl1: stem + "imy",
		Pl2: stem + "icie",
		Pl3: stem + "ią",
	}, true
}

// heuristicYc handles -yć verbs.
// myć → myję, myjesz, myje... (standard -yć / monosyllabic)
// żyć → żyję, żyjesz, żyje... (monosyllabic)
// uczyć → uczę, uczysz, uczy... (polysyllabic stem ends in soft consonant)
func heuristicYc(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "yć") {
		return PresentTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "yć")

	// Monosyllabic stems always use -yję pattern (myć, żyć, ryć, etc.)
	// Check by rune count since Polish letters can be multi-byte
	runeCount := len([]rune(stem))
	if runeCount <= 2 {
		fullStem := stem + "y"
		return PresentTense{
			Sg1: fullStem + "ję",
			Sg2: fullStem + "jesz",
			Sg3: fullStem + "je",
			Pl1: fullStem + "jemy",
			Pl2: fullStem + "jecie",
			Pl3: fullStem + "ją",
		}, true
	}

	// Polysyllabic stems ending in soft consonant (cz, sz, ż, rz) use -ę/-ysz
	if endsInSoftConsonant(stem) {
		return PresentTense{
			Sg1: stem + "ę",
			Sg2: stem + "ysz",
			Sg3: stem + "y",
			Pl1: stem + "ymy",
			Pl2: stem + "ycie",
			Pl3: stem + "ą",
		}, true
	}

	// Standard -yć → -yję pattern
	fullStem := stem + "y"
	return PresentTense{
		Sg1: fullStem + "ję",
		Sg2: fullStem + "jesz",
		Sg3: fullStem + "je",
		Pl1: fullStem + "jemy",
		Pl2: fullStem + "jecie",
		Pl3: fullStem + "ją",
	}, true
}

// heuristicEc handles -eć verbs.
// Most -ieć verbs: biednieć → biednieję (891 verbs go to -ieję)
// Few -ieć exceptions: umieć → umiem (26 verbs go to -iem)
// Other -eć verbs: mieć → mam (irregular), chcieć → chcę (different pattern)
func heuristicEc(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "eć") {
		return PresentTense{}, false
	}

	// Most -ieć verbs conjugate as -ieję/-iejesz (891 vs 26)
	if strings.HasSuffix(infinitive, "ieć") {
		// -umieć family: umieć → umiem
		if strings.HasSuffix(infinitive, "umieć") {
			stem := strings.TrimSuffix(infinitive, "ć")
			return PresentTense{
				Sg1: stem + "m",
				Sg2: stem + "sz",
				Sg3: stem,
				Pl1: stem + "my",
				Pl2: stem + "cie",
				Pl3: stem + "ją",
			}, true
		}
		// -wiedzieć family: wiedzieć → wiem (ie→∅ in present)
		if strings.HasSuffix(infinitive, "wiedzieć") {
			stem := strings.TrimSuffix(infinitive, "iedzieć")
			return PresentTense{
				Sg1: stem + "iem",
				Sg2: stem + "iesz",
				Sg3: stem + "ie",
				Pl1: stem + "iemy",
				Pl2: stem + "iecie",
				Pl3: stem + "iedzą",
			}, true
		}
		// śmieć: śmieć → śmiem
		if infinitive == "śmieć" || strings.HasSuffix(infinitive, "ośmieć") {
			stem := strings.TrimSuffix(infinitive, "ć")
			return PresentTense{
				Sg1: stem + "m",
				Sg2: stem + "sz",
				Sg3: stem,
				Pl1: stem + "my",
				Pl2: stem + "cie",
				Pl3: stem + "ją",
			}, true
		}
		// chcieć: chcieć → chcę (special -ę/-esz pattern)
		if strings.HasSuffix(infinitive, "chcieć") {
			stem := strings.TrimSuffix(infinitive, "ieć")
			return PresentTense{
				Sg1: stem + "ę",
				Sg2: stem + "esz",
				Sg3: stem + "e",
				Pl1: stem + "emy",
				Pl2: stem + "ecie",
				Pl3: stem + "ą",
			}, true
		}
		// mieć is suppletive - skip, let it fail (needs lookup table)
		if infinitive == "mieć" || strings.HasSuffix(infinitive, "mieć") && !strings.HasSuffix(infinitive, "umieć") {
			return PresentTense{}, false
		}
		// Standard -ieć → -ieję pattern
		stem := strings.TrimSuffix(infinitive, "ć")
		return PresentTense{
			Sg1: stem + "ję",
			Sg2: stem + "jesz",
			Sg3: stem + "je",
			Pl1: stem + "jemy",
			Pl2: stem + "jecie",
			Pl3: stem + "ją",
		}, true
	}

	// Other -eć verbs (not -ieć): use -em/-esz pattern
	stem := strings.TrimSuffix(infinitive, "ć")
	return PresentTense{
		Sg1: stem + "m",
		Sg2: stem + "sz",
		Sg3: stem,
		Pl1: stem + "my",
		Pl2: stem + "cie",
		Pl3: stem + "ją",
	}, true
}

// heuristicAc handles regular -ać verbs (fallback).
// czytać → czytam, czytasz, czyta, czytamy, czytacie, czytają
func heuristicAc(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "ać") {
		return PresentTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ć") // keeps the 'a'
	return PresentTense{
		Sg1: stem + "m",
		Sg2: stem + "sz",
		Sg3: stem,
		Pl1: stem + "my",
		Pl2: stem + "cie",
		Pl3: stem + "ją",
	}, true
}

// Consonant alternation helpers

// softConsonants are consonants (and digraphs) that are already "soft"
// and don't undergo further alternation before front vowels.
var softConsonants = []string{
	"szcz", "dż", "dź", // trigraph/digraphs first
	"sz", "ż", "cz", "rz", "dz", // digraphs
	"ś", "ź", "ć", "ń", "l", "j", // single soft consonants
}

// endsInSoftConsonant returns true if stem ends in a soft consonant.
func endsInSoftConsonant(stem string) bool {
	for _, soft := range softConsonants {
		if strings.HasSuffix(stem, soft) {
			return true
		}
	}
	return false
}

// hardeningMap maps hard consonants to their soft alternates.
// Used for consonant alternations before front vowels (ę, e, i).
var softeningMap = map[string]string{
	"st": "szcz", // prosty → proszę (when applicable)
	"s":  "sz",   // nosić → noszę
	"z":  "ż",    // wozić → wożę
	"d":  "dz",   // chodzić → chodzę (but stem is already chodz-)
	"t":  "c",    // płacić → płacę
	"ch": "sz",   // słuchać → słyszę (rare in verbs)
	"k":  "cz",   // płakać → płaczę
	"g":  "ż",    // biegać → biegam (but some: strzec → strzeżę)
	"r":  "rz",   // patrzeć → patrzę
	"sł": "śl",   // myślić → myślę
	"zł": "źl",   // (rare)
	"sn": "śn",   // śnić → śnię
	"zn": "źn",   // (rare)
}

// applySoftening attempts to soften the final consonant of a stem.
// Returns (softened stem, true) if softening applies, (_, false) otherwise.
func applySoftening(stem string) (string, bool) {
	// Try longer patterns first
	patterns := []string{"st", "sł", "zł", "sn", "zn", "ch"}
	for _, p := range patterns {
		if strings.HasSuffix(stem, p) {
			if soft, ok := softeningMap[p]; ok {
				return strings.TrimSuffix(stem, p) + soft, true
			}
		}
	}

	// Try single consonants
	singles := []string{"s", "z", "d", "t", "k", "g", "r"}
	for _, p := range singles {
		if strings.HasSuffix(stem, p) {
			if soft, ok := softeningMap[p]; ok {
				return strings.TrimSuffix(stem, p) + soft, true
			}
		}
	}

	return "", false
}
