#!/usr/bin/env bats

load test_helper

@test "relax build framework" {
  run relax build framework
  assert_success
  [[ ! "${lines[${#lines[@]}-3]}" =~ "\[[ \*]*\]" ]]
  [[ "${lines[${#lines[@]}-3]}" =~ "Time:" ]]
}

@test "relax build framework --progress" {
  run relax build framework
  assert_success
  [[ "${lines[${#lines[@]}-3]}" =~ "\[[ \*]*\]" ]]
  [[ "${lines[${#lines[@]}-3]}" =~ "Time:" ]]
}

@test "relax build framework --framework" {
  run relax build framework --framework
  assert_success

  [[ -f ./SampleFramework.framework/SampleFramework ]]
  file ./SampleFramework.framework/SampleFramework | grep i386
  file ./SampleFramework.framework/SampleFramework | grep x86_64
  file ./SampleFramework.framework/SampleFramework | grep arm.*
}

@test "relax build staticlib" {
  run relax build staticlib --framework-with-static Sample
  assert_success
  [[ -d ./Sample.framework ]]
  [[ -f ./Sample.framework.zip ]]
}

@test "relax build: check workspace restoration" {
  run git diff --exit-code --quiet
  assert_success
}
