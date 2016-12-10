#!/usr/bin/env bats

load test_helper

@test "relax validate ipa" {
  run relax validate "$(relax show development ipa)"
  assert_success
}

@test "relax validate archive" {
  run relax validate "$(relax show development archive)"
  assert_success
}

