#!/usr/bin/env bats

load test_helper

@test "relax archive development" {
  export COUNTRY=ja
  run relax archive development
  assert_success
  [[ "${lines[${#lines[@]}-2]}" =~ "Time:" ]]
}

@test "relax archive development --progress" {
  export COUNTRY=ja
  run relax archive development --progress
  assert_success
  echo "$output" > bats.log
  [[ "${lines[${#lines[@]}-2]}" =~ "\[[ \*]*\]" ]]
  [[ "${lines[${#lines[@]}-2]}" =~ "Time:" ]]
}

@test "relax archive development2" {
  export BUNDLE_SUFFIX=debug
  export VERSION="0.0.1"
  run  relax archive development2
  assert_success
  [[ ! "${output}" =~ "Clean DerivedData" ]]
}

@test "relax archive development2 --no-cache " {
  export BUNDLE_SUFFIX=debug
  export VERSION="0.0.1"
  run  relax archive development2 --no-cache
  assert_success
  [[ "${output}" =~ "Clean DerivedData" ]]
}

@test "relax archive: check workspace restoration" {
  run git diff --exit-code --quiet
  assert_success
}
