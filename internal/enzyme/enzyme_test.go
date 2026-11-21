package enzyme

import "testing"

func TestDatabaseValid(t *testing.T) {
	if len(Database) < 20 {
		t.Fatalf("database has %d enzymes, want at least 20", len(Database))
	}
	seen := map[string]bool{}
	for _, e := range Database {
		if err := e.Validate(); err != nil {
			t.Errorf("invalid enzyme %v: %v", e, err)
		}
		if seen[e.Name] {
			t.Errorf("duplicate enzyme name %q", e.Name)
		}
		seen[e.Name] = true
		if !e.IsPalindromic() {
			t.Errorf("enzyme %s site %s is not palindromic", e.Name, e.Site)
		}
	}
}

func TestOverhangs(t *testing.T) {
	cases := []struct {
		name       string
		overhangTy string
		overhangLn int
		overhangSq string
	}{
		{"EcoRI", "5'", 4, "AATT"},
		{"BamHI", "5'", 4, "GATC"},
		{"HindIII", "5'", 4, "AGCT"},
		{"PstI", "3'", 4, "TGCA"},
		{"KpnI", "3'", 4, "GTAC"},
		{"SacI", "3'", 4, "AGCT"},
		{"SmaI", "blunt", 0, ""},
		{"EcoRV", "blunt", 0, ""},
		{"HpaI", "blunt", 0, ""},
		{"PvuII", "blunt", 0, ""},
		{"AccI", "5'", 2, "MK"}, // degenerate overhang: M=A/C, K=G/T
		{"NotI", "5'", 4, "GGCC"},
		{"NdeI", "5'", 2, "TA"},
		{"XmaI", "5'", 4, "CCGG"},
		{"SphI", "3'", 4, "CATG"},
		{"ApaI", "3'", 4, "GGCC"},
	}
	for _, c := range cases {
		e, ok := Find(Database, c.name)
		if !ok {
			t.Fatalf("enzyme %s not in database", c.name)
		}
		if got := e.OverhangType(); got != c.overhangTy {
			t.Errorf("%s OverhangType = %q, want %q", c.name, got, c.overhangTy)
		}
		if got := e.OverhangLen(); got != c.overhangLn {
			t.Errorf("%s OverhangLen = %d, want %d", c.name, got, c.overhangLn)
		}
		if got := e.OverhangSeq(); got != c.overhangSq {
			t.Errorf("%s OverhangSeq = %q, want %q", c.name, got, c.overhangSq)
		}
	}
}

func TestRecognitionNotation(t *testing.T) {
	cases := map[string]string{
		"EcoRI": "G/AATTC",
		"BamHI": "G/GATCC",
		"PstI":  "CTGCA/G",
		"SmaI":  "CCC/GGG",
		"NotI":  "GC/GGCCGC",
		"AccI":  "GT/MKAC",
	}
	for name, want := range cases {
		e, _ := Find(Database, name)
		if got := e.Recognition(); got != want {
			t.Errorf("%s Recognition = %q, want %q", name, got, want)
		}
	}
}

// TestIUPACMatching checks the degenerate enzymes against hand-listed
// variants: AccI (GTMKAC) must match GTAGAC, GTATAC, GTCGAC, GTCTAC
// and nothing else among the 16 GT..AC combinations.
func TestIUPACMatching(t *testing.T) {
	accI, _ := Find(Database, "AccI")
	match := []string{"GTAGAC", "GTATAC", "GTCGAC", "GTCTAC"}
	for _, s := range match {
		if !accI.MatchesAt(s, 0) {
			t.Errorf("AccI should match %s", s)
		}
	}
	noMatch := []string{"GTGGAC", "GTAAAC", "GTCGAG", "ATAGAC", "GTAGAT"}
	for _, s := range noMatch {
		if accI.MatchesAt(s, 0) {
			t.Errorf("AccI should NOT match %s", s)
		}
	}

	// HinfI GANTC: N matches any base.
	hinfI, _ := Find(Database, "HinfI")
	for _, s := range []string{"GAATC", "GATTC", "GAGTC", "GACTC"} {
		if !hinfI.MatchesAt(s, 0) {
			t.Errorf("HinfI should match %s", s)
		}
	}
	if hinfI.MatchesAt("GATTC", 1) {
		t.Error("HinfI MatchesAt out of range should be false")
	}

	// AvaI CYCGRG: Y=C/T, R=A/G.
	avaI, _ := Find(Database, "AvaI")
	for _, s := range []string{"CTCGAG", "CCCGGG", "CCCGAG"} {
		if !avaI.MatchesAt(s, 0) {
			t.Errorf("AvaI should match %s", s)
		}
	}
	if avaI.MatchesAt("CACGAG", 0) {
		t.Error("AvaI should NOT match CACGAG (position 2 must be C/T)")
	}
	if avaI.MatchesAt("CTCGTG", 0) {
		t.Error("AvaI should NOT match CTCGTG (position 5 is T, R = A/G)")
	}
}

func TestExpectedFrequency(t *testing.T) {
	ecoRI, _ := Find(Database, "EcoRI")
	if got := ecoRI.ExpectedFrequency(); got != 1.0/4096.0 {
		t.Errorf("EcoRI ExpectedFrequency = %v, want 1/4096", got)
	}
	taqI, _ := Find(Database, "TaqI")
	if got := taqI.ExpectedFrequency(); got != 1.0/256.0 {
		t.Errorf("TaqI ExpectedFrequency = %v, want 1/256", got)
	}
	notI, _ := Find(Database, "NotI")
	if got := notI.ExpectedFrequency(); got != 1.0/65536.0 {
		t.Errorf("NotI ExpectedFrequency = %v, want 1/65536", got)
	}
	accI, _ := Find(Database, "AccI")
	// GTMKAC: 4 concrete variants out of 4^6
	if got := accI.ExpectedFrequency(); got != 4.0/4096.0 {
		t.Errorf("AccI ExpectedFrequency = %v, want 4/4096", got)
	}
	if got := ecoRI.ExpectedCount(4096); got != 1.0 {
		t.Errorf("EcoRI ExpectedCount(4096) = %v, want 1", got)
	}
}

func TestFindCaseInsensitive(t *testing.T) {
	if _, ok := Find(Database, "ecori"); !ok {
		t.Error("Find(ecori) should be case-insensitive")
	}
	if _, ok := Find(Database, "NoSuchEnzyme"); ok {
		t.Error("Find(NoSuchEnzyme) should fail")
	}
}

func TestValidate(t *testing.T) {
	bad := []Enzyme{
		{Name: "", Site: "GAATTC", TopCut: 1},
		{Name: "X", Site: "", TopCut: 0},
		{Name: "X", Site: "GAXTTC", TopCut: 1},
		{Name: "X", Site: "GAATTC", TopCut: 7},
	}
	for _, e := range bad {
		if err := e.Validate(); err == nil {
			t.Errorf("Validate(%v) expected error, got nil", e)
		}
	}
}
