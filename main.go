package main

import (
	"goEVTXWatcher/watcher"
)

func main() {
	// if !watcher.CheckAdministratorPrivilege() {
	// 	watcher.RequestAdministratorPrivilege()
	// 	fmt.Println("You were requested to run the program as administrator")
	// } else {
	// 	fmt.Println("You have administrator privilege")
	// }

	guid := watcher.GetLogmanGUID("Microsoft-Windows-Sysmon")
	watcher.RunETWByGuid(guid)
}
