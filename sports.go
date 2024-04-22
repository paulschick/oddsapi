// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
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

func (s *SportsParams) BuildPath(baseUrl *url.URL) (string, error) {
	return buildPath(s, "v4/sports", baseUrl)
}

type SportsService struct {
	c BaseRequestClient
}

func NewSportsService(c *Client) *SportsService {
	return &SportsService{c: c}
}

func (s *SportsService) GetSports() ([]*Sports, *Response, error) {
	params := &SportsParams{ApiToken: s.c.GetApiToken()}
	var data []*Sports
	return requestHandler(params, s.c, data)
}
