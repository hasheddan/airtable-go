// Copyright 2019 The airtable-go AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build integration

package integration

import (
	"os"

	"github.com/hasheddan/airtable-go/airtable"
)

var client *airtable.Client

func init() {
	key := os.Getenv("AIRTABLE_KEY")
	base := os.Getenv("AIRTABLE_BASE")
	if key == "" || base == "" {
		print("Integration tests will be unable to run successfully without . \n\n")
	}
	client = airtable.NewClient(nil, base, key)
}
