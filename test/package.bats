#!/usr/bin/env bats

load test_helper

@test "relax package development artifact" {
  run relax export development
  run relax package development artifact
  assert_success
}

