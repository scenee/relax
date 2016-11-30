#!/usr/bin/env bats

load test_helper

@test "relax archive development" {
  cd test/SampleApp
  run relax archive development
  assert_success
}

