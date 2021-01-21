#!/usr/bin/env bats

load test_helper

@test "relax package development artifact" {
  run relax export development
  run relax package development artifact
  assert_success
  [[ $(find artifact -name "Sample*.ipa") ]]; \
	  [[ -f artifact/Sample\ App.xcarchive.zip ]]
}

@test "relax package development2 artifact" {
  run relax package development2 artifact
  assert_failure
}

