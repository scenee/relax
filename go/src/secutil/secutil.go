package main

import (
	"certutil"
	"fmt"
	"os"
	"os/exec"
)

//
// Utils
//

func usage() {
	fmt.Println("usage: secutil install <commonname>")
	os.Exit(0)
}

//
// Main
//
func init() {
}

func main() {
	var (
		cmd string
		cn  string
	)

	if len(os.Args) != 3 {
		usage()
	}

	cmd = os.Args[1]
	cn = os.Args[2]

	if cmd == "" {
		usage()
	}

	switch cmd {
	case "install":
		if _, err := exec.Command("/usr/bin/security", "find-certificate", "-c", cn).Output(); err == nil {
			return
		}
		certutil.InstallCertificate(cn)
	default:
		usage()
	}
}
