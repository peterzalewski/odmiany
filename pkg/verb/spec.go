package verb

import (
	"strings"
	"unicode/utf8"
)

// verbSpec unifies all irregular data for a single verb.
// nil fields mean "derive from heuristics".
type verbSpec struct {
	pres       *presSpec // nil = derive from heuristics
	past       *pastSpec // nil = derive from heuristics
	verbalNoun []string  // nil = derive from heuristics
}

// irregularSpecs is the unified map of all verbs with any irregular component.
// Populated at package init time from the spec builders below.
var irregularSpecs map[string]verbSpec

// prefixableVerbs lists verbs that can take prefixes productively.
// Used by lookupIrregularSpec for prefix-stripping lookup.
var prefixableVerbs = map[string]bool{
	// Present tense prefixable
	"pisać": true, "brać": true, "jechać": true, "dać": true,
	"wziąć": true, "iść": true, "jeść": true, "prać": true,
	"czesać": true, "kasać": true, "ciosać": true, "ciesać": true,
	"skakać": true, "płakać": true, "wiązać": true, "kazać": true,
	"mazać": true, "lizać": true, "kołysać": true, "krzesać": true,
	"naleźć": true, "spać": true, "bać": true, "dziać": true,
	"podobać": true,
	// Monosyllabic verbs
	"bić": true, "lić": true, "pić": true, "żyć": true, "myć": true,
	"ryć": true, "szyć": true, "wyć": true, "kryć": true,
	// Other prefixable present bases
	"pomnieć": true, "mrzeć": true, "ciec": true, "woleć": true,
	"jąć": true, "cząć": true, "patrzeć": true,
	"rwać": true, "zwać": true, "dbać": true, "śmiać": true,
	"cierpieć": true, "wisieć": true, "jeździć": true,
	"pachnieć": true, "strzec": true, "chować": true,
	"grzmieć": true, "szumieć": true, "tłumieć": true,
	"okazać": true, "karać": true, "kraść": true, "kłaść": true,
	"lać": true, "grześć": true, "przeć": true, "wrzeć": true,
	"śnić": true, "rzec": true, "wiać": true, "krajać": true,
	"słać": true, "nająć": true, "tłuc": true, "pleść": true, "kląć": true,
	"żreć": true, "chwiać": true,
	"starzeć": true, "gorzeć": true, "dorzeć": true, "dobrzeć": true,
	"czcić": true, "kpić": true, "ulec": true, "wściec": true,
	"dojrzeć": true, "boleć": true, "swędzieć": true,
	"tajać": true, "ćpać": true, "wić": true,
	"bimbać": true, "gabać": true, "chybać": true, "gnić": true,
	"siać": true, "gibać": true, "siorbać": true, "stąpać": true,
	"pchlać": true, "rychlać": true, "gdybać": true,
	"użyć": true,
	// Inchoative -eć verbs
	"chorzeć": true, "tężeć": true, "dumieć": true, "goreć": true,
	"śniedzieć": true, "srebrzeć": true, "cukrzeć": true,
	// Additional prefixable bases
	"łajać": true, "bajać": true, "pierdzieć": true, "skomleć": true,
	"strzeliwać": true, "myśliwać": true, "boliwać": true, "mgliwać": true,
	"kpać": true, "tlić": true, "clić": true, "dlić": true,
	"kasłać": true, "mieszywać": true, "supływać": true, "bazgrywać": true,
	"podobywać": true, "cierpać": true, "siąpać": true, "tyrpać": true,
	"ściubać": true, "ślipać": true, "bombać": true,
	"szedzieć": true, "piać": true, "spiać": true, "skuliwać": true,
	"kaszliwać": true, "pyskiwać": true, "ziajać": true,
	"śmierdzieć": true,

	// Past tense prefixable (additions not already present)
	"być": true, "ciąć": true,
	"dąć": true, "giąć": true, "piąć": true, "miąć": true,
	"żąć": true,
	"siąść": true, "paść": true, "prząść": true,
	"gryźć": true, "leźć": true, "wieźć": true, "nieść": true,
	"trzeć": true, "drzeć": true,
	"stać": true, "mieć": true,
	"wiedzieć": true, "siedzieć": true, "widzieć": true,
	"biec": true, "lec": true, "piec": true, "wlec": true,
	"rosnąć": true, "rość": true, "schnąć": true, "przysięgnąć": true,
	"umrzeć": true,
	"mleć": true, "pleć": true, "żec": true,
	"musieć": true, "słyszeć": true, "móc": true,

	// Verbal noun prefixable (additions not already present)
	"tyć": true,
	"powić": true,
	"chrzcić": true,
	"mierzić": true, "gałęzić": true, "więzić": true,
	"francuzić": true, "lesić": true, "tłamsić": true,
	"chrzęścieć": true,
	"gzić": true,
	"siec": true, "strzyc": true, "prząc": true, "siąc": true, "ląc": true,
	"bość": true, "bóść": true, "gnieść": true,
	"mieść": true, "róść": true, "trząść": true,
	"jść": true, "nijść": true, "niść": true,
	"grząźć": true, "liźć": true,
	"słonić": true,
	"wieść": true, "oblec": true, "sieść": true,
	"pomóc": true, "domóc": true,
	"upaść": true, "podnieść": true,
	"przysiąc": true, "niemóc": true,
	"postrzec": true, "wsiąść": true,
	"strząść": true,
	// Compound prefix bases for VN
	"zbyć": true, "dobyć": true, "pożyć": true,
	"poszyć": true, "najść": true,
}

