// Copyright 2012 The Summstat Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package summstat

import (
	"math"
	"sort"
)

// The type of samples we track statistics for
type Sample float64

// Implements sort.Interface so that we can sort samples for percentiles
type sampleSlice []Sample

func (s sampleSlice) Len() int {
	return len(s)
}

func (s sampleSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sampleSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// A Stats represents descriptive statistics about Samples which are being
// added incrementally.
type Stats struct {
	count     int
	sum       Sample
	sum2      Sample
	max       Sample
	min       Sample
	samples   []Sample
	sorted    bool
	bins      []Sample
	binCounts []int
}

// NewStats returns a new Stats
func NewStats() *Stats {
	return &Stats{
		max: -math.MaxFloat64,
		min: math.MaxFloat64,
	}
}

// AddSample adds a sample value and updates the statistics.
func (s *Stats) AddSample(val Sample) {
	s.count++
	s.sum += val
	s.sum2 += val * val
	if val > s.max {
		s.max = val
	}
	if val < s.min {
		s.min = val
	}
	if len(s.bins) > 0 {
		// TODO: use faster lookup method for large bin counts
		for bin, binVal := range s.bins {
			if val <= binVal {
				s.binCounts[bin]++
				break
			}
		}
	} else {
		s.samples = append(s.samples, val)
		s.sorted = false
	}
}

// Count returns the number of samples added.
func (s Stats) Count() int {
	return s.count
}

// Min returns minimal sample value added.
func (s Stats) Min() Sample {
	if s.min > s.max {
		return 0
	}
	return s.min
}

// Max returns the maximal sample value added.
func (s Stats) Max() Sample {
	if s.min > s.max {
		return 0
	}
	return s.max
}

func (s *Stats) sortSamples() {
	if !s.sorted {
		sort.Sort(sampleSlice(s.samples))
		s.sorted = true
	}
}

// Percentile returns the sample value at the given percentile.
//
// It may not be called after CreateBins, which discards the samples from
// which the percentile is calculated.
func (s Stats) Percentile(pct float64) Sample {
	if len(s.bins) > 0 {
		panic("cannot call Percentile() after CreateBins()")
	}
	if len(s.samples) == 0 {
		return 0
	}
	if pct < 0 {
		panic("pct too small")
	}
	if pct > 1 {
		panic("pct too large")
	}
	s.sortSamples()
	// scale pct into int in [0, len-1]
	// Adding 0.5 turns the implicit floor operation of int() into a rounding operation
	i := int(float64(len(s.samples)-1)*pct + 0.5)
	return s.samples[i]
}

// Median returns the median of the samples.
//
// It may not be called after CreateBins, which discards the samples from
// which the percentile is calculated.
func (s Stats) Median() float64 {
	if len(s.bins) > 0 {
		panic("cannot call Percentile() after CreateBins()")
	}
	l := len(s.samples)
	if l == 0 {
		return 0
	}
	s.sortSamples()
	half, rem := l/2, l%2
	if rem == 0 {
		return (float64(s.samples[half]) + float64(s.samples[half-1])) / 2
	}
	return float64(s.samples[half])
}

// Mean returns the mean of the samples.
func (s Stats) Mean() float64 {
	return float64(s.sum) / float64(s.count)
}

// Stddev returns the standard deviation of the samples.
func (s Stats) Stddev() float64 {
	m := s.Mean()
	return math.Sqrt(float64(s.sum2)/float64(s.count) - m*m)
}

// Spread returns the difference of the maximal and minimal sample values.
func (s Stats) Spread() Sample {
	if s.min > s.max {
		return 0
	}
	return s.max - s.min
}

// CreateBins divides the sample space into nbins bins for tracking counts.
//
// As samples are added, the count for the corresponding bin will be
// incremented and the sample value will not be stored.
//
// This saves memory at the expense of granularity. Percentile() and Median()
// cannot be called after CreateBins() because they are no longer meaningful.
// Instead, use Bin(i) to inspect the distribution of data by bin. Any existing
// stored samples are discarded.
//
// The bins created will be:
//   (-Inf,low], (low, s/nmid+low], (s/nmid+low, 2s/nmid], ..., (high,+Inf)
//   where:
//     s = high - low
//     nmid = nbins-2
//
// Thus, the space (high-low) is divided into nbins-2 equally sized pieces
// and the remaining two bins extend from -math.MaxFloat64 to low and high to
// math.MaxFloat64.
//
// Low must be strictly less than high, so nbins must be at least 3.
func (s *Stats) CreateBins(nbins int, low, high Sample) {
	if high <= low {
		panic("high must be greater than low")
	}
	if nbins < 3 {
		panic("Not enough bins")
	}
	spread := high - low
	s.bins = make([]Sample, nbins)
	s.binCounts = make([]int, nbins)
	for i := 0; i < nbins-1; i++ {
		s.bins[i] = Sample(i)*spread/Sample(nbins-2) + low
	}
	s.bins[nbins-1] = math.MaxFloat64
	// save memory: stop storing samples now that we track by bins
	s.samples = []Sample{}
}

// CreateBinsDiscard is shorthand for calling CreateBins(nbins, ...) with low
// value s.Percentile(discardPct) and high value s.Percentile(1-discardPct)
// with a check to make sure enough samples have been collected to make
// discardPct meaningful (1/discardPct samples are required).
func (s *Stats) CreateBinsDiscard(nbins int, discardPct float64) {
	if len(s.samples) < int(1.0/discardPct) {
		panic("Not enough samples")
	}
	s.CreateBins(nbins, s.Percentile(discardPct), s.Percentile(1.0-discardPct))
}

// Returns the count and low and high ends of the i'th bin.
//
// The bin interval is (low,high]
func (s Stats) Bin(i int) (count int, low, high Sample) {
	count = s.binCounts[i]
	high = s.bins[i]
	if i == 0 {
		low = -math.MaxFloat64
	} else {
		low = s.bins[i-1]
	}
	return
}
