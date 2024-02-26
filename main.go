package main

import (
	"fmt"
	"goEVTXWatcher/watcher"
)

func main() {
	if !watcher.CheckAdministratorPrivilege() {
		watcher.RequestAdministratorPrivilege()
		fmt.Println("You were requested to run the program as administrator")
	} else {
		fmt.Println("You have administrator privilege")
	}

	const filePath = "C:\\Windows\\System32\\Winevt\\Logs\\Microsoft-Windows-Sysmon%4Operational.evtx"

	watcher.WatchFile(filePath)
}
