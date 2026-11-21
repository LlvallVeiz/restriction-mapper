// Package tests holds end-to-end tests: they build the real rmap
// binary and run it against the example FASTA files, checking the
// headline numbers that the README quotes.
package tests

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(t.TempDir(), "rmap")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/rmap")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}
	return bin
}

func run(t *testing.T, bin string, args ...string) string {
	t.Helper()
	cmd := exec.Command(bin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("rmap %v failed: %v\n%s", args, err, out)
	}
	return string(out)
}

func TestE2EEnzymes(t *testing.T) {
	bin := buildBinary(t)
	out := run(t, bin, "enzymes")
	for _, want := range []string{"EcoRI", "G/AATTC", "BamHI", "NotI", "32 enzymes"} {
		if !strings.Contains(out, want) {
			t.Errorf("enzymes output missing %q:\n%s", want, out)
		}
	}
}

func TestE2EPuc19Scan(t *testing.T) {
	bin := buildBinary(t)
	out := run(t, bin, "scan", "-in", "../examples/input/puc19.fasta", "-circular")
	// Headline numbers, cross-checked independently in Python
	// (scripts/verify-puc19.py): 27 sites, 18 enzymes cut.
	if !strings.Contains(out, "27 sites from 18 enzymes") {
		t.Errorf("pUC19 scan summary wrong:\n%s", out)
	}
	// The MCS cluster in the L09137.2 orientation: EcoRI 396, BamHI
	// 417, HindIII 447.
	for _, want := range []string{"396   EcoRI", "417   BamHI", "447   HindIII"} {
		if !strings.Contains(out, want) {
			t.Errorf("pUC19 scan missing %q:\n%s", want, out)
		}
	}
}

func TestE2EPuc19DoubleDigest(t *testing.T) {
	bin := buildBinary(t)
	out := run(t, bin, "digest", "-in", "../examples/input/puc19.fasta",
		"-enzymes", "EcoRI,HindIII", "-circular")
	// EcoRI cut 396, HindIII cut 447 -> fragments 51 bp (MCS) + 2635 bp.
	if !strings.Contains(out, "51 bp") || !strings.Contains(out, "2635 bp") {
		t.Errorf("double digest fragments wrong:\n%s", out)
	}
	if !strings.Contains(out, "Total: 2686 bp") {
		t.Errorf("double digest total wrong:\n%s", out)
	}
}

func TestE2ESyntheticDigest(t *testing.T) {
	bin := buildBinary(t)
	out := run(t, bin, "digest", "-in", "../examples/input/synthetic-plasmid.fasta",
		"-enzymes", "EcoRI,BamHI", "-circular")
	// Designed cuts at 101 (EcoRI) and 201 (BamHI) -> 100 + 3100 bp.
	if !strings.Contains(out, "100 bp") || !strings.Contains(out, "3100 bp") {
		t.Errorf("synthetic digest fragments wrong:\n%s", out)
	}
}

func TestE2EBatchScan(t *testing.T) {
	bin := buildBinary(t)
	out := run(t, bin, "scan", "-in", "../examples/input/multi.fasta")
	for _, want := range []string{
		"frag-a: 24 bp, 2 sites",
		"frag-b: 18 bp, 3 sites",
		"frag-c: 12 bp, 0 sites",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("batch scan missing %q:\n%s", want, out)
		}
	}
}

func TestE2EUnknownEnzyme(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "scan", "-in", "../examples/input/puc19.fasta",
		"-enzymes", "NoSuchEnzyme")
	if err := cmd.Run(); err == nil {
		t.Error("scan with unknown enzyme should fail, but exited 0")
	}
}

func TestE2EMissingFile(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "scan", "-in", "../examples/input/does-not-exist.fasta")
	if err := cmd.Run(); err == nil {
		t.Error("scan with missing file should fail, but exited 0")
	}
}
