package relfile

import (
	"fmt"
	"reflect"
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

	val := reflect.ValueOf(opt).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		switch valueField.Interface().(type) {
		case string:
			//fmt.Printf("string: Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("yaml"))
			if valueField.Interface() != "" {
				fmt.Printf("-%v\n\"%s\"\n", strings.Split(tag.Get("yaml"), ",")[0], valueField.Interface())
			}
		case bool:
			//fmt.Printf("bool: Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("yaml"))
			fmt.Printf("-%v\n%v\n",
				strings.Split(tag.Get("yaml"), ",")[0],
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
