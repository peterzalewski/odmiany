package verb

import (
	"cmp"
	"strings"
)

// pastSpec compactly represents a past tense paradigm via stems.
// All 13 forms are derivable from at most 4 stems.
type pastSpec struct {
	stem   string // base stem (default for all positions)
	masc   string // masculine sg override (sg1m, sg2m) - defaults to stem
	sg3m   string // complete sg3m form - defaults to resolved_masc+"ł"
	fem    string // feminine/neuter/non-virile override - defaults to stem
	virile string // virile plural override - defaults to resolved_fem
}

func (s pastSpec) build() PastTense {
	masc := cmp.Or(s.masc, s.stem)
	fem := cmp.Or(s.fem, s.stem)
	vir := cmp.Or(s.virile, fem)
	sg3m := cmp.Or(s.sg3m, masc+"ł")
	return PastTense{
		Sg1M: masc + "łem", Sg1F: fem + "łam",
		Sg2M: masc + "łeś", Sg2F: fem + "łaś",
		Sg3M: sg3m, Sg3F: fem + "ła", Sg3N: fem + "ło",
		Pl1V: vir + "liśmy", Pl1NV: fem + "łyśmy",
		Pl2V: vir + "liście", Pl2NV: fem + "łyście",
		Pl3V: vir + "li", Pl3NV: fem + "ły",
	}
}

// pastHomographs contains verbs with multiple valid past tense paradigms.
var pastHomographs = map[string][]PastParadigm{
	// wlec: "to drag" has two valid sg3m forms (wlekł/wlókł), but all other forms use wlek-
	"wlec": {
		{PastTense: pastSpec{stem: "wlek"}.build(), Gloss: "sg3m wlekł variant"},
		{PastTense: pastSpec{stem: "wlek", sg3m: "wlókł"}.build(), Gloss: "sg3m wlókł variant"},
	},

	// paść: "to fall" (padł) vs "to graze" (pasł)
	"paść": {
		{PastTense: pastSpec{stem: "pas", virile: "paś"}.build(), Gloss: "to graze (animals)"},
		{PastTense: pastSpec{stem: "pad"}.build(), Gloss: "to fall"},
	},
}

// buildPascHomograph creates homograph entries for prefixed -paść verbs.
// These verbs have two valid paradigms (fall vs graze pattern) plus a mixed form.
func buildPascHomograph(prefix string) []PastParadigm {
	return []PastParadigm{
		// Pure "fall" pattern
		{PastTense: pastSpec{stem: prefix + "pad"}.build(), Gloss: "to fall"},
		// Mixed pattern: masc sg "graze", fem/pl "fall"
		{PastTense: pastSpec{masc: prefix + "pas", fem: prefix + "pad"}.build(), Gloss: "to fall (variant)"},
	}
}

// buildWlecHomograph creates homograph entries for prefixed -wlec verbs.
// These verbs have two valid sg3m forms (wlekł/wlókł) but all other forms use wlek-.
func buildWlecHomograph(prefix string) []PastParadigm {
	return []PastParadigm{
		{PastTense: pastSpec{stem: prefix + "wlek"}.build(), Gloss: "sg3m wlekł variant"},
		{PastTense: pastSpec{stem: prefix + "wlek", sg3m: prefix + "wlókł"}.build(), Gloss: "sg3m wlókł variant"},
	}
}

func init() {
	// Add homographs for prefixed -paść verbs
	pascPrefixes := []string{"do", "na", "od", "o", "pod", "po", "prze", "przy", "roz", "s", "u", "w", "wy", "za", "zaprze"}
	for _, p := range pascPrefixes {
		pastHomographs[p+"paść"] = buildPascHomograph(p)
	}

	// Add homographs for prefixed -wlec verbs
	// wlec has two valid sg3m forms (wlekł/wlókł), all other forms use wlek-
	wlecPrefixes := []string{"do", "na", "ob", "od", "o", "pod", "po", "prze", "przy", "roz", "u", "we", "w", "wy", "za", "ze", "z"}
	for _, p := range wlecPrefixes {
		pastHomographs[p+"wlec"] = buildWlecHomograph(p)
	}
}

// lookupPastHomograph returns all paradigms for a past tense homograph verb.
func lookupPastHomograph(infinitive string) ([]PastParadigm, bool) {
	if paradigms, ok := pastHomographs[infinitive]; ok {
		return paradigms, true
	}
	return nil, false
}

