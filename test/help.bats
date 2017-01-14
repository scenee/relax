#!/usr/bin/env bats

load test_helper

@test "relax help" {
	run relax help
	assert_success
}

@test "relax help profile" {
	run relax help profile
	assert_success
}

@test "relax help keychain" {
	run relax help keychain
	assert_success
}

@test "relax help init" {
	run relax help init
	assert_success
}

@test "relax help archive" {
	run relax help archive
	assert_success
}

@test "relax help build" {
	run relax help build
	assert_success
}

@test "relax help export" {
	run relax help export
	assert_success
}

@test "relax help package" {
	run relax help package
	assert_success
}

@test "relax help show" {
	run relax help show
	assert_success
}

@test "relax help upload" {
	run relax help upload
	assert_success
}

