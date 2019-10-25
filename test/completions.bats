#!/usr/bin/env bats

load test_helper

@test "relax completions resign" {
  run relax completions resign
  [[ ${output} = "-i --identifer -k --keychain -p -c" ]]
}

# relax resign -i |
@test "relax completions resign -i" {
  run relax completions resign -i
  [[ -z ${output} ]]
}

# relax resign -i foo |
@test "relax completions resign foo" {
  run relax completions resign foo
  [[ ${output} = "-i --identifer -k --keychain -p -c" ]]
}

# relax resign -i foo -|
@test "relax completions resign foo -" {
  run relax completions resign foo -
  [[ ${output} = "-i --identifer -k --keychain -p -c" ]]
}

@test "relax completiens foo" {
  run relax completions foo
  [[ -z ${output} ]]
}

@test "relax completions archive" {
  run relax completions archive
  [[ ${output} =~ "development" ]]
}

@test "relax completions show" {
  run relax completions show
  [[ ${output} =~ "development" ]]
}

@test "relax completions show development" {
  run relax completions show development
  [[ ${output} =~ "build" ]]
}

@test "relax completions show foo" {
  run relax completions show foo
  [[ -z ${output} ]]
}


@test "relax completions export" {
  run relax completions export
  [[ ${output} =~ "development" ]]
}

@test "relax completions build" {
  run relax completions build
  [[ ${output} =~ "development" ]]
}

@test "relax completions build foo" {
  run relax completions build foo
  [[ -z ${output} ]]
}

@test "relax completions build development" {
  run relax completions build development
  [[ ${output} = "-c --framework --progress" ]]
}

@test "relax completions package" {
  run relax completions package
  [[ ${output} =~ "development" ]]
}

@test "relax completions package development" {
  run relax completions package development
  [[ -z ${output} ]]
}

