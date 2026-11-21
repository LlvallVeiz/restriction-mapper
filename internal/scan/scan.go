// Package scan finds restriction sites in a DNA sequence. It scans
// both strands (forward sequence and its reverse complement), reports
// overlapping sites, and handles circular molecules where a site may
// span the origin. Positions are 1-based on the top strand.
package scan

import (
	"sort"

	"github.com/LlvallVeiz/restriction-mapper/internal/dna"
	"github.com/LlvallVeiz/restriction-mapper/internal/enzyme"
)

// Site is one occurrence of a restriction site.
type Site struct {
	// Enzyme is the enzyme name, e.g. "EcoRI".
	Enzyme string
	// Start is the 1-based position of the first site base on the top strand.
	Start int
	// Strand is "+" if the site was found on the top strand, "-" if on
	// the bottom strand (read as its reverse complement).
	Strand string
	// CutTop is the number of bases from the 5' end of the molecule to
	// the left of the top-strand cut (so the cut falls between bases
	// CutTop and CutTop+1, 1-based). For a circular molecule it is
	// taken modulo the sequence length.
	CutTop int
	// CutBottom is the same count for the bottom-strand cut.
	CutBottom int
}

// End returns the 1-based position of the last site base on the top
// strand (wrapping for circular molecules is handled by callers).
func (s Site) End(siteLen int) int { return s.Start + siteLen - 1 }

// Scan finds all sites of the given enzymes in seq. seq must already
// be cleaned (upper case, no white space). If circular is true, sites
// spanning the origin are found and Start wraps into 1..len(seq).
// Results are sorted by Start, then by enzyme name. Palindromic sites
// are reported once (strand "+"): a reverse-complement hit at the same
// position is the same physical site, not a second one.
func Scan(seq string, enzymes []enzyme.Enzyme, circular bool) []Site {
	n := len(seq)
	if n == 0 {
		return nil
	}
	maxLen := 0
	for _, e := range enzymes {
		if e.Length() > maxLen {
			maxLen = e.Length()
		}
	}
	// For circular molecules, append a prefix so sites spanning the
	// origin are visible; only starts < n are reported.
	search := seq
	if circular && maxLen > 1 {
		search = seq + seq[:maxLen-1]
	}
	rc := dna.ReverseComplement(seq)
	rcSearch := rc
	if circular && maxLen > 1 {
		rcSearch = rc + rc[:maxLen-1]
	}

	seen := map[[2]interface{}]bool{}
	var sites []Site
	add := func(e enzyme.Enzyme, start0 int, strand string) {
		// start0 is 0-based on the top strand.
		key := [2]interface{}{e.Name, start0}
		if seen[key] {
			return
		}
		seen[key] = true
		top := start0 + e.TopCut         // bases left of top cut
		bottom := start0 + e.BottomCut() // bases left of bottom cut
		if circular {
			top = mod(top, n)
			bottom = mod(bottom, n)
		}
		sites = append(sites, Site{
			Enzyme:    e.Name,
			Start:     start0 + 1,
			Strand:    strand,
			CutTop:    top,
			CutBottom: bottom,
		})
	}
	for _, e := range enzymes {
		L := e.Length()
		for i := 0; i+L <= len(search); i++ {
			if i >= n {
				break // only starts inside the molecule
			}
			if e.MatchesAt(search, i) {
				add(e, i, "+")
			}
		}
		// Bottom strand: a match at rc index j starts at top-strand
		// 0-based position n-(j+L).
		for j := 0; j+L <= len(rcSearch); j++ {
			if j >= n {
				break
			}
			if e.MatchesAt(rcSearch, j) {
				start0 := n - (j + L)
				if circular {
					start0 = mod(start0, n)
				}
				if start0 < 0 {
					continue
				}
				add(e, start0, "-")
			}
		}
	}
	sort.Slice(sites, func(i, j int) bool {
		if sites[i].Start != sites[j].Start {
			return sites[i].Start < sites[j].Start
		}
		return sites[i].Enzyme < sites[j].Enzyme
	})
	return sites
}

func mod(a, n int) int {
	m := a % n
	if m < 0 {
		m += n
	}
	return m
}

// SitesFor returns only the sites of one enzyme (by name).
func SitesFor(sites []Site, name string) []Site {
	var out []Site
	for _, s := range sites {
		if s.Enzyme == name {
			out = append(out, s)
		}
	}
	return out
}

// CountPerEnzyme counts sites per enzyme name.
func CountPerEnzyme(sites []Site) map[string]int {
	counts := map[string]int{}
	for _, s := range sites {
		counts[s.Enzyme]++
	}
	return counts
}

// NonCutters returns the names of enzymes with no site in the sequence.
func NonCutters(sites []Site, enzymes []enzyme.Enzyme) []string {
	counts := CountPerEnzyme(sites)
	var out []string
	for _, e := range enzymes {
		if counts[e.Name] == 0 {
			out = append(out, e.Name)
		}
	}
	sort.Strings(out)
	return out
}

// TopCutPositions returns the sorted, de-duplicated top-strand cut
// positions (bases left of the cut) of the given sites. These are the
// cut coordinates a digest consumes.
func TopCutPositions(sites []Site) []int {
	seen := map[int]bool{}
	var cuts []int
	for _, s := range sites {
		if !seen[s.CutTop] {
			seen[s.CutTop] = true
			cuts = append(cuts, s.CutTop)
		}
	}
	sort.Ints(cuts)
	return cuts
}
