package verb

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// VerbalNoun derives the verbal noun (rzeczownik odsłownikowy) from a Polish
// verb infinitive. Returns a slice because some verbs have multiple valid forms.
// Examples: czytać → ["czytanie"], pić → ["picie"], ciec → ["cieczenie", "cieknięcie"]
func VerbalNoun(infinitive string) ([]string, error) {
	// 1. Check irregular lookup (with prefix support)
	if forms, ok := lookupIrregularVerbalNoun(infinitive); ok {
		return forms, nil
	}

	// 2. -ać → -anie
	if strings.HasSuffix(infinitive, "ać") && !strings.HasSuffix(infinitive, "nąć") {
		stem := strings.TrimSuffix(infinitive, "ać")
		return []string{stem + "anie"}, nil
	}

	// 3. -nąć → soften + nięcie
	if strings.HasSuffix(infinitive, "nąć") {
		return verbalNounNac(infinitive), nil
	}

	// 4. Non-nąć -ąć → -ęcie
	if strings.HasSuffix(infinitive, "ąć") {
		stem := strings.TrimSuffix(infinitive, "ąć")
		return []string{stem + "ęcie"}, nil
	}

	// 5. -ić → softened stem + enie
	if strings.HasSuffix(infinitive, "ić") {
		return verbalNounIc(infinitive), nil
	}

	// 6. -uć → -ucie
	if strings.HasSuffix(infinitive, "uć") {
		stem := strings.TrimSuffix(infinitive, "uć")
		return []string{stem + "ucie"}, nil
	}

	// 7. -yć → -enie or -ycie
	if strings.HasSuffix(infinitive, "yć") {
		return verbalNounYc(infinitive), nil
	}

	// 8. -eć → -enie (with special cases)
	if strings.HasSuffix(infinitive, "eć") {
		return verbalNounEc(infinitive), nil
	}

	// 9. -c / -ść / -źć → should have been caught by irregular lookup
	return nil, fmt.Errorf("cannot derive verbal noun for %q", infinitive)
}

// verbalNounNac handles -nąć verbs: strip -nąć, soften before ń, add -nięcie.
func verbalNounNac(infinitive string) []string {
	stem := strings.TrimSuffix(infinitive, "nąć")
	softStem := softenBeforeNForGerund(stem)
	return []string{softStem + "nięcie"}
}

// softenBeforeNForGerund softens the final consonant of a stem before ń
// in verbal noun derivation.
//   - s → ś unless preceded by p, k, or m (ps, ks, ms clusters don't soften)
//   - z → ź unless z is part of rz or cz digraph, or łz cluster
func softenBeforeNForGerund(stem string) string {
	if strings.HasSuffix(stem, "s") {
		if len(stem) >= 2 {
			before := stem[len(stem)-2]
			if before == 'p' || before == 'k' || before == 'm' {
				return stem
			}
		}
		return stem[:len(stem)-1] + "ś"
	}

	if strings.HasSuffix(stem, "z") {
		if strings.HasSuffix(stem, "rz") || strings.HasSuffix(stem, "cz") ||
			strings.HasSuffix(stem, "łz") {
			return stem
		}
		return stem[:len(stem)-1] + "ź"
	}

	return stem
}

// verbalNounIc handles -ić verbs.
func verbalNounIc(infinitive string) []string {
	stem := strings.TrimSuffix(infinitive, "ić")

	// Vowel-ending stems: j-insertion → stem + jenie
	if endsInVowel(stem) {
		return []string{stem + "jenie"}
	}

	// Short stems (monosyllabic with a vowel): stem + icie
	// Consonant-only clusters like ćm, kp, tl are NOT monosyllabic
	runeCount := utf8.RuneCountInString(stem)
	if runeCount <= 2 && containsVowel(stem) {
		return []string{stem + "icie"}
	}

	// źdź softening: jeździć → jeżdżenie
	if strings.HasSuffix(stem, "źdz") {
		softened := strings.TrimSuffix(stem, "źdz") + "żdż"
		return []string{softened + "enie"}
	}

	// Try standard softening (but not for s in ks/ps clusters)
	if softStem, ok := applySofteningForGerund(stem); ok {
		return []string{softStem + "enie"}
	}

	// Soft consonant or non-softenable c: stem + enie
	if endsInSoftConsonant(stem) || endsInNonSoftenableC(stem) {
		return []string{stem + "enie"}
	}

	// Hard consonant without softening: keep i → stem + ienie
	return []string{stem + "ienie"}
}

