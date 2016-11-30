#!/usr/bin/env bats

load test_helper

@test "relax resign" {
  pushd test/SampleApp
  run relax ls-certs
  assert_success
  echo "$output" >> ../../bats.log
  [[ $output =~ "iPhone Distribution: Shin Yamamoto (J3D7L9FHSS)" ]]
  run relax ls-profiles
  echo "$output" >> ../../bats.log
  assert_success
  [[ $output =~ "Scenee Wild Card Adhoc" ]]
  run relax resign -i "com.scenee.SampleApp" -p "Scenee Wild Card Adhoc" -c "iPhone Distribution: Shin Yamamoto (J3D7L9FHSS)" $(relax show enterprise ipa)
  echo "$output" >> ../../bats.log
  assert_success
  run relax validate SampleApp-resigned.ipa
  echo "$output" >> ../../bats.log
  assert_success
  popd
}

