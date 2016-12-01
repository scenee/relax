# Relax

relax is a comfortable release tool for your iOS Application. 

It's hard to understand stuff of `xcodebuild` and codesigning mechanism.
It takes care of much of the hassle of them, so you can focus on development.

- Don't need to write a huge build script
- Easy to check/validate your built product for codesigining
- Reduce build time just to change a code signature


# Installation

```bash
$brew tap SCENEE/homebrew-formulae
$brew install relax
```
# Requirement

- Xcode8.0+

NOTE: Relax might be working on Xcode 7.3.1.

# Getting Started

## Set up `Relfile`

Run this command and set up each configurations in `Relfile`.

```bash
$relax init
```

An example of Relfile is here.

```yaml
workspace: SampleApp

development: # Release Type
  scheme: SampleApp
  configuration: Debug
  build_settings:
    - OTHER_SWIFT_FLAGS: -DMOCK

adhoc:
  scheme: SampleApp
  team_id: COMPANY_TEAM_ID
  bundle_version_format:  %R-%C
  export_options:
    method:  ad-hoc

enterprise:
  scheme: SampleApp
  team_id: ENTERPRISE_TEAM_ID
  export_options:
    method:  enterprise

appstore:
  scheme: SampleApp
  sdk: iphoneos
  configuration: Release
  team_id: COMPANY_TEAM_ID 
  export_options:
    method:  appstore

framework:
  scheme: Sample Framework
  configuration: Release

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

## Bundle Identifier Format

The characters and their meanings are as follows.

| Character | Meaning |
|:---------|:-------|
|%V| Version number|
|%v| Bundle number|
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


