#!/usr/bin/env bats

load test_helper

@test "relax archive development" {
  export COUNTRY=ja
  run relax archive development
  assert_success
  [[ ! "${output}" =~ "\[[ \*]*\]" ]]
  [[ ! "${output}" =~ "Time:" ]]
}

@test "relax archive development --progress" {
  export COUNTRY=ja
  run relax archive development --progress
  assert_success
  [[ ! "${output}" =~ "\[[ \*]*\]" ]]
  [[ "${output}" =~ "Time:" ]]
}

@test "relax archive development2" {
  export BUNDLE_SUFFIX=debug
  export VERSION="0.0.1"
  run  relax -v archive development2
  assert_success
  [[ ! "${output}" =~ "-derivedDataPath" ]]
}

@test "relax archive development2 --no-cache " {
  export BUNDLE_SUFFIX=debug
  export VERSION="0.0.1"
  run  relax -v archive development2 --no-cache
  assert_success
  [[ "${output}" =~ "-derivedDataPath" ]]
}

@test "relax archive adhoc" {
  run  relax archive adhoc
  assert_success
  [[ "$(grep '==> xcodebuild' <<< "$output")" =~ "iPhone Distribution" ]]
}


@test "relax archive: check workspace restoration" {
  run git diff --exit-code --quiet
  assert_success
}
