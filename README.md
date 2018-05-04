[![Build Status](https://travis-ci.org/SCENEE/relax.svg?branch=support-travis-ci)](https://travis-ci.org/SCENEE/relax)


# Relax

Relax is a **declarative** and **CI oriented** release tool for iOS App distributions that encourages rapid distribution **under your control**.

You don't need to write the same build script any more when you deliver your apps. Relax will save your time. You just write a declarative configuration file(aka. Relfile) for the distributions.

It's hard to understand `xcodebuild` stuff, for example, codesigning mechanism. Relax takes care of much of the hassle of them. so you can focus on development.

Relax can

- **Declarative**
    - Ralax uses the declarative file(Relfile) to configure your each build. It makes you easy to read and understand your build.
- **CI oriented**
    - Relax is easy to integrate CI tools(Jenkins, Bamboo, etc) and services(TravisCI, CircleCI, etc).
- **Under control(Reproducible)**: 
    - Manual Signing in free: Relax runs a build only on manual signing mode, but don't worry, it automatically switches a codesigning mode from Automatic to Manual only in build time. Therefore you can use Automatic Signing in your development.
    - Custom Keychain store: Relax makes and switch a temp keychain. You won't affected from a keychain settings in a build machine and reproduce your build on your machine.
    - Validation: Check a IPA file if it has a correct codesigning and entitlements.
    - Resign: Resign a IPA file for a distribution with a different bundle identifier, cetificate and provisioning profile
- **Safe**
    - You don't need to share your AppleID and **password**.
- **Easy to customize**: Relax is easy to customize Info.plist and Build Settings of your Xcode project only in build time.
    - Multi distributions(i.e. adhoc, enterprise & appstore): You can setup each configuration like code signing, Info.plist, Build Settings, etc.

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
You don't need to take care of a host environment(i.e. ruby version and gem settings).

As a result, You can set up iOS build environment on a new machine quickly
including keychain and provisioning profiles. 

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

# Relfile

Relfile is a configuration file for `relax`.

The declarative file will really make you easy to understand what build settings you use to build a distribution and customize them. See [this Refile](https://github.com/SCENEE/relax/blob/master/sample/Relfile) for detail.

Here is an example.

```yaml
version: '2'

workspace: SampleApp
uploader:
  crashlytics:
    token:  __MY_TOKEN__
    secret: __MY_SECRET__
distributions:  # Define a distribution
  dev:
    # Required
    scheme: SampleApp
    team_id: ABCDEFGHIJ
    provisioning_profile: 'Relax Adhoc'

    # Optional
    configuration: Debug

    version: '1.0.1'

    bundle_identifier: com.scenee.SampleApp.dev
    bundle_version: '%b-%h-$c'  # See 'Bundle Version Format section'

    info_plist:
      CFBundleName: 'SmapleApp(Debug)'
      UISupportedExternalAccessoryProtocols:
        - com.example.test-accessory
    build_settings:
      OTHER_SWIFT_FLAGS:
        - '-DMOCK'

    export_options:
      compileBitcode: false

  prod:
    # Required
    scheme: SampleApp
    team_id: ABCDEFGHIJ
    provisioning_profile: 'Relax Enterprise'

    # Optional
    bundle_identifier: com.scenee.SampleApp
    bundle_version: '$BUILD_NUMBER'  # You can use shell environment variables!

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


# Advanced

## Validate the ipa

```bash
# Validate the ipa file
$ relax validate "$(relax show adhoc ipa)"
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

## Upload ipa

```bash
# Upload the ipa file (Need to add a token and secret in Relfile)
$ relax upload crashlytics "$(relax show adhoc ipa)"
```

## `keychain` commands

The `keychain` module commands make you free from keychain stuff and prevent a codesign build break!
Actually this is an useful wrapper of `security` command.

Run here and see [this script](https://github.com/SCENEE/relax/blob/master/test/run.sh) for detail.
```
$ relax help keychain
```

## `profile` commands

The `profile` module commands make it easy to find, use or remove provisioning profiles without Xcode Preferences.

Run here and see [this script](https://github.com/SCENEE/relax/blob/master/test/run.sh) for detail.
```
$ relax help profile
```

# Known Issues

- Homebrew(0.9.x) failed to update Relax. Please use Homebrew(1.1.2+) with `brewe update`.
- Relax hasn't yet support Carthage. If you use it, Relax might not be working well. I'm glad for you to make a pull request to support it!
- `stty: stdin isn't a terminal` can be printed on a CI build server, but Relax is working well.
