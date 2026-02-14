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
	// boleć: "physical pain" vs "to grieve/worry" (inchoative)
	"boleć": {
		{
			PresentTense: PresentTense{
				Sg1: "bolę", Sg2: "bolisz", Sg3: "boli",
				Pl1: "bolimy", Pl2: "bolicie", Pl3: "bolą",
			},
			Gloss: "to hurt (physical pain)",
		},
		{
			PresentTense: PresentTense{
				Sg1: "boleję", Sg2: "bolejesz", Sg3: "boleje",
				Pl1: "bolejemy", Pl2: "bolejecie", Pl3: "boleją",
			},
			Gloss: "to grieve, to worry",
		},
	},
	// stajać: frequentative of stać (both patterns attested)
	"stajać": {
		{
			PresentTense: PresentTense{
				Sg1: "staję", Sg2: "stajesz", Sg3: "staje",
				Pl1: "stajemy", Pl2: "stajecie", Pl3: "stają",
			},
			Gloss: "to keep standing/stopping (frequentative)",
		},
		{
			PresentTense: PresentTense{
				Sg1: "stajam", Sg2: "stajasz", Sg3: "staja",
				Pl1: "stajamy", Pl2: "stajacie", Pl3: "stajają",
			},
			Gloss: "to keep standing/stopping (variant)",
		},
	},
	// chlać: vulgar "to gulp" (both patterns attested)
	"chlać": {
		{
			PresentTense: PresentTense{
				Sg1: "chlam", Sg2: "chlasz", Sg3: "chla",
				Pl1: "chlamy", Pl2: "chlacie", Pl3: "chlają",
			},
			Gloss: "to gulp/slurp (vulgar)",
		},
		{
			PresentTense: PresentTense{
				Sg1: "chleję", Sg2: "chlejesz", Sg3: "chleje",
				Pl1: "chlejemy", Pl2: "chlejecie", Pl3: "chleją",
			},
			Gloss: "to gulp/slurp (variant)",
		},
	},
}

