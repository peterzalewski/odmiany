package verb

import (
	"fmt"
	"strings"
)

// ConjugatePast returns all valid past tense paradigms for a verb.
// Most verbs return a single paradigm; homographs and dual-form verbs return multiple.
func ConjugatePast(infinitive string) ([]PastParadigm, error) {
	// Check homographs first (verbs with multiple valid paradigms)
	if paradigms, ok := lookupPastHomograph(infinitive); ok {
		return paradigms, nil
	}

	// Check irregular verbs (including prefixed forms)
	if p, ok := lookupPastIrregularWithPrefix(infinitive); ok {
		return []PastParadigm{{PastTense: p}}, nil
	}

	// Check for dual-form -nąć verbs (both n-dropping and n-keeping valid)
	if isDualFormNacVerb(infinitive) {
		return buildDualFormNacParadigms(infinitive), nil
	}

	// Try heuristics in order of specificity
	for _, h := range pastHeuristics {
		if p, ok := h(infinitive); ok {
			return []PastParadigm{{PastTense: p}}, nil
		}
	}
	return nil, fmt.Errorf("no past tense heuristic matched: %s", infinitive)
}

// buildDualFormNacParadigms returns both paradigms for verbs that can use
// either n-dropped or n-kept sg3m in standard Polish.
//
// These verbs have a HYBRID paradigm where:
//   - sg1m, sg2m: n-keeping with ą (kwitnąłem, kwitnąłeś)
//   - sg1f, sg2f, sg3f, sg3n: n-keeping with ę (kwitnęłam, kwitnęła, etc.)
//   - sg3m: VARIES - either n-dropped (kwitł) or n-kept (kwitnął)
//   - pl*nv (non-virile): n-keeping with ę (kwitnęłyśmy, kwitnęły)
//
// The virile plural depends on the verb type:
//   - VirileDropped: n-dropped (kwitliśmy, kwitli)
//   - VirileKept: n-kept with ę (trzasnęliśmy, trzasnęli)
//
// The two returned paradigms differ only in sg3m.
func buildDualFormNacParadigms(infinitive string) []PastParadigm {
	stemWithoutNac := strings.TrimSuffix(infinitive, "nąć") // "kwit" or "klęk"
	baseStem := strings.TrimSuffix(infinitive, "ąć")        // "kwitn" or "klękn"
	mascNKeptStem := baseStem + "ą"                         // "kwitną" or "klękną"
	femStem := baseStem + "ę"                               // "kwitnę" or "klęknę"

	// Determine virile plural type
	var virileStem string
	if isDualFormVirileKept(infinitive) {
		// Virile uses n-kept stem: trzasnęli
		virileStem = femStem
	} else {
		// Virile uses n-dropped stem: kwitli, klękli
		virileStem = palatalizeForVirile(stemWithoutNac, infinitive)
	}

	// Apply masculine vowel alternation (ę→ą for klęknąć→kląkł)
	mascStemDropped := applyMascSgAlternation(stemWithoutNac, infinitive)

	// Check if this verb has masculine alternation (determines sg1m/sg2m pattern)
	// Verbs with alternation use n-dropped stem for sg1m/sg2m (kląkłem)
	// Verbs without alternation use n-kept stem for sg1m/sg2m (kwitnąłem)
	hasAlternation := mascStemDropped != stemWithoutNac

	var sg1m, sg2m string
	if hasAlternation {
		// klęknąć: sg1m=kląkłem (n-dropped with alternation)
		sg1m = mascStemDropped + "łem"
		sg2m = mascStemDropped + "łeś"
	} else {
		// kwitnąć: sg1m=kwitnąłem (n-kept)
		sg1m = mascNKeptStem + "łem"
		sg2m = mascNKeptStem + "łeś"
	}

	// Masculine sg3m n-dropped
	mascSg3MDropped := mascStemDropped + "ł"

	// Build the hybrid paradigm with n-dropped sg3m
	paradigm1 := PastTense{
		Sg1M:  sg1m,
		Sg1F:  femStem + "łam",
		Sg2M:  sg2m,
		Sg2F:  femStem + "łaś",
		Sg3M:  mascSg3MDropped, // n-dropped: kwitł or kląkł
		Sg3F:  femStem + "ła",
		Sg3N:  femStem + "ło",
		Pl1V:  virileStem + "liśmy",
		Pl1NV: femStem + "łyśmy",
		Pl2V:  virileStem + "liście",
		Pl2NV: femStem + "łyście",
		Pl3V:  virileStem + "li",
		Pl3NV: femStem + "ły",
	}

	// Build the hybrid paradigm with n-kept sg3m
	paradigm2 := PastTense{
		Sg1M:  sg1m, // same as paradigm1
		Sg1F:  femStem + "łam",
		Sg2M:  sg2m, // same as paradigm1
		Sg2F:  femStem + "łaś",
		Sg3M:  mascNKeptStem + "ł", // n-kept: kwitnął or klęknął
		Sg3F:  femStem + "ła",
		Sg3N:  femStem + "ło",
		Pl1V:  virileStem + "liśmy",
		Pl1NV: femStem + "łyśmy",
		Pl2V:  virileStem + "liście",
		Pl2NV: femStem + "łyście",
		Pl3V:  virileStem + "li",
		Pl3NV: femStem + "ły",
	}

	return []PastParadigm{
		{PastTense: paradigm1, Gloss: "sg3m n-dropped variant"},
		{PastTense: paradigm2, Gloss: "sg3m n-kept variant"},
	}
}

// pastHeuristic is a function that attempts to conjugate a verb in past tense.
type pastHeuristic func(infinitive string) (PastTense, bool)

// pastHeuristics is the ordered list of past tense conjugation heuristics.
var pastHeuristics = []pastHeuristic{
	// -ąść/-ąźć verbs: trząść → trząsł, prząść → prządł
	heuristicPastAsc,
	// -cząć verbs: począć → począł/poczęła
	heuristicPastCzac,
	// -strzyc verbs: strzyc → strzygł/strzygła
	heuristicPastStrzyc,
	// -bość/-bóść verbs: bość → bódł/bodła
	heuristicPastBosc,
	// -nąć verbs: two patterns based on stem
	heuristicPastNac,
	// -ść/-źć verbs: nieść → niósł
	heuristicPastSc,
	// -c verbs (móc, piec): móc → mógł
	heuristicPastC,
	// -ować verbs: pracować → pracował
	heuristicPastOwac,
	// -ywać/-iwać verbs: pokazywać → pokazywał
	heuristicPastYwacIwac,
	// -awać verbs: dawać → dawał
	heuristicPastAwac,
	// -ić verbs: robić → robił
	heuristicPastIc,
	// -yć verbs: myć → mył
	heuristicPastYc,
	// -uć verbs: czuć → czuł
	heuristicPastUc,
	// -eć verbs: umieć → umiał
	heuristicPastEc,
	// -ać verbs (fallback): czytać → czytał
	heuristicPastAc,
}

