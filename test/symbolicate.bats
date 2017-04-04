#!/usr/bin/env bats

load test_helper

# Test an ipa file detection
@test "relax symbolicate ipa" {
	touch temp
	run relax symbolicate temp "$(relax show development ipa)"
	assert_success
	rm temp
}

# Test a xcarchive file detection
@test "relax symbolicate archive" {
	touch temp
	run relax symbolicate temp "$(relax show development archive)"
	assert_success
	rm temp
}

# Test a unexpected file detection
@test "relax symbolicate temp(Unexpected file)" {
	touch temp
	run relax symbolicate temp temp
	assert_failure
	rm temp
}

