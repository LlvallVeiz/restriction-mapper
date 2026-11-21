// Package dna provides basic operations on DNA sequences: complement,
// reverse complement, GC content and validation. All functions accept
// upper- or lower-case input and treat sequences as plain strings of
// IUPAC base codes (A, C, G, T, U, N and the ambiguity codes).
package dna

import (
	"fmt"
	"strings"
)

// complement maps each IUPAC base code to its complement (upper case).
// Ambiguity codes complement to the code of the complementary set, e.g.
// R (A/G) <-> Y (C/T), S (G/C) <-> S, W (A/T) <-> W, K (G/T) <-> M
// (A/C), B (C/G/T) <-> V (A/C/G), D (A/G/T) <-> H (A/C/T), N <-> N.
var complement = map[byte]byte{
	'A': 'T', 'T': 'A', 'U': 'A',
	'C': 'G', 'G': 'C',
	'R': 'Y', 'Y': 'R',
	'S': 'S', 'W': 'W',
	'K': 'M', 'M': 'K',
	'B': 'V', 'V': 'B',
	'D': 'H', 'H': 'D',
	'N': 'N',
}

// IUPACSet maps each IUPAC code to the set of concrete bases it stands for.
var IUPACSet = map[byte]string{
	'A': "A", 'C': "C", 'G': "G", 'T': "T", 'U': "T",
	'R': "AG", 'Y': "CT", 'S': "GC", 'W': "AT",
	'K': "GT", 'M': "AC",
	'B': "CGT", 'D': "AGT", 'H': "ACT", 'V': "ACG",
	'N': "ACGT",
}

// Complement returns the complement of a single base (upper case).
// It returns 0 for characters that are not IUPAC codes.
func Complement(b byte) byte {
	if b >= 'a' && b <= 'z' {
		b -= 'a' - 'A'
	}
	c, ok := complement[b]
	if !ok {
		return 0
	}
	return c
}

// ReverseComplement returns the reverse complement of s, upper case.
// Characters that are not IUPAC codes are passed through unchanged so
// callers can decide how to handle them (use Validate first for strict
// checking).
func ReverseComplement(s string) string {
	s = strings.ToUpper(s)
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := complement[s[len(s)-1-i]]
		if c == 0 {
			c = s[len(s)-1-i]
		}
		out[i] = c
	}
	return string(out)
}

// Clean upper-cases s and removes all white space.
func Clean(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		b.WriteByte(c)
	}
	return b.String()
}

// Validate reports whether s contains only IUPAC base codes (after
// cleaning). It returns a descriptive error naming the first bad
// character and its 1-based position.
func Validate(s string) error {
	s = Clean(s)
	for i := 0; i < len(s); i++ {
		if _, ok := IUPACSet[s[i]]; !ok {
			return fmt.Errorf("invalid base %q at position %d", s[i], i+1)
		}
	}
	if len(s) == 0 {
		return fmt.Errorf("empty sequence")
	}
	return nil
}

// GCContent returns the fraction of G+C bases in s (0..1). Ambiguity
// codes contribute fractionally: S counts as 1, R/Y/W/K/M as 0.5,
// B/V as 2/3, D/H as 1/3, N as 0.5. Pure ACGT sequences are exact.
func GCContent(s string) float64 {
	s = Clean(s)
	if len(s) == 0 {
		return 0
	}
	var gc float64
	for i := 0; i < len(s); i++ {
		set := IUPACSet[s[i]]
		if set == "" {
			continue
		}
		g := 0
		for j := 0; j < len(set); j++ {
			if set[j] == 'G' || set[j] == 'C' {
				g++
			}
		}
		gc += float64(g) / float64(len(set))
	}
	return gc / float64(len(s))
}

// IsPalindromic reports whether s reads the same on both strands, i.e.
// s equals its own reverse complement (the defining property of most
// type II restriction sites).
func IsPalindromic(s string) bool {
	s = Clean(s)
	return s == ReverseComplement(s)
}
