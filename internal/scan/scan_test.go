package scan

import (
	"testing"

	"github.com/LlvallVeiz/restriction-mapper/internal/enzyme"
)

func mustEnzyme(t *testing.T, name string) enzyme.Enzyme {
	t.Helper()
	e, ok := enzyme.Find(enzyme.Database, name)
	if !ok {
		t.Fatalf("enzyme %s not in database", name)
	}
	return e
}

// testSeq is hand-built: positions (1-based) of the embedded sites are
// written out explicitly so the expectations below are hand-computed,
// not copied from the code under test.
//
//	1         11        21
//	123456789012345678901234
//	AAAAGAATTC CCCCGGATCC TTTT
//	     EcoRI^    BamHI^
//
// EcoRI site GAATTC at 5..10, top cut after base 5, bottom cut after 9.
// BamHI site GGATCC at 15..20, top cut after 15, bottom cut after 19.
const testSeq = "AAAAGAATTCCCCCGGATCCTTTT"

func TestScanKnownSites(t *testing.T) {
	enzymes := []enzyme.Enzyme{mustEnzyme(t, "EcoRI"), mustEnzyme(t, "BamHI")}
	sites := Scan(testSeq, enzymes, false)
	if len(sites) != 2 {
		t.Fatalf("got %d sites, want 2: %+v", len(sites), sites)
	}
	eco := sites[0]
	if eco.Enzyme != "EcoRI" || eco.Start != 5 || eco.Strand != "+" {
		t.Errorf("site 0 = %+v, want EcoRI at 5 +", eco)
	}
	if eco.CutTop != 5 || eco.CutBottom != 9 {
		t.Errorf("EcoRI cuts = %d/%d, want 5/9", eco.CutTop, eco.CutBottom)
	}
	bam := sites[1]
	if bam.Enzyme != "BamHI" || bam.Start != 15 {
		t.Errorf("site 1 = %+v, want BamHI at 15", bam)
	}
	if bam.CutTop != 15 || bam.CutBottom != 19 {
		t.Errorf("BamHI cuts = %d/%d, want 15/19", bam.CutTop, bam.CutBottom)
	}
}

func TestScanNoDoubleCountPalindrome(t *testing.T) {
	// A palindromic site must be reported once, not once per strand.
	sites := Scan("GAATTC", []enzyme.Enzyme{mustEnzyme(t, "EcoRI")}, false)
	if len(sites) != 1 {
		t.Fatalf("got %d sites, want 1: %+v", len(sites), sites)
	}
}

func TestScanOverlapping(t *testing.T) {
	// "GAATTCGAATTC" holds two adjacent EcoRI sites (11 and 12 share
	// no bases, but the scan must not skip the second one).
	sites := Scan("GAATTCGAATTC", []enzyme.Enzyme{mustEnzyme(t, "EcoRI")}, false)
	if len(sites) != 2 || sites[0].Start != 1 || sites[1].Start != 7 {
		t.Fatalf("overlapping scan = %+v, want starts 1 and 7", sites)
	}
	// Truly overlapping windows: a 4 bp site "AAAA" occurs at every
	// offset of "AAAAA". This guards against a matcher that skips
	// ahead past a match (regex-style non-overlapping search).
	aaaa := enzyme.Enzyme{Name: "TestAAAA", Site: "AAAA", TopCut: 2}
	sites = Scan("AAAAA", []enzyme.Enzyme{aaaa}, false)
	if len(sites) != 2 || sites[0].Start != 1 || sites[1].Start != 2 {
		t.Fatalf("AAAA scan = %+v, want starts 1 and 2", sites)
	}
	// Two different enzymes sharing the same bases: SmaI (CCC/GGG)
	// and XmaI (C/CCGGG) both cut "CCCGGG" at position 1.
	sites = Scan("CCCGGG", []enzyme.Enzyme{mustEnzyme(t, "SmaI"), mustEnzyme(t, "XmaI")}, false)
	if len(sites) != 2 {
		t.Fatalf("SmaI/XmaI scan = %+v, want 2 sites", sites)
	}
	if sites[0].Enzyme != "SmaI" || sites[1].Enzyme != "XmaI" {
		t.Errorf("shared-site enzymes = %q, %q, want SmaI, XmaI", sites[0].Enzyme, sites[1].Enzyme)
	}
}

