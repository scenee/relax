# Relax

relax is a comfortable release tool for your iOS Application. 

It's hard to understand stuff of `xcodebuild` and codesigning mechanism.
It takes care of much of the hassle of them, so you can focus on development.

Relax can be

- Manage your multi-deployment in a Relax configuration file(Relfile)
    - Deploy your app quickly to a tester, a client and an end user in your professional project
    - Easy to set and revert a specific property of Info.plist and xcode build settings on each release
    - Switch codesign mode(Automatic, Manual) automatically as your Relfile
- Validate your product's codesigin for a deployment
- Resign your xarchive and ipa file.

# Installation

```bash
$brew tap SCENEE/homebrew-formulae
$brew install relax
```
# Requirements

Relax must depend on only pre-install command line tools in macOS and ones of Xcode.

Because Relax aims to get rid of installation and environment stuff when you manage a build server for your teaem.

- macOS 10.11+
- Xcode8.0+

NOTE: Relax might be working on Xcode 7.3.1

# Getting Started

## Set up `Relfile`

Run this command and set up each configurations in `Relfile`.

```bash
$relax init
```

An example of Relfile is here.
Check [this Refile section](#Relfile) and [the reference Refile](https://github.com/SCENEE/relax/blob/master/etc/Relfile) for detail.

```yaml
workspace: SampleApp

development: # Define a Release Type
  scheme: SampleApp
  configuration: Debug
  build_settings:
    OTHER_SWIFT_FLAGS: 
      - "-DMOCK"
      - "-DDEBUG" 
  info_plist:
    UISupportedExternalAccessoryProtocols:
      - com.example.SampleApp.dev

adhoc:
  scheme: SampleApp
  configuration: Debug
  team_id: __MY_COMPANY_TEAM_ID__
  bundle_version_format:  "%R-%C"
  build_settings:
    OTHER_SWIFT_FLAGS: -DDEBUG
  info_plist:
    CFBundleName: SmapleApp (DEBUG)
    UISupportedExternalAccessoryProtocols:
      - com.example.SampleApp
      - com.example.SampleApp2
  export_options:
    method:  ad-hoc

enterprise:
  scheme: SampleApp
  team_id: __MY_ENTERPRISE_TEAM_ID__
  export_options:
    method:  enterprise

appstore:
  scheme: SampleApp
  sdk: iphoneos
  configuration: Release
  team_id: __MY_COMPANY_TEAM_ID__
  export_options:
    method:  appstore

framework:
  scheme: Sample Framework
  configuration: Release


crashlytics:
  token:  __MY_TOKEN__
  secret: __MY_SECRET__
  group:  __MY_GROUP__

log_formatter: xcpretty
```

## Build an archive for your product

```bash
$relax -v archive adhoc
```

If you would like to print xcodebuild stdout, please use with '-v' option.

## Export an .ipa file

```bash
$relax export adhoc
```

You can specify a xcarchive file path after a release type.

## Check the mobileprovision, entitlements and version of an .ipa file

```bash
$relax validate /path/to/SampleApp.ipa
```

## Upload an .ipa file to Crashlytics

```bash
$relax upload crashlytics /path/to/SampleApp.ipa
```

## Resign an .ipa file with other provisioning profile and a certificate

```bash
$relax resign -p "<my-provisioning-profile>" -c "iPhone Distribution: <Me>" /path/to/SampleApp.ipa
```
## Other commands

```bash
$relax commands
```

# Relfile

Relfile is a relax configuration file. The reference file is [Here](https://github.com/SCENEE/relax/blob/master/etc/Relfile)

## Use an Environment variable

You can use an Environment variable in Relfile.
An example is here.

```yaml
development2:
  scheme: Sample App
  bundle_identifier: com.scenee.SampleApp.$BUNDLE_SUFFIX
  ....
```

```bash
$BUNDLE_SUFFIX=debug relax archive development2
```
or

```bash
$export BUNDLE_SUFFIX=debug 
$relax archive development2
```
But you can't use Xcode build setting variables (PRODUCT_NAME etc.) in Relfile.
Because they can be overridden by Relfile's definitions.

## Bundle Version Format

The characters and their meanings are as follows.

| Character | Meaning |
|:---------|:-------|
|%V| Release version number|
|%v| Build version number|
|%C| Build configuration|
|%R| SCM commit ref|
|%B| SCM branch name|

## Export Option Support

| Option | Response status |
|:---------|:-------|
| method | OK |
| uploadSymbols | OK |
| compileBitcode | OK |
| team_id | OK(Automatic assigned) |
| thinning | OK |
| embedOnDemandResourcesAssetPacksInBundle | No |
| manifest | No |
| onDemandResourcesAssetPacksBaseURL | No |
| iCloudContainerEnvironment | No |


