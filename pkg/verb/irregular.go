package verb

// Conjugation classes for Polish present tense.
// Named after standard Polish linguistics conventions.
const (
	ConjI   byte = 'e' // Class I:   -ę, -esz, -e, -emy, -ecie, -ą
	ConjIIa byte = 'i' // Class IIa: -ę, -isz, -i, -imy, -icie, -ą
	ConjIIb byte = 'y' // Class IIb: -ę, -ysz, -y, -ymy, -ycie, -ą
	ConjIII byte = 'a' // Class III: -am, -asz, -a, -amy, -acie, -ają
	ConjIV  byte = 'E' // Class IV:  -em, -esz, -e, -emy, -ecie, -(j)ą
)

// presentSpec compactly describes a present tense paradigm.
// Uses stem + class byte, with optional per-cell overrides; expanded by build().
type presentSpec struct {
	stem  string // stem for Sg2/Sg3/Pl1/Pl2 (and Sg1/Pl3 if sg13 is empty)
	sg13  string // Sg1/Pl3 stem override (defaults to stem)
	class byte   // ConjI, ConjIIa, ConjIIb, ConjIII, or ConjIV
	sg1   string // complete Sg1 form override
	sg2   string // complete Sg2 form override
	sg3   string // complete Sg3 form override
	pl1   string // complete Pl1 form override
	pl2   string // complete Pl2 form override
	pl3   string // complete Pl3 form override
}

func (s presentSpec) build() PresentTense {
	sg13 := s.sg13
	if sg13 == "" {
		sg13 = s.stem
	}

	var pt PresentTense
	switch s.class {
	case ConjI:
		pt = PresentTense{
			Sg1: sg13 + "ę", Sg2: s.stem + "esz", Sg3: s.stem + "e",
			Pl1: s.stem + "emy", Pl2: s.stem + "ecie", Pl3: sg13 + "ą",
		}
	case ConjIIa:
		pt = PresentTense{
			Sg1: sg13 + "ę", Sg2: s.stem + "isz", Sg3: s.stem + "i",
			Pl1: s.stem + "imy", Pl2: s.stem + "icie", Pl3: sg13 + "ą",
		}
	case ConjIIb:
		pt = PresentTense{
			Sg1: sg13 + "ę", Sg2: s.stem + "ysz", Sg3: s.stem + "y",
			Pl1: s.stem + "ymy", Pl2: s.stem + "ycie", Pl3: sg13 + "ą",
		}
	case ConjIII:
		pt = PresentTense{
			Sg1: sg13 + "am", Sg2: s.stem + "asz", Sg3: s.stem + "a",
			Pl1: s.stem + "amy", Pl2: s.stem + "acie", Pl3: sg13 + "ają",
		}
	case ConjIV:
		pt = PresentTense{
			Sg1: s.stem + "m", Sg2: s.stem + "sz", Sg3: s.stem,
			Pl1: s.stem + "my", Pl2: s.stem + "cie", Pl3: sg13 + "ją",
		}
	}

	if s.sg1 != "" {
		pt.Sg1 = s.sg1
	}
	if s.sg2 != "" {
		pt.Sg2 = s.sg2
	}
	if s.sg3 != "" {
		pt.Sg3 = s.sg3
	}
	if s.pl1 != "" {
		pt.Pl1 = s.pl1
	}
	if s.pl2 != "" {
		pt.Pl2 = s.pl2
	}
	if s.pl3 != "" {
		pt.Pl3 = s.pl3
	}
	return pt
}

