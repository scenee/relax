package main

import (
	"flag"
	"fmt"
	"github.com/DHowett/go-plist"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const PREFIX string = "REL_CONFIG"

func get_source_line(key, value string) string {
	k := strings.Join([]string{PREFIX, key}, "_")
	return fmt.Sprintf("export %v=\"%v\"\n", k, value)
}

func get_source_line2(name string, key string, value interface{}) string {
	k := strings.Join([]string{PREFIX, name, key}, "_")
	return fmt.Sprintf("export %v=\"%v\"\n", k, value)
}

type Relfile struct {
	Xcodeproj     string
	Workspace     string
	Uploader      string
	Log_formatter string
	Distributions map[string]Distribution
}

func (r Relfile) list() {
	for k, _ := range r.Distributions {
		fmt.Println(k)
	}
}

func (r Relfile) prep_out(out string) (f *os.File) {
	_, err := os.Stat(out)
	if os.IsExist(err) {
		os.Remove(out)
	}
	f, err = os.Create(out)
	checkError(err)
	f.Close()

	f, err = os.OpenFile(out, os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func (r Relfile) gen_plist(dist string, in string, out string) {
	d := r.Distributions[dist]

	of := r.prep_out(out)
	defer of.Close()

	d.make_plist(in, of)
}

func (r Relfile) gen_source(dist string, out string) {
	d := r.Distributions[dist]

	of := r.prep_out(out)
	defer of.Close()

	r.make_source(of)
	d.make_source(dist, of)
}

func (r Relfile) make_source(out *os.File) {
	var source string

	source += get_source_line("xcodeproj", r.Xcodeproj)
	source += get_source_line("workspace", r.Workspace)
	source += get_source_line("uploader", r.Uploader)
	source += get_source_line("log_formatter", r.Log_formatter)

	if _, err := out.WriteString(source); err != nil {
		panic(err)
	}
}

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
	Export_options       map[string]interface{}
}

func (d Distribution) make_plist(base string, out *os.File) {
	var decoder *plist.Decoder

	if f, err := os.Open(base); err != nil {
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

func (d Distribution) make_source(name string, out *os.File) {

	var source string
	source += get_source_line2(name, "scheme", d.Scheme)
	source += get_source_line2(name, "sdk", d.Sdk)
	source += get_source_line2(name, "configuration", d.Configuration)
	source += get_source_line2(name, "provisioning_profile", d.Provisioning_profile)
	source += get_source_line2(name, "team_id", d.Team_id)
	source += get_source_line2(name, "bundle_identifier", d.Bundle_identifier)
	source += get_source_line2(name, "bundle_version", d.Bundle_version)
	source += get_source_line2(name, "version", d.Version)

	// "--- Build settings\n"
	build_settings := strings.Join([]string{PREFIX, name, "build_settings"}, "_")

	source += fmt.Sprintf("%v=()\n", build_settings)

	for _, vars := range []map[string]interface{}{d.Build_settings, d.Info_plist, d.Export_options} {
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

	for k, v := range d.Export_options {
		source += get_source_line2(name, "export_options_"+k, v)
	}

	if _, err := out.WriteString(source); err != nil {
		panic(err)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func main() {
	cur, _ := os.Getwd()

	var (
		in     string
		out    string
		plbase string
	)
	flag.StringVar(&in, "f", cur+"/Relfile", "A output path")
	flag.StringVar(&out, "o", cur+"/out_relparser", "A output path")
	flag.StringVar(&plbase, "plist", "", "A Info plist")

	flag.Parse()

	cmd := flag.Arg(0)
	dist := flag.Arg(1)

	//fmt.Println("in", in)
	//fmt.Println("out", out)
	//fmt.Println("cmd", cmd)
	//fmt.Println("plist", plbase)
	//fmt.Println("dist", dist)

	var rel Relfile
	var data []byte
	var err error

	if data, err = ioutil.ReadFile(in); err != nil {
		log.Fatalf("error: %v", err)
	}

	if err = yaml.Unmarshal(data, &rel); err != nil {
		log.Fatalf("error: %v", err)
	}

	switch cmd {
	case "gen_source":
		rel.gen_source(dist, out)
	case "gen_plist":
		rel.gen_plist(dist, plbase, out)
	case "list":
		rel.list()
	}
}