// buildPastTense creates a full past paradigm from the l-participle stem.
// The stem is what you get after removing -ł from the masculine singular 3rd person form.
// e.g., for czytać: stem = "czyta" → czytał, czytała, czytało, czytali, czytały
func buildPastTense(stem string) PastTense {
	// Masculine-personal plural uses -li (ł→l before i)
	// The virile forms soften ł to l
	virileStem := stem // stem already ends without ł, we add -li directly

	return PastTense{
		Sg1M:  stem + "łem",
		Sg1F:  stem + "łam",
		Sg2M:  stem + "łeś",
		Sg2F:  stem + "łaś",
		Sg3M:  stem + "ł",
		Sg3F:  stem + "ła",
		Sg3N:  stem + "ło",
		Pl1V:  virileStem + "liśmy",
		Pl1NV: stem + "łyśmy",
		Pl2V:  virileStem + "liście",
		Pl2NV: stem + "łyście",
		Pl3V:  virileStem + "li",
		Pl3NV: stem + "ły",
	}
}

// Verbs where e→a alternation applies to ALL forms (not just masculine).
// blednąć → bladłem, bladł, bladła (all use "blad" stem)
var allFormsEToAVerbs = map[string]bool{
	"blednąć": true, "bladnąć": true,
}

// buildPastTenseNDropped creates a past paradigm for n-dropping -nąć verbs.
// The stem is what remains after removing -nąć (e.g., "gas" from "gasnąć").
// These verbs need consonant palatalization in virile plural: s→ś, z→ź, etc.
// The infinitive is passed to check for vowel alternation patterns.
func buildPastTenseNDropped(stem, infinitive string) PastTense {
	// Check if this verb uses e→a alternation in ALL forms
	base := extractBase(infinitive)
	useAltForAll := allFormsEToAVerbs[infinitive] || (base != infinitive && allFormsEToAVerbs[base])

	// Get the virile stem with palatalized final consonant
	virileStem := palatalizeForVirile(stem, infinitive)

	// Check for masculine sg vowel alternation (ę→ą, e→a)
	// This applies to sg1m, sg2m, sg3m
	mascStem := applyMascSgAlternation(stem, infinitive)

	// Apply sg3m-only alternations (o→ó, epenthetic e)
	// These only affect sg3m, not sg1m/sg2m
	sg3mStem := applySg3MOnlyAlternation(mascStem, infinitive)

	// Determine feminine/neuter/non-virile stem
	var otherStem string
	if useAltForAll {
		// e→a applies to ALL forms: blednąć → bladła (not bledła)
		otherStem = mascStem
	} else {
		// Alternation only in masculine: klęknąć → kląkł but klękła
		otherStem = stem
	}

	return PastTense{
		Sg1M:  mascStem + "łem",
		Sg1F:  otherStem + "łam",
		Sg2M:  mascStem + "łeś",
		Sg2F:  otherStem + "łaś",
		Sg3M:  sg3mStem + "ł",
		Sg3F:  otherStem + "ła",
		Sg3N:  otherStem + "ło",
		Pl1V:  virileStem + "liśmy",
		Pl1NV: otherStem + "łyśmy",
		Pl2V:  virileStem + "liście",
		Pl2NV: otherStem + "łyście",
		Pl3V:  virileStem + "li",
		Pl3NV: otherStem + "ły",
	}
}

// applyMascSgAlternation applies vowel alternation for masculine singular forms (sg1m, sg2m, sg3m).
// This applies to alternations that affect ALL masculine forms:
// ę→ą or e→a: blednąć→bladł, więdnąć→wiądł, klęknąć→kląkł, etc.
func applyMascSgAlternation(stem, infinitive string) string {
	// Verbs with ę→ą or e→a alternation in ALL masculine sg forms
	eToAVerbs := map[string]bool{
		"blednąć": true, "bladnąć": true,
		"więdnąć": true, "zwiędnąć": true,
		"ziębnąć": true,
		"klęknąć": true, "klęsnąć": true,
		"lęgnąć": true, "lęknąć": true,
		"grzęznąć": true, "gręznąć": true, "grząznąć": true, "grąznąć": true,
		"przęgnąć": true, "strzęgnąć": true, "sięgnąć": true,
		"więznąć": true, "więzgnąć": true,
		"wiąznąć": true,
	}

	// Check if the infinitive itself is in the list first
	if eToAVerbs[infinitive] {
		return applyEToA(stem)
	}

	// Then check for prefixed forms (e.g., nadwiędnąć)
	base := extractBase(infinitive)
	if base != infinitive && eToAVerbs[base] {
		return applyEToA(stem)
	}

	return stem
}

// applyEToA applies ę→ą or e→a alternation to the rightmost occurrence.
func applyEToA(stem string) string {
	runes := []rune(stem)
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] == 'ę' {
			runes[i] = 'ą'
			return string(runes)
		}
		if runes[i] == 'e' {
			runes[i] = 'a'
			return string(runes)
		}
	}
	return stem
}

// applySg3MOnlyAlternation applies alternations that ONLY affect sg3m (not sg1m/sg2m).
// o→ó: moknąć → mókł (sg3m) but mokłem (sg1m)
// epenthetic e: schnąć → sechł (sg3m) but schłem (sg1m)
func applySg3MOnlyAlternation(stem, infinitive string) string {
	// Verbs with o→ó alternation ONLY in sg3m
	oToOKreskaVerbs := map[string]bool{
		"moknąć":   true,
		"chłodnąć": true,
	}

	// Verbs with epenthetic e ONLY in sg3m (consonant cluster before -ł)
	// schnąć → sechł (sg3m), schłem (sg1m)
	epentheticEVerbs := map[string]bool{
		"schnąć": true,
	}

	// Check o→ó (only in sg3m)
	if oToOKreskaVerbs[infinitive] || (extractBase(infinitive) != infinitive && oToOKreskaVerbs[extractBase(infinitive)]) {
		runes := []rune(stem)
		for i := len(runes) - 1; i >= 0; i-- {
			if runes[i] == 'o' {
				runes[i] = 'ó'
				return string(runes)
			}
		}
	}

	// Check epenthetic e (only in sg3m)
	if epentheticEVerbs[infinitive] || (extractBase(infinitive) != infinitive && epentheticEVerbs[extractBase(infinitive)]) {
		// Insert 'e' before the final consonant cluster
		// sch → sech
		runes := []rune(stem)
		if len(runes) >= 2 {
			// Find the consonant cluster and insert e before the last consonant
			// For "sch": insert e before 'h' → "sech"
			return string(runes[:len(runes)-1]) + "e" + string(runes[len(runes)-1:])
		}
	}

	return stem
}

