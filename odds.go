// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
	"strings"
	"time"
)

type Outcome struct {
	Name  string   `json:"name" csv:"name"`
	Price float64  `json:"price" csv:"price"`
	Point *float64 `json:"point,omitempty" csv:"point,omitempty"`
}

type Market struct {
	Key        string     `json:"key" csv:"key"`
	LastUpdate string     `json:"last_update" csv:"last_update"`
	Outcomes   []*Outcome `json:"outcomes" csv:"outcomes"`
}

type BookMaker struct {
	Key        string    `json:"key" csv:"key"`
	Title      string    `json:"title" csv:"title"`
	LastUpdate string    `json:"last_update" csv:"last_update"`
	Markets    []*Market `json:"markets" csv:"markets"`
}

type Odds struct {
	Id           string       `json:"id" csv:"id"`
	SportKey     string       `json:"sport_key" csv:"sport_key"`
	SportTitle   string       `json:"sport_title" csv:"sport_title"`
	CommenceTime string       `json:"commence_time" csv:"commence_time"`
	HomeTeam     string       `json:"home_team" csv:"home_team"`
	AwayTeam     string       `json:"away_team" csv:"away_team"`
	BookMakers   []*BookMaker `json:"bookmakers" csv:"bookmakers"`
}

type OddsParams struct {
	SportKey string `url:"-"`
	ApiToken string `url:"apiKey"`

	// This can be one or more Regions separated by commas
	Region string `url:"regions"`

	// Markets is optional. This should be one MarketKey or multiple
	// in a single string, comma-separated
	Markets    string     `url:"markets,omitempty"`
	DateFormat DateFormat `url:"dateFormat,omitempty"`
	OddsFormat OddsFormat `url:"oddsFormat,omitempty"`

	// Optional event ids passed as a comma-separated string
	EventIds *string `url:"eventIds,omitempty"`

	// Optional bookmakers to return as a comma-separated string
	Bookmakers *string `url:"bookmakers,omitempty"`

	// Optional for games that commence on or after
	// No effect if sport is upcoming
	// ISO 8601
	CommenceTimeFrom *string `url:"commenceTimeFrom,omitempty"`

	// Commence on or before this option. Optional
	// No effect if sport is upcoming
	// ISO 8601
	CommenceTimeTo *string `url:"commentTimeTo,omitempty"`
}

func NewOddsParams(apiKey, sportKey string) *OddsParams {
	return &OddsParams{
		ApiToken: apiKey,
		SportKey: sportKey,
	}
}

func NewOddsParamsUpcoming(apiKey string) *OddsParams {
	return &OddsParams{
		ApiToken: apiKey,
		SportKey: DefaultSports,
	}
}

func (o *OddsParams) SetRegions(regions ...Region) error {
	if regions == nil {
		return nil
	}
	r := make([]string, len(regions))
	for i, region := range regions {
		if !region.Valid() {
			o.Region = DefaultRegion.String()
			return fmt.Errorf("invalid region provided: %s", region)
		}
		r[i] = region.String()
	}
	regionStr := strings.Join(r, ",")
	o.Region = regionStr
	return nil
}

func (o *OddsParams) SetMarkets(markets ...MarketKey) error {
	if markets == nil {
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
	o.Markets = marketStr
	return nil
}

func (o *OddsParams) SetEventIds(eventIds ...string) {
	if eventIds == nil {
		o.EventIds = nil
		return
	}
	eventStr := strings.Join(eventIds, ",")
	o.EventIds = &eventStr
}

func (o *OddsParams) SetBookmakers(bookmakers ...string) {
	if bookmakers == nil {
		o.Bookmakers = nil
		return
	}
	bStr := strings.Join(bookmakers, ",")
	o.Bookmakers = &bStr
}

func (o *OddsParams) SetCommenceTimeFromISO(timeFrom string) error {
	if _, err := time.Parse(time.RFC3339, timeFrom); err != nil {
		return err
	}
	o.CommenceTimeFrom = &timeFrom
	return nil
}

func (o *OddsParams) SetCommenceTimeFrom(timeFrom time.Time) {
	iso := timeFrom.Format(time.RFC3339)
	o.CommenceTimeFrom = &iso
}

func (o *OddsParams) SetCommenceTimeToISO(timeTo string) error {
	if _, err := time.Parse(time.RFC3339, timeTo); err != nil {
		return err
	}
	o.CommenceTimeTo = &timeTo
	return nil
}

func (o *OddsParams) SetCommenceTimeTo(timeTo time.Time) {
	iso := timeTo.Format(time.RFC3339)
	o.CommenceTimeTo = &iso
}

func (o *OddsParams) GetEncoded() (string, error) {
	if o.Markets == "" {
		o.Markets = MarketH2H.String()
	}
	if o.DateFormat == "" || !o.DateFormat.Valid() {
		o.DateFormat = DefaultDateFormat
	}
	if o.OddsFormat == "" || !o.OddsFormat.Valid() {
		o.OddsFormat = DefaultOddsFormat
	}
	if o.EventIds != nil && *o.EventIds == "" {
		o.EventIds = nil
	}
	if o.Bookmakers != nil && *o.Bookmakers == "" {
		o.Bookmakers = nil
	}
	if o.CommenceTimeFrom != nil && *o.CommenceTimeFrom == "" {
		o.CommenceTimeFrom = nil
	}
	if o.CommenceTimeTo != nil && *o.CommenceTimeTo == "" {
		o.CommenceTimeTo = nil
	}
	q, err := query.Values(o)
	if err != nil {
		return "", err
	}
	return q.Encode(), nil
}

func (o *OddsParams) BuildPath(baseUrl *url.URL) (string, error) {
	basePath := fmt.Sprintf("v4/sports/%s/odds/", o.SportKey)
	bURLCopy := *baseUrl
	bURL := &bURLCopy
	bURL = bURL.JoinPath(basePath)
	encoded, err := o.GetEncoded()
	if err != nil {
		return "", err
	}
	bURL.RawQuery = encoded
	return bURL.String(), nil
}

type OddsService struct {
	c *Client
}

func NewOddsService(c *Client) *OddsService {
	return &OddsService{c: c}
}

func (o *OddsService) NewOddsParams(sportKey string) *OddsParams {
	return NewOddsParams(o.c.apiToken, sportKey)
}

func (o *OddsService) NewOddsParamsUpcoming() *OddsParams {
	return NewOddsParamsUpcoming(o.c.apiToken)
}

func (o *OddsService) GetOdds(params *OddsParams) ([]*Odds, *Response, error) {
	reqUrl, err := params.BuildPath(o.c.GetBaseUrl())
	if err != nil {
		return nil, nil, err
	}

	req, err := o.c.NewGetRequest(reqUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	var data []*Odds
	resp, err := o.c.Do(req, &data)
	if err != nil {
		return nil, resp, err
	}

	return data, resp, nil
}
