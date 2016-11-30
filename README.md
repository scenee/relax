# Relax

relax is a stable and comfortable release tool for Xcode. It makes you freed from xcodebuild stuffs.

# Installation

```bash
$brew tap SCENEE/homebrew-formulae
$brew install relax
```
# Requirement

- Xcode8.0+ (relax might be working as your project settings on Xcode 7.3.1)

# Getting Started

## Set up `Relfile`

Run this command and set up each configurations in `Relfile`.

```bash
$relax init
```

An example of Relfile is here.

```yaml
workspace: SampleApp

development:
  scheme: SampleApp
  configuration: Debug
  build_settings:
    - OTHER_SWIFT_FLAGS: "-DMOCK"

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

log_formatter: xcpretty
```

## Build a xcarchive

```bash
$relax -v archive adhoc
```

## Export an .ipa file

```bash
$relax export adhoc --latest
```

## Check the provisioning profile and the entitlements of an .ipa file

```bash
$relax validate /path/to/SampleApp.ipa
```

## Upload your ipa to Crashlytics

```bash
$relax upload crashlytics /path/to/SampleApp.ipa
```

## Resign an .ipa file with a provisioning profile and a certificate

```bash
$relax resign -p "<my-provisioning-profile>" -c "iPhone Distribution: <Me>" /path/to/SampleApp.ipa
```

# Help

```bash
$relax help
```
