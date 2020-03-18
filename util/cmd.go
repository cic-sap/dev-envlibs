package util

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

type ShellExecuteResult struct {
	Command  string
	Elapsed  time.Duration
	ExitCode int
	Stdout   string
	Stderr   string
	Error    error /* returned by cmd.Run() */
}

func (r *ShellExecuteResult) String() string {
	return fmt.Sprintf("exec: %s exited with %d, duration: %s, stdout: %s, stderr: %s, error: %s",
		r.Command, r.ExitCode, r.Elapsed, r.Stdout, r.Stderr, r.Error)
}



type ExecCmd struct {
	currentPath string
}

func NewExecCmd(cp string) *ExecCmd{
	return &ExecCmd{currentPath:cp}
}

func (e *ExecCmd)CD(path string) {
	e.currentPath = path
}
func (e *ExecCmd) RestPath() {
	currentPath = ""
}

func (e *ExecCmd)ExecWithEnv(ctx context.Context, commands string, env []string) *ShellExecuteResult {

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "bash", "-ec", commands)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if e.currentPath != "" {
		cmd.Dir = e.currentPath
	}

	if env != nil {
		cmd.Env = env
	}

	start := time.Now()
	err := cmd.Run() /* cmd will be killed if ctx is canceled or reaches timeout */
	elapsed := time.Since(start)

	rc := &ShellExecuteResult{}
	if timedout := ctx.Err(); timedout != nil {
		rc.Error = timedout
	} else {
		rc.Error = err
	}
	rc.Command = commands
	rc.Elapsed = elapsed
	if err != nil {
		rc.ExitCode = 255 /* unknown error */
		if exit, ok := err.(*exec.ExitError); ok {
			if status, ok := exit.Sys().(syscall.WaitStatus); ok {
				rc.ExitCode = status.ExitStatus()
			}
		}
	}

	rc.Stdout = stdout.String()
	rc.Stderr = stderr.String()

	return rc
}

func(e *ExecCmd) Exec(ctx context.Context, commands string) *ShellExecuteResult {
	return e.ExecWithEnv(ctx, commands, nil)
}

var currentPath string
func CD(path string) {
	currentPath = path
}
func RestPath() {
	currentPath = ""
}
func ExecWithEnv(ctx context.Context, commands string, env []string) *ShellExecuteResult {

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "bash", "-ec", commands)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if currentPath != "" {
		cmd.Dir = currentPath
	}

	if env != nil {
		cmd.Env = env
	}

	start := time.Now()
	err := cmd.Run() /* cmd will be killed if ctx is canceled or reaches timeout */
	elapsed := time.Since(start)

	rc := &ShellExecuteResult{}
	if timedout := ctx.Err(); timedout != nil {
		rc.Error = timedout
	} else {
		rc.Error = err
	}
	rc.Command = commands
	rc.Elapsed = elapsed
	if err != nil {
		rc.ExitCode = 255 /* unknown error */
		if exit, ok := err.(*exec.ExitError); ok {
			if status, ok := exit.Sys().(syscall.WaitStatus); ok {
				rc.ExitCode = status.ExitStatus()
			}
		}
	}

	rc.Stdout = stdout.String()
	rc.Stderr = stderr.String()

	return rc
}

func Exec(ctx context.Context, commands string) *ShellExecuteResult {
	return ExecWithEnv(ctx, commands, nil)
}