// verbalNounYc handles -yć verbs.
func verbalNounYc(infinitive string) []string {
	stem := strings.TrimSuffix(infinitive, "yć")

	// Monosyllabic stems with a vowel: żyć → życie, myć → mycie
	// Consonant-only clusters like lż, mż are NOT monosyllabic
	runeCount := utf8.RuneCountInString(stem)
	if runeCount <= 2 && containsVowel(stem) {
		return []string{stem + "ycie"}
	}

	// Standard: uczyć → uczenie, burzyć → burzenie
	return []string{stem + "enie"}
}

// verbalNounEc handles -eć verbs.
func verbalNounEc(infinitive string) []string {
	// -Cieć pattern: consonant + ieć
	// Strip -ieć, check soft/hard, add -enie or -ienie.
	// Note: softening (s→sz etc.) is NOT productive for -eC-ieć gerunds —
	// the few exceptions (musieć, wisieć, chrzęścieć) are handled as irregulars.
	if strings.HasSuffix(infinitive, "ieć") && len(infinitive) > 3 {
		stem := strings.TrimSuffix(infinitive, "ieć")

		// Soft consonant or non-softenable c: stem + enie
		if endsInSoftConsonant(stem) || endsInNonSoftenableC(stem) {
			return []string{stem + "enie"}
		}

		// Hard consonant: keep i → stem + ienie
		return []string{stem + "ienie"}
	}

	// Plain -eć: strip -eć, add -enie
	stem := strings.TrimSuffix(infinitive, "eć")
	return []string{stem + "enie"}
}

// applySofteningForGerund applies consonant softening for -ić verbal nouns.
// Unlike present tense softening, s in consonant clusters (ks, ps) doesn't soften.
func applySofteningForGerund(stem string) (string, bool) {
	if endsInSoftConsonant(stem) {
		return "", false
	}

	// Don't soften s in ks or ps consonant clusters
	if strings.HasSuffix(stem, "s") {
		runes := []rune(stem)
		if len(runes) >= 2 {
			before := runes[len(runes)-2]
			if before == 'k' || before == 'p' {
				return "", false
			}
		}
	}

	return applySoftening(stem)
}

func isPolishVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u', 'y', 'ą', 'ę', 'ó':
		return true
	}
	return false
}

// containsVowel returns true if the string contains at least one Polish vowel.
func containsVowel(s string) bool {
	for _, r := range s {
		if isPolishVowel(r) {
			return true
		}
	}
	return false
}