// extractBase extracts the base verb from a potentially prefixed infinitive.
// Returns the infinitive itself if no valid prefix is found.
func extractBase(infinitive string) string {
	// Try longer prefixes first to avoid false positives
	// (e.g., "przy" before "prz", "prze" before "pr")
	sortedPrefixes := []string{
		"prze", "przy", "roze", "nade", "pode", "pode",
		"roz", "nad", "pod", "obe", "ode",
		"za", "na", "po", "do", "od", "ob", "wy",
		"u", "s", "z", "w", "o",
	}
	for _, prefix := range sortedPrefixes {
		if strings.HasPrefix(infinitive, prefix) {
			candidate := infinitive[len(prefix):]
			// Only accept if the candidate is at least 4 characters (minimum: "nąć" + 1 char stem)
			// This prevents false positives like więdnąć → iędnąć (from stripping "w")
			if len(candidate) >= 4 && strings.HasSuffix(candidate, "nąć") {
				return candidate
			}
		}
	}
	return infinitive
}

// palatalizeForVirile applies consonant palatalization for virile plural forms.
// This affects the final consonant before -li endings.
// s→ś, n→ń always. z→ź only when the stem contains ę (grzęźli but not marźli).
func palatalizeForVirile(stem, infinitive string) string {
	if stem == "" {
		return stem
	}
	runes := []rune(stem)
	last := runes[len(runes)-1]

	// Always palatalize s→ś and n→ń
	if last == 's' {
		runes[len(runes)-1] = 'ś'
		return string(runes)
	}
	if last == 'n' {
		runes[len(runes)-1] = 'ń'
		return string(runes)
	}

	// Only palatalize z→ź if the stem contains ę
	// grzęznąć → grzęźli (has ę)
	// marznąć → marzli (no ę, no palatalization)
	if last == 'z' && strings.ContainsRune(infinitive, 'ę') {
		runes[len(runes)-1] = 'ź'
		return string(runes)
	}

	return stem
}

// buildPastTenseMixedNDrop creates a past paradigm for verbs with mixed n-dropping.
// Pattern: singular/non-virile retain n (with ą→ę), virile plural drops n.
// Example: cuchnąć → cuchnął/cuchnęła (sg retain), cuchli (pl virile drop)
func buildPastTenseMixedNDrop(stemWithoutNac, infinitive string) PastTense {
	// Virile stem drops the n
	virileStem := palatalizeForVirile(stemWithoutNac, infinitive)

	// Other forms retain n with ą→ę alternation
	mascStem := stemWithoutNac + "ną"
	femStem := stemWithoutNac + "nę"

	return PastTense{
		Sg1M:  mascStem + "łem",
		Sg1F:  femStem + "łam",
		Sg2M:  mascStem + "łeś",
		Sg2F:  femStem + "łaś",
		Sg3M:  mascStem + "ł",
		Sg3F:  femStem + "ła",
		Sg3N:  femStem + "ło",
		Pl1V:  virileStem + "liśmy",
		Pl1NV: femStem + "łyśmy",
		Pl2V:  virileStem + "liście",
		Pl2NV: femStem + "łyście",
		Pl3V:  virileStem + "li",
		Pl3NV: femStem + "ły",
	}
}

// buildPastTenseWithAlternation creates a past paradigm with ą→ę alternation.
// Used for -nąć verbs: ciągnąć → ciągnął/ciągnęła (ą stays in masc sg, ę elsewhere)
func buildPastTenseWithAlternation(mascStem, femStem string) PastTense {
	return PastTense{
		Sg1M:  mascStem + "łem",
		Sg1F:  femStem + "łam",
		Sg2M:  mascStem + "łeś",
		Sg2F:  femStem + "łaś",
		Sg3M:  mascStem + "ł",
		Sg3F:  femStem + "ła",
		Sg3N:  femStem + "ło",
		Pl1V:  femStem + "liśmy",
		Pl1NV: femStem + "łyśmy",
		Pl2V:  femStem + "liście",
		Pl2NV: femStem + "łyście",
		Pl3V:  femStem + "li",
		Pl3NV: femStem + "ły",
	}
}

// heuristicPastAsc handles -ąść and -ąźć verbs.
// trząść → trząsł/trzęsła/trzęśli (ą→ę in fem/neut/plural, ś softening in virile)
// prząść → prządł/przędła/przędli
func heuristicPastAsc(infinitive string) (PastTense, bool) {
	// -ąść verbs: trząść → trząsł/trzęsła
	if strings.HasSuffix(infinitive, "ząść") {
		prefix := strings.TrimSuffix(infinitive, "ząść")
		return PastTense{
			Sg1M:  prefix + "ząsłem",
			Sg1F:  prefix + "zęsłam",
			Sg2M:  prefix + "ząsłeś",
			Sg2F:  prefix + "zęsłaś",
			Sg3M:  prefix + "ząsł",
			Sg3F:  prefix + "zęsła",
			Sg3N:  prefix + "zęsło",
			Pl1V:  prefix + "zęśliśmy",
			Pl1NV: prefix + "zęsłyśmy",
			Pl2V:  prefix + "zęśliście",
			Pl2NV: prefix + "zęsłyście",
			Pl3V:  prefix + "zęśli",
			Pl3NV: prefix + "zęsły",
		}, true
	}

	// prząść → prządł/przędła
	if strings.HasSuffix(infinitive, "prząść") {
		prefix := strings.TrimSuffix(infinitive, "prząść")
		return PastTense{
			Sg1M:  prefix + "prządłem",
			Sg1F:  prefix + "przędłam",
			Sg2M:  prefix + "prządłeś",
			Sg2F:  prefix + "przędłaś",
			Sg3M:  prefix + "prządł",
			Sg3F:  prefix + "przędła",
			Sg3N:  prefix + "przędło",
			Pl1V:  prefix + "przędliśmy",
			Pl1NV: prefix + "przędłyśmy",
			Pl2V:  prefix + "przędliście",
			Pl2NV: prefix + "przędłyście",
			Pl3V:  prefix + "przędli",
			Pl3NV: prefix + "przędły",
		}, true
	}

	// -siąść verbs: siąść → siadł/siadła/siedli
	if strings.HasSuffix(infinitive, "siąść") {
		prefix := strings.TrimSuffix(infinitive, "siąść")
		return PastTense{
			Sg1M:  prefix + "siadłem",
			Sg1F:  prefix + "siadłam",
			Sg2M:  prefix + "siadłeś",
			Sg2F:  prefix + "siadłaś",
			Sg3M:  prefix + "siadł",
			Sg3F:  prefix + "siadła",
			Sg3N:  prefix + "siadło",
			Pl1V:  prefix + "siedliśmy",
			Pl1NV: prefix + "siadłyśmy",
			Pl2V:  prefix + "siedliście",
			Pl2NV: prefix + "siadłyście",
			Pl3V:  prefix + "siedli",
			Pl3NV: prefix + "siadły",
		}, true
	}

	return PastTense{}, false
}

