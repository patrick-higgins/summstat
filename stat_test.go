// Copyright 2012 The Summstat Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package summstat

import (
	"math"
	"testing"
)

const (
	epsilon = 0.0000000000000000001
)

type statTest struct {
	count   int
	min     Sample
	max     Sample
	median  float64
	mean    float64
	stddev  float64
	spread  Sample
	samples []Sample
}

var (
	tests = []statTest{
		{ // 0
			samples: []Sample{},
			count:   0,
			min:     0,
			max:     0,
			median:  0,
			mean:    0,
			stddev:  0,
			spread:  0,
		},
		{ // 1
			samples: []Sample{1},
			count:   1,
			min:     1,
			max:     1,
			median:  1,
			mean:    1,
			stddev:  0,
			spread:  0,
		},
		{ // 2
			samples: []Sample{1, 2},
			count:   2,
			min:     1,
			max:     2,
			median:  1.5,
			mean:    1.5,
			stddev:  .5,
			spread:  1,
		},
		{ // 3
			samples: []Sample{1, 2, 3},
			count:   3,
			min:     1,
			max:     3,
			median:  2,
			mean:    2,
			stddev:  0.8164965809277263,
			spread:  2,
		},
		{ // 4
			samples: []Sample{0, 1, 2, 3, 4, 5},
			count:   6,
			min:     0,
			max:     5,
			median:  2.5,
			mean:    2.5,
			stddev:  1.707825127659933,
			spread:  5,
		},
		{ // 5
			samples: []Sample{-10, -9, -8, -7},
			count:   4,
			min:     -10,
			max:     -7,
			median:  -8.5,
			mean:    -8.5,
			stddev:  1.118033988749895,
			spread:  3,
		},
		{ // 6
			samples: []Sample{-1, 1},
			count:   2,
			min:     -1,
			max:     1,
			median:  0,
			mean:    0,
			stddev:  1,
			spread:  2,
		},
		{ // 7
			samples: []Sample{-1, 0, 1},
			count:   3,
			min:     -1,
			max:     1,
			median:  0,
			mean:    0,
			stddev:  0.816496580927726,
			spread:  2,
		},
	}
)

func insertSamples(s *Stats, samples []Sample) {
	for _, sample := range samples {
		s.AddSample(sample)
	}
}

func Test(t *testing.T) {
	for i, test := range tests {
		s := NewStats()
		insertSamples(s, test.samples)
		if s.Count() != test.count {
			t.Errorf("[%d] Invalid count: %d, expected: %d", i, s.Count(), test.count)
		}
		if s.Min() != test.min {
			t.Errorf("[%d] Invalid min: %v, expected: %v", i, s.Min(), test.min)
		}
		if s.Max() != test.max {
			t.Errorf("[%d] Invalid max: %v, expected: %v", i, s.Max(), test.max)
		}
		if s.Median() != test.median {
			t.Errorf("[%d] Invalid median: %v, expected: %v", i, s.Median(), test.median)
		}
		if math.Abs(s.Mean()-test.mean) > epsilon {
			t.Errorf("[%d] Invalid mean: %f, expected: %f", i, s.Mean(), test.mean)
		}
		if math.Abs(s.Stddev()-test.stddev) > epsilon {
			t.Errorf("[%d] Invalid stddev: %v, expected: %v", i, s.Stddev(), test.stddev)
		}
		if s.Spread() != test.spread {
			t.Errorf("[%d] Invalid spread: %v, expected: %v", i, s.Spread(), test.spread)
		}
	}
}

func chkPct(t *testing.T, s *Stats, pct float64, exp Sample) {
	if s.Percentile(pct) != exp {
		t.Errorf("%.1f%% != %v: %v", 100*pct, exp, s.Percentile(pct))
	}
}

func TestPercentile(t *testing.T) {
	s := NewStats()
	insertSamples(s, []Sample{0, 1, 10, 25, 100})
	chkPct(t, s, 0, 0)
	chkPct(t, s, .25, 1)
	chkPct(t, s, .5, 10)
	chkPct(t, s, .75, 25)
	chkPct(t, s, 1, 100)

	s = NewStats()
	insertSamples(s, []Sample{25})
	chkPct(t, s, 0, 25)
	chkPct(t, s, .25, 25)
	chkPct(t, s, .5, 25)
	chkPct(t, s, .75, 25)
	chkPct(t, s, 1, 25)

	s = NewStats()
	insertSamples(s, []Sample{1, 2})
	chkPct(t, s, 0, 1)
	chkPct(t, s, .25, 1)
	chkPct(t, s, .5, 2)
	chkPct(t, s, .75, 2)
	chkPct(t, s, 1, 2)
}

type binTest struct {
	samples  []Sample
	binCount int
	low      Sample
	high     Sample
	bins     []Sample
}

var binTests = []binTest{
	{
		samples:  []Sample{1, 2, 3},
		binCount: 3,
		low:      1,
		high:     3,
		bins:     []Sample{1, 3, math.MaxFloat64},
	},
	{
		samples:  []Sample{1, 2, 3, 4},
		binCount: 4,
		low:      1,
		high:     4,
		bins:     []Sample{1, 2.5, 4, math.MaxFloat64},
	},
	{
		samples:  []Sample{1, 2, 3, 4},
		binCount: 5,
		low:      1,
		high:     4,
		bins:     []Sample{1, 2, 3, 4, math.MaxFloat64},
	},
	{
		samples:  []Sample{-100, -75, -50},
		binCount: 4,
		low:      -100,
		high:     -50,
		bins:     []Sample{-100, -75, -50, math.MaxFloat64},
	},
}

type discardBinTest struct {
	samples  []Sample
	binCount int
	discard  float64
	bins     []Sample
}

var discardBinTests = []discardBinTest{
	{
		samples: []Sample{
			-1000000,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			1000000,
		},
		binCount: 11,
		discard:  0.01,
		bins:     []Sample{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, math.MaxFloat64},
	},
	{
		samples: []Sample{
			-1000000,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			1000000,
		},
		binCount: 11,
		discard:  0.01,
		bins:     []Sample{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, math.MaxFloat64},
	},
}

func TestDiscardBins(t *testing.T) {
	for _, test := range discardBinTests {
		s := NewStats()
		insertSamples(s, test.samples)
		s.CreateBinsDiscard(test.binCount, test.discard)
		if len(s.bins) != len(test.bins) {
			t.Errorf("not enough bins: %d, expected %d", len(s.bins), len(test.bins))
		}
		for i, bin := range s.bins {
			if bin != test.bins[i] {
				t.Errorf("s.bins[%d] != %v, expected %v", i, bin, test.bins[i])
			}
		}
	}
}
