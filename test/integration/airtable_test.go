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
		print("Must provide key and base to run integration tests. \n\n")
		os.Exit(1)
	}
	client = airtable.NewClient(nil, base, key)
}
