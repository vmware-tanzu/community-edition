// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package standalone

import (
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// GetDebugLevel default is 2 (aka DefaultLogLevel)
func GetDebugLevel(s string) string {
	_, err := strconv.Atoi(s)
	if err != nil {
		return DefaultLogLevel
	}
	return s
}
