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
			paradigms, err := ConjugatePresent(tt.infinitive)
			if err != nil {
				t.Fatalf("ConjugatePresent(%q) error: %v", tt.infinitive, err)
			}
			if len(paradigms) == 0 {
				t.Fatalf("ConjugatePresent(%q) returned no paradigms", tt.infinitive)
			}
			got := paradigms[0].PresentTense
			if got != tt.want {
				t.Errorf("ConjugatePresent(%q) =\n%+v\nwant:\n%+v", tt.infinitive, got, tt.want)
			}
		})
	}
}

func TestConjugatePresentSupported(t *testing.T) {
	// These verbs should now be supported via irregulars or heuristics.
	supported := []string{"robić", "nieść", "być"}
	for _, v := range supported {
		t.Run(v, func(t *testing.T) {
			paradigms, err := ConjugatePresent(v)
			if err != nil {
				t.Errorf("ConjugatePresent(%q) returned error: %v", v, err)
			}
			if len(paradigms) == 0 {
				t.Errorf("ConjugatePresent(%q) returned no paradigms", v)
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

func TestHomographs(t *testing.T) {
	// Test that homographs return multiple paradigms
	tests := []struct {
		infinitive string
		wantCount  int
	}{
		{"stać", 2},  // to stand vs to become
		{"słać", 2},  // to send vs to spread (bedding)
	}

	for _, tt := range tests {
		t.Run(tt.infinitive, func(t *testing.T) {
			paradigms, err := ConjugatePresent(tt.infinitive)
			if err != nil {
				t.Fatalf("ConjugatePresent(%q) error: %v", tt.infinitive, err)
			}
			if len(paradigms) != tt.wantCount {
				t.Errorf("ConjugatePresent(%q) returned %d paradigms, want %d",
					tt.infinitive, len(paradigms), tt.wantCount)
			}
		})
	}
}
