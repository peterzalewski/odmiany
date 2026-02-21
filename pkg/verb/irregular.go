package verb

// Conjugation classes for Polish present tense.
const (
	ConjE byte = 'e' // -ę, -esz, -e, -emy, -ecie, -ą
	ConjI byte = 'i' // -ę, -isz, -i, -imy, -icie, -ą
	ConjY byte = 'y' // -ę, -ysz, -y, -ymy, -ycie, -ą
	ConjA byte = 'a' // -am, -asz, -a, -amy, -acie, -ają
)

// presSpec compactly describes a present tense paradigm.
// At most 2 stems + a class byte; expanded to full PresentTense by build().
type presSpec struct {
	stem  string // stem for Sg2/Sg3/Pl1/Pl2 (and Sg1/Pl3 if sg13 is empty)
	sg13  string // Sg1/Pl3 stem override (defaults to stem)
	class byte   // ConjE, ConjI, ConjY, or ConjA
	sg1   string // complete Sg1 form override
	pl3   string // complete Pl3 form override
}

func (s presSpec) build() PresentTense {
	sg13 := s.sg13
	if sg13 == "" {
		sg13 = s.stem
	}

	var pt PresentTense
	switch s.class {
	case ConjE:
		pt = PresentTense{
			Sg1: sg13 + "ę", Sg2: s.stem + "esz", Sg3: s.stem + "e",
			Pl1: s.stem + "emy", Pl2: s.stem + "ecie", Pl3: sg13 + "ą",
		}
	case ConjI:
		pt = PresentTense{
			Sg1: sg13 + "ę", Sg2: s.stem + "isz", Sg3: s.stem + "i",
			Pl1: s.stem + "imy", Pl2: s.stem + "icie", Pl3: sg13 + "ą",
		}
	case ConjY:
		pt = PresentTense{
			Sg1: sg13 + "ę", Sg2: s.stem + "ysz", Sg3: s.stem + "y",
			Pl1: s.stem + "ymy", Pl2: s.stem + "ycie", Pl3: sg13 + "ą",
		}
	case ConjA:
		pt = PresentTense{
			Sg1: sg13 + "am", Sg2: s.stem + "asz", Sg3: s.stem + "a",
			Pl1: s.stem + "amy", Pl2: s.stem + "acie", Pl3: sg13 + "ają",
		}
	}

	if s.sg1 != "" {
		pt.Sg1 = s.sg1
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
			PresentTense: presSpec{sg13: "stoj", stem: "sto", class: ConjI}.build(),
			Gloss:        "to stand",
		},
		{
			PresentTense: presSpec{sg13: "stan", stem: "stani", class: ConjE}.build(),
			Gloss:        "to become, to afford",
		},
	},
	// słać: "to send" vs "to spread (bedding)"
	"słać": {
		{
			PresentTense: presSpec{stem: "śl", class: ConjE}.build(),
			Gloss:        "to send",
		},
		{
			PresentTense: presSpec{stem: "ściel", class: ConjE}.build(),
			Gloss:        "to spread (bedding)",
		},
	},
	// boleć: "physical pain" vs "to grieve/worry" (inchoative)
	"boleć": {
		{
			PresentTense: presSpec{stem: "bol", class: ConjI}.build(),
			Gloss:        "to hurt (physical pain)",
		},
		{
			PresentTense: presSpec{stem: "bolej", class: ConjE}.build(),
			Gloss:        "to grieve, to worry",
		},
	},
	// stajać: frequentative of stać (both patterns attested)
	"stajać": {
		{
			PresentTense: presSpec{stem: "staj", class: ConjE}.build(),
			Gloss:        "to keep standing/stopping (frequentative)",
		},
		{
			PresentTense: presSpec{stem: "staj", class: ConjA}.build(),
			Gloss:        "to keep standing/stopping (variant)",
		},
	},
	// chlać: vulgar "to gulp" (both patterns attested)
	"chlać": {
		{
			PresentTense: presSpec{stem: "chl", class: ConjA}.build(),
			Gloss:        "to gulp/slurp (vulgar)",
		},
		{
			PresentTense: presSpec{stem: "chlej", class: ConjE}.build(),
			Gloss:        "to gulp/slurp (variant)",
		},
	},
	// ziajać: both -am and -ę patterns attested
	"ziajać": {
		{
			PresentTense: presSpec{stem: "ziaj", class: ConjA}.build(),
			Gloss:        "to pant/gasp",
		},
		{
			PresentTense: presSpec{stem: "ziaj", class: ConjE}.build(),
			Gloss:        "to pant/gasp (variant)",
		},
	},
	// bajać: both -am and -ę patterns attested
	"bajać": {
		{
			PresentTense: presSpec{stem: "baj", class: ConjA}.build(),
			Gloss:        "to tell fairy tales",
		},
		{
			PresentTense: presSpec{stem: "baj", class: ConjE}.build(),
			Gloss:        "to tell fairy tales (variant)",
		},
	},
	// przytajać: both -am and -ę patterns attested
	"przytajać": {
		{
			PresentTense: presSpec{stem: "przytaj", class: ConjA}.build(),
			Gloss:        "to crouch/hide",
		},
		{
			PresentTense: presSpec{stem: "przytaj", class: ConjE}.build(),
			Gloss:        "to crouch/hide (variant)",
		},
	},
	// kaszliwać: both -iwuję and -uję patterns attested
	"kaszliwać": {
		{
			PresentTense: presSpec{stem: "kaszliwuj", class: ConjE}.build(),
			Gloss:        "to cough (frequentative)",
		},
		{
			PresentTense: presSpec{stem: "kaszluj", class: ConjE}.build(),
			Gloss:        "to cough (frequentative, variant)",
		},
	},
	// połajać: both -am and -ę patterns attested (different from other łajać prefixes)
	"połajać": {
		{
			PresentTense: presSpec{stem: "połaj", class: ConjA}.build(),
			Gloss:        "to scold",
		},
		{
			PresentTense: presSpec{stem: "połaj", class: ConjE}.build(),
			Gloss:        "to scold (variant)",
		},
	},
	// pyskiwać: both -iwuję and -uję patterns attested
	"pyskiwać": {
		{
			PresentTense: presSpec{stem: "pyskiwuj", class: ConjE}.build(),
			Gloss:        "to talk back",
		},
		{
			PresentTense: presSpec{stem: "pyskuj", class: ConjE}.build(),
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

// irregularPresSpecs defines present tense stems compactly.
// Expanded to full PresentTense paradigms by init().
var irregularPresSpecs = map[string]presSpec{
	// Suppletive / highly irregular
	"mieć":   {stem: "m", class: ConjA},
	"musieć": {sg13: "musz", stem: "mus", class: ConjI},
	"jeść":   {stem: "j", sg13: "jedz", class: ConjE, sg1: "jem"},
	"dać":    {stem: "d", class: ConjA, pl3: "dadzą"},
	"sprzedać": {stem: "sprzed", class: ConjA, pl3: "sprzedadzą"},
	"wziąć":  {sg13: "wezm", stem: "weźmi", class: ConjE},
	"ciąć":   {sg13: "tn", stem: "tni", class: ConjE},
	"iść":    {sg13: "id", stem: "idzi", class: ConjE},

	// Stem-changing -ać verbs (jechać type)
	"jechać": {sg13: "jad", stem: "jedzi", class: ConjE},

	// Stem-changing -ać verbs (brać type: a→io/ie)
	"brać": {sg13: "bior", stem: "bierz", class: ConjE},
	"prać": {sg13: "pior", stem: "pierz", class: ConjE},

	// Minority -sać verbs that alternate (s→sz)
	"pisać":   {stem: "pisz", class: ConjE},
	"czesać":  {stem: "czesz", class: ConjE},
	"kołysać": {stem: "kołysz", class: ConjE},

	// -sać verbs that stay regular -am
	"kasać":   {stem: "kas", class: ConjA},
	"ciosać":  {stem: "cios", class: ConjA},
	"ciesać":  {stem: "cies", class: ConjA},
	"krzesać": {stem: "krzes", class: ConjA},

	// Minority -kać verbs that alternate (k→cz)
	"skakać":  {stem: "skacz", class: ConjE},
	"płakać":  {stem: "płacz", class: ConjE},

	// Minority -zać verbs that alternate (z→ż)
	"wiązać":  {stem: "wiąż", class: ConjE},
	"kazać":   {stem: "każ", class: ConjE},
	"mazać":   {stem: "maż", class: ConjE},
	"lizać":   {stem: "liż", class: ConjE},
	"okazać":  {stem: "okaż", class: ConjE},

	// Minority -rać that alternates (r→rz)
	"karać":   {stem: "karz", class: ConjE},

	// naleźć - suppletive stem najd-
	"naleźć":  {sg13: "najd", stem: "najdzi", class: ConjE},

	// spać - stem change s→ś before palatal
	"spać":    {sg13: "śpi", stem: "śp", class: ConjI},

	// bać się - suppletive stem boj-
	"bać":     {sg13: "boj", stem: "bo", class: ConjI},

	// dziać się - special -ać → -eję pattern
	"dziać":   {stem: "dziej", class: ConjE},

	// podobać się - regular -am
	"podobać": {stem: "podob", class: ConjA},

	// Monosyllabic -ić/-yć verbs with j-insertion
	"bić":   {stem: "bij", class: ConjE},
	"lić":   {stem: "lij", class: ConjE},
	"pić":   {stem: "pij", class: ConjE},
	"żyć":   {stem: "żyj", class: ConjE},
	"myć":   {stem: "myj", class: ConjE},
	"ryć":   {stem: "ryj", class: ConjE},
	"szyć":  {stem: "szyj", class: ConjE},
	"wyć":   {stem: "wyj", class: ConjE},
	"kryć":  {stem: "kryj", class: ConjE},
	"wić":   {stem: "wij", class: ConjE},
	"gnić":  {stem: "gnij", class: ConjE},

	// Prefixed -żyć with compound prefixes
	"użyć":          {stem: "użyj", class: ConjE},
	"spożyć":        {stem: "spożyj", class: ConjE},
	"współżyć":      {stem: "współżyj", class: ConjE},
	"współprzeżyć":  {stem: "współprzeżyj", class: ConjE},

	// Prefixed wić forms that don't match simple prefix patterns
	"opowić":    {stem: "opowij", class: ConjE},
	"rozpowić":  {stem: "rozpowij", class: ConjE},
	"spowić":    {stem: "spowij", class: ConjE},
	"upowić":    {stem: "upowij", class: ConjE},

	// -pomnieć verbs
	"pomnieć":   {stem: "pomn", class: ConjI},

	// -mrzeć verbs (rz→r in sg1/pl3)
	"mrzeć":     {sg13: "mr", stem: "mrz", class: ConjE},

	// ciec verbs: k-insertion
	"ciec":      {sg13: "ciekn", stem: "ciekni", class: ConjE},

	// woleć
	"woleć":     {stem: "wol", class: ConjI},

	// -jąć verbs: suppletive stem -jm-
	"jąć":    {sg13: "jm", stem: "jmi", class: ConjE},
	"zdjąć":  {sg13: "zdejm", stem: "zdejmi", class: ConjE},
	"podjąć": {sg13: "podejm", stem: "podejmi", class: ConjE},
	"odjąć":  {sg13: "odejm", stem: "odejmi", class: ConjE},
	"objąć":  {sg13: "obejm", stem: "obejmi", class: ConjE},
	"nająć":  {sg13: "najm", stem: "najmi", class: ConjE},

	// -cząć verbs: czn- stem
	"cząć":      {sg13: "czn", stem: "czni", class: ConjE},
	"począć":    {sg13: "poczn", stem: "poczni", class: ConjE},
	"odpocząć":  {sg13: "odpoczn", stem: "odpoczni", class: ConjE},
	"rozpocząć": {sg13: "rozpoczn", stem: "rozpoczni", class: ConjE},
	"spocząć":   {sg13: "spoczn", stem: "spoczni", class: ConjE},
	"wypocząć":  {sg13: "wypoczn", stem: "wypoczni", class: ConjE},
	"wszcząć":   {sg13: "wszczn", stem: "wszczni", class: ConjE},
	"poczęć":    {sg13: "poczn", stem: "poczni", class: ConjE},

	// Action verb -mieć patterns (grzmieć → grzmię, not grzmieję)
	"grzmieć":  {sg13: "grzmi", stem: "grzm", class: ConjI},
	"szumieć":  {sg13: "szumi", stem: "szum", class: ConjI},
	"tłumieć":  {sg13: "tłumi", stem: "tłum", class: ConjI},

	// patrzeć - action verb (class y)
	"patrzeć":  {stem: "patrz", class: ConjY},

	// Inchoative -rzeć/-eć verbs (use -eję pattern)
	"starzeć":  {stem: "starzej", class: ConjE},
	"gorzeć":   {stem: "gorzej", class: ConjE},
	"dorzeć":   {stem: "dorzej", class: ConjE},
	"dobrzeć":  {stem: "dobrzej", class: ConjE},
	"dojrzeć":  {stem: "dojrzej", class: ConjE},
	"doźrzeć":  {stem: "doźrzej", class: ConjE},
	"przejrzeć": {stem: "przejrzej", class: ConjE},

	// -rwać verbs: -ę/-ie pattern
	"rwać":      {sg13: "rw", stem: "rwi", class: ConjE},

	// -zwać verbs: -ę/-ie pattern
	"zwać":      {sg13: "zw", stem: "zwi", class: ConjE},

	// dbać - regular -am
	"dbać":      {stem: "db", class: ConjA},

	// śmiać się - special pattern (ać → eję)
	"śmiać":     {stem: "śmiej", class: ConjE},

	// -ieć action verbs
	"cierpieć":   {sg13: "cierpi", stem: "cierp", class: ConjI},
	"wisieć":     {sg13: "wisz", stem: "wis", class: ConjI},
	"tkwieć":     {sg13: "tkwi", stem: "tkw", class: ConjI},
	"śmierdzieć": {stem: "śmierdz", class: ConjI},
	"swędzieć":   {stem: "swędz", class: ConjI},
	"pierdzieć":  {stem: "pierdz", class: ConjI},
	"skomleć":    {stem: "skoml", class: ConjI},

	// jeździć - correct softening źdź → żdż
	"jeździć":    {sg13: "jeżdż", stem: "jeźdź", class: ConjI},

	// -nieć action verbs
	"pachnieć":   {sg13: "pachn", stem: "pachni", class: ConjE},

	// -strzec verbs: c→g alternation
	"strzec":     {sg13: "strzeg", stem: "strzeż", class: ConjE},

	// -chować verbs: -owam
	"chować":     {stem: "chow", class: ConjA},

	// -kraść verbs: suppletive kradn- stem
	"kraść":      {sg13: "kradn", stem: "kradni", class: ConjE},

	// -kłaść verbs: suppletive kład- stem
	"kłaść":      {sg13: "kład", stem: "kładzi", class: ConjE},

	// uczcić/czcić - needs szcz
	"uczcić":     {sg13: "uczcz", stem: "uczc", class: ConjI},
	"czcić":      {sg13: "czcz", stem: "czc", class: ConjI},

	// kpić - no j-insertion
	"kpić":       {sg13: "kpi", stem: "kp", class: ConjI},

	// ulec - gn-insertion
	"ulec":       {sg13: "ulegn", stem: "ulegni", class: ConjE},

	// wściec - kn-insertion
	"wściec":     {sg13: "wściekn", stem: "wściekni", class: ConjE},

	// boleć prefixed forms (base is homograph)
	// bolę pattern (physical pain)
	"poboleć":  {stem: "pobol", class: ConjI},
	"rozboleć": {stem: "rozbol", class: ConjI},
	"zaboleć":  {stem: "zabol", class: ConjI},
	// boleję pattern (inchoative/emotional)
	"oboleć":      {stem: "obolej", class: ConjE},
	"odboleć":     {stem: "odbolej", class: ConjE},
	"przeboleć":   {stem: "przebolej", class: ConjE},
	"współboleć":  {stem: "współbolej", class: ConjE},
	"wyboleć":     {stem: "wybolej", class: ConjE},

	// wspomnieć - special prefix form
	"wspomnieć":   {stem: "wspomn", class: ConjI},

	// opisać - minority alternating -sać
	"opisać":      {stem: "opisz", class: ConjE},

	// wskazać - minority alternating -zać
	"wskazać":     {stem: "wskaż", class: ConjE},

	// brać prefix verbs with vowel elision
	"odebrać":  {sg13: "odbior", stem: "odbierz", class: ConjE},
	"zebrać":   {sg13: "zbior", stem: "zbierz", class: ConjE},
	"rozebrać": {sg13: "rozbior", stem: "rozbierz", class: ConjE},

	// lać verbs (j-insertion)
	"lać":       {stem: "lej", class: ConjE},

	// grześć - suppletive grzeb- stem
	"grześć":    {sg13: "grzeb", stem: "grzebi", class: ConjE},

	// żreć
	"żreć":      {stem: "żr", class: ConjE},

	// -przeć/-wrzeć verbs
	"przeć":     {sg13: "pr", stem: "prz", class: ConjE},
	"wrzeć":     {sg13: "wr", stem: "wrz", class: ConjE},

	// śnić - no j-insertion
	"śnić":      {sg13: "śni", stem: "śn", class: ConjI},

	// rzec - k-insertion
	"rzec":      {sg13: "rzekn", stem: "rzecz", class: ConjE},

	// tłuc - suppletive stem tłuk/tłucz
	"tłuc":      {sg13: "tłuk", stem: "tłucz", class: ConjE},

	// pleść - suppletive stem plot/plec
	"pleść":     {sg13: "plot", stem: "pleci", class: ConjE},

	// kląć - suppletive stem kln
	"kląć":      {sg13: "kln", stem: "klni", class: ConjE},

	// piąć - suppletive stem pn
	"piąć":      {sg13: "pn", stem: "pni", class: ConjE},
	"wspiąć":    {sg13: "wespn", stem: "wespni", class: ConjE},
	"zapiąć":    {sg13: "zapn", stem: "zapni", class: ConjE},
	"przypiąć":  {sg13: "przypn", stem: "przypni", class: ConjE},
	"odpiąć":    {sg13: "odpn", stem: "odpni", class: ConjE},
	"dopiąć":    {sg13: "dopn", stem: "dopni", class: ConjE},
	"spiąć":     {sg13: "spn", stem: "spni", class: ConjE},
	"wpiąć":     {sg13: "wpn", stem: "wpni", class: ConjE},
	"napiąć":    {sg13: "napn", stem: "napni", class: ConjE},
	"rozpiąć":   {sg13: "rozpn", stem: "rozpni", class: ConjE},
	"wypiąć":    {sg13: "wypn", stem: "wypni", class: ConjE},

	// wiać - special pattern (wieję)
	"wiać":      {stem: "wiej", class: ConjE},

	// chwiać - sway (chwieję)
	"chwiać":    {stem: "chwiej", class: ConjE},

	// krajać - j-insertion (kraję not krajam)
	"krajać":    {stem: "kraj", class: ConjE},

	// tajać - minority -ajać with -ję pattern
	"tajać":     {stem: "taj", class: ConjE},

	// ćpać - regular -am
	"ćpać":      {stem: "ćp", class: ConjA},

	// Regular -am -bać verbs (not alternating)
	"bimbać":    {stem: "bimb", class: ConjA},
	"gabać":     {stem: "gab", class: ConjA},
	"chybać":    {stem: "chyb", class: ConjA},
	"gibać":     {stem: "gib", class: ConjA},
	"gdybać":    {stem: "gdyb", class: ConjA},
	"zaniedbać": {stem: "zaniedb", class: ConjA},

	// Regular -am misc verbs
	"siorbać":   {stem: "siorb", class: ConjA},
	"stąpać":    {stem: "stąp", class: ConjA},
	"pchlać":    {stem: "pchl", class: ConjA},
	"rychlać":   {stem: "rychl", class: ConjA},
	"kpać":      {stem: "kp", class: ConjA},
	"kasłać":    {stem: "kasł", class: ConjA},
	"cierpać":   {stem: "cierp", class: ConjA},
	"siąpać":    {stem: "siąp", class: ConjA},
	"tyrpać":    {stem: "tyrp", class: ConjA},
	"ściubać":   {stem: "ściub", class: ConjA},
	"ślipać":    {stem: "ślip", class: ConjA},
	"bombać":    {stem: "bomb", class: ConjA},

	// Inchoative -eć verbs (use -eję pattern)
	"chorzeć":     {stem: "chorzej", class: ConjE},
	"tężeć":       {stem: "tężej", class: ConjE},
	"dumieć":      {stem: "dumiej", class: ConjE},
	"goreć":       {stem: "gorej", class: ConjE},
	"śniedzieć":   {stem: "śniedziej", class: ConjE},
	"srebrzeć":    {stem: "srebrzej", class: ConjE},
	"cukrzeć":     {stem: "cukrzej", class: ConjE},
	"dorośleć":    {stem: "doroślej", class: ConjE},
	"wydorośleć":  {stem: "wydoroślej", class: ConjE},
	"zelżeć":      {stem: "zelżej", class: ConjE},
	"wilżeć":      {stem: "wilżej", class: ConjE},
	"wężeć":       {stem: "wężej", class: ConjE},
	"rzedzieć":    {stem: "rzedziej", class: ConjE},
	"sfolżeć":     {stem: "sfolżej", class: ConjE},
	"szlachcieć":  {stem: "szlachciej", class: ConjE},
	"ochujeć":     {stem: "ochujej", class: ConjE},
	"ociężeć":     {stem: "ociężej", class: ConjE},
	"ściężeć":     {stem: "ściężej", class: ConjE},
	"oszedzieć":   {stem: "oszedziej", class: ConjE},
	"szedzieć":    {stem: "szedziej", class: ConjE},
	"sposążeć":    {stem: "sposążej", class: ConjE},
	"wyryżeć":     {stem: "wyryżej", class: ConjE},

	// siać - to sow (ia → ie + ję)
	"siać":        {stem: "siej", class: ConjE},

	// -tajać verbs meaning "to conceal" (use -tajam not -taję)
	"utajać":      {stem: "utaj", class: ConjA},
	"zatajać":     {stem: "zataj", class: ConjA},

	// łajać - minority -ajać that uses -ę/-esz
	"łajać":       {stem: "łaj", class: ConjE},

	// knajać - uses -ę/-esz
	"knajać":      {stem: "knaj", class: ConjE},

	// -iwać verbs that use -iwuję pattern
	"strzeliwać": {stem: "strzeliwuj", class: ConjE},
	"myśliwać":   {stem: "myśliwuj", class: ConjE},
	"boliwać":    {stem: "boliwuj", class: ConjE},
	"mgliwać":    {stem: "mgliwuj", class: ConjE},
	"skuliwać":   {stem: "skuliwuj", class: ConjE},

	// -ić verbs with no j-insertion
	"tlić":        {stem: "tl", class: ConjI},
	"clić":        {stem: "cl", class: ConjI},
	"dlić":        {stem: "dl", class: ConjI},

	// -ywać verbs that use -uję pattern
	"mieszywać":   {stem: "mieszuj", class: ConjE},
	"supływać":    {stem: "supłuj", class: ConjE},
	"bazgrywać":   {stem: "bazgruj", class: ConjE},
	"podobywać":   {stem: "podobuj", class: ConjE},

	// dziamdziać - uses -am pattern
	"dziamdziać":  {stem: "dziamdzi", class: ConjA},

	// -piać verbs: uses -eję pattern
	"piać":        {stem: "piej", class: ConjE},
	"spiać":       {stem: "spiej", class: ConjE},
	"dośpiać":     {stem: "dośpiej", class: ConjE},
	"uśpiać":      {stem: "uśpiej", class: ConjE},

	// pomieć - uses -am pattern
	"pomieć":      {stem: "pom", class: ConjA},

	// sposzyć - j-insertion
	"sposzyć":     {stem: "sposzyj", class: ConjE},

	// źreć/źrzeć - inchoative
	"źreć":        {stem: "źrej", class: ConjE},
	"źrzeć":       {stem: "źrzej", class: ConjE},

	// oziać - uses -eję pattern
	"oziać":       {stem: "oziej", class: ConjE},

	// Compound prefix -strzeliwać verbs
	"porozstrzeliwać": {stem: "porozstrzeliwuj", class: ConjE},
	"powystrzeliwać":  {stem: "powystrzeliwuj", class: ConjE},
}

// irregularPresOverrides contains verbs that are fully suppletive
// and cannot be described by presSpec.
var irregularPresOverrides = map[string]PresentTense{
	"być": {
		Sg1: "jestem", Sg2: "jesteś", Sg3: "jest",
		Pl1: "jesteśmy", Pl2: "jesteście", Pl3: "są",
	},
}

// irregularVerbs is the runtime lookup table, populated by init().
var irregularVerbs map[string]PresentTense

func init() {
	irregularVerbs = make(map[string]PresentTense, len(irregularPresSpecs)+len(irregularPresOverrides))
	for verb, spec := range irregularPresSpecs {
		irregularVerbs[verb] = spec.build()
	}
	for verb, pt := range irregularPresOverrides {
		irregularVerbs[verb] = pt
	}
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
