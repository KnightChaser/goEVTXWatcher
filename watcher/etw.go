package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/bi-zone/etw"
	"golang.org/x/sys/windows"
)

func RunETWByGuid(guidInput string) {
	// Subscribe to the event via the given GUID
	guid, _ := windows.GUIDFromString(fmt.Sprintf("{%s}", guidInput))
	etwSession, err := etw.NewSession(guid)
	log.Printf("Session is being created (via GUID %v)\n", guid)
	if err != nil {
		log.Panicf("Error creating ETW session: %s", err)
		return
	}

	// A callback function for ETW
	etwSessionCallBack := func(etwEvent *etw.Event) {
		// Print the event
		eventProperties, err := etwEvent.EventProperties()
		if err != nil {
			log.Printf("Error getting event properties: %s", err)
		}

		jsonifiedEventProperties, err := json.Marshal(eventProperties)
		if err != nil {
			log.Printf("Error marshalling event properties: %s", err)
			return
		}

		log.Printf("Event: %s", string(jsonifiedEventProperties))
	}

	var etwWaitGroup sync.WaitGroup
	etwWaitGroup.Add(1)

	// Goroutine to handle termination signal
	go func() {
		trapCancelSignal := make(chan os.Signal, 1)
		signal.Notify(trapCancelSignal, os.Interrupt)
		<-trapCancelSignal // Wait for the termination signal
		log.Println("Received termination signal. Closing ETW session.")
		if err := etwSession.Close(); err != nil {
			log.Printf("Error closing ETW session: %s", err)
		}
		etwWaitGroup.Done()
	}()

	go func() {
		if err := etwSession.Process(etwSessionCallBack); err != nil {
			log.Panicf("Error processing ETW session: %s", err)
		}
		etwWaitGroup.Done()
	}()

	etwWaitGroup.Wait()
}
