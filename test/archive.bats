#!/usr/bin/env bats

load test_helper

@test "relax archive development" {
  export COUNTRY=ja
  run relax archive development
  echo "${output}" >> bats.log
  assert_success
}

@test "relax archive development2" {
  export BUNDLE_SUFFIX=debug
  run  relax archive development2
  echo "${output}" >> bats.log
  assert_success
}

