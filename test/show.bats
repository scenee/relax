#!/usr/bin/env bats

load test_helper

@test "relax show development build" {
  run relax show development build
  assert_success
}

@test "relax show development ipa" {
  run relax show development ipa
  assert_success
}

@test "relax show development archive" {
  run relax show development archive
  assert_success
}
