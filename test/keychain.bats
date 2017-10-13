#!/usr/bin/env bats

load test_helper

@test "relax keychain info" {
  run relax keychain info certificates/RelaxCertificates.p12 -P "$CERTS_PASS"
  assert_success
}

