package watcher

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"syscall"
)

// GetLogmanGUID returns the GUID of the specified provider
func GetLogmanGUID(provider string) string {
	command := fmt.Sprintf("logman query providers | findstr %s", provider)
	cmd := exec.Command("cmd", "/C", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	logmanOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error getting GUID: %s", err)
	}

	logmanOutputString := string(logmanOutput[:])
	logmanOutputGUIDRegex := regexp.MustCompile(`\{([0-9A-Fa-f\-]+)\}`)
	matches := logmanOutputGUIDRegex.FindStringSubmatch(logmanOutputString)
	if len(matches) > 1 {
		guidValue := matches[1]
		return guidValue
	} else {
		log.Panicf("Error getting GUID: %s", logmanOutputString)
		return ""
	}
}
