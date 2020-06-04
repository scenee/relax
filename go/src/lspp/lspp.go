package main

import (
	"fmt"
	"log"
	"os"
	"relparser/relfile"
	"strings"
)

var logErr *log.Logger

func init() {
	logErr = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

func findProvisioningProfile(pattern string, team string, isLatest bool, isVerbose bool) (outs []string) {
	infos := relfile.FindProvisioningProfile(pattern, team)

	names := ""
	var ret []string
	for _, info := range infos {
		if isLatest && strings.Contains(names, info.Pp.Name) {
			continue
		}
		names += "\t" + info.Pp.Name
		var s string
		teamID := info.Pp.TeamID()
		if isVerbose {
			s = fmt.Sprintf("%v -- %v %v %v %v(%v)",
				info.Name,
				info.Pp.ExpirationDate.Local().Format("2006-01-02 15:04:05"),
				teamID,
				info.Pp.ProvisioningType(),
				info.Pp.Name,
				info.Pp.TeamName,
			)
		} else {
			s = fmt.Sprintf("%v", info.Name)
		}
		ret = append(ret, s)
	}
	return ret
}

func main() {
	var isLatest, isVerbose bool
	var team, pattern string

	args := os.Args
	for i, arg := range args {
		switch arg {
		case "--latest":
			isLatest = true
		case "-v":
			isVerbose = true
		case "--team":
			team = "found"
		default:
			if i > 0 {
				if team == "found" {
					team = arg
				} else {
					pattern = arg
				}
			}
		}
	}

	if team == "found" {
		team = ""
	}

	for _, out := range findProvisioningProfile(pattern, team, isLatest, isVerbose) {
		fmt.Printf("%s\n", out)
	}

}
