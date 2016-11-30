#!/usr/bin/env bats

load test_helper

@test "relax export development --latest" {
  cd test/SampleApp
  run relax -v export development --latest
  assert_success
}

@test "relax export enterprise <development-archive>" {
  cd test/SampleApp
  run relax -v export enterprise $(relax show development archive)
  assert_success
  run relax validate $(relax show enterprise ipa)
  assert_success
  [[ "${output}" =~ "TGKEN7XA5C" ]]
}
