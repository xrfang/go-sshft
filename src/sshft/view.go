package sshft

import (
	"fmt"
	"strings"
)

func (c *Client) Cat(path string) (string, error) {
	path = strings.Replace(path, `"`, `\"`, -1)
	xr := c.sshExec(`cat "`+path+`"`, nil)
	if xr.err != nil {
		return "", xr
	}
	return xr.out, nil
}

func (c *Client) Head(path string, lines int) (string, error) {
	path = strings.Replace(path, `"`, `\"`, -1)
	head := `head "` + path + `"`
	if lines > 0 {
		head = fmt.Sprintf("%s -n %d", head, lines)
	}
	xr := c.sshExec(head, nil)
	if xr.err != nil {
		return "", xr
	}
	return xr.out, nil
}

func (c *Client) Tail(path string, lines int) (string, error) {
	path = strings.Replace(path, `"`, `\"`, -1)
	tail := `tail "` + path + `"`
	if lines > 0 {
		tail = fmt.Sprintf("%s -n %d", tail, lines)
	}
	xr := c.sshExec(tail, nil)
	if xr.err != nil {
		return "", xr
	}
	return xr.out, nil
}
