// Copyright 2012 The Summstat Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package summstat allows one to incrementally compute summary statistics
for a data set.

It allows accurate median and percentiles to be returned, though these require
the entire dataset to be retained. If the dataset cannot be kept in memory, a
subset of it can be stored initially to determine an approximate range of
interest to divide into bins for which counts will be tracked. This allows one
to collect approximate percentile data without the memory overhead.
*/
package summstat
