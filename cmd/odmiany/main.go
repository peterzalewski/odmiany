package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"petezalew.ski/odmiany/pkg/verb"
)

func main() {
	past := flag.Bool("past", false, "show past tense conjugation")
	gerund := flag.Bool("gerund", false, "show verbal noun (rzeczownik odsłownikowy)")
	flag.Parse()

	verbs := flag.Args()
	if len(verbs) < 1 {
		fmt.Fprintln(os.Stderr, "usage: odmiany [-past|-gerund] <verb> [verb2] [verb3] ...")
		os.Exit(1)
	}

	compact := len(verbs) > 1

	for i, infinitive := range verbs {
		switch {
		case *gerund:
			showVerbalNoun(infinitive)
		case *past:
			showPastTense(infinitive, compact)
		default:
			showPresentTense(infinitive, compact)
		}

		if !compact && i < len(verbs)-1 {
			fmt.Println()
		}
	}
}

func showVerbalNoun(infinitive string) {
	forms, err := verb.VerbalNoun(infinitive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", infinitive, err)
		return
	}
	fmt.Printf("%s: %s\n", infinitive, strings.Join(forms, ", "))
}

func showPresentTense(infinitive string, compact bool) {
	paradigms, err := verb.ConjugatePresent(infinitive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", infinitive, err)
		return
	}

	if compact {
		// Compact format for multiple verbs
		for _, p := range paradigms {
			fmt.Printf("%s: %s, %s, %s, %s, %s, %s\n",
				infinitive, p.Sg1, p.Sg2, p.Sg3, p.Pl1, p.Pl2, p.Pl3)
		}
	} else {
		// Detailed format for single verb
		fmt.Printf("Present tense of %s:\n", infinitive)
		for j, p := range paradigms {
			if len(paradigms) > 1 {
				if p.Gloss != "" {
					fmt.Printf("\n  [%d] %s:\n", j+1, p.Gloss)
				} else {
					fmt.Printf("\n  [%d]:\n", j+1)
				}
			}
			fmt.Printf("  ja      %s\n", p.Sg1)
			fmt.Printf("  ty      %s\n", p.Sg2)
			fmt.Printf("  on/ona  %s\n", p.Sg3)
			fmt.Printf("  my      %s\n", p.Pl1)
			fmt.Printf("  wy      %s\n", p.Pl2)
			fmt.Printf("  oni/one %s\n", p.Pl3)
		}
	}
}

func showPastTense(infinitive string, compact bool) {
	paradigms, err := verb.ConjugatePast(infinitive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", infinitive, err)
		return
	}

	if compact {
		// Compact format for multiple verbs
		for _, p := range paradigms {
			fmt.Printf("%s: %s/%s, %s/%s, %s/%s/%s, %s/%s, %s/%s, %s/%s\n",
				infinitive,
				p.Sg1M, p.Sg1F, p.Sg2M, p.Sg2F, p.Sg3M, p.Sg3F, p.Sg3N,
				p.Pl1V, p.Pl1NV, p.Pl2V, p.Pl2NV, p.Pl3V, p.Pl3NV)
		}
	} else {
		// Detailed format for single verb
		fmt.Printf("Past tense of %s:\n", infinitive)
		for j, p := range paradigms {
			if len(paradigms) > 1 {
				if p.Gloss != "" {
					fmt.Printf("\n  [%d] %s:\n", j+1, p.Gloss)
				} else {
					fmt.Printf("\n  [%d]:\n", j+1)
				}
			}
			fmt.Printf("  ja (m)      %s\n", p.Sg1M)
			fmt.Printf("  ja (f)      %s\n", p.Sg1F)
			fmt.Printf("  ty (m)      %s\n", p.Sg2M)
			fmt.Printf("  ty (f)      %s\n", p.Sg2F)
			fmt.Printf("  on          %s\n", p.Sg3M)
			fmt.Printf("  ona         %s\n", p.Sg3F)
			fmt.Printf("  ono         %s\n", p.Sg3N)
			fmt.Printf("  my (v)      %s\n", p.Pl1V)
			fmt.Printf("  my (nv)     %s\n", p.Pl1NV)
			fmt.Printf("  wy (v)      %s\n", p.Pl2V)
			fmt.Printf("  wy (nv)     %s\n", p.Pl2NV)
			fmt.Printf("  oni         %s\n", p.Pl3V)
			fmt.Printf("  one         %s\n", p.Pl3NV)
		}
	}
}
