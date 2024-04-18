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

	service := client.OddsService
	params := service.NewOddsParamsUpcoming()
	_ = params.SetRegions(oddsapi.RegionUs)
	_ = params.SetMarkets(oddsapi.MarketH2H)

	data, response, err := service.GetOdds(params)
	if err != nil {
		log.Fatalf("error requesting odds: %s | status code: %d, status: %s", err, response.StatusCode, response.Status)
	}

	if response.StatusCode != 200 {
		log.Fatalf("bad response for odds request %d: %s", response.StatusCode, response.Status)
	}

	log.Printf("Returned %d odds", len(data))

	fp := "odds.json"
	f, err := os.Create(fp)
	if err != nil {
		log.Fatalf("error creating file %s: %s", fp, err)
	}

	dataBytes, err := json.Marshal(&data)
	if err != nil {
		log.Fatalf("error marshalling data: %s", err)
	}

	_, err = f.Write(dataBytes)
	if err != nil {
		log.Fatalf("error writing file %s: %s", fp, err)
	}
}
