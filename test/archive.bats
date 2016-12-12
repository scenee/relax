#!/usr/bin/env bats

load test_helper

@test "relax archive development" {
  pushd test/SampleApp
  export COUNTRY=ja
  run relax archive development
  assert_success
  popd
  echo "${output}" >> bats.log
}

@test "relax archive development2" {
  pushd test/SampleApp
  export BUNDLE_SUFFIX=debug
  run  relax archive development2
  assert_success
  popd
  echo "${output}" >> bats.log
}