// irregularVerbalNouns maps infinitives to their verbal noun form(s).
var irregularVerbalNouns = map[string][]string{
	// -rzeć → -arcie family (dual forms: warcie from rzeć-stem, wrzenie from plain -eć)
	"drzeć": {"darcie"},
	"mrzeć": {"marcie"},
	"przeć": {"parcie"},
	"trzeć": {"tarcie"},
	"wrzeć": {"warcie", "wrzenie"},
	"żreć":  {"żarcie"},

	// mleć/pleć with ie-insertion
	"mleć": {"mielenie"},
	"pleć": {"pielenie"},

	// -tworzyć → -twarcie (o→a vowel alternation)
	"otworzyć":    {"otwarcie"},
	"przetworzyć": {"przetwarcie"},
	"roztworzyć":  {"roztwarcie"},

	// słonić → słonięcie (ń + ięcie instead of nienie)
	"słonić": {"słonięcie"},

	// przychrzanić — Polimorf data artifact
	"przychrzanić": {"przychrzanienie"},

	// susnąć — irregular (doesn't soften s→ś)
	"susnąć": {"susnięcie"},

	// Monosyllabic -ić base verbs
	"bić": {"bicie"}, "gnić": {"gnicie"}, "pić": {"picie"}, "wić": {"wicie"},

	// Monosyllabic -yć base verbs
	"być": {"bycie"}, "żyć": {"życie"}, "myć": {"mycie"}, "ryć": {"rycie"},
	"szyć": {"szycie"}, "kryć": {"krycie"}, "wyć": {"wycie"}, "tyć": {"tycie"},

	// powić — powicie (not powienie; po+wić, keeps monosyllabic -icie ending)
	"powić": {"powicie"},

	// gzić — irregular z→ż (not caught by cluster rule)
	"gzić": {"gżenie"},

	// śnić — irregular (śnienie, not śnicie despite short stem)
	"śnić": {"śnienie"},

	// czcić/chrzcić — c→cz softening in rzc cluster
	"czcić":   {"czczenie"},
	"chrzcić": {"chrzczenie"},

	// -ić verbs where r+z is NOT the digraph rz (like marznąć)
	// The z softens independently to ż
	"mierzić": {"mierżenie"},

	// -ęzić verbs — z does NOT soften to ż (lexical exceptions)
	"gałęzić": {"gałęzienie"},
	"więzić":  {"więzienie"},

	// -uzić verbs — z does NOT soften to ż (lexical exceptions)
	"francuzić": {"francuzienie"},
	"kniazić":   {"kniazienie"},

	// -lesić — s does NOT soften to sz (lexical exception)
	"lesić": {"lesienie"},

	// -tłamsić — s DOES soften to sz (ms cluster is productive unlike ks/ps)
	"tłamsić": {"tłamszenie"},

	// -eć softening exceptions (non-inchoative verbs where s→sz IS correct)
	"musieć": {"muszenie"},
	"wisieć": {"wiszenie"},
	// chrzęścieć — the only -ścieć verb where śc→szcz
	"chrzęścieć": {"chrzęszczenie"},

	// -c verbs (present-tense stem based)
	"biec":   {"biegnięcie"},
	"ciec":   {"cieczenie", "cieknięcie"},
	"lec":    {"legnięcie", "lężenie"},
	"ląc":    {"lęgnięcie", "lęknięcie", "lężenie"},
	"móc":    {"możenie"},
	"piec":   {"pieczenie"},
	"rzec":   {"rzeczenie"},
	"siec":   {"sieczenie"},
	"strzec": {"strzeżenie"},
	"strzyc": {"strzyżenie"},
	"tłuc":   {"tłuczenie"},
	"wlec":   {"wleczenie"},
	"prząc":  {"przęgnięcie", "przężenie"},
	"siąc":   {"sięgnięcie", "siężenie"},

	// -ść verbs
	"bość":   {"bodzenie"},
	"bóść":   {"bodzenie"},
	"gnieść": {"gniecenie"},
	"grześć": {"grzebienie"},
	"iść":    {"iście"},
	"jeść":   {"jedzenie"},
	"kraść":  {"kradzenie"},
	"kłaść":  {"kładzenie"},
	"mieść":  {"miecenie"},
	"nieść":  {"niesienie"},
	"paść":   {"padnięcie", "pasienie"},
	"pleść":  {"plecenie"},
	"prząść": {"przędzenie"},
	"róść":   {"rośnięcie"},
	"siąść":  {"siądnięcie"},
	"trząść": {"trzęsienie"},

	// -jść (prefixed iść): the verbal noun stem is "jście"
	"jść": {"jście"},
	// -niść: wniść/wyniść/wzniść/zniść
	"niść": {"niście"},

	// pójść — special prefix (ó)
	"pójść": {"pójście"},

	// -źć verbs
	"gryźć":  {"gryzienie"},
	"grząźć": {"grzęzienie", "grzęźnięcie"},
	"leźć":   {"lezienie"},
	"liźć":   {"lezienie"},
	"wieźć":  {"wiezienie"},

	// Additional base verbs for -c/-ść/-źć
	"wieść":  {"wiedzenie"},
	"żec":    {"żegnięcie", "żżenie"},
	"wściec": {"wścieknięcie", "wścieczenie"},
	"oblec":  {"obleczenie"},
	"sieść":  {"siędnięcie"},

	// Compound-prefix verbs (double/triple prefix base forms)
	"pomóc":    {"pomożenie"},
	"domóc":    {"domożenie"},
	"naleźć":   {"nalezienie"},
	"najść":    {"najście"},
	"upaść":    {"upadnięcie"},
	"podnieść": {"podniesienie"},
	"przysiąc": {"przysięgnięcie", "przysiężenie"},
	"niemóc":   {"niemożenie"},
	"postrzec": {"postrzeżenie"},
	"wsiąść":   {"wsiądnięcie"},
	"strząść":  {"strzęsienie"},

	// Compound prefix bases with their own verbal noun forms
	"zbyć":  {"zbycie"},
	"dobyć": {"dobycie"},
	"użyć":  {"użycie"},

	// Prefixed żyć compounds
	"pożyć":  {"pożycie"},
	"spożyć": {"spożycie"},

	// współprzeżyć — too many prefixes for stripping
	"współprzeżyć": {"współprzeżycie"},

	// poszyć/sposzyć — szyć compounds
	"poszyć":  {"poszycie"},
	"sposzyć": {"sposzycie"},

	// Verbs with non-standard prefixes
	"ściec":            {"ścieczenie", "ścieknięcie"},
	"spostrzec":        {"spostrzeżenie"},
	"złorzec":          {"złorzeczenie", "złorzeknięcie"},
	"zapobiec":         {"zapobiegnięcie"},
	"współubiec":       {"współubiegnięcie"},
	"współposiąść":     {"współposiądnięcie"},
	"wspomóc":          {"wspomożenie"},
	"krzywoprzysiąc":   {"krzywoprzysięgnięcie", "krzywoprzysiężenie"},
	"zaprzepaść":       {"zaprzepadnięcie"},
	"nadojeść":         {"nadojedzenie"},
	"półwisieć":        {"półwiszenie"},
	"przesiąc":         {"przesiąknięcie"},
	"niedomóc":         {"niedomożenie"},
	"współżyć":         {"współżycie"},
	"zbezeczcić":       {"zbezeczczenie"},
	"zeźreć":           {"zziarcie"},
	"osić":             {"oszenie"},
	"zażyznić":         {"zażyznienie"},

	// Voicing assimilation with z- prefix (z+t→st, z+p→sp in spelling)
	"zetrzeć": {"starcie"},
	"zeprzeć": {"sparcie"},

	// sprzeć/wesprzeć — distinct verbal noun stems
	"sprzeć":   {"sprzenie"},
	"wesprzeć": {"wsparcie"},

	// ode- prefix kept for wrzeć
	"odewrzeć": {"odewarcie"},

	// rozpostrzeć — compound trzeć
	"rozpostrzeć": {"rozpostarcie"},

	// Prefixed mrzeć with -u- infix (obumrzeć, odumrzeć, zaumrzeć)
	"obumrzeć": {"obumarcie"},
	"odumrzeć": {"odumarcie"},
	"zaumrzeć": {"zaumarcie"},

	// zeźrzeć/zeżreć — suppletive
	"zeźrzeć": {"zziarcie"},
	"zeżreć":  {"zżarcie"},
	"zrzeć":   {"żarcie"},
}

