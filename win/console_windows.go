//go:build windows

package main

import (
	"os"

	"golang.org/x/sys/windows"
)

func attachConsole() {
	kernel := windows.NewLazySystemDLL("kernel32.dll")
	attached, _, _ := kernel.NewProc("AttachConsole").Call(^uintptr(0))
	if attached == 0 {
		kernel.NewProc("AllocConsole").Call()
	}
	handle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err == nil && handle != 0 {
		os.Stdout = os.NewFile(uintptr(handle), "stdout")
		os.Stderr = os.Stdout
	}
}
