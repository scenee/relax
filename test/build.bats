#!/usr/bin/env bats

load test_helper

@test "relax build framework" {
  run relax build framework
  assert_success
  [[ ! "${output}" =~ "\[[ \*]*\]" ]]
  [[ "${output}" =~ "Time:" ]]
  file ./SampleFramework.framework/SampleFramework | grep -q x86_64
  file ./SampleFramework.framework/SampleFramework | grep -q arm.*
}

@test "relax build framework --progress" {
  run relax build framework --progress
  assert_success
  [[ "${output}" =~ "\[[ \*]*\]" ]]
  [[ "${output}" =~ "Time:" ]]

  [[ -f ./SampleFramework.framework/SampleFramework ]]
  file ./SampleFramework.framework/SampleFramework | grep -q x86_64
  file ./SampleFramework.framework/SampleFramework | grep -q arm.*
}

@test "relax build framework --no-simulator" {
  run relax build framework --no-simulator
  assert_success

  [[ -f ./SampleFramework.framework/SampleFramework ]]
  file ./SampleFramework.framework/SampleFramework | grep -q -v x86_64
  file ./SampleFramework.framework/SampleFramework | grep -q arm.*
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
  file ./Sample.framework/Sample | grep -q x86_64
  file ./Sample.framework/Sample | grep -q arm.*
}

@test "relax build staticlib --framework Sample --no-simulator" {
  run relax build staticlib --framework Sample --no-simulator
  assert_success

  [[ -d ./Sample.framework ]]
  [[ -f ./Sample.framework.zip ]]

  file ./Sample.framework/Sample | grep -q -v x86_64
  file ./Sample.framework/Sample | grep -q arm.*
}

@test "relax build: check workspace restoration" {
  run git diff --exit-code --quiet
  assert_success
}
