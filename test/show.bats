#!/usr/bin/env bats

load test_helper

@test "relax show development build" {
  cd test/SampleApp
  run relax show development build
  assert_success
}

@test "relax show development ipa" {
  cd test/SampleApp
  run relax show development ipa
  assert_success
}

@test "relax show development archive" {
  cd test/SampleApp
  run relax show development archive
  assert_success
}
