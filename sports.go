// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
	"github.com/google/go-querystring/query"
	"net/url"
)

type Sports struct {
	Key          string `json:"key" csv:"key"`
	Active       bool   `json:"active" csv:"active"`
	Group        string `json:"group" csv:"group"`
	Description  string `json:"description" csv:"description"`
	Title        string `json:"title" csv:"title"`
	HasOutrights bool   `json:"has_outrights" csv:"has_outrights"`
}

type SportsParams struct {
	ApiToken string `url:"apiKey"`
}

func (s *SportsParams) GetEncoded() (string, error) {
	q, err := query.Values(s)
	if err != nil {
		return "", err
	}
	return q.Encode(), nil
}

func (s *SportsParams) BuildPath(baseUrl *url.URL) (string, error) {
	basePath := "v4/sports"
	bURLCopy := *baseUrl
	bURL := &bURLCopy
	bURL = bURL.JoinPath(basePath)
	encoded, err := s.GetEncoded()
	if err != nil {
		return "", err
	}
	bURL.RawQuery = encoded
	return bURL.String(), nil
}

type SportsService struct {
	c *Client
}

func NewSportsService(c *Client) *SportsService {
	return &SportsService{c: c}
}

func (s *SportsService) GetSports() ([]*Sports, *Response, error) {
	params := &SportsParams{ApiToken: s.c.apiToken}
	reqUrl, err := params.BuildPath(s.c.GetBaseUrl())
	if err != nil {
		return nil, nil, err
	}

	req, err := s.c.NewGetRequest(reqUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	var data []*Sports
	resp, err := s.c.Do(req, &data)
	if err != nil {
		return nil, resp, err
	}

	return data, resp, nil
}
