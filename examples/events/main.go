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

	service := client.EventService

	p := service.NewEventParams("americanfootball_nfl")

	data, resp, err := service.GetEvents(p)
	if err != nil {
		log.Printf("error: received status code: %d\nMessage: %s", resp.StatusCode, resp.Status)
		log.Fatalf("error requesting events: %s", err)
	}

	log.Printf("received %d events", len(data))

	b, err := json.Marshal(&data)
	if err != nil {
		log.Fatalf("error marshalling data: %s", err)
	}

	fp := "events.json"
	f, err := os.Create(fp)
	if err != nil {
		log.Fatalf("error creating file %s: %s", fp, err)
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = f.Write(b)
	if err != nil {
		log.Fatalf("error writing file %s: %s", fp, err)
	}
}
