#!/usr/bin/env bats

load test_helper

@test "relax apply development" {
  run relax apply development
  assert_success

  grep -q "DEVELOPMENT_TEAM = J3D7L9FHSS;" SampleApp.xcodeproj/project.pbxproj
  grep -q "ProvisioningStyle = Manual;" SampleApp.xcodeproj/project.pbxproj

  git checkout .
}
