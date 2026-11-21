// Command rmap maps restriction-enzyme sites on DNA sequences and
// simulates digests.
//
// Usage:
//
//	rmap enzymes [-format table|json]
//	rmap scan    -in FILE [-enzymes EcoRI,BamHI] [-all] [-circular] [-format text|table|json] [-out FILE]
//	rmap digest  -in FILE -enzymes EcoRI,HindIII [-circular] [-seq] [-format text|json]
//	rmap map     -in FILE [-enzymes EcoRI,BamHI] [-all] [-circular] [-width N]
//
// Positions are 1-based on the top strand. A FASTA file may hold
// several records; scan and map process each one, digest processes
// the first record unless -id selects another.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/LlvallVeiz/restriction-mapper/internal/digest"
	"github.com/LlvallVeiz/restriction-mapper/internal/enzyme"
	"github.com/LlvallVeiz/restriction-mapper/internal/fasta"
	"github.com/LlvallVeiz/restriction-mapper/internal/report"
	"github.com/LlvallVeiz/restriction-mapper/internal/scan"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "enzymes":
		err = cmdEnzymes(os.Args[2:])
	case "scan":
		err = cmdScan(os.Args[2:])
	case "digest":
		err = cmdDigest(os.Args[2:])
	case "map":
		err = cmdMap(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	case "-v", "--version", "version":
		fmt.Printf("rmap %s\n", version)
	default:
		fmt.Fprintf(os.Stderr, "rmap: unknown subcommand %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "rmap: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Print(`rmap ` + version + ` — restriction site mapper

Usage:
  rmap enzymes [-format table|json]
      List the built-in enzyme database (32 type II enzymes).
  rmap scan -in FILE [-enzymes NAME,...] [-all] [-circular]
            [-format text|table|json] [-out FILE]
      Find restriction sites in a FASTA file (multi-record = batch).
  rmap digest -in FILE -enzymes NAME,... [-circular] [-seq]
            [-format text|json] [-id RECORD]
      Simulate a single or double digest and list the fragments.
  rmap map -in FILE [-enzymes NAME,...] [-all] [-circular] [-width N]
      Draw an ASCII restriction map for each record.
  rmap version
      Print the version.

Flags:
  -in FILE        input FASTA file (required for scan/digest/map)
  -enzymes LIST   comma-separated enzyme names (default: all enzymes
                  for scan/map; required for digest)
  -all            scan/map with every enzyme in the database
  -circular       treat the sequence as circular (plasmid): sites may
                  span the origin, and digests wrap around
  -seq            include fragment sequences in digest output
  -format FORMAT  output format (see each subcommand)
  -out FILE       write output to FILE instead of stdout
  -width N        map width in columns (default 100)
  -id RECORD      digest only this FASTA record ID (default: first)
`)
}

// resolveEnzymes turns a comma-separated -enzymes value into enzymes.
// An empty value means "all enzymes"; unknown names are an error.
func resolveEnzymes(value string) ([]enzyme.Enzyme, error) {
	if strings.TrimSpace(value) == "" {
		return enzyme.Default(), nil
	}
	var out []enzyme.Enzyme
	for _, name := range strings.Split(value, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		e, ok := enzyme.Find(enzyme.Default(), name)
		if !ok {
			return nil, fmt.Errorf("unknown enzyme %q (see `rmap enzymes`)", name)
		}
		out = append(out, e)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no enzymes selected")
	}
	return out, nil
}

func readFASTA(path string) ([]fasta.Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return fasta.Parse(f)
}

func writeOut(path, content string) error {
	if path == "" {
		fmt.Print(content)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// --- enzymes -----------------------------------------------------------

func cmdEnzymes(args []string) error {
	fs := flag.NewFlagSet("enzymes", flag.ContinueOnError)
	format := fs.String("format", "table", "output format: table|json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	enzymes := enzyme.ByName(enzyme.Default())
	switch *format {
	case "table":
		var b strings.Builder
		b.WriteString("  Enzyme  Site      Length  Ends          Organism\n")
		b.WriteString("  ------- --------- ------- ------------- --------------------------------\n")
		for _, e := range enzymes {
			ends := e.OverhangType()
			if ends != "blunt" {
				ends = fmt.Sprintf("%s %d nt", ends, e.OverhangLen())
			}
			fmt.Fprintf(&b, "  %-7s %-9s %-7d %-13s %s\n",
				e.Name, e.Recognition(), e.Length(), ends, e.Organism)
		}
		fmt.Fprintf(&b, "\n  %d enzymes. Sites and cut positions: NEB / REBASE.\n", len(enzymes))
		fmt.Print(b.String())
	case "json":
		return json.NewEncoder(os.Stdout).Encode(enzymes)
	default:
		return fmt.Errorf("unknown format %q (want table|json)", *format)
	}
	return nil
}

// --- scan --------------------------------------------------------------

// scanResult is the JSON shape of one record's scan.
type scanResult struct {
	Record     string      `json:"record"`
	Length     int         `json:"length_bp"`
	GCPercent  float64     `json:"gc_percent"`
	Circular   bool        `json:"circular"`
	Sites      []scan.Site `json:"sites"`
	NonCutters []string    `json:"non_cutters"`
}

func cmdScan(args []string) error {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	in := fs.String("in", "", "input FASTA file")
	enz := fs.String("enzymes", "", "comma-separated enzyme names (default: all)")
	all := fs.Bool("all", false, "use all enzymes")
	circular := fs.Bool("circular", false, "circular DNA")
	format := fs.String("format", "text", "output format: text|table|json")
	out := fs.String("out", "", "output file (default: stdout)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *in == "" {
		return fmt.Errorf("scan requires -in FILE")
	}
	if *all {
		*enz = ""
	}
	enzymes, err := resolveEnzymes(*enz)
	if err != nil {
		return err
	}
	records, err := readFASTA(*in)
	if err != nil {
		return err
	}
	switch *format {
	case "text", "table":
		var b strings.Builder
		for _, rec := range records {
			sites := scan.Scan(rec.Sequence, enzymes, *circular)
			fmt.Fprintf(&b, "== %s (%d bp, GC %.2f%%, %s) ==\n",
				rec.ID, rec.Length(), 100*rec.GC(), topology(*circular))
			b.WriteString(report.SummaryLine(rec.ID, rec.Length(), sites, enzymes) + "\n")
			b.WriteString(report.SitesTable(sites, enzymes))
			b.WriteString("Statistics (observed vs expected in random DNA):\n")
			b.WriteString(report.StatsTable(sites, enzymes, rec.Length()))
			nc := scan.NonCutters(sites, enzymes)
			fmt.Fprintf(&b, "Non-cutters (%d): %s\n\n", len(nc), strings.Join(nc, ", "))
		}
		return writeOut(*out, b.String())
	case "json":
		var results []scanResult
		for _, rec := range records {
			sites := scan.Scan(rec.Sequence, enzymes, *circular)
			if sites == nil {
				sites = []scan.Site{}
			}
			nc := scan.NonCutters(sites, enzymes)
			if nc == nil {
				nc = []string{}
			}
			results = append(results, scanResult{
				Record:     rec.ID,
				Length:     rec.Length(),
				GCPercent:  100 * rec.GC(),
				Circular:   *circular,
				Sites:      sites,
				NonCutters: nc,
			})
		}
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return err
		}
		return writeOut(*out, string(data)+"\n")
	default:
		return fmt.Errorf("unknown format %q (want text|table|json)", *format)
	}
}

// --- digest ------------------------------------------------------------

// digestResult is the JSON shape of a digest.
type digestResult struct {
	Record    string            `json:"record"`
	Length    int               `json:"length_bp"`
	Enzymes   []string          `json:"enzymes"`
	Circular  bool              `json:"circular"`
	Cuts      []int             `json:"cut_positions"`
	Fragments []digest.Fragment `json:"fragments"`
}

func cmdDigest(args []string) error {
	fs := flag.NewFlagSet("digest", flag.ContinueOnError)
	in := fs.String("in", "", "input FASTA file")
	enz := fs.String("enzymes", "", "comma-separated enzyme names (1 or 2)")
	circular := fs.Bool("circular", false, "circular DNA")
	withSeq := fs.Bool("seq", false, "include fragment sequences")
	format := fs.String("format", "text", "output format: text|json")
	out := fs.String("out", "", "output file (default: stdout)")
	id := fs.String("id", "", "FASTA record ID (default: first record)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *in == "" {
		return fmt.Errorf("digest requires -in FILE")
	}
	if strings.TrimSpace(*enz) == "" {
		return fmt.Errorf("digest requires -enzymes NAME[,NAME]")
	}
	enzymes, err := resolveEnzymes(*enz)
	if err != nil {
		return err
	}
	if len(enzymes) > 2 {
		return fmt.Errorf("digest supports at most 2 enzymes (single or double digest), got %d", len(enzymes))
	}
	records, err := readFASTA(*in)
	if err != nil {
		return err
	}
	rec := records[0]
	if *id != "" {
		found := false
		for _, r := range records {
			if r.ID == *id {
				rec = r
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("record %q not found in %s", *id, *in)
		}
	}
	sites := scan.Scan(rec.Sequence, enzymes, *circular)
	cuts := scan.TopCutPositions(sites)
	frags := digest.Digest(rec.Length(), cuts, *circular)
	if *withSeq {
		digest.FillSequences(frags, rec.Sequence, *circular)
	}
	names := make([]string, len(enzymes))
	for i, e := range enzymes {
		names[i] = e.Name
	}
	switch *format {
	case "text":
		var b strings.Builder
		fmt.Fprintf(&b, "Digest of %s (%d bp, %s) with %s\n",
			rec.ID, rec.Length(), topology(*circular), strings.Join(names, " + "))
		if len(sites) == 0 {
			b.WriteString("No sites found: the DNA is not cut (1 fragment, uncut).\n")
		} else {
			for _, e := range enzymes {
				n := len(scan.SitesFor(sites, e.Name))
				fmt.Fprintf(&b, "  %s: %d site(s)\n", e.Name, n)
			}
		}
		fmt.Fprintf(&b, "Cut positions (bases left of cut): %v\n", cuts)
		fmt.Fprintf(&b, "Fragments (%d):\n", len(frags))
		b.WriteString("  #    Start   End     Length\n")
		b.WriteString("  ----  ------  ------  ------\n")
		for _, f := range frags {
			fmt.Fprintf(&b, "  %-4d  %-6d  %-6d  %d bp", f.Index, f.Start, f.End, f.Length)
			if *withSeq {
				fmt.Fprintf(&b, "  %s", f.Sequence)
			}
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, "Total: %d bp (check: %d bp)\n", digest.TotalLength(frags), rec.Length())
		return writeOut(*out, b.String())
	case "json":
		data, err := json.MarshalIndent(digestResult{
			Record:    rec.ID,
			Length:    rec.Length(),
			Enzymes:   names,
			Circular:  *circular,
			Cuts:      cuts,
			Fragments: frags,
		}, "", "  ")
		if err != nil {
			return err
		}
		return writeOut(*out, string(data)+"\n")
	default:
		return fmt.Errorf("unknown format %q (want text|json)", *format)
	}
}

// --- map ---------------------------------------------------------------

func cmdMap(args []string) error {
	fs := flag.NewFlagSet("map", flag.ContinueOnError)
	in := fs.String("in", "", "input FASTA file")
	enz := fs.String("enzymes", "", "comma-separated enzyme names (default: all)")
	all := fs.Bool("all", false, "use all enzymes")
	circular := fs.Bool("circular", false, "circular DNA")
	width := fs.Int("width", 100, "map width in columns")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *in == "" {
		return fmt.Errorf("map requires -in FILE")
	}
	if *all {
		*enz = ""
	}
	enzymes, err := resolveEnzymes(*enz)
	if err != nil {
		return err
	}
	records, err := readFASTA(*in)
	if err != nil {
		return err
	}
	for _, rec := range records {
		sites := scan.Scan(rec.Sequence, enzymes, *circular)
		fmt.Printf("== Restriction map: %s (%d bp, %s) ==\n", rec.ID, rec.Length(), topology(*circular))
		fmt.Print(report.MapASCII(rec.Length(), sites, enzymes, *width))
		fmt.Println()
	}
	return nil
}

func topology(circular bool) string {
	if circular {
		return "circular"
	}
	return "linear"
}
