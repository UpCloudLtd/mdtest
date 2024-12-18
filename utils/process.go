//go:build !windows

package utils

import (
	"fmt"
	"os/exec"
	"syscall"
)

func UseProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func Terminate(cmd *exec.Cmd) error {
	err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("failed to terminate process group: %w", err)
	}
	return nil
}
