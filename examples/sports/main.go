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

	ss := client.SportsService
	sports, res, err := ss.GetSports()
	if err != nil {
		log.Fatalf("error requesting sports: %s", err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("invalid status code: %d %s", res.StatusCode, res.Status)
	}

	log.Printf("returned %d sports", len(sports))

	sportsBytes, err := json.Marshal(&sports)
	if err != nil {
		log.Fatalf("error marshalling sports: %s", err)
	}

	outFilePath := "sports.json"

	f, err := os.Create(outFilePath)
	if err != nil {
		log.Fatalf("error creating file %s: %s", outFilePath, err)
	}

	_, err = f.Write(sportsBytes)
	if err != nil {
		log.Fatalf("error writing to file %s: %s", outFilePath, err)
	}
}
