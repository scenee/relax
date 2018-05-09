package main

import (
	"fmt"
	"log"
	"os"
	"relparser/relfile"
	"sort"
	"strings"
)

type ByNameLatest []*relfile.ProvisioningProfileInfo

func (p ByNameLatest) Len() int      { return len(p) }
func (p ByNameLatest) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ByNameLatest) Less(i, j int) bool {
	if p[i].Pp.TeamName == p[i].Pp.TeamName {
		if p[i].Pp.Name == p[j].Pp.Name {
			return p[i].Pp.CreationDate.Unix() >= p[j].Pp.CreationDate.Unix()
		} else {
			return p[i].Pp.Name < p[j].Pp.Name
		}
	} else {
		return p[i].Pp.TeamName < p[i].Pp.TeamName
	}
}

var logErr *log.Logger

func init() {
	logErr = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

func findProvisioningProfile(pattern string, team string, isLatest bool, isVerbose bool) (outs []string) {
	infos := relfile.FindProvisioningProfile(pattern, team)

	sort.Sort(ByNameLatest(infos))

	names := ""
	var ret []string
	for _, info := range infos {
		if isLatest && strings.Contains(names, info.Pp.Name) {
			continue
		}
		names += "\t" + info.Pp.Name
		var s string
		teamId := info.Pp.TeamID()
		if isVerbose {
			s = fmt.Sprintf("%v -- %v %v %v %v", info.Name, info.Pp.CreationDate.Local().Format("2006-01-02 15:04:05"), teamId, info.Pp.TeamName, info.Pp.Name)
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
