#!/usr/bin/env bats

load test_helper

@test "relax init" {
  mv Relfile Relfile.bak
  run relax init
  assert_success

  [[ "${lines[5]}" =~ "version: '2'" ]] \
  && [ "$(sed -n 's/  *compileBitcode\: \(.*\)/\1/p' Relfile)" == 'false' ]

  mv Relfile.bak Relfile
}

