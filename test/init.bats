#!/usr/bin/env bats

load test_helper

@test "relax init" {
  mv Relfile Relfile.bak
  echo -e "1\n1" | relax init

  cat Relfile | grep -q -e "version: '2'"
  cat Relfile | grep -q -e "xcodeproj: SampleApp"
  cat Relfile | grep -q -e "scheme: Sample App"

  mv Relfile.bak Relfile
}