// nDroppingNacVerbs contains -nąć verbs that drop the n COMPLETELY in past tense.
// These are typically inchoative (state-change) verbs.
// gasnąć → gasł (not gasnął), schnąć → schł (not schnął)
// IMPORTANT: This list is for verbs that drop n in ALL forms (base + prefixed).
// Verbs that only drop n in virile plural go in mixedNDropNacVerbs.
// NOTE: Many Polish verbs have BOTH n-dropping and n-keeping variants in usage.
// We follow the more common/standard pattern for each verb.
var nDroppingNacVerbs = map[string]bool{
	// Verified n-dropping verbs (state change / inchoative)
	// NOTE: Many verbs have BOTH n-dropping and n-keeping variants in the corpus.
	// This list contains verbs where n-dropping is more common/standard.
	"blednąć": true, "bladnąć": true, "blaknąć": true,
	"brzęknąć": true, "brzydnąć": true,
	"cienknąć": true, "chłodnąć": true, "chrzypnąć": true, "chrypnąć": true,
	"chudnąć": true, "cichnąć": true, "ciemnąć": true, "cieknąć": true, "cierpnąć": true,
	"czeznąć": true, "ćwirknąć": true,
	"duchnąć": true,
	"gadnąć": true, "gasnąć": true, "gnuśnąć": true, "głuchnąć": true, "gorknąć": true,
	"grzęznąć": true, "grząznąć": true, "grąznąć": true, "gręznąć": true,
	"jaśnąć": true,
	"kisnąć": true, "klęknąć": true, "klęsnąć": true, "kostnąć": true, "kraśnąć": true,
	"krzepnąć": true, "krzesnąć": true, "kwaśnąć": true, "kwitnąć": true, "kładnąć": true,
	"lepnąć": true, "lęgnąć": true, "lęknąć": true,
	"marznąć": true, "mdlnąć": true, "mierznąć": true, "mierżnąć": true, "mięknąć": true,
	"milknąć": true, "moknąć": true,
	"pełznąć": true, "pęknąć": true, "pierzchnąć": true, "puchnąć": true,
	"przycichnąć": true, "przęgnąć": true,
	"rymsnąć": true, "rypnąć": true, "rzadnąć": true, "rzednąć": true,
	"sieknąć": true, "skrzepnąć": true, "słabnąć": true,
	"stęchnąć": true, "stęgnąć": true, "strzęgnąć": true, "stygnąć": true, "świrknąć": true,
	"ścichnąć": true, "ścierpnąć": true, "ślepnąć": true,
	"śmiardnąć": true, "śmierdnąć": true, "świerknąć": true,
	"tęchnąć": true, "twardnąć": true,
	"usechnąć": true, "usychnąć": true,
	"wiąznąć": true, "więzgnąć": true, "więznąć": true, "więdnąć": true, "wilgnąć": true, "wyknąć": true,
	"zbadnąć": true, "zdechnąć": true, "ziębnąć": true,
	"zmierzchnąć": true, "zwiędnąć": true,
	"żółknąć": true,
	// Note: smoknąć is NOT prefixed s+moknąć - it's a separate verb that keeps n
}

// mixedNDropNacVerbs contains -nąć verbs with MIXED n-dropping pattern.
// Singular/non-virile plural retain n, virile plural drops n.
// cuchnąć → cuchnął/cuchnęła/cuchli (sg retain n, pl virile drop)
// NOTE: If a verb appears in both this list AND nDroppingNacVerbs, this takes
// precedence, so make sure verbs that should fully drop n are NOT in this list.
var mixedNDropNacVerbs = map[string]bool{
	"buchnąć":    true,
	"cuchnąć":    true,
	"gęstnąć":    true,
	"mierzchnąć": true,
	// NOTE: niknąć and pachnąć removed - corpus shows they keep n in virile
	// (uniknęli, pachnęli, not unikli, pachli)
}

