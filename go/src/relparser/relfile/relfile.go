package relfile

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
)

//
// Utils
//

func LoadRelfile(path string) (*Relfile, error) {
	var (
		err    error
		c      *exec.Cmd
		f      *os.File
		buffer bytes.Buffer
		rel    Relfile
	)

	c = exec.Command("perl", "-pe", `s/\$(\{)?([a-zA-Z_]\w*)(?(1)\})/$ENV{$2}/g`)

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

	err = yaml.Unmarshal(buffer.Bytes(), &rel)
	if err != nil {
		return nil, err
	}

	return &rel, nil
}

//
// Relfile
//

type Relfile struct {
	Version       string                  `yaml:"version"`
	Xcodeproj     string                  `yaml:"xcodeproj"`
	Workspace     string                  `yaml:"workspace"`
	Uploader      map[string]interface{}  `yaml:"uploader"`
	LogFormatter  string                  `yaml:"log_formatter"`
	Distributions map[string]Distribution `yaml:"distributions"`
}

func (r Relfile) List() {
	for k, _ := range r.Distributions {
		fmt.Println(k)
	}
}

func (r Relfile) Check(dist string) {
	d := r.Distributions[dist]
	d.Check()
}

func (r Relfile) GenOptionsPlist(dist string, infoPlist, out string) {
	d := r.Distributions[dist]

	of := prepareFile(out)
	defer of.Close()

	d.writeExportOptions(infoPlist, of)
}

func (r Relfile) GenPlist(dist string, in string, out string) {
	d := r.Distributions[dist]

	of := prepareFile(out)
	defer of.Close()

	d.WriteInfoPlist(in, of)
}

func (r Relfile) GenSource(dist string) *os.File {
	d := r.Distributions[dist]

	of := prepareFile(os.TempDir() + "/relax/" + dist + "_source")
	defer of.Close()

	r.writeSource(of)
	d.writeSource(dist, of)

	return of
}

func (r Relfile) writeSource(out *os.File) {
	var (
		err    error
		source string
	)

	source += genSourceline("xcodeproj", r.Xcodeproj)
	source += genSourceline("workspace", r.Workspace)
	source += genSourceline("log_formatter", r.LogFormatter)

	for k, v := range r.Uploader {
		up := v.(map[interface{}]interface{})
		for name, value := range up {
			//fmt.Printf("---\t%v: %v\n", name, value)
			source += genSourceLine2("uploader_"+k, name.(string), value.(string))
		}
	}

	_, err = out.WriteString(source)
	if err != nil {
		logger.Fatal(err)
	}
}

func (r Relfile) PrintBuildOptions(dist string) {
	d := r.Distributions[dist]
	d.BuildOptions.PrintBuildOptions()
}
