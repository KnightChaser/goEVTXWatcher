package watcher

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

// RequestAdministratorPrivilege requests the administrator privilege
func RequestAdministratorPrivilege() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	showCmd := int32(windows.SW_NORMAL)

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

// CheckAdministratorPrivilege checks if the program has administrator privilege
// by conducting a simple test of opening a handle to the physical drive that is not allowed to be opened without administrator privilege
func CheckAdministratorPrivilege() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")

	return err == nil
}
