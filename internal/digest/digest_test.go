package digest

import (
	"reflect"
	"testing"
)

// Hand-computed digest on the 24 bp test sequence
// "AAAAGAATTCCCCCGGATCCTTTT" (EcoRI cut after base 5 and 9,
// BamHI cut after 15 and 19). Using only the TOP-strand cuts
// [5, 15] (as a single digest per enzyme would give):
// linear fragments: 5, 10, 9.
func TestDigestLinearSingle(t *testing.T) {
	frags := Digest(24, []int{5}, false)
	if len(frags) != 2 {
		t.Fatalf("got %d fragments, want 2", len(frags))
	}
	if frags[0].Length != 5 || frags[1].Length != 19 {
		t.Errorf("lengths = %d, %d, want 5, 19", frags[0].Length, frags[1].Length)
	}
	if frags[0].Start != 1 || frags[0].End != 5 || frags[1].Start != 6 || frags[1].End != 24 {
		t.Errorf("coords = %+v", frags)
	}
	if TotalLength(frags) != 24 {
		t.Errorf("total = %d, want 24", TotalLength(frags))
	}
}

func TestDigestLinearDouble(t *testing.T) {
	// EcoRI + BamHI top cuts [5, 15] -> fragments 5, 10, 9.
	frags := Digest(24, []int{5, 15}, false)
	if got := Lengths(frags); !reflect.DeepEqual(got, []int{10, 9, 5}) {
		t.Errorf("lengths = %v, want [10 9 5]", got)
	}
	// Unsorted / duplicated input must not change the result.
	frags2 := Digest(24, []int{15, 5, 5, 15}, false)
	if got := Lengths(frags2); !reflect.DeepEqual(got, []int{10, 9, 5}) {
		t.Errorf("unsorted/dup lengths = %v, want [10 9 5]", got)
	}
}

func TestDigestLinearNoCut(t *testing.T) {
	frags := Digest(100, nil, false)
	if len(frags) != 1 || frags[0].Length != 100 {
		t.Errorf("no-cut linear = %+v, want 1 fragment of 100", frags)
	}
}

func TestDigestCircular(t *testing.T) {
	// 0 cuts -> whole plasmid.
	if frags := Digest(2686, nil, true); len(frags) != 1 || frags[0].Length != 2686 {
		t.Errorf("circular no-cut = %+v", frags)
	}
	// 1 cut linearises the plasmid: still one fragment of full length.
	if frags := Digest(2686, []int{400}, true); len(frags) != 1 || frags[0].Length != 2686 {
		t.Errorf("circular 1-cut = %+v", frags)
	}
	// 2 cuts at 400 and 1400 -> fragments 1000 and 1686.
	frags := Digest(2686, []int{400, 1400}, true)
	if got := Lengths(frags); !reflect.DeepEqual(got, []int{1686, 1000}) {
		t.Errorf("circular 2-cut lengths = %v, want [1686 1000]", got)
	}
	// The wrapping fragment runs from 1401 to 400 (End < Start).
	var wrap *Fragment
	for i := range frags {
		if frags[i].End < frags[i].Start {
			wrap = &frags[i]
		}
	}
	if wrap == nil {
		t.Fatal("no wrapping fragment found")
	}
	if wrap.Start != 1401 || wrap.End != 400 || wrap.Length != 1686 {
		t.Errorf("wrapping fragment = %+v, want Start 1401 End 400 Length 1686", *wrap)
	}
	// 3 cuts, one at the origin (0): fragments 500, 700, 1486.
	frags = Digest(2686, []int{0, 500, 1200}, true)
	if got := Lengths(frags); !reflect.DeepEqual(got, []int{1486, 700, 500}) {
		t.Errorf("circular origin-cut lengths = %v, want [1486 700 500]", got)
	}
	if TotalLength(frags) != 2686 {
		t.Errorf("total = %d, want 2686", TotalLength(frags))
	}
}

func TestDigestOutOfRangeIgnored(t *testing.T) {
	frags := Digest(100, []int{-5, 50, 101, 200}, false)
	if len(frags) != 2 || frags[0].Length != 50 || frags[1].Length != 50 {
		t.Errorf("out-of-range cuts gave %+v", frags)
	}
}

func TestFillSequences(t *testing.T) {
	seq := "AAAAGAATTCCCCCGGATCCTTTT"
	frags := Digest(24, []int{5, 15}, false)
	FillSequences(frags, seq, false)
	if frags[0].Sequence != "AAAAG" {
		t.Errorf("frag 0 seq = %q, want AAAAG", frags[0].Sequence)
	}
	if frags[1].Sequence != "AATTCCCCCG" {
		t.Errorf("frag 1 seq = %q, want AATTCCCCCG", frags[1].Sequence)
	}
	if frags[2].Sequence != "GATCCTTTT" {
		t.Errorf("frag 2 seq = %q, want GATCCTTTT", frags[2].Sequence)
	}
	// Wrapping circular fragment.
	frags = Digest(10, []int{3, 7}, true)
	FillSequences(frags, "ACGTACGTAC", true)
	// Fragment 1: bases 4..7 = "TACG"; fragment 2 wraps 8..10+1..3 = "TACACG".
	if frags[0].Sequence != "TACG" {
		t.Errorf("circular frag 0 seq = %q, want TACG", frags[0].Sequence)
	}
	if frags[1].Sequence != "TACACG" {
		t.Errorf("wrapping frag seq = %q, want TACACG", frags[1].Sequence)
	}
}
