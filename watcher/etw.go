package watcher

import (
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
	log.Printf("Session is being created: %v(via GUID %v)", etwSession, guid)
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
		log.Printf("Event: %v", eventProperties)
	}

	var etwWaitGroup sync.WaitGroup
	etwWaitGroup.Add(1)
	go func() {
		if err := etwSession.Process(etwSessionCallBack); err != nil {
			log.Panicf("Error processing ETW session: %s", err)
		}
		etwWaitGroup.Done()
	}()

	// Trap cancellation signal
	trapSignalChannel := make(chan os.Signal, 1)
	signal.Notify(trapSignalChannel, os.Interrupt)
	<-trapSignalChannel

	if err := etwSession.Close(); err != nil {
		log.Printf("Error closing ETW session: %s", err)
	}

	etwWaitGroup.Wait()
}
