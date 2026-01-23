package verb

import "testing"

func TestConjugatePresentAc(t *testing.T) {
	tests := []struct {
		infinitive string
		want       PresentTense
	}{
		{
			infinitive: "czytać",
			want: PresentTense{
				Sg1: "czytam",
				Sg2: "czytasz",
				Sg3: "czyta",
				Pl1: "czytamy",
				Pl2: "czytacie",
				Pl3: "czytają",
			},
		},
		{
			infinitive: "grać",
			want: PresentTense{
				Sg1: "gram",
				Sg2: "grasz",
				Sg3: "gra",
				Pl1: "gramy",
				Pl2: "gracie",
				Pl3: "grają",
			},
		},
		{
			infinitive: "kochać",
			want: PresentTense{
				Sg1: "kocham",
				Sg2: "kochasz",
				Sg3: "kocha",
				Pl1: "kochamy",
				Pl2: "kochacie",
				Pl3: "kochają",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.infinitive, func(t *testing.T) {
			got, err := ConjugatePresent(tt.infinitive)
			if err != nil {
				t.Fatalf("ConjugatePresent(%q) error: %v", tt.infinitive, err)
			}
			if got != tt.want {
				t.Errorf("ConjugatePresent(%q) =\n%+v\nwant:\n%+v", tt.infinitive, got, tt.want)
			}
		})
	}
}

func TestConjugatePresentUnsupported(t *testing.T) {
	// These verb classes are not yet supported.
	// Note: pisać ends in -ać but is conjugation III (piszę), not I (czytam).
	// We don't yet distinguish these, so pisać would be incorrectly conjugated.
	unsupported := []string{"robić", "nieść", "być"}
	for _, v := range unsupported {
		t.Run(v, func(t *testing.T) {
			_, err := ConjugatePresent(v)
			if err == nil {
				t.Errorf("ConjugatePresent(%q) expected error, got nil", v)
			}
		})
	}
}

func TestPresentTenseGet(t *testing.T) {
	p := PresentTense{
		Sg1: "czytam",
		Sg2: "czytasz",
		Sg3: "czyta",
		Pl1: "czytamy",
		Pl2: "czytacie",
		Pl3: "czytają",
	}

	tests := []struct {
		person Person
		number Number
		want   string
	}{
		{First, Singular, "czytam"},
		{Second, Singular, "czytasz"},
		{Third, Singular, "czyta"},
		{First, Plural, "czytamy"},
		{Second, Plural, "czytacie"},
		{Third, Plural, "czytają"},
	}

	for _, tt := range tests {
		got := p.Get(tt.person, tt.number)
		if got != tt.want {
			t.Errorf("Get(%d, %d) = %q, want %q", tt.person, tt.number, got, tt.want)
		}
	}
}
