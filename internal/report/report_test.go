package report

import (
	"strings"
	"testing"

	"github.com/LlvallVeiz/restriction-mapper/internal/enzyme"
	"github.com/LlvallVeiz/restriction-mapper/internal/scan"
)

func testSites() ([]scan.Site, []enzyme.Enzyme) {
	enzymes := []enzyme.Enzyme{
		{Name: "EcoRI", Site: "GAATTC", TopCut: 1},
		{Name: "BamHI", Site: "GGATCC", TopCut: 1},
	}
	sites := []scan.Site{
		{Enzyme: "EcoRI", Start: 5, Strand: "+", CutTop: 5, CutBottom: 9},
		{Enzyme: "BamHI", Start: 15, Strand: "+", CutTop: 15, CutBottom: 19},
	}
	return sites, enzymes
}

func TestSitesTable(t *testing.T) {
	sites, enzymes := testSites()
	out := SitesTable(sites, enzymes)
	if !strings.Contains(out, "EcoRI") || !strings.Contains(out, "G/AATTC") {
		t.Errorf("SitesTable missing EcoRI row:\n%s", out)
	}
	if !strings.Contains(out, "BamHI") || !strings.Contains(out, "G/GATCC") {
		t.Errorf("SitesTable missing BamHI row:\n%s", out)
	}
	if got := SitesTable(nil, enzymes); !strings.Contains(got, "no sites found") {
		t.Errorf("empty SitesTable should say so:\n%s", got)
	}
}

func TestStatsTable(t *testing.T) {
	sites, enzymes := testSites()
	out := StatsTable(sites, enzymes, 24)
	if !strings.Contains(out, "EcoRI") || !strings.Contains(out, "BamHI") {
		t.Errorf("StatsTable missing rows:\n%s", out)
	}
	// Expected EcoRI count in 24 bp of random DNA: 24/4096 = 0.01.
	if !strings.Contains(out, "0.01") {
		t.Errorf("StatsTable should show expected 0.01:\n%s", out)
	}
}

func TestMapASCII(t *testing.T) {
	sites, enzymes := testSites()
	out := MapASCII(24, sites, enzymes, 24)
	if !strings.Contains(out, "EcoRI") || !strings.Contains(out, "BamHI") {
		t.Errorf("MapASCII missing enzyme rows:\n%s", out)
	}
	if !strings.Contains(out, "24 bp") {
		t.Errorf("MapASCII missing length label:\n%s", out)
	}
	// Each enzyme row must carry exactly one '|' marker.
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "EcoRI") {
			if strings.Count(line, "|") != 1 {
				t.Errorf("EcoRI map row should have 1 marker: %q", line)
			}
		}
	}
	if got := MapASCII(0, nil, nil, 80); !strings.Contains(got, "empty sequence") {
		t.Errorf("empty MapASCII should say so: %q", got)
	}
}

func TestSummaryLine(t *testing.T) {
	sites, enzymes := testSites()
	got := SummaryLine("test", 24, sites, enzymes)
	want := "test: 24 bp, 2 sites from 2 enzymes, 0 enzymes do not cut"
	if got != want {
		t.Errorf("SummaryLine = %q, want %q", got, want)
	}
}
