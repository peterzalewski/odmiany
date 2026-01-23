# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`odmiany` is a Polish morphology engine for conjugating verbs and declining nouns. The name means "conjugations/declensions" in Polish.

## Build and Test Commands

```bash
go build ./...              # Build all packages
go test ./...               # Run all tests
go test -run TestName ./pkg # Run a specific test
go run ./cmd/odmiany        # Run the CLI
```

## Architecture

Rule-based morphological engine (not ML). Polish morphology is systematic with finite, enumerable patterns.

**Key design principles:**
- Rules for regular patterns, lookup tables for irregulars
- Stem + ending composition for word form generation
- Validation against Morfeusz2/Polimorf or Wiktionary data

**Planned structure:**
- `cmd/odmiany/` - CLI entrypoint
- `pkg/verb/` - Verb conjugation (aspect, tense, person, number, gender in past)
- `pkg/noun/` - Noun declension (7 cases, 2 numbers, 3 genders)
- `pkg/stem/` - Stem extraction and alternation rules (vowel/consonant changes)

## Polish Morphology Notes

- Verbs: aspect pairs (imperfective/perfective), 3 persons × 2 numbers, tenses, gender distinction in past tense
- Nouns: 7 cases (mianownik, dopełniacz, celownik, biernik, narzędnik, miejscownik, wołacz), singular/plural, masculine/feminine/neuter with masculine subcategories (personal, animate, inanimate)
- Common alternations: ą→ę, o→ó, stem consonant softening