// homographs contains verbs with multiple valid paradigms (different meanings).
// These are checked first before irregular verbs.
var homographs = map[string][]Paradigm{
	// stać: "to stand" (imperfective) vs "to become/afford" (perfective)
	"stać": {
		{
			PresentTense: presentSpec{sg13: "stoj", stem: "sto", class: ConjIIa}.build(),
			Gloss:        "to stand",
		},
		{
			PresentTense: presentSpec{sg13: "stan", stem: "stani", class: ConjI}.build(),
			Gloss:        "to become, to afford",
		},
	},
	// słać: "to send" vs "to spread (bedding)"
	"słać": {
		{
			PresentTense: presentSpec{stem: "śl", class: ConjI}.build(),
			Gloss:        "to send",
		},
		{
			PresentTense: presentSpec{stem: "ściel", class: ConjI}.build(),
			Gloss:        "to spread (bedding)",
		},
	},
	// boleć: "physical pain" vs "to grieve/worry" (inchoative)
	"boleć": {
		{
			PresentTense: presentSpec{stem: "bol", class: ConjIIa}.build(),
			Gloss:        "to hurt (physical pain)",
		},
		{
			PresentTense: presentSpec{stem: "bolej", class: ConjI}.build(),
			Gloss:        "to grieve, to worry",
		},
	},
	// stajać: frequentative of stać (both patterns attested)
	"stajać": {
		{
			PresentTense: presentSpec{stem: "staj", class: ConjI}.build(),
			Gloss:        "to keep standing/stopping (frequentative)",
		},
		{
			PresentTense: presentSpec{stem: "staj", class: ConjIII}.build(),
			Gloss:        "to keep standing/stopping (variant)",
		},
	},
	// chlać: vulgar "to gulp" (both patterns attested)
	"chlać": {
		{
			PresentTense: presentSpec{stem: "chl", class: ConjIII}.build(),
			Gloss:        "to gulp/slurp (vulgar)",
		},
		{
			PresentTense: presentSpec{stem: "chlej", class: ConjI}.build(),
			Gloss:        "to gulp/slurp (variant)",
		},
	},
	// ziajać: both -am and -ę patterns attested
	"ziajać": {
		{
			PresentTense: presentSpec{stem: "ziaj", class: ConjIII}.build(),
			Gloss:        "to pant/gasp",
		},
		{
			PresentTense: presentSpec{stem: "ziaj", class: ConjI}.build(),
			Gloss:        "to pant/gasp (variant)",
		},
	},
	// bajać: both -am and -ę patterns attested
	"bajać": {
		{
			PresentTense: presentSpec{stem: "baj", class: ConjIII}.build(),
			Gloss:        "to tell fairy tales",
		},
		{
			PresentTense: presentSpec{stem: "baj", class: ConjI}.build(),
			Gloss:        "to tell fairy tales (variant)",
		},
	},
	// przytajać: both -am and -ę patterns attested
	"przytajać": {
		{
			PresentTense: presentSpec{stem: "przytaj", class: ConjIII}.build(),
			Gloss:        "to crouch/hide",
		},
		{
			PresentTense: presentSpec{stem: "przytaj", class: ConjI}.build(),
			Gloss:        "to crouch/hide (variant)",
		},
	},
	// kaszliwać: both -iwuję and -uję patterns attested
	"kaszliwać": {
		{
			PresentTense: presentSpec{stem: "kaszliwuj", class: ConjI}.build(),
			Gloss:        "to cough (frequentative)",
		},
		{
			PresentTense: presentSpec{stem: "kaszluj", class: ConjI}.build(),
			Gloss:        "to cough (frequentative, variant)",
		},
	},
	// połajać: both -am and -ę patterns attested (different from other łajać prefixes)
	"połajać": {
		{
			PresentTense: presentSpec{stem: "połaj", class: ConjIII}.build(),
			Gloss:        "to scold",
		},
		{
			PresentTense: presentSpec{stem: "połaj", class: ConjI}.build(),
			Gloss:        "to scold (variant)",
		},
	},
	// pyskiwać: both -iwuję and -uję patterns attested
	"pyskiwać": {
		{
			PresentTense: presentSpec{stem: "pyskiwuj", class: ConjI}.build(),
			Gloss:        "to talk back",
		},
		{
			PresentTense: presentSpec{stem: "pyskuj", class: ConjI}.build(),
			Gloss:        "to talk back (variant)",
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

	// These homographs support prefix expansion because their
	// prefixed forms retain both conjugation patterns.
	// stać prefixed forms like "dostać" are NOT homographs.
	expandableHomographs := map[string]bool{
		"słać": true, "chlać": true, "ziajać": true, "bajać": true,
		"przytajać": true, "kaszliwać": true, "pyskiwać": true,
	}
	for _, prefix := range verbPrefixes {
		if len(infinitive) > len(prefix) && infinitive[:len(prefix)] == prefix {
			base := infinitive[len(prefix):]
			// Expand homographs that support prefix expansion
			if expandableHomographs[base] {
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

// irregularPresentSpecs defines present tense stems compactly.
// Merged into irregularSpecs by buildIrregularSpecs() in spec.go.
var irregularPresentSpecs = map[string]presentSpec{
	// Suppletive / highly irregular
	"być": {
		class: ConjIV, // suppletive - all forms overridden
		sg1: "jestem", sg2: "jesteś", sg3: "jest",
		pl1: "jesteśmy", pl2: "jesteście", pl3: "są",
	},
	"mieć":   {stem: "m", class: ConjIII},
	"umieć":  {stem: "umie", class: ConjIV},
	"musieć": {sg13: "musz", stem: "mus", class: ConjIIa},
	"jeść":   {stem: "je", sg13: "jedz", class: ConjIV},
	"dać":    {stem: "d", class: ConjIII, pl3: "dadzą"},
	"sprzedać": {stem: "sprzed", class: ConjIII, pl3: "sprzedadzą"},
	"wziąć":  {sg13: "wezm", stem: "weźmi", class: ConjI},
	"ciąć":   {sg13: "tn", stem: "tni", class: ConjI},
	"iść":    {sg13: "id", stem: "idzi", class: ConjI},

	// Stem-changing -ać verbs (jechać type)
	"jechać": {sg13: "jad", stem: "jedzi", class: ConjI},

	// Stem-changing -ać verbs (brać type: a→io/ie)
	"brać": {sg13: "bior", stem: "bierz", class: ConjI},
	"prać": {sg13: "pior", stem: "pierz", class: ConjI},

	// Minority -sać verbs that alternate (s→sz)
	"pisać":   {stem: "pisz", class: ConjI},
	"czesać":  {stem: "czesz", class: ConjI},
	"kołysać": {stem: "kołysz", class: ConjI},

	// -sać verbs that stay regular -am
	"kasać":   {stem: "kas", class: ConjIII},
	"ciosać":  {stem: "cios", class: ConjIII},
	"ciesać":  {stem: "cies", class: ConjIII},
	"krzesać": {stem: "krzes", class: ConjIII},

	// Minority -kać verbs that alternate (k→cz)
	"skakać":  {stem: "skacz", class: ConjI},
	"płakać":  {stem: "płacz", class: ConjI},

	// Minority -zać verbs that alternate (z→ż)
	"wiązać":  {stem: "wiąż", class: ConjI},
	"kazać":   {stem: "każ", class: ConjI},
	"mazać":   {stem: "maż", class: ConjI},
	"lizać":   {stem: "liż", class: ConjI},
	"okazać":  {stem: "okaż", class: ConjI},

	// Minority -rać that alternates (r→rz)
	"karać":   {stem: "karz", class: ConjI},

	// naleźć - suppletive stem najd-
	"naleźć":  {sg13: "najd", stem: "najdzi", class: ConjI},

	// spać - stem change s→ś before palatal
	"spać":    {sg13: "śpi", stem: "śp", class: ConjIIa},

	// bać się - suppletive stem boj-
	"bać":     {sg13: "boj", stem: "bo", class: ConjIIa},

	// dziać się - special -ać → -eję pattern
	"dziać":   {stem: "dziej", class: ConjI},

	// podobać się - regular -am
	"podobać": {stem: "podob", class: ConjIII},

	// Monosyllabic -ić/-yć verbs with j-insertion
	"bić":   {stem: "bij", class: ConjI},
	"lić":   {stem: "lij", class: ConjI},
	"pić":   {stem: "pij", class: ConjI},
	"żyć":   {stem: "żyj", class: ConjI},
	"myć":   {stem: "myj", class: ConjI},
	"ryć":   {stem: "ryj", class: ConjI},
	"szyć":  {stem: "szyj", class: ConjI},
	"wyć":   {stem: "wyj", class: ConjI},
	"kryć":  {stem: "kryj", class: ConjI},
	"wić":   {stem: "wij", class: ConjI},
	"gnić":  {stem: "gnij", class: ConjI},

	// Prefixed -żyć with compound prefixes
	"użyć":          {stem: "użyj", class: ConjI},
	"spożyć":        {stem: "spożyj", class: ConjI},
	"współżyć":      {stem: "współżyj", class: ConjI},
	"współprzeżyć":  {stem: "współprzeżyj", class: ConjI},

	// Prefixed wić forms that don't match simple prefix patterns
	"opowić":    {stem: "opowij", class: ConjI},
	"rozpowić":  {stem: "rozpowij", class: ConjI},
	"spowić":    {stem: "spowij", class: ConjI},
	"upowić":    {stem: "upowij", class: ConjI},

	// -pomnieć verbs
	"pomnieć":   {stem: "pomn", class: ConjIIa},

	// -mrzeć verbs (rz→r in sg1/pl3)
	"mrzeć":     {sg13: "mr", stem: "mrz", class: ConjI},

	// ciec verbs: k-insertion
	"ciec":      {sg13: "ciekn", stem: "ciekni", class: ConjI},

	// woleć
	"woleć":     {stem: "wol", class: ConjIIa},

	// -jąć verbs: suppletive stem -jm-
	"jąć":    {sg13: "jm", stem: "jmi", class: ConjI},
	"zdjąć":  {sg13: "zdejm", stem: "zdejmi", class: ConjI},
	"podjąć": {sg13: "podejm", stem: "podejmi", class: ConjI},
	"odjąć":  {sg13: "odejm", stem: "odejmi", class: ConjI},
	"objąć":  {sg13: "obejm", stem: "obejmi", class: ConjI},
	"nająć":  {sg13: "najm", stem: "najmi", class: ConjI},

	// -cząć verbs: czn- stem
	"cząć":      {sg13: "czn", stem: "czni", class: ConjI},
	"począć":    {sg13: "poczn", stem: "poczni", class: ConjI},
	"odpocząć":  {sg13: "odpoczn", stem: "odpoczni", class: ConjI},
	"rozpocząć": {sg13: "rozpoczn", stem: "rozpoczni", class: ConjI},
	"spocząć":   {sg13: "spoczn", stem: "spoczni", class: ConjI},
	"wypocząć":  {sg13: "wypoczn", stem: "wypoczni", class: ConjI},
	"wszcząć":   {sg13: "wszczn", stem: "wszczni", class: ConjI},
	"poczęć":    {sg13: "poczn", stem: "poczni", class: ConjI},

	// Action verb -mieć patterns (grzmieć → grzmię, not grzmieję)
	"grzmieć":  {sg13: "grzmi", stem: "grzm", class: ConjIIa},
	"szumieć":  {sg13: "szumi", stem: "szum", class: ConjIIa},
	"tłumieć":  {sg13: "tłumi", stem: "tłum", class: ConjIIa},

	// patrzeć - action verb (class y)
	"patrzeć":  {stem: "patrz", class: ConjIIb},

	// Inchoative -rzeć/-eć verbs (use -eję pattern)
	"starzeć":  {stem: "starzej", class: ConjI},
	"gorzeć":   {stem: "gorzej", class: ConjI},
	"dorzeć":   {stem: "dorzej", class: ConjI},
	"dobrzeć":  {stem: "dobrzej", class: ConjI},
	"dojrzeć":  {stem: "dojrzej", class: ConjI},
	"doźrzeć":  {stem: "doźrzej", class: ConjI},
	"przejrzeć": {stem: "przejrzej", class: ConjI},

	// -rwać verbs: -ę/-ie pattern
	"rwać":      {sg13: "rw", stem: "rwi", class: ConjI},

	// -zwać verbs: -ę/-ie pattern
	"zwać":      {sg13: "zw", stem: "zwi", class: ConjI},

	// dbać - regular -am
	"dbać":      {stem: "db", class: ConjIII},

	// śmiać się - special pattern (ać → eję)
	"śmiać":     {stem: "śmiej", class: ConjI},

	// -ieć action verbs
	"cierpieć":   {sg13: "cierpi", stem: "cierp", class: ConjIIa},
	"wisieć":     {sg13: "wisz", stem: "wis", class: ConjIIa},
	"tkwieć":     {sg13: "tkwi", stem: "tkw", class: ConjIIa},
	"śmierdzieć": {stem: "śmierdz", class: ConjIIa},
	"swędzieć":   {stem: "swędz", class: ConjIIa},
	"pierdzieć":  {stem: "pierdz", class: ConjIIa},
	"skomleć":    {stem: "skoml", class: ConjIIa},

	// jeździć - correct softening źdź → żdż
	"jeździć":    {sg13: "jeżdż", stem: "jeźdź", class: ConjIIa},

	// -nieć action verbs
	"pachnieć":   {sg13: "pachn", stem: "pachni", class: ConjI},

	// -strzec verbs: c→g alternation
	"strzec":     {sg13: "strzeg", stem: "strzeż", class: ConjI},

	// -chować verbs: -owam
	"chować":     {stem: "chow", class: ConjIII},

	// -kraść verbs: suppletive kradn- stem
	"kraść":      {sg13: "kradn", stem: "kradni", class: ConjI},

	// -kłaść verbs: suppletive kład- stem
	"kłaść":      {sg13: "kład", stem: "kładzi", class: ConjI},

	// uczcić/czcić - needs szcz
	"uczcić":     {sg13: "uczcz", stem: "uczc", class: ConjIIa},
	"czcić":      {sg13: "czcz", stem: "czc", class: ConjIIa},

	// kpić - no j-insertion
	"kpić":       {sg13: "kpi", stem: "kp", class: ConjIIa},

	// ulec - gn-insertion
	"ulec":       {sg13: "ulegn", stem: "ulegni", class: ConjI},

	// wściec - kn-insertion
	"wściec":     {sg13: "wściekn", stem: "wściekni", class: ConjI},

	// boleć prefixed forms (base is homograph)
	// bolę pattern (physical pain)
	"poboleć":  {stem: "pobol", class: ConjIIa},
	"rozboleć": {stem: "rozbol", class: ConjIIa},
	"zaboleć":  {stem: "zabol", class: ConjIIa},
	// boleję pattern (inchoative/emotional)
	"oboleć":      {stem: "obolej", class: ConjI},
	"odboleć":     {stem: "odbolej", class: ConjI},
	"przeboleć":   {stem: "przebolej", class: ConjI},
	"współboleć":  {stem: "współbolej", class: ConjI},
	"wyboleć":     {stem: "wybolej", class: ConjI},

	// wspomnieć - special prefix form
	"wspomnieć":   {stem: "wspomn", class: ConjIIa},

	// opisać - minority alternating -sać
	"opisać":      {stem: "opisz", class: ConjI},

	// wskazać - minority alternating -zać
	"wskazać":     {stem: "wskaż", class: ConjI},

	// brać prefix verbs with vowel elision
	"odebrać":  {sg13: "odbior", stem: "odbierz", class: ConjI},
	"zebrać":   {sg13: "zbior", stem: "zbierz", class: ConjI},
	"rozebrać": {sg13: "rozbior", stem: "rozbierz", class: ConjI},

	// lać verbs (j-insertion)
	"lać":       {stem: "lej", class: ConjI},

	// grześć - suppletive grzeb- stem
	"grześć":    {sg13: "grzeb", stem: "grzebi", class: ConjI},

	// żreć
	"żreć":      {stem: "żr", class: ConjI},

	// -przeć/-wrzeć verbs
	"przeć":     {sg13: "pr", stem: "prz", class: ConjI},
	"wrzeć":     {sg13: "wr", stem: "wrz", class: ConjI},

	// śnić - no j-insertion
	"śnić":      {sg13: "śni", stem: "śn", class: ConjIIa},

	// rzec - k-insertion
	"rzec":      {sg13: "rzekn", stem: "rzecz", class: ConjI},

	// tłuc - suppletive stem tłuk/tłucz
	"tłuc":      {sg13: "tłuk", stem: "tłucz", class: ConjI},

	// pleść - suppletive stem plot/plec
	"pleść":     {sg13: "plot", stem: "pleci", class: ConjI},

	// kląć - suppletive stem kln
	"kląć":      {sg13: "kln", stem: "klni", class: ConjI},

	// piąć - suppletive stem pn
	"piąć":      {sg13: "pn", stem: "pni", class: ConjI},
	"wspiąć":    {sg13: "wespn", stem: "wespni", class: ConjI},
	"zapiąć":    {sg13: "zapn", stem: "zapni", class: ConjI},
	"przypiąć":  {sg13: "przypn", stem: "przypni", class: ConjI},
	"odpiąć":    {sg13: "odpn", stem: "odpni", class: ConjI},
	"dopiąć":    {sg13: "dopn", stem: "dopni", class: ConjI},
	"spiąć":     {sg13: "spn", stem: "spni", class: ConjI},
	"wpiąć":     {sg13: "wpn", stem: "wpni", class: ConjI},
	"napiąć":    {sg13: "napn", stem: "napni", class: ConjI},
	"rozpiąć":   {sg13: "rozpn", stem: "rozpni", class: ConjI},
	"wypiąć":    {sg13: "wypn", stem: "wypni", class: ConjI},

	// wiać - special pattern (wieję)
	"wiać":      {stem: "wiej", class: ConjI},

	// chwiać - sway (chwieję)
	"chwiać":    {stem: "chwiej", class: ConjI},

	// krajać - j-insertion (kraję not krajam)
	"krajać":    {stem: "kraj", class: ConjI},

	// tajać - minority -ajać with -ję pattern
	"tajać":     {stem: "taj", class: ConjI},

	// ćpać - regular -am
	"ćpać":      {stem: "ćp", class: ConjIII},

	// Regular -am -bać verbs (not alternating)
	"bimbać":    {stem: "bimb", class: ConjIII},
	"gabać":     {stem: "gab", class: ConjIII},
	"chybać":    {stem: "chyb", class: ConjIII},
	"gibać":     {stem: "gib", class: ConjIII},
	"gdybać":    {stem: "gdyb", class: ConjIII},
	"zaniedbać": {stem: "zaniedb", class: ConjIII},

	// Regular -am misc verbs
	"siorbać":   {stem: "siorb", class: ConjIII},
	"stąpać":    {stem: "stąp", class: ConjIII},
	"pchlać":    {stem: "pchl", class: ConjIII},
	"rychlać":   {stem: "rychl", class: ConjIII},
	"kpać":      {stem: "kp", class: ConjIII},
	"kasłać":    {stem: "kasł", class: ConjIII},
	"cierpać":   {stem: "cierp", class: ConjIII},
	"siąpać":    {stem: "siąp", class: ConjIII},
	"tyrpać":    {stem: "tyrp", class: ConjIII},
	"ściubać":   {stem: "ściub", class: ConjIII},
	"ślipać":    {stem: "ślip", class: ConjIII},
	"bombać":    {stem: "bomb", class: ConjIII},

	// Inchoative -eć verbs (use -eję pattern)
	"chorzeć":     {stem: "chorzej", class: ConjI},
	"tężeć":       {stem: "tężej", class: ConjI},
	"dumieć":      {stem: "dumiej", class: ConjI},
	"goreć":       {stem: "gorej", class: ConjI},
	"śniedzieć":   {stem: "śniedziej", class: ConjI},
	"srebrzeć":    {stem: "srebrzej", class: ConjI},
	"cukrzeć":     {stem: "cukrzej", class: ConjI},
	"dorośleć":    {stem: "doroślej", class: ConjI},
	"wydorośleć":  {stem: "wydoroślej", class: ConjI},
	"zelżeć":      {stem: "zelżej", class: ConjI},
	"wilżeć":      {stem: "wilżej", class: ConjI},
	"wężeć":       {stem: "wężej", class: ConjI},
	"rzedzieć":    {stem: "rzedziej", class: ConjI},
	"sfolżeć":     {stem: "sfolżej", class: ConjI},
	"szlachcieć":  {stem: "szlachciej", class: ConjI},
	"ochujeć":     {stem: "ochujej", class: ConjI},
	"ociężeć":     {stem: "ociężej", class: ConjI},
	"ściężeć":     {stem: "ściężej", class: ConjI},
	"oszedzieć":   {stem: "oszedziej", class: ConjI},
	"szedzieć":    {stem: "szedziej", class: ConjI},
	"sposążeć":    {stem: "sposążej", class: ConjI},
	"wyryżeć":     {stem: "wyryżej", class: ConjI},

	// siać - to sow (ia → ie + ję)
	"siać":        {stem: "siej", class: ConjI},

	// -tajać verbs meaning "to conceal" (use -tajam not -taję)
	"utajać":      {stem: "utaj", class: ConjIII},
	"zatajać":     {stem: "zataj", class: ConjIII},

	// łajać - minority -ajać that uses -ę/-esz
	"łajać":       {stem: "łaj", class: ConjI},

	// knajać - uses -ę/-esz
	"knajać":      {stem: "knaj", class: ConjI},

	// -iwać verbs that use -iwuję pattern
	"strzeliwać": {stem: "strzeliwuj", class: ConjI},
	"myśliwać":   {stem: "myśliwuj", class: ConjI},
	"boliwać":    {stem: "boliwuj", class: ConjI},
	"mgliwać":    {stem: "mgliwuj", class: ConjI},
	"skuliwać":   {stem: "skuliwuj", class: ConjI},

	// -ić verbs with no j-insertion
	"tlić":        {stem: "tl", class: ConjIIa},
	"clić":        {stem: "cl", class: ConjIIa},
	"dlić":        {stem: "dl", class: ConjIIa},

	// -ywać verbs that use -uję pattern
	"mieszywać":   {stem: "mieszuj", class: ConjI},
	"supływać":    {stem: "supłuj", class: ConjI},
	"bazgrywać":   {stem: "bazgruj", class: ConjI},
	"podobywać":   {stem: "podobuj", class: ConjI},

	// dziamdziać - uses -am pattern
	"dziamdziać":  {stem: "dziamdzi", class: ConjIII},

	// -piać verbs: uses -eję pattern
	"piać":        {stem: "piej", class: ConjI},
	"spiać":       {stem: "spiej", class: ConjI},
	"dośpiać":     {stem: "dośpiej", class: ConjI},
	"uśpiać":      {stem: "uśpiej", class: ConjI},

	// pomieć - uses -am pattern
	"pomieć":      {stem: "pom", class: ConjIII},

	// sposzyć - j-insertion
	"sposzyć":     {stem: "sposzyj", class: ConjI},

	// źreć/źrzeć - inchoative
	"źreć":        {stem: "źrej", class: ConjI},
	"źrzeć":       {stem: "źrzej", class: ConjI},

	// oziać - uses -eję pattern
	"oziać":       {stem: "oziej", class: ConjI},

	// Compound prefix -strzeliwać verbs
	"porozstrzeliwać": {stem: "porozstrzeliwuj", class: ConjI},
	"powystrzeliwać":  {stem: "powystrzeliwuj", class: ConjI},
}

// Common prefixes in Polish
var verbPrefixes = []string{
	"prze", "przy", "roz", "roze", "wy", "za", "na", "po", "do", "od", "ode", "ob", "obe",
	"pod", "pode", "nad", "nade", "wz", "wze", "u", "s", "z", "ze", "w", "we", "o",
}
