#!/usr/bin/env bats

load test_helper

@test "relax build framework" {
  run relax build framework
  assert_success
}

@test "relax build staticlib" {
  run relax build staticlib --framework Sample
  assert_success
  [[ -d ./Sample.framework ]] \
  && [[ -f ./Sample.framework.zip ]]
}

@test "relax build: check workspace restoration" {
  run git diff --exit-code --quiet
  assert_success
}
