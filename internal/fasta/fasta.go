// Package fasta parses FASTA files: one or more records, each a header
// line starting with '>' followed by sequence lines. Headers keep the
// full description; the ID is the first whitespace-separated token.
package fasta

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/LlvallVeiz/restriction-mapper/internal/dna"
)

// Record is one FASTA entry.
type Record struct {
	// ID is the first token of the header line (without '>').
	ID string
	// Description is the full header line (without '>').
	Description string
	// Sequence is the concatenated, upper-cased sequence (white space removed).
	Sequence string
}

// Length returns the sequence length in base pairs.
func (r Record) Length() int { return len(r.Sequence) }

// GC returns the GC fraction of the sequence (0..1).
func (r Record) GC() float64 { return dna.GCContent(r.Sequence) }

// Parse reads all records from r. It accepts lower-case bases, blank
// lines and CRLF line endings. It returns an error for a missing
// header, a record with an empty sequence, or invalid base characters.
func Parse(r io.Reader) ([]Record, error) {
	var records []Record
	var cur *Record
	sc := bufio.NewScanner(r)
	// Allow long sequence lines (default 64 KiB is too small for FASTA).
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, ">") {
			header := strings.TrimSpace(line[1:])
			if header == "" {
				return nil, fmt.Errorf("line %d: empty FASTA header", lineNo)
			}
			if cur != nil {
				if err := finish(cur, lineNo); err != nil {
					return nil, err
				}
				records = append(records, *cur)
			}
			id := header
			if i := strings.IndexAny(header, " \t"); i >= 0 {
				id = header[:i]
			}
			cur = &Record{ID: id, Description: header}
			continue
		}
		if cur == nil {
			return nil, fmt.Errorf("line %d: sequence data before any header", lineNo)
		}
		cur.Sequence += dna.Clean(line)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("reading FASTA: %w", err)
	}
	if cur != nil {
		if err := finish(cur, lineNo); err != nil {
			return nil, err
		}
		records = append(records, *cur)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("no FASTA records found")
	}
	return records, nil
}

func finish(r *Record, lineNo int) error {
	if len(r.Sequence) == 0 {
		return fmt.Errorf("record %q: empty sequence", r.ID)
	}
	if err := dna.Validate(r.Sequence); err != nil {
		return fmt.Errorf("record %q: %w", r.ID, err)
	}
	return nil
}

// ParseString parses FASTA content from a string (convenience for tests).
func ParseString(s string) ([]Record, error) {
	return Parse(strings.NewReader(s))
}
