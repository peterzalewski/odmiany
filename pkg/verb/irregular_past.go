package verb

import "strings"

// pastHomographs contains verbs with multiple valid past tense paradigms.
var pastHomographs = map[string][]PastParadigm{
	// paść: "to fall" (padł) vs "to graze" (pasł)
	"paść": {
		{
			PastTense: PastTense{
				Sg1M: "pasłem", Sg1F: "pasłam",
				Sg2M: "pasłeś", Sg2F: "pasłaś",
				Sg3M: "pasł", Sg3F: "pasła", Sg3N: "pasło",
				Pl1V: "paśliśmy", Pl1NV: "pasłyśmy",
				Pl2V: "paśliście", Pl2NV: "pasłyście",
				Pl3V: "paśli", Pl3NV: "pasły",
			},
			Gloss: "to graze (animals)",
		},
		{
			PastTense: PastTense{
				Sg1M: "padłem", Sg1F: "padłam",
				Sg2M: "padłeś", Sg2F: "padłaś",
				Sg3M: "padł", Sg3F: "padła", Sg3N: "padło",
				Pl1V: "padliśmy", Pl1NV: "padłyśmy",
				Pl2V: "padliście", Pl2NV: "padłyście",
				Pl3V: "padli", Pl3NV: "padły",
			},
			Gloss: "to fall",
		},
	},
}

// buildPascHomograph creates homograph entries for prefixed -paść verbs.
// These verbs have two valid paradigms (fall vs graze pattern) plus a mixed form.
func buildPascHomograph(prefix string) []PastParadigm {
	return []PastParadigm{
		// Pure "fall" pattern
		{
			PastTense: PastTense{
				Sg1M:  prefix + "padłem", Sg1F: prefix + "padłam",
				Sg2M:  prefix + "padłeś", Sg2F: prefix + "padłaś",
				Sg3M:  prefix + "padł", Sg3F: prefix + "padła", Sg3N: prefix + "padło",
				Pl1V:  prefix + "padliśmy", Pl1NV: prefix + "padłyśmy",
				Pl2V:  prefix + "padliście", Pl2NV: prefix + "padłyście",
				Pl3V:  prefix + "padli", Pl3NV: prefix + "padły",
			},
			Gloss: "to fall",
		},
		// Mixed pattern: masc sg "graze", fem/pl "fall"
		{
			PastTense: PastTense{
				Sg1M:  prefix + "pasłem", Sg1F: prefix + "padłam",
				Sg2M:  prefix + "pasłeś", Sg2F: prefix + "padłaś",
				Sg3M:  prefix + "pasł", Sg3F: prefix + "padła", Sg3N: prefix + "padło",
				Pl1V:  prefix + "padliśmy", Pl1NV: prefix + "padłyśmy",
				Pl2V:  prefix + "padliście", Pl2NV: prefix + "padłyście",
				Pl3V:  prefix + "padli", Pl3NV: prefix + "padły",
			},
			Gloss: "to fall (variant)",
		},
	}
}

