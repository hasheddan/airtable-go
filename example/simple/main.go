// Copyright 2019 The airtable-go AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hasheddan/airtable-go/airtable"
)

func main() {
	client := airtable.NewClient(nil, os.Getenv("AIRTABLE_BASE"), os.Getenv("AIRTABLE_KEY"))
	t, err := client.Table("Schedule").Get(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(t.Records); i++ {
		fmt.Println(t.Records[i])
	}
}
