package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	cmd := exec.Command(command, args...)

	// to avoid an error `open /dev/null: no such file or directory`
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	tmpDir := "/tmp/myfakehost"
	err := os.MkdirAll(tmpDir, 0744)
	if err != nil {
		panic("Failed to create tmp dir: " + err.Error())
	}

	err = exec.Command("mkdir", "-p", filepath.Join(tmpDir, filepath.Dir(command))).Run()
	if err != nil {
		panic("Command mkdir failed: " + err.Error())
	}

	err = exec.Command("cp", command, filepath.Join(tmpDir, command)).Run()
	if err != nil {
		panic("Command copy failed: " + err.Error())
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: tmpDir,
	}

	err = cmd.Run()
	if err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		os.Exit(exitCode)
	}
}
