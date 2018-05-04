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
	// Required
	Scheme              string `yaml:"scheme""`
	TeamID              string `yaml:"team_id"`
	ProvisioningProfile string `yaml:"provisioning_profile"`

	// Optional
	Sdk           string                 `yaml:"sdk",omitempty"`
	Configuration string                 `yaml:"configuration,omitempty"`
	Version       string                 `yaml:"version,omitempty"`
	BundleID      string                 `yaml:"bundle_identifier,omitempty"`
	BundleVersion string                 `yaml:"bundle_version,omitempty"`
	InfoPlist     map[string]interface{} `yaml:"info_plist,omitempty"`
	BuildSettings map[string]interface{} `yaml:"build_settings,omitempty"`
	BuildOptions  BuildOptions           `yaml:"build_options,omitempty"`
	ExportOptions ExportOptions          `yaml:"export_options,omitempty"`
}

//
// Utils
//

func (d *Distribution) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type typeAlias Distribution
	var t = &typeAlias{
		Sdk: "iphoneos",
	}

	err = unmarshal(t)
	if err != nil {
		return err
	}

	*d = Distribution(*t)
	return nil
}

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
		logger.Fatalf("open error: %v", err)
	} else {
		defer f.Close()
		decoder = plist.NewDecoder(f)
	}

	err = decoder.Decode(&data)
	if err != nil {
		logger.Fatalf("decode error: %v", err)
	}

	props, ok := data["ApplicationProperties"].(map[string]interface{})
	if ok {
		return props["CFBundleIdentifier"].(string)
	} else {
		return data["CFBundleIdentifier"].(string)
	}
}

func cleanupInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanupInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanupInterfaceMap(v)
	default:
		return v
	}
}

func (d Distribution) WriteInfoPlist(basePlistPath string, out *os.File) {
	//fmt.Println("--- WriteInfoPlist")
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
		data[k] = cleanupMapValue(v)
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

	if d.ProvisioningProfile == "" {
		logger.Fatalf("`provisioning_profile` field is required in Relfile")
	}

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
	opts.SetProvisioningProfiles(bundleID, d.TeamID, d.ProvisioningProfile)
	opts.Encode(encoder)
}

func (d Distribution) writeSource(name string, out *os.File) {
	var (
		err            error
		source         string
		build_settings string
	)

	source += fmt.Sprintf("export %v=\"%v\"\n", "_SCHEME", d.Scheme)
	source += fmt.Sprintf("export %v=\"%v\"\n", "_SDK", d.Sdk)
	source += fmt.Sprintf("export %v=\"%v\"\n", "_CONFIGURATION", d.Configuration)
	source += fmt.Sprintf("export %v=\"%v\"\n", "_PROVISIONING_PROFILE", d.ProvisioningProfile)
	source += fmt.Sprintf("export %v=\"%v\"\n", "_TEAM_ID", d.TeamID)

	infos := FindProvisioningProfile(d.ProvisioningProfile, d.TeamID)
	pp := infos[0].Pp
	source += fmt.Sprintf("export %v=\"%v\"\n", "_IDENTITY", pp.CertificateType())

	source += fmt.Sprintf("export %v=\"%v\"\n", "_BUNDLE_ID", d.BundleID)
	source += fmt.Sprintf("export %v=\"%v\"\n", "_BUNDLE_VERSION", d.BundleVersion)
	source += fmt.Sprintf("export %v=\"%v\"\n", "_VERSION", d.Version)

	// fmt/Println("--- Build settings\n")
	build_settings = strings.Join([]string{PREFIX, name, "build_settings"}, "_")

	source += fmt.Sprintf("%v=()\n", build_settings)

	// FIXME: Improve here
	for k, v := range d.BuildSettings {
		switch v := v.(type) {
		case []interface{}:
			var ss []string
			for _, s := range v {
				ss = append(ss, fmt.Sprintf("%v", s))
			}
			source += fmt.Sprintf("%v+=(%v='%v')\n", build_settings, k, strings.Join(ss, "{}"))
		default:
			source += fmt.Sprintf("%v+=(%v='%v')\n", build_settings, k, v)
		}
	}
	source += fmt.Sprintf("export %v\n", build_settings)

	_, err = out.WriteString(source)
	if err != nil {
		logger.Fatal(err)
	}
}