// verbalNounPrefixableVerbs lists base verbs whose verbal noun can be derived
// by prefixing.
var verbalNounPrefixableVerbs = map[string]string{
	// Monosyllabic -ić verbs
	"bić": "bić", "gnić": "gnić", "pić": "pić", "wić": "wić",
	// powić: po+wić, keeps monosyllabic ending
	"powić": "powić",
	// Monosyllabic -yć verbs
	"być": "być", "żyć": "żyć", "myć": "myć", "ryć": "ryć",
	"szyć": "szyć", "kryć": "kryć", "wyć": "wyć", "tyć": "tyć",
	// -c verbs
	"biec": "biec", "ciec": "ciec", "lec": "lec", "móc": "móc",
	"piec": "piec", "rzec": "rzec", "siec": "siec", "strzec": "strzec",
	"strzyc": "strzyc", "tłuc": "tłuc", "wlec": "wlec",
	"prząc": "prząc", "siąc": "siąc", "ląc": "ląc",
	// -ść verbs
	"bość": "bość", "bóść": "bóść", "gnieść": "gnieść", "grześć": "grześć",
	"iść": "iść", "jeść": "jeść", "kraść": "kraść", "kłaść": "kłaść",
	"mieść": "mieść", "nieść": "nieść", "paść": "paść", "pleść": "pleść",
	"prząść": "prząść", "róść": "róść", "siąść": "siąść", "trząść": "trząść",
	// -jść (prefixed iść, e.g. dojść, przejść)
	"jść": "jść",
	// -nijść (e.g. wnijść, wynijść)
	"nijść": "nijść",
	// -niść (e.g. wniść, wyniść)
	"niść": "niść",
	// -źć verbs
	"gryźć": "gryźć", "grząźć": "grząźć", "leźć": "leźć", "liźć": "liźć",
	"wieźć": "wieźć",
	// -eć → -arcie family
	"drzeć": "drzeć", "mrzeć": "mrzeć", "przeć": "przeć",
	"trzeć": "trzeć", "wrzeć": "wrzeć", "żreć": "żreć",
	// mleć/pleć
	"mleć": "mleć", "pleć": "pleć",
	// słonić
	"słonić": "słonić",
	// -eć softening exceptions
	"musieć": "musieć", "wisieć": "wisieć",
	"chrzęścieć": "chrzęścieć",
	// gzić
	"gzić": "gzić",
	// -ić softening exceptions (prefixable)
	"mierzić":   "mierzić",
	"gałęzić":   "gałęzić",
	"więzić":     "więzić",
	"francuzić":  "francuzić",
	"lesić":      "lesić",
	"tłamsić":    "tłamsić",
	// Compound prefix bases with their own verbal noun forms
	"zbyć":  "zbyć",
	"dobyć": "dobyć",
	"użyć":  "użyć",
	"pożyć": "pożyć",
	// poszyć/sposzyć — szyć compounds
	"poszyć": "poszyć",
	// Compound bases for -c/-ść/-źć
	"wieść": "wieść", "żec": "żec", "wściec": "wściec",
	"oblec": "oblec", "sieść": "sieść",
	"pomóc": "pomóc", "domóc": "domóc",
	"naleźć": "naleźć", "najść": "najść",
	"upaść": "upaść", "podnieść": "podnieść",
	"przysiąc": "przysiąc", "niemóc": "niemóc",
	"postrzec": "postrzec", "wsiąść": "wsiąść",
	"strząść": "strząść",
	// śnić: irregular — śnienie
	"śnić": "śnić",
	// czcić/chrzcić — c→cz in rzc cluster
	"czcić": "czcić", "chrzcić": "chrzcić",
}

