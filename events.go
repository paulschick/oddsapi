// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Event struct {
	Id           string `json:"id"`
	SportKey     string `json:"sport_key"`
	SportTitle   string `json:"sport_title"`
	CommenceTime string `json:"commence_time"`
	HomeTeam     string `json:"home_team"`
	AwayTeam     string `json:"away_team"`
}

type EventParams struct {
	SportKey   string     `url:"-"`
	ApiToken   string     `url:"apiKey"`
	DateFormat DateFormat `url:"dateFormat,omitempty"`

	// Optional event ids passed as a comma-separated string
	EventIds *string `url:"eventIds,omitempty"`

	// Optional for games that commence on or after
	// No effect if sport is upcoming
	// ISO 8601
	CommenceTimeFrom *string `url:"commenceTimeFrom,omitempty"`

	// Commence on or before this option. Optional
	// No effect if sport is upcoming
	// ISO 8601
	CommenceTimeTo *string `url:"commentTimeTo,omitempty"`
}

func (e *EventParams) SetEventIds(eventIds ...string) {
	if eventIds == nil {
		e.EventIds = nil
		return
	}
	eventStr := strings.Join(eventIds, ",")
	e.EventIds = &eventStr
}

func (e *EventParams) SetCommenceTimeFromISO(timeFrom string) error {
	if _, err := time.Parse(time.RFC3339, timeFrom); err != nil {
		return err
	}
	e.CommenceTimeFrom = &timeFrom
	return nil
}

func (e *EventParams) SetCommenceTimeFrom(timeFrom time.Time) {
	iso := timeFrom.Format(time.RFC3339)
	e.CommenceTimeFrom = &iso
}

func (e *EventParams) SetCommenceTimeToISO(timeTo string) error {
	if _, err := time.Parse(time.RFC3339, timeTo); err != nil {
		return err
	}
	e.CommenceTimeTo = &timeTo
	return nil
}

func (e *EventParams) SetCommenceTimeTo(timeTo time.Time) {
	iso := timeTo.Format(time.RFC3339)
	e.CommenceTimeTo = &iso
}

func (e *EventParams) checkSetDateFormat() {
	if e.DateFormat == "" || !e.DateFormat.Valid() {
		e.DateFormat = DefaultDateFormat
	}
}

func (e *EventParams) checkSetEventIds() {
	if e.EventIds != nil && *e.EventIds == "" {
		e.EventIds = nil
	}
}

func (e *EventParams) checkSetCommenceTimes() {
	if e.CommenceTimeFrom != nil && *e.CommenceTimeFrom == "" {
		e.CommenceTimeFrom = nil
	}
	if e.CommenceTimeTo != nil && *e.CommenceTimeTo == "" {
		e.CommenceTimeTo = nil
	}
}

func (e *EventParams) BuildPath(baseUrl *url.URL) (string, error) {
	if e.SportKey == "" {
		return "", errors.New("no sports key provided")
	}
	basePath := fmt.Sprintf("v4/sports/%s/events", e.SportKey)
	return buildPath(e, basePath, baseUrl, e.checkSetDateFormat, e.checkSetEventIds, e.checkSetCommenceTimes)
}

type EventService struct {
	c *Client
}

func NewEventService(c *Client) *EventService {
	return &EventService{c: c}
}

func (e *EventService) NewEventParams(sportKey string) *EventParams {
	return &EventParams{
		ApiToken:   e.c.apiToken,
		SportKey:   sportKey,
		DateFormat: DefaultDateFormat,
	}
}

func (e *EventService) GetEvents(params *EventParams) ([]*Event, *Response, error) {
	var data []*Event
	return requestHandler[[]*Event](params, e.c, data)
}
