package dna

import "testing"

func TestReverseComplement(t *testing.T) {
	cases := []struct{ in, want string }{
		{"GAATTC", "GAATTC"}, // EcoRI site is a palindrome
		{"AAAACCCC", "GGGGTTTT"},
		{"ACGT", "ACGT"},
		{"ATGC", "GCAT"},
		{"gaattc", "GAATTC"}, // lower case accepted
		{"", ""},
		{"RY", "RY"},           // R<->Y, palindrome
		{"AAGCTT", "AAGCTT"},   // HindIII site is a palindrome
		{"GGATCC", "GGATCC"},   // BamHI site is a palindrome
		{"GATTACA", "TGTAATC"}, // hand-computed: complement CTTAAT reversed
	}
	for _, c := range cases {
		if got := ReverseComplement(c.in); got != c.want {
			t.Errorf("ReverseComplement(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestComplement(t *testing.T) {
	if got := Complement('A'); got != 'T' {
		t.Errorf("Complement(A) = %c, want T", got)
	}
	if got := Complement('r'); got != 'Y' {
		t.Errorf("Complement(r) = %c, want Y", got)
	}
	if got := Complement('N'); got != 'N' {
		t.Errorf("Complement(N) = %c, want N", got)
	}
	if got := Complement('X'); got != 0 {
		t.Errorf("Complement(X) = %c, want 0", got)
	}
}

func TestGCContent(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"GCGC", 1.0},
		{"ATAT", 0.0},
		{"GATC", 0.5},
		{"GGCCAA", 4.0 / 6.0},
		{"", 0.0},
		{"NNNN", 0.5}, // N is 50% GC by definition here
		{"SSSS", 1.0}, // S = G/C
		{"WWWW", 0.0}, // W = A/T
	}
	for _, c := range cases {
		if got := GCContent(c.in); got != c.want {
			t.Errorf("GCContent(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestValidate(t *testing.T) {
	if err := Validate("ACGTNacgtn"); err != nil {
		t.Errorf("Validate(ACGTN...) unexpected error: %v", err)
	}
	if err := Validate("ACGTX"); err == nil {
		t.Error("Validate(ACGTX) expected error, got nil")
	}
	if err := Validate(""); err == nil {
		t.Error("Validate(\"\") expected error, got nil")
	}
	if err := Validate("ACGT 123"); err == nil {
		t.Error("Validate with digits expected error, got nil")
	}
}

func TestIsPalindromic(t *testing.T) {
	for _, s := range []string{"GAATTC", "GGATCC", "AAGCTT", "GATATC", "CCCGGG"} {
		if !IsPalindromic(s) {
			t.Errorf("IsPalindromic(%q) = false, want true", s)
		}
	}
	if IsPalindromic("GATTACA") {
		t.Error("IsPalindromic(GATTACA) = true, want false")
	}
}

func TestClean(t *testing.T) {
	if got := Clean("ac gt\nAC"); got != "ACGTAC" {
		t.Errorf("Clean = %q, want ACGTAC", got)
	}
}