// heuristicPastCzac handles -cząć verbs (począć, rozpocząć).
// począć → począł/poczęła/poczęli (ą→ę alternation)
func heuristicPastCzac(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "cząć") {
		return PastTense{}, false
	}
	prefix := strings.TrimSuffix(infinitive, "cząć")
	return PastTense{
		Sg1M:  prefix + "cząłem",
		Sg1F:  prefix + "częłam",
		Sg2M:  prefix + "cząłeś",
		Sg2F:  prefix + "częłaś",
		Sg3M:  prefix + "czął",
		Sg3F:  prefix + "częła",
		Sg3N:  prefix + "częło",
		Pl1V:  prefix + "częliśmy",
		Pl1NV: prefix + "częłyśmy",
		Pl2V:  prefix + "częliście",
		Pl2NV: prefix + "częłyście",
		Pl3V:  prefix + "częli",
		Pl3NV: prefix + "częły",
	}, true
}

// heuristicPastStrzyc handles -strzyc verbs (strzyc, ostrzyc).
// strzyc → strzygł/strzygła/strzygli (c→g alternation)
func heuristicPastStrzyc(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "strzyc") {
		return PastTense{}, false
	}
	prefix := strings.TrimSuffix(infinitive, "strzyc")
	return PastTense{
		Sg1M:  prefix + "strzygłem",
		Sg1F:  prefix + "strzygłam",
		Sg2M:  prefix + "strzygłeś",
		Sg2F:  prefix + "strzygłaś",
		Sg3M:  prefix + "strzygł",
		Sg3F:  prefix + "strzygła",
		Sg3N:  prefix + "strzygło",
		Pl1V:  prefix + "strzygli" + "śmy",
		Pl1NV: prefix + "strzygłyśmy",
		Pl2V:  prefix + "strzygli" + "ście",
		Pl2NV: prefix + "strzygłyście",
		Pl3V:  prefix + "strzygli",
		Pl3NV: prefix + "strzygły",
	}, true
}

// heuristicPastBosc handles -bość and -bóść verbs.
// bość/bóść → bódł/bodła (ó→o alternation, ść→dł)
func heuristicPastBosc(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "bość") && !strings.HasSuffix(infinitive, "bóść") {
		return PastTense{}, false
	}

	var prefix string
	if strings.HasSuffix(infinitive, "bóść") {
		prefix = strings.TrimSuffix(infinitive, "bóść")
	} else {
		prefix = strings.TrimSuffix(infinitive, "bość")
	}

	// Pattern: ó only in sg3m, o elsewhere
	return PastTense{
		Sg1M:  prefix + "bodłem",
		Sg1F:  prefix + "bodłam",
		Sg2M:  prefix + "bodłeś",
		Sg2F:  prefix + "bodłaś",
		Sg3M:  prefix + "bódł",
		Sg3F:  prefix + "bodła",
		Sg3N:  prefix + "bodło",
		Pl1V:  prefix + "bodliśmy",
		Pl1NV: prefix + "bodłyśmy",
		Pl2V:  prefix + "bodliście",
		Pl2NV: prefix + "bodłyście",
		Pl3V:  prefix + "bodli",
		Pl3NV: prefix + "bodły",
	}, true
}

// heuristicPastOwac handles -ować verbs.
// pracować → pracował, pracowała, pracowało...
func heuristicPastOwac(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "ować") {
		return PastTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ć") // pracowa
	return buildPastTense(stem), true
}

// heuristicPastYwacIwac handles -ywać and -iwać verbs.
// pokazywać → pokazywał
func heuristicPastYwacIwac(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "ywać") && !strings.HasSuffix(infinitive, "iwać") {
		return PastTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ć") // pokazywa
	return buildPastTense(stem), true
}

// heuristicPastAwac handles -awać verbs (not -ywać or -iwać).
// dawać → dawał
func heuristicPastAwac(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "awać") {
		return PastTense{}, false
	}
	if strings.HasSuffix(infinitive, "ywać") || strings.HasSuffix(infinitive, "iwać") {
		return PastTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ć") // dawa
	return buildPastTense(stem), true
}

