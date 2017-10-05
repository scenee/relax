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
	Sdk                 string                 `yaml:"sdk",omitempty"`
	Scheme              string                 `yaml:"scheme""`
	TeamID              string                 `yaml:"team_id"`
	ProvisioningProfile string                 `yaml:"provisioning_profile,omitempty"`
	Configuration       string                 `yaml:"configuration,omitempty"`
	Version             string                 `yaml:"version,omitempty"`
	BundleID            string                 `yaml:"bundle_identifier,omitempty"`
	BundleVersion       string                 `yaml:"bundle_version,omitempty"`
	InfoPlist           map[string]interface{} `yaml:"info_plist,omitempty"`
	BuildSettings       map[string]interface{} `yaml:"build_settings,omitempty"`
	BuildOptions        BuildOptions           `yaml:"build_options,omitempty"`
	ExportOptions       ExportOptions          `yaml:"export_options,omitempty"`
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
	for k, v := range d.InfoPlist {
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
	if d.BundleID == "" {
		bundleID = getBundleID(infoPlist)
	} else {
		bundleID = d.BundleID
	}

	encoder = plist.NewEncoder(out)
	encoder.Indent("\t")

	// fmt.Println("export options:", d.ExportOptions)
	opts = d.ExportOptions
	opts.TeamID = d.TeamID
	if d.ProvisioningProfile != "" {
		opts.SetProvisioningProfiles(bundleID, d.ProvisioningProfile)
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
	source += genSourceLine2(name, "provisioning_profile", d.ProvisioningProfile)
	source += genSourceLine2(name, "team_id", d.TeamID)
	source += genSourceLine2(name, "bundle_identifier", d.BundleID)
	source += genSourceLine2(name, "bundle_version", d.BundleVersion)
	source += genSourceLine2(name, "version", d.Version)

	// "--- Build settings\n"
	build_settings = strings.Join([]string{PREFIX, name, "build_settings"}, "_")

	source += fmt.Sprintf("%v=()\n", build_settings)

	for _, vars := range []map[string]interface{}{d.BuildSettings, d.InfoPlist} {
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
