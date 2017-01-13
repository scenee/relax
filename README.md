[![Build Status](https://travis-ci.org/SCENEE/relax.svg?branch=support-travis-ci)](https://travis-ci.org/SCENEE/relax)

# Relax

Relax is a comfortable release tool for your iOS Application. 

It's hard to understand stuff of `xcodebuild` and codesigning mechanism.
It takes care of much of the hassle of them, so you can focus on development

You can

- **Manage** multi distributions(i.e. adhoc, enterprise and appstore) in a Relax configuration file(aka. Relfile)
- **Validate** your product(ipa/archive) has a correct codesign, entitlements and mobileprovision.
- **Resign** your product(ipa, archive) with a different bundle id and provisioning profile

Relax can

- Modify and revert specific properties in Info.plist and xcode build settings on each deployment
- Switch codesign modes(Automatic or Manual) automatically as your preferences
- Make it easy to manage your keychain with 'relax keychain'

# Installation

## Homebrew

```bash
$ brew tap SCENEE/homebrew-formulae
$ brew install relax
```

## Install script

```bash
$ curl -fsSL https://raw.githubusercontent.com/SCENEE/relax/master/install.sh | bash
```

# Requirements

Relax must only depend on pre-installed command line tools in macOS and Xcode.
Because it aims to get rid of any stuff of a host envornment to make it easy to manage a build server.
As a result, You can set it up on a macOS build machine quickly even if it's a virtual machine.

- macOS 10.11+
- Xcode8.0+

NOTE: Relax might be working on Xcode 7.3.1


# Known Issues

- Homebrew(0.9.x) failed to update Relax. Please use Homebrew(1.1.2+) with `brewe update`.
- Relax hasn't yet support Carthage. If you use it, Relax might not be working well. I'm glad for you to make a pull request to support it!
- `stty: stdin isn't a terminal` can be printed on a CI build server, but Relax is working well.

# Getting Started

## Quick start

```bash
$ cd /path/to/your/project

# Create a Relfile template
$ relax init 

# Build a xcarchive file
$ relax archive development 

# Print a path to a built archive
$ relax show development archive 

# Export an ipa file
$ relax export development 

# Print a path to a exported ipa file
$ relax show development ipa 

# Validate the ipa file
$ relax validate "$(relax show development ipa)"

# Upload the ipa file (It's necessary to add a token and secret in Relfile)
$ relax upload crashlytics "$(relax show development ipa)"

```

## Set up Relfile

Run this command and set up each configurations in `Relfile`.

```bash
$ relax init
```

An example of Relfile is here.
And also check [this 'Refile' section](#relfile) and [the reference Refile](https://github.com/SCENEE/relax/blob/master/sample/Relfile) for detail.

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
  export_options:
    method:  development

adhoc:
  scheme: SampleApp
  team_id: __MY_COMPANY_TEAM_ID__
  bundle_version:  "%h-%C" # See 'Bundle Version Format section'
  build_settings:
    OTHER_SWIFT_FLAGS: -DDEBUG
  info_plist: # You can change Info.plist settings for a deployment.
    CFBundleName: SmapleApp (DEBUG)
    CFBundleDevelopmentRegion: en
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
$ relax -v archive adhoc
```

`xcodebuild` stdout is always written to a log file. 
If you would like to print it in your console, please use with '-v' option.

## Export an .ipa file

```bash
$ relax export adhoc
```

You can specify a xcarchive file path after a release type like here.

```bash
$ relax export adhoc /path/to/archive
```

Relax can export it on a different team and certificate from one signed xcarchive.

## Check the mobileprovision, entitlements and version of an .ipa file

```bash
$ relax validate /path/to/SampleApp.ipa
```

## Upload an .ipa file to Crashlytics

```bash
$ relax upload crashlytics /path/to/SampleApp.ipa
```

## Resign an .ipa file with other provisioning profile and a certificate

```bash
$ relax resign -p "<my-provisioning-profile>" -c "iPhone Distribution: <Me>" /path/to/SampleApp.ipa
```
## Other commands

```bash
$ relax commands
```

# Relfile

Relfile is a relax configuration file. The reference file is [Here](https://github.com/SCENEE/relax/blob/master/samples/Relfile)

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
$ BUNDLE_SUFFIX=debug relax archive development2
```
or

```bash
$ export BUNDLE_SUFFIX=debug 
$ relax archive development2
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
