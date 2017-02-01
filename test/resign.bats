#!/usr/bin/env bats

load test_helper

@test "relax resign" {

  run relax ls-certs
  assert_success
  [[ $output =~ "iPhone Distribution: Shin Yamamoto (J3D7L9FHSS)" ]]

  run relax ls-profiles
  assert_success
  [[ $output =~ "Relax AdHoc" ]]

  rm -rf *-resigned.ipa
  run relax resign -i "com.scenee.SampleApp" -p "Relax AdHoc" -c "iPhone Distribution: Shin Yamamoto (J3D7L9FHSS)" "$(relax show development ipa)"
  echo "$output" >> bats.log
  assert_success

  run relax validate "Sample App-resigned.ipa"
  echo "$output" >> bats.log
  assert_success
}

@test "relax resign failed if resigned ipa is existing" {
  run relax resign -i "com.scenee.SampleApp" -p "Relax Adhoc" -c "iPhone Distribution: Shin Yamamoto (J3D7L9FHSS)" "$(relax show development ipa)"
  assert_failure
}

