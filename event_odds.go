// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
	"strings"
)

type EventOddsParams struct {
	SportKey string `url:"-"`
	EventKey string `url:"-"`
	ApiToken string `url:"apiKey"`

	// This can be one or more Regions separated by commas
	Region string `url:"regions"`

	// Markets is optional. This should be one MarketKey or multiple
	// in a single string, comma-separated
	Markets string `url:"markets,omitempty"`

	DateFormat DateFormat `url:"dateFormat,omitempty"`
	OddsFormat OddsFormat `url:"oddsFormat,omitempty"`

	// Optional bookmakers to return as a comma-separated string
	Bookmakers *string `url:"bookmakers,omitempty"`
}

func (e *EventOddsParams) SetRegions(regions ...Region) error {
	if regions == nil {
		e.Region = "us"
		return nil
	}
	r := make([]string, len(regions))
	for i, region := range regions {
		if !region.Valid() {
			e.Region = DefaultRegion.String()
			return fmt.Errorf("invalid region provided: %s", region)
		}
		r[i] = region.String()
	}
	regionStr := strings.Join(r, ",")
	e.Region = regionStr
	return nil
}

func (e *EventOddsParams) SetMarkets(markets ...MarketKey) error {
	if markets == nil {
		e.Markets = "h2h"
		return nil
	}
	m := make([]string, len(markets))
	for i, market := range markets {
		if !market.Valid() {
			return fmt.Errorf("invalid market provided: %s", market)
		}
		m[i] = market.String()
	}
	marketStr := strings.Join(m, ",")
	e.Markets = marketStr
	return nil
}

func (e *EventOddsParams) SetBookmakers(bookmakers ...string) {
	if bookmakers == nil {
		e.Bookmakers = nil
		return
	}
	bStr := strings.Join(bookmakers, ",")
	e.Bookmakers = &bStr
}

func (e *EventOddsParams) SetDateFormat(dateFormat DateFormat) bool {
	if !dateFormat.Valid() {
		e.DateFormat = DefaultDateFormat
		return false
	}
	e.DateFormat = dateFormat
	return true
}

func (e *EventOddsParams) SetOddsFormat(oddsFormat OddsFormat) bool {
	if !oddsFormat.Valid() {
		e.OddsFormat = DefaultOddsFormat
		return false
	}
	e.OddsFormat = oddsFormat
	return true
}

func (e *EventOddsParams) GetEncoded() (string, error) {
	if e.DateFormat == "" || !e.DateFormat.Valid() {
		e.DateFormat = DefaultDateFormat
	}
	if e.OddsFormat == "" || !e.OddsFormat.Valid() {
		e.OddsFormat = DefaultOddsFormat
	}
	if e.Bookmakers != nil && *e.Bookmakers == "" {
		e.Bookmakers = nil
	}
	if e.Region == "" {
		e.Region = DefaultRegion.String()
	}

	q, err := query.Values(e)
	if err != nil {
		return "", err
	}
	return q.Encode(), nil
}

func (e *EventOddsParams) BuildPath(baseUrl *url.URL) (string, error) {
	basePath := fmt.Sprintf("v4/sports/%s/events/%s/odds", e.SportKey, e.EventKey)
	bURLCopy := *baseUrl
	bURL := &bURLCopy
	bURL = bURL.JoinPath(basePath)
	encoded, err := e.GetEncoded()
	if err != nil {
		return "", err
	}
	bURL.RawQuery = encoded
	return bURL.String(), nil
}

type EventOddsService struct{ c *Client }

func NewEventOddsService(c *Client) *EventOddsService {
	return &EventOddsService{c: c}
}

func (e *EventOddsService) NewParams(sportKey, eventKey string) *EventOddsParams {
	return &EventOddsParams{
		SportKey: sportKey,
		EventKey: eventKey,
		ApiToken: e.c.apiToken,
	}
}

func (e *EventOddsService) GetOdds(params *EventOddsParams) (*Odds, *Response, error) {
	reqUrl, err := params.BuildPath(e.c.GetBaseUrl())
	if err != nil {
		return nil, nil, err
	}

	req, err := e.c.NewGetRequest(reqUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	var data *Odds
	resp, err := e.c.Do(req, &data)
	if err != nil {
		return nil, resp, err
	}

	return data, resp, nil
}
