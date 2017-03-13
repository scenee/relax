#!/usr/bin/env bats

load test_helper

@test "relax archive development" {
  export COUNTRY=ja
  run relax archive development
  assert_success
}

@test "relax archive development2" {
  export BUNDLE_SUFFIX=debug
  export VERSION="0.0.1"
  run  relax archive development2
  assert_success
  [[ ! "${output}" =~ "Clean DerivedData" ]]
}

@test "relax archive --no-cache development2" {
  export BUNDLE_SUFFIX=debug
  export VERSION="0.0.1"
  run  relax archive --no-cache development2
  assert_success
  [[ "${output}" =~ "Clean DerivedData" ]]
}
