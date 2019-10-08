# Relax

Relax is a tiny release tool for iOS developers who don't want to be bothering with code signing and Xcode stuffs!!

You just configure `scheme` and `provisioning_profile` to build archive and IPA files.

Relax will save your time. It's hard to understand `xcodebuild` stuff, for example, code signing mechanism. Relax takes care of much of the hassle of them. so you can focus on development.

## Features

**Mutli distribution**

* Eanble to configure multi distribution for each code signing identity, info plist and build settings(i.e. ad-hoc builds for different clients).

**Less configuration**

* One configuration file instead of many xcconfig files or build configurations in your project
* Detect and check the related identity from a provisioning profile.
* Automatically generate ExportOptions.plist for each distribution

**Easy & Simple** 

* Easy to install into macOS
* Provide simple CLI to resolve issues of IPA distribution support and CI environment.
  * Validate and Resign IPA
  * Inspector of keychain and provisioning profiles

## Installation

### Homebrew

```bash
$ brew install scenee/formulae/relax
```

### Install script

```bash
$ curl -fsSL https://raw.githubusercontent.com/SCENEE/relax/master/install.sh | bash
```

NOTE: You don't need to take care of a host environment(i.e. ruby version and gem settings).

## Requirements

- Xcode9.4.1+


## Create a IPA file

```bash
$ relax dist /path/to/xcodeproj_or_xcworkspace --scheme scheme_name --profile profile_name
$ # OR
$ relax dist /path/to/xcodeproj_or_xcworkspace -s scheme_name -p profile_name
```

Or use Relfile.

```bash
$ relax init
$ relax dist adhoc
```

### Notes

* You need to create a provisioning profile for your identity(certificate) and install them to a build machine by yourself because Relax doesn't access to Apple Developer Center for security reasons.
* **`relax profile add` and `relax keychain add`** help you to install them and resolve permissions for your identities in your keychain. I highly recommend to use those commands. See [here](https://github.com/SCENEE/relax/blob/master/test/setup.sh#L22) and [here](https://github.com/SCENEE/relax/blob/master/test/setup.sh#L33).


## Relfile

Relfile is a configuration file for Relax. The declarative file will really make you easy to understand how to customize Info.plist and build settings for a distribution. See [here](https://github.com/SCENEE/relax/blob/master/sample/Relfile) for detail.

Here is an example.

```yaml
version: '2'

workspace: SampleApp
distributions:
  adhoc:
    # Required
    scheme: SampleApp
    provisioning_profile: 'Relax Adhoc'

    # Optional
    version: '1.0.1'
    configuration: Debug
    bundle_identifier: com.scenee.SampleApp.dev
    bundle_version: '$BUILD_NUMBER'  # You can use shell environment variables!
    info_plist:
      CFBundleName: 'SmapleApp(Debug)'
      UISupportedExternalAccessoryProtocols:
        - com.example.test-accessory
    build_settings:
      OTHER_SWIFT_FLAGS:
        - '-DMOCK'
    export_options:
      compileBitcode: false

  ent:
    # Required
    scheme: SampleApp
    provisioning_profile: 'Relax Enterprise'

    # Optional
    bundle_identifier: com.scenee.SampleApp
    bundle_version: '%b-%h-$c'  # See 'Bundle Version Format section'
    info_plist:
      UISupportedExternalAccessoryProtocols:
        - com.example.accessory

  framework:
    # Required
    scheme: Sample Framework

    # Optional
    configuration: Release

log_formatter: xcpretty # Optional
```

### Use Environment variables in Relfile

You can use Environment variables in Relfile. That's much useful in CI services. For example,

```yaml
development2:
  scheme: Sample App
  bundle_version: $BUILD_NUMBER
  ....
```

```bash
$ BUILD_NUMBER=11 relax archive development2
```

or

```bash
$ export BUILD_NUMBER=11
$ relax archive development2
```

But, you know, you can't use Xcode build setting variables (like PRODUCT_NAME etc.) in Relfile because they can be overridden by Relfile's definitions.

### Export Option Support

| Option                                   | Response status                                                                                 |
| :--------------------------------------- | :------------------------------------------------------------------                             |
| compileBitcode                           | OK                                                                                              |
| embedOnDemandResourcesAssetPacksInBundle | Not supported                                                                                   |
| iCloudContainerEnvironment               | Not supported                                                                                   |
| manifest                                 | Not supported                                                                                   |
| method                                   | Auto-assigned 'ad-hoc', 'app-store', 'development' or 'enterprise' from `provisioning_profile`. |
| onDemandResourcesAssetPacksBaseURL       | Not supported                                                                                   |
| teamID                                   | Auto-assigned from `provisioning_profile`.                                                      |
| provisioningProfiles                     | Auto-assigned from `provisioning_profile`.                                                      |
| signingCertificate                       | Auto-assigned 'iPhone Developer' or 'iPhone Distribution' from `provisioning_profile`.          |
| signingStyle                             | Auto-assigned 'automatic' or 'manual' determined from `provisioning_profile`.                   |
| thinning                                 | OK                                                                                              |
| uploadBitcode                            | OK                                                                                              |
| uploadSymbols                            | OK                                                                                              |

### Bundle Version Format

You can use specific format characters in a value of `bundle_version` field.
The characters and their meanings are as follows.

| Character | Meaning |
|:---------|:-------|
|%c| Build configuration|
|%h| Git abbreviated commit hash|
|%b| Git branch name|

## Advanced usages

### Create an archive file

```bash
$ relax archive dev
```

### Create a IPA file

```bash
$ relax export "/path/to/xcarchive"
$ # OR
$ relax export dev
```

### Validate a IPA file

Check a IPA file if it has a correct code signing and entitlements.

```bash
$ relax validate "$(relax show adhoc ipa)"
```

You can also validate an archive file.

### Resign a IPA file for an enterprise distribution

Resign a IPA file for a distribution by a provisioning profile. The related identity is selected automatically from active keychains.

```bash
$ relax resign -i "com.mycompany.SampleApp" -p "<enterprise-provisioning-profile>"  /path/to/ipa
```

### `keychain` commands

The `keychain` module commands make you free from keychain stuff and prevent a build break!
Actually this is an useful wrapper of `security` command.

Run here and see [this script](https://github.com/SCENEE/relax/blob/master/test/run.sh#L24) for detail.
```
$ relax help keychain
```

### `profile` commands

The `profile` module commands make it easy to find, use or remove provisioning profiles without Xcode Preferences.

Run here and see [this script](https://github.com/SCENEE/relax/blob/master/test/run.sh#L37) for detail.
```
$ relax help profile
```

### Symbolicate a crash log

```bash
$ relax symbolicate sampleapp.crash SampleApp.xcarchive
```

## Known Issues

- Homebrew(0.9.x) failed to update Relax. Please use Homebrew(1.1.2+) with `brewe update`.
- Relax hasn't yet support Carthage. If you use it, Relax might not be working well. I'm glad for you to make a pull request to support it!
- `stty: stdin isn't a terminal` can be printed on a CI build server, but Relax is working well.
