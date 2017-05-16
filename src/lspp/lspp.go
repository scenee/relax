package main

import (
	"fmt"
	"github.com/DHowett/go-plist"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"sort"
	"strings"
	"time"
)

type ProvisioningProfile struct {
	AppIDName      string    `plist:AppIDName`
	TeamName       string    `plist:TeamName`
	UUID           string    `plist:UUID`
	Name           string    `plist:Name`
	TeamIdentifier []string  `plist:TeamIdentifier`
	CreationDate   time.Time `plist:CreationDate`
}

func (p ProvisioningProfile) GetTeamIdentifier() (s string) {
	if len(p.TeamIdentifier) == 0 {
		return ""
	} else {
		return p.TeamIdentifier[0]
	}
}

type Result struct {
	Data ProvisioningProfile
	Path string
	Name string
}

type ByNameLatest []*Result

func (p ByNameLatest) Len() int      { return len(p) }
func (p ByNameLatest) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ByNameLatest) Less(i, j int) bool {
	if p[i].Data.TeamName == p[i].Data.TeamName {
		if p[i].Data.Name == p[j].Data.Name {
			return p[i].Data.CreationDate.Unix() >= p[j].Data.CreationDate.Unix()
		} else {
			return p[i].Data.Name < p[j].Data.Name
		}
	} else {
		return p[i].Data.TeamName < p[i].Data.TeamName
	}
}

func parse(in string, out string, c chan *Result) {
	defer func() {
		if err := os.Remove(out); err != nil {
			logErr.Println("remove error:", err)
		}
	}()

	if _, err := exec.Command("/usr/bin/security", "cms", "-D", "-i", in, "-o", out).Output(); err != nil {
		logErr.Println("command error:", err)
		c <- nil
		return
	}

	var err error
	var file *os.File
	var decoder *plist.Decoder

	file, err = os.Open(out)
	if err != nil {
		logErr.Println("open error:", err)
		c <- nil
		return
	}
	defer file.Close()

	decoder = plist.NewDecoder(file)

	var data ProvisioningProfile

	if err := decoder.Decode(&data); err != nil {
		logErr.Println("decode error:", err)
		c <- nil
		return
	}

	comps := strings.Split(in, "/")
	name := comps[len(comps)-1]

	ret := Result{Data: data, Path: in, Name: name}
	c <- &ret
}

var logErr *log.Logger

func init() {
	logErr = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

func main() {

	usr, _ := user.Current()
	var PP_PATH string = usr.HomeDir + "/Library/MobileDevice/Provisioning Profiles"

	var isLatest, isVerborse bool
	var team, pattern string

	args := os.Args
	for i, arg := range args {
		switch arg {
		case "--latest":
			isLatest = true
		case "-v":
			isVerborse = true
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

	//fmt.Println(isLatest)
	//fmt.Println(isVerborse)
	//fmt.Println(team)
	//fmt.Println(pattern)

	infos, err := ioutil.ReadDir(PP_PATH)
	if err != nil {
		logErr.Fatalf("error: %v", err)
	}

	temp, _ := ioutil.TempDir("", "relax/lspp")
	err = os.MkdirAll(temp, 0700)
	if err != nil {
		logErr.Fatalf("error: %v", err)
	}

	s := make(chan bool, 32)
	c := make(chan *Result, len(infos))
	count := 0
	for _, info := range infos {
		s <- true
		name := info.Name()
		if false == strings.HasSuffix(name, "mobileprovision") {
			continue
		}
		in := PP_PATH + "/" + name
		out := temp + "/" + name
		go func() {
			defer func() { <-s }()
			parse(in, out, c)
		}()
		count++
	}

	var rets []*Result

	for i := 0; i < count; i++ {
		ret := <-c
		if ret == nil {
			continue
		}

		if team != "" && team != ret.Data.GetTeamIdentifier() {
			continue
		}

		if pattern != "" {
			if matched, err := regexp.MatchString(pattern, ret.Data.Name); err != nil || !matched {
				continue
			}
		}
		rets = append(rets, ret)
	}

	sort.Sort(ByNameLatest(rets))

	names := ""
	for _, ret := range rets {
		data := ret.Data
		name := ret.Name
		if isLatest && strings.Contains(names, ret.Data.Name) {
			continue
		}
		names += "\t" + ret.Data.Name
		var s string
		teamId := data.GetTeamIdentifier()
		if isVerborse {
			s = fmt.Sprintf("%v/%v -- %v %v %v %v", PP_PATH, name, data.CreationDate.Local().Format("2006-01-02 15:04:05"), teamId, data.TeamName, data.Name)
		} else {
			s = fmt.Sprintf("%v/%v", PP_PATH, name)
		}
		fmt.Printf("%s\n", s)
	}
}
