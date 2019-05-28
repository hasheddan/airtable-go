package airtable

import (
	"context"
	"time"
)

// TableService allows connection to table methods of the API
type TableService service

// Table represents and Airtable Table
type Table struct {
	Name    string
	Records []Record `json:"records,omitempty"`
	Offset  string   `json:"offset,omitempty"`
}

// Record represents a record in an Airtable table
type Record struct {
	ID          string                 `json:"id,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
	CreatedTime time.Time              `json:"createdTime,omitempty"`
}

// Get gets records from a table
func (s *TableService) Get(ctx context.Context) (*Table, error) {
	req, err := s.client.NewRequest("GET", s.Selected, nil)
	if err != nil {
		return nil, err
	}

	var t Table
	_, err = s.client.Do(ctx, req, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
