// Copyright 2019 The airtable-go AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// TableGetParams defines the available options for returning records of a table
type TableGetParams struct {
	Fields        string `param:"fields,omitempty"`
	Formula       string `param:"filterByFormula,omitempty"`
	MaxRecords    int    `param:"maxRecords,omitempty"`
	PageSize      int    `param:"pageSize,omitempty"`
	SortField     string `param:"sortField,omitempty"`
	SortDirection string `param:"sortDirection,omitempty"`
	View          string `param:"view,omitempty"`
	CellFormat    string `param:"cellFormat,omitempty"`
	TimeZone      string `param:"timeZone,omitempty"`
	UserLocale    string `param:"userLocale,omitempty"`
}

// Get gets records from a table
func (s *TableService) Get(ctx context.Context) (*Table, error) {
	req, err := s.client.NewRequest("GET", s.Selected, &TableGetParams{}, nil)
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
