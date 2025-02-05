package relfile

import (
	"certutil"
	"crypto/sha1"
	"crypto/x509"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/DHowett/go-plist"
)

/*
Distribution struct
*/
type Distribution struct {
	// Required
	Scheme              string `yaml:"scheme"`
	ProvisioningProfile string `yaml:"provisioning_profile"`

	// Optional
	Sdk           string                 `yaml:"sdk,omitempty"`
	Configuration string                 `yaml:"configuration,omitempty"`
	Version       string                 `yaml:"version,omitempty"`
	BundleID      string                 `yaml:"bundle_identifier,omitempty"`
	BundleVersion string                 `yaml:"bundle_version,omitempty"`
	InfoPlist     map[string]interface{} `yaml:"info_plist,omitempty"`
	BuildSettings map[string]interface{} `yaml:"build_settings,omitempty"`
	BuildOptions  BuildOptions           `yaml:"build_options,omitempty"`
	ExportOptions ExportOptions          `yaml:"export_options,omitempty"`
}

// UnmarshalYAML Distribution
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

// Check Distribution
func (d *Distribution) Check() {
	// Check the ProvisioningProfile existence
	infos := FindProvisioningProfile("^"+d.ProvisioningProfile+"$", "")
	if len(infos) == 0 {
		logger.Fatalf("\"%s\" not found.", d.ProvisioningProfile)
	}

	// Check the related Certificate existence
	pp := infos[0].Pp

	ok := false
	for _, data := range pp.DeveloperCertificates {
		var (
			cert *x509.Certificate
			err  error
		)
		// fmt.Printf("%q\n", data)
		if cert, err = x509.ParseCertificate(data); err != nil {
			fmt.Println("error:", err)
			return
		}
		issuerCN := cert.Issuer.CommonName
		if data, err = exec.Command("/usr/bin/security", "find-certificate", "-c", issuerCN).Output(); err != nil {
			logger.Printf("\"%s\" certificate doesn't installed in Keychain.", issuerCN)
			certutil.InstallCertificate(issuerCN, "")
		}

		sha1Fingerprint := sha1.Sum(cert.Raw)

		if data, err = exec.Command("/usr/bin/security", "find-identity", "-v", "-p", "codesigning").Output(); err != nil {
			logger.Fatalln(err)
		}
		re := regexp.MustCompile(fmt.Sprintf("%X", sha1Fingerprint))
		matches := re.FindStringSubmatch(string(data[:]))
		if len(matches) > 0 {
			ok = true
		}
	}

	if !ok {
		logger.Fatalf("No valid set of certificate and identity found for \"%s\". Please check 'My Certificates' in Keychain Access.app.", d.ProvisioningProfile)
	}

	if d.BundleID != "" {
		re := regexp.MustCompile(strings.Replace(pp.AppID(), "*", ".*", -1))
		matches := re.FindStringSubmatch(d.BundleID)
		if len(matches) == 0 {
			logger.Fatalf("\"%s\" doesn't match AppID of \"%s\".", d.BundleID, d.ProvisioningProfile)
		}
	}
}

// WriteInfoPlist Output Info.plist
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
					new := mergeMap(old, new)
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
		logger.Fatalf("`provisioning_profile` field is required in Relfile.")
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
		err           error
		source        string
		buildSettings string
	)

	source += fmt.Sprintf("export %v=\"%v\"\n", "_SCHEME", d.Scheme)

	source += fmt.Sprintf("export %v=\"%v\"\n", "_SDK", d.Sdk)
	source += fmt.Sprintf("export %v=\"%v\"\n", "_CONFIGURATION", d.Configuration)

	source += fmt.Sprintf("export %v=\"%v\"\n", "_VERSION", d.Version)
	source += fmt.Sprintf("export %v=\"%v\"\n", "_BUNDLE_ID", d.BundleID)
	source += fmt.Sprintf("export %v=\"%v\"\n", "_BUNDLE_VERSION", d.BundleVersion)

	infos := FindProvisioningProfile("^"+d.ProvisioningProfile+"$", "")
	if len(infos) > 0 {
		pp := infos[0].Pp
		source += fmt.Sprintf("export %v=\"%v\"\n", "_PROVISIONING_PROFILE", pp.Name)
		source += fmt.Sprintf("export %v=\"%v\"\n", "_TEAM_ID", pp.TeamID())
		source += fmt.Sprintf("export %v=\"%v\"\n", "_IDENTITY", pp.CertificateType())
	}

	// FIXME: Improve here
	// Build Settings
	buildSettings = strings.Join([]string{PREFIX, name, "build_settings"}, "_")
	source += fmt.Sprintf("%v=()\n", buildSettings)
	for k, v := range d.BuildSettings {
		switch v := v.(type) {
		case []interface{}:
			var ss []string
			for _, s := range v {
				ss = append(ss, fmt.Sprintf("%v", s))
			}
			source += fmt.Sprintf("%v+=(%v='%v')\n", buildSettings, k, strings.Join(ss, "{}"))
		default:
			source += fmt.Sprintf("%v+=(%v='%v')\n", buildSettings, k, v)
		}
	}
	source += fmt.Sprintf("export %v\n", buildSettings)

	_, err = out.WriteString(source)
	if err != nil {
		logger.Fatal(err)
	}
}
