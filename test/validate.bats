#!/usr/bin/env bats

load test_helper

@test "relax validate ipa" {
  cd test/SampleApp
  run relax validate "$(relax show development ipa)"
  assert_success
}

@test "relax validate archive" {
  cd test/SampleApp
  run relax validate "$(relax show development archive)"
  assert_success
}

