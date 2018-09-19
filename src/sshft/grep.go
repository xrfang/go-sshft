package sshft

import (
	"bytes"
	"errors"
	"strings"
)

type GrepOption struct {
	IgnoreCase  bool   //-i
	InvertMatch bool   //-v
	Matcher     string //basic(-G), fixed(-F), extended(-E)
	Pattern     string
	Recursive   bool //-r
	SkipBinary  bool //-I
	WholeWord   bool //-w
}

func (g GrepOption) cmdLine() string {
	if g.Pattern == "" {
		panic(errors.New("GrepOption: missing pattern"))
	}
	args := []string{"grep", "-l"} //always output filename only
	if g.IgnoreCase {
		args = append(args, "-i")
	}
	if g.InvertMatch {
		args = append(args, "-v")
	}
	switch g.Matcher {
	case "fixed":
		args = append(args, "-G")
	case "extended":
		args = append(args, "-E")
	}
	if g.Recursive {
		args = append(args, "-r")
	}
	if g.SkipBinary {
		args = append(args, "-I")
	}
	if g.WholeWord {
		args = append(args, "-w")
	}
	pattern := strings.Replace(g.Pattern, `"`, `\"`, -1)
	args = append(args, `"`+pattern+`"`)
	return strings.Join(args, " ")
}

func (c *Client) Grep(path string, searches ...GrepOption) (matches []string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	var sh bytes.Buffer
	if len(searches) == 0 {
		panic(errors.New("Grep: no searches given"))
	}
	path = strings.Replace(path, `"`, `\"`, -1)
	cmds := []string{searches[0].cmdLine() + ` "` + path + `"`}
	if len(searches) > 1 {
		for _, s := range searches[1:] {
			cmds = append(cmds, "xargs "+s.cmdLine())
		}
	}
	sh.WriteString(strings.Join(cmds, "|"))
	xr := c.sshExec("bash -s", &sh)
	assert(xr.err)
	for _, m := range strings.Split(xr.out, "\n") {
		m = strings.TrimSpace(m)
		if m != "" {
			matches = append(matches, m)
		}
	}
	return
}