// heuristicPastNac handles -nąć verbs.
// Main patterns:
// 1. N-retaining with ą→ę alternation: ciągnąć → ciągnął/ciągnęła (has ą in stem)
// 2. N-dropping: schnąć → schł/schła (inchoative/state-change verbs)
// 3. Mixed: cuchnąć → cuchnął/cuchnęła/cuchli (retain n in sg, drop in virile pl)
// 4. Prefixed dual-form: przekwitnąć → przekwitł/przekwitnęła (n-dropped sg3m, n-kept fem)
func heuristicPastNac(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "nąć") {
		return PastTense{}, false
	}

	// Base stem without -nąć
	stemWithoutNac := strings.TrimSuffix(infinitive, "nąć")

	// Check the mixed n-dropping verbs first
	// These retain n in singular/non-virile, drop n in virile plural
	if isKnownMixedNDropVerb(infinitive) {
		return buildPastTenseMixedNDrop(stemWithoutNac, infinitive), true
	}

	// Check the fully n-dropping verbs list
	// Also check for prefixed forms of n-dropping verbs
	if isKnownNDroppingVerb(infinitive) {
		// Use the n-dropped builder which handles consonant palatalization
		return buildPastTenseNDropped(stemWithoutNac, infinitive), true
	}

	// Check if the stem contains ą - if so, it's n-retaining with ą→ę alternation
	// ciągnąć → ciągnął/ciągnęła (ą in stem "ciąg")
	if strings.Contains(stemWithoutNac, "ą") {
		baseStem := strings.TrimSuffix(infinitive, "ąć") // ciągn
		mascStem := baseStem + "ą"                       // ciągną
		femStem := baseStem + "ę"                        // ciągnę
		return buildPastTenseWithAlternation(mascStem, femStem), true
	}

	// Default: N-retaining with ą→ę alternation (like kopnąć → kopnął/kopnęła)
	baseStem := strings.TrimSuffix(infinitive, "ąć") // kopn
	mascStem := baseStem + "ą"                       // kopną
	femStem := baseStem + "ę"                        // kopnę
	return buildPastTenseWithAlternation(mascStem, femStem), true
}

// Dual-form verbs whose prefixed forms use N-DROPPED sg3m.
// Example: przekwitnąć → przekwitł (sg3m), przekwitnęła (sg3f)
var dualBasesPrefixedNDropped = map[string]bool{
	"kwitnąć":    true, // bloom
	"brzęknąć":   true, // clang
	"pierzchnąć": true, // scatter
}

// Dual-form verbs whose prefixed forms use N-KEPT sg3m.
// Example: zatrzasnąć → zatrzasnął (sg3m), zatrzasnęła (sg3f)
var dualBasesPrefixedNKept = map[string]bool{
	"trzasnąć": true, // slam
	"śliznąć":  true, // slip
	"niknąć":   true, // vanish
	"siągnąć":  true, // reach
	"siąknąć":  true, // seep
	"sięknąć":  true, // seep (variant)
}

// isPrefixedDualFormVerb checks if this is a prefixed form of a dual-form base verb.
// Returns the base verb if found, empty string otherwise.
func getPrefixedDualFormBase(infinitive string) string {
	// Don't match if it's directly in the dual list (those return both paradigms)
	if isDualFormNacVerb(infinitive) {
		return ""
	}

	// Check if base is a dual-form verb
	base := extractBase(infinitive)
	if base != infinitive && isDualFormNacVerb(base) {
		return base
	}

	return ""
}

// isPrefixedDualFormNDropped returns true if this prefixed verb should use n-dropped sg3m.
func isPrefixedDualFormNDropped(infinitive string) bool {
	base := getPrefixedDualFormBase(infinitive)
	return base != "" && dualBasesPrefixedNDropped[base]
}

// isPrefixedDualFormNKept returns true if this prefixed verb should use n-kept sg3m.
func isPrefixedDualFormNKept(infinitive string) bool {
	base := getPrefixedDualFormBase(infinitive)
	return base != "" && dualBasesPrefixedNKept[base]
}

// buildPastTensePrefixedDualFormNDropped builds past tense for prefixed forms that use n-dropped sg3m.
// Pattern: n-dropped sg3m, n-kept feminine/neuter/non-virile, n-dropped virile.
// Example: przekwitnąć → przekwitł (sg3m), przekwitnęła (sg3f), przekwitli (pl3v)
func buildPastTensePrefixedDualFormNDropped(stemWithoutNac, infinitive string) PastTense {
	// Get the virile stem with palatalized final consonant (n-dropped)
	virileStem := palatalizeForVirile(stemWithoutNac, infinitive)

	// Check for masculine sg vowel alternation (ę→ą, e→a)
	mascStem := applyMascSgAlternation(stemWithoutNac, infinitive)

	// Apply sg3m-only alternations (o→ó, epenthetic e)
	sg3mStem := applySg3MOnlyAlternation(mascStem, infinitive)

	// N-kept stem for feminine/neuter/non-virile
	baseStem := strings.TrimSuffix(infinitive, "ąć") // przekwitn
	femStem := baseStem + "ę"                        // przekwitnę

	return PastTense{
		Sg1M:  mascStem + "łem",
		Sg1F:  femStem + "łam",
		Sg2M:  mascStem + "łeś",
		Sg2F:  femStem + "łaś",
		Sg3M:  sg3mStem + "ł",
		Sg3F:  femStem + "ła",
		Sg3N:  femStem + "ło",
		Pl1V:  virileStem + "liśmy",
		Pl1NV: femStem + "łyśmy",
		Pl2V:  virileStem + "liście",
		Pl2NV: femStem + "łyście",
		Pl3V:  virileStem + "li",
		Pl3NV: femStem + "ły",
	}
}

// buildPastTensePrefixedDualFormNKept builds past tense for prefixed forms that use n-kept sg3m.
// Pattern: n-kept sg3m, n-kept feminine/neuter/non-virile, n-kept virile.
// Example: zatrzasnąć → zatrzasnął (sg3m), zatrzasnęła (sg3f), zatrzasnęli (pl3v)
func buildPastTensePrefixedDualFormNKept(stemWithoutNac, infinitive string) PastTense {
	// N-kept stems
	baseStem := strings.TrimSuffix(infinitive, "ąć") // zatrzasn
	mascStem := baseStem + "ą"                       // zatrzasną
	femStem := baseStem + "ę"                        // zatrzasnę

	return PastTense{
		Sg1M:  mascStem + "łem",
		Sg1F:  femStem + "łam",
		Sg2M:  mascStem + "łeś",
		Sg2F:  femStem + "łaś",
		Sg3M:  mascStem + "ł",
		Sg3F:  femStem + "ła",
		Sg3N:  femStem + "ło",
		Pl1V:  femStem + "liśmy",
		Pl1NV: femStem + "łyśmy",
		Pl2V:  femStem + "liście",
		Pl2NV: femStem + "łyście",
		Pl3V:  femStem + "li",
		Pl3NV: femStem + "ły",
	}
}

// nRetainingVerbs are verbs that look like prefixed forms of n-dropping verbs
// but are actually separate lexemes that retain n in past tense.
// smoknąć (to get slapped) looks like s+moknąć but is NOT - it keeps n: smoknął
var nRetainingVerbs = map[string]bool{
	"smoknąć": true, // to get slapped, NOT s+moknąć (to get wet)
}

