#!/usr/bin/env bats

load test_helper

@test "relax archive development" {
  cd test/SampleApp
  run relax archive development
  assert_success
}

@test "relax archive development2" {
  cd test/SampleApp
  export BUNDLE_SUFFIX=debug
  run  relax archive development2
  assert_success
}

