#!/usr/bin/env python3
"""Generate the synthetic example FASTA files for restriction-mapper.

Deterministic (fixed seed): re-running this script reproduces the
example inputs byte-for-byte.

- synthetic-plasmid.fasta: 3200 bp circular plasmid, random background
  (seed 20261120) with 29 designed MCS sites planted at known
  1-based positions.
- lambda-fragment.fasta: 5000 bp linear lambda-like fragment, random
  background with 5 designed sites.
- multi.fasta: 3 short hand-built records for batch mode.

The designed sites are the ground truth the tests and README quote;
any additional sites the scanner finds on these files come from the
random background, which is expected (see README).
"""
import random

PLASMID_SITES = {
    101: "GAATTC", 201: "GGATCC", 301: "AAGCTT", 401: "TCTAGA",
    501: "GTCGAC", 601: "CTGCAG", 701: "GGTACC", 801: "GAGCTC",
    901: "CCCGGG", 1001: "CCATGG", 1101: "CTCGAG", 1201: "ACTAGT",
    1301: "AGATCT", 1401: "ATCGAT", 1501: "GATATC", 1601: "CATATG",
    1701: "GCATGC", 1801: "GGGCCC", 1901: "GCTAGC", 2001: "ACGCGT",
    2101: "CAGCTG", 2201: "AGGCCT", 2301: "AGTACT", 2401: "CCCGGG",
    2501: "CTTAAG", 2601: "ACCGGT", 2701: "TCGA", 2801: "GATTC",
    2901: "GTAGAC",
}
LAMBDA_SITES = {
    500: "GAATTC", 1500: "GAATTC", 2500: "GGATCC",
    3500: "AAGCTT", 4500: "GAATTC",
}


def plant(length, sites, rng):
    seq = [rng.choice("ACGT") for _ in range(length)]
    for pos, s in sites.items():
        for i, b in enumerate(s):
            seq[pos - 1 + i] = b
    return "".join(seq)


def write_fasta(path, header, seq, width=70):
    with open(path, "w") as f:
        f.write(f">{header}\n")
        for i in range(0, len(seq), width):
            f.write(seq[i : i + width] + "\n")


def main():
    rng = random.Random(20261120)
    plasmid = plant(3200, PLASMID_SITES, rng)
    write_fasta(
        "examples/input/synthetic-plasmid.fasta",
        "pSYN3200 synthetic circular plasmid (designed MCS, 3200 bp)",
        plasmid,
    )
    lam = plant(5000, LAMBDA_SITES, rng)
    write_fasta(
        "examples/input/lambda-fragment.fasta",
        "lambda-frag synthetic lambda-like linear fragment (5000 bp)",
        lam,
    )
    with open("examples/input/multi.fasta", "w") as f:
        f.write(">frag-a\nAAAAGAATTCCCCCGGATCCTTTT\n")
        f.write(">frag-b\nGGATCCGAATTCGGATCC\n")
        f.write(">frag-c\nACGTACGTACGT\n")
    print("generated: synthetic-plasmid (3200 bp), lambda-fragment (5000 bp), multi.fasta")


if __name__ == "__main__":
    main()
