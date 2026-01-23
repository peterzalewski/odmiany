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
	// -eptać verbs: szeptać → szepczę
	heuristicEptac,
	// -łamać verbs: łamać → łamię
	heuristicLamac,
	// -ać verbs with consonant alternations: pisać → piszę
	heuristicAcAlternating,
	// -nąć verbs: ciągnąć → ciągnę
	heuristicNac,
	// -ąść verbs: trząść → trzęsę, siąść → siądę
	heuristicAsc,
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

	// Check if this verb conjugates as -wam/-wasz (not -uję)
	// Pattern-based: certain stem endings indicate -wam conjugation
	if usesYwacWamPattern(infinitive, stem) {
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

// usesYwacWamPattern determines if a -ywać verb conjugates as -wam/-wasz
// instead of the standard -uję/-ujesz pattern.
// The key insight: verbs derived from monosyllabic roots (być→bywać, myć→mywać,
// żyć→żywać, pływać, szyć→szywać) use -wam, while verbs with -ywać as a
// derivational suffix (pokazywać from pokazać) use -uję.
func usesYwacWamPattern(infinitive, stem string) bool {
	// -bywać: check if it's a prefixed form of bywać (from być)
	// e.g., odbywać, przebywać, zdobywać → -wam
	// but udziobywać (from dziobać), odrąbywać (from rąbać) → -uję
	if strings.HasSuffix(infinitive, "bywać") {
		if isPrefixedBywac(infinitive) {
			return true
		}
		return false
	}

	// -mywać: check if it's from myć (umywać, wymywać, obmywać)
	// but zatrzymywać, ułamywać (from trzymać, łamać) → -uję
	if strings.HasSuffix(infinitive, "mywać") {
		if isPrefixedMywac(infinitive) {
			return true
		}
		return false
	}

	// -rywać: check if it's from rwać/grać/kryć (zrywać, grywać, krywać)
	// but patrywać, orywać (from patrzeć, orać) → -uję
	if strings.HasSuffix(infinitive, "rywać") {
		if isPrefixedRywac(infinitive) {
			return true
		}
		return false
	}

	// -ływać: from pływać → always -wam
	if strings.HasSuffix(infinitive, "ływać") {
		return true
	}

	// -żywać: from żyć → always -wam (używać, zażywać, nadużywać)
	if strings.HasSuffix(infinitive, "żywać") {
		return true
	}

	// -czywać: from -czyć roots → always -wam (odpoczywać)
	if strings.HasSuffix(infinitive, "czywać") {
		return true
	}

	// -szywać: from szyć → always -wam (doszywać, przeszywać)
	if strings.HasSuffix(infinitive, "szywać") {
		return true
	}

	// -zywać: wzywać, odzywać from zew/zyw root → -wam
	// but związywać, okazywać, etc. → -uję
	if strings.HasSuffix(infinitive, "zywać") {
		// Only simple prefixes + zywać go to -wam
		if isPrefixedZywac(infinitive) {
			return true
		}
		return false
	}

	return false
}

// Common verbal prefixes in Polish
var verbalPrefixes = []string{
	"prze", "przy", "po", "pod", "podo", "od", "ode", "do", "za", "na", "nad", "nade",
	"u", "w", "we", "wy", "z", "ze", "s", "roz", "roze", "o", "ob", "obe",
}

// isPrefixedBywac checks if the verb is (prefixes) + bywać from być
func isPrefixedBywac(infinitive string) bool {
	base := strings.TrimSuffix(infinitive, "bywać")
	if base == "" {
		return true // bywać itself
	}
	// Strip prefixes repeatedly
	return canStripAllPrefixes(base)
}

// isPrefixedMywac checks if the verb is (prefixes) + mywać from myć
func isPrefixedMywac(infinitive string) bool {
	// mywać derivatives: [prefix]mywać (umywać, wymywać, obmywać, podmywać)
	// NOT zatrzymywać, wstrzymywać (from trzymać), ułamywać (from łamać)
	base := strings.TrimSuffix(infinitive, "mywać")
	if base == "" {
		return true // mywać itself
	}
	// If there's content before mywać, it should only be prefixes
	// trzymywać, łamywać patterns have content that's not just prefixes
	return canStripAllPrefixes(base)
}

// isPrefixedRywac checks if the verb is (prefixes) + rywać/grywać/krywać/srywać
func isPrefixedRywac(infinitive string) bool {
	// From rwać: zrywać, odrywać, porywać, urywać, wyrywać, etc.
	// From grać: grywać, zagrywać, rozgrywać, etc.
	// From kryć: krywać, ukrywać, odkrywać, etc.
	// From srać: srywać, zasrywać, etc. (vulgar)
	// NOT from orać: orywać, zaorywać → -uję
	// NOT from patrzeć: patrywać, przypatrywać → -uję

	// Check for grywać, krywać, srywać patterns
	if strings.HasSuffix(infinitive, "grywać") ||
		strings.HasSuffix(infinitive, "krywać") ||
		strings.HasSuffix(infinitive, "srywać") ||
		strings.HasSuffix(infinitive, "drywać") {
		return true
	}

	// For plain -rywać, check if it's prefixed rwać
	base := strings.TrimSuffix(infinitive, "rywać")
	if base == "" {
		return true // rywać itself
	}
	return canStripAllPrefixes(base)
}

// isPrefixedZywac checks if the verb is (prefixes) + zywać from zew/zyw
func isPrefixedZywac(infinitive string) bool {
	// wzywać, odzywać, przyzywać → -wam
	// związywać, pokazywać, etc. → -uję
	base := strings.TrimSuffix(infinitive, "zywać")
	if base == "" {
		return true
	}
	return canStripAllPrefixes(base)
}

// canStripAllPrefixes returns true if the string consists only of valid prefixes
func canStripAllPrefixes(s string) bool {
	if s == "" {
		return true
	}
	// Try each prefix
	for _, p := range verbalPrefixes {
		if strings.HasPrefix(s, p) {
			rest := strings.TrimPrefix(s, p)
			if canStripAllPrefixes(rest) {
				return true
			}
		}
	}
	return false
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
	// -otać verbs: t→cz in 1sg/3pl, t→c elsewhere
	// chichotać → chichoczę, chichocesz, chichoce, chichocemy, chichocecie, chichoczą
	return PresentTense{
		Sg1: stem + "czę",
		Sg2: stem + "cesz",
		Sg3: stem + "ce",
		Pl1: stem + "cemy",
		Pl2: stem + "cecie",
		Pl3: stem + "czą",
	}, true
}

// heuristicEptac handles -eptać verbs (and similar -ptać patterns).
// szeptać → szepczę, szepcesz, szepce (pt→pcz in 1sg/3pl, pt→pc elsewhere)
func heuristicEptac(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "ptać") {
		return PresentTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "tać")
	// pt→pcz in 1sg/3pl, pt→pc in others
	return PresentTense{
		Sg1: stem + "czę",
		Sg2: stem + "cesz",
		Sg3: stem + "ce",
		Pl1: stem + "cemy",
		Pl2: stem + "cecie",
		Pl3: stem + "czą",
	}, true
}

