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

func (d *Distribution) IsMatchBundleID(infoPlist string) bool {
	// get bundle identifier
	var bundleID string
	if d.BundleID == "" {
		bundleID = getBundleID(infoPlist)
	} else {
		bundleID = d.BundleID
	}

	infos := FindProvisioningProfile(d.ProvisioningProfile, "")
	pp := infos[0].Pp

	// TODO:
	fmt.Printf("%v == %v", pp.AppID(), bundleID)
	return true
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
		new := cleanupMapValue(v)
		if old, ok := data[k]; ok {
			switch old := old.(type) {
			case map[string]interface{}:
				if new, ok := new.(map[string]interface{}); ok {
					new := updatedMap(old, new)
					data[k] = new
					continue
				}
			}
		}
		data[k] = new
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
		encoder *plist.Encoder
		opts    ExportOptions
	)

	if d.ProvisioningProfile == "" {
		logger.Fatalf("`provisioning_profile` field is required in Relfile")
	}

	encoder = plist.NewEncoder(out)
	encoder.Indent("\t")

	// fmt.Println("export options:", d.ExportOptions)

	// Info.plist is one in xcarchive file. So we have to use the BundleID, not in Relfile.
	bundleID := getBundleID(infoPlist)

	opts = d.ExportOptions
	opts.SetProvisioningProfiles(d.ProvisioningProfile, bundleID)
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

	source += fmt.Sprintf("export %v=\"%v\"\n", "_VERSION", d.Version)
	source += fmt.Sprintf("export %v=\"%v\"\n", "_BUNDLE_ID", d.BundleID)
	source += fmt.Sprintf("export %v=\"%v\"\n", "_BUNDLE_VERSION", d.BundleVersion)

	infos := FindProvisioningProfile(d.ProvisioningProfile, "")

	if len(infos) == 0 {
		logger.Fatalf("Not installed \"%s\"", d.ProvisioningProfile)
	}

	pp := infos[0].Pp
	source += fmt.Sprintf("export %v=\"%v\"\n", "_TEAM_ID", pp.TeamID())
	source += fmt.Sprintf("export %v=\"%v\"\n", "_IDENTITY", pp.CertificateType())
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
