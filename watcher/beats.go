package watcher

import (
	"fmt"
	"log"

	"github.com/KnightChaser/sentinela"
)

func EVTXBeats(filePath string) {
	stats, err := sentinela.ParseEVTX(filePath)
	if err != nil {
		log.Fatalf("Error parsing evtx: %s", err)
	}

	fmt.Printf("# of events: %v", len(stats.Event))
}
