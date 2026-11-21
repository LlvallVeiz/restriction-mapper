// Package report renders scan results as human-readable text: a site
// table, per-enzyme statistics and an ASCII restriction map.
package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/LlvallVeiz/restriction-mapper/internal/enzyme"
	"github.com/LlvallVeiz/restriction-mapper/internal/scan"
)

// SitesTable renders all sites as a text table, one row per site.
// Columns: position, enzyme, strand, recognition, top cut, bottom cut.
func SitesTable(sites []scan.Site, enzymes []enzyme.Enzyme) string {
	var b strings.Builder
	b.WriteString("  Pos   Enzyme  Strand  Site      TopCut  BotCut\n")
	b.WriteString("  ----- ------- ------- --------- ------- -------\n")
	byName := map[string]enzyme.Enzyme{}
	for _, e := range enzymes {
		byName[e.Name] = e
	}
	for _, s := range sites {
		e := byName[s.Enzyme]
		fmt.Fprintf(&b, "  %-5d %-7s %-7s %-9s %-7d %-7d\n",
			s.Start, s.Enzyme, s.Strand, e.Recognition(), s.CutTop, s.CutBottom)
	}
	if len(sites) == 0 {
		b.WriteString("  (no sites found)\n")
	}
	return b.String()
}

// StatsTable renders per-enzyme statistics: observed site count,
// expected count in random DNA of the same length and GC, and the
// overhang the enzyme leaves. Enzymes are sorted by observed count
// (descending), then by name.
func StatsTable(sites []scan.Site, enzymes []enzyme.Enzyme, seqLen int) string {
	counts := scan.CountPerEnzyme(sites)
	type row struct {
		e        enzyme.Enzyme
		observed int
		expected float64
	}
	var rows []row
	for _, e := range enzymes {
		rows = append(rows, row{e, counts[e.Name], e.ExpectedCount(seqLen)})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].observed != rows[j].observed {
			return rows[i].observed > rows[j].observed
		}
		return rows[i].e.Name < rows[j].e.Name
	})
	var b strings.Builder
	b.WriteString("  Enzyme  Sites  Expected  Overhang   Site\n")
	b.WriteString("  ------- ------ --------- ---------- ---------\n")
	for _, r := range rows {
		over := r.e.OverhangType()
		if over != "blunt" {
			over = fmt.Sprintf("%s %d nt", over, r.e.OverhangLen())
		}
		fmt.Fprintf(&b, "  %-7s %-6d %-9.2f %-10s %s\n",
			r.e.Name, r.observed, r.expected, over, r.e.Recognition())
	}
	return b.String()
}

// MapASCII draws a simple restriction map: a ruler line plus one row
// per enzyme that cuts, with '|' at each site start (a cut inside the
// site is shown at the site start for readability; exact cut
// positions are in the site table). width is the number of columns
// used for the sequence; positions scale linearly.
func MapASCII(seqLen int, sites []scan.Site, enzymes []enzyme.Enzyme, width int) string {
	if width < 20 {
		width = 20
	}
	if seqLen <= 0 {
		return "(empty sequence)\n"
	}
	byEnzyme := map[string][]scan.Site{}
	for _, s := range sites {
		byEnzyme[s.Enzyme] = append(byEnzyme[s.Enzyme], s)
	}
	var names []string
	for name := range byEnzyme {
		names = append(names, name)
	}
	sort.Strings(names)

	var b strings.Builder
	// Ruler: ticks every 10% of the sequence.
	ruler := make([]byte, width)
	for i := range ruler {
		ruler[i] = '-'
	}
	for f := 0; f <= 10; f++ {
		pos := f * (width - 1) / 10
		ruler[pos] = '+'
	}
	fmt.Fprintf(&b, "  %s\n", string(ruler))
	fmt.Fprintf(&b, "  1%s%d bp\n", strings.Repeat(" ", width-1-len(fmt.Sprint(seqLen))), seqLen)
	for _, name := range names {
		row := make([]byte, width)
		for i := range row {
			row[i] = ' '
		}
		for _, s := range byEnzyme[name] {
			col := (s.Start - 1) * (width - 1) / seqLen
			if col < 0 {
				col = 0
			}
			if col >= width {
				col = width - 1
			}
			row[col] = '|'
		}
		fmt.Fprintf(&b, "  %-7s %s (%d)\n", name, string(row), len(byEnzyme[name]))
	}
	return b.String()
}

// SummaryLine returns a one-line summary of a scan.
func SummaryLine(recordID string, seqLen int, sites []scan.Site, enzymes []enzyme.Enzyme) string {
	nonCutters := scan.NonCutters(sites, enzymes)
	return fmt.Sprintf("%s: %d bp, %d sites from %d enzymes, %d enzymes do not cut",
		recordID, seqLen, len(sites), len(enzymes)-len(nonCutters), len(nonCutters))
}