// Dual-form -nąć verbs have BOTH n-dropping and n-keeping sg3m variants.
// They split into two types based on their virile plural:
//
// Type 1 (virileDropped): virile plural uses n-dropped stem
//   kwitnąć → kwitł/kwitnął (sg3m) but kwitli (virile always n-dropped)
//
// Type 2 (virileKept): virile plural uses n-kept stem with ę
//   trzasnąć → trzasł/trzasnął (sg3m) but trzasnęli (virile always n-kept)

// dualFormNacVerbsVirileDropped - dual-form verbs with n-dropped virile plural
// Includes both base verbs with dual entries AND prefixed forms with dual entries
var dualFormNacVerbsVirileDropped = map[string]bool{
	// Base verbs with dual entries
	"buchnąć": true, "cuchnąć": true, "gęstnąć": true, "głuchnąć": true,
	"klęknąć": true, "kwitnąć": true, "mierzchnąć": true, "niknąć": true,
	"pełznąć": true, "pierzchnąć": true, "pizdnąć": true, "rymsnąć": true,
	"rypnąć": true, "sieknąć": true, "siągnąć": true, "siąknąć": true,
	"sięknąć": true, "spełgnąć": true, "stęgnąć": true,
	// Prefixed verbs with dual entries (base has single entry)
	"dosięgnąć": true, "napuchnąć": true, "ochlapnąć": true, "oklapnąć": true,
	"ostygnąć": true, "przesięgnąć": true, "przywyknąć": true, "spuchnąć": true,
	"ubodnąć": true, "wyziębnąć": true, "zgorzknąć": true,
	// Prefixed forms of dual bases that also have dual entries
	// From buchnąć:
	"wybuchnąć": true,
	// From cuchnąć:
	"zacuchnąć": true,
	// From gęstnąć:
	"zgęstnąć": true,
	// From klęknąć:
	"poklęknąć": true, "przyklęknąć": true, "uklęknąć": true, "zaklęknąć": true,
	// From kwitnąć:
	"dokwitnąć": true, "okwitnąć": true, "przekwitnąć": true,
	"rozkwitnąć": true, "wykwitnąć": true, "zakwitnąć": true,
	// From mierzchnąć:
	"pomierzchnąć": true,
	// From niknąć:
	"poniknąć": true, "wyniknąć": true, "zaniknąć": true, "zniknąć": true,
	// From pełznąć:
	"dopełznąć": true, "nadpełznąć": true, "odpełznąć": true, "opełznąć": true,
	"podpełznąć": true, "popełznąć": true, "przepełznąć": true, "przypełznąć": true,
	"rozpełznąć": true, "spełznąć": true, "wpełznąć": true, "wypełznąć": true, "zapełznąć": true,
	// From pierzchnąć:
	"rozpierzchnąć": true, "spierzchnąć": true,
	// From siągnąć:
	"dosiągnąć": true,
	// From siąknąć:
	"nasiąknąć": true, "osiąknąć": true, "podsiąknąć": true,
	"przesiąknąć": true, "wsiąknąć": true, "wysiąknąć": true,
	// From sięknąć:
	"przesięknąć": true, "wsięknąć": true,
}

// dualFormNacVerbsVirileKept - dual-form verbs with n-kept virile plural
var dualFormNacVerbsVirileKept = map[string]bool{
	// Base verbs with dual entries
	"brzęknąć": true, "chrypnąć": true, "prysnąć": true, "trysnąć": true,
	"trzasnąć": true, "wisnąć": true, "śliznąć": true,
	// Prefixed verbs with dual entries (base has single entry)
	"rozbłysnąć": true, "rozplasnąć": true, "rozplusnąć": true, "zabłysnąć": true,
	// Prefixed forms of dual bases that also have dual entries
	// From brzęknąć:
	"zabrzęknąć": true,
	// From prysnąć:
	"odprysnąć": true, "rozprysnąć": true, "sprysnąć": true, "wprysnąć": true, "wyprysnąć": true,
	// From trysnąć:
	"natrysnąć": true, "roztrysnąć": true, "wtrysnąć": true, "wytrysnąć": true,
	// From wisnąć:
	"nawisnąć": true, "obwisnąć": true, "owisnąć": true, "rozwisnąć": true,
	"uwisnąć": true, "zawisnąć": true, "zwisnąć": true,
	// From śliznąć:
	"obśliznąć": true, "ośliznąć": true,
}

// isDualFormNacVerb checks if a verb has dual n-drop/n-keep forms.
// Only returns true for verbs DIRECTLY listed - no prefix matching.
// This is because prefixed forms often have only one variant in the corpus
// even when the base verb has both.
func isDualFormNacVerb(infinitive string) bool {
	if !strings.HasSuffix(infinitive, "nąć") {
		return false
	}
	return dualFormNacVerbsVirileDropped[infinitive] || dualFormNacVerbsVirileKept[infinitive]
}

// isDualFormVirileKept returns true if this dual-form verb uses n-kept virile plural.
func isDualFormVirileKept(infinitive string) bool {
	return dualFormNacVerbsVirileKept[infinitive]
}

// isKnownNDroppingVerb checks if a verb is a known n-dropping -nąć verb.
// Checks both the base form and prefixed forms.
func isKnownNDroppingVerb(infinitive string) bool {
	// Check exclusion list first - these look like prefixed n-dropping verbs
	// but are separate lexemes that retain n
	if nRetainingVerbs[infinitive] {
		return false
	}

	// Direct lookup
	if nDroppingNacVerbs[infinitive] {
		return true
	}

	// Check for prefixed forms
	for _, prefix := range verbPrefixes {
		if strings.HasPrefix(infinitive, prefix) {
			base := infinitive[len(prefix):]
			if nDroppingNacVerbs[base] {
				return true
			}
		}
	}

	return false
}

// isKnownMixedNDropVerb checks if a verb has mixed n-dropping pattern.
// These verbs retain n in singular/non-virile, drop n in virile plural.
func isKnownMixedNDropVerb(infinitive string) bool {
	// Direct lookup
	if mixedNDropNacVerbs[infinitive] {
		return true
	}

	// Check for prefixed forms
	for _, prefix := range verbPrefixes {
		if strings.HasPrefix(infinitive, prefix) {
			base := infinitive[len(prefix):]
			if mixedNDropNacVerbs[base] {
				return true
			}
		}
	}

	return false
}

