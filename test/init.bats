#!/usr/bin/env bats

load test_helper

@test "relax init" {
  mv Relfile Relfile.bak
  run relax init
  assert_success
  mv Relfile.bak Relfile
}

