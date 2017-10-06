package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"relparser/relfile"
)

//
// Global variables
//

var logger *log.Logger

//
// Utils
//

func usage() {
	fmt.Println("usage: relparser [-f <Relfile>] list")
	fmt.Println("       relparser [-f <Relfile>] [-o <output>] source <distribution>")
	fmt.Println("       relparser [-f <Relfile>] build_options <distribution>")
	fmt.Println("       relparser [-f <Relfile>] [-o <output>] [-plist <Info.plist>] plist <distribution>")
	fmt.Println("       relparser [-f <Relfile>] [-o <output>] -plist <Info.plist> export_options <distribution>")
	os.Exit(0)
}

//
// Main
//
func init() {
	logger = log.New(os.Stderr, "error: ", 0)
}

func main() {
	cur, _ := os.Getwd()

	var (
		in        string
		out       string
		infoPlist string
		cmd       string
		dist      string

		rel *relfile.Relfile
		err error
	)
	flag.StringVar(&in, "f", cur+"/Relfile", "A Relfile path")
	flag.StringVar(&out, "o", cur+"/out", "An output path")
	flag.StringVar(&infoPlist, "plist", "", "An Info plist")

	flag.Parse()

	cmd = flag.Arg(0)
	dist = flag.Arg(1)

	//fmt.Println("in", in)
	//fmt.Println("out", out)
	//fmt.Println("command", command)
	//fmt.Println("infoPlist", infoPlist)
	//fmt.Println("dist", dist)

	if cmd == "" {
		usage()
	}

	rel, err = relfile.LoadRelfile(in)
	if err != nil {
		logger.Fatal(err)
	}

	if rel.Version == "" || len(rel.Distributions) == 0 {
		logger.Fatalf("Please update your Relfile format to 2.x. See https://github.com/SCENEE/relax#relfile")
	}

	switch cmd {
	case "source":
		rel.GenSrouce(dist, out)

	case "plist":
		rel.GenPlist(dist, infoPlist, out)

	case "build_options":
		rel.PrintBuildOptions(dist)

	case "export_options":
		if infoPlist == "" {
			logger.Fatalf("Pass a Info.plist path using '-plist' option")
		}
		rel.GenOptionsPlist(dist, infoPlist, out)

	case "list":
		rel.List()

	default:
		usage()

	}
}
