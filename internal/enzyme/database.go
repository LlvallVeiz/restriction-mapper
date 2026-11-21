package enzyme

// Database is the built-in catalogue of common type II restriction
// enzymes. Site and TopCut follow the NEB Alphabetized List of
// Recognition Sequences ("/" marks the top-strand cut); every entry
// was checked against that list on 2026-11-20. All 32 enzymes are
// palindromic and cut inside their recognition site, so a single
// TopCut value fully describes the cut on both strands. Source
// organisms follow REBASE. Type IIS enzymes that cut outside their
// site (e.g. BsaI, GGTCTC(1/5)) are deliberately NOT included: they
// need a different cut model, see docs/RESEARCH.md.
var Database = []Enzyme{
	{Name: "EcoRI", Site: "GAATTC", TopCut: 1, Organism: "Escherichia coli RY13"},
	{Name: "BamHI", Site: "GGATCC", TopCut: 1, Organism: "Bacillus amyloliquefaciens H"},
	{Name: "HindIII", Site: "AAGCTT", TopCut: 1, Organism: "Haemophilus influenzae Rd"},
	{Name: "XbaI", Site: "TCTAGA", TopCut: 1, Organism: "Xanthomonas badrii"},
	{Name: "SalI", Site: "GTCGAC", TopCut: 1, Organism: "Streptomyces albus G"},
	{Name: "PstI", Site: "CTGCAG", TopCut: 5, Organism: "Providencia stuartii"},
	{Name: "KpnI", Site: "GGTACC", TopCut: 5, Organism: "Klebsiella pneumoniae OK8"},
	{Name: "SacI", Site: "GAGCTC", TopCut: 5, Organism: "Streptomyces achromogenes"},
	{Name: "SmaI", Site: "CCCGGG", TopCut: 3, Organism: "Serratia marcescens Sb"},
	{Name: "NcoI", Site: "CCATGG", TopCut: 1, Organism: "Nocardia corallina"},
	{Name: "NdeI", Site: "CATATG", TopCut: 2, Organism: "Neisseria denitrificans"},
	{Name: "XhoI", Site: "CTCGAG", TopCut: 1, Organism: "Xanthomonas holcicola"},
	{Name: "SpeI", Site: "ACTAGT", TopCut: 1, Organism: "Sphaerotilus natans"},
	{Name: "NotI", Site: "GCGGCCGC", TopCut: 2, Organism: "Nocardia otitidiscaviarum"},
	{Name: "BglII", Site: "AGATCT", TopCut: 1, Organism: "Bacillus globigii"},
	{Name: "ClaI", Site: "ATCGAT", TopCut: 2, Organism: "Caryophanon latum L"},
	{Name: "EcoRV", Site: "GATATC", TopCut: 3, Organism: "Escherichia coli J62 pLG74"},
	{Name: "HpaI", Site: "GTTAAC", TopCut: 3, Organism: "Haemophilus parainfluenzae"},
	{Name: "SphI", Site: "GCATGC", TopCut: 5, Organism: "Streptomyces phaeochromogenes"},
	{Name: "ApaI", Site: "GGGCCC", TopCut: 5, Organism: "Acetobacter pasteurianus"},
	{Name: "AccI", Site: "GTMKAC", TopCut: 2, Organism: "Acinetobacter calcoaceticus"},
	{Name: "HinfI", Site: "GANTC", TopCut: 1, Organism: "Haemophilus influenzae Rf"},
	{Name: "TaqI", Site: "TCGA", TopCut: 1, Organism: "Thermus aquaticus YT1"},
	{Name: "NheI", Site: "GCTAGC", TopCut: 1, Organism: "Neisseria mucosa var. heidelbergensis"},
	{Name: "AgeI", Site: "ACCGGT", TopCut: 1, Organism: "Agrobacterium gelatinovorum"},
	{Name: "MluI", Site: "ACGCGT", TopCut: 1, Organism: "Micrococcus luteus"},
	{Name: "PvuII", Site: "CAGCTG", TopCut: 3, Organism: "Proteus vulgaris"},
	{Name: "StuI", Site: "AGGCCT", TopCut: 3, Organism: "Streptomyces tubercidicus"},
	{Name: "ScaI", Site: "AGTACT", TopCut: 3, Organism: "Streptomyces caespitosus"},
	{Name: "XmaI", Site: "CCCGGG", TopCut: 1, Organism: "Xanthomonas malvacearum"},
	{Name: "AvaI", Site: "CYCGRG", TopCut: 1, Organism: "Anabaena variabilis"},
	{Name: "AflII", Site: "CTTAAG", TopCut: 1, Organism: "Anabaena flos-aquae"},
}

// Default returns a copy of the full built-in database.
func Default() []Enzyme {
	out := make([]Enzyme, len(Database))
	copy(out, Database)
	return out
}
