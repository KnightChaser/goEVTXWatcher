package watcher

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"sync"

	"github.com/bi-zone/etw"
	"golang.org/x/sys/windows"
)

func RunETWByGuid(guidInput string) {
	// Subscribe to the event via the given GUID
	guid, err := windows.GUIDFromString(fmt.Sprintf("{%s}", guidInput))
	if err != nil {
		log.Panicf("Error parsing GUID: %s", err)
		return
	}

	fmt.Printf("Trying to create ETW session for GUID: %s\n", guidInput)
	etwSession, err := etw.NewSession(guid)
	if err != nil {
		log.Panicf("Error creating ETW session: %s", err)
		return
	}
	etwSessionName, err := extractEtwSessionName(etwSession)
	if err != nil {
		log.Panicf("Error extracting ETW session name: %s", err)
	}
	log.Printf("ETW session created and registered: %s\n", etwSessionName)

	// A callback function for ETW
	etwSessionCallBack := etwSessionCallback
	var etwWaitGroup sync.WaitGroup
	etwWaitGroup.Add(1)

	// Goroutine to handle termination signal
	go func() {
		// handling signal
		trapCancelSignal := make(chan os.Signal, 1)
		signal.Notify(trapCancelSignal, os.Interrupt)
		<-trapCancelSignal // Block until a signal is received

		// Terminate the ETW session
		log.Println("Received termination signal. Closing ETW session.")
		log.Printf("Terminate ETW session: %s\n", etwSessionName)
		terminateEtwSession(etwSessionName)
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

// Directly extract the ETW session name from the etw.Session object
// It looks like the official method is absent, so we have to use a workaround
func extractEtwSessionName(etwSession *etw.Session) (string, error) {
	etwSessionString := fmt.Sprintf("%v", etwSession)
	// According to the package implementation, the etw.Session object is printed as "go-etw-{random_guid}"
	etwSessionNameMatchPattern := `go-etw-\{[^\}]+\}`

	re, err := regexp.Compile(etwSessionNameMatchPattern)
	if err != nil {
		return "", err
	}

	match := re.FindString(etwSessionString)
	if match == "" {
		return "", fmt.Errorf("no match found for pattern: %s", etwSessionNameMatchPattern)
	}

	return match, nil
}

// A callback function for ETW
func etwSessionCallback(etwEvent *etw.Event) {
	log.Printf("=> Event #%v just has been captured\n", etwEvent.Header.ID)
}

// Terminate the ETW session by the given etw session name
func terminateEtwSession(etwSessionName string) error {
	logmanCommand := fmt.Sprintf("logman stop \"%s\" -ets", etwSessionName)
	terminateCommand := exec.Command("powershell", "-Command", logmanCommand)
	err := terminateCommand.Run()
	if err != nil {
		return fmt.Errorf("error terminating ETW session: %s", err)
	}

	return nil
}