// heuristicLamac handles -łamać and -kłamać verbs.
// łamać → łamię, łamiesz, łamie (m→mi alternation)
func heuristicLamac(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "łamać") && !strings.HasSuffix(infinitive, "kłamać") {
		return PresentTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ać")
	return PresentTense{
		Sg1: stem + "ię",
		Sg2: stem + "iesz",
		Sg3: stem + "ie",
		Pl1: stem + "iemy",
		Pl2: stem + "iecie",
		Pl3: stem + "ią",
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
// przysnąć → przysnę, przyśniesz, przyśnie (s→ś before front vowels)
func heuristicNac(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "nąć") {
		return PresentTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ąć") // keeps the 'n'

	// For 1sg and 3pl, use hard stem (ends in n)
	sg1Stem := stem
	pl3Stem := stem

	// For other forms, soften consonant before 'n' if applicable
	// sn → śn, zn → źn before front vowels
	softStem := softenBeforeN(stem)

	return PresentTense{
		Sg1: sg1Stem + "ę",
		Sg2: softStem + "iesz",
		Sg3: softStem + "ie",
		Pl1: softStem + "iemy",
		Pl2: softStem + "iecie",
		Pl3: pl3Stem + "ą",
	}, true
}

// softenBeforeN softens consonants before 'n' at end of stem
// Rules:
// - sn → śn always
// - zn → źn only when preceded by front vowel (i, ę)
//   e.g., grzęznąć → grzęźnie, obliznąć → obliźnie
//   but pełznąć → pełznie (ł is not a front vowel)
func softenBeforeN(stem string) string {
	if strings.HasSuffix(stem, "sn") {
		return strings.TrimSuffix(stem, "sn") + "śn"
	}
	// For zn, check if preceded by front vowel
	if strings.HasSuffix(stem, "zn") {
		runes := []rune(stem)
		if len(runes) >= 3 {
			vowelBefore := runes[len(runes)-3]
			// Front vowels that trigger softening: i, ę
			if vowelBefore == 'i' || vowelBefore == 'ę' {
				return strings.TrimSuffix(stem, "zn") + "źn"
			}
		}
	}
	return stem
}

// heuristicAsc handles -ąść verbs.
// Three main patterns:
// - siąść type: usiąść → usiądę, usiądziesz, usiądzie (ą stays, ść→dzie)
// - trząść type: potrząść → potrzęsę, potrzęsiesz, potrzęsie (ą→ę, ść→s)
// - prząść type: uprząść → uprzędę, uprzędziesz, uprzędzie (ą→ę, ść→dzie)
func heuristicAsc(infinitive string) (PresentTense, bool) {
	if !strings.HasSuffix(infinitive, "ąść") {
		return PresentTense{}, false
	}

	// -siąść type (siedzieć family): ą stays, ść→dź
	if strings.HasSuffix(infinitive, "siąść") {
		stem := strings.TrimSuffix(infinitive, "ąść")
		return PresentTense{
			Sg1: stem + "ądę",
			Sg2: stem + "ądziesz",
			Sg3: stem + "ądzie",
			Pl1: stem + "ądziemy",
			Pl2: stem + "ądziecie",
			Pl3: stem + "ądą",
		}, true
	}

	// -trząść type (shake): ą→ę, ść→s
	if strings.HasSuffix(infinitive, "trząść") {
		stem := strings.TrimSuffix(infinitive, "ąść")
		return PresentTense{
			Sg1: stem + "ęsę",
			Sg2: stem + "ęsiesz",
			Sg3: stem + "ęsie",
			Pl1: stem + "ęsiemy",
			Pl2: stem + "ęsiecie",
			Pl3: stem + "ęsą",
		}, true
	}

	// -prząść type (spin): ą→ę, ść→dzie
	if strings.HasSuffix(infinitive, "prząść") {
		stem := strings.TrimSuffix(infinitive, "ąść")
		return PresentTense{
			Sg1: stem + "ędę",
			Sg2: stem + "ędziesz",
			Sg3: stem + "ędzie",
			Pl1: stem + "ędziemy",
			Pl2: stem + "ędziecie",
			Pl3: stem + "ędą",
		}, true
	}

	// Default: don't match, let other heuristics try
	return PresentTense{}, false
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

	// Determine 1sg and 3pl forms based on stem-final consonant
	// Both 1sg and 3pl undergo softening in Polish -ić verbs
	var sg1, pl3 string
	if softStem, ok := applySoftening(stem); ok {
		// Stem ends in consonant that softens: nosić → noszę, noszą
		// gościć → goszczę, goszczą
		sg1 = softStem + "ę"
		pl3 = softStem + "ą"
	} else if endsInSoftConsonant(stem) || endsInNonSoftenableC(stem) {
		// Stem already ends in soft consonant: chodzić → chodzę, chodzą
		// Or ends in c (non-softenable): cucić → cucę, kształcić → kształcę
		sg1 = stem + "ę"
		pl3 = stem + "ą"
	} else {
		// No softening needed: robić → robię, robią
		sg1 = stem + "ię"
		pl3 = stem + "ią"
	}

	return PresentTense{
		Sg1: sg1,
		Sg2: stem + "isz",
		Sg3: stem + "i",
		Pl1: stem + "imy",
		Pl2: stem + "icie",
		Pl3: pl3,
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

	// Verbs ending in soft consonant + -eć have two patterns:
	// 1. -ę/-ysz for action verbs: leżeć → leżę, krzyczeć → krzyczę
	// 2. -eję/-ejesz for inchoative verbs: maleć → maleję, boleć → boleję
	//
	// Pattern by ending (based on corpus statistics):
	// -żeć, -czeć, -rzeć: mostly -ę/-ysz (action verbs)
	// -leć, -szeć: mostly -eję/-ejesz (inchoative verbs)
	stem := strings.TrimSuffix(infinitive, "eć")

	// -leć and -szeć verbs are mostly inchoative → -eję pattern
	if strings.HasSuffix(infinitive, "leć") || strings.HasSuffix(infinitive, "szeć") {
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

	// -żeć, -czeć, -rzeć verbs are mostly action verbs → -ę/-ysz pattern
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

	// Other -eć verbs: use -em/-esz pattern (rare)
	stem = strings.TrimSuffix(infinitive, "ć")
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

// endsInNonSoftenableC returns true if stem ends in c that doesn't undergo softening.
// In Polish -ić verbs, stems ending in c take -ę (not -ię), unless the c is part
// of a softenable cluster like śc → szcz.
// Examples: cucić → cucę, kształcić → kształcę, but gościć → goszczę (via śc→szcz)
func endsInNonSoftenableC(stem string) bool {
	if !strings.HasSuffix(stem, "c") {
		return false
	}
	// Check if this c is part of a softenable pattern
	// śc → szcz is handled by applySoftening
	if strings.HasSuffix(stem, "śc") || strings.HasSuffix(stem, "źc") {
		return false // these go through applySoftening
	}
	return true
}

// hardeningMap maps hard consonants to their soft alternates.
// Used for consonant alternations before front vowels (ę, e, i).
var softeningMap = map[string]string{
	"śc": "szcz", // gościć → goszczę, czyścić → czyszczę (stem is gośc-, not gość-)
	"źc": "żdż",  // rare - if it exists
	"st": "szcz", // prosty → proszę (when applicable)
	"s":  "sz",   // nosić → noszę
	"z":  "ż",    // wozić → wożę
	"ź":  "ż",    // woźić → wożę (rare but exists)
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
	// If stem already ends in a soft consonant, no softening needed
	if endsInSoftConsonant(stem) {
		return "", false
	}

	// Check if the stem ends in a soft consonant + n cluster (like czn, żn, szn)
	// These should not be softened - the n is part of a soft cluster
	for _, soft := range []string{"cz", "sz", "ż", "rz", "dz"} {
		if strings.HasSuffix(stem, soft+"n") {
			return "", false
		}
	}

	// Try longer patterns first
	patterns := []string{"śc", "źc", "st", "sł", "zł", "sn", "zn", "ch"}
	for _, p := range patterns {
		if strings.HasSuffix(stem, p) {
			if soft, ok := softeningMap[p]; ok {
				return strings.TrimSuffix(stem, p) + soft, true
			}
		}
	}

	// Try single consonants
	singles := []string{"s", "z", "ź", "d", "t", "k", "g", "r"}
	for _, p := range singles {
		if strings.HasSuffix(stem, p) {
			if soft, ok := softeningMap[p]; ok {
				return strings.TrimSuffix(stem, p) + soft, true
			}
		}
	}

	return "", false
}
