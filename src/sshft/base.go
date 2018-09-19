package sshft

import (
	"bytes"
	"io"
	"os/exec"
	"regexp"
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

//Client is the SSH file transfer client.  This client works via the
//ssh command installed on system, which is presumably Linux.
type Client struct {
	user string
	addr string
	pkey string
	path string
	rxs  []*regexp.Regexp
}

//NewClient returns a *Client with its properties populated by the arguments
//given. addr is the hostname or ip address to connect, an optional port could
//be specified using host:port notation. If args are given, the first one is
//ssh login name, second one is identity, i.e. private key.
func NewClient(addr string, args ...string) *Client {
	var login, pkey string
	if len(args) > 0 {
		login = args[0]
	}
	if len(args) > 1 {
		pkey = args[1]
	}
	rxs := []*regexp.Regexp{
		regexp.MustCompile(`\s+`),
		regexp.MustCompile(` -> `),
	}
	return &Client{user: login, addr: addr, pkey: pkey, rxs: rxs}
}

func (c *Client) sshArgs(cmd string) []string {
	var cmdline []string
	if c.pkey != "" {
		cmdline = append(cmdline, "-i", c.pkey)
	}
	if c.user == "" {
		cmdline = append(cmdline, c.addr)
	} else {
		cmdline = append(cmdline, c.user+"@"+c.addr)
	}
	return append(cmdline, cmd)
}

type execResult struct {
	err error
	out string
}

func (xr execResult) Error() string {
	return xr.err.Error() + "\n" + xr.out
}

func (c *Client) sshExec(command string, stdin io.Reader) execResult {
	var buf bytes.Buffer
	cmd := exec.Command("ssh", c.sshArgs(command)...)
	cmd.Stdin = stdin
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return execResult{err, buf.String()}
}
