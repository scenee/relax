#!/usr/bin/env bats

load test_helper

@test "relax export development" {
  run relax export development
  assert_success
}

@test "relax export adhoc <development-archive>" {
  run relax export adhoc "$(relax show development archive)"
  assert_success
  run relax validate "$(relax show adhoc ipa)"
  assert_success
  [[ "${output}" =~ "Relax Adhoc" ]]
  [[ "${output}" =~ "J3D7L9FHSS" ]]
}