func init() {
	irregularSpecs = buildIrregularSpecs()
}

// buildIrregularSpecs merges the three legacy maps into a unified verbSpec map.
func buildIrregularSpecs() map[string]verbSpec {
	specs := make(map[string]verbSpec, 600)

	// Helper to get-or-create a spec entry
	get := func(verb string) verbSpec {
		return specs[verb]
	}

	// 1. Populate from present tense specs
	for verb, ps := range irregularPresSpecs {
		s := get(verb)
		ps := ps // copy
		s.pres = &ps
		specs[verb] = s
	}

	// 2. Populate from past tense specs
	for verb, ps := range irregularPastSpecs {
		s := get(verb)
		ps := ps // copy
		s.past = &ps
		specs[verb] = s
	}

	// 3. Populate from verbal noun map
	for verb, forms := range irregularVerbalNouns {
		s := get(verb)
		formsCopy := make([]string, len(forms))
		copy(formsCopy, forms)
		s.verbalNoun = formsCopy
		specs[verb] = s
	}

	return specs
}

// lookupIrregularPres looks up a verb's present tense spec in the unified map,
// including prefix-stripping for known prefixable bases.
func lookupIrregularPres(infinitive string) (ps presSpec, prefix string, found bool) {
	// Direct lookup first
	if s, ok := irregularSpecs[infinitive]; ok && s.pres != nil {
		return *s.pres, "", true
	}

	// Try stripping prefixes to find base irregular verb
	for _, pfx := range verbPrefixes {
		if len(infinitive) > len(pfx) && infinitive[:len(pfx)] == pfx {
			base := infinitive[len(pfx):]
			if prefixableVerbs[base] {
				if s, ok := irregularSpecs[base]; ok && s.pres != nil {
					return *s.pres, pfx, true
				}
			}
		}
	}

	return presSpec{}, "", false
}

// lookupIrregularPast looks up a verb's past tense spec in the unified map,
// including prefix-stripping for known prefixable bases.
func lookupIrregularPast(infinitive string) (ps pastSpec, prefix string, found bool) {
	// Direct lookup first
	if s, ok := irregularSpecs[infinitive]; ok && s.past != nil {
		return *s.past, "", true
	}

	// Try stripping prefixes to find base irregular verb
	for _, pfx := range verbPrefixes {
		if len(infinitive) > len(pfx) && infinitive[:len(pfx)] == pfx {
			base := infinitive[len(pfx):]
			if prefixableVerbs[base] {
				if s, ok := irregularSpecs[base]; ok && s.past != nil {
					return *s.past, pfx, true
				}
			}
		}
	}

	return pastSpec{}, "", false
}

// lookupIrregularVN looks up a verb's verbal noun forms in the unified map,
// including prefix-stripping for known prefixable bases.
func lookupIrregularVN(infinitive string) (forms []string, prefix string, found bool) {
	// Direct lookup first
	if s, ok := irregularSpecs[infinitive]; ok && s.verbalNoun != nil {
		return s.verbalNoun, "", true
	}

	// Try stripping prefixes to find base irregular verb
	for _, pfx := range verbPrefixes {
		if len(infinitive) > len(pfx) && infinitive[:len(pfx)] == pfx {
			base := infinitive[len(pfx):]
			if prefixableVerbs[base] {
				if s, ok := irregularSpecs[base]; ok && s.verbalNoun != nil {
					return s.verbalNoun, pfx, true
				}
			}
		}
	}

	return nil, "", false
}

// applyPrefixToPresent applies a prefix to all forms of a present tense paradigm.
func applyPrefixToPresent(prefix string, pt PresentTense) PresentTense {
	return PresentTense{
		Sg1: prefix + pt.Sg1,
		Sg2: prefix + pt.Sg2,
		Sg3: prefix + pt.Sg3,
		Pl1: prefix + pt.Pl1,
		Pl2: prefix + pt.Pl2,
		Pl3: prefix + pt.Pl3,
	}
}

// applyPrefixToVerbalNoun applies a prefix to verbal noun forms,
// handling epenthetic vowel stripping.
func applyPrefixToVerbalNoun(prefix string, baseForms []string) []string {
	p := stripEpentheticVowelForVN(prefix, baseForms[0])
	forms := make([]string, len(baseForms))
	for i, f := range baseForms {
		forms[i] = p + f
	}
	return forms
}

// epentheticPrefixes maps epenthetic prefix forms to their short forms.
// Used by verbal noun prefix stripping.
var epentheticPrefixes = map[string]string{
	"ode": "od", "pode": "pod", "nade": "nad", "roze": "roz",
	"wze": "wz", "obe": "ob", "we": "w", "ze": "z",
}

// stripEpentheticVowelForVN strips the trailing 'e' from prefixes like
// "ode", "pode", etc. for verbal noun derivation.
// For single-consonant short prefixes (z, w), the epenthetic vowel is kept
// before s-family sibilants (s, ś, ź, z, sz) and before consonant clusters
// that would be unpronounceable (ww, zwl-, etc.).
func stripEpentheticVowelForVN(prefix, baseForm string) string {
	short, ok := epentheticPrefixes[prefix]
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