// irregularPastSpecs contains compact stem specs for verbs that cannot
// be conjugated by heuristics alone. Expanded to irregularPastVerbs in init().
var irregularPastSpecs = map[string]pastSpec{
	// === -jść verbs (special stem: pójść → poszedł, not pószedł) ===
	"pójść":  {masc: "poszed", fem: "posz"},
	"obejść": {masc: "obszed", fem: "obesz"},
	"odejść": {masc: "odszed", fem: "odesz"},
	"podejść": {masc: "podszed", fem: "podesz"},
	"nadejść": {masc: "nadszed", fem: "nadesz"},
	"rozejść": {masc: "rozszed", fem: "rozesz"},

	// === -naleźć verbs (special stem: odnaleźć → odnalazł) ===
	"znaleźć":  {stem: "znalaz", virile: "znaleź"},
	"odnaleźć": {stem: "odnalaz", virile: "odnaleź"},
	"wynaleźć": {stem: "wynalaz", virile: "wynaleź"},

	// === -sieść verbs (sieść → siadł, special alternation) ===
	"sieść":   {stem: "siad", virile: "sied"},
	"osiąść":  {stem: "osiad", virile: "osied"},
	"posieść": {stem: "posiad", virile: "posied"},
	"osieść":  {stem: "osiad", virile: "osied"},

	// Note: -paść verbs (dopaść, napaść, etc.) are handled via pastHomographs

	// podupaść - pure "fall" only (no graze variant)
	"podupaść": {stem: "podupad"},

	// Suppletive verbs - completely irregular stems
	"być": {stem: "by"},
	"iść": {masc: "szed", fem: "sz"},

	// jeść → jadł (suppletive stem jad-/jedl-)
	"jeść":     {stem: "jad", virile: "jed"},
	"nadojeść": {stem: "nadojad", virile: "nadojed"},

	// wziąć → wziął/wzięła (ą→ę alternation)
	"wziąć":          {masc: "wzią", fem: "wzię"},
	"przedsięwziąć":  {masc: "przedsięwzią", fem: "przedsięwzię"},

	// jąć → jął/jęła (ą→ę alternation)
	"jąć":     {masc: "ją", fem: "ję"},
	"zdjąć":   {masc: "zdją", fem: "zdję"},
	"rozdjąć": {masc: "rozdją", fem: "rozdję"},

	// miąć → miął/mięła (ą→ę alternation)
	"miąć": {masc: "mią", fem: "mię"},

	// nająć → najął/najęła
	"nająć": {masc: "nają", fem: "naję"},

	// dąć → dął/dęła
	"dąć": {masc: "dą", fem: "dę"},

	// ciąć → ciął/cięła
	"ciąć":  {masc: "cią", fem: "cię"},
	"ściąć": {masc: "ścią", fem: "ścię"},

	// giąć → giął/gięła
	"giąć": {masc: "gią", fem: "gię"},

	// piąć → piął/pięła
	"piąć":   {masc: "pią", fem: "pię"},
	"wspiąć": {masc: "wspią", fem: "wspię"},

	// żąć → żął/żęła
	"żąć": {masc: "żą", fem: "żę"},

	// kląć → klął/klęła
	"kląć": {masc: "klą", fem: "klę"},

	// siąść → siadł/siadła (special, ą→a)
	"siąść": {stem: "siad", virile: "sied"},

	// Note: paść is handled via pastHomographs (both "graze" and "fall" patterns)

	// kraść → kradł (suppletive stem krad-)
	"kraść": {stem: "krad"},

	// kłaść → kładł (suppletive stem kład-)
	"kłaść": {stem: "kład"},

	// prząść → prządł/przędła (ą→ę in masculine vs feminine)
	"prząść": {masc: "prząd", fem: "przęd"},

	// gryźć special virile forms
	"gryźć": {stem: "gryz", virile: "gryź"},

	// leźć special virile forms
	"leźć": {stem: "laz", virile: "leź"},
	"liźć": {stem: "laz", virile: "leź"},

	// wieźć - ó→o alternation (ó only in sg3m)
	"wieźć": {stem: "wioz", virile: "wieź", sg3m: "wiózł"},

	// nieść special ó→o alternation
	"nieść": {stem: "nios", virile: "nieś", sg3m: "niósł"},

	// pleść → plótł/plotła (ó only in sg3m)
	"pleść": {stem: "plot", virile: "plet", sg3m: "plótł"},

	// grześć → grzebł/grzebła (suppletive stem grzeb-)
	"grześć": {stem: "grzeb"},

	// tłuc → tłukł/tłukła
	"tłuc": {stem: "tłuk"},

	// przeć → parł/parła (suppletive stem par-)
	"przeć": {stem: "par"},

	// wrzeć → wrzał/wrzała (for boiling)
	"wrzeć": {stem: "wrza", virile: "wrze"},

	// zawrzeć → zawarł (suppletive stem war-)
	"zawrzeć":  {stem: "zawar"},
	"wywrzeć":  {stem: "wywar"},
	"dowrzeć":  {stem: "dowar"},
	"zewrzeć":  {stem: "zwar"},
	"odewrzeć": {stem: "odewar"},

	// trzeć → tarł/tarła (e→a alternation)
	"trzeć": {stem: "tar"},
	"drzeć": {stem: "dar"},
	"mrzeć": {stem: "mar"},
	"umrzeć": {stem: "umar"},

	// mleć → mełł/mełła (grind - suppletive stem with łł gemination)
	"mleć": {stem: "meł"},
	"pleć": {stem: "peł"},

	// żec → żegł/żegła (burn/sting)
	"żec": {stem: "żeg"},

	// podżec: asymmetric epenthetic - sg3m strips, others keep
	"podżec": {stem: "podeżg", sg3m: "podżegł"},
	"rozżec": {stem: "rozeżg", sg3m: "rozżegł"},
	"zżec":   {stem: "zeżg", sg3m: "zżegł"},

	// Compound prefixed verbs that need explicit entries
	"spostrzec":  {stem: "spostrzeg"},
	"zapobiec":   {stem: "zapobieg"},
	"współubiec": {stem: "współubieg"},

	// sprzeć → sprzał (NOT sparł - different from s+przeć)
	"sprzeć": {stem: "sprza", virile: "sprze"},

	// zeprzeć → sparł (ze- assimilates to s-)
	"zeprzeć": {stem: "spar"},

	// zetrzeć → starł (ze- assimilates to s-)
	"zetrzeć": {stem: "star"},

	// wesprzeć → wsparł (we- assimilates to ws-)
	"wesprzeć": {stem: "wspar"},

	// wetrzeć → wtarł (we- assimilates to w-)
	"wetrzeć": {stem: "wtar"},

	// otworzyć family: special stem twar-
	"otworzyć":    {stem: "otwar"},
	"przetworzyć": {stem: "przetwar"},
	"roztworzyć":  {stem: "roztwar"},

	// prać epenthetic forms: keep epenthetic vowel throughout
	"obeprać":  {stem: "obepra"},
	"odeprać":  {stem: "odepra"},
	"podeprać": {stem: "podepra"},

	// wejść → wszedł (special we- prefix)
	"wejść": {masc: "wszed", fem: "wesz"},

	// wznijść (archaic) → wzeszedł
	"wznijść": {masc: "wzeszed", fem: "wzesz"},

	// żreć → żarł/żarła (e→a alternation)
	"żreć": {stem: "żar"},

	// źreć → ziarł (suppletive stem ziar-)
	"źreć":    {stem: "ziar"},
	"zeźreć":  {stem: "zziar"},
	"zeźrzeć": {stem: "zziar"},
	"zrzeć":   {stem: "żar"},

	// przywrzeć → przywarł (uses war- stem)
	"przywrzeć": {stem: "przywar"},

	// rozeprzeć → rozeprzał (keeps epenthetic e)
	"rozeprzeć": {stem: "rozeprza", virile: "rozeprze"},

	// rozewrzeć → rozewrzał (keeps epenthetic e)
	"rozewrzeć": {stem: "rozewrza", virile: "rozewrze"},

	// rozpostrzeć → rozpostarł (post- + trzeć → -tar- stem)
	"rozpostrzeć": {stem: "rozpostar"},

	// krzywoprzysięgnąć → krzywoprzysiągł (drops n, ę→ą in masculine)
	"krzywoprzysięgnąć": {masc: "krzywoprzysiąg", fem: "krzywoprzysięg"},

	// nagadnąć → nagadnął/nagadnęła (n-KEPT)
	"nagadnąć": {masc: "nagadną", fem: "nagadnę"},
	"zagadnąć": {masc: "zagadną", fem: "zagadnę"},

	// rymsnąć → rymsnął/rymsnęła but rymsli (MIXED n-drop)
	"rymsnąć": {masc: "rymsną", fem: "rymsnę", sg3m: "rymsł", virile: "ryms"},

	// zastrzęgnąć → zastrzęgł (n-dropped, NO ę→ą alternation in masculine)
	"zastrzęgnąć": {stem: "zastrzęg"},

	// przeschnąć → przesechł (epenthetic 'e' for prze+schnąć)
	"przeschnąć": {stem: "przesch", sg3m: "przesechł"},

	// wpółgasnąć → wpółgasł (n-dropped)
	"wpółgasnąć": {stem: "wpółgas", virile: "wpółgaś"},

	// wskrzesnąć → wskrzesł (n-dropped)
	"wskrzesnąć": {stem: "wskrzes", virile: "wskrześ"},

	// przyosłabnąć → przyosłabł (n-dropped)
	"przyosłabnąć": {stem: "przyosłab"},

	// zmierzchnąć → zmierzchł (n-dropped)
	"zmierzchnąć": {stem: "zmierzch"},

	// zabrzęknąć → zabrzęknął/zabrzęknęła but zabrzękli (MIXED n-drop)
	"zabrzęknąć": {masc: "zabrzękną", fem: "zabrzęknę", sg3m: "zabrzękł", virile: "zabrzęk"},

	// oślizgnąć → oślizgnął/oślizgnęła but oślizgli (MIXED n-drop)
	"oślizgnąć": {masc: "oślizgną", fem: "oślizgnę", sg3m: "oślizgł", virile: "oślizg"},

	// obślizgnąć → obślizgnął/obślizgnęła but obślizgli (MIXED n-drop)
	"obślizgnąć": {masc: "obślizgną", fem: "obślizgnę", virile: "obślizg"},

	// brać → brał (simple -ać, but included as base for prefixes)
	"brać": {stem: "bra"},
	"prać": {stem: "pra"},

	// dać → dał
	"dać": {stem: "da"},

	// stać → stał (to stand)
	"stać": {stem: "sta"},

	// mieć → miał
	"mieć": {stem: "mia", virile: "mie"},

	// chcieć → chciał
	"chcieć": {stem: "chcia", virile: "chcie"},

	// wiedzieć → wiedział
	"wiedzieć": {stem: "wiedzia", virile: "wiedzie"},

	// siedzieć → siedział
	"siedzieć": {stem: "siedzia", virile: "siedzie"},

	// widzieć → widział
	"widzieć": {stem: "widzia", virile: "widzie"},

	// słyszeć → słyszał
	"słyszeć": {stem: "słysza", virile: "słysze"},

	// musieć → musiał
	"musieć": {stem: "musia", virile: "musie"},

	// móc → mógł/mogła (ó→o alternation)
	"móc": {stem: "mog", sg3m: "mógł"},

	// biec → biegł (suppletive stem bieg-)
	"biec": {stem: "bieg"},

	// lec → legł
	"lec": {stem: "leg"},

	// rzec → rzekł
	"rzec": {stem: "rzek"},

	// ciec → ciekł
	"ciec": {stem: "ciek"},

	// strzec → strzegł
	"strzec": {stem: "strzeg"},

	// przesiąc → przesiąkł (archaic shortened form of przesiąknąć)
	"przesiąc": {stem: "przesiąk"},

	// schnąć → sechł/schła (epenthetic 'e' ONLY in sg3m)
	"schnąć": {stem: "sch", sg3m: "sechł"},

	// przysięgnąć → przysiągł/przysięgła (ę→ą alternation in masculine)
	"przysięgnąć": {masc: "przysiąg", fem: "przysięg"},

	// piec → piekł
	"piec": {stem: "piek"},

	// wlec: handled via pastHomographs

	// siąść family
	"usiąść":  {stem: "usiad", virile: "usied"},
	"wysieść": {stem: "wysiad", virile: "wysied"},
	"zsieść":  {stem: "zsiad", virile: "zsied"},

	// rosnąć → rósł/rosła (special n-dropping with ó→o)
	"rosnąć": {stem: "ros", virile: "roś", sg3m: "rósł"},
	"rość":   {stem: "ros", virile: "roś", sg3m: "rósł"},
}

