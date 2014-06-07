// Copyright 2014 The Google Genomics API Client in Go Authors.
// All rights reserved.
// Use of this source code is governed by a Apache license that can be
// found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
)

func indent(level int) {
	fmt.Print(strings.Repeat("\t", level))
}
