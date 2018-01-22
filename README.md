[![Build Status](https://travis-ci.org/SCENEE/relax.svg?branch=support-travis-ci)](https://travis-ci.org/SCENEE/relax)


# Relax

Relax is a declarative release tool for iOS App distributions that encourages rapid distribution.

You don't need to write the same build script any more when you deliver your apps. Relax will save your time. You just write a declarative configuration file(aka. Relfile) for the distributions.

It's hard to understand `xcodebuild` stuff, for example, codesigning mechanism. Relax takes care of much of the hassle of them. so you can focus on development.

Relax can

- **Easy to release** multi distributions(i.e. adhoc, enterprise & appstore) with each build settings like code signing, bundle id and bundle version, etc.
- **Validate an ipa** to check if it has a correct codesign, a mobileprovision and entitlements
- **Resign an ipa** for a ditribution with a different bundle identifier, cetificate and provisioning profile
- **Set up a custom keychain** not to depend on a keychain settings in a build machine

In background, it works as below.

- Modify & Revert Info.plist properties and build settings in a xcodeproj file
- Switch codesign modes(Automatic or Manual) implicity if you specify a provisioning profile in Relfile


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

- macOS 10.11+
- Xcode8.0+

NOTE: Relax might be working on Xcode 7.3.1

Relax depends on only command line tools pre-installed in macOS and Xcode.
You don't need to take care of a host envornment(i.e. ruby version and gem settings).

As a result, You can set up iOS build environment on a new machine quickly
including keychain and provisionig profiles. 

For example, it's easy to manage a build server with a provisioning tool like Ansible.

# Getting Started

## Set up Relfile

```bash
$ cd /path/to/your/project
$ relax init
```

## Archive and Export

```bash
# Build a xcarchive file
# `xcodebuild` stdout is always written to a log file.
# If you would like to print logs in your console, please use with '-v' option.
$ relax -v archive adhoc

# Print a path to a built archive
$ relax show adhoc archive

# Export an ipa file
# Relax can export it on a different team and certificate from one signed xcarchive.
$ relax export adhoc

# Print a path to a exported ipa file
$ relax show adhoc ipa
```

## Validate the ipa

```bash
# Validate the ipa file
$ relax validate "$(relax show adhoc ipa)"
```

## Upload ipa

```bash
# Upload the ipa file (Need to add a token and secret in Relfile)
$ relax upload crashlytics "$(relax show adhoc ipa)"

```

## Resign an ipa for an enterprise distribution

```bash

$ relax resign -m "com.mycompany.SampleApp" -p "<enterprise-provisioning-profile>" -c "iPhone Distribution: My Company"  "$(relax show dev ipa)"
$ relax validate SampleApp-resigned.ipa
```

## Symbolicate a crash log

```bash
$ relax symbolicate sampleapp.crash SampleApp.xcarchive
```


# Relfile

Relfile is a configuration file for `relax`.

The declarative file will really make you easy to understand what build settings you use to build a distribution and customize them. See [this Refile](https://github.com/SCENEE/relax/blob/master/sample/Relfile) for detail.

Here is an example.

```yaml
version: '2'

workspace: SampleApp
log_formatter: xcpretty
uploader:
  crashlytics:
    token:  __MY_TOKEN__
    secret: __MY_SECRET__
distributions:  # Define a distribution
  adhoc:
    scheme: SampleApp
    team_id: ABCDEFGHIJ
    provisioning_profile: "Relax Adhoc"
    configuration: Debug
    version: 0.1.0
    bundle_version: "%b-%h-$c"  # See 'Bundle Version Format section'
    bundle_identifier: com.scenee.SampleApp.adhoc
    info_plist:
      CFBundleName: "SmapleApp(Debug)"
      UISupportedExternalAccessoryProtocols:
        - com.example.test-accessory
    build_settings:
      OTHER_SWIFT_FLAGS:
        - "-DMOCK"
        - "-DDEBUG" 
      OTHER_LINKER_FLAGS:
        - "$(inherited)"
        - "-ObjC"
    export_options:
      method:  ad-hoc
      compileBitcode: false

  appstore:
    scheme: SampleApp
    team_id: ABCDEFGHIJ
    provisioning_profile: "Relax AppStore"
    version: 1.0
    bundle_version: "$BUILD_NUMBER"  # You can use shell environment variables!
    bundle_identifier: com.scenee.SampleApp
    info_plist:
      UISupportedExternalAccessoryProtocols:
        - com.example.accessory
    export_options:
      method:  appstore

  framework:
    scheme: Sample Framework
    configuration: Release
```

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
|%c| Build configuration|
|%h| Git abbreviated commit hash|
|%b| Git branch name|

## Export Option Support

| Option                                   | Response status                                           |
| :--------------------------------------- | :-------------------------------------------------------- |
| compileBitcode                           | OK                                                        |
| embedOnDemandResourcesAssetPacksInBundle | Not yet                                                   |
| iCloudContainerEnvironment               | Not yet                                                   |
| manifest                                 | Not yet                                                   |
| method                                   | OK                                                        |
| onDemandResourcesAssetPacksBaseURL       | Not yet                                                   |
| teamID                                   | Applied `team_id` prop in Relfile                         |
| provisioningProfiles                     | Generated from `provisioning_profile` prop in Relfile     |
| signingCertificate                       | Auto-assigned 'iPhone Developer' or 'iPhone Distribution' |
| signingStyle                             | Auto-assigned 'automatic' or 'manual'                     |
| thinning                                 | OK                                                        |
| uploadBitcode                            | OK                                                        |
| uploadSymbols                            | OK                                                        |

## Uploader Support 

### Crashlytics

| Option                                   | Response status                                           |
| :--------------------------------------- | :-------------------------------------------------------- |
| token                                    | A API token                                               |
| secret                                   | A build secret                                            |
| group                                    | A group name                                              |


# CI Utilities

## relax keychain

The `keychain` module commands make you free from keychain stuff and prevent a codesign build break!
Actually this is an usefull wrapper of `security` command.

## relax profile

The `profile` module commands make it easy to find, use or remove provisioning profiles without Xcode Preferences.


# What's different from GYM?

- Multi disbribution management
- Focus on use cases on a macOS build machine
- Less dependancies 
- Support ipa resign and validation
- Support keychain managment


# Known Issues

- Homebrew(0.9.x) failed to update Relax. Please use Homebrew(1.1.2+) with `brewe update`.
- Relax hasn't yet support Carthage. If you use it, Relax might not be working well. I'm glad for you to make a pull request to support it!
- `stty: stdin isn't a terminal` can be printed on a CI build server, but Relax is working well.
