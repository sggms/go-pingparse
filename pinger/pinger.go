package pinger

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/sggms/go-pingparse/pinger/parser"
)

// Ping will ping the specified IPv4 address wit the provided timeout and interval settings.
func Ping(ipV4Address string, interval, timeout time.Duration) (*parser.PingOutput, error) {
	var (
		output, errorOutput bytes.Buffer
		exitCode            int
	)

	cmd := exec.Command("ping", "-n", "-v", "-w", fmt.Sprintf("%d", int(timeout.Seconds())), "-i", fmt.Sprintf("%d", int(interval.Seconds())), ipV4Address)
	cmd.Stdout = &output
	cmd.Stderr = &errorOutput

	err := cmd.Run()
	if err == nil {
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	} else {
		exitCode, err = parseExitCode(err)
		if err != nil {
			return nil, err
		}
	}

	// try to parse output also in case of failure
	po, err := parser.Parse(output.String())
	if err == nil {
		return po, nil
	}

	// in case of error, use also the execution context errors (if any)
	return nil, fmt.Errorf("exit code: %d\nparse error: %v\nstdout:\n%s\nstderr:\n%s", exitCode, err, output.String(), errorOutput.String())
}

func parseExitCode(err error) (int, error) {
	// try to get the exit code
	if exitError, ok := err.(*exec.ExitError); ok {
		ws := exitError.Sys().(syscall.WaitStatus)
		return ws.ExitStatus(), nil
	}

	// This will happen (in OSX) if `name` is not available in $PATH,
	// in this situation, exit code could not be get, and stderr will be
	// empty string very likely, so we use the default fail code, and format err
	// to string and set to stderr
	return 0, fmt.Errorf("could not get exit code for failed program: %v", err)
}
