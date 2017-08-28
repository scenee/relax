#!/usr/bin/env bats

load test_helper

@test "relax build framework" {
  run relax build framework
  assert_success
  [[ ! "${lines[${#lines[@]}-4]}" =~ "\[[ \*]*\]" ]]
  [[ "${lines[${#lines[@]}-4]}" =~ "Time:" ]]
}

@test "relax build framework --progress" {
  run relax build framework --progress
  assert_success
  [[ "${lines[${#lines[@]}-4]}" =~ "\[[ \*]*\]" ]]
  [[ "${lines[${#lines[@]}-4]}" =~ "Time:" ]]

  [[ -f ./SampleFramework.framework/SampleFramework ]]
  file ./SampleFramework.framework/SampleFramework | grep arm.*
}

@test "relax build framework --with-simulator" {
  run relax build framework --with-simulator
  assert_success

  [[ -f ./SampleFramework.framework/SampleFramework ]]
  file ./SampleFramework.framework/SampleFramework | grep i386
  file ./SampleFramework.framework/SampleFramework | grep x86_64
  file ./SampleFramework.framework/SampleFramework | grep arm.*
}

@test "relax build framework --framework" {
  run relax build framework --framework
  assert_failure
}

@test "relax build staticlib --framework Sample" {
  run relax build staticlib --framework Sample
  assert_success

  [[ -d ./Sample.framework ]]
  [[ -f ./Sample.framework.zip ]]
  file ./Sample.framework/Sample | grep arm.*
}

@test "relax build staticlib --framework Sample --with-simulator" {
  run relax build staticlib --framework Sample --with-simulator
  assert_success

  [[ -d ./Sample.framework ]]
  [[ -f ./Sample.framework.zip ]]

  file ./Sample.framework/Sample | grep i386
  file ./Sample.framework/Sample | grep x86_64
  file ./Sample.framework/Sample | grep arm.*
}

@test "relax build: check workspace restoration" {
  run git diff --exit-code --quiet
  assert_success
}
