package fasta

import (
	"strings"
	"testing"
)

const multiFASTA = `>seq1 first record
ACGTACGTAC
GTAC
>seq2 second record
ggatccgaattc
>seq3
NNNNACGT
`

func TestParseMultiRecord(t *testing.T) {
	recs, err := ParseString(multiFASTA)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(recs) != 3 {
		t.Fatalf("got %d records, want 3", len(recs))
	}
	if recs[0].ID != "seq1" || recs[0].Description != "seq1 first record" {
		t.Errorf("record 0 ID/desc = %q/%q", recs[0].ID, recs[0].Description)
	}
	if recs[0].Sequence != "ACGTACGTACGTAC" {
		t.Errorf("record 0 sequence = %q", recs[0].Sequence)
	}
	if recs[1].Sequence != "GGATCCGAATTC" {
		t.Errorf("record 1 sequence = %q, want upper-cased", recs[1].Sequence)
	}
	if recs[2].Sequence != "NNNNACGT" {
		t.Errorf("record 2 sequence = %q", recs[2].Sequence)
	}
}

func TestParseCRLFAndBlanks(t *testing.T) {
	recs, err := ParseString(">a\r\nACGT\r\n\r\nACGT\r\n")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if recs[0].Sequence != "ACGTACGT" {
		t.Errorf("sequence = %q, want ACGTACGT", recs[0].Sequence)
	}
}

func TestParseErrors(t *testing.T) {
	bad := []string{
		"",
		"ACGT\n",
		">only-header\n",
		">x\nACGTX\n",
		">\nACGT\n",
	}
	for _, s := range bad {
		if _, err := ParseString(s); err == nil {
			t.Errorf("ParseString(%q) expected error, got nil", s)
		}
	}
}

func TestRecordGC(t *testing.T) {
	recs, _ := ParseString(">x\nGGCCAA\n")
	if got := recs[0].GC(); got != 4.0/6.0 {
		t.Errorf("GC = %v, want 0.6667", got)
	}
	if recs[0].Length() != 6 {
		t.Errorf("Length = %d, want 6", recs[0].Length())
	}
}

func TestParseLongLine(t *testing.T) {
	seq := strings.Repeat("ACGT", 50000) // 200 kb on one line
	recs, err := ParseString(">big\n" + seq + "\n")
	if err != nil {
		t.Fatalf("Parse long line failed: %v", err)
	}
	if recs[0].Length() != 200000 {
		t.Errorf("Length = %d, want 200000", recs[0].Length())
	}
}
