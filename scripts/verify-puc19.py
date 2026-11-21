#!/usr/bin/env python3
"""Independent cross-check for rmap: re-implements IUPAC site counting
in plain Python (no shared code with the Go tool) and compares the
per-enzyme site counts on pUC19 against the Go scan output.

Also verifies the EcoRI+HindIII double digest of circular pUC19 by
recomputing fragment lengths from the cut positions.
Usage: python3 verify-puc19.py <puc19.fasta> <go-scan-json>
"""
import json
import sys

IUPAC = {
    "A": "A", "C": "C", "G": "G", "T": "T",
    "R": "AG", "Y": "CT", "S": "GC", "W": "AT",
    "K": "GT", "M": "AC", "B": "CGT", "D": "AGT",
    "H": "ACT", "V": "ACG", "N": "ACGT",
}

# (name, site, top_cut) — transcribed from NEB, same source as the Go DB.
ENZYMES = [
    ("EcoRI", "GAATTC", 1), ("BamHI", "GGATCC", 1), ("HindIII", "AAGCTT", 1),
    ("XbaI", "TCTAGA", 1), ("SalI", "GTCGAC", 1), ("PstI", "CTGCAG", 5),
    ("KpnI", "GGTACC", 5), ("SacI", "GAGCTC", 5), ("SmaI", "CCCGGG", 3),
    ("NcoI", "CCATGG", 1), ("NdeI", "CATATG", 2), ("XhoI", "CTCGAG", 1),
    ("SpeI", "ACTAGT", 1), ("NotI", "GCGGCCGC", 2), ("BglII", "AGATCT", 1),
    ("ClaI", "ATCGAT", 2), ("EcoRV", "GATATC", 3), ("HpaI", "GTTAAC", 3),
    ("SphI", "GCATGC", 5), ("ApaI", "GGGCCC", 5), ("AccI", "GTMKAC", 2),
    ("HinfI", "GANTC", 1), ("TaqI", "TCGA", 1), ("NheI", "GCTAGC", 1),
    ("AgeI", "ACCGGT", 1), ("MluI", "ACGCGT", 1), ("PvuII", "CAGCTG", 3),
    ("StuI", "AGGCCT", 3), ("ScaI", "AGTACT", 3), ("XmaI", "CCCGGG", 1),
    ("AvaI", "CYCGRG", 1), ("AflII", "CTTAAG", 1),
]

COMP = str.maketrans("ACGTURYSWKMBDHVN", "TGCAAYRWSMKHVBND")


def revcomp(s):
    return s.translate(COMP)[::-1]


def read_fasta(path):
    seq = []
    with open(path) as f:
        for line in f:
            if line.startswith(">"):
                continue
            seq.append(line.strip())
    return "".join(seq).upper()


def count_sites_circular(seq, site):
    """Count (possibly wrapping) occurrences on both strands, deduplicated."""
    n = len(seq)
    ext = seq + seq[: len(site) - 1]
    rcext = revcomp(seq) + revcomp(seq)[: len(site) - 1]
    found = set()
    for strand, s in (("+", ext), ("-", rcext)):
        for i in range(n):
            window = s[i : i + len(site)]
            if len(window) < len(site):
                continue
            if all(b in IUPAC[p] for b, p in zip(window, site)):
                if strand == "+":
                    found.add(i + 1)
                else:
                    found.add((n - (i + len(site))) % n + 1)
    return sorted(found)


def main():
    fasta_path, json_path = sys.argv[1], sys.argv[2]
    seq = read_fasta(fasta_path)
    print(f"pUC19 length: {len(seq)} bp")
    assert len(seq) == 2686, "pUC19 must be 2686 bp"

    with open(json_path) as f:
        go = json.load(f)[0]
    go_counts = {}
    go_starts = {}
    for site in go["sites"]:
        go_counts[site["Enzyme"]] = go_counts.get(site["Enzyme"], 0) + 1
        go_starts.setdefault(site["Enzyme"], []).append(site["Start"])

    mismatches = 0
    total_py = 0
    for name, site, _ in ENZYMES:
        py = count_sites_circular(seq, site)
        total_py += len(py)
        g = go_counts.get(name, 0)
        if len(py) != g or sorted(py) != sorted(go_starts.get(name, [])):
            print(f"MISMATCH {name}: python {py} vs go {sorted(go_starts.get(name, []))}")
            mismatches += 1
    print(f"python total sites: {total_py}, go total sites: {len(go['sites'])}")
    print(f"non-cutters (python): {sorted(n for n, _, _ in ENZYMES if go_counts.get(n, 0) == 0)}")

    # Double digest EcoRI + HindIII on circular pUC19, from python cuts.
    cuts = set()
    for name, site, top in ENZYMES:
        if name not in ("EcoRI", "HindIII"):
            continue
        for start in count_sites_circular(seq, site):
            cuts.add((start - 1 + top) % len(seq))
    cuts = sorted(cuts)
    print(f"EcoRI+HindIII top cuts: {cuts}")
    frags = []
    for i, c in enumerate(cuts):
        nxt = cuts[(i + 1) % len(cuts)]
        frags.append((nxt - c) % len(seq) or len(seq))
    print(f"EcoRI+HindIII circular fragments (sorted desc): {sorted(frags, reverse=True)}")
    print(f"sum check: {sum(frags)} == {len(seq)}")

    if mismatches == 0 and total_py == len(go["sites"]):
        print("VERIFY OK: python and go agree on all 32 enzymes.")
    else:
        print(f"VERIFY FAILED: {mismatches} mismatches.")
        sys.exit(1)


if __name__ == "__main__":
    main()
