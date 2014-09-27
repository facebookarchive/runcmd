// Package runcmd provides a convenience Run function for exec.Cmd which
// includes the original command along with the output and error streams.
package runcmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CommandError is returned by run and provides an error string which includes
// the full command and the contents of the output and error streams.
type CommandError struct {
	fullCommand string
	streams     *Streams
}

// Error provides a helpful error text.
func (e *CommandError) Error() string {
	return fmt.Sprintf(
		"error executing: %s:\n%s\n%s",
		e.fullCommand,
		e.streams.Stderr().Bytes(),
		bytes.TrimSpace(e.streams.Stdout().Bytes()),
	)
}

// Streams provides access to the output and error buffers.
type Streams struct {
	out *bytes.Buffer
	err *bytes.Buffer
}

// Stderr returns the underlying buffer with the contents of the error stream.
func (s *Streams) Stderr() *bytes.Buffer {
	return s.err
}

// Stdout returns the underlying buffer with the contents of the output stream.
func (s *Streams) Stdout() *bytes.Buffer {
	return s.out
}

// Run the command and return the associated streams, and error if any. The
// error may be a CommandError.
func Run(cmd *exec.Cmd) (*Streams, error) {
	var bout, berr bytes.Buffer
	streams := &Streams{
		out: &bout,
		err: &berr,
	}
	cmd.Stdout = &bout
	cmd.Stderr = &berr
	if err := cmd.Run(); err != nil {
		return streams, &CommandError{
			fullCommand: cmd.Path + " " + strings.Join(cmd.Args, " "),
			streams:     streams,
		}
	}
	return streams, nil
}
