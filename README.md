# Relax

relax is a comfortable release tool for your iOS Application. 

It's hard to understand stuff of `xcodebuild` and codesigning mechanism.
It takes care of much of the hassle of them, so you can focus on development.

Relax can be

- Manage your multi-deployment in a Relax configuration file(aka. Relfile)
    - Deploy your app quickly to testers, clients and end users in your professional project
    - Easy to set and revert a specific property of Info.plist and xcode build settings on each deployment.
    - Switch codesign mode(Automatic, Manual) automatically as your Relfile
- Validate your product's codesigin for a deployment
- Resign your xarchive and ipa file.

# Installation

```bash
$brew tap SCENEE/homebrew-formulae
$brew install relax
```
# Requirements

- macOS 10.11+
- Xcode8.0+

NOTE: Relax might be working on Xcode 7.3.1

Relax must have dependencies only on pre-installed command line tools in macOS and Xcode.

Because Relax aims to get rid of environment stuff when you manage a build server for your teaem.

# Getting Started

## Set up `Relfile`

Run this command and set up each configurations in `Relfile`.

```bash
$relax init
```

An example of Relfile is here.
And also check [this 'Refile' section](#relfile) and [the reference Refile](https://github.com/SCENEE/relax/blob/master/etc/Relfile) for detail.

```yaml
workspace: SampleApp

development: # Define a deployment type
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
  bundle_version:  "%R-%C"
  build_settings:
    OTHER_SWIFT_FLAGS: -DDEBUG
  info_plist: # You can change Info.plist settings for a deployment.
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

`xcodebuild` stdout is always written to a log file. 
If you would like to print it in your console, please use with '-v' option.

## Export an .ipa file

```bash
$relax export adhoc
```

You can specify a xcarchive file path after a release type like here.

```bash
$relax export adhoc /path/to/archive
```

Relax can export it on a different team and certificate from one signed xcarchive.

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
|%B| Build version number|
|%C| Build configuration|
|%h| Git abbreviated commit hash|
|%D| Git branch name|

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

# Known Issues

- Homebrew(0.9.x) failed to update Relax. Please use Homebrew(1.1.2+) with `brewe update`.
- Relax hasn't yet support Carthage. If you use it, Relax might not be working well. I'm glad for you to make a pull request to support it!