// isNDroppingVerb checks if a stem (without -nąć) indicates n-dropping.
// N-dropping occurs when there's a consonant cluster at the end of the stem
// without any vowel between (the -nąć is an inchoative suffix, not part of root).
// Examples: sch- (schnąć), marz- (marznąć), gas- (gasnąć)
// Counter-examples: kop- (kopnąć) has a clear vowel, n is part of the verbal root
func isNDroppingVerb(stem string) bool {
	if stem == "" {
		return false
	}

	runes := []rune(stem)
	if len(runes) < 2 {
		return false
	}

	// Check if stem ends in a consonant cluster (no vowel in last 2+ chars)
	// This indicates the -nąć is a suffix being added to a consonantal stem
	vowels := map[rune]bool{
		'a': true, 'e': true, 'i': true, 'o': true, 'u': true, 'y': true,
		'ą': true, 'ę': true, 'ó': true,
	}

	// Check last two characters for vowels
	hasRecentVowel := false
	for i := len(runes) - 1; i >= 0 && i >= len(runes)-2; i-- {
		if vowels[runes[i]] {
			hasRecentVowel = true
			break
		}
	}

	// If no vowel in last 2 chars, it's likely n-dropping
	// e.g., "sch" in schnąć, "marz" in marznąć
	return !hasRecentVowel
}

// heuristicPastSc handles -ść and -źć verbs.
// These have vowel alternations: ie→io/io, ó→o
// nieść → niósł/niosła, wieźć → wiózł/wiozła
func heuristicPastSc(infinitive string) (PastTense, bool) {
	// -mieść verbs: mieść → miótł/miotła (ie→ió/io, ść→tł)
	// Note: ó only in sg3m, o elsewhere (miotłem not miótłem)
	if strings.HasSuffix(infinitive, "mieść") {
		prefix := strings.TrimSuffix(infinitive, "mieść")
		return PastTense{
			Sg1M:  prefix + "miotłem",
			Sg1F:  prefix + "miotłam",
			Sg2M:  prefix + "miotłeś",
			Sg2F:  prefix + "miotłaś",
			Sg3M:  prefix + "miótł",
			Sg3F:  prefix + "miotła",
			Sg3N:  prefix + "miotło",
			Pl1V:  prefix + "mietliśmy",
			Pl1NV: prefix + "miotłyśmy",
			Pl2V:  prefix + "mietliście",
			Pl2NV: prefix + "miotłyście",
			Pl3V:  prefix + "mietli",
			Pl3NV: prefix + "miotły",
		}, true
	}

	// -gnieść verbs: gnieść → gniótł/gniotła
	// Note: ó only in sg3m, o elsewhere
	if strings.HasSuffix(infinitive, "gnieść") {
		prefix := strings.TrimSuffix(infinitive, "gnieść")
		return PastTense{
			Sg1M:  prefix + "gniotłem",
			Sg1F:  prefix + "gniotłam",
			Sg2M:  prefix + "gniotłeś",
			Sg2F:  prefix + "gniotłaś",
			Sg3M:  prefix + "gniótł",
			Sg3F:  prefix + "gniotła",
			Sg3N:  prefix + "gniotło",
			Pl1V:  prefix + "gnietliśmy",
			Pl1NV: prefix + "gniotłyśmy",
			Pl2V:  prefix + "gnietliście",
			Pl2NV: prefix + "gniotłyście",
			Pl3V:  prefix + "gnietli",
			Pl3NV: prefix + "gniotły",
		}, true
	}

	// -wieść verbs: wieść → wiódł/wiodła (lead)
	// Note: ó only in sg3m, o elsewhere
	if strings.HasSuffix(infinitive, "wieść") {
		prefix := strings.TrimSuffix(infinitive, "wieść")
		return PastTense{
			Sg1M:  prefix + "wiodłem",
			Sg1F:  prefix + "wiodłam",
			Sg2M:  prefix + "wiodłeś",
			Sg2F:  prefix + "wiodłaś",
			Sg3M:  prefix + "wiódł",
			Sg3F:  prefix + "wiodła",
			Sg3N:  prefix + "wiodło",
			Pl1V:  prefix + "wiedliśmy",
			Pl1NV: prefix + "wiodłyśmy",
			Pl2V:  prefix + "wiedliście",
			Pl2NV: prefix + "wiodłyście",
			Pl3V:  prefix + "wiedli",
			Pl3NV: prefix + "wiodły",
		}, true
	}

	// -ieść verbs (nieść type): ie→ió/io alternation
	// Note: ó only in sg3m, o elsewhere
	if strings.HasSuffix(infinitive, "ieść") {
		prefix := strings.TrimSuffix(infinitive, "ieść")
		return PastTense{
			Sg1M:  prefix + "iosłem",
			Sg1F:  prefix + "iosłam",
			Sg2M:  prefix + "iosłeś",
			Sg2F:  prefix + "iosłaś",
			Sg3M:  prefix + "iósł",
			Sg3F:  prefix + "iosła",
			Sg3N:  prefix + "iosło",
			Pl1V:  prefix + "ieśliśmy",
			Pl1NV: prefix + "iosłyśmy",
			Pl2V:  prefix + "ieśliście",
			Pl2NV: prefix + "iosłyście",
			Pl3V:  prefix + "ieśli",
			Pl3NV: prefix + "iosły",
		}, true
	}

	// -ieźć verbs (wieźć type): ie→ió/io alternation
	// Note: ó only in sg3m, o elsewhere
	if strings.HasSuffix(infinitive, "ieźć") {
		prefix := strings.TrimSuffix(infinitive, "ieźć")
		return PastTense{
			Sg1M:  prefix + "iozłem",
			Sg1F:  prefix + "iozłam",
			Sg2M:  prefix + "iozłeś",
			Sg2F:  prefix + "iozłaś",
			Sg3M:  prefix + "iózł",
			Sg3F:  prefix + "iozła",
			Sg3N:  prefix + "iozło",
			Pl1V:  prefix + "ieźliśmy",
			Pl1NV: prefix + "iozłyśmy",
			Pl2V:  prefix + "ieźliście",
			Pl2NV: prefix + "iozłyście",
			Pl3V:  prefix + "ieźli",
			Pl3NV: prefix + "iozły",
		}, true
	}

	// -yźć verbs (gryźć type): no vowel alternation
	if strings.HasSuffix(infinitive, "yźć") {
		stem := strings.TrimSuffix(infinitive, "źć") // gryz
		return PastTense{
			Sg1M:  stem + "złem",
			Sg1F:  stem + "złam",
			Sg2M:  stem + "złeś",
			Sg2F:  stem + "złaś",
			Sg3M:  stem + "zł",
			Sg3F:  stem + "zła",
			Sg3N:  stem + "zło",
			Pl1V:  stem + "źliśmy",
			Pl1NV: stem + "złyśmy",
			Pl2V:  stem + "źliście",
			Pl2NV: stem + "złyście",
			Pl3V:  stem + "źli",
			Pl3NV: stem + "zły",
		}, true
	}

	// -eźć verbs (leźć type): no vowel alternation
	if strings.HasSuffix(infinitive, "eźć") {
		stem := strings.TrimSuffix(infinitive, "źć") // lez
		return PastTense{
			Sg1M:  stem + "złem",
			Sg1F:  stem + "złam",
			Sg2M:  stem + "złeś",
			Sg2F:  stem + "złaś",
			Sg3M:  stem + "zł",
			Sg3F:  stem + "zła",
			Sg3N:  stem + "zło",
			Pl1V:  stem + "źliśmy",
			Pl1NV: stem + "złyśmy",
			Pl2V:  stem + "źliście",
			Pl2NV: stem + "złyście",
			Pl3V:  stem + "źli",
			Pl3NV: stem + "zły",
		}, true
	}

	return PastTense{}, false
}

