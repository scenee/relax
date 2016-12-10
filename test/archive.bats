#!/usr/bin/env bats

load test_helper

@test "relax archive development" {
  export COUNTRY=ja
  run relax archive development
  assert_success
  echo "${output}" >> bats.log
}

@test "relax archive development2" {
  export BUNDLE_SUFFIX=debug
  run  relax archive development2
  assert_success
  echo "${output}" >> bats.log
}

