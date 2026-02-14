package main

import (
	"fmt"
	"os"

	"petezalew.ski/odmiany/pkg/verb"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: odmiany <verb>")
		os.Exit(1)
	}

	infinitive := os.Args[1]
	paradigms, err := verb.ConjugatePresent(infinitive)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Present tense of %s:\n", infinitive)
	for i, p := range paradigms {
		if len(paradigms) > 1 {
			if p.Gloss != "" {
				fmt.Printf("\n  [%d] %s:\n", i+1, p.Gloss)
			} else {
				fmt.Printf("\n  [%d]:\n", i+1)
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