func init() {
	// Add homographs for prefixed -paść verbs
	prefixes := []string{"do", "na", "od", "o", "pod", "po", "prze", "przy", "roz", "s", "u", "w", "wy", "za", "zaprze"}
	for _, p := range prefixes {
		pastHomographs[p+"paść"] = buildPascHomograph(p)
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
	"niknąć":     true,
	"pachnąć":    true,
}

// irregularPastVerbs contains past tense paradigms for verbs that cannot
// be conjugated by heuristics alone.
var irregularPastVerbs = map[string]PastTense{
	// === -jść verbs (special stem: pójść → poszedł, not pószedł) ===
	"pójść": {
		Sg1M: "poszedłem", Sg1F: "poszłam",
		Sg2M: "poszedłeś", Sg2F: "poszłaś",
		Sg3M: "poszedł", Sg3F: "poszła", Sg3N: "poszło",
		Pl1V: "poszliśmy", Pl1NV: "poszłyśmy",
		Pl2V: "poszliście", Pl2NV: "poszłyście",
		Pl3V: "poszli", Pl3NV: "poszły",
	},
	"obejść": {
		Sg1M: "obszedłem", Sg1F: "obeszłam",
		Sg2M: "obszedłeś", Sg2F: "obeszłaś",
		Sg3M: "obszedł", Sg3F: "obeszła", Sg3N: "obeszło",
		Pl1V: "obeszliśmy", Pl1NV: "obeszłyśmy",
		Pl2V: "obeszliście", Pl2NV: "obeszłyście",
		Pl3V: "obeszli", Pl3NV: "obeszły",
	},
	"odejść": {
		Sg1M: "odszedłem", Sg1F: "odeszłam",
		Sg2M: "odszedłeś", Sg2F: "odeszłaś",
		Sg3M: "odszedł", Sg3F: "odeszła", Sg3N: "odeszło",
		Pl1V: "odeszliśmy", Pl1NV: "odeszłyśmy",
		Pl2V: "odeszliście", Pl2NV: "odeszłyście",
		Pl3V: "odeszli", Pl3NV: "odeszły",
	},
	"podejść": {
		Sg1M: "podszedłem", Sg1F: "podeszłam",
		Sg2M: "podszedłeś", Sg2F: "podeszłaś",
		Sg3M: "podszedł", Sg3F: "podeszła", Sg3N: "podeszło",
		Pl1V: "podeszliśmy", Pl1NV: "podeszłyśmy",
		Pl2V: "podeszliście", Pl2NV: "podeszłyście",
		Pl3V: "podeszli", Pl3NV: "podeszły",
	},
	"nadejść": {
		Sg1M: "nadszedłem", Sg1F: "nadeszłam",
		Sg2M: "nadszedłeś", Sg2F: "nadeszłaś",
		Sg3M: "nadszedł", Sg3F: "nadeszła", Sg3N: "nadeszło",
		Pl1V: "nadeszliśmy", Pl1NV: "nadeszłyśmy",
		Pl2V: "nadeszliście", Pl2NV: "nadeszłyście",
		Pl3V: "nadeszli", Pl3NV: "nadeszły",
	},
	"rozejść": {
		Sg1M: "rozszedłem", Sg1F: "rozeszłam",
		Sg2M: "rozszedłeś", Sg2F: "rozeszłaś",
		Sg3M: "rozszedł", Sg3F: "rozeszła", Sg3N: "rozeszło",
		Pl1V: "rozeszliśmy", Pl1NV: "rozeszłyśmy",
		Pl2V: "rozeszliście", Pl2NV: "rozeszłyście",
		Pl3V: "rozeszli", Pl3NV: "rozeszły",
	},

	// === -naleźć verbs (special stem: odnaleźć → odnalazł) ===
	"znaleźć": {
		Sg1M: "znalazłem", Sg1F: "znalazłam",
		Sg2M: "znalazłeś", Sg2F: "znalazłaś",
		Sg3M: "znalazł", Sg3F: "znalazła", Sg3N: "znalazło",
		Pl1V: "znaleźliśmy", Pl1NV: "znalazłyśmy",
		Pl2V: "znaleźliście", Pl2NV: "znalazłyście",
		Pl3V: "znaleźli", Pl3NV: "znalazły",
	},
	"odnaleźć": {
		Sg1M: "odnalazłem", Sg1F: "odnalazłam",
		Sg2M: "odnalazłeś", Sg2F: "odnalazłaś",
		Sg3M: "odnalazł", Sg3F: "odnalazła", Sg3N: "odnalazło",
		Pl1V: "odnaleźliśmy", Pl1NV: "odnalazłyśmy",
		Pl2V: "odnaleźliście", Pl2NV: "odnalazłyście",
		Pl3V: "odnaleźli", Pl3NV: "odnalazły",
	},
	"wynaleźć": {
		Sg1M: "wynalazłem", Sg1F: "wynalazłam",
		Sg2M: "wynalazłeś", Sg2F: "wynalazłaś",
		Sg3M: "wynalazł", Sg3F: "wynalazła", Sg3N: "wynalazło",
		Pl1V: "wynaleźliśmy", Pl1NV: "wynalazłyśmy",
		Pl2V: "wynaleźliście", Pl2NV: "wynalazłyście",
		Pl3V: "wynaleźli", Pl3NV: "wynalazły",
	},

	// === -sieść verbs (sieść → siadł, special alternation) ===
	"sieść": {
		Sg1M: "siadłem", Sg1F: "siadłam",
		Sg2M: "siadłeś", Sg2F: "siadłaś",
		Sg3M: "siadł", Sg3F: "siadła", Sg3N: "siadło",
		Pl1V: "siedliśmy", Pl1NV: "siadłyśmy",
		Pl2V: "siedliście", Pl2NV: "siadłyście",
		Pl3V: "siedli", Pl3NV: "siadły",
	},
	"osiąść": {
		Sg1M: "osiadłem", Sg1F: "osiadłam",
		Sg2M: "osiadłeś", Sg2F: "osiadłaś",
		Sg3M: "osiadł", Sg3F: "osiadła", Sg3N: "osiadło",
		Pl1V: "osiedliśmy", Pl1NV: "osiadłyśmy",
		Pl2V: "osiedliście", Pl2NV: "osiadłyście",
		Pl3V: "osiedli", Pl3NV: "osiadły",
	},
	"posieść": {
		Sg1M: "posiadłem", Sg1F: "posiadłam",
		Sg2M: "posiadłeś", Sg2F: "posiadłaś",
		Sg3M: "posiadł", Sg3F: "posiadła", Sg3N: "posiadło",
		Pl1V: "posiedliśmy", Pl1NV: "posiadłyśmy",
		Pl2V: "posiedliście", Pl2NV: "posiadłyście",
		Pl3V: "posiedli", Pl3NV: "posiadły",
	},
	"osieść": {
		Sg1M: "osiadłem", Sg1F: "osiadłam",
		Sg2M: "osiadłeś", Sg2F: "osiadłaś",
		Sg3M: "osiadł", Sg3F: "osiadła", Sg3N: "osiadło",
		Pl1V: "osiedliśmy", Pl1NV: "osiadłyśmy",
		Pl2V: "osiedliście", Pl2NV: "osiadłyście",
		Pl3V: "osiedli", Pl3NV: "osiadły",
	},

	// Note: -paść verbs (dopaść, napaść, etc.) are handled via pastHomographs

	// podupaść - pure "fall" only (no graze variant)
	"podupaść": {
		Sg1M: "podupadłem", Sg1F: "podupadłam",
		Sg2M: "podupadłeś", Sg2F: "podupadłaś",
		Sg3M: "podupadł", Sg3F: "podupadła", Sg3N: "podupadło",
		Pl1V: "podupadliśmy", Pl1NV: "podupadłyśmy",
		Pl2V: "podupadliście", Pl2NV: "podupadłyście",
		Pl3V: "podupadli", Pl3NV: "podupadły",
	},

	// Suppletive verbs - completely irregular stems

	// być → był (suppletive)
	"być": {
		Sg1M: "byłem", Sg1F: "byłam",
		Sg2M: "byłeś", Sg2F: "byłaś",
		Sg3M: "był", Sg3F: "była", Sg3N: "było",
		Pl1V: "byliśmy", Pl1NV: "byłyśmy",
		Pl2V: "byliście", Pl2NV: "byłyście",
		Pl3V: "byli", Pl3NV: "były",
	},

	// iść → szedł (suppletive stem szed-/sz-)
	"iść": {
		Sg1M: "szedłem", Sg1F: "szłam",
		Sg2M: "szedłeś", Sg2F: "szłaś",
		Sg3M: "szedł", Sg3F: "szła", Sg3N: "szło",
		Pl1V: "szliśmy", Pl1NV: "szłyśmy",
		Pl2V: "szliście", Pl2NV: "szłyście",
		Pl3V: "szli", Pl3NV: "szły",
	},

	// jeść → jadł (suppletive stem jad-/jedl-)
	"jeść": {
		Sg1M: "jadłem", Sg1F: "jadłam",
		Sg2M: "jadłeś", Sg2F: "jadłaś",
		Sg3M: "jadł", Sg3F: "jadła", Sg3N: "jadło",
		Pl1V: "jedliśmy", Pl1NV: "jadłyśmy",
		Pl2V: "jedliście", Pl2NV: "jadłyście",
		Pl3V: "jedli", Pl3NV: "jadły",
	},

	// nadojeść → nadojadł (nad + o + jeść)
	"nadojeść": {
		Sg1M: "nadojadłem", Sg1F: "nadojadłam",
		Sg2M: "nadojadłeś", Sg2F: "nadojadłaś",
		Sg3M: "nadojadł", Sg3F: "nadojadła", Sg3N: "nadojadło",
		Pl1V: "nadojedliśmy", Pl1NV: "nadojadłyśmy",
		Pl2V: "nadojedliście", Pl2NV: "nadojadłyście",
		Pl3V: "nadojedli", Pl3NV: "nadojadły",
	},

	// wziąć → wziął/wzięła (ą→ę alternation)
	"wziąć": {
		Sg1M: "wziąłem", Sg1F: "wzięłam",
		Sg2M: "wziąłeś", Sg2F: "wzięłaś",
		Sg3M: "wziął", Sg3F: "wzięła", Sg3N: "wzięło",
		Pl1V: "wzięliśmy", Pl1NV: "wzięłyśmy",
		Pl2V: "wzięliście", Pl2NV: "wzięłyście",
		Pl3V: "wzięli", Pl3NV: "wzięły",
	},

	// przedsięwziąć → przedsięwziął/przedsięwzięła (przed + się + wziąć)
	"przedsięwziąć": {
		Sg1M: "przedsięwziąłem", Sg1F: "przedsięwzięłam",
		Sg2M: "przedsięwziąłeś", Sg2F: "przedsięwzięłaś",
		Sg3M: "przedsięwziął", Sg3F: "przedsięwzięła", Sg3N: "przedsięwzięło",
		Pl1V: "przedsięwzięliśmy", Pl1NV: "przedsięwzięłyśmy",
		Pl2V: "przedsięwzięliście", Pl2NV: "przedsięwzięłyście",
		Pl3V: "przedsięwzięli", Pl3NV: "przedsięwzięły",
	},

	// jąć → jął/jęła (ą→ę alternation)
	"jąć": {
		Sg1M: "jąłem", Sg1F: "jęłam",
		Sg2M: "jąłeś", Sg2F: "jęłaś",
		Sg3M: "jął", Sg3F: "jęła", Sg3N: "jęło",
		Pl1V: "jęliśmy", Pl1NV: "jęłyśmy",
		Pl2V: "jęliście", Pl2NV: "jęłyście",
		Pl3V: "jęli", Pl3NV: "jęły",
	},

	// zdjąć → zdjął/zdjęła (z + djąć variant of jąć)
	"zdjąć": {
		Sg1M: "zdjąłem", Sg1F: "zdjęłam",
		Sg2M: "zdjąłeś", Sg2F: "zdjęłaś",
		Sg3M: "zdjął", Sg3F: "zdjęła", Sg3N: "zdjęło",
		Pl1V: "zdjęliśmy", Pl1NV: "zdjęłyśmy",
		Pl2V: "zdjęliście", Pl2NV: "zdjęłyście",
		Pl3V: "zdjęli", Pl3NV: "zdjęły",
	},

	// rozdjąć → rozdjął/rozdjęła (roz + djąć variant of jąć)
	"rozdjąć": {
		Sg1M: "rozdjąłem", Sg1F: "rozdjęłam",
		Sg2M: "rozdjąłeś", Sg2F: "rozdjęłaś",
		Sg3M: "rozdjął", Sg3F: "rozdjęła", Sg3N: "rozdjęło",
		Pl1V: "rozdjęliśmy", Pl1NV: "rozdjęłyśmy",
		Pl2V: "rozdjęliście", Pl2NV: "rozdjęłyście",
		Pl3V: "rozdjęli", Pl3NV: "rozdjęły",
	},

	// miąć → miął/mięła (ą→ę alternation)
	"miąć": {
		Sg1M: "miąłem", Sg1F: "mięłam",
		Sg2M: "miąłeś", Sg2F: "mięłaś",
		Sg3M: "miął", Sg3F: "mięła", Sg3N: "mięło",
		Pl1V: "mięliśmy", Pl1NV: "mięłyśmy",
		Pl2V: "mięliście", Pl2NV: "mięłyście",
		Pl3V: "mięli", Pl3NV: "mięły",
	},

	// nająć → najął/najęła (ą→ę alternation)
	"nająć": {
		Sg1M: "nająłem", Sg1F: "najęłam",
		Sg2M: "nająłeś", Sg2F: "najęłaś",
		Sg3M: "najął", Sg3F: "najęła", Sg3N: "najęło",
		Pl1V: "najęliśmy", Pl1NV: "najęłyśmy",
		Pl2V: "najęliście", Pl2NV: "najęłyście",
		Pl3V: "najęli", Pl3NV: "najęły",
	},

	// dąć → dął/dęła (ą→ę alternation, but also dął/dęła)
	"dąć": {
		Sg1M: "dąłem", Sg1F: "dęłam",
		Sg2M: "dąłeś", Sg2F: "dęłaś",
		Sg3M: "dął", Sg3F: "dęła", Sg3N: "dęło",
		Pl1V: "dęliśmy", Pl1NV: "dęłyśmy",
		Pl2V: "dęliście", Pl2NV: "dęłyście",
		Pl3V: "dęli", Pl3NV: "dęły",
	},

	// ciąć → ciął/cięła (ą→ę alternation)
	"ciąć": {
		Sg1M: "ciąłem", Sg1F: "cięłam",
		Sg2M: "ciąłeś", Sg2F: "cięłaś",
		Sg3M: "ciął", Sg3F: "cięła", Sg3N: "cięło",
		Pl1V: "cięliśmy", Pl1NV: "cięłyśmy",
		Pl2V: "cięliście", Pl2NV: "cięłyście",
		Pl3V: "cięli", Pl3NV: "cięły",
	},

	// ściąć → ściął/ścięła (ś + ciąć)
	"ściąć": {
		Sg1M: "ściąłem", Sg1F: "ścięłam",
		Sg2M: "ściąłeś", Sg2F: "ścięłaś",
		Sg3M: "ściął", Sg3F: "ścięła", Sg3N: "ścięło",
		Pl1V: "ścięliśmy", Pl1NV: "ścięłyśmy",
		Pl2V: "ścięliście", Pl2NV: "ścięłyście",
		Pl3V: "ścięli", Pl3NV: "ścięły",
	},

	// giąć → giął/gięła (ą→ę alternation)
	"giąć": {
		Sg1M: "giąłem", Sg1F: "gięłam",
		Sg2M: "giąłeś", Sg2F: "gięłaś",
		Sg3M: "giął", Sg3F: "gięła", Sg3N: "gięło",
		Pl1V: "gięliśmy", Pl1NV: "gięłyśmy",
		Pl2V: "gięliście", Pl2NV: "gięłyście",
		Pl3V: "gięli", Pl3NV: "gięły",
	},

	// piąć → piął/pięła (ą→ę alternation)
	"piąć": {
		Sg1M: "piąłem", Sg1F: "pięłam",
		Sg2M: "piąłeś", Sg2F: "pięłaś",
		Sg3M: "piął", Sg3F: "pięła", Sg3N: "pięło",
		Pl1V: "pięliśmy", Pl1NV: "pięłyśmy",
		Pl2V: "pięliście", Pl2NV: "pięłyście",
		Pl3V: "pięli", Pl3NV: "pięły",
	},

	// wspiąć → wspiął/wspięła (w + spiąć variant of piąć)
	"wspiąć": {
		Sg1M: "wspiąłem", Sg1F: "wspięłam",
		Sg2M: "wspiąłeś", Sg2F: "wspięłaś",
		Sg3M: "wspiął", Sg3F: "wspięła", Sg3N: "wspięło",
		Pl1V: "wspięliśmy", Pl1NV: "wspięłyśmy",
		Pl2V: "wspięliście", Pl2NV: "wspięłyście",
		Pl3V: "wspięli", Pl3NV: "wspięły",
	},

	// żąć → żął/żęła (ą→ę alternation)
	"żąć": {
		Sg1M: "żąłem", Sg1F: "żęłam",
		Sg2M: "żąłeś", Sg2F: "żęłaś",
		Sg3M: "żął", Sg3F: "żęła", Sg3N: "żęło",
		Pl1V: "żęliśmy", Pl1NV: "żęłyśmy",
		Pl2V: "żęliście", Pl2NV: "żęłyście",
		Pl3V: "żęli", Pl3NV: "żęły",
	},

	// kląć → klął/klęła (ą→ę alternation)
	"kląć": {
		Sg1M: "kląłem", Sg1F: "klęłam",
		Sg2M: "kląłeś", Sg2F: "klęłaś",
		Sg3M: "klął", Sg3F: "klęła", Sg3N: "klęło",
		Pl1V: "klęliśmy", Pl1NV: "klęłyśmy",
		Pl2V: "klęliście", Pl2NV: "klęłyście",
		Pl3V: "klęli", Pl3NV: "klęły",
	},

	// siąść → siadł/siadła (special, ą→a)
	"siąść": {
		Sg1M: "siadłem", Sg1F: "siadłam",
		Sg2M: "siadłeś", Sg2F: "siadłaś",
		Sg3M: "siadł", Sg3F: "siadła", Sg3N: "siadło",
		Pl1V: "siedliśmy", Pl1NV: "siadłyśmy",
		Pl2V: "siedliście", Pl2NV: "siadłyście",
		Pl3V: "siedli", Pl3NV: "siadły",
	},

	// Note: paść is handled via pastHomographs (both "graze" and "fall" patterns)

	// kraść → kradł (suppletive stem krad-)
	"kraść": {
		Sg1M: "kradłem", Sg1F: "kradłam",
		Sg2M: "kradłeś", Sg2F: "kradłaś",
		Sg3M: "kradł", Sg3F: "kradła", Sg3N: "kradło",
		Pl1V: "kradliśmy", Pl1NV: "kradłyśmy",
		Pl2V: "kradliście", Pl2NV: "kradłyście",
		Pl3V: "kradli", Pl3NV: "kradły",
	},

	// kłaść → kładł (suppletive stem kład-)
	"kłaść": {
		Sg1M: "kładłem", Sg1F: "kładłam",
		Sg2M: "kładłeś", Sg2F: "kładłaś",
		Sg3M: "kładł", Sg3F: "kładła", Sg3N: "kładło",
		Pl1V: "kładliśmy", Pl1NV: "kładłyśmy",
		Pl2V: "kładliście", Pl2NV: "kładłyście",
		Pl3V: "kładli", Pl3NV: "kładły",
	},

	// prząść → prządł/przędła (suppletive stem prząd-/przęd-, ą→ę in fem/neut/virile)
	"prząść": {
		Sg1M: "prządłem", Sg1F: "przędłam",
		Sg2M: "prządłeś", Sg2F: "przędłaś",
		Sg3M: "prządł", Sg3F: "przędła", Sg3N: "przędło",
		Pl1V: "przędliśmy", Pl1NV: "przędłyśmy",
		Pl2V: "przędliście", Pl2NV: "przędłyście",
		Pl3V: "przędli", Pl3NV: "przędły",
	},

	// gryźć special virile forms
	"gryźć": {
		Sg1M: "gryzłem", Sg1F: "gryzłam",
		Sg2M: "gryzłeś", Sg2F: "gryzłaś",
		Sg3M: "gryzł", Sg3F: "gryzła", Sg3N: "gryzło",
		Pl1V: "gryźliśmy", Pl1NV: "gryzłyśmy",
		Pl2V: "gryźliście", Pl2NV: "gryzłyście",
		Pl3V: "gryźli", Pl3NV: "gryzły",
	},

	// leźć special virile forms
	"leźć": {
		Sg1M: "lazłem", Sg1F: "lazłam",
		Sg2M: "lazłeś", Sg2F: "lazłaś",
		Sg3M: "lazł", Sg3F: "lazła", Sg3N: "lazło",
		Pl1V: "leźliśmy", Pl1NV: "lazłyśmy",
		Pl2V: "leźliście", Pl2NV: "lazłyście",
		Pl3V: "leźli", Pl3NV: "lazły",
	},

	// liźć - archaic variant of leźć
	"liźć": {
		Sg1M: "lazłem", Sg1F: "lazłam",
		Sg2M: "lazłeś", Sg2F: "lazłaś",
		Sg3M: "lazł", Sg3F: "lazła", Sg3N: "lazło",
		Pl1V: "leźliśmy", Pl1NV: "lazłyśmy",
		Pl2V: "leźliście", Pl2NV: "lazłyście",
		Pl3V: "leźli", Pl3NV: "lazły",
	},

	// wieźć - ó→o alternation (ó only in sg3m)
	"wieźć": {
		Sg1M: "wiozłem", Sg1F: "wiozłam",
		Sg2M: "wiozłeś", Sg2F: "wiozłaś",
		Sg3M: "wiózł", Sg3F: "wiozła", Sg3N: "wiozło",
		Pl1V: "wieźliśmy", Pl1NV: "wiozłyśmy",
		Pl2V: "wieźliście", Pl2NV: "wiozłyście",
		Pl3V: "wieźli", Pl3NV: "wiozły",
	},

	// nieść special ó→o alternation
	"nieść": {
		Sg1M: "niosłem", Sg1F: "niosłam",
		Sg2M: "niosłeś", Sg2F: "niosłaś",
		Sg3M: "niósł", Sg3F: "niosła", Sg3N: "niosło",
		Pl1V: "nieśliśmy", Pl1NV: "niosłyśmy",
		Pl2V: "nieśliście", Pl2NV: "niosłyście",
		Pl3V: "nieśli", Pl3NV: "niosły",
	},

	// pleść → plótł/plotła (ó only in sg3m)
	"pleść": {
		Sg1M: "plotłem", Sg1F: "plotłam",
		Sg2M: "plotłeś", Sg2F: "plotłaś",
		Sg3M: "plótł", Sg3F: "plotła", Sg3N: "plotło",
		Pl1V: "pletliśmy", Pl1NV: "plotłyśmy",
		Pl2V: "pletliście", Pl2NV: "plotłyście",
		Pl3V: "pletli", Pl3NV: "plotły",
	},

	// grześć → grzebł/grzebła (suppletive stem grzeb-)
	"grześć": {
		Sg1M: "grzebłem", Sg1F: "grzebłam",
		Sg2M: "grzebłeś", Sg2F: "grzebłaś",
		Sg3M: "grzebł", Sg3F: "grzebła", Sg3N: "grzebło",
		Pl1V: "grzebliśmy", Pl1NV: "grzebłyśmy",
		Pl2V: "grzebliście", Pl2NV: "grzebłyście",
		Pl3V: "grzebli", Pl3NV: "grzebły",
	},

	// tłuc → tłukł/tłukła
	"tłuc": {
		Sg1M: "tłukłem", Sg1F: "tłukłam",
		Sg2M: "tłukłeś", Sg2F: "tłukłaś",
		Sg3M: "tłukł", Sg3F: "tłukła", Sg3N: "tłukło",
		Pl1V: "tłukliśmy", Pl1NV: "tłukłyśmy",
		Pl2V: "tłukliście", Pl2NV: "tłukłyście",
		Pl3V: "tłukli", Pl3NV: "tłukły",
	},

	// przeć → parł/parła (suppletive stem par-)
	"przeć": {
		Sg1M: "parłem", Sg1F: "parłam",
		Sg2M: "parłeś", Sg2F: "parłaś",
		Sg3M: "parł", Sg3F: "parła", Sg3N: "parło",
		Pl1V: "parliśmy", Pl1NV: "parłyśmy",
		Pl2V: "parliście", Pl2NV: "parłyście",
		Pl3V: "parli", Pl3NV: "parły",
	},

	// wrzeć → wrał/wrała (for boiling) or wrzał/wrzała
	"wrzeć": {
		Sg1M: "wrzałem", Sg1F: "wrzałam",
		Sg2M: "wrzałeś", Sg2F: "wrzałaś",
		Sg3M: "wrzał", Sg3F: "wrzała", Sg3N: "wrzało",
		Pl1V: "wrzeliśmy", Pl1NV: "wrzałyśmy",
		Pl2V: "wrzeliście", Pl2NV: "wrzałyście",
		Pl3V: "wrzeli", Pl3NV: "wrzały",
	},

	// zawrzeć → zawarł (suppletive stem war-, different from wrzeć "to boil")
	"zawrzeć": {
		Sg1M: "zawarłem", Sg1F: "zawarłam",
		Sg2M: "zawarłeś", Sg2F: "zawarłaś",
		Sg3M: "zawarł", Sg3F: "zawarła", Sg3N: "zawarło",
		Pl1V: "zawarliśmy", Pl1NV: "zawarłyśmy",
		Pl2V: "zawarliście", Pl2NV: "zawarłyście",
		Pl3V: "zawarli", Pl3NV: "zawarły",
	},

	// wywrzeć → wywarł
	"wywrzeć": {
		Sg1M: "wywarłem", Sg1F: "wywarłam",
		Sg2M: "wywarłeś", Sg2F: "wywarłaś",
		Sg3M: "wywarł", Sg3F: "wywarła", Sg3N: "wywarło",
		Pl1V: "wywarliśmy", Pl1NV: "wywarłyśmy",
		Pl2V: "wywarliście", Pl2NV: "wywarłyście",
		Pl3V: "wywarli", Pl3NV: "wywarły",
	},

	// dowrzeć → dowarł
	"dowrzeć": {
		Sg1M: "dowarłem", Sg1F: "dowarłam",
		Sg2M: "dowarłeś", Sg2F: "dowarłaś",
		Sg3M: "dowarł", Sg3F: "dowarła", Sg3N: "dowarło",
		Pl1V: "dowarliśmy", Pl1NV: "dowarłyśmy",
		Pl2V: "dowarliście", Pl2NV: "dowarłyście",
		Pl3V: "dowarli", Pl3NV: "dowarły",
	},

	// zewrzeć → zwarł
	"zewrzeć": {
		Sg1M: "zwarłem", Sg1F: "zwarłam",
		Sg2M: "zwarłeś", Sg2F: "zwarłaś",
		Sg3M: "zwarł", Sg3F: "zwarła", Sg3N: "zwarło",
		Pl1V: "zwarliśmy", Pl1NV: "zwarłyśmy",
		Pl2V: "zwarliście", Pl2NV: "zwarłyście",
		Pl3V: "zwarli", Pl3NV: "zwarły",
	},

	// odewrzeć → odewarł (keeps epenthetic e before w)
	"odewrzeć": {
		Sg1M: "odewarłem", Sg1F: "odewarłam",
		Sg2M: "odewarłeś", Sg2F: "odewarłaś",
		Sg3M: "odewarł", Sg3F: "odewarła", Sg3N: "odewarło",
		Pl1V: "odewarliśmy", Pl1NV: "odewarłyśmy",
		Pl2V: "odewarliście", Pl2NV: "odewarłyście",
		Pl3V: "odewarli", Pl3NV: "odewarły",
	},

	// trzeć → tarł/tarła (e→a alternation)
	"trzeć": {
		Sg1M: "tarłem", Sg1F: "tarłam",
		Sg2M: "tarłeś", Sg2F: "tarłaś",
		Sg3M: "tarł", Sg3F: "tarła", Sg3N: "tarło",
		Pl1V: "tarliśmy", Pl1NV: "tarłyśmy",
		Pl2V: "tarliście", Pl2NV: "tarłyście",
		Pl3V: "tarli", Pl3NV: "tarły",
	},

	// drzeć → darł/darła (e→a alternation)
	"drzeć": {
		Sg1M: "darłem", Sg1F: "darłam",
		Sg2M: "darłeś", Sg2F: "darłaś",
		Sg3M: "darł", Sg3F: "darła", Sg3N: "darło",
		Pl1V: "darliśmy", Pl1NV: "darłyśmy",
		Pl2V: "darliście", Pl2NV: "darłyście",
		Pl3V: "darli", Pl3NV: "darły",
	},

	// mrzeć → marł/marła (e→a alternation)
	"mrzeć": {
		Sg1M: "marłem", Sg1F: "marłam",
		Sg2M: "marłeś", Sg2F: "marłaś",
		Sg3M: "marł", Sg3F: "marła", Sg3N: "marło",
		Pl1V: "marliśmy", Pl1NV: "marłyśmy",
		Pl2V: "marliście", Pl2NV: "marłyście",
		Pl3V: "marli", Pl3NV: "marły",
	},

	// żreć → żarł/żarła (e→a alternation)
	"żreć": {
		Sg1M: "żarłem", Sg1F: "żarłam",
		Sg2M: "żarłeś", Sg2F: "żarłaś",
		Sg3M: "żarł", Sg3F: "żarła", Sg3N: "żarło",
		Pl1V: "żarliśmy", Pl1NV: "żarłyśmy",
		Pl2V: "żarliście", Pl2NV: "żarłyście",
		Pl3V: "żarli", Pl3NV: "żarły",
	},

	// brać → brał (simple -ać, but included as base for prefixes)
	"brać": {
		Sg1M: "brałem", Sg1F: "brałam",
		Sg2M: "brałeś", Sg2F: "brałaś",
		Sg3M: "brał", Sg3F: "brała", Sg3N: "brało",
		Pl1V: "braliśmy", Pl1NV: "brałyśmy",
		Pl2V: "braliście", Pl2NV: "brałyście",
		Pl3V: "brali", Pl3NV: "brały",
	},

	// prać → prał
	"prać": {
		Sg1M: "prałem", Sg1F: "prałam",
		Sg2M: "prałeś", Sg2F: "prałaś",
		Sg3M: "prał", Sg3F: "prała", Sg3N: "prało",
		Pl1V: "praliśmy", Pl1NV: "prałyśmy",
		Pl2V: "praliście", Pl2NV: "prałyście",
		Pl3V: "prali", Pl3NV: "prały",
	},

	// -jść verbs (prefixed iść): przejść → przeszedł
	// These use the szedł stem with prefix modifications

	// dać → dał
	"dać": {
		Sg1M: "dałem", Sg1F: "dałam",
		Sg2M: "dałeś", Sg2F: "dałaś",
		Sg3M: "dał", Sg3F: "dała", Sg3N: "dało",
		Pl1V: "daliśmy", Pl1NV: "dałyśmy",
		Pl2V: "daliście", Pl2NV: "dałyście",
		Pl3V: "dali", Pl3NV: "dały",
	},

	// stać → stał (to stand)
	"stać": {
		Sg1M: "stałem", Sg1F: "stałam",
		Sg2M: "stałeś", Sg2F: "stałaś",
		Sg3M: "stał", Sg3F: "stała", Sg3N: "stało",
		Pl1V: "staliśmy", Pl1NV: "stałyśmy",
		Pl2V: "staliście", Pl2NV: "stałyście",
		Pl3V: "stali", Pl3NV: "stały",
	},

	// mieć → miał
	"mieć": {
		Sg1M: "miałem", Sg1F: "miałam",
		Sg2M: "miałeś", Sg2F: "miałaś",
		Sg3M: "miał", Sg3F: "miała", Sg3N: "miało",
		Pl1V: "mieliśmy", Pl1NV: "miałyśmy",
		Pl2V: "mieliście", Pl2NV: "miałyście",
		Pl3V: "mieli", Pl3NV: "miały",
	},

	// chcieć → chciał
	"chcieć": {
		Sg1M: "chciałem", Sg1F: "chciałam",
		Sg2M: "chciałeś", Sg2F: "chciałaś",
		Sg3M: "chciał", Sg3F: "chciała", Sg3N: "chciało",
		Pl1V: "chcieliśmy", Pl1NV: "chciałyśmy",
		Pl2V: "chcieliście", Pl2NV: "chciałyście",
		Pl3V: "chcieli", Pl3NV: "chciały",
	},

	// wiedzieć → wiedział
	"wiedzieć": {
		Sg1M: "wiedziałem", Sg1F: "wiedziałam",
		Sg2M: "wiedziałeś", Sg2F: "wiedziałaś",
		Sg3M: "wiedział", Sg3F: "wiedziała", Sg3N: "wiedziało",
		Pl1V: "wiedzieliśmy", Pl1NV: "wiedziałyśmy",
		Pl2V: "wiedzieliście", Pl2NV: "wiedziałyście",
		Pl3V: "wiedzieli", Pl3NV: "wiedziały",
	},

	// siedzieć → siedział
	"siedzieć": {
		Sg1M: "siedziałem", Sg1F: "siedziałam",
		Sg2M: "siedziałeś", Sg2F: "siedziałaś",
		Sg3M: "siedział", Sg3F: "siedziała", Sg3N: "siedziało",
		Pl1V: "siedzieliśmy", Pl1NV: "siedziałyśmy",
		Pl2V: "siedzieliście", Pl2NV: "siedziałyście",
		Pl3V: "siedzieli", Pl3NV: "siedziały",
	},

	// widzieć → widział
	"widzieć": {
		Sg1M: "widziałem", Sg1F: "widziałam",
		Sg2M: "widziałeś", Sg2F: "widziałaś",
		Sg3M: "widział", Sg3F: "widziała", Sg3N: "widziało",
		Pl1V: "widzieliśmy", Pl1NV: "widziałyśmy",
		Pl2V: "widzieliście", Pl2NV: "widziałyście",
		Pl3V: "widzieli", Pl3NV: "widziały",
	},

	// słyszeć → słyszał
	"słyszeć": {
		Sg1M: "słyszałem", Sg1F: "słyszałam",
		Sg2M: "słyszałeś", Sg2F: "słyszałaś",
		Sg3M: "słyszał", Sg3F: "słyszała", Sg3N: "słyszało",
		Pl1V: "słyszeliśmy", Pl1NV: "słyszałyśmy",
		Pl2V: "słyszeliście", Pl2NV: "słyszałyście",
		Pl3V: "słyszeli", Pl3NV: "słyszały",
	},

	// musieć → musiał
	"musieć": {
		Sg1M: "musiałem", Sg1F: "musiałam",
		Sg2M: "musiałeś", Sg2F: "musiałaś",
		Sg3M: "musiał", Sg3F: "musiała", Sg3N: "musiało",
		Pl1V: "musieliśmy", Pl1NV: "musiałyśmy",
		Pl2V: "musieliście", Pl2NV: "musiałyście",
		Pl3V: "musieli", Pl3NV: "musiały",
	},

	// móc → mógł/mogła (ó→o alternation)
	"móc": {
		Sg1M: "mogłem", Sg1F: "mogłam",
		Sg2M: "mogłeś", Sg2F: "mogłaś",
		Sg3M: "mógł", Sg3F: "mogła", Sg3N: "mogło",
		Pl1V: "mogliśmy", Pl1NV: "mogłyśmy",
		Pl2V: "mogliście", Pl2NV: "mogłyście",
		Pl3V: "mogli", Pl3NV: "mogły",
	},

	// biec → biegł (suppletive stem bieg-)
	"biec": {
		Sg1M: "biegłem", Sg1F: "biegłam",
		Sg2M: "biegłeś", Sg2F: "biegłaś",
		Sg3M: "biegł", Sg3F: "biegła", Sg3N: "biegło",
		Pl1V: "biegliśmy", Pl1NV: "biegłyśmy",
		Pl2V: "biegliście", Pl2NV: "biegłyście",
		Pl3V: "biegli", Pl3NV: "biegły",
	},

	// lec → legł
	"lec": {
		Sg1M: "ległem", Sg1F: "ległam",
		Sg2M: "ległeś", Sg2F: "ległaś",
		Sg3M: "legł", Sg3F: "legła", Sg3N: "legło",
		Pl1V: "legliśmy", Pl1NV: "ległyśmy",
		Pl2V: "legliście", Pl2NV: "ległyście",
		Pl3V: "legli", Pl3NV: "legły",
	},

	// rzec → rzekł
	"rzec": {
		Sg1M: "rzekłem", Sg1F: "rzekłam",
		Sg2M: "rzekłeś", Sg2F: "rzekłaś",
		Sg3M: "rzekł", Sg3F: "rzekła", Sg3N: "rzekło",
		Pl1V: "rzekliśmy", Pl1NV: "rzekłyśmy",
		Pl2V: "rzekliście", Pl2NV: "rzekłyście",
		Pl3V: "rzekli", Pl3NV: "rzekły",
	},

	// ciec → ciekł (suppletive stem ciek-)
	"ciec": {
		Sg1M: "ciekłem", Sg1F: "ciekłam",
		Sg2M: "ciekłeś", Sg2F: "ciekłaś",
		Sg3M: "ciekł", Sg3F: "ciekła", Sg3N: "ciekło",
		Pl1V: "ciekliśmy", Pl1NV: "ciekłyśmy",
		Pl2V: "ciekliście", Pl2NV: "ciekłyście",
		Pl3V: "ciekli", Pl3NV: "ciekły",
	},

	// strzec → strzegł
	"strzec": {
		Sg1M: "strzegłem", Sg1F: "strzegłam",
		Sg2M: "strzegłeś", Sg2F: "strzegłaś",
		Sg3M: "strzegł", Sg3F: "strzegła", Sg3N: "strzegło",
		Pl1V: "strzegliśmy", Pl1NV: "strzegłyśmy",
		Pl2V: "strzegliście", Pl2NV: "strzegłyście",
		Pl3V: "strzegli", Pl3NV: "strzegły",
	},

	// przesiąc → przesiąkł (archaic shortened form of przesiąknąć)
	"przesiąc": {
		Sg1M: "przesiąkłem", Sg1F: "przesiąkłam",
		Sg2M: "przesiąkłeś", Sg2F: "przesiąkłaś",
		Sg3M: "przesiąkł", Sg3F: "przesiąkła", Sg3N: "przesiąkło",
		Pl1V: "przesiąkliśmy", Pl1NV: "przesiąkłyśmy",
		Pl2V: "przesiąkliście", Pl2NV: "przesiąkłyście",
		Pl3V: "przesiąkli", Pl3NV: "przesiąkły",
	},

	// schnąć → sechł/schła (asymmetric: masc sg has 'e', others don't)
	"schnąć": {
		Sg1M: "sechłem", Sg1F: "schłam",
		Sg2M: "sechłeś", Sg2F: "schłaś",
		Sg3M: "sechł", Sg3F: "schła", Sg3N: "schło",
		Pl1V: "schliśmy", Pl1NV: "schłyśmy",
		Pl2V: "schliście", Pl2NV: "schłyście",
		Pl3V: "schli", Pl3NV: "schły",
	},

	// przysięgnąć → przysiągł/przysięgła (to swear an oath)
	// This is a separate lexeme from sięgnąć (to reach), which keeps n.
	// The przysiąg- stem has ę→ą alternation in masculine forms.
	"przysięgnąć": {
		Sg1M: "przysiągłem", Sg1F: "przysięgłam",
		Sg2M: "przysiągłeś", Sg2F: "przysięgłaś",
		Sg3M: "przysiągł", Sg3F: "przysięgła", Sg3N: "przysięgło",
		Pl1V: "przysięgliśmy", Pl1NV: "przysięgłyśmy",
		Pl2V: "przysięgliście", Pl2NV: "przysięgłyście",
		Pl3V: "przysięgli", Pl3NV: "przysięgły",
	},

	// piec → piekł
	"piec": {
		Sg1M: "piekłem", Sg1F: "piekłam",
		Sg2M: "piekłeś", Sg2F: "piekłaś",
		Sg3M: "piekł", Sg3F: "piekła", Sg3N: "piekło",
		Pl1V: "piekliśmy", Pl1NV: "piekłyśmy",
		Pl2V: "piekliście", Pl2NV: "piekłyście",
		Pl3V: "piekli", Pl3NV: "piekły",
	},

	// wlec → wlókł/wlokła (ó→o alternation)
	"wlec": {
		Sg1M: "wlokłem", Sg1F: "wlokłam",
		Sg2M: "wlokłeś", Sg2F: "wlokłaś",
		Sg3M: "wlókł", Sg3F: "wlokła", Sg3N: "wlokło",
		Pl1V: "wlekliśmy", Pl1NV: "wlokłyśmy",
		Pl2V: "wlekliście", Pl2NV: "wlokłyście",
		Pl3V: "wlekli", Pl3NV: "wlokły",
	},

	// siąść family need special handling
	// usiąść → usiadł
	"usiąść": {
		Sg1M: "usiadłem", Sg1F: "usiadłam",
		Sg2M: "usiadłeś", Sg2F: "usiadłaś",
		Sg3M: "usiadł", Sg3F: "usiadła", Sg3N: "usiadło",
		Pl1V: "usiedliśmy", Pl1NV: "usiadłyśmy",
		Pl2V: "usiedliście", Pl2NV: "usiadłyście",
		Pl3V: "usiedli", Pl3NV: "usiadły",
	},

	// -sieść family (prefixed forms of sieść = to sit down)
	// wysieść → wysiadł (wy + sieść)
	"wysieść": {
		Sg1M: "wysiadłem", Sg1F: "wysiadłam",
		Sg2M: "wysiadłeś", Sg2F: "wysiadłaś",
		Sg3M: "wysiadł", Sg3F: "wysiadła", Sg3N: "wysiadło",
		Pl1V: "wysiedliśmy", Pl1NV: "wysiadłyśmy",
		Pl2V: "wysiedliście", Pl2NV: "wysiadłyście",
		Pl3V: "wysiedli", Pl3NV: "wysiadły",
	},
	// zsieść → zsiadł (z + sieść)
	"zsieść": {
		Sg1M: "zsiadłem", Sg1F: "zsiadłam",
		Sg2M: "zsiadłeś", Sg2F: "zsiadłaś",
		Sg3M: "zsiadł", Sg3F: "zsiadła", Sg3N: "zsiadło",
		Pl1V: "zsiedliśmy", Pl1NV: "zsiadłyśmy",
		Pl2V: "zsiedliście", Pl2NV: "zsiadłyście",
		Pl3V: "zsiedli", Pl3NV: "zsiadły",
	},

	// rosnąć → rósł/rosła (special n-dropping with ó→o)
	"rosnąć": {
		Sg1M: "rosłem", Sg1F: "rosłam",
		Sg2M: "rosłeś", Sg2F: "rosłaś",
		Sg3M: "rósł", Sg3F: "rosła", Sg3N: "rosło",
		Pl1V: "rośliśmy", Pl1NV: "rosłyśmy",
		Pl2V: "rośliście", Pl2NV: "rosłyście",
		Pl3V: "rośli", Pl3NV: "rosły",
	},

	// rość - shortened form of rosnąć
	"rość": {
		Sg1M: "rosłem", Sg1F: "rosłam",
		Sg2M: "rosłeś", Sg2F: "rosłaś",
		Sg3M: "rósł", Sg3F: "rosła", Sg3N: "rosło",
		Pl1V: "rośliśmy", Pl1NV: "rosłyśmy",
		Pl2V: "rośliście", Pl2NV: "rosłyście",
		Pl3V: "rośli", Pl3NV: "rosły",
	},
}

// pastPrefixableVerbs lists verbs that can take prefixes in past tense.
var pastPrefixableVerbs = map[string]bool{
	"być": true, "iść": true, "jeść": true, "brać": true, "prać": true,
	"jąć": true, "dąć": true, "ciąć": true, "giąć": true, "piąć": true, "miąć": true, "nająć": true,
	"żąć": true, "kląć": true, "wziąć": true,
	"siąść": true, "paść": true, "kraść": true, "kłaść": true, "prząść": true,
	"gryźć": true, "leźć": true, "wieźć": true, "nieść": true,
	"pleść": true, "grześć": true, "tłuc": true,
	"przeć": true, "wrzeć": true, "trzeć": true, "drzeć": true, "mrzeć": true, "żreć": true,
	"dać": true, "stać": true, "mieć": true,
	"wiedzieć": true, "siedzieć": true, "widzieć": true,
	"biec": true, "lec": true, "rzec": true, "ciec": true, "strzec": true, "piec": true, "wlec": true,
	"rosnąć": true, "rość": true, "schnąć": true, "przysięgnąć": true,
}

// lookupPastIrregular checks if a verb has an irregular past tense paradigm.
func lookupPastIrregular(infinitive string) (PastTense, bool) {
	p, ok := irregularPastVerbs[infinitive]
	return p, ok
}

// lookupPastIrregularWithPrefix tries to find an irregular past tense verb,
// including checking if it's a prefixed form of a known irregular.
func lookupPastIrregularWithPrefix(infinitive string) (PastTense, bool) {
	// Direct lookup first
	if p, ok := irregularPastVerbs[infinitive]; ok {
		return p, ok
	}

	// Handle -nijść verbs (archaic variant): wnijść → wszedł (must come before -jść)
	if strings.HasSuffix(infinitive, "nijść") {
		prefix := strings.TrimSuffix(infinitive, "nijść")
		if prefix != "" {
			return buildJscPast(prefix), true
		}
	}

	// Handle -jść verbs (prefixed iść): przejść → przeszedł
	if strings.HasSuffix(infinitive, "jść") {
		prefix := strings.TrimSuffix(infinitive, "jść")
		if prefix != "" {
			// przejść → przeszedł, wyjść → wyszedł, etc.
			return buildJscPast(prefix), true
		}
	}

	// Handle -niść verbs (archaic/dialectal variants of -jść): wniść → wszedł
	if strings.HasSuffix(infinitive, "niść") {
		prefix := strings.TrimSuffix(infinitive, "niść")
		if prefix != "" {
			// wniść → wszedł, wyniść → wyszedł, zniść → zszedł, etc.
			return buildJscPast(prefix), true
		}
	}

	// Handle prefixed -schnąć verbs: obeschnąć → obsechł/obeschła
	// These have asymmetric stems (masc sg vs others) and complex epenthetic handling.
	if strings.HasSuffix(infinitive, "schnąć") && infinitive != "schnąć" {
		return buildSchnacPast(infinitive), true
	}

	// Try stripping prefixes to find base irregular verb
	for _, prefix := range verbPrefixes {
		if len(infinitive) > len(prefix) && infinitive[:len(prefix)] == prefix {
			base := infinitive[len(prefix):]
			if pastPrefixableVerbs[base] {
				if baseParadigm, ok := irregularPastVerbs[base]; ok {
					// Apply prefix to all forms
					return applyPrefixToPast(prefix, baseParadigm), true
				}
			}
		}
	}

	return PastTense{}, false
}

// buildJscPast builds past tense for -jść verbs (prefixed iść).
// przejść → przeszedł/przeszła, wyjść → wyszedł/wyszła
func buildJscPast(prefix string) PastTense {
	return PastTense{
		Sg1M:  prefix + "szedłem",
		Sg1F:  prefix + "szłam",
		Sg2M:  prefix + "szedłeś",
		Sg2F:  prefix + "szłaś",
		Sg3M:  prefix + "szedł",
		Sg3F:  prefix + "szła",
		Sg3N:  prefix + "szło",
		Pl1V:  prefix + "szliśmy",
		Pl1NV: prefix + "szłyśmy",
		Pl2V:  prefix + "szliście",
		Pl2NV: prefix + "szłyście",
		Pl3V:  prefix + "szli",
		Pl3NV: prefix + "szły",
	}
}

// buildSchnacPast builds past tense for prefixed schnąć verbs.
// schnąć has an asymmetric stem: masc sg uses "sechł", others use "schł".
// Additionally, prefixes with epenthetic vowels (obe-, pode-, roze-, ze-)
// strip the vowel in masc sg but keep it in other forms.
// obeschnąć: sg3m=obsechł (ob+sechł), sg3f=obeschła (obe+schła)
func buildSchnacPast(infinitive string) PastTense {
	// Determine the prefix from the infinitive
	prefix := strings.TrimSuffix(infinitive, "schnąć")

	// For masculine singular: strip epenthetic from prefix, use sechł stem
	mascPrefix := prefix
	epenthetic := map[string]string{
		"obe": "ob", "pode": "pod", "roze": "roz", "ze": "z",
	}
	for full, stripped := range epenthetic {
		if strings.HasSuffix(prefix, full) {
			mascPrefix = prefix[:len(prefix)-len(full)] + stripped
			break
		}
	}

	// For other forms: keep original prefix, use schł stem
	otherPrefix := prefix

	// Special case: zeschnąć has ze→s in masc (ssechł, not zsechł),
	// but ze→ze in others (zeschła, not sschła)
	// Actually the corpus shows ze→z for masc: zeschnąć → ssechł = z+sechł
	// Wait, that's s+sechł = ssechł. So ze→s? No, z+sechł = zsechł, but corpus shows ssechł.
	// Looking closer: zeschnąć: sg3m=ssechł. This must be a special case.
	// Actually, s+sechł wouldn't make sense either. Let me handle zeschnąć specially.
	if infinitive == "zeschnąć" {
		// zeschnąć is weird: sg3m = ssechł (ze → s before s-stem?)
		// Actually, historically ze- before s- assimilates: ze+schnąć → zeschnąć,
		// past: ze+sechł → the ze is dropped and s+sechł = ssechł
		return PastTense{
			Sg1M: "ssechłem", Sg1F: "zeschłam",
			Sg2M: "ssechłeś", Sg2F: "zeschłaś",
			Sg3M: "ssechł", Sg3F: "zeschła", Sg3N: "zeschło",
			Pl1V: "zeschliśmy", Pl1NV: "zeschłyśmy",
			Pl2V: "zeschliście", Pl2NV: "zeschłyście",
			Pl3V: "zeschli", Pl3NV: "zeschły",
		}
	}

	return PastTense{
		Sg1M:  mascPrefix + "sechłem",
		Sg1F:  otherPrefix + "schłam",
		Sg2M:  mascPrefix + "sechłeś",
		Sg2F:  otherPrefix + "schłaś",
		Sg3M:  mascPrefix + "sechł",
		Sg3F:  otherPrefix + "schła",
		Sg3N:  otherPrefix + "schło",
		Pl1V:  otherPrefix + "schliśmy",
		Pl1NV: otherPrefix + "schłyśmy",
		Pl2V:  otherPrefix + "schliście",
		Pl2NV: otherPrefix + "schłyście",
		Pl3V:  otherPrefix + "schli",
		Pl3NV: otherPrefix + "schły",
	}
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
	// The epenthetic vowel is kept before: s, ś, z, ź, ż, b, p, m, w
	if prefix == "ze" {
		keepVowel := map[rune]bool{
			's': true, 'ś': true, 'z': true, 'ź': true, 'ż': true,
			'b': true, 'p': true, 'm': true, 'w': true,
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