// buildJscPast builds past tense for -jść verbs (prefixed iść).
// przejść → przeszedł/przeszła, wyjść → wyszedł/wyszła
func buildJscPast(prefix string) PastTense {
	return pastSpec{masc: prefix + "szed", fem: prefix + "sz"}.build()
}

// buildSchnacPast builds past tense for prefixed schnąć verbs.
// schnąć has an asymmetric stem: sg3m uses "sechł" (epenthetic e), others use "schł".
// Additionally, prefixes with epenthetic vowels (obe-, pode-, roze-, ze-)
// strip the vowel in sg3m but keep it in other forms.
// obeschnąć: sg1m=obeschłem, sg3m=obsechł, sg3f=obeschła
func buildSchnacPast(infinitive string) PastTense {
	// Determine the prefix from the infinitive
	prefix := strings.TrimSuffix(infinitive, "schnąć")

	// For sg3m only: strip epenthetic from prefix, use sechł stem
	sg3mPrefix := prefix
	epenthetic := map[string]string{
		"obe": "ob", "pode": "pod", "roze": "roz", "ze": "z",
	}
	for full, stripped := range epenthetic {
		if strings.HasSuffix(prefix, full) {
			sg3mPrefix = prefix[:len(prefix)-len(full)] + stripped
			break
		}
	}

	// Special case: zeschnąć has unusual assimilation
	// sg3m = ssechł (ze → s before s-stem, then s+sechł = ssechł)
	if infinitive == "zeschnąć" {
		return pastSpec{stem: "zesch", sg3m: "ssechł"}.build()
	}

	return pastSpec{stem: prefix + "sch", sg3m: sg3mPrefix + "sechł"}.build()
}

