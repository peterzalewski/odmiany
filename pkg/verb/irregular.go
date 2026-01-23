package verb

// irregularVerbs contains present tense paradigms for verbs that cannot
// be conjugated by heuristics alone. These are either:
// - Suppletive verbs (stem changes completely: być → jestem)
// - Minority pattern verbs (e.g., pisać → piszę when most -sać → -sam)
var irregularVerbs = map[string]PresentTense{
	// Suppletive verbs - completely irregular stems
	"być": {
		Sg1: "jestem", Sg2: "jesteś", Sg3: "jest",
		Pl1: "jesteśmy", Pl2: "jesteście", Pl3: "są",
	},
	"mieć": {
		Sg1: "mam", Sg2: "masz", Sg3: "ma",
		Pl1: "mamy", Pl2: "macie", Pl3: "mają",
	},
	"jeść": {
		Sg1: "jem", Sg2: "jesz", Sg3: "je",
		Pl1: "jemy", Pl2: "jecie", Pl3: "jedzą",
	},
	"dać": {
		Sg1: "dam", Sg2: "dasz", Sg3: "da",
		Pl1: "damy", Pl2: "dacie", Pl3: "dadzą",
	},
	"wziąć": {
		Sg1: "wezmę", Sg2: "weźmiesz", Sg3: "weźmie",
		Pl1: "weźmiemy", Pl2: "weźmiecie", Pl3: "wezmą",
	},
	"ciąć": {
		Sg1: "tnę", Sg2: "tniesz", Sg3: "tnie",
		Pl1: "tniemy", Pl2: "tniecie", Pl3: "tną",
	},
	"iść": {
		Sg1: "idę", Sg2: "idziesz", Sg3: "idzie",
		Pl1: "idziemy", Pl2: "idziecie", Pl3: "idą",
	},

	// Stem-changing -ać verbs (jechać type: ech→ad/edz)
	"jechać": {
		Sg1: "jadę", Sg2: "jedziesz", Sg3: "jedzie",
		Pl1: "jedziemy", Pl2: "jedziecie", Pl3: "jadą",
	},

	// Stem-changing -ać verbs (brać type: a→io/ie)
	"brać": {
		Sg1: "biorę", Sg2: "bierzesz", Sg3: "bierze",
		Pl1: "bierzemy", Pl2: "bierzecie", Pl3: "biorą",
	},
	"prać": {
		Sg1: "piorę", Sg2: "pierzesz", Sg3: "pierze",
		Pl1: "pierzemy", Pl2: "pierzecie", Pl3: "piorą",
	},

	// Minority -sać verbs that alternate (s→sz)
	// Most -sać verbs are regular (-sam), but these go to -szę
	"pisać": {
		Sg1: "piszę", Sg2: "piszesz", Sg3: "pisze",
		Pl1: "piszemy", Pl2: "piszecie", Pl3: "piszą",
	},
	"czesać": {
		Sg1: "czeszę", Sg2: "czeszesz", Sg3: "czesze",
		Pl1: "czeszemy", Pl2: "czeszecie", Pl3: "czeszą",
	},
	"kasać": {
		Sg1: "kaszę", Sg2: "kaszesz", Sg3: "kasze",
		Pl1: "kaszemy", Pl2: "kaszecie", Pl3: "kaszą",
	},
	"kołysać": {
		Sg1: "kołyszę", Sg2: "kołyszesz", Sg3: "kołysze",
		Pl1: "kołyszemy", Pl2: "kołyszecie", Pl3: "kołyszą",
	},
	"ciosać": {
		Sg1: "cioszę", Sg2: "cioszesz", Sg3: "ciosze",
		Pl1: "cioszemy", Pl2: "cioszecie", Pl3: "cioszą",
	},
	"ciesać": {
		Sg1: "cieszę", Sg2: "cieszesz", Sg3: "ciesze",
		Pl1: "cieszemy", Pl2: "cieszecie", Pl3: "cieszą",
	},
	"krzesać": {
		Sg1: "krzeszę", Sg2: "krzeszesz", Sg3: "krzesze",
		Pl1: "krzeszemy", Pl2: "krzeszecie", Pl3: "krzeszą",
	},
	"skakać": {
		Sg1: "skaczę", Sg2: "skaczesz", Sg3: "skacze",
		Pl1: "skaczemy", Pl2: "skaczecie", Pl3: "skaczą",
	},
	"płakać": {
		Sg1: "płaczę", Sg2: "płaczesz", Sg3: "płacze",
		Pl1: "płaczemy", Pl2: "płaczecie", Pl3: "płaczą",
	},

	// Minority -zać verbs that alternate (z→ż)
	"wiązać": {
		Sg1: "wiążę", Sg2: "wiążesz", Sg3: "wiąże",
		Pl1: "wiążemy", Pl2: "wiążecie", Pl3: "wiążą",
	},
	"kazać": {
		Sg1: "każę", Sg2: "każesz", Sg3: "każe",
		Pl1: "każemy", Pl2: "każecie", Pl3: "każą",
	},
	"mazać": {
		Sg1: "mażę", Sg2: "mażesz", Sg3: "maże",
		Pl1: "mażemy", Pl2: "mażecie", Pl3: "mażą",
	},
	"lizać": {
		Sg1: "liżę", Sg2: "liżesz", Sg3: "liże",
		Pl1: "liżemy", Pl2: "liżecie", Pl3: "liżą",
	},

	// stać - has two meanings with different conjugations
	// stać (to stand) → stoję, stoisz... (imperfective)
	"stać": {
		Sg1: "stoję", Sg2: "stoisz", Sg3: "stoi",
		Pl1: "stoimy", Pl2: "stoicie", Pl3: "stoją",
	},
}

