#!/usr/bin/env bats

load test_helper

@test "relax build framework" {
  run relax build framework
  assert_success
}

@test "relax build staticlib" {
  run relax build staticlib --framework Sample
  assert_success
}

