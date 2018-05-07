[![Build Status](https://travis-ci.org/SCENEE/relax.svg?branch=support-travis-ci)](https://travis-ci.org/SCENEE/relax)


# Relax

Relax is a lazy(yes, just for laziness) release tool for iOS developers who don't want to be bothering with codesign and Xcode stuffs!!

You just configure `scheme` and `provisioning_profile` to build archive and ipa files.

You don't need to waste your time for codesigning problem especially on your CI workflow. Relax will save your time. It's hard to understand `xcodebuild` stuff, for example, codesigning mechanism. Relax takes care of much of the hassle of them. so you can focus on development.

Relax is not faster than xcodebuild, but more simple.

Relax is..

- **Easy**
    - You just configure `scheme`, `provisioning_profile`.
    - Relax generates an ExportOptions.plist from `provisioning_profile`.
    - Relax switches the codesigning mode **from Automatic to Manual automatically**.
- **Reproducible** 
    - Relax builds an app **only on 'Manual' signing mode** to prevent codesigning problems and make a build reproducible.
    - You can create an **isolated keychain db** with Relax. You could reproduce a build in your local machine as well if you use it. See `keychain` command.
- **Fine-tunable**
    - Relfile helps you to configure and understand multi distributions(i.e. adhoc, enterprise & appstore) having a few differences on code signing, Info.plist, Build Settings.
    - You no longer need to use many xcconfig files or build configurations in your project for the differences.


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

# Build an IPA in oneline

```bash
$ relax dist --scheme 'Sample App' --profile 'Relax AdHoc'
$ # OR
$ relax dist -s 'Sample App' -p 'Relax AdHoc'
```

# Build an IPA w/ Relfile

```bash
$ cd /path/to/your/project
$ relax init
$ # Update Relfile
$ relax dist adhoc
```


# Relfile

Relfile is a configuration file for Relax. The declarative file will really make you easy to understand what build settings you use to build a distribution and customize them. See [this Refile](https://github.com/SCENEE/relax/blob/master/sample/Relfile) for detail.

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

## Use Environment variables in Relfile

You can use Environment variables in Relfile. For example,

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

But you can't use Xcode build setting variables (i.e. PRODUCT_NAME etc.) in Relfile because they can be overridden by Relfile's definitions.

## Export Option Support

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

## Bundle Version Format

The characters and their meanings are as follows.

| Character | Meaning |
|:---------|:-------|
|%c| Build configuration|
|%h| Git abbreviated commit hash|
|%b| Git branch name|


# Advanced

## Create an archive file

```bash
$ relax archive dev
```

## Export an archive file to an IPA file

```bash
$ relax export "/path/to/xcarchive"
$ # OR
$ relax export dev
```

## Validate an IPA

Check a IPA file if it has a correct codesigning and entitlements.

```bash
$ relax validate "$(relax show adhoc ipa)"
```

## Resign an IPA for an enterprise distribution

Resign a IPA file for a distribution with a different bundle identifier, cetificate and provisioning profile

```bash
$ relax resign -m "com.mycompany.SampleApp" -p "<enterprise-provisioning-profile>" -c "iPhone Distribution: My Company"  "$(relax show dev ipa)"
$ relax validate SampleApp-resigned.ipa
```

## Symbolicate a crash log

```bash
$ relax symbolicate sampleapp.crash SampleApp.xcarchive
```

## `keychain` commands

The `keychain` module commands make you free from keychain stuff and prevent a codesign build break!
Actually this is an useful wrapper of `security` command.

Run here and see [this script](<F28>https://github.com/SCENEE/relax/blob/master/test/run.sh#L24) for detail.
```
$ relax help keychain
```

## `profile` commands

The `profile` module commands make it easy to find, use or remove provisioning profiles without Xcode Preferences.

Run here and see [this script](https://github.com/SCENEE/relax/blob/master/test/run.sh#L37) for detail.
```
$ relax help profile
```


# Known Issues

- Homebrew(0.9.x) failed to update Relax. Please use Homebrew(1.1.2+) with `brewe update`.
- Relax hasn't yet support Carthage. If you use it, Relax might not be working well. I'm glad for you to make a pull request to support it!
- `stty: stdin isn't a terminal` can be printed on a CI build server, but Relax is working well.