// lookupIrregularVerbalNoun checks the irregular map, including prefix stripping.
// Handles epenthetic vowels in prefixes (ode+przeć → odeprzeć → odparcie).
func lookupIrregularVerbalNoun(infinitive string) ([]string, bool) {
	// Direct lookup
	if forms, ok := irregularVerbalNouns[infinitive]; ok {
		return forms, true
	}

	// Try stripping prefixes
	for _, prefix := range verbPrefixes {
		if len(infinitive) > len(prefix) && infinitive[:len(prefix)] == prefix {
			base := infinitive[len(prefix):]
			if baseKey, ok := verbalNounPrefixableVerbs[base]; ok {
				if baseForms, ok := irregularVerbalNouns[baseKey]; ok {
					p := stripEpentheticVowelForGerund(prefix, baseForms[0])
					forms := make([]string, len(baseForms))
					for i, f := range baseForms {
						forms[i] = p + f
					}
					return forms, true
				}
			}
		}
	}

	return nil, false
}

// stripEpentheticVowelForGerund strips the trailing 'e' from prefixes like
// "ode", "pode", etc. for verbal noun derivation.
// For single-consonant short prefixes (z, w), the epenthetic vowel is kept
// before sibilants (s, ś, ź, z, ż, sz) and w to avoid unpronounceable clusters.
func stripEpentheticVowelForGerund(prefix, baseForm string) string {
	epenthetic := map[string]string{
		"ode": "od", "pode": "pod", "nade": "nad", "roze": "roz",
		"wze": "wz", "obe": "ob", "we": "w", "ze": "z",
	}

	short, ok := epenthetic[prefix]
	if !ok {
		return prefix
	}

	if len(baseForm) > 0 {
		firstRune, _ := utf8.DecodeRuneInString(baseForm)
		// Keep epenthetic vowel if base starts with a vowel
		if isPolishVowel(firstRune) {
			return prefix
		}
		// Keep epenthetic vowel for -jść forms (nadejść → nadejście)
		if strings.HasPrefix(baseForm, "jście") || strings.HasPrefix(baseForm, "jść") {
			return prefix
		}
		// For single-consonant short prefixes (z, w), keep epenthetic vowel
		// to avoid unpronounceable clusters at the prefix boundary.
		if len(short) == 1 {
			// Keep before s-family sibilants (zs, zś, zsz, wś are bad clusters)
			if firstRune == 's' || firstRune == 'ś' || firstRune == 'ź' ||
				firstRune == 'z' {
				return prefix
			}
			if strings.HasPrefix(baseForm, "sz") {
				return prefix
			}
			// Keep we- before w (avoid ww doubling)
			if short == "w" && firstRune == 'w' {
				return prefix
			}
			// Keep ze- before w+consonant (avoid clusters like zwl-)
			if short == "z" && firstRune == 'w' && len(baseForm) >= 2 {
				secondRune, _ := utf8.DecodeRuneInString(baseForm[utf8.RuneLen(firstRune):])
				if !isPolishVowel(secondRune) {
					return prefix
				}
			}
		}
	}

	return short
}

// nijść is a special compound base: ni + jść. We handle it specially
// to avoid needing it in the irregular map as a separate entry.
func init() {
	// Register "nijść" as a prefixable base
	irregularVerbalNouns["nijść"] = []string{"nijście"}
	verbalNounPrefixableVerbs["nijść"] = "nijść"
}
