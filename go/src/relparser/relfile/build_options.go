package relfile

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
)

type BuildOptions struct {
	EnableAddressSanitizer           bool   `yaml:"enableAddressSanitizer"`
	EnableThreadSanitizer            bool   `yaml:"enableThreadSanitizer"`
	EnableUndefinedBehaviorSanitizer bool   `yaml:"enableUndefinedBehaviorSanitizer"`
	Toolchain                        string `yaml:"toolchain,omitempty"`
	Xcconfig                         string `yaml:"xcconfig,omitempty"`
}

func (opt *BuildOptions) PrintBuildOptions() {

	version := GetXcodeVersion()
	verComps := strings.Split(version, ".")
	major, _ := strconv.Atoi(verComps[0])

	val := reflect.ValueOf(opt).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		optName := strings.Split(tag.Get("yaml"), ",")[0]
		optValue := valueField.Interface()

		switch optValue.(type) {
		case string:
			//fmt.Printf("string: Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("yaml"))
			if optValue != "" {
				fmt.Printf("-%v\n\"%s\"\n", optName, optValue)
			}
		case bool:
			//fmt.Printf("bool: Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("yaml"))
			if optName == "enableUndefinedBehaviorSanitizer" && major <= 8 {
				continue
			}
			fmt.Printf("-%v\n%v\n",
				optName,
				func() string {
					if valueField.Interface().(bool) {
						return "YES"
					} else {
						return "NO"
					}
				}())
		}
	}
}

func GetXcodeVersion() string {
	var (
		err    error
		c      *exec.Cmd
		buffer bytes.Buffer
	)

	c = exec.Command("xcodebuild", "-version")
	c.Stdout = &buffer
	err = c.Run()
	if err != nil {
		return ""
	}

	scanner := bufio.NewScanner(bytes.NewReader(buffer.Bytes()))
	var version string
	for scanner.Scan() {
		version = strings.Split(scanner.Text(), " ")[1]
		break
	}
	if err := scanner.Err(); err != nil {
		logger.Fatalf("Failed to read Xcode verion")
	}

	return version
}
