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

  run relax resign -k relax.keychain -i "com.scenee.SampleApp" -p "Relax AdHoc" -c "iPhone Distribution: Shin Yamamoto (J3D7L9FHSS)" "$(relax show development ipa)"
  assert_success

  run relax validate "$(find . -name "*-resigned.ipa")"
  assert_success
}

@test "relax resign failed if resigned ipa is existing" {
  run relax resign -i "com.scenee.SampleApp" -p "Relax Adhoc" -c "iPhone Distribution: Shin Yamamoto (J3D7L9FHSS)" "$(relax show development ipa)"
  assert_failure
}

@test "relax resign with invalid keychain" {

  relax keychain create mock.keychain -p mock
  relax keychain use mock.keychain -p mock

  rm -rf *-resigned.ipa
  run relax resign -k mock.keychain -i "com.scenee.SampleApp" -p "Relax AdHoc" -c "iPhone Distribution: Shin Yamamoto (J3D7L9FHSS)" "$(relax show development ipa)"
  assert_failure


  relax keychain use relax.keychain -p relax

  relax keychain delete mock.keychain
}
