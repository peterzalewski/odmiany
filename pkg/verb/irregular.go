package verb

// homographs contains verbs with multiple valid paradigms (different meanings).
// These are checked first before irregular verbs.
var homographs = map[string][]Paradigm{
	// stać: "to stand" (imperfective) vs "to become/afford" (perfective)
	"stać": {
		{
			PresentTense: PresentTense{
				Sg1: "stoję", Sg2: "stoisz", Sg3: "stoi",
				Pl1: "stoimy", Pl2: "stoicie", Pl3: "stoją",
			},
			Gloss: "to stand",
		},
		{
			PresentTense: PresentTense{
				Sg1: "stanę", Sg2: "staniesz", Sg3: "stanie",
				Pl1: "staniemy", Pl2: "staniecie", Pl3: "staną",
			},
			Gloss: "to become, to afford",
		},
	},
	// słać: "to send" vs "to spread (bedding)"
	"słać": {
		{
			PresentTense: PresentTense{
				Sg1: "ślę", Sg2: "ślesz", Sg3: "śle",
				Pl1: "ślemy", Pl2: "ślecie", Pl3: "ślą",
			},
			Gloss: "to send",
		},
		{
			PresentTense: PresentTense{
				Sg1: "ścielę", Sg2: "ścielesz", Sg3: "ściele",
				Pl1: "ścielemy", Pl2: "ścielecie", Pl3: "ścielą",
			},
			Gloss: "to spread (bedding)",
		},
	},
}

