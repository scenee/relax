package relfile

import (
	"fmt"
	"github.com/DHowett/go-plist"
	"os"
	"strings"
)

//
// Distribution
//

type Distribution struct {
	Sdk                  string
	Scheme               string
	Team_id              string
	Provisioning_profile string
	Configuration        string
	Version              string
	Bundle_version       string
	Bundle_identifier    string
	Build_settings       map[string]interface{}
	Info_plist           map[string]interface{}
	Build_options        BuildOptions
	Export_options       ExportOptions
}

//
// Utils
//

const PREFIX string = "REL_CONFIG"

func genSourceline(key, value string) string {
	k := strings.Join([]string{PREFIX, key}, "_")
	return fmt.Sprintf("export %v=\"%v\"\n", k, value)
}

func genSourceLine2(name string, key string, value interface{}) string {
	k := strings.Join([]string{PREFIX, name, key}, "_")
	return fmt.Sprintf("export %v=\"%v\"\n", k, value)
}

func getBundleID(path string) string {
	var (
		err     error
		decoder *plist.Decoder
		f       *os.File
		data    map[string]interface{}
	)

	f, err = os.Open(path)
	if err != nil {
		logger.Fatalf("open error:", err)
	} else {
		defer f.Close()
		decoder = plist.NewDecoder(f)
	}

	err = decoder.Decode(&data)
	if err != nil {
		logger.Fatalf("decode error:", err)
	}

	props, ok := data["ApplicationProperties"].(map[string]interface{})
	if ok {
		return props["CFBundleIdentifier"].(string)
	} else {
		return data["CFBundleIdentifier"].(string)
	}
}

func (d Distribution) WriteInfoPlist(basePlistPath string, out *os.File) {
	var (
		err     error
		decoder *plist.Decoder
		f       *os.File
		data    map[string]interface{}
	)

	f, err = os.Open(basePlistPath)
	if err != nil {
		return
	}
	defer f.Close()

	decoder = plist.NewDecoder(f)

	err = decoder.Decode(&data)
	if err != nil {
		logger.Fatal(err)
	}

	//fmt.Println(data)
	//fmt.Println("--- Info.plist")

	/* Update Info.plist data */
	for k, v := range d.Info_plist {
		//fmt.Printf("---\t%v: %v\n", k, v)
		data[k] = v
	}

	encoder := plist.NewEncoder(out)
	encoder.Indent("\t")
	err = encoder.Encode(data)
	if err != nil {
		logger.Fatal(err)
	}
}

func (d Distribution) writeExportOptions(infoPlist string, out *os.File) {

	var (
		bundleID string
		encoder  *plist.Encoder
		opts     ExportOptions
	)

	// get bundle identifier
	if d.Bundle_identifier == "" {
		bundleID = getBundleID(infoPlist)
	} else {
		bundleID = d.Bundle_identifier
	}

	encoder = plist.NewEncoder(out)
	encoder.Indent("\t")

	// fmt.Println("export options:", d.Export_options)
	opts = d.Export_options
	opts.TeamID = d.Team_id
	if d.Provisioning_profile != "" {
		opts.SetProvisioningProfiles(bundleID, d.Provisioning_profile)
	}
	opts.Encode(encoder)
}

func (d Distribution) writeSource(name string, out *os.File) {

	var (
		err            error
		source         string
		build_settings string
	)

	source += genSourceLine2(name, "scheme", d.Scheme)
	source += genSourceLine2(name, "sdk", d.Sdk)
	source += genSourceLine2(name, "configuration", d.Configuration)
	source += genSourceLine2(name, "provisioning_profile", d.Provisioning_profile)
	source += genSourceLine2(name, "team_id", d.Team_id)
	source += genSourceLine2(name, "bundle_identifier", d.Bundle_identifier)
	source += genSourceLine2(name, "bundle_version", d.Bundle_version)
	source += genSourceLine2(name, "version", d.Version)

	// "--- Build settings\n"
	build_settings = strings.Join([]string{PREFIX, name, "build_settings"}, "_")

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

	_, err = out.WriteString(source)
	if err != nil {
		logger.Fatal(err)
	}
}
