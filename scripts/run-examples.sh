#!/bin/sh
# Regenerate every file in examples/output/ from examples/input/.
# Run from the repository root: ./scripts/run-examples.sh
set -e
cd "$(dirname "$0")/.."

if [ ! -x ./rmap ]; then
  go build -o rmap ./cmd/rmap
fi

mkdir -p examples/output

./rmap enzymes > examples/output/enzymes-list.txt
./rmap scan -in examples/input/puc19.fasta -circular > examples/output/puc19-scan.txt
./rmap scan -in examples/input/puc19.fasta -circular -format json -out examples/output/puc19-scan.json
./rmap digest -in examples/input/puc19.fasta -enzymes EcoRI,HindIII -circular > examples/output/puc19-digest-ecori-hindiii.txt
./rmap digest -in examples/input/puc19.fasta -enzymes EcoRI -circular > examples/output/puc19-digest-ecori.txt
./rmap map -in examples/input/puc19.fasta -circular -width 100 > examples/output/puc19-map.txt
./rmap scan -in examples/input/synthetic-plasmid.fasta -circular > examples/output/synthetic-scan.txt
./rmap digest -in examples/input/synthetic-plasmid.fasta -enzymes EcoRI,BamHI -circular > examples/output/synthetic-digest.txt
./rmap scan -in examples/input/lambda-fragment.fasta > examples/output/lambda-scan.txt
./rmap scan -in examples/input/multi.fasta > examples/output/batch-scan.txt

echo "examples/output/ regenerated."
