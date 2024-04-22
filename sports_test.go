// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
	"errors"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"net/url"
	"testing"
)

func TestSportsParams_BuildPath(t *testing.T) {
	var (
		apiToken = "api-token"
		expected = "https://api.the-odds-api.com/v4/sports?apiKey=api-token"
	)

	p := &SportsParams{ApiToken: apiToken}

	bURL, _ := url.Parse(DefaultBaseUrl)
	result, err := p.BuildPath(bURL)
	if err != nil {
		t.Error(err)
	}

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestSportsService_GetSports_NewGetRequestError(t *testing.T) {
	var (
		apiToken = "api-token"
	)

	c := newTestRequestClient(apiToken)

	getReqFn := func(requestUrl string, headers *map[string]string) (*retryablehttp.Request, error) {
		return nil, errors.New("get request creation error")
	}

	DoFn := func(req *retryablehttp.Request, data interface{}) (*Response, error) {
		return nil, nil
	}

	c.SetNewGetRequest(getReqFn)
	c.SetDo(DoFn)

	ss := &SportsService{c: c}

	sports, resp, err := ss.GetSports()
	if err == nil {
		t.Error("expected error not to be nil")
	}
	if resp != nil {
		t.Error("expected response to be nil")
	}
	if sports != nil {
		t.Error("expected sports to be nil")
	}
}

func TestSportsService_GetSports_DoError(t *testing.T) {
	var (
		apiToken = "api-token"
	)

	c := newTestRequestClient(apiToken)

	getReqFn := func(requestUrl string, headers *map[string]string) (*retryablehttp.Request, error) {
		return nil, nil
	}

	DoFn := func(req *retryablehttp.Request, data interface{}) (*Response, error) {
		r := &Response{
			&http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     "bad request",
			},
		}
		return r, errors.New("error executing GET request")
	}

	c.SetNewGetRequest(getReqFn)
	c.SetDo(DoFn)

	ss := &SportsService{c: c}

	sports, resp, err := ss.GetSports()
	if err == nil {
		t.Error("expected error not to be nil")
	}
	if resp == nil {
		t.Error("expected response not to be nil")
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected status code to be 400, got %d", resp.StatusCode)
	}
	if sports != nil {
		t.Error("expected sports to be nil")
	}
}