// lookupHomograph returns all paradigms for a homograph verb.
func lookupHomograph(infinitive string) ([]Paradigm, bool) {
	// Direct lookup
	if paradigms, ok := homographs[infinitive]; ok {
		return paradigms, true
	}

	// Check for prefixed forms of homographs
	for _, prefix := range verbPrefixes {
		if len(infinitive) > len(prefix) && infinitive[:len(prefix)] == prefix {
			base := infinitive[len(prefix):]
			if baseParadigms, ok := homographs[base]; ok {
				// Apply prefix to all paradigms
				result := make([]Paradigm, len(baseParadigms))
				for i, bp := range baseParadigms {
					result[i] = Paradigm{
						PresentTense: PresentTense{
							Sg1: prefix + bp.Sg1,
							Sg2: prefix + bp.Sg2,
							Sg3: prefix + bp.Sg3,
							Pl1: prefix + bp.Pl1,
							Pl2: prefix + bp.Pl2,
							Pl3: prefix + bp.Pl3,
						},
						Gloss: bp.Gloss,
					}
				}
				return result, true
			}
		}
	}

	return nil, false
}

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
	"musieć": {
		Sg1: "muszę", Sg2: "musisz", Sg3: "musi",
		Pl1: "musimy", Pl2: "musicie", Pl3: "muszą",
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

	// naleźć - suppletive stem najd- (base for znaleźć, odnaleźć, etc.)
	// naleźć → najdę, najdziesz, najdzie...
	"naleźć": {
		Sg1: "najdę", Sg2: "najdziesz", Sg3: "najdzie",
		Pl1: "najdziemy", Pl2: "najdziecie", Pl3: "najdą",
	},

	// spać - stem change s→ś before palatal
	// spać → śpię, śpisz, śpi...
	"spać": {
		Sg1: "śpię", Sg2: "śpisz", Sg3: "śpi",
		Pl1: "śpimy", Pl2: "śpicie", Pl3: "śpią",
	},

	// bać się - suppletive stem boj-
	// bać → boję, boisz, boi...
	"bać": {
		Sg1: "boję", Sg2: "boisz", Sg3: "boi",
		Pl1: "boimy", Pl2: "boicie", Pl3: "boją",
	},

	// dziać się - special -ać → -eję pattern
	// dziać → dzieję, dziejesz, dzieje...
	"dziać": {
		Sg1: "dzieję", Sg2: "dziejesz", Sg3: "dzieje",
		Pl1: "dziejemy", Pl2: "dziejecie", Pl3: "dzieją",
	},

	// podobać się - regular -am (not alternating like most -bać)
	"podobać": {
		Sg1: "podobam", Sg2: "podobasz", Sg3: "podoba",
		Pl1: "podobamy", Pl2: "podobacie", Pl3: "podobają",
	},

	// Monosyllabic -ić verbs with j-insertion
	// bić → biję (zabić → zabiję, pobić → pobiję, etc.)
	"bić": {
		Sg1: "biję", Sg2: "bijesz", Sg3: "bije",
		Pl1: "bijemy", Pl2: "bijecie", Pl3: "biją",
	},
	// lić → liję (oblić → obleję, wylić → wyleję, etc.)
	"lić": {
		Sg1: "liję", Sg2: "lijesz", Sg3: "lije",
		Pl1: "lijemy", Pl2: "lijecie", Pl3: "liją",
	},
	// pić → piję (wypić, napić, etc.)
	"pić": {
		Sg1: "piję", Sg2: "pijesz", Sg3: "pije",
		Pl1: "pijemy", Pl2: "pijecie", Pl3: "piją",
	},
	// żyć → żyję (przeżyć, wyżyć, etc.)
	"żyć": {
		Sg1: "żyję", Sg2: "żyjesz", Sg3: "żyje",
		Pl1: "żyjemy", Pl2: "żyjecie", Pl3: "żyją",
	},
	// myć → myję (umyć, wymyć, etc.)
	"myć": {
		Sg1: "myję", Sg2: "myjesz", Sg3: "myje",
		Pl1: "myjemy", Pl2: "myjecie", Pl3: "myją",
	},
	// ryć → ryję (wyryć, zaryć, etc.)
	"ryć": {
		Sg1: "ryję", Sg2: "ryjesz", Sg3: "ryje",
		Pl1: "ryjemy", Pl2: "ryjecie", Pl3: "ryją",
	},
	// szyć → szyję (zszyć, uszyć, etc.)
	"szyć": {
		Sg1: "szyję", Sg2: "szyjesz", Sg3: "szyje",
		Pl1: "szyjemy", Pl2: "szyjecie", Pl3: "szyją",
	},
	// wyć → wyję (zawyć, etc.)
	"wyć": {
		Sg1: "wyję", Sg2: "wyjesz", Sg3: "wyje",
		Pl1: "wyjemy", Pl2: "wyjecie", Pl3: "wyją",
	},
	// kryć → kryję (ukryć, odkryć, etc.)
	"kryć": {
		Sg1: "kryję", Sg2: "kryjesz", Sg3: "kryje",
		Pl1: "kryjemy", Pl2: "kryjecie", Pl3: "kryją",
	},

	// -pomnieć verbs: zapomnieć → zapomnę (ie drops)
	"pomnieć": {
		Sg1: "pomnę", Sg2: "pomnisz", Sg3: "pomni",
		Pl1: "pomnimy", Pl2: "pomnicie", Pl3: "pomną",
	},

	// -mrzeć verbs: umrzeć → umrę (rz→r)
	"mrzeć": {
		Sg1: "mrę", Sg2: "mrzesz", Sg3: "mrze",
		Pl1: "mrzemy", Pl2: "mrzecie", Pl3: "mrą",
	},

	// ciec verbs: uciec → ucieknę (k-insertion like biec)
	"ciec": {
		Sg1: "cieknę", Sg2: "ciekniesz", Sg3: "cieknie",
		Pl1: "ciekniemy", Pl2: "ciekniecie", Pl3: "ciekną",
	},

	// woleć - special -eć verb
	"woleć": {
		Sg1: "wolę", Sg2: "wolisz", Sg3: "woli",
		Pl1: "wolimy", Pl2: "wolicie", Pl3: "wolą",
	},

	// -jąć verbs: suppletive stem -jm-/-m-
	// zająć → zajmę, przyjąć → przyjmę, etc.
	"jąć": {
		Sg1: "jmę", Sg2: "jmiesz", Sg3: "jmie",
		Pl1: "jmiemy", Pl2: "jmiecie", Pl3: "jmą",
	},
	// Special -jąć verbs with e-insertion (consonant cluster + jąć)
	// zdjąć → zdejmę (z + d + jąć, d insertion)
	"zdjąć": {
		Sg1: "zdejmę", Sg2: "zdejmiesz", Sg3: "zdejmie",
		Pl1: "zdejmiemy", Pl2: "zdejmiecie", Pl3: "zdejmą",
	},
	// podjąć → podejmę (pod + jąć, e insertion)
	"podjąć": {
		Sg1: "podejmę", Sg2: "podejmiesz", Sg3: "podejmie",
		Pl1: "podejmiemy", Pl2: "podejmiecie", Pl3: "podejmą",
	},
	// odjąć → odejmę
	"odjąć": {
		Sg1: "odejmę", Sg2: "odejmiesz", Sg3: "odejmie",
		Pl1: "odejmiemy", Pl2: "odejmiecie", Pl3: "odejmą",
	},
	// wziąć is already in the list as suppletive

	// począć → pocznę (czn- stem, different from -jąć)
	"cząć": {
		Sg1: "cznę", Sg2: "czniesz", Sg3: "cznie",
		Pl1: "czniemy", Pl2: "czniecie", Pl3: "czną",
	},
	// począć family - add as direct entries since prefix analysis is complex
	"począć": {
		Sg1: "pocznę", Sg2: "poczniesz", Sg3: "pocznie",
		Pl1: "poczniemy", Pl2: "poczniecie", Pl3: "poczną",
	},
	"odpocząć": {
		Sg1: "odpocznę", Sg2: "odpoczniesz", Sg3: "odpocznie",
		Pl1: "odpoczniemy", Pl2: "odpoczniecie", Pl3: "odpoczną",
	},
	"rozpocząć": {
		Sg1: "rozpocznę", Sg2: "rozpoczniesz", Sg3: "rozpocznie",
		Pl1: "rozpoczniemy", Pl2: "rozpoczniecie", Pl3: "rozpoczną",
	},

	// Action verb -mieć patterns (grzmieć → grzmię, not grzmieję)
	"grzmieć": {
		Sg1: "grzmię", Sg2: "grzmisz", Sg3: "grzmi",
		Pl1: "grzmimy", Pl2: "grzmicie", Pl3: "grzmią",
	},
	"szumieć": {
		Sg1: "szumię", Sg2: "szumisz", Sg3: "szumi",
		Pl1: "szumimy", Pl2: "szumicie", Pl3: "szumią",
	},
	"tłumieć": {
		Sg1: "tłumię", Sg2: "tłumisz", Sg3: "tłumi",
		Pl1: "tłumimy", Pl2: "tłumicie", Pl3: "tłumią",
	},

	// patrzeć - action verb (patrzę, not patrzeję)
	"patrzeć": {
		Sg1: "patrzę", Sg2: "patrzysz", Sg3: "patrzy",
		Pl1: "patrzymy", Pl2: "patrzycie", Pl3: "patrzą",
	},

	// -rwać verbs: use -ę/-ie pattern (not -am/-asz)
	// rwać → rwę, derwać → derwę, wyrwać → wyrwę
	"rwać": {
		Sg1: "rwę", Sg2: "rwiesz", Sg3: "rwie",
		Pl1: "rwiemy", Pl2: "rwiecie", Pl3: "rwą",
	},

	// -zwać verbs: use -ę/-ie pattern
	// zwać → zwę, nazwać → nazwę, wezwać → wezwę
	"zwać": {
		Sg1: "zwę", Sg2: "zwiesz", Sg3: "zwie",
		Pl1: "zwiemy", Pl2: "zwiecie", Pl3: "zwą",
	},

	// dbać - regular -am (not alternating like other -bać)
	"dbać": {
		Sg1: "dbam", Sg2: "dbasz", Sg3: "dba",
		Pl1: "dbamy", Pl2: "dbacie", Pl3: "dbają",
	},

	// śmiać się - special pattern (ać → eję)
	"śmiać": {
		Sg1: "śmieję", Sg2: "śmiejesz", Sg3: "śmieje",
		Pl1: "śmiejemy", Pl2: "śmiejecie", Pl3: "śmieją",
	},

	// -ieć action verbs (cierpieć, pachnieć, wisieć, tkwieć, etc.)
	"cierpieć": {
		Sg1: "cierpię", Sg2: "cierpisz", Sg3: "cierpi",
		Pl1: "cierpimy", Pl2: "cierpicie", Pl3: "cierpią",
	},
	"wisieć": {
		Sg1: "wiszę", Sg2: "wisisz", Sg3: "wisi",
		Pl1: "wisimy", Pl2: "wisicie", Pl3: "wiszą",
	},
	"tkwieć": {
		Sg1: "tkwię", Sg2: "tkwisz", Sg3: "tkwi",
		Pl1: "tkwimy", Pl2: "tkwicie", Pl3: "tkwią",
	},
	"śmierdzieć": {
		Sg1: "śmierdzę", Sg2: "śmierdzisz", Sg3: "śmierdzi",
		Pl1: "śmierdzimy", Pl2: "śmierdzicie", Pl3: "śmierdzą",
	},

	// jeździć - correct softening źdź → żdż
	"jeździć": {
		Sg1: "jeżdżę", Sg2: "jeździsz", Sg3: "jeździ",
		Pl1: "jeździmy", Pl2: "jeździcie", Pl3: "jeżdżą",
	},

	// -nieć action verbs (pachnieć, etc.)
	"pachnieć": {
		Sg1: "pachnę", Sg2: "pachniesz", Sg3: "pachnie",
		Pl1: "pachniemy", Pl2: "pachniecie", Pl3: "pachną",
	},

	// -strzec verbs: c→g alternation (not c→k)
	"strzec": {
		Sg1: "strzegę", Sg2: "strzeżesz", Sg3: "strzeże",
		Pl1: "strzeżemy", Pl2: "strzeżecie", Pl3: "strzegą",
	},

	// -chować verbs: use -owam (not -uję)
	"chować": {
		Sg1: "chowam", Sg2: "chowasz", Sg3: "chowa",
		Pl1: "chowamy", Pl2: "chowacie", Pl3: "chowają",
	},

	// okazać - minority alternating -zać (z→ż)
	"okazać": {
		Sg1: "okażę", Sg2: "okażesz", Sg3: "okaże",
		Pl1: "okażemy", Pl2: "okażecie", Pl3: "okażą",
	},

	// karać - minority alternating -rać (r→rz)
	"karać": {
		Sg1: "karzę", Sg2: "karzesz", Sg3: "karze",
		Pl1: "karzemy", Pl2: "karzecie", Pl3: "karzą",
	},

	// -kraść verbs: suppletive kradn- stem
	"kraść": {
		Sg1: "kradnę", Sg2: "kradniesz", Sg3: "kradnie",
		Pl1: "kradniemy", Pl2: "kradniecie", Pl3: "kradną",
	},

	// -kłaść verbs: suppletive kład- stem
	"kłaść": {
		Sg1: "kładę", Sg2: "kładziesz", Sg3: "kładzie",
		Pl1: "kładziemy", Pl2: "kładziecie", Pl3: "kładą",
	},

	// uczcić - needs szcz (not szc)
	"uczcić": {
		Sg1: "uczczę", Sg2: "uczcisz", Sg3: "uczci",
		Pl1: "uczcimy", Pl2: "uczcicie", Pl3: "uczczą",
	},

	// wspomnieć - special prefix form (w + historical root)
	"wspomnieć": {
		Sg1: "wspomnę", Sg2: "wspomnisz", Sg3: "wspomni",
		Pl1: "wspomnimy", Pl2: "wspomnicie", Pl3: "wspomną",
	},

	// opisać - minority alternating -sać
	"opisać": {
		Sg1: "opiszę", Sg2: "opiszesz", Sg3: "opisze",
		Pl1: "opiszemy", Pl2: "opiszecie", Pl3: "opiszą",
	},

	// wskazać - minority alternating -zać
	"wskazać": {
		Sg1: "wskażę", Sg2: "wskażesz", Sg3: "wskaże",
		Pl1: "wskażemy", Pl2: "wskażecie", Pl3: "wskażą",
	},

	// brać prefix verbs with vowel elision
	// ode+brać → odbiorę (not odebiorę)
	"odebrać": {
		Sg1: "odbiorę", Sg2: "odbierzesz", Sg3: "odbierze",
		Pl1: "odbierzemy", Pl2: "odbierzecie", Pl3: "odbiorą",
	},
	"zebrać": {
		Sg1: "zbiorę", Sg2: "zbierzesz", Sg3: "zbierze",
		Pl1: "zbierzemy", Pl2: "zbierzecie", Pl3: "zbiorą",
	},
	"rozebrać": {
		Sg1: "rozbiorę", Sg2: "rozbierzesz", Sg3: "rozbierze",
		Pl1: "rozbierzemy", Pl2: "rozbierzecie", Pl3: "rozbiorą",
	},

	// lać verbs (j-insertion like myć)
	"lać": {
		Sg1: "leję", Sg2: "lejesz", Sg3: "leje",
		Pl1: "lejemy", Pl2: "lejecie", Pl3: "leją",
	},

	// pogrześć - suppletive grzeb- stem
	"grześć": {
		Sg1: "grzebę", Sg2: "grzebiesz", Sg3: "grzebie",
		Pl1: "grzebiemy", Pl2: "grzebiecie", Pl3: "grzebą",
	},

	// -przeć/-wrzeć verbs: special patterns
	// oprzeć → oprę (not oprzę)
	"przeć": {
		Sg1: "prę", Sg2: "przesz", Sg3: "prze",
		Pl1: "przemy", Pl2: "przecie", Pl3: "prą",
	},
	// zawrzeć → zawrę
	"wrzeć": {
		Sg1: "wrę", Sg2: "wrzesz", Sg3: "wrze",
		Pl1: "wrzemy", Pl2: "wrzecie", Pl3: "wrą",
	},

	// śnić - no j-insertion (śnię not śniję)
	"śnić": {
		Sg1: "śnię", Sg2: "śnisz", Sg3: "śni",
		Pl1: "śnimy", Pl2: "śnicie", Pl3: "śnią",
	},

	// rzec - k-insertion like biec
	"rzec": {
		Sg1: "rzeknę", Sg2: "rzeczesz", Sg3: "rzecze",
		Pl1: "rzeczemy", Pl2: "rzeczecie", Pl3: "rzekną",
	},

	// wiać - special pattern (wieję not wiam)
	"wiać": {
		Sg1: "wieję", Sg2: "wiejesz", Sg3: "wieje",
		Pl1: "wiejemy", Pl2: "wiejecie", Pl3: "wieją",
	},

	// krajać - j-insertion (kraję not krajam)
	"krajać": {
		Sg1: "kraję", Sg2: "krajesz", Sg3: "kraje",
		Pl1: "krajemy", Pl2: "krajecie", Pl3: "krają",
	},

	// nająć - jm- stem like jąć
	"nająć": {
		Sg1: "najmę", Sg2: "najmiesz", Sg3: "najmie",
		Pl1: "najmiemy", Pl2: "najmiecie", Pl3: "najmą",
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
	"pod", "pode", "nad", "nade", "wz", "wze", "u", "s", "z", "ze", "w", "we", "o",
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
		"naleźć": true, "spać": true, "bać": true, "dziać": true,
		"podobać": true,
		// Monosyllabic verbs
		"bić": true, "lić": true, "pić": true, "żyć": true, "myć": true,
		"ryć": true, "szyć": true, "wyć": true, "kryć": true,
		// Other prefixable bases
		"pomnieć": true, "mrzeć": true, "ciec": true, "woleć": true,
		"jąć": true, "cząć": true, "patrzeć": true,
		"rwać": true, "zwać": true, "dbać": true, "śmiać": true,
		"cierpieć": true, "wisieć": true, "jeździć": true,
		"pachnieć": true, "strzec": true, "chować": true,
		"okazać": true, "karać": true, "kraść": true, "kłaść": true,
		"lać": true, "grześć": true, "przeć": true, "wrzeć": true,
		"śnić": true, "rzec": true, "wiać": true, "krajać": true,
		"słać": true, "nająć": true,
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
