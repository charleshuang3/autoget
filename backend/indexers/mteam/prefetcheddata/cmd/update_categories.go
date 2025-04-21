package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/charleshuang3/autoget/backend/indexers/mteam/prefetcheddata"
)

// This is a tools used to update categories.
func main() {
	var outputFile string
	var apiKey string

	flag.StringVar(&outputFile, "o", "", "Output file path for categories JSON")
	flag.StringVar(&apiKey, "k", "", "M-Team API Key")
	flag.Parse()

	if outputFile == "" {
		log.Fatal().Msg("Output file path (-o) is required")
	}
	if apiKey == "" {
		log.Fatal().Msg("API Key (-k) is required")
	}

	log.Info().Msg("Fetching categories from M-Team API...")
	categories, err := prefetcheddata.FetchCategories(apiKey, true)
	if err != nil {
		log.Fatal().Msgf("Failed to fetch categories: %v", err)
	}
	log.Info().Msg("Successfully fetched categories.")

	log.Info().Msgf("Encoding categories to JSON and writing to %s...", outputFile)
	jsonData, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		log.Fatal().Msgf("Failed to marshal categories to JSON: %v", err)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatal().Msgf("Failed to write categories to file %s: %v", outputFile, err)
	}

	log.Info().Msgf("Successfully wrote categories to %s", outputFile)
}
