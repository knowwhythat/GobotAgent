//go:build windows

package sys_exec

import (
	"os/exec"
	"syscall"
)

func BuildCmd(cmdStr ...string) *exec.Cmd {
	command := exec.Command(cmdStr[0], cmdStr[1:]...)
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return command
}

func BuildBashCmd(cmdStr ...string) *exec.Cmd {
	cmdArr := []string{"/c"}
	cmdArr = append(cmdArr, cmdStr...)
	command := exec.Command("cmd", cmdArr...)
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return command
}

func Setpgid(cmd *exec.Cmd) {
}

func KillProcess(cmd *exec.Cmd) error {
	if cmd != nil && cmd.Process != nil {
		if err := cmd.Process.Signal(syscall.SIGKILL); err != nil {
			if err := cmd.Process.Kill(); err != nil {
				return err
			}
		}
	}
	return nil
}
