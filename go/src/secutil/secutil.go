package main

import (
	"certutil"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

//
// Utils
//

func usage() {
	fmt.Println("usage: secutil [-k <keychain>] install <certifcate>")
	os.Exit(0)
}

//
// Main
//
func init() {
}

func main() {
	var (
		cmd      string
		cn       string
		keychain string
	)

	flag.StringVar(&keychain, "k", "", "Use the keychain rather than the default keychain")

	flag.Parse()

	cmd = flag.Arg(0)
	cn = flag.Arg(1)

	if cmd == "" || cn == "" {
		usage()
	}

	switch cmd {
	case "install":
		if _, err := exec.Command("/usr/bin/security", "find-certificate", "-c", cn, keychain).Output(); err == nil {
			return
		}
		certutil.InstallCertificate(cn, keychain)
	default:
		usage()
	}
}
