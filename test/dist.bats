#!/usr/bin/env bats

load test_helper

@test "relax dist  --scheme \"Sample App\" --profile \"Relax AdHoc\"" {
  run relax dist --scheme "Sample App" --profile "Relax AdHoc"
  assert_success
  [[ "${output}" =~ "xcarchive" ]]
  [[ "${output}" =~ "ipa" ]]
}

@test "relax dist adhoc" {
  run relax dist adhoc
  assert_success
  [[ "${output}" =~ "xcarchive" ]]
  [[ "${output}" =~ "ipa" ]]
}

