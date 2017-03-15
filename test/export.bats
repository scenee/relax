#!/usr/bin/env bats

load test_helper

@test "relax export development" {
  run relax export development
  assert_success
  echo "$output" > bats.log
  [[ "${lines[${#lines[@]}-2]}" =~ "Time:" ]]
}

@test "relax export development --progress" {
  run relax export development --progress
  assert_success
  [[ "${lines[${#lines[@]}-2]}" =~ "\[[ \*]*\]" ]]
  [[ "${lines[${#lines[@]}-2]}" =~ "Time:" ]]
}

@test "relax export adhoc <development-archive>" {
  run relax export adhoc "$(relax show development archive)"
  run relax validate "$(relax show adhoc ipa)"

  [[ "${lines[9]}" =~ "Bundle Version: 1.0" ]] \
  && [[ "${lines[9]}" =~ "-debug)" ]] \
  && [[ "${lines[11]}" =~ "com.scenee.SampleApp" ]] \
  && [[ "${lines[31]}" =~ "Relax AdHoc" ]] \
  && [[ "${lines[33]}" =~ "J3D7L9FHSS" ]]
}

@test "relax export adhoc <development2-archive>" {
  run relax export adhoc "$(relax show development2 archive)"
  run relax validate "$(relax show adhoc ipa)"

  [[ "${lines[9]}" =~ "Bundle Version: 0.0.1" ]] \
  && [[ "${lines[9]}" =~ "-debug-internal)" ]] \
  && [[ "${lines[10]}" =~ "Bundle Name: Sample App" ]] \
  && [[  "${lines[11]}" =~ "Bundle Identifier: com.scenee.SampleApp.debug" ]] \
  && [[ "${lines[31]}" =~ "Relax AdHoc" ]] \
  && [[ "${lines[33]}" =~ "J3D7L9FHSS" ]]
}
