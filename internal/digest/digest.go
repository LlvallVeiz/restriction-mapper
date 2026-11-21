// Package digest simulates restriction digests: it turns a sorted list
// of cut positions into the fragments a gel would show, for linear and
// circular DNA, single or double digests.
package digest

import "sort"

// Fragment is one piece of DNA produced by a digest.
type Fragment struct {
	// Index is the fragment number in 5'->3' order along the top strand.
	Index int
	// Start is the 1-based position of the first base (top strand).
	Start int
	// End is the 1-based position of the last base (top strand). For a
	// circular fragment that wraps the origin, End < Start.
	End int
	// Length is the fragment length in base pairs.
	Length int
	// Sequence holds the fragment sequence when requested ("" otherwise).
	Sequence string
}

// Digest cuts a molecule of length seqLen at the given cut positions
// and returns the fragments. Each cut position is the number of bases
// to the left of the cut (0..seqLen). Positions are de-duplicated and
// sorted; out-of-range positions are ignored.
//
//   - Linear, no cuts: one fragment, the whole molecule.
//   - Linear, k cuts: k+1 fragments.
//   - Circular, 0 or 1 cuts: one fragment, the whole molecule
//     (a single cut linearises a plasmid without splitting it).
//   - Circular, k >= 2 cuts: k fragments, the last one wrapping the origin.
func Digest(seqLen int, cuts []int, circular bool) []Fragment {
	if seqLen <= 0 {
		return nil
	}
	seen := map[int]bool{}
	var cs []int
	for _, c := range cuts {
		if c < 0 || c > seqLen {
			continue
		}
		// A cut exactly at the end of a linear molecule (c == seqLen)
		// produces no fragment; drop it for linear DNA.
		if !circular && c == seqLen {
			continue
		}
		if !seen[c] {
			seen[c] = true
			cs = append(cs, c)
		}
	}
	sort.Ints(cs)

	if circular {
		return digestCircular(seqLen, cs)
	}
	return digestLinear(seqLen, cs)
}

func digestLinear(seqLen int, cuts []int) []Fragment {
	bounds := append([]int{0}, cuts...)
	bounds = append(bounds, seqLen)
	frags := make([]Fragment, 0, len(bounds)-1)
	for i := 0; i+1 < len(bounds); i++ {
		frags = append(frags, Fragment{
			Index:  i + 1,
			Start:  bounds[i] + 1,
			End:    bounds[i+1],
			Length: bounds[i+1] - bounds[i],
		})
	}
	return frags
}

func digestCircular(seqLen int, cuts []int) []Fragment {
	if len(cuts) <= 1 {
		return []Fragment{{Index: 1, Start: 1, End: seqLen, Length: seqLen}}
	}
	frags := make([]Fragment, 0, len(cuts))
	for i, c := range cuts {
		next := cuts[(i+1)%len(cuts)]
		length := next - c
		if length <= 0 {
			length += seqLen
		}
		start := c + 1
		if start > seqLen {
			start -= seqLen
		}
		end := next
		if end == 0 {
			end = seqLen
		}
		frags = append(frags, Fragment{
			Index:  i + 1,
			Start:  start,
			End:    end,
			Length: length,
		})
	}
	return frags
}

// FillSequences fills Fragment.Sequence for each fragment from seq.
// For circular fragments that wrap the origin the sequence is taken
// across the junction. Fragments with inconsistent coordinates are
// left with an empty Sequence.
func FillSequences(frags []Fragment, seq string, circular bool) {
	n := len(seq)
	for i := range frags {
		f := &frags[i]
		if f.Length <= 0 || f.Length > n {
			continue
		}
		if !circular || f.End >= f.Start {
			if f.Start >= 1 && f.End <= n {
				f.Sequence = seq[f.Start-1 : f.End]
			}
			continue
		}
		// Wrapping fragment: Start..n then 1..End.
		if f.Start >= 1 && f.Start <= n && f.End >= 1 && f.End <= n {
			f.Sequence = seq[f.Start-1:] + seq[:f.End]
		}
	}
}

// Lengths returns just the fragment lengths, longest first.
func Lengths(frags []Fragment) []int {
	out := make([]int, len(frags))
	for i, f := range frags {
		out[i] = f.Length
	}
	sort.Sort(sort.Reverse(sort.IntSlice(out)))
	return out
}

// TotalLength sums fragment lengths (should equal the molecule length).
func TotalLength(frags []Fragment) int {
	t := 0
	for _, f := range frags {
		t += f.Length
	}
	return t
}
