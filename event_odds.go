// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
	"errors"
	"fmt"
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

func (e *EventOddsParams) ValidateDateFormat() {
	if e.DateFormat == "" || !e.DateFormat.Valid() {
		e.DateFormat = DefaultDateFormat
	}
}

func (e *EventOddsParams) ValidateOddsFormat() {
	if e.OddsFormat == "" || !e.OddsFormat.Valid() {
		e.OddsFormat = DefaultOddsFormat
	}
}

func (e *EventOddsParams) ValidateBookmakers() {
	if e.Bookmakers != nil && *e.Bookmakers == "" {
		e.Bookmakers = nil
	}
}

func (e *EventOddsParams) ValidateRegion() {
	if e.Region == "" {
		e.Region = DefaultRegion.String()
	}
}

func (e *EventOddsParams) BuildPath(baseUrl *url.URL) (string, error) {
	if e.SportKey == "" {
		return "", errors.New("sports key is empty")
	} else if e.EventKey == "" {
		return "", errors.New("event key is empty")
	}
	basePath := fmt.Sprintf("v4/sports/%s/events/%s/odds", e.SportKey, e.EventKey)
	return buildPath(e, basePath, baseUrl,
		e.ValidateDateFormat, e.ValidateOddsFormat, e.ValidateBookmakers, e.ValidateRegion)
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
	var data *Odds
	return requestHandler(params, e.c, data)
}