// lookupIrregular checks if a verb has an irregular paradigm.
// Returns the paradigm and true if found, zero value and false otherwise.
func lookupIrregular(infinitive string) (PresentTense, bool) {
	p, ok := irregularVerbs[infinitive]
	return p, ok
}

// For prefixed verbs, we derive from the base form.
// This allows "napisać" to use "pisać" paradigm with prefix.
var irregularBases = map[string]string{
	// -pisać derivatives
	"pisać": "pisać",
	// -brać derivatives
	"brać": "brać",
	// -jechać derivatives
	"jechać": "jechać",
	// -dać derivatives
	"dać": "dać",
	// -wziąć derivatives
	"wziąć": "wziąć",
	// -iść derivatives
	"iść": "iść",
	// -jeść derivatives
	"jeść": "jeść",
}

// Common prefixes in Polish
var verbPrefixes = []string{
	"prze", "przy", "roz", "roze", "wy", "za", "na", "po", "do", "od", "ode", "ob", "obe",
	"pod", "pode", "nad", "nade", "wz", "wze", "u", "s", "z", "ze", "w", "we",
}

// lookupIrregularWithPrefix tries to find an irregular verb,
// including checking if it's a prefixed form of a known irregular.
func lookupIrregularWithPrefix(infinitive string) (PresentTense, bool) {
	// Direct lookup first
	if p, ok := irregularVerbs[infinitive]; ok {
		return p, ok
	}

	// Try stripping prefixes to find base irregular verb
	// Only for verbs that are known to take prefixes productively
	prefixableVerbs := map[string]bool{
		"pisać": true, "brać": true, "jechać": true, "dać": true,
		"wziąć": true, "iść": true, "jeść": true, "prać": true,
		"czesać": true, "kasać": true, "ciosać": true, "ciesać": true,
		"skakać": true, "płakać": true, "wiązać": true, "kazać": true,
		"mazać": true, "lizać": true, "kołysać": true, "krzesać": true,
	}

	for _, prefix := range verbPrefixes {
		if len(infinitive) > len(prefix) && infinitive[:len(prefix)] == prefix {
			base := infinitive[len(prefix):]
			if prefixableVerbs[base] {
				if baseParadigm, ok := irregularVerbs[base]; ok {
					// Apply prefix to all forms
					return PresentTense{
						Sg1: prefix + baseParadigm.Sg1,
						Sg2: prefix + baseParadigm.Sg2,
						Sg3: prefix + baseParadigm.Sg3,
						Pl1: prefix + baseParadigm.Pl1,
						Pl2: prefix + baseParadigm.Pl2,
						Pl3: prefix + baseParadigm.Pl3,
					}, true
				}
			}
		}
	}

	return PresentTense{}, false
}
