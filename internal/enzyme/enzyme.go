// Package enzyme describes type II restriction endonucleases: their
// recognition sites (IUPAC), where they cut each strand, and the kind
// of ends they leave. All site and cut data in database.go were taken
// from the NEB Alphabetized List of Recognition Sequences and REBASE
// (accessed 2026-11-20); see docs/ENZYMES.md for the full table.
package enzyme

import (
	"fmt"
	"sort"
	"strings"

	"github.com/LlvallVeiz/restriction-mapper/internal/dna"
)

// Enzyme is one type II restriction endonuclease that cuts inside its
// (palindromic) recognition site.
//
// Site is written 5'->3' with IUPAC ambiguity codes, without any cut
// marker. TopCut is the number of site bases to the LEFT of the cut on
// the top strand (0-based count, so 1 means "after the first base",
// the NEB notation G/AATTC). The cut on the bottom strand is the
// mirror image: BottomCut = len(Site) - TopCut, counted from the left
// on the top strand. This mirror rule holds for palindromic sites,
// which is every enzyme in the built-in database.
type Enzyme struct {
	// Name is the common name, e.g. "EcoRI".
	Name string
	// Site is the recognition sequence, e.g. "GAATTC".
	Site string
	// TopCut is the offset of the top-strand cut within Site.
	TopCut int
	// Organism is the source bacterium, e.g. "Escherichia coli RY13".
	Organism string
}

// Length returns the length of the recognition site in base pairs.
func (e Enzyme) Length() int { return len(e.Site) }

// BottomCut returns the number of site bases to the left of the cut on
// the bottom strand, counted from the left on the top strand.
func (e Enzyme) BottomCut() int { return len(e.Site) - e.TopCut }

// OverhangLen returns the length of the single-stranded overhang in
// nucleotides (0 for blunt cutters).
func (e Enzyme) OverhangLen() int {
	d := e.TopCut - e.BottomCut()
	if d < 0 {
		d = -d
	}
	return d
}

// OverhangType returns "5'", "3'" or "blunt".
func (e Enzyme) OverhangType() string {
	switch {
	case e.TopCut < e.BottomCut():
		return "5'"
	case e.TopCut > e.BottomCut():
		return "3'"
	default:
		return "blunt"
	}
}

// OverhangSeq returns the single-stranded overhang sequence 5'->3' on
// the top strand ("" for blunt cutters).
func (e Enzyme) OverhangSeq() string {
	switch e.OverhangType() {
	case "5'":
		return e.Site[e.TopCut:e.BottomCut()]
	case "3'":
		return e.Site[e.BottomCut():e.TopCut]
	default:
		return ""
	}
}

// Recognition returns the site in NEB notation, with "/" marking the
// top-strand cut, e.g. "G/AATTC".
func (e Enzyme) Recognition() string {
	return e.Site[:e.TopCut] + "/" + e.Site[e.TopCut:]
}

// IsPalindromic reports whether the recognition site reads the same on
// both strands. Every enzyme in the built-in database satisfies this.
func (e Enzyme) IsPalindromic() bool { return dna.IsPalindromic(e.Site) }

// MatchesAt reports whether the recognition site matches the sequence
// starting at 0-based position pos. IUPAC codes in the site match any
// concrete base in their set; N matches anything.
func (e Enzyme) MatchesAt(seq string, pos int) bool {
	if pos < 0 || pos+e.Length() > len(seq) {
		return false
	}
	for i := 0; i < e.Length(); i++ {
		set := dna.IUPACSet[e.Site[i]]
		if set == "" {
			return false
		}
		if !strings.ContainsRune(set, rune(seq[pos+i])) {
			return false
		}
	}
	return true
}

// ExpectedFrequency returns the probability that the site occurs at a
// random position in random DNA (each base equally likely), i.e. the
// product over site positions of (degeneracy / 4). For a 6-cutter
// without ambiguity codes this is 1/4096.
func (e Enzyme) ExpectedFrequency() float64 {
	p := 1.0
	for i := 0; i < e.Length(); i++ {
		set := dna.IUPACSet[e.Site[i]]
		p *= float64(len(set)) / 4.0
	}
	return p
}

// ExpectedCount returns the expected number of sites in a random
// sequence of length seqLen (linear molecule).
func (e Enzyme) ExpectedCount(seqLen int) float64 {
	return float64(seqLen) * e.ExpectedFrequency()
}

// Validate checks the enzyme definition: non-empty name and site, site
// of IUPAC codes only, cut inside the site.
func (e Enzyme) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("enzyme with empty name")
	}
	if e.Site == "" {
		return fmt.Errorf("enzyme %s: empty recognition site", e.Name)
	}
	for i := 0; i < len(e.Site); i++ {
		if _, ok := dna.IUPACSet[e.Site[i]]; !ok {
			return fmt.Errorf("enzyme %s: invalid code %q in site %q", e.Name, e.Site[i], e.Site)
		}
	}
	if e.TopCut < 0 || e.TopCut > len(e.Site) {
		return fmt.Errorf("enzyme %s: TopCut %d outside site of length %d", e.Name, e.TopCut, len(e.Site))
	}
	return nil
}

// ByName returns the enzymes sorted alphabetically by name.
func ByName(enzymes []Enzyme) []Enzyme {
	out := make([]Enzyme, len(enzymes))
	copy(out, enzymes)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Find looks an enzyme up by name (case-insensitive) in the given list.
// The second return value reports whether it was found.
func Find(enzymes []Enzyme, name string) (Enzyme, bool) {
	for _, e := range enzymes {
		if strings.EqualFold(e.Name, name) {
			return e, true
		}
	}
	return Enzyme{}, false
}
