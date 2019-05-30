// Copyright 2019 The airtable-go AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package airtable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// TODO: Make version dynamic
	defaultAPIURL        = "https://api.airtable.com/v0/"
	defaultRetryWaitMin  = 30 * time.Second
	defaultRetryWaitMax  = 60 * time.Second
	defaultRetryAttempts = 4
)

// Client is a client for interacting with the Airtable API
type Client struct {
	client        *http.Client
	APIURL        *url.URL
	APIKey        string
	Base          string
	RetryWaitMin  time.Duration
	RetryWaitMax  time.Duration
	RetryAttempts int

	common service

	Tables *TableService
}

type service struct {
	Selected string
	client   *Client
}

// NewClient returns a new Client
func NewClient(httpClient *http.Client, base string, apiKey string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	apiURL, _ := url.Parse(defaultAPIURL)

	c := &Client{client: httpClient, APIURL: apiURL, APIKey: apiKey, Base: base, RetryWaitMin: defaultRetryWaitMin, RetryWaitMax: defaultRetryWaitMax, RetryAttempts: defaultRetryAttempts}
	c.common.client = c

	c.Tables = (*TableService)(&c.common)

	return c
}

// Table allows selection of a specific table in the base
func (c *Client) Table(name string) *TableService {
	c.Tables.Selected = name
	return c.Tables
}

// NewRequest makes authenticated requests to the Airtable API
func (c *Client) NewRequest(method, endpoint string, params interface{}, body interface{}) (*http.Request, error) {

	u, err := c.APIURL.Parse(c.Base + "/" + endpoint)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	var bearer = "Bearer " + c.APIKey

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)

	return req, nil
}

// Response is a Airtable API response
type Response struct {
	*http.Response

	Offset string
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	// TODO: handle paging
	return response
}

// Do sends an API request and returns the API response.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

// CheckResponse checks if the response from the API
func CheckResponse(r *http.Response) error {
	if r.StatusCode == http.StatusOK {
		return nil
	}

	switch {
	case r.StatusCode == http.StatusTooManyRequests:
		return &RateLimitError{
			Response: r,
			Message:  fmt.Sprintf("API rate limit exceeded, must wait at least %d seconds before retry.", defaultRetryWaitMin),
		}
	default:
		return &ResponseError{
			Response: r,
			Message:  fmt.Sprint("Error processing request."),
		}
	}

}

// ResponseError is the error returned when request fails
type ResponseError struct {
	Response *http.Response
	Message  string `json:"message"`
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

// RateLimitError is the error returned when request fails due to rate limiting - status code 429
type RateLimitError struct {
	Response *http.Response
	Message  string `json:"message"`
}

func (r *RateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

// TODO: backoff and ratelimiting
