// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding/json"
	"log"
	"oddsapi"
	"os"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalf("no api key provided")
	}
	client, err := oddsapi.NewClient(apiKey, 10)
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}

	service := client.EventOddsService
	p := service.NewParams("americanfootball_nfl", "eca3b71919531e7ae0b4f3f501157e6c")
	p.SetDateFormat(oddsapi.DefaultDateFormat)
	p.SetOddsFormat(oddsapi.DecimalOddsFormat)

	data, resp, err := service.GetOdds(p)
	if err != nil {
		log.Printf("bad response getting event odds: %d\nmessage: %s", resp.StatusCode, resp.Status)
		log.Fatalf("error getting event odds: %s", err)
	}

	fp := "event-odds.json"
	f, err := os.Create(fp)
	if err != nil {
		log.Fatalf("error creating file %s: %s", fp, err)
	}

	b, err := json.Marshal(&data)
	if err != nil {
		log.Fatalf("err marshalling data: %s", err)
	}

	_, err = f.Write(b)
	if err != nil {
		log.Fatalf("error writing data: %s", err)
	}
}
