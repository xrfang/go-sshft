package main

import (
	"fmt"
	"os"
	"sshft"
)

func main() {
	path := "/"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	_ = path
	c := sshft.NewClient("pj")
	matches, err := c.Grep("go", sshft.GrepOption{
		SkipBinary: true,
		Recursive:  true,
		Pattern:    "strings",
	}, sshft.GrepOption{
		IgnoreCase: true,
		Pattern:    "print",
	})
	if err != nil {
		panic(err)
	}
	for i, m := range matches {
		fmt.Printf("%d: %s\n", i, m)
	}
}
