#!/usr/bin/env bats

load test_helper

@test "relax build framework" {
  cd test/SampleApp
  run relax build framework
  assert_success
}

@test "relax build staticlib" {
  cd test/SampleApp
  run relax build staticlib --framework Sample
  assert_success
}