// heuristicPastC handles -c verbs (móc, piec, etc.).
// móc → mógł/mogła (ó→o alternation)
// piec → piekł/piekła
func heuristicPastC(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "c") {
		return PastTense{}, false
	}
	// Skip -ść/-źć (handled above) and vowel+ć patterns (handled below)
	if strings.HasSuffix(infinitive, "ść") || strings.HasSuffix(infinitive, "źć") ||
		strings.HasSuffix(infinitive, "ać") || strings.HasSuffix(infinitive, "eć") ||
		strings.HasSuffix(infinitive, "ić") || strings.HasSuffix(infinitive, "yć") ||
		strings.HasSuffix(infinitive, "uć") || strings.HasSuffix(infinitive, "ąć") {
		return PastTense{}, false
	}

	// móc type: ó→o alternation, c→g (ó only in sg3m)
	if strings.HasSuffix(infinitive, "óc") {
		prefix := strings.TrimSuffix(infinitive, "óc")
		return PastTense{
			Sg1M:  prefix + "ogłem",
			Sg1F:  prefix + "ogłam",
			Sg2M:  prefix + "ogłeś",
			Sg2F:  prefix + "ogłaś",
			Sg3M:  prefix + "ógł",
			Sg3F:  prefix + "ogła",
			Sg3N:  prefix + "ogło",
			Pl1V:  prefix + "ogliśmy",
			Pl1NV: prefix + "ogłyśmy",
			Pl2V:  prefix + "ogliście",
			Pl2NV: prefix + "ogłyście",
			Pl3V:  prefix + "ogli",
			Pl3NV: prefix + "ogły",
		}, true
	}

	// piec type: c→k
	if strings.HasSuffix(infinitive, "ec") {
		stem := strings.TrimSuffix(infinitive, "c") // pie
		return PastTense{
			Sg1M:  stem + "kłem",
			Sg1F:  stem + "kłam",
			Sg2M:  stem + "kłeś",
			Sg2F:  stem + "kłaś",
			Sg3M:  stem + "kł",
			Sg3F:  stem + "kła",
			Sg3N:  stem + "kło",
			Pl1V:  stem + "kliśmy",
			Pl1NV: stem + "kłyśmy",
			Pl2V:  stem + "kliście",
			Pl2NV: stem + "kłyście",
			Pl3V:  stem + "kli",
			Pl3NV: stem + "kły",
		}, true
	}

	return PastTense{}, false
}

// heuristicPastIc handles -ić verbs.
// robić → robił, robiła, robiło, robili, robiły
func heuristicPastIc(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "ić") {
		return PastTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ć") // robi
	return buildPastTense(stem), true
}

// heuristicPastYc handles -yć verbs.
// myć → mył, myła, myło, myli, myły
func heuristicPastYc(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "yć") {
		return PastTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ć") // my
	return buildPastTense(stem), true
}

// heuristicPastUc handles -uć verbs.
// czuć → czuł, czuła, czuło, czuli, czuły
func heuristicPastUc(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "uć") {
		return PastTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ć") // czu
	return buildPastTense(stem), true
}

// heuristicPastEc handles -eć verbs.
// Pattern: singular/non-virile uses -ał, virile uses -eli
// umieć → umiał/umiała/umieli (not umiali!)
// leżeć → leżał/leżała/leżeli
func heuristicPastEc(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "eć") {
		return PastTense{}, false
	}

	// Get base stem without -eć
	baseStem := strings.TrimSuffix(infinitive, "eć")

	// Singular and non-virile plural use -ał- stem
	var aStem string
	if strings.HasSuffix(infinitive, "ieć") {
		// -ieć verbs: ie → ia (umieć → umia + ł)
		aStem = strings.TrimSuffix(infinitive, "ieć") + "ia"
	} else {
		// Other -eć verbs: e → a (leżeć → leża + ł)
		aStem = baseStem + "a"
	}

	// Virile plural uses -eli stem (keeps the e)
	eStem := baseStem + "e"

	return PastTense{
		Sg1M:  aStem + "łem",
		Sg1F:  aStem + "łam",
		Sg2M:  aStem + "łeś",
		Sg2F:  aStem + "łaś",
		Sg3M:  aStem + "ł",
		Sg3F:  aStem + "ła",
		Sg3N:  aStem + "ło",
		Pl1V:  eStem + "liśmy",
		Pl1NV: aStem + "łyśmy",
		Pl2V:  eStem + "liście",
		Pl2NV: aStem + "łyście",
		Pl3V:  eStem + "li",
		Pl3NV: aStem + "ły",
	}, true
}

// heuristicPastAc handles regular -ać verbs (fallback).
// czytać → czytał, czytała, czytało, czytali, czytały
func heuristicPastAc(infinitive string) (PastTense, bool) {
	if !strings.HasSuffix(infinitive, "ać") {
		return PastTense{}, false
	}
	stem := strings.TrimSuffix(infinitive, "ć") // czyta
	return buildPastTense(stem), true
}
