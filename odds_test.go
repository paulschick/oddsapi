// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
	"testing"
	"time"
)

func TestOddsParams_Encode(t *testing.T) {
	p := &OddsParams{
		SportKey:   "upcoming",
		ApiToken:   "api-key",
		Region:     "us",
		Markets:    "spreads",
		DateFormat: "iso",
		OddsFormat: "decimal",
	}

	q, err := query.Values(p)
	if err != nil {
		t.Error(err)
	}

	encoded := q.Encode()
	basePath := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/%s/odds/", p.SportKey)
	baseURL, _ := url.Parse(basePath)

	baseURL.RawQuery = encoded

	result := baseURL.String()
	expected := "https://api.the-odds-api.com/v4/sports/upcoming/odds/?apiKey=api-key&dateFormat=iso&markets=spreads&oddsFormat=decimal&regions=us"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestOddsParams_SetRegions(t *testing.T) {
	p := NewOddsParams("api-key", "upcoming")
	regions := []Region{RegionUs, RegionUs2}
	err := p.SetRegions(regions...)
	if err != nil {
		t.Error(err)
	}
	expectedRegions := "us,us2"
	resultRegions := p.Region
	if resultRegions != expectedRegions {
		t.Errorf("expected regions '%s', got '%s'", expectedRegions, resultRegions)
	}

	regions2 := []Region{RegionUs}
	p = NewOddsParams("api-key", "upcoming")
	err = p.SetRegions(regions2...)
	if err != nil {
		t.Error(err)
	}

	if p.Region != "us" {
		t.Errorf("expected region 'us', got '%s'", p.Region)
	}
}

func TestOddsParams_SetMarkets(t *testing.T) {
	p := NewOddsParams("api-key", "upcoming")
	err := p.SetMarkets("h2h", "spreads")
	if err != nil {
		t.Error(err)
	}
	expected := "h2h,spreads"
	if p.Markets != expected {
		t.Errorf("expected '%s', got '%s'", expected, p.Markets)
	}

	err = p.SetMarkets("spreads")
	if err != nil {
		t.Error(err)
	}
	exp2 := "spreads"
	if p.Markets != exp2 {
		t.Errorf("expected '%s', got '%s'", exp2, p.Markets)
	}

	err = p.SetMarkets("us", "us2")
	if err == nil {
		t.Errorf("expected an error when setting invalid markets, did not get an error")
	}
}

func TestOddsParams_SetEventIds(t *testing.T) {
	p := NewOddsParams("api-key", "upcoming")
	p.SetEventIds("event1", "event2")
	expected := "event1,event2"
	if *p.EventIds != expected {
		t.Errorf("expected event IDs '%s', got '%s'", expected, *p.EventIds)
	}

	p.SetEventIds()
	if p.EventIds != nil {
		t.Errorf("expected event IDs to be nil, got '%s'", *p.EventIds)
	}
}

func TestOddsParams_SetBookmakers(t *testing.T) {
	p := NewOddsParams("api-key", "upcoming")
	p.SetBookmakers("bm1", "bm2")
	expected := "bm1,bm2"
	if *p.Bookmakers != expected {
		t.Errorf("expected Bookmakers '%s', got '%s'", expected, *p.Bookmakers)
	}

	p.SetBookmakers()
	if p.Bookmakers != nil {
		t.Errorf("expected Bookmakers to be nil, got '%s'", *p.Bookmakers)
	}
}

func TestOddsParams_SetCommenceTimeFromISO(t *testing.T) {
	testDateGood := "2024-01-01T00:00:00Z"
	testDateBad := "2024-01-01"

	p := NewOddsParams("api-key", "upcoming")
	err := p.SetCommenceTimeFromISO(testDateGood)
	if err != nil {
		t.Errorf("expected commence time from to pass, failed with error: %s", err.Error())
	}

	if *p.CommenceTimeFrom != testDateGood {
		t.Errorf("expected CommenceTimeFrom to be '%s', got '%s'", testDateGood, *p.CommenceTimeFrom)
	}

	err = p.SetCommenceTimeFromISO(testDateBad)
	if err == nil {
		t.Errorf("expected date '%s' to fail, it did not.", testDateBad)
	}
}

func TestOddsParams_SetCommenceTimeFrom(t *testing.T) {
	p := NewOddsParams("api-key", "upcoming")
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := "2024-01-01T00:00:00Z"
	p.SetCommenceTimeFrom(date)
	if *p.CommenceTimeFrom != expected {
		t.Errorf("expected CommenceTimeFrom to be '%s', got '%s'", expected, *p.CommenceTimeFrom)
	}
}

func TestOddsParams_SetCommenceTimeToISO(t *testing.T) {
	testDateGood := "2024-01-01T00:00:00Z"
	testDateBad := "2024-01-01"

	p := NewOddsParams("api-key", "upcoming")
	err := p.SetCommenceTimeToISO(testDateGood)
	if err != nil {
		t.Errorf("expected commence time to to pass, failed with error: %s", err.Error())
	}

	if *p.CommenceTimeTo != testDateGood {
		t.Errorf("expected CommenceTimeTo to be '%s', got '%s'", testDateGood, *p.CommenceTimeTo)
	}

	err = p.SetCommenceTimeToISO(testDateBad)
	if err == nil {
		t.Errorf("expected date '%s' to fail, it did not.", testDateBad)
	}
}

func TestOddsParams_SetCommenceTimeTo(t *testing.T) {
	p := NewOddsParams("api-key", "upcoming")
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := "2024-01-01T00:00:00Z"
	p.SetCommenceTimeTo(date)
	if *p.CommenceTimeTo != expected {
		t.Errorf("expected CommenceTimeTo to be '%s', got '%s'", expected, *p.CommenceTimeTo)
	}
}