// lookupHomograph returns all paradigms for a homograph verb.
func lookupHomograph(infinitive string) ([]Paradigm, bool) {
	// Direct lookup - only the bare form, not prefixed forms
	// Prefixed forms like "dostać", "przestać" are NOT homographs:
	// they only use one paradigm (stanę), handled by heuristics.
	if paradigms, ok := homographs[infinitive]; ok {
		return paradigms, true
	}

	// Only słać and chlać support prefix expansion for homographs
	// because their prefixed forms retain both conjugation patterns.
	// stać prefixed forms like "dostać" are NOT homographs.
	for _, prefix := range verbPrefixes {
		if len(infinitive) > len(prefix) && infinitive[:len(prefix)] == prefix {
			base := infinitive[len(prefix):]
			// Expand słać and chlać homographs
			if base == "słać" || base == "chlać" {
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
	// sprzedać - compound prefix sprzedać (s+prze+dać)
	"sprzedać": {
		Sg1: "sprzedam", Sg2: "sprzedasz", Sg3: "sprzeda",
		Pl1: "sprzedamy", Pl2: "sprzedacie", Pl3: "sprzedadzą",
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
		Sg1: "kasam", Sg2: "kasasz", Sg3: "kasa",
		Pl1: "kasamy", Pl2: "kasacie", Pl3: "kasają",
	},
	"kołysać": {
		Sg1: "kołyszę", Sg2: "kołyszesz", Sg3: "kołysze",
		Pl1: "kołyszemy", Pl2: "kołyszecie", Pl3: "kołyszą",
	},
	"ciosać": {
		Sg1: "ciosam", Sg2: "ciosasz", Sg3: "ciosa",
		Pl1: "ciosamy", Pl2: "ciosacie", Pl3: "ciosają",
	},
	"ciesać": {
		Sg1: "ciesam", Sg2: "ciesasz", Sg3: "ciesa",
		Pl1: "ciesamy", Pl2: "ciesacie", Pl3: "ciesają",
	},
	"krzesać": {
		Sg1: "krzesam", Sg2: "krzesasz", Sg3: "krzesa",
		Pl1: "krzesamy", Pl2: "krzesacie", Pl3: "krzesają",
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
	// wić → wiję (powić, spowić, etc.) - j-insertion
	"wić": {
		Sg1: "wiję", Sg2: "wijesz", Sg3: "wije",
		Pl1: "wijemy", Pl2: "wijecie", Pl3: "wiją",
	},
	// Prefixed wić forms that don't match simple prefix patterns
	"opowić": {
		Sg1: "opowiję", Sg2: "opowijesz", Sg3: "opowije",
		Pl1: "opowijemy", Pl2: "opowijecie", Pl3: "opowiją",
	},
	"rozpowić": {
		Sg1: "rozpowiję", Sg2: "rozpowijesz", Sg3: "rozpowije",
		Pl1: "rozpowijemy", Pl2: "rozpowijecie", Pl3: "rozpowiją",
	},
	"spowić": {
		Sg1: "spowiję", Sg2: "spowijesz", Sg3: "spowije",
		Pl1: "spowijemy", Pl2: "spowijecie", Pl3: "spowiją",
	},
	"upowić": {
		Sg1: "upowiję", Sg2: "upowijesz", Sg3: "upowije",
		Pl1: "upowijemy", Pl2: "upowijecie", Pl3: "upowiją",
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
	"spocząć": {
		Sg1: "spocznę", Sg2: "spoczniesz", Sg3: "spocznie",
		Pl1: "spoczniemy", Pl2: "spoczniecie", Pl3: "spoczną",
	},
	"wypocząć": {
		Sg1: "wypocznę", Sg2: "wypoczniesz", Sg3: "wypocznie",
		Pl1: "wypoczniemy", Pl2: "wypoczniecie", Pl3: "wypoczną",
	},
	"wszcząć": {
		Sg1: "wszcznę", Sg2: "wszczniesz", Sg3: "wszcznie",
		Pl1: "wszczniemy", Pl2: "wszczniecie", Pl3: "wszczną",
	},
	"poczęć": {
		Sg1: "pocznę", Sg2: "poczniesz", Sg3: "pocznie",
		Pl1: "poczniemy", Pl2: "poczniecie", Pl3: "poczną",
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

	// Inchoative -rzeć verbs (starzeć się, gorzeć) - use -eję pattern
	"starzeć": {
		Sg1: "starzeję", Sg2: "starzejesz", Sg3: "starzeje",
		Pl1: "starzejemy", Pl2: "starzejecie", Pl3: "starzeją",
	},
	"gorzeć": {
		Sg1: "gorzeję", Sg2: "gorzejesz", Sg3: "gorzeje",
		Pl1: "gorzejemy", Pl2: "gorzejecie", Pl3: "gorzeją",
	},
	"dorzeć": {
		Sg1: "dorzeję", Sg2: "dorzejesz", Sg3: "dorzeje",
		Pl1: "dorzejemy", Pl2: "dorzejecie", Pl3: "dorzeją",
	},
	"dobrzeć": {
		Sg1: "dobrzeję", Sg2: "dobrzejesz", Sg3: "dobrzeje",
		Pl1: "dobrzejemy", Pl2: "dobrzejecie", Pl3: "dobrzeją",
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

	// czcić - needs szcz pattern
	"czcić": {
		Sg1: "czczę", Sg2: "czcisz", Sg3: "czci",
		Pl1: "czcimy", Pl2: "czcicie", Pl3: "czczą",
	},

	// kpić - no j-insertion (kpię not kpiję)
	"kpić": {
		Sg1: "kpię", Sg2: "kpisz", Sg3: "kpi",
		Pl1: "kpimy", Pl2: "kpicie", Pl3: "kpią",
	},

	// objąć - e-insertion (obejmę not objmę)
	"objąć": {
		Sg1: "obejmę", Sg2: "obejmiesz", Sg3: "obejmie",
		Pl1: "obejmiemy", Pl2: "obejmiecie", Pl3: "obejmą",
	},

	// ulec - gn-insertion (ulegnę not ulekę)
	"ulec": {
		Sg1: "ulegnę", Sg2: "ulegniesz", Sg3: "ulegnie",
		Pl1: "ulegniemy", Pl2: "ulegniecie", Pl3: "ulegną",
	},

	// wściec - kn-insertion (wścieknę not wściekę)
	"wściec": {
		Sg1: "wścieknę", Sg2: "wściekniesz", Sg3: "wścieknie",
		Pl1: "wściekniemy", Pl2: "wściekniecie", Pl3: "wściekną",
	},

	// dojrzeć - inchoative "to mature" (dojrzeję not dojrzę)
	"dojrzeć": {
		Sg1: "dojrzeję", Sg2: "dojrzejesz", Sg3: "dojrzeje",
		Pl1: "dojrzejemy", Pl2: "dojrzejecie", Pl3: "dojrzeją",
	},

	// boleć prefixed forms (base is homograph, handled separately)
	// bolę pattern (physical pain): poboleć, rozboleć, zaboleć
	"poboleć": {
		Sg1: "pobolę", Sg2: "pobolisz", Sg3: "poboli",
		Pl1: "pobolimy", Pl2: "pobolicie", Pl3: "pobolą",
	},
	"rozboleć": {
		Sg1: "rozbolę", Sg2: "rozbolisz", Sg3: "rozboli",
		Pl1: "rozbolimy", Pl2: "rozbolicie", Pl3: "rozbolą",
	},
	"zaboleć": {
		Sg1: "zabolę", Sg2: "zabolisz", Sg3: "zaboli",
		Pl1: "zabolimy", Pl2: "zabolicie", Pl3: "zabolą",
	},
	// boleję pattern (inchoative/emotional): oboleć, odboleć, przeboleć, współboleć, wyboleć
	"oboleć": {
		Sg1: "oboleję", Sg2: "obolejesz", Sg3: "oboleje",
		Pl1: "obolejemy", Pl2: "obolejecie", Pl3: "oboleją",
	},
	"odboleć": {
		Sg1: "odboleję", Sg2: "odbolejesz", Sg3: "odboleje",
		Pl1: "odbolejemy", Pl2: "odbolejecie", Pl3: "odboleją",
	},
	"przeboleć": {
		Sg1: "przeboleję", Sg2: "przebolejesz", Sg3: "przeboleje",
		Pl1: "przebolejemy", Pl2: "przebolejecie", Pl3: "przeboleją",
	},
	"współboleć": {
		Sg1: "współboleję", Sg2: "współbolejesz", Sg3: "współboleje",
		Pl1: "współbolejemy", Pl2: "współbolejecie", Pl3: "współboleją",
	},
	"wyboleć": {
		Sg1: "wyboleję", Sg2: "wybolejesz", Sg3: "wyboleje",
		Pl1: "wybolejemy", Pl2: "wybolejecie", Pl3: "wyboleją",
	},

	// swędzieć - action verb (swędzę not swędzieję)
	"swędzieć": {
		Sg1: "swędzę", Sg2: "swędzisz", Sg3: "swędzi",
		Pl1: "swędzimy", Pl2: "swędzicie", Pl3: "swędzą",
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

	// żreć - special -reć pattern (not -rem)
	"żreć": {
		Sg1: "żrę", Sg2: "żresz", Sg3: "żre",
		Pl1: "żremy", Pl2: "żrecie", Pl3: "żrą",
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

	// tłuc - suppletive stem tłuk/tłucz
	"tłuc": {
		Sg1: "tłukę", Sg2: "tłuczesz", Sg3: "tłucze",
		Pl1: "tłuczemy", Pl2: "tłuczecie", Pl3: "tłuką",
	},

	// pleść - suppletive stem plot/plec
	"pleść": {
		Sg1: "plotę", Sg2: "pleciesz", Sg3: "plecie",
		Pl1: "pleciemy", Pl2: "pleciecie", Pl3: "plotą",
	},

	// kląć - suppletive stem kln
	"kląć": {
		Sg1: "klnę", Sg2: "klniesz", Sg3: "klnie",
		Pl1: "klniemy", Pl2: "klniecie", Pl3: "klną",
	},

	// piąć - suppletive stem pn (with e-insertion for consonant clusters)
	"piąć": {
		Sg1: "pnę", Sg2: "pniesz", Sg3: "pnie",
		Pl1: "pniemy", Pl2: "pniecie", Pl3: "pną",
	},
	// Prefixed piąć verbs with e-insertion
	"wspiąć": {
		Sg1: "wespnę", Sg2: "wespniesz", Sg3: "wespnie",
		Pl1: "wespniemy", Pl2: "wespniecie", Pl3: "wespną",
	},
	"zapiąć": {
		Sg1: "zapnę", Sg2: "zapniesz", Sg3: "zapnie",
		Pl1: "zapniemy", Pl2: "zapniecie", Pl3: "zapną",
	},
	"przypiąć": {
		Sg1: "przypnę", Sg2: "przypniesz", Sg3: "przypnie",
		Pl1: "przypniemy", Pl2: "przypniecie", Pl3: "przypną",
	},
	"odpiąć": {
		Sg1: "odpnę", Sg2: "odpniesz", Sg3: "odpnie",
		Pl1: "odpniemy", Pl2: "odpniecie", Pl3: "odpną",
	},
	"dopiąć": {
		Sg1: "dopnę", Sg2: "dopniesz", Sg3: "dopnie",
		Pl1: "dopniemy", Pl2: "dopniecie", Pl3: "dopną",
	},
	"spiąć": {
		Sg1: "spnę", Sg2: "spniesz", Sg3: "spnie",
		Pl1: "spniemy", Pl2: "spniecie", Pl3: "spną",
	},
	"wpiąć": {
		Sg1: "wpnę", Sg2: "wpniesz", Sg3: "wpnie",
		Pl1: "wpniemy", Pl2: "wpniecie", Pl3: "wpną",
	},
	"napiąć": {
		Sg1: "napnę", Sg2: "napniesz", Sg3: "napnie",
		Pl1: "napniemy", Pl2: "napniecie", Pl3: "napną",
	},
	"rozpiąć": {
		Sg1: "rozpnę", Sg2: "rozpniesz", Sg3: "rozpnie",
		Pl1: "rozpniemy", Pl2: "rozpniecie", Pl3: "rozpną",
	},
	"wypiąć": {
		Sg1: "wypnę", Sg2: "wypniesz", Sg3: "wypnie",
		Pl1: "wypniemy", Pl2: "wypniecie", Pl3: "wypną",
	},

	// wiać - special pattern (wieję not wiam)
	"wiać": {
		Sg1: "wieję", Sg2: "wiejesz", Sg3: "wieje",
		Pl1: "wiejemy", Pl2: "wiejecie", Pl3: "wieją",
	},

	// chwiać - sway (chwieję not chwiam)
	"chwiać": {
		Sg1: "chwieję", Sg2: "chwiejesz", Sg3: "chwieje",
		Pl1: "chwiejemy", Pl2: "chwiejecie", Pl3: "chwieją",
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

	// tajać - minority -ajać with -ję pattern (taję not tajam)
	"tajać": {
		Sg1: "taję", Sg2: "tajesz", Sg3: "taje",
		Pl1: "tajemy", Pl2: "tajecie", Pl3: "tają",
	},

	// ćpać - slang for drug use (regular -am pattern)
	"ćpać": {
		Sg1: "ćpam", Sg2: "ćpasz", Sg3: "ćpa",
		Pl1: "ćpamy", Pl2: "ćpacie", Pl3: "ćpają",
	},

	// bimbać - regular -am (not alternating like other -bać)
	"bimbać": {
		Sg1: "bimbam", Sg2: "bimbasz", Sg3: "bimba",
		Pl1: "bimbamy", Pl2: "bimbacie", Pl3: "bimbają",
	},

	// gabać - regular -am (not alternating)
	"gabać": {
		Sg1: "gabam", Sg2: "gabasz", Sg3: "gaba",
		Pl1: "gabamy", Pl2: "gabacie", Pl3: "gabają",
	},

	// chybać - regular -am (pochybać, przychybać, wychybać)
	"chybać": {
		Sg1: "chybam", Sg2: "chybasz", Sg3: "chyba",
		Pl1: "chybamy", Pl2: "chybacie", Pl3: "chybają",
	},

	// gibać - regular -am
	"gibać": {
		Sg1: "gibam", Sg2: "gibasz", Sg3: "giba",
		Pl1: "gibamy", Pl2: "gibacie", Pl3: "gibają",
	},

	// siorbać - regular -am (to slurp)
	"siorbać": {
		Sg1: "siorbam", Sg2: "siorbasz", Sg3: "siorba",
		Pl1: "siorbamy", Pl2: "siorbacie", Pl3: "siorbają",
	},

	// stąpać - regular -am (to step)
	"stąpać": {
		Sg1: "stąpam", Sg2: "stąpasz", Sg3: "stąpa",
		Pl1: "stąpamy", Pl2: "stąpacie", Pl3: "stąpają",
	},

	// pchlać - related to pchła (flea), not chlać
	"pchlać": {
		Sg1: "pchlam", Sg2: "pchlasz", Sg3: "pchla",
		Pl1: "pchlamy", Pl2: "pchlacie", Pl3: "pchlają",
	},

	// rychlać - not related to chlać
	"rychlać": {
		Sg1: "rychlam", Sg2: "rychlasz", Sg3: "rychla",
		Pl1: "rychlamy", Pl2: "rychlacie", Pl3: "rychlają",
	},

	// gdybać - regular -am (to speculate)
	"gdybać": {
		Sg1: "gdybam", Sg2: "gdybasz", Sg3: "gdyba",
		Pl1: "gdybamy", Pl2: "gdybacie", Pl3: "gdybają",
	},

	// gnić - to rot (j-insertion: gniję)
	"gnić": {
		Sg1: "gniję", Sg2: "gnijesz", Sg3: "gnije",
		Pl1: "gnijemy", Pl2: "gnijecie", Pl3: "gniją",
	},

	// siać - to sow (ia → ie + ję)
	"siać": {
		Sg1: "sieję", Sg2: "siejesz", Sg3: "sieje",
		Pl1: "siejemy", Pl2: "siejecie", Pl3: "sieją",
	},

	// -tajać verbs meaning "to conceal" (from "taja"), use -tajam not -taję
	// Different from tajać meaning "to thaw" which uses taję
	"utajać": {
		Sg1: "utajam", Sg2: "utajasz", Sg3: "utaja",
		Pl1: "utajamy", Pl2: "utajacie", Pl3: "utajają",
	},
	"zatajać": {
		Sg1: "zatajam", Sg2: "zatajasz", Sg3: "zataja",
		Pl1: "zatajamy", Pl2: "zatajacie", Pl3: "zatajają",
	},
	"przytajać": {
		Sg1: "przytajam", Sg2: "przytajasz", Sg3: "przytaja",
		Pl1: "przytajamy", Pl2: "przytajacie", Pl3: "przytajają",
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
