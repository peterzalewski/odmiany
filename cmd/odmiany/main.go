package main

import (
	"fmt"
	"os"

	"petezalew.ski/odmiany/pkg/verb"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: odmiany <verb> [verb2] [verb3] ...")
		os.Exit(1)
	}

	verbs := os.Args[1:]
	compact := len(verbs) > 1

	for i, infinitive := range verbs {
		paradigms, err := verb.ConjugatePresent(infinitive)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", infinitive, err)
			continue
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

		if !compact && i < len(verbs)-1 {
			fmt.Println()
		}
	}
}
