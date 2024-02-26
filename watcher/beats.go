package watcher

import (
	"fmt"
	"log"

	"github.com/KnightChaser/sentinela"
	"github.com/tidwall/gjson"
)

func EVTXBeats(filePath string) {
	stats, err := sentinela.ParseEVTX(filePath)
	if err != nil {
		log.Fatalf("Error parsing evtx: %s", err)
	}

	for _, stat := range stats.Event {
		fmt.Printf("Event: %v\n detected", gjson.Get(stat, "System.EventRecordID"))
	}
}
