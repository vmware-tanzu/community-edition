//nolint
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
// A slightly modified version of the golang/x/mod/semver package tests
// that removes the leading `v`

package semver

import (
	"math/rand"
	"sort"
	"strings"
	"testing"
)

var tests = []struct {
	in  string
	out string
}{
	{"bad", ""},
	{"1-alpha.beta.gamma", ""},
	{"1-pre", ""},
	{"1+meta", ""},
	{"1-pre+meta", ""},
	{"1.2-pre", ""},
	{"1.2+meta", ""},
	{"1.2-pre+meta", ""},
	{"1.0.0-alpha", "1.0.0-alpha"},
	{"1.0.0-alpha.1", "1.0.0-alpha.1"},
	{"1.0.0-alpha.beta", "1.0.0-alpha.beta"},
	{"1.0.0-beta", "1.0.0-beta"},
	{"1.0.0-beta.2", "1.0.0-beta.2"},
	{"1.0.0-beta.11", "1.0.0-beta.11"},
	{"1.0.0-rc.1", "1.0.0-rc.1"},
	{"1", "1.0.0"},
	{"1.0", "1.0.0"},
	{"1.0.0", "1.0.0"},
	{"1.2", "1.2.0"},
	{"1.2.0", "1.2.0"},
	{"1.2.3-456", "1.2.3-456"},
	{"1.2.3-456.789", "1.2.3-456.789"},
	{"1.2.3-456-789", "1.2.3-456-789"},
	{"1.2.3-456a", "1.2.3-456a"},
	{"1.2.3-pre", "1.2.3-pre"},
	{"1.2.3-pre+meta", "1.2.3-pre"},
	{"1.2.3-pre.1", "1.2.3-pre.1"},
	{"1.2.3-zzz", "1.2.3-zzz"},
	{"1.2.3", "1.2.3"},
	{"1.2.3+meta", "1.2.3"},
	{"1.2.3+meta-pre", "1.2.3"},
	{"1.2.3+meta-pre.sha.256a", "1.2.3"},
}

func TestCompare(t *testing.T) {
	for i, ti := range tests {
		for j, tj := range tests {
			cmp := Compare(ti.in, tj.in)
			var want int
			if ti.out == tj.out {
				want = 0
			} else if i < j {
				want = -1
			} else {
				want = +1
			}
			if cmp != want {
				t.Errorf("Compare(%q, %q) = %d, want %d", ti.in, tj.in, cmp, want)
			}
		}
	}
}

func TestSort(t *testing.T) {
	versions := make([]string, len(tests))
	for i, test := range tests {
		versions[i] = test.in
	}
	rand.Shuffle(len(versions), func(i, j int) { versions[i], versions[j] = versions[j], versions[i] })
	Sort(versions)
	if !sort.IsSorted(ByVersion(versions)) {
		t.Errorf("list is not sorted:\n%s", strings.Join(versions, "\n"))
	}
}