// stripEpentheticVowel removes the epenthetic 'e' from prefixes when applying
// them to past tense forms. The epenthetic vowel appears in infinitives before
// consonant clusters but disappears in conjugated forms.
// ze + drzeć → infinitive zedrzeć, but past zdarł (not zedarł)
// However, some clusters require keeping the vowel:
// ze + siąść → zesiąść, past zesiadł (not zsiadł, 'zs' is unpronounceable)
func stripEpentheticVowel(prefix string, baseForm string) string {
	epenthetic := map[string]string{
		"ze": "z", "we": "w", "ode": "od", "obe": "ob",
		"pode": "pod", "nade": "nad", "roze": "roz", "wze": "wz",
	}
	stripped, ok := epenthetic[prefix]
	if !ok {
		return prefix
	}

	baseFirstChar := rune(0)
	if len(baseForm) > 0 {
		baseFirstChar = []rune(baseForm)[0]
	}

	// Special case: schnąć → sechł. The past stem "sech" is pronounceable
	// after prefixes (obsechł, podsechł, rozsechł) even though it starts with 's'.
	// The epenthetic vowel was needed for the infinitive (obeschnąć) because
	// "obschnąć" would have an unpronounceable "bschn" cluster.
	if strings.HasPrefix(baseForm, "sech") {
		return stripped
	}

	// Don't strip if it would create an unpronounceable or unusual cluster
	// e.g., ze + siadł → zesiadł (not zsiadł)
	// e.g., ze + brał → zebrał (not zbrał)
	// The epenthetic vowel is kept before: s, ś, z, ź, ż, b, p, w
	// NOTE: 'm' is NOT in this list because "zm" is a common, easy cluster (zmarł, zmełł)
	if prefix == "ze" {
		keepVowel := map[rune]bool{
			's': true, 'ś': true, 'z': true, 'ź': true, 'ż': true,
			'b': true, 'p': true, 'w': true,
		}
		if keepVowel[baseFirstChar] {
			return prefix
		}
	}
	// Similar for other prefixes with epenthetic vowels before 'b' or 's'
	if prefix == "ode" || prefix == "pode" || prefix == "nade" || prefix == "obe" || prefix == "we" || prefix == "roze" {
		if baseFirstChar == 'b' || baseFirstChar == 's' || baseFirstChar == 'ś' {
			return prefix
		}
	}
	return stripped
}

// applyPrefixToPast applies a prefix to all forms of a past tense paradigm.
// Strips epenthetic vowels from prefixes before applying.
func applyPrefixToPast(prefix string, base PastTense) PastTense {
	// Pass the base sg3m form to decide about epenthetic vowel
	p := stripEpentheticVowel(prefix, base.Sg3M)
	return PastTense{
		Sg1M:  p + base.Sg1M,
		Sg1F:  p + base.Sg1F,
		Sg2M:  p + base.Sg2M,
		Sg2F:  p + base.Sg2F,
		Sg3M:  p + base.Sg3M,
		Sg3F:  p + base.Sg3F,
		Sg3N:  p + base.Sg3N,
		Pl1V:  p + base.Pl1V,
		Pl1NV: p + base.Pl1NV,
		Pl2V:  p + base.Pl2V,
		Pl2NV: p + base.Pl2NV,
		Pl3V:  p + base.Pl3V,
		Pl3NV: p + base.Pl3NV,
	}
}
