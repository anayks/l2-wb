package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra/cobra/cmd"
)

type Command interface {
	execute(string) (string, error)
}

type CommandPipe []Command

type BasicCommand struct {
	name string
	args []string
}

func (c CommandPipe) execute(string) (stdout string, err error) {
	for _, v := range c {
		stdout, err = v.execute(stdout)
		if err != nil {
			return "", err
		}
	}
	return stdout, nil
}

type CdCommand struct {
	BasicCommand
}

func (c *CdCommand) execute(string) (string, error) {
	for _, v := range c.args {
		if err := os.Chdir(v); err != nil {
			return "", err
		}
	}

	return "", nil
}

type PwdCommand struct {
	BasicCommand
}

func (c *PwdCommand) execute(s string) (string, error) {
	cmd := exec.Command("pwd")

	initStdin(cmd, s)

	stdout, err := cmd.Output()

	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

type EchoCommand struct {
	BasicCommand
}

func (c *EchoCommand) execute(s string) (string, error) {
	cmd := exec.Command("echo", c.args...)

	initStdin(cmd, s)

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

type KillCommand struct {
	BasicCommand
}

func (c *KillCommand) execute(s string) (string, error) {
	cmd := exec.Command("kill", c.args...)

	initStdin(cmd, s)

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

type ProcessesCommand struct {
	BasicCommand
}

func (c *ProcessesCommand) execute(s string) (string, error) {
	cmd := exec.Command("ps", c.args...)

	initStdin(cmd, s)

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

type CatCommand struct {
	BasicCommand
}

func (c *CatCommand) execute(s string) (string, error) {
	cmd := exec.Command("cat", c.args...)

	initStdin(cmd, s)

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

func initStdin(cmd *exec.Cmd, s string) error {
	var b2 bytes.Buffer
	cmd.Stdin = &b2

	_, err := b2.WriteString(s)
	if err != nil {
		return err
	}
	return nil
}

type ForkCommand struct {
	BasicCommand
}

func (c *ForkCommand) execute(s string) (string, error) {
	id, err := syscall.ForkExec(os.Args[0], os.Args, &syscall.ProcAttr{
		Env: append(os.Environ()),
		Sys: &syscall.SysProcAttr{
			Setsid: true,
		},
		Files: []uintptr{0, 1, 2}, // print message to the same pty
	})

	if err != nil {
		return "", err
	}

	return string(id), nil
}

type ExecCommand struct {
	BasicCommand
}

func (c *ExecCommand) execute(s string) (string, error) {
	cmd := cmd.Execute(c.name, c.args)

	initStdin(cmd, s)

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return stdout, nil
}

func NewCommand(name string, args []string) (Command, error) {
	cmd := BasicCommand{
		name: name,
		args: args,
	}

	if name == "cd" {
		return &CdCommand{cmd}, nil
	} else if name == "pwd" {
		return &PwdCommand{cmd}, nil
	} else if name == "echo" {
		return &EchoCommand{cmd}, nil
	} else if name == "kill" {
		return &KillCommand{cmd}, nil
	} else if name == "ps" {
		return &ProcessesCommand{cmd}, nil
	} else if name == "cat" {
		return &CatCommand{cmd}, nil
	} else if name == "fork" {
		return &ForkCommand{cmd}, nil
	} else if name == "exec" {
		return &ExecCommand{cmd}, nil
	}

	return nil, fmt.Errorf("Command \"%v\" doesn't exist!", name)
}
