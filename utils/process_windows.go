//go:build windows

package utils

import (
	"fmt"
	"os/exec"
	"syscall"
)

func UseProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func Terminate(cmd *exec.Cmd) error {
	// TODO: Ensure that also the child processes are terminated.
	// From https://go.googlesource.com/go/+/refs/heads/master/src/os/signal/signal_windows_test.go#17
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		return fmt.Errorf("loading kernel32.dll failed: %w", e)
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		return fmt.Errorf("finding GenerateConsoleCtrlEvent process failed: %w", e)
	}
	r, _, e := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(cmd.Process.Pid))
	if r == 0 {
		return fmt.Errorf("calling GenerateConsoleCtrlEvent failed: %w", e)
	}
	return nil
}
