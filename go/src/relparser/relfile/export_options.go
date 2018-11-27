package relfile

import (
	"github.com/DHowett/go-plist"
)

//
// ExportOptions
//
/**
  compileBitcode : Bool

          For non-App Store exports, should Xcode re-compile the app from bitcode? Defaults to YES.

  embedOnDemandResourcesAssetPacksInBundle : Bool

          For non-App Store exports, if the app uses On Demand Resources and this is YES, asset packs are embedded in the app bundle so that the app can be tested without a server to host asset packs. Defaults to YES unless onDemandResourcesAssetPacksBaseURL is specified.

  iCloudContainerEnvironment : String

          If the app is using CloudKit, this configures the "com.apple.developer.icloud-container-environment" entitlement. Available options vary depending on the type of provisioning profile used, but may include: Development and Production.

  installerSigningCertificate : String

          For manual signing only. Provide a certificate name, SHA-1 hash, or automatic selector to use for signing. Automatic selectors allow Xcode to pick the newest installed certificate of a particular type. The available automatic selectors are "Mac Installer Distribution" and "Developer ID Installer". Defaults to an automatic certificate selector matching the current distribution method.

  manifest : Dictionary

          For non-App Store exports, users can download your app over the web by opening your distribution manifest file in a web browser. To generate a distribution manifest, the value of this key should be a dictionary with three sub-keys: appURL, displayImageURL, fullSizeImageURL. The additional sub-key assetPackManifestURL is required when using on-demand resources.

  method : String

          Describes how Xcode should export the archive. Available options: app-store, package, ad-hoc, enterprise, development, developer-id, and mac-application. The list of options varies based on the type of archive. Defaults to development.

  onDemandResourcesAssetPacksBaseURL : String

          For non-App Store exports, if the app uses On Demand Resources and embedOnDemandResourcesAssetPacksInBundle isn't YES, this should be a base URL specifying where asset packs are going to be hosted. This configures the app to download asset packs from the specified URL.

  provisioningProfiles : Dictionary

          For manual signing only. Specify the provisioning profile to use for each executable in your app. Keys in this dictionary are the bundle identifiers of executables; values are the provisioning profile name or UUID to use.

  signingCertificate : String

          For manual signing only. Provide a certificate name, SHA-1 hash, or automatic selector to use for signing. Automatic selectors allow Xcode to pick the newest installed certificate of a particular type. The available automatic selectors are "Mac App Distribution", "iOS Distribution", "iOS Developer", "Developer ID Application", and "Mac Developer". Defaults to an automatic certificate selector matching the current distribution method.

  signingStyle : String

          The signing style to use when re-signing the app for distribution. Options are manual or automatic. Apps that were automatically signed when archived can be signed manually or automatically during distribution, and default to automatic. Apps that were manually signed when archived must be manually signed during distribtion, so the value of signingStyle is ignored.

  stripSwiftSymbols : Bool

          Should symbols be stripped from Swift libraries in your IPA? Defaults to YES.

  teamID : String

          The Developer Portal team to use for this export. Defaults to the team used to build the archive.

  thinning : String

          For non-App Store exports, should Xcode thin the package for one or more device variants? Available options: <none> (Xcode produces a non-thinned universal app), <thin-for-all-variants> (Xcode produces a universal app and all available thinned variants), or a model identifier for a specific device (e.g. "iPhone7,1"). Defaults to <none>.

  uploadBitcode : Bool

          For App Store exports, should the package include bitcode? Defaults to YES.

  uploadSymbols : Bool
*/

// ExportOptions struct
type ExportOptions struct {
	Method            string `yaml:"method"                      plist:"method"`
	Thinning          string `yaml:"thinning,omitempty"          plist:"thinning,omitempty"`
	CompileBitcode    bool   `yaml:"compileBitcode,omitempty"    plist:"compileBitcode"`
	UploadSymbols     bool   `yaml:"uploadSymbols,omitempty"     plist:"uploadSymbols"`
	UploadBitcode     bool   `yaml:"uploadBitcode,omitempty"     plist:"uploadBitcode"`
	StripSwiftSymbols bool   `yaml:"stripSwiftSymbols,omitempty" plist:"stripSwiftSymbols"`

	// Auto filled-in properties
	TeamID               string            `plist:"teamID,omitempty"`
	ProvisioningProfiles map[string]string `plist:"provisioningProfiles,omitempty"`
	SigningCertificate   string            `plist:"signingCertificate,omitempty"`
	SigningStyle         string            `plist:"signingStyle,omitempty"`
}

// UnmarshalYAML gets a ExportOptions object from YAML
func (opts *ExportOptions) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type typeAlias ExportOptions
	var t = &typeAlias{
		CompileBitcode:    true,
		StripSwiftSymbols: true,
		UploadSymbols:     true,
		UploadBitcode:     true,
	}

	err = unmarshal(t)
	if err != nil {
		return err
	}

	*opts = ExportOptions(*t)
	if opts.CompileBitcode == false {
		opts.UploadBitcode = false
	}
	return nil
}

// SetProvisioningProfiles set export options from a provisioning profile
func (opts *ExportOptions) SetProvisioningProfiles(provisioningProfile string, bundleID string) {
	if provisioningProfile == "" {
		logger.Fatalf("provisioningProfile is empty")
	}

	opts.SigningStyle = "manual"

	infos := FindProvisioningProfile("^"+provisioningProfile+"$", "")
	if len(infos) == 0 {
		logger.Fatalf("Not installed \"%s\"", provisioningProfile)
	}

	pp := infos[0].Pp

	opts.ProvisioningProfiles = map[string]string{bundleID: pp.Name}
	opts.TeamID = pp.TeamID()
	opts.Method = pp.ProvisioningType()
	opts.SigningCertificate = pp.CertificateType()
}

// Encode encodes plist
func (opts *ExportOptions) Encode(encoder *plist.Encoder) {
	err := encoder.Encode(opts)
	if err != nil {
		logger.Fatal(err)
	}
}
