package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/charleshuang3/autoget/organizer/agents/javactor"
)

var (
	apiKey = os.Getenv("GEMINI_API_KEY")
)

func main() {
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set.")
	}

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <directory> [--dry-run]", os.Args[0])
	}

	dir := os.Args[1]
	dryRun := false
	if len(os.Args) > 2 && os.Args[2] == "--dry-run" {
		dryRun = true
		log.Println("Dry run mode enabled. No changes will be made.")
	}

	actorFilePath := filepath.Join(dir, "actor.json")

	// Check if dir/actor.json exists
	_, err := os.Stat(actorFilePath)
	if os.IsNotExist(err) {
		log.Fatalf("Error: %s does not exist.", actorFilePath)
	} else if err != nil {
		log.Fatalf("Error checking %s: %v", actorFilePath, err)
	}

	// Read dir/actor.json
	actorData := make(map[string][]string)
	fileContent, err := os.ReadFile(actorFilePath)
	if err != nil {
		log.Fatalf("Error reading %s: %v", actorFilePath, err)
	}

	if len(fileContent) > 0 {
		err = json.Unmarshal(fileContent, &actorData)
		if err != nil {
			log.Fatalf("Error unmarshaling %s: %v", actorFilePath, err)
		}
	}

	// List all folders in dir
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading directory %s: %v", dir, err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex for protecting actorData during concurrent writes

	var processedCount int32
	var skippedCount int32
	var errorCount int32
	var totalFolders int32

	foldersToProcess := []string{}
	for _, file := range files {
		if file.IsDir() {
			totalFolders++
			folderName := file.Name()
			if _, exists := actorData[folderName]; exists {
				atomic.AddInt32(&skippedCount, 1)
				fmt.Printf("actor %q: skipped\n", folderName)
			} else {
				foldersToProcess = append(foldersToProcess, folderName)
			}
		}
	}

	if dryRun {
		for _, folderName := range foldersToProcess {
			fmt.Printf("Dry run: %q\n", folderName)
		}
		fmt.Printf("\nDry run complete. Total folders: %d, Would process: %d, Skipped: %d\n",
			totalFolders, len(foldersToProcess), skippedCount)
		return
	}

	for i, folderName := range foldersToProcess {
		wg.Add(1)
		go func(idx int, name string) {
			defer wg.Done()
			atomic.AddInt32(&processedCount, 1)
			result, runErr := javactor.Run(apiKey, name)
			if runErr != nil {
				atomic.AddInt32(&errorCount, 1)
				fmt.Printf("actor %q error [%d/%d] %v\n", name, processedCount, totalFolders, runErr)
				return
			}

			mu.Lock()
			actorData[name] = result
			mu.Unlock()
			fmt.Printf("actor %q done [%d/%d]\n", name, processedCount, totalFolders)
		}(i, folderName)
	}

	wg.Wait() // Wait for all goroutines to finish

	// check if duplicated name in json
	exist := make(map[string]string)
	for actor, names := range actorData {
		for _, name := range names {
			if another, ok := exist[name]; ok {
				fmt.Printf("Confilct found: %q -> %q and %q -> %q\n", actor, name, another, name)
			}
			exist[name] = actor
		}
	}

	// After all finish write json
	updatedContent, err := json.MarshalIndent(actorData, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling updated actor data: %v", err)
	}

	err = os.WriteFile(actorFilePath, updatedContent, 0644)
	if err != nil {
		log.Fatalf("Error writing updated %s: %v", actorFilePath, err)
	}

	fmt.Printf("Successfully updated %s\n", actorFilePath)
	fmt.Printf("\n--- Processing Summary ---\n")
	fmt.Printf("Total folders scanned: %d\n", totalFolders)
	fmt.Printf("Folders processed: %d\n", processedCount)
	fmt.Printf("Folders skipped (already in actor.json): %d\n", skippedCount)
	fmt.Printf("Folders with errors during processing: %d\n", errorCount)
	fmt.Printf("--------------------------\n")
}
