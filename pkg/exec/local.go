package exec

import (
	"os"
	osExec "os/exec"

	"bytes"
)

// The idea behind this pkg is that there's a common
// mechanism for executing commands, a transport can
// be added later like SSH or something else
// ToDo: add a way to interrupt the command, via a chanel + signals
type Exec struct {
	Shell string
}

// Run executes a command, returns stderr, stdout as strings and error
// For now it does not support sudo.
func (e *Exec) Run(cmd string) (string, string, error) {
	env := os.Environ()
	c := osExec.Command(e.Shell, "-c", cmd)
	c.Stderr = nil
	c.Env = env

	var outbuf, errbuf bytes.Buffer
	c.Stdout = &outbuf
	c.Stderr = &errbuf

	err := c.Run()

	return outbuf.String(), errbuf.String(), err
}