func TestScanCircularWrap(t *testing.T) {
	// "ATTC" + "CCG" + "GGA": the EcoRI site GAATTC wraps the origin
	// (bases 9-10 = GA, bases 1-4 = ATTC).
	seq := "ATTCCCGGGA" // length 10; site occupies 9,10,1,2,3,4
	linear := Scan(seq, []enzyme.Enzyme{mustEnzyme(t, "EcoRI")}, false)
	if len(linear) != 0 {
		t.Fatalf("linear scan found wrapping site: %+v", linear)
	}
	circ := Scan(seq, []enzyme.Enzyme{mustEnzyme(t, "EcoRI")}, true)
	if len(circ) != 1 {
		t.Fatalf("circular scan = %+v, want 1 site", circ)
	}
	if circ[0].Start != 9 {
		t.Errorf("wrapping site Start = %d, want 9", circ[0].Start)
	}
	// CutTop: start0=8, TopCut=1 -> 9 bases left of the cut.
	if circ[0].CutTop != 9 || circ[0].CutBottom != 3 {
		t.Errorf("wrapping site cuts = %d/%d, want 9/3", circ[0].CutTop, circ[0].CutBottom)
	}
}

func TestScanReverseComplementOnly(t *testing.T) {
	// A genuinely non-palindromic site can only be found on the bottom
	// strand: custom enzyme "GAAACC" (reverse complement "GGTTTC").
	// The top strand "GGTTTC" holds no GAAACC, but its bottom strand
	// (read 5'->3') is GAAACC. Note: HinfI GANTC cannot be used here
	// because GANTC is self-reverse-complementary at the IUPAC level.
	nonPal := enzyme.Enzyme{Name: "TestNP", Site: "GAAACC", TopCut: 2}
	sites := Scan("GGTTTC", []enzyme.Enzyme{nonPal}, false)
	if len(sites) != 1 {
		t.Fatalf("got %d sites, want 1: %+v", len(sites), sites)
	}
	if sites[0].Strand != "-" || sites[0].Start != 1 {
		t.Errorf("site = %+v, want strand - at 1", sites[0])
	}
	if sites[0].CutTop != 2 || sites[0].CutBottom != 4 {
		t.Errorf("site cuts = %d/%d, want 2/4", sites[0].CutTop, sites[0].CutBottom)
	}
}

func TestScanDegenerate(t *testing.T) {
	// AccI GTMKAC matches GTAGAC (M=A) and GTCTAC (M=C, K=T).
	sites := Scan("GTAGACGTCTAC", []enzyme.Enzyme{mustEnzyme(t, "AccI")}, false)
	if len(sites) != 2 || sites[0].Start != 1 || sites[1].Start != 7 {
		t.Fatalf("AccI scan = %+v, want starts 1 and 7", sites)
	}
}

func TestNonCutters(t *testing.T) {
	sites := Scan(testSeq, enzyme.Database, false)
	nc := NonCutters(sites, enzyme.Database)
	for _, name := range []string{"EcoRI", "BamHI"} {
		for _, n := range nc {
			if n == name {
				t.Errorf("%s reported as non-cutter", name)
			}
		}
	}
	found := false
	for _, n := range nc {
		if n == "NotI" {
			found = true
		}
	}
	if !found {
		t.Error("NotI should be a non-cutter on testSeq")
	}
}

func TestCountPerEnzyme(t *testing.T) {
	sites := Scan("GAATTCGAATTCGGATCC", []enzyme.Enzyme{mustEnzyme(t, "EcoRI"), mustEnzyme(t, "BamHI")}, false)
	counts := CountPerEnzyme(sites)
	if counts["EcoRI"] != 2 || counts["BamHI"] != 1 {
		t.Errorf("counts = %v, want EcoRI 2, BamHI 1", counts)
	}
}

func TestScanEmpty(t *testing.T) {
	if got := Scan("", enzyme.Database, false); len(got) != 0 {
		t.Errorf("empty sequence gave %d sites", len(got))
	}
}

func TestTopCutPositions(t *testing.T) {
	sites := Scan(testSeq, []enzyme.Enzyme{mustEnzyme(t, "EcoRI"), mustEnzyme(t, "BamHI")}, false)
	cuts := TopCutPositions(sites)
	if len(cuts) != 2 || cuts[0] != 5 || cuts[1] != 15 {
		t.Errorf("cuts = %v, want [5 15]", cuts)
	}
}
