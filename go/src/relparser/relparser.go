package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/DHowett/go-plist"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"strings"
)

//
// Global variables
//

const PREFIX string = "REL_CONFIG"

var logger *log.Logger

//
// Utils
//

func usage() {
	fmt.Println("usage: relparser [-f <Relfile>] list")
	fmt.Println("       relparser [-f <Relfile>] [-o <output>] source <distribution>")
	fmt.Println("       relparser [-f <Relfile>] [-o <output>] [-plist <Info.plist>] plist <distribution>")
	fmt.Println("       relparser [-f <Relfile>] [-o <output>] -plist <Info.plist> export_options <distribution>")
	os.Exit(0)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func loadRelfile(path string) (*Relfile, error) {
	var (
		c      *exec.Cmd
		f      *os.File
		err    error
		buffer bytes.Buffer
		rel    Relfile
	)

	c = exec.Command("envsubst")

	f, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	c.Stdin = f

	c.Stdout = &buffer
	err = c.Run()
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(buffer.Bytes(), &rel); err != nil {
		return nil, err
	}

	return &rel, nil
}

func genSourceline(key, value string) string {
	k := strings.Join([]string{PREFIX, key}, "_")
	return fmt.Sprintf("export %v=\"%v\"\n", k, value)
}

func genSourceLine2(name string, key string, value interface{}) string {
	k := strings.Join([]string{PREFIX, name, key}, "_")
	return fmt.Sprintf("export %v=\"%v\"\n", k, value)
}

func getBundleID(path string) string {
	var decoder *plist.Decoder

	if f, err := os.Open(path); err != nil {
		logger.Fatalf("open error:", err)
	} else {
		defer f.Close()
		decoder = plist.NewDecoder(f)
	}

	var data map[string]interface{}

	if err := decoder.Decode(&data); err != nil {
		logger.Fatalf("decode error:", err)
	}

	if props, ok := data["ApplicationProperties"].(map[string]interface{}); ok {
		return props["CFBundleIdentifier"].(string)
	} else {
		return data["CFBundleIdentifier"].(string)
	}
}

//
// Relfile
//

type Relfile struct {
	Version       string
	Xcodeproj     string
	Workspace     string
	Uploader      map[string]interface{}
	Log_formatter string
	Distributions map[string]Distribution
}

func (r Relfile) list() {
	for k, _ := range r.Distributions {
		fmt.Println(k)
	}
}

func (r Relfile) prepare_out(out string) (f *os.File) {
	_, err := os.Stat(out)
	if os.IsExist(err) {
		os.Remove(out)
	}
	f, err = os.Create(out)
	checkErr(err)
	f.Close()

	f, err = os.OpenFile(out, os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func (r Relfile) genOptionsPlist(dist string, infoPlist, out string) {
	d := r.Distributions[dist]

	of := r.prepare_out(out)
	defer of.Close()

	d.writeExportOptions(infoPlist, of)
}

func (r Relfile) genPlist(dist string, in string, out string) {
	d := r.Distributions[dist]

	of := r.prepare_out(out)
	defer of.Close()

	d.writeInfoPlist(in, of)
}

func (r Relfile) genSrouce(dist string, out string) {
	d := r.Distributions[dist]

	of := r.prepare_out(out)
	defer of.Close()

	r.writeSource(of)
	d.writeSource(dist, of)
}

func (r Relfile) writeSource(out *os.File) {
	var source string

	source += genSourceline("xcodeproj", r.Xcodeproj)
	source += genSourceline("workspace", r.Workspace)
	source += genSourceline("log_formatter", r.Log_formatter)

	for k, v := range r.Uploader {
		up := v.(map[interface{}]interface{})
		for name, value := range up {
			//fmt.Printf("---\t%v: %v\n", name, value)
			source += genSourceLine2("uploader_"+k, name.(string), value.(string))
		}
	}

	if _, err := out.WriteString(source); err != nil {
		panic(err)
	}
}

//
// Distribution
//

type Distribution struct {
	Scheme               string
	Sdk                  string
	Configuration        string
	Provisioning_profile string
	Team_id              string
	Bundle_identifier    string
	Bundle_version       string
	Version              string
	Build_settings       map[string]interface{}
	Info_plist           map[string]interface{}
	Export_options       ExportOptions
}

func (d Distribution) writeInfoPlist(basePlistPath string, out *os.File) {
	var decoder *plist.Decoder

	if f, err := os.Open(basePlistPath); err != nil {
		//fmt.Println("open error:", err)
		return
	} else {
		defer f.Close()
		decoder = plist.NewDecoder(f)
	}

	var data map[string]interface{}

	if err := decoder.Decode(&data); err != nil {
		//fmt.Println("decode error:", err)
		panic(err)
	}

	//fmt.Println(data)
	//fmt.Println("--- Info.plist")

	for k, v := range d.Info_plist {
		data[k] = v
		//fmt.Printf("---\t%v: %v\n", k, v)
	}

	encoder := plist.NewEncoder(out)
	encoder.Indent("\t")
	if err := encoder.Encode(data); err != nil {
		//fmt.Println("encode error: ", err)
		panic(err)
	}
}

func (d Distribution) writeExportOptions(infoPlist string, out *os.File) {

	fmt.Println("export options:", d.Export_options)

	var bundleID string
	// get bundle identifier
	if d.Bundle_identifier == "" {
		bundleID = getBundleID(infoPlist)
	} else {
		bundleID = d.Bundle_identifier
	}

	encoder := plist.NewEncoder(out)
	encoder.Indent("\t")

	opts := d.Export_options
	opts.TeamID = d.Team_id
	opts.Encode(encoder, infoPlist, bundleID)

}

func (d Distribution) writeSource(name string, out *os.File) {

	var source string
	source += genSourceLine2(name, "scheme", d.Scheme)
	source += genSourceLine2(name, "sdk", d.Sdk)
	source += genSourceLine2(name, "configuration", d.Configuration)
	source += genSourceLine2(name, "provisioning_profile", d.Provisioning_profile)
	source += genSourceLine2(name, "team_id", d.Team_id)
	source += genSourceLine2(name, "bundle_identifier", d.Bundle_identifier)
	source += genSourceLine2(name, "bundle_version", d.Bundle_version)
	source += genSourceLine2(name, "version", d.Version)

	// "--- Build settings\n"
	build_settings := strings.Join([]string{PREFIX, name, "build_settings"}, "_")

	source += fmt.Sprintf("%v=()\n", build_settings)

	for _, vars := range []map[string]interface{}{d.Build_settings, d.Info_plist} {
		for k, v := range vars {
			switch _v := v.(type) {
			default:
				source += fmt.Sprintf("%v+=(%v='%v')\n", build_settings, k, v)
			case []interface{}:
				var ss []string
				for _, s := range _v {
					ss = append(ss, fmt.Sprintf("%v", s))
				}
				source += fmt.Sprintf("%v+=(%v='%v')\n", build_settings, k, strings.Join(ss, "{}"))
			}
		}
	}
	source += fmt.Sprintf("export %v\n", build_settings)

	if _, err := out.WriteString(source); err != nil {
		panic(err)
	}
}

//
// ExportOptions
//

type ExportOptions struct {
	ProvisioningProfiles map[string]string `plist:"provisioningProfiles", omitempty`
	CompileBitcode       bool              `plist:"compileBitcode", omitempty`
	Method               string            `plist:"method"`
	SigningCertificate   string            `plist:"signingCertificate", omitempty`
	SigningStyle         string            `plist:"signingStyle", omitempty`
	TeamID               string            `plist:"teamID"`
	Provisioning_profile string            `plist:"-"`
}

func (opts ExportOptions) Encode(encoder *plist.Encoder, infoPlist string, bundleID string) {
	if bundleID == "" {
		logger.Fatalf("bundleID is empty")
	}

	if opts.Provisioning_profile != "" {
		opts.ProvisioningProfiles = map[string]string{bundleID: opts.Provisioning_profile}
		opts.SigningStyle = "manual"

		switch opts.Method {
		case "development":
			opts.SigningCertificate = "iPhone Developer"
		default:
			opts.SigningCertificate = "iPhone Distribution"
		}
	}

	if err := encoder.Encode(&opts); err != nil {
		//fmt.Println("encode error: ", err)
		panic(err)
	}
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

		rel *Relfile
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

	rel, err = loadRelfile(in)
	if err != nil {
		logger.Fatal(err)
	}

	if rel.Version == "" || len(rel.Distributions) == 0 {
		logger.Fatalf("Please update your Relfile format to 2.x. See https://github.com/SCENEE/relax#relfile")
	}

	switch cmd {
	case "source":
		rel.genSrouce(dist, out)
	case "plist":
		rel.genPlist(dist, infoPlist, out)
	case "export_options":
		if infoPlist == "" {
			logger.Fatalf("Pass a Info.plist path using '-plist' option")
		}
		rel.genOptionsPlist(dist, infoPlist, out)
	case "list":
		rel.list()
	default:
		usage()
	}
}
