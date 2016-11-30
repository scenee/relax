#!/usr/bin/env bats

load test_helper

@test "relax package development artifact" {
  cd test/SampleApp
  run relax package development artifact
  assert_success
}

