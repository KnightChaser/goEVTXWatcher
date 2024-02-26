package watcher

import (
	"fmt"
	"log"
	"os"
	"time"
)

func WatchFile(filepath string) {
	initialStat, err := os.Stat(filepath)
	if err != nil {
		log.Panicf("Error getting file info: %s", err)
	}

	go func() {
		for {
			stat, err := os.Stat(filepath)
			if err != nil {
				log.Fatalf("Error getting file info: %s", err)
			}

			if stat.ModTime() != initialStat.ModTime() || stat.Size() != initialStat.Size() {
				fmt.Printf("File has been modified. Size: %d, ModTime: %v\n", stat.Size(), stat.ModTime())
				initialStat = stat
				// EVTXBeats(filepath)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	// Prevent the upper go routine from exiting
	select {}
}
