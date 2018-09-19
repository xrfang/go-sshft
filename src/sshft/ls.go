package sshft

import (
	"path/filepath"
	"strconv"
	"strings"
)

type fsAuth struct {
	Read    bool
	Write   bool
	Execute bool
}

func (a *fsAuth) parse(rwx string) {
	a.Read = rwx[0] != '-'
	a.Write = rwx[1] != '-'
	a.Execute = rwx[2] != '-'
}

type FSEntry struct {
	Type      string
	Auths     [3]fsAuth
	Items     int
	Owner     string
	Group     string
	Size      int
	Timestamp int
	Name      string
	Target    string //for symbol-links only
	MimeInfo  string
}

func (c Client) parseEntry(entry []string) FSEntry {
	var fse FSEntry
	switch entry[0][0] {
	case '-':
		fse.Type = "file"
	case 'd':
		fse.Type = "directory"
	case 'l':
		fse.Type = "symlink"
	}
	fse.Auths[0].parse(entry[0][1:4])
	fse.Auths[1].parse(entry[0][4:7])
	fse.Auths[2].parse(entry[0][7:10])
	fse.Items, _ = strconv.Atoi(entry[1])
	fse.Owner = entry[2]
	fse.Group = entry[3]
	fse.Size, _ = strconv.Atoi(entry[4])
	fse.Timestamp, _ = strconv.Atoi(entry[5])
	if fse.Type == "symlink" {
		nt := c.rxs[1].Split(entry[6], 2)
		fse.Name = nt[0]
		if len(nt) > 1 {
			fse.Target = nt[1]
		}
	} else {
		fse.Name = entry[6]
	}
	return fse
}

func (c Client) getMimeInfo(path string) (map[string]string, error) {
	xr := c.sshExec("file -i "+path, nil)
	if xr.err != nil {
		return nil, xr
	}
	mi := make(map[string]string)
	for _, m := range strings.Split(xr.out, "\n") {
		v := strings.SplitN(m, ":", 2)
		if len(v) == 2 {
			mi[v[0]] = strings.TrimSpace(v[1])
		}
	}
	if strings.HasPrefix(mi[path], "inode/directory") {
		sub, err := c.getMimeInfo(filepath.Join(path, "*"))
		if err == nil {
			for k, v := range sub {
				mi[filepath.Base(k)] = v
			}
		}
	}
	return mi, nil
}

func (c *Client) List(path string) ([]FSEntry, error) {
	xr := c.sshExec(`ls -l --time-style="+%s" `+path, nil)
	if xr.err != nil {
		return nil, xr
	}
	var es []FSEntry
	for _, l := range strings.Split(xr.out, "\n") {
		s := c.rxs[0].Split(l, 7)
		if len(s) != 7 {
			continue
		}
		e := c.parseEntry(s)
		es = append(es, e)
	}
	if len(es) == 0 {
		return nil, nil
	}
	mi, err := c.getMimeInfo(path)
	if err == nil {
		for i := range es {
			es[i].MimeInfo = mi[es[i].Name]
		}
	}
	return es, nil
}
